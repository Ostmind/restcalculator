package main

import (
	"github.com/labstack/gommon/log"
	"log/slog"
	"os"
	"os/signal"
	"restcalculator/internal/application"
	"restcalculator/internal/config"
	"restcalculator/internal/logger"
	"restcalculator/internal/model"
)

func main() {

	cfg, err := config.LoadConfig("internal/config/config.json")
	if err != nil {
		log.Fatal("No config cannot start server", slog.Any("error", err))
	}

	sloger := logger.SetupLogger(cfg.EnvType)
	sloger.Info("starting calculator",
		slog.String("env", cfg.EnvType))

	result := model.NewResult()

	app := application.New(sloger, cfg.Port, result)

	app.Run()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	log.Info("Recieved interrupt signal")
	app.Stop(cfg.ShutdownTimeout)
}
