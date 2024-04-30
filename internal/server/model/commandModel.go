package model

import (
	"github.com/chronotrax/go-c2/pkg/msgqueue"
	"github.com/google/uuid"
)

// A Command is a message sent to an agent with its response.
type Command struct {
	AgentID          uuid.UUID `json:"agentID"`
	msgqueue.Message `json:"message"`
	Output           string `json:"output"`
}

// NewCommand is a Command constructor.
func NewCommand(agentID uuid.UUID, msgID uuid.UUID, command string, args []string, output string) *Command {
	return &Command{
		AgentID: agentID,
		Message: msgqueue.Message{
			MsgID:   msgID,
			Command: command,
			Args:    args,
		},
		Output: output,
	}
}

// CommandStore is a database abstraction interface for a Command database.
type CommandStore interface {
	Insert(*Command) (rowsAffected int64, err error)
	Get(agentID, msgID uuid.UUID) (*Command, error)
}
