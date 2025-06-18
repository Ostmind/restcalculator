package application

import (
	"context"
	"log/slog"
	"restcalculator/internal/http-server/handlers"
	"restcalculator/internal/model"
	srv "restcalculator/internal/server"
	"time"
)

type App struct {
	server *srv.Server
	logger *slog.Logger
}

func New(logger *slog.Logger, port int, result *model.Results) *App {

	handler := handlers.New(logger, result)

	server := srv.New(handler, logger, port)

	return &App{
		server: server,
		logger: logger,
	}
}

func (a App) Run() {
	a.logger.Info("Starting app...")
	go a.server.Run()
}

func (a App) Stop(shutdownTimeout time.Duration) {
	a.logger.Info("Stopping app...")

	timeout := shutdownTimeout * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	doneCh := make(chan error)
	go func() {
		doneCh <- a.server.Stop(ctx)
	}()

	select {
	case err := <-doneCh:
		if err != nil {
			a.logger.Error("Error while stopping server: %v", err)
		}
		a.logger.Info("App has been stopped gracefully")

	case <-ctx.Done():
		a.logger.Warn("App stopped forced")
	}
}
