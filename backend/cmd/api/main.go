package main

import (
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/app"
	_ "github.com/lib/pq"
)

func main() {
	var cfg app.Config

	flag.IntVar(&cfg.Port, "port", 3000, "API http server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development | staging | production)")

	flag.StringVar(&cfg.DB.DSN, "dsn", os.Getenv("OTH_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.RPS, "limiter-rps", getEnvFloat("LIMITER_RPS", 2), "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", getEnvInt("LIMITER_BURST", 4), "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", getEnvBool("LIMITER_ENABLED", true), "Enable rate limiter")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := application.Server(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
