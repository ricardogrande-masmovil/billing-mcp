package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthController struct {}

func NewHealthController() healthController {
	return healthController{}
}

func (h healthController) IsHealthy(ectx echo.Context) error {
	return ectx.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}