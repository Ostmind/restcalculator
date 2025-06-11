package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"restcalculator/internal/http-server/controller"
	"restcalculator/internal/model"
)

type Controller struct {
	logger *slog.Logger
	result *model.Results
}

func New(log *slog.Logger, res *model.Results) *Controller {
	return &Controller{log, res}
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

	res, err := controller.MakeCalc(req.Value, operation, ctr.result, userCookie.Value)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, res)
}
