package models

import (
	"time"

	"github.com/guregu/null"
)

type CreateUserParams struct {
	Name           string
	Email          string
	HashedPassword string
	Role           int64
}

type GetUserResponse struct {
	UserID    int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      int64     `db:"role_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UpdateUserParams struct {
	UserID int64
	Name   null.String
	Role   int64
}
