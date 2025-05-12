package main

import (
	"github.com/labstack/echo/v4"
	"restcalculator/handler"
)

func main() {
	e := echo.New()
	e.Use(handler.LogRequest)
	e.POST("/", handler.Handler)
	e.Logger.Fatal(e.Start(":8080"))
}
