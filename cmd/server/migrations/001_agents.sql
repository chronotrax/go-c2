-- +goose up
CREATE TABLE agents
(
    id STRING PRIMARY KEY,
    ip STRING
);

-- +goose down
DROP TABLE agents;