-- +goose Up
CREATE SCHEMA users;

CREATE TABLE users.user_role (
    id integer GENERATED ALWAYS AS IDENTITY,
    title text NOT NULL,

    PRIMARY KEY(id)
);

CREATE TABLE users.user (
    id integer GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    hashed_password text NOT NULL, 
    role_id integer NOT NULL,
    created_at timestamp DEFAULT now(),
    updated_at timestamp DEFAULT now(),

    PRIMARY KEY(id),
    FOREIGN KEY(role_id) REFERENCES users.user_role(id)
);

-- +goose Down
DROP TABLE users.user;
DROP TABLE users.user_role;

DROP SCHEMA users;
