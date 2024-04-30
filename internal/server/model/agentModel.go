package model

import (
	"net"

	"github.com/google/uuid"
)

// An Agent is a c2 client that registers itself with the c2 server and takes commands.
type Agent struct {
	ID uuid.UUID `json:"id"`
	IP net.IP    `json:"ip"`
}

// NewAgent is an Agent constructor.
func NewAgent(id uuid.UUID, ip net.IP) *Agent {
	return &Agent{
		ID: id,
		IP: ip,
	}
}

// AgentStore is a database abstraction interface for an Agent database.
type AgentStore interface {
	// Insert inserts Agent into the database.
	Insert(*Agent) (rowsAffected int64, err error)

	// Exists checks if Agent with id exists in the database.
	Exists(uuid.UUID) (exists bool, err error)

	// Get gets the Agent with id from the database.
	Get(uuid.UUID) (*Agent, error)

	// Delete deletes Agent from the database.
	Delete(id uuid.UUID) (rowsAffected int64, err error)
}
