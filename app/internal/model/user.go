package model

import (
	"time"
)

// CreateUserParams holds the parameters for creating a new user.
type CreateUserParams struct {
	Name           string
	Email          string
	HashedPassword string
	Role           int64
}

// CreateUserResponse represents the response after creating a user, containing the user's ID.
type CreateUserResponse struct {
	UserID int64
}

// GetUserParams holds the parameters for retrieving a user by ID.
type GetUserParams struct {
	UserID int64
}

// GetUserResponse represents the details of a user retrieved from the database.
type GetUserResponse struct {
	UserID    int64
	Name      string
	Email     string
	Role      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpdateUserParams holds the parameters for updating an existing user.
type UpdateUserParams struct {
	UserID int64
	Name   *string
	Role   int64
}

// DeleteUserParams holds the parameters for deleting a user by ID.
type DeleteUserParams struct {
	UserID int64
}
