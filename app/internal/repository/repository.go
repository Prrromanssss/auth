package repository

import (
	"context"

	"github.com/Prrromanssss/auth/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params models.CreateUserParams) (userID int64, err error)
	GetUser(ctx context.Context, userID int64) (resp models.GetUserResponse, err error)
	UpdateUser(ctx context.Context, params models.UpdateUserParams) (err error)
	DeleteUser(ctx context.Context, userID int64) (err error)
}
