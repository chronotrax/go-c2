-- +goose up
CREATE TABLE commands
(
    agentID TEXT,
    msgID   TEXT,
    command TEXT,
    args    TEXT,
    output  TEXT,
    PRIMARY KEY (agentID, msgID)
);

-- +goose down
DROP TABLE commands;