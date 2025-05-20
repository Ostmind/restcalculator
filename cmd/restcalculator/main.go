package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"log/slog"
	"net/http"
	"os"
	"restcalculator/internal/config"
	"restcalculator/internal/http-server/handlers"
	"restcalculator/internal/http-server/middleware"
	"restcalculator/internal/logger"
	"restcalculator/internal/model"
	"time"
)

func main() {
	echoInstance := echo.New()

	config, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		log.Fatal("No config cannot start server", slog.Any("error", err))
	}

	sloger := logger.SetupLogger(config.EnvType)
	sloger.Info("starting calculator",
		slog.String("env", config.EnvType))

	result := model.NewResult()

	controller := handlers.New(sloger, result)

	echoInstance.Use(middleware.LogRequest)
	echoInstance.POST("/sum/{operation}", controller.Calculation)

	server := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  config.Timeout * time.Second,
		WriteTimeout: config.Timeout * time.Second,
	}

	startErr := make(chan error)

	go func() {
		if err := echoInstance.StartServer(server); err != nil {
			log.Fatal("server error", slog.Any("error", err))
			close(startErr)
		}
	}()

	stop := make(chan os.Signal, 1)
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout*time.Second)
	defer cancel()

	select {
	case <-startErr:
	case <-stop:
		log.Info("stopping server")
		if err := echoInstance.Shutdown(ctx); err != nil {
			log.Error("failed to gracefully shutdown server", slog.Any("error", err))
		}

		log.Info("server stopped")
	}

}
