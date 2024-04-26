package model

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/uuid"
)

// An Agent is a c2 client that registers itself with the c2 server and takes commands.
type Agent struct {
	ID uuid.UUID `json:"id"`
	IP net.IP    `json:"ip"`
}

// newAgent is an Agent constructor.
func newAgent(id uuid.UUID, ip net.IP) *Agent {
	return &Agent{
		ID: id,
		IP: ip,
	}
}

// ParseAgentModel Parses and validates Agent from string inputs.
func ParseAgentModel(id, ip string) (*Agent, error) {
	newID, err := uuid.Parse(id)
	if err != nil || newID == uuid.Nil || newID.String() == "" {
		return nil, fmt.Errorf("failed to parse agent ID: %w", err)
	}

	newIP := net.ParseIP(ip)
	if newIP == nil {
		return nil, errors.New("invalid agent IP")
	}

	return newAgent(newID, newIP), nil
}

// AgentDB is a database abstraction interface for an Agent database.
type AgentDB interface {
	// Insert inserts Agent into the database.
	Insert(*Agent) error
}
