package model

// User represents the user data structure used for caching.
type User struct {
	UserID    int64  `redis:"id"`
	Name      string `redis:"name"`
	Email     string `redis:"email"`
	Role      int64  `redis:"role"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt int64  `redis:"updated_at"`
}

// GetUserParams defines the parameters needed to retrieve user information.
type GetUserParams struct {
	UserID int64 `redis:"id"`
}

// DeleteUserParams defines the parameters needed to delete a user.
type DeleteUserParams struct {
	UserID int64 `redis:"id"`
}
