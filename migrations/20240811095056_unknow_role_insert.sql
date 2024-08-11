-- +goose Up
INSERT INTO users.user_role
    (title)
VALUES
    ('unknown');

-- +goose Down
DELETE FROM users.user_role
WHERE title IN ('unknown');


