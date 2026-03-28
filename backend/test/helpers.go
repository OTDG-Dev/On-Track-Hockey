package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/app"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	tclog "github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	MIGRATIONS_PATH = "file://../db/migrations"

	sharedSetupMu   sync.Mutex
	sharedSetupOnce sync.Once
	sharedApp       *app.Application
	sharedBaseURL   string
	sharedContainer testcontainers.Container
)

// Discard verbose logs from test containers unless a test fails.
func init() {
	tclog.SetDefault(log.New(io.Discard, "", 0))
}

// setup starts a postgres container, runs migrations, and starts the API server.
// Optional flag: setup(t, true) forces a clean recreate of shared infra.
// Returns the application instance and the base URL for making requests.
func setup(t *testing.T, cleanRecreate ...bool) (*app.Application, string) {
	t.Helper()

	forceRecreate := len(cleanRecreate) > 0 && cleanRecreate[0]

	sharedSetupMu.Lock()
	defer sharedSetupMu.Unlock()

	if forceRecreate {
		teardownShared()
		sharedApp = nil
		sharedBaseURL = ""
		sharedContainer = nil
		sharedSetupOnce = sync.Once{}
	}

	// Initialize expensive infrastructure only once.
	sharedSetupOnce.Do(func() {
		port := getFreePort(t)
		sharedBaseURL = fmt.Sprintf("http://%s:%d", "localhost", port)

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var dsn string
		sharedContainer, dsn = startPostgres(t, ctx)
		runMigrations(t, dsn)

		sharedApp = startAPI(t, dsn, port)
		waitForAPI(t, sharedBaseURL+"/v1/healthcheck")
	})

	// Wipe the database clean before each test to ensure test isolation.
	resetDB(t, sharedApp)

	return sharedApp, sharedBaseURL
}

func teardownShared() {
	if sharedApp != nil {
		sharedApp.Shutdown()
	}

	if sharedContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		_ = sharedContainer.Terminate(ctx)
	}
}

// startPostgres starts a postgres container and returns the DSN for connecting to it.
func startPostgres(t *testing.T, ctx context.Context) (testcontainers.Container, string) {
	t.Helper()

	c, err := postgres.Run(ctx, "postgres:18.1-alpine",
		postgres.WithDatabase("ontrackhockey"),
		postgres.WithUsername("othadmin"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	if err != nil {
		t.Fatal(err)
	}

	dsn, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return c, dsn
}

func runMigrations(t *testing.T, dsn string) {
	t.Helper()

	m, err := migrate.New(MIGRATIONS_PATH, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}
}

func startAPI(t *testing.T, dsn string, port int) *app.Application {
	t.Helper()

	app, err := app.New(app.Config{
		Port: port,
		Env:  "test",
		DB: app.DBConfig{
			DSN:          dsn,
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  15 * time.Minute,
		},
		Limiter: app.LimiterConfig{
			RPS: 0, Burst: 0, Enabled: false,
		},
	}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}

	go app.Server()

	return app
}

func waitForAPI(t *testing.T, url string) {
	t.Helper()

	client := http.Client{Timeout: 500 * time.Millisecond}
	deadline := time.Now().Add(10 * time.Second)

	var lastErr error
	lastStatus := 0

	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if resp != nil {
			lastStatus = resp.StatusCode
			_ = resp.Body.Close()
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		if err != nil {
			lastErr = err
		}
		time.Sleep(100 * time.Millisecond)
	}

	t.Fatalf("server did not start: last_status=%d last_error=%v", lastStatus, lastErr)
}

// seedLeague creates a league via POST and returns its ID.
func seedLeague(t *testing.T, url, name string) int {
	t.Helper()

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(struct {
		Name string `json:"name"`
	}{Name: name}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/leagues", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedLeague: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		League struct {
			ID int `json:"id"`
		} `json:"league"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.League.ID
}

// seedDivision creates a division via POST and returns its ID.
func seedDivision(t *testing.T, url string, leagueID int, name string) int {
	t.Helper()

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(struct {
		Name     string `json:"name"`
		LeagueID int    `json:"league_id"`
	}{Name: name, LeagueID: leagueID}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/divisions", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedDivision: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		Division struct {
			ID int `json:"id"`
		} `json:"division"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.Division.ID
}

// seedTeam creates a team via POST and returns its ID.
func seedTeam(t *testing.T, url string, divisionID int, fullName, shortName string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(struct {
		FullName   string `json:"full_name"`
		ShortName  string `json:"short_name"`
		DivisionID int    `json:"division_id"`
		IsActive   bool   `json:"is_active"`
	}{
		FullName:   fullName,
		ShortName:  shortName,
		DivisionID: divisionID,
		IsActive:   true,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/teams", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedTeam: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		Team struct {
			ID int `json:"id"`
		} `json:"team"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.Team.ID
}

// resetDB wipes the database clean between tests by truncating all tables and resetting IDs.
// Rather than recreating the entire container, this is much faster and allows us to keep the API server running.
func resetDB(t *testing.T, application *app.Application) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := application.DB.QueryContext(ctx, `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		  AND table_name IN ('game_event_participants', 'game_events', 'games', 'roster', 'players', 'teams', 'divisions', 'leagues')
	`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Fatal(err)
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if len(tables) == 0 {
		// Nothing to reset yet (safe no-op for partially migrated schemas).
		return
	}

	quotedTables := make([]string, 0, len(tables))
	for _, table := range tables {
		quotedTables = append(quotedTables, fmt.Sprintf("\"%s\"", table))
	}

	_, err = application.DB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(quotedTables, ", ")))
	if err != nil {
		t.Fatal(err)
	}
}

// getFreePort asks the kernel for a free open port that is ready to use.
func getFreePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("failed to resolve tcp address")
	}

	return addr.Port
}
