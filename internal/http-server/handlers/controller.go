package handlers

import (
	"encoding/json"
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
	SUM  = "+"
	MULT = "*"
)

func makeCalc(num []float64, operation string, res *model.Results, cookie string, err chan struct{}, ok chan model.Calculation) {
	var calculation model.Calculation

	defer close(err)
	defer close(ok)

	switch operation {
	case SUM:
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

		ok <- calculation

	case MULT:
		var mult float64 = 1

		for i := range num {
			mult += num[i]

			calculation.Numbers = append(calculation.Numbers, num[i])
		}

		calculation.Operation = operation
		calculation.Result = mult

		res.Mutex.Lock()
		res.UserValues[cookie] = append(res.UserValues[cookie], calculation)
		res.Mutex.Unlock()

		ok <- calculation

	default:
		err <- struct{}{}
	}
}

func (ctr Controller) Calculation(c echo.Context) error {

	var (
		req          model.Request
		errorChannel chan struct{}
		resulChannel chan model.Calculation
	)

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

	go makeCalc(req.Value, operation, ctr.result, userCookie.Value, errorChannel, resulChannel)

	select {
	case calculation := <-resulChannel:
		return c.JSON(http.StatusOK, calculation)

	case <-errorChannel:
		return c.NoContent(http.StatusInternalServerError)

	}

}
