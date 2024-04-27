package route

import (
	"github.com/chronotrax/go-c2/internal/server/handler"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, d *handler.Depends) {
	e.GET("/ping", handler.Handle(d, handler.PingGet))

	e.POST("/agent/register/:id", handler.Handle(d, handler.RegisterPost))
	e.GET("/agent/command/:id", handler.Handle(d, handler.CommandGet))

	e.DELETE("/server/agent/:id", handler.Handle(d, handler.AgentDelete))
	e.POST("/server/command/:id", handler.Handle(d, handler.CommandPost))
	e.POST("/server/command", handler.Handle(d, handler.CommandAllPost))
	e.DELETE("/server/command/:id", handler.Handle(d, handler.CommandDelete))
	e.DELETE("/server/command", handler.Handle(d, handler.CommandDeleteAll))
}
