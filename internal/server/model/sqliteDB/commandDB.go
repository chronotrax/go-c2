package sqliteDB

import (
	"log/slog"
	"strings"

	"github.com/chronotrax/go-c2/internal/server/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// CommandDB is a SQLite3 database for [model.Command]
type CommandDB struct {
	DB *sqlx.DB
}

// NewCommandDB is a CommandDB constructor.
func NewCommandDB(db *sqlx.DB) *CommandDB {
	return &CommandDB{DB: db}
}

// command is an internal sqlite schema of a [model.Agent].
type command struct {
	AgentID string
	MsgID   string
	Command string
	Args    string
	Output  string
}

// toDBCommand converts [model.Command] to command.
// Assumes [model.Command] is a valid command.
func toDBCommand(c *model.Command) *command {
	return &command{
		AgentID: c.AgentID.String(),
		MsgID:   c.MsgID.String(),
		Command: c.Command,
		Args:    strings.Join(c.Args, " "),
		Output:  c.Output,
	}
}

// toModelCommand converts command to [model.Command].
// Assumes command is a valid [model.Command].
func toModelCommand(c *command) *model.Command {
	return &model.Command{
		AgentID: uuid.MustParse(c.AgentID),
		MsgID:   uuid.MustParse(c.MsgID),
		Command: c.Command,
		Args:    strings.Split(c.Args, " "),
		Output:  c.Output,
	}
}

// Insert inserts the [model.Command] into the database.
func (d *CommandDB) Insert(cmdModel *model.Command) (rowsAffected int64, err error) {
	c := toDBCommand(cmdModel)

	const InsertCommand = `INSERT INTO commands (agentID, msgID, command, args, output) VALUES (?, ?, ?, ?, ?)`
	result, err := d.DB.Exec(InsertCommand, c.AgentID, c.MsgID, c.Command, c.Args, c.Output)
	if err != nil {
		slog.Error("failed to insert command into database", slog.String("error", err.Error()))
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Get gets the [model.Command] with agentID and messageID.
func (d *CommandDB) Get(agentID, msgID uuid.UUID) (*model.Command, error) {
	const SelectCommand = `SELECT * FROM commands WHERE agentID = ? AND msgID = ? LIMIT 1`
	c := &command{}
	err := d.DB.Get(c, SelectCommand, agentID, msgID)
	if err != nil {
		return nil, err
	}

	return toModelCommand(c), nil
}
