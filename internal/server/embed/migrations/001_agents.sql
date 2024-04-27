-- +goose up
CREATE TABLE agents
(
    id STRING NOT NULL PRIMARY KEY,
    ip STRING NOT NULL 
);

-- +goose down
DROP TABLE agents;