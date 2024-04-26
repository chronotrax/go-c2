package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//goland:noinspection GoUnusedParameter
func PingGet(d *Depends, c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"response": "pong"})
}
