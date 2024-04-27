package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/chronotrax/go-c2/internal/server/model"
	"github.com/chronotrax/go-c2/internal/util"

	"github.com/labstack/echo/v4"
)

// POST /agent/register/:id
func RegisterPost(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Check if ID already exists
	exists, err := d.AgentStore.Exists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, internalServerError)
	}
	if exists {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "duplicate UUID, please generate a new UUID"})
	}

	// Get agent's IP address
	ipStr := c.RealIP()
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Convert to model
	agent := model.NewAgent(id, ip)

	// Insert in db
	rows, err := d.AgentStore.Insert(agent)
	if err != nil || rows != 1 {
		return c.JSON(http.StatusInternalServerError, internalServerError)
	}

	// Register with message queue
	d.MsgQueue.Register(agent.ID)
	return c.NoContent(http.StatusNoContent)
}

// DELETE /server/agent/:id
func AgentDelete(d *Depends, c echo.Context) error {
	// Get ID from URL
	idStr := c.Param("id")
	id, err := util.ValidateUUID(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newParseError(err))
	}

	// Check if ID exists
	exists, err := d.AgentStore.Exists(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, internalServerError)
	}
	if !exists {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": fmt.Sprintf("agent does not exist with id: %v", id)})
	}

	// Delete agent from db
	rows, err := d.AgentStore.Delete(id)
	if err != nil || rows != 1 {
		return c.JSON(http.StatusInternalServerError, internalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
