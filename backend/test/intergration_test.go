package test

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/app"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestAPI(t *testing.T) {
	_, connStr := createDB(t)

	if err := runMigrations(connStr); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	startAPI(connStr)
	waitForServer(t)

	t.Run("healthcheck", func(t *testing.T) {
		testHealthcheck(t, connStr)
	})

	t.Run("leagues_e2e", func(t *testing.T) {
		testLeaguesE2E(t, connStr)
	})
}

func waitForServer(t *testing.T) {
	t.Helper()

	for i := 0; i < 20; i++ {
		resp, err := http.Get("http://localhost:3000/v1/healthcheck")
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatal("server did not start")
}

func createDB(t *testing.T) (*postgres.PostgresContainer, string) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("users"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		testcontainers.TerminateContainer(container)
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return container, connStr
}

func runMigrations(dsn string) error {
	m, err := migrate.New(
		"file://../db/migrations",
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func startAPI(dsn string) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var cfg app.Config
	cfg.Port = 3000
	cfg.Env = "development"

	cfg.DB.DSN = dsn
	cfg.DB.MaxOpenConns = 25
	cfg.DB.MaxIdleConns = 25
	cfg.DB.MaxIdleTime = 15 * time.Minute

	cfg.Limiter.RPS = 20
	cfg.Limiter.Burst = 50
	cfg.Limiter.Enabled = true

	application, err := app.New(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := application.Server(); err != nil {
			log.Fatal(err)
		}
	}()
}

func testHealthcheck(t *testing.T, dsn string) {
	t.Helper()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	var cfg app.Config
	cfg.DB.DSN = dsn
	cfg.DB.MaxOpenConns = 25
	cfg.DB.MaxIdleConns = 25
	cfg.DB.MaxIdleTime = 15 * time.Minute

	application, err := app.New(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	rr := httptest.NewRecorder()

	application.Routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "status") {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func testLeaguesE2E(t *testing.T, dsn string) {
	t.Helper()

	t.Log("running E2E on leagues")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	var cfg app.Config
	cfg.DB.DSN = dsn
	cfg.DB.MaxOpenConns = 25
	cfg.DB.MaxIdleConns = 25
	cfg.DB.MaxIdleTime = 15 * time.Minute

	application, err := app.New(cfg, logger)
	if err != nil {
		t.Fatal(err)
	}

	err = application.Models.League.Insert(&data.League{Name: "NHL"})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Get("http://localhost:3000/v1/leagues")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}

	if !strings.Contains(body, "NHL") {
		t.Fatalf("expected NHL in response: %s", body)
	}

	rows, err := application.Models.League.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	for _, l := range rows {
		t.Logf("leagues: name=%s", l.Name)
	}
}
