package log

const (
	queryCreateAPILog = `
		INSERT INTO users.api_user_log
			(action_type, request_data, response_data)
		VALUES
			($1, $2, $3);
	`
)
