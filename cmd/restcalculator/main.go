package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"restcalculator/internal/http-server/handlers/multiplying"
	"restcalculator/internal/http-server/handlers/sum"
	"restcalculator/internal/http-server/logger"
	"restcalculator/internal/http-server/middleware"
	"syscall"
	"time"
)

func main() {
	echoInstance := echo.New()
	envType := os.Getenv("EnvType")

	log := logger.SetupLogger(envType)
	log.Info("starting calculator",
		slog.String("env", envType))
	sumController := sum.New(log)
	multController := multiplying.New(log)
	echoInstance.Use(middleware.LogRequest)
	echoInstance.POST("/sum", sumController.Sum)
	echoInstance.POST("/mult", multController.Multiplying)

	startErr := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	go func() {
		if err := echoInstance.Start(":8080"); err != nil && err != http.ErrServerClosed {
			echoInstance.Logger.Fatal("shutting down the server")
			close(startErr)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-startErr:
	case <-stop:
		log.Info("stopping server")
		if err := echoInstance.Shutdown(ctx); err != nil {
			log.Error("failed to gracefully shutdown server")
		}

		log.Info("server stopped")
	}
}
