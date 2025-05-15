package sum

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"restcalculator/internal/model"

	"github.com/labstack/echo/v4"
)

type SumController struct {
	logger *slog.Logger
	result *model.Results
}

func New(log *slog.Logger, res *model.Results) *SumController {
	return &SumController{log, res}
}

func (ctr SumController) Sum(c echo.Context) error {
	select {
	case <-c.Request().Context().Done():
		return c.NoContent(http.StatusRequestTimeout)
	default:
		var (
			req         model.Request
			sum         float64 = 0
			calculation model.Calculation
			operation   = "+"
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

		for i := range req.Value {

			sum += req.Value[i]

			calculation.Numbers = append(calculation.Numbers, req.Value[i])
		}

		calculation.Operation = operation
		calculation.Result = sum

		ctr.result.UserValues[userCookie.Value] = append(ctr.result.UserValues[userCookie.Value], calculation)

		return c.JSON(http.StatusOK, calculation)
	}
}
