package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

type Calculator struct {
	Numbers   []float64
	Operation string
	Result    float64
}

type Results struct {
	Calc []Calculator
}

type Request struct {
	Value float64
}

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()
		fmt.Printf("Request: %s %s %s %d\n", c.Request().Method, c.Request().URL, stop.Sub(start), c.Response().Status)
		return err
	}
}

func Handler(c echo.Context) error {
	var (
		req []Request
		res Results
	)

	var sum float64
	var calculation Calculator
	_ = json.NewDecoder(c.Request().Body).Decode(&req)
	for i := 0; i < len(req); i++ {

		sum += req[i].Value
		calculation.Numbers = append(calculation.Numbers, req[i].Value)
	}
	calculation.Operation = "+"
	calculation.Result = sum

	res.Calc = append(res.Calc, calculation)

	return c.JSON(200, calculation)

}
