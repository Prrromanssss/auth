package user

const (
	queryUpdateUser = `
		UPDATE users.user
		SET name = COALESCE($2, name), role_id = $3, updated_at = now()
		WHERE id = $1
	`

	queryCreateUser = `
		INSERT INTO users.user
			(name, email, hashed_password, role_id)
		VALUES
			($1, $2, $3, $4)
		RETURNING id;
	`

	queryGetUser = `
		SELECT
			id
			, name
			, email
			, role_id
			, created_at
			, updated_at
		FROM users.user
		WHERE id = $1
	`

	queryDeleteUser = `
		DELETE FROM users.user
		WHERE id = $1
	`
)
