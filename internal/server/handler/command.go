package handler

import (
	"net/http"

	"github.com/chronotrax/go-c2/internal/util"
	"github.com/chronotrax/go-c2/pkg/msgqueue"
	"github.com/labstack/echo/v4"
)

// POST /server/command/:id
func CommandPost(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Get command from body
	params := new(struct {
		Command string `json:"command"`
	})
	if err = c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Append command to message queue
	err = d.MsgQueue.AddMsg(id, msgqueue.NewMsg(params.Command))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// POST /server/command
func CommandAllPost(d *Depends, c echo.Context) error {
	// Get command from body
	params := new(struct {
		Command string `json:"command"`
	})
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Append command to message queues
	d.MsgQueue.AddMsgAll(msgqueue.NewMsg(params.Command))
	return c.NoContent(http.StatusNoContent)
}

// GET /agent/command/:id
func CommandGet(d *Depends, c echo.Context) error {
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

// DELETE /server/command/:id
func CommandDelete(d *Depends, c echo.Context) error {
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

	return c.NoContent(http.StatusNoContent)
}

// DELETE /server/command
func CommandDeleteAll(d *Depends, c echo.Context) error {
	// Delete most recent command from message queues
	d.MsgQueue.DeleteMsgAll()
	return c.NoContent(http.StatusNoContent)
}
