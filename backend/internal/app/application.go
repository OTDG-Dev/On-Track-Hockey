package app

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

type Config struct {
	Port    int
	Env     string
	DB      DBConfig
	Limiter LimiterConfig
}

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type LimiterConfig struct {
	RPS     float64
	Burst   int
	Enabled bool
}

type Application struct {
	Config Config
	Logger *slog.Logger
	DB     *sql.DB
	Models data.Models
}

func New(cfg Config, logger *slog.Logger) (*Application, error) {
	db, err := OpenDB(cfg)
	if err != nil {
		return nil, err
	}

	app := &Application{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Models: data.NewModel(db),
	}

	return app, nil
}

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxOpenConns)
	db.SetConnMaxIdleTime(cfg.DB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
