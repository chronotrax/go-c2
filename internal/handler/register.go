package handler

import (
	"errors"
	"fmt"
	"modernc.org/sqlite"
	"net/http"

	"github.com/chronotrax/go-c2/internal/model"

	"github.com/labstack/echo/v4"
)

func RegisterPost(d *Depends, c echo.Context) error {
	params := new(struct {
		ID string `json:"id"`
	})
	if err := c.Bind(params); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("failed to parse request: %s", err)})
	}

	ipStr := c.RealIP()

	agent, err := model.ParseAgentModel(params.ID, ipStr)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": err.Error()})
	}

	err = d.AgentStore.Insert(agent)

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == 1555 { // duplicate UUID primary key
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": "duplicate UUID, please generate a new UUID"})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return c.NoContent(http.StatusCreated)
}
