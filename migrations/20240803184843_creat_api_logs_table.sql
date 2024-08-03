-- +goose Up
CREATE TABLE users.api_user_log (
    id integer GENERATED ALWAYS AS IDENTITY,
    user_id integer,
    action_type VARCHAR(50) NOT NULL,
    request_data JSONB NOT NULL,
    response_data JSONB,
    timestamp TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users.api_user_log;
