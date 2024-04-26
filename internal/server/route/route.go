package route

import (
	"github.com/chronotrax/go-c2/internal/server/handler"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, d *handler.Depends) {
	e.GET("/ping", handler.Handle(d, handler.PingGet))

	e.POST("/agent/register", handler.Handle(d, handler.RegisterPost))
}
