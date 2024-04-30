-- +goose up
CREATE TABLE agents
(
    id TEXT NOT NULL PRIMARY KEY,
    ip TEXT NOT NULL
);

-- +goose down
DROP TABLE agents;