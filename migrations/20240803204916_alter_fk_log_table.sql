-- +goose Up
ALTER TABLE users.api_user_log
ADD CONSTRAINT api_user_log_user_id_fkey
FOREIGN KEY (user_id)
REFERENCES users.user (id);

-- +goose Down
ALTER TABLE users.api_user_log
DROP CONSTRAINT api_user_log_user_id_fkey;
