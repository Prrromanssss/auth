package cache

import (
	"context"

	"github.com/Prrromanssss/auth/internal/model"
)

// UserCache defines the interface for user cache operations.
type UserCache interface {
	// Create adds a new user to the cache.
	Create(ctx context.Context, params model.User) (err error)
	// Get retrieves user information from the cache.
	Get(ctx context.Context, params model.GetUserParams) (resp model.GetUserResponse, err error)
	// Delete removes a user from the cache.
	Delete(ctx context.Context, params model.DeleteUserParams) (err error)
}
