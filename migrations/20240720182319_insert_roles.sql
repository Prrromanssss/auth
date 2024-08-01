-- +goose Up
INSERT INTO users.user_role
    (title)
VALUES
    ('user'),
    ('admin');

-- +goose Down
DELETE FROM users.user_role
WHERE title IN ('user', 'admin');

