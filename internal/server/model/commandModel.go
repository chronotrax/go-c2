package model

import "github.com/google/uuid"

// A Command is a message sent to an agent with its response.
type Command struct {
	AgentID uuid.UUID `json:"agentID"`
	MsgID   uuid.UUID `json:"msgID"`
	Command string    `json:"command"`
	Args    []string  `json:"args"`
	Output  string    `json:"output"`
}

// NewCommand is a Command constructor.
func NewCommand(agentID uuid.UUID, msgID uuid.UUID, command string, args []string, output string) *Command {
	return &Command{
		AgentID: agentID,
		MsgID:   msgID,
		Command: command,
		Args:    args,
		Output:  output,
	}
}

// CommandStore is a database abstraction interface for a Command database.
type CommandStore interface {
	Insert(*Command) (rowsAffected int64, err error)
	Get(agentID, msgID uuid.UUID) (*Command, error)
}
