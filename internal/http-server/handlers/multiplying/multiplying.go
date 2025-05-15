package multiplying

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"restcalculator/internal/model"

	"github.com/labstack/echo/v4"
)

type MultiplyingController struct {
	logger *slog.Logger
	result *model.Results
}

func New(log *slog.Logger, res *model.Results) *MultiplyingController {
	return &MultiplyingController{log, res}
}

func (ctr MultiplyingController) Multiplying(c echo.Context) error {
	select {
	case <-c.Request().Context().Done():
		return c.NoContent(http.StatusRequestTimeout)
	default:
		var (
			req         model.Request
			mult        float64 = 1
			calculation model.Calculation
			operation   = "*"
		)

		ctr.logger.Debug("Get Request for Multi")

		userCookie, err := c.Cookie("Token")
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		err = json.NewDecoder(c.Request().Body).Decode(&req)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		for i := range req.Value {

			mult *= req.Value[i]

			calculation.Numbers = append(calculation.Numbers, req.Value[i])
		}

		calculation.Operation = operation
		calculation.Result = mult

		ctr.result.UserValues[userCookie.Value] = append(ctr.result.UserValues[userCookie.Value], calculation)

		return c.JSON(http.StatusOK, calculation)

	}
}
