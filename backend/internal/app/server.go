package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (app *Application) Server() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
	}

	app.Logger.Info("starting api server", "addr", srv.Addr, "env", app.Config.Env)

	return srv.ListenAndServe()
}
