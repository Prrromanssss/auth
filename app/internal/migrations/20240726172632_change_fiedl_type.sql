-- +goose Up
ALTER TABLE users.user
    ALTER COLUMN email TYPE VARCHAR(255),
    ALTER COLUMN hashed_password TYPE VARCHAR(255),
    ALTER COLUMN name TYPE VARCHAR(255);

-- +goose Down
ALTER TABLE users.user
    ALTER COLUMN email TYPE TEXT,
    ALTER COLUMN hashed_password TYPE TEXT,
    ALTER COLUMN name TYPE TEXT;

