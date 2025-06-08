package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"restcalculator/internal/config"
	"restcalculator/internal/http-server/handlers"
	"restcalculator/internal/http-server/middleware"
	"restcalculator/internal/logger"
	"restcalculator/internal/model"
	"time"
)

func main() {
	echoInstance := echo.New()

	cfg, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		log.Fatal("No config cannot start server", slog.Any("error", err))
	}

	sloger := logger.SetupLogger(cfg.EnvType)
	sloger.Info("starting calculator",
		slog.String("env", cfg.EnvType))

	result := model.NewResult()

	controller := handlers.New(sloger, result)

	echoInstance.Use(middleware.LogRequest)
	echoInstance.POST("/{operation}", controller.Calculation)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ShutdownTimeout * time.Second,
		WriteTimeout: cfg.ShutdownTimeout * time.Second,
	}

	startErr := make(chan error)

	go func() {
		if err := echoInstance.StartServer(server); err != nil {
			startErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout*time.Second)
	defer cancel()

	select {
	case err = <-startErr:
		log.Fatal("server error", slog.Any("error", err))
	case <-stop:
		log.Info("stopping server")
		if err := echoInstance.Shutdown(ctx); err != nil {
			log.Error("failed to gracefully shutdown server", slog.Any("error", err))
		}

		log.Info("server stopped")
	}

}
