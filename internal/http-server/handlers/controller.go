package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"restcalculator/internal/model"

	"github.com/labstack/echo/v4"
)

type Controller struct {
	logger *slog.Logger
	result *model.Results
}

func New(log *slog.Logger, res *model.Results) *Controller {
	return &Controller{log, res}
}

const (
	SUM_OPERATION  = "+"
	MULT_OPERATION = "*"
)

func makeCalc(num []float64, operation string, res *model.Results, cookie string) (model.Calculation, error) {
	var (
		err         error
		calculation model.Calculation
	)

	switch operation {
	case SUM_OPERATION:
		var sum float64

		for i := range num {
			sum += num[i]

			calculation.Numbers = append(calculation.Numbers, num[i])
		}

		calculation.Operation = operation
		calculation.Result = sum

		res.Mutex.Lock()
		res.UserValues[cookie] = append(res.UserValues[cookie], calculation)
		res.Mutex.Unlock()

	case MULT_OPERATION:
		var mult float64 = 1

		for i := range num {
			mult *= num[i]

			calculation.Numbers = append(calculation.Numbers, num[i])
		}

		calculation.Operation = operation
		calculation.Result = mult

		res.Mutex.Lock()
		res.UserValues[cookie] = append(res.UserValues[cookie], calculation)
		res.Mutex.Unlock()

	default:
		err = errors.New("No Operation")
	}

	return calculation, err
}

func (ctr Controller) Calculation(c echo.Context) error {

	var req model.Request

	ctr.logger.Debug("Get Request for Sum")

	userCookie, err := c.Cookie("Token")
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	operation := c.Param("operation")

	res, err := makeCalc(req.Value, operation, ctr.result, userCookie.Value)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}
