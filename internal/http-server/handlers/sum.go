package handlers

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"restcalculator/internal/http-server/logger"
	"restcalculator/internal/http-server/model"
)

type SumController struct {
	ctx    context.Context
	logger *logger.Logger
}

func NewSumController(ctx context.Context, slogger *logger.Logger) *SumController {
	return &SumController{
		ctx:    ctx,
		logger: slogger,
	}
}

func (ctr *SumController) Sum(c echo.Context) error {
	var (
		req []model.Request
		res model.Results
	)
	ctr.logger.Info("Get Request on Sum")
	var sum float64
	var calculation model.Calculator

	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	for i := range req {

		sum += req[i].Value

		calculation.Numbers = append(calculation.Numbers, req[i].Value)
	}

	calculation.Operation = "+"
	calculation.Result = sum
	ctr.logger.Info("Get Sum: ", slog.Float64("Sum", sum))

	res.Calc = append(res.Calc, calculation)

	return c.JSON(http.StatusOK, calculation)
}
