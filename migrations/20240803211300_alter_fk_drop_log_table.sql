-- +goose Up
ALTER TABLE users.api_user_log
DROP COLUMN user_id;

-- +goose Down
ALTER TABLE users.api_user_log
ADD COLUMN user_id integer;

ALTER TABLE users.api_user_log
ADD CONSTRAINT api_user_log_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users.user (id);
