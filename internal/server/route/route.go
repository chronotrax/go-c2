package route

import (
	"github.com/chronotrax/go-c2/internal/server/handler"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, d *handler.Depends) {
	e.GET("/ping", handler.Handle(d, handler.PingGet))

	e.POST("/agent/register/:id", handler.Handle(d, handler.AgentRegisterPost))
	e.GET("/agent/command/:id", handler.Handle(d, handler.AgentCommandGet))
	e.POST("/agent/command/:id", handler.Handle(d, handler.AgentCommandPost))

	e.DELETE("/server/agent/:id", handler.Handle(d, handler.AgentDelete))
	e.POST("/server/command/:id", handler.Handle(d, handler.ServerCommandPost))
	e.POST("/server/command", handler.Handle(d, handler.ServerCommandPostAll))
	e.DELETE("/server/command/:id", handler.Handle(d, handler.ServerCommandDelete))
	e.DELETE("/server/command", handler.Handle(d, handler.ServerCommandDeleteAll))
}
