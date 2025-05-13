package sum

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"restcalculator/internal/http-server/model"

	"github.com/labstack/echo/v4"
)

type SumController struct {
	logger *slog.Logger
}

func New(log *slog.Logger) *SumController {
	return &SumController{log}
}

func (ctr SumController) Sum(c echo.Context) error {
	var (
		req []model.Request
		res model.Results
	)

	var sum float64
	var calculation model.Calculator
	ctr.logger.Info("Get Request for Sum")
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

	res.Calc = append(res.Calc, calculation)

	return c.JSON(http.StatusOK, calculation)
}
