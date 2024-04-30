package handler

import (
	"github.com/chronotrax/go-c2/internal/server/model"
	"github.com/chronotrax/go-c2/pkg/msgqueue"

	"github.com/labstack/echo/v4"
)

// Depends lists the dependencies needed for handlers.
type Depends struct {
	MsgQueue     msgqueue.MessageQueue
	AgentStore   model.AgentStore
	CommandStore model.CommandStore
}

// NewDepends is a Depends constructor.
func NewDepends(msgQ msgqueue.MessageQueue, agentStore model.AgentStore, cmdStore model.CommandStore) *Depends {
	return &Depends{MsgQueue: msgQ, AgentStore: agentStore, CommandStore: cmdStore}
}

// AppHandler is a custom [echo.HandlerFunc] handler that includes this app's dependencies.
type AppHandler func(d *Depends, c echo.Context) error

// Handle calls an AppHandler's [echo.HandlerFunc] with Depends.
func Handle(depends *Depends, handler AppHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(depends, c)
	}
}
