package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"restcalculator/internal/http-server/handlers"
	"restcalculator/internal/http-server/logger"
)

const (
	envLocal = "local"
)

func main() {
	ctx := context.Background()
	e := echo.New()
	l := logger.Get()
	l.Info("Starting Server")
	sumController := handlers.NewSumController(ctx, l)
	e.POST("/sum", sumController.Sum)
	e.Logger.Fatal(e.Start(":8080"))
}
