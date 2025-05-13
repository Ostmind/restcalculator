package multiplying

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"restcalculator/internal/http-server/model"

	"github.com/labstack/echo/v4"
)

type MultiplyingController struct {
	logger *slog.Logger
}

func New(log *slog.Logger) *MultiplyingController {
	return &MultiplyingController{log}
}

func multiply(numbers *[]model.Request) {

}

func (ctr MultiplyingController) Multiplying(c echo.Context) error {
	var (
		req []model.Request
		res model.Results
	)

	var mult float64 = 1
	var calculation model.Calculator
	ctr.logger.Info("Get Request for Multi")
	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	for i := range req {

		mult *= req[i].Value

		calculation.Numbers = append(calculation.Numbers, req[i].Value)
	}

	calculation.Operation = "+"
	calculation.Result = mult

	res.Calc = append(res.Calc, calculation)

	return c.JSON(http.StatusOK, calculation)
}
