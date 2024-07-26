package models

import (
	"time"
)

// CreateUserParams holds the parameters for creating a new user.
type CreateUserParams struct {
	Name           string // User's name
	Email          string // User's email
	HashedPassword string // User's hashed password
	Role           int64  // User's role
}

// GetUserResponse represents the details of a user retrieved from the database.
type GetUserResponse struct {
	UserID    int64     `db:"id"`         // User identifier
	Name      string    `db:"name"`       // User's name
	Email     string    `db:"email"`      // User's email
	Role      int64     `db:"role_id"`    // User's role
	CreatedAt time.Time `db:"created_at"` // Creation timestamp
	UpdatedAt time.Time `db:"updated_at"` // Last update timestamp
}

// UpdateUserParams holds the parameters for updating an existing user.
type UpdateUserParams struct {
	UserID int64  // User identifier
	Name   string // New name for the user
	Role   int64  // New role for the user
}
