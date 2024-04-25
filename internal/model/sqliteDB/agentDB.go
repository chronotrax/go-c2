package sqliteDB

import (
	"log/slog"

	"github.com/chronotrax/go-c2/internal/model"
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

// Insert inserts the [model.Agent] into the database.
func (s *AgentDB) Insert(model *model.Agent) error {
	schema := toDBAgent(model)

	const InsertAgent = `INSERT INTO agents (id, ip) VALUES (?, ?)`
	result, err := s.DB.Exec(InsertAgent, schema.ID, schema.IP)
	if err != nil {
		slog.Error("error inserting agent", slog.String("error", err.Error()))
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil || affected <= 0 {
		return err
	}
	return nil
}
