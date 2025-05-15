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
	"restcalculator/internal/http-server/handlers/multiplying"
	"restcalculator/internal/http-server/handlers/sum"
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

	sumController := sum.New(sloger, result)

	multController := multiplying.New(sloger, result)

	echoInstance.Use(middleware.LogRequest)
	echoInstance.POST("/sum", sumController.Sum)
	echoInstance.POST("/mult", multController.Multiplying)

	server := &http.Server{
		Addr:         ":" + config.Port,
		ReadTimeout:  config.Timeout * time.Second,
		WriteTimeout: config.Timeout * time.Second,
	}

	startErr := make(chan error)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := echoInstance.StartServer(server); err != nil {
			log.Fatal("server error", slog.Any("error", err))
			startErr <- err
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout*time.Second)
	defer cancel()
	if err := echoInstance.Shutdown(ctx); err != nil {
		echoInstance.Logger.Fatal(err)
	}

}
