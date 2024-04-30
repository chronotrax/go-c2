package handler

import (
	"net/http"

	"github.com/chronotrax/go-c2/internal/server/model"
	"github.com/chronotrax/go-c2/internal/util"
	"github.com/chronotrax/go-c2/pkg/msgqueue"
	"github.com/labstack/echo/v4"
)

// POST /server/command/:id
func ServerCommandPost(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Get command from body
	params := new(struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	})
	if err = c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Append command to message queue
	err = d.MsgQueue.AddMsg(id, msgqueue.NewMessage(params.Command, params.Args...))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// POST /server/command
func ServerCommandPostAll(d *Depends, c echo.Context) error {
	// Get command from body
	params := new(struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	})
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Append command to all message queues
	d.MsgQueue.AddMsgAll(msgqueue.NewMessage(params.Command, params.Args...))
	return c.NoContent(http.StatusOK)
}

// DELETE /server/command/:id
func ServerCommandDelete(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Delete most recent command from message queue
	err = d.MsgQueue.DeleteMsg(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

// DELETE /server/command
func ServerCommandDeleteAll(d *Depends, c echo.Context) error {
	// Delete most recent command from all message queues
	d.MsgQueue.DeleteMsgAll()
	return c.NoContent(http.StatusOK)
}

// GET /agent/command/:id
func AgentCommandGet(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Get command from message queue
	msg, err := d.MsgQueue.GetMsg(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, msg)
}

// POST /agent/command/:id
func AgentCommandPost(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Get command output from body
	params := new(struct {
		msgqueue.Message `json:"message"`
		Output           string `json:"output"`
	})
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Convert to model
	cmd := model.NewCommand(id, params.MsgID, params.Command, params.Args, params.Output)

	// Insert in db
	rows, err := d.CommandStore.Insert(cmd)
	if err != nil || rows != 1 {
		return c.JSON(http.StatusInternalServerError, internalServerError)
	}

	return c.NoContent(http.StatusOK)
}
