package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Prrromanssss/auth/internal/models"
	"github.com/Prrromanssss/auth/internal/repository"
)

type userPGRepo struct {
	db *sqlx.DB
}

func NewPGRepo(db *sqlx.DB) repository.UserRepository {
	return &userPGRepo{db: db}
}

func (p *userPGRepo) CreateUser(ctx context.Context, params models.CreateUserParams) (userID int64, err error) {
	err = p.db.GetContext(ctx, &userID, queryCreateUser, params.Name, params.Email, params.HashedPassword, params.Role)
	if err != nil {
		return 0, errors.Wrap(
			err,
			"userPGRepo.CreateUser.GetContext.queryCreateUser",
		)
	}

	return userID, nil
}

func (p *userPGRepo) GetUser(ctx context.Context, userID int64) (resp models.GetUserResponse, err error) {
	err = p.db.GetContext(ctx, &resp, queryGetUser, userID)
	if err != nil {
		return resp, errors.Wrapf(
			err,
			"userPGRepo.GetUser.GetContext.queryGetUser(userID: %d)",
			userID,
		)
	}

	return resp, nil
}

func (p *userPGRepo) UpdateUser(ctx context.Context, params models.UpdateUserParams) (err error) {
	_, err = p.db.ExecContext(ctx, queryUpdateUser, params.UserID, params.Name, params.Role)
	if err != nil {
		return errors.Wrapf(
			err,
			"userPGRepo.UpdateUser.ExecContext.queryUpdateUser(userID: %d)",
			params.UserID,
		)
	}

	return nil
}

func (p *userPGRepo) DeleteUser(ctx context.Context, userID int64) (err error) {
	_, err = p.db.ExecContext(ctx, queryDeleteUser, userID)
	if err != nil {
		return errors.Wrapf(
			err,
			"userPGRepo.DeleteUser.ExecContext.queryDeleteUser(userID: %d)",
			userID,
		)
	}

	return nil
}
