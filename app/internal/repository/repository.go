package repository

import (
	"context"

	"github.com/Prrromanssss/auth/internal/models"
)

// UserRepository defines methods for user data operations.
type UserRepository interface {
	// CreateUser adds a new user and returns the user ID and any error.
	CreateUser(ctx context.Context, params models.CreateUserParams) (userID int64, err error)

	// GetUser retrieves a user by ID and returns user details and any error.
	GetUser(ctx context.Context, userID int64) (resp models.GetUserResponse, err error)

	// UpdateUser updates user details by ID and returns any error.
	UpdateUser(ctx context.Context, params models.UpdateUserParams) (err error)

	// DeleteUser removes a user by ID and returns any error.
	DeleteUser(ctx context.Context, userID int64) (err error)
}
