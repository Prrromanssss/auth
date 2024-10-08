package model

import (
	"database/sql"
	"time"
)

// CreateUserParams holds the parameters for creating a new user.
type CreateUserParams struct {
	Name           string `db:"name"`
	Email          string `db:"email"`
	HashedPassword string `db:"hashed_password"`
	Role           int64  `db:"role_id"`
}

// CreateUserResponse represents the response after creating a new user.
type CreateUserResponse struct {
	UserID    int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// GetUserParams holds the parameters for retrieving a user by ID.
type GetUserParams struct {
	UserID int64 `db:"id"`
}

// GetUserResponse represents the details of a user retrieved from the database.
type GetUserResponse struct {
	UserID    int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// UpdateUserParams holds the parameters for updating an existing user.
type UpdateUserParams struct {
	UserID int64          `db:"id"`
	Name   sql.NullString `db:"name"`
	Role   int64          `db:"role_id"`
}

// UpdateUserResponse represents the response after creating a new user.
type UpdateUserResponse struct {
	UserID    int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DeleteUserParams holds the parameters for deleting a user by ID.
type DeleteUserParams struct {
	UserID int64 `db:"id"`
}
