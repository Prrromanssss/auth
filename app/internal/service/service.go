package service

import (
	"context"

	"github.com/Prrromanssss/auth/internal/model"
)

// UserRepository defines methods for user data operations.
type UserService interface {
	// CreateUser adds a new user and returns the user ID and any error.
	CreateUser(ctx context.Context, params model.CreateUserParams) (resp model.CreateUserResponse, err error)

	// GetUser retrieves a user by ID and returns user details and any error.
	GetUser(ctx context.Context, params model.GetUserParams) (resp model.GetUserResponse, err error)

	// UpdateUser updates user details by ID and returns any error.
	UpdateUser(ctx context.Context, params model.UpdateUserParams) (err error)

	// DeleteUser removes a user by ID and returns any error.
	DeleteUser(ctx context.Context, params model.DeleteUserParams) (err error)
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
