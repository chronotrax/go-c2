package sqliteDB

import (
	"database/sql"
	"errors"
	"log/slog"
	"net"

	"github.com/chronotrax/go-c2/internal/server/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AgentDB is a SQLite3 database for [model.Agent].
type AgentDB struct {
	DB *sqlx.DB
}

// NewAgentDB is an AgentDB constructor.
func NewAgentDB(db *sqlx.DB) *AgentDB {
	return &AgentDB{DB: db}
}

// agent is an internal sqlite schema of a [model.Agent].
type agent struct {
	ID string
	IP string
}

// toDBAgent converts [model.Agent] to agent.
// Assumes [model.Agent] is a valid agent.
func toDBAgent(a *model.Agent) *agent {
	return &agent{a.ID.String(), a.IP.String()}
}

// toModelAgent converts agent to [model.Agent].
// Assumes agent is a valid [model.Agent].
func toModelAgent(a *agent) *model.Agent {
	return &model.Agent{
		ID: uuid.MustParse(a.ID),
		IP: net.ParseIP(a.IP),
	}
}

// Insert inserts the [model.Agent] into the database.
func (d *AgentDB) Insert(agentModel *model.Agent) (rowsAffected int64, err error) {
	a := toDBAgent(agentModel)

	const InsertAgent = `INSERT INTO agents (id, ip) VALUES (?, ?)`
	result, err := d.DB.Exec(InsertAgent, a.ID, a.IP)
	if err != nil {
		slog.Error("failed to insert agent", slog.String("error", err.Error()))
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}

// Exists checks if [model.Agent] with id exists in the database.
func (d *AgentDB) Exists(id uuid.UUID) (bool, error) {
	const SelectAgent = `SELECT * FROM agents WHERE id = ? LIMIT 1`
	a := &agent{}
	err := d.DB.Get(a, SelectAgent, id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Get gets the [model.Agent] with id.
func (d *AgentDB) Get(id uuid.UUID) (*model.Agent, error) {
	const SelectAgent = `SELECT * FROM agents WHERE id = ? LIMIT 1`
	a := &agent{}
	err := d.DB.Get(a, SelectAgent, id)
	if err != nil {
		return nil, err
	}

	return toModelAgent(a), nil
}

// Delete deletes the [model.Agent] with id.
func (d *AgentDB) Delete(id uuid.UUID) (rowsAffected int64, err error) {
	const DeleteAgent = `DELETE FROM agents WHERE id = ?`
	result, err := d.DB.Exec(DeleteAgent, id)
	if err != nil {
		slog.Error("failed to delete agent", slog.String("error", err.Error()))
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return affected, nil
}
