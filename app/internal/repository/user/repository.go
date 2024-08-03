package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/internal/repository/user/converter"
	modelRepo "github.com/Prrromanssss/auth/internal/repository/user/model"
)

type userPGRepo struct {
	db *sqlx.DB
}

// NewPGRepo creates a new instance of userPGRepo with the provided database connection.
func NewPGRepo(db *sqlx.DB) repository.UserRepository {
	return &userPGRepo{db: db}
}

// CreateUser inserts a new user into the database with the provided parameters.
func (p *userPGRepo) CreateUser(
	ctx context.Context,
	params *model.CreateUserParams,
) (resp *model.CreateUserResponse, err error) {
	paramsRepo := converter.ConvertCreateUserParamsFromServiceToRepo(params)

	var respRepo modelRepo.CreateUserResponse

	err = p.db.GetContext(ctx, &respRepo, queryCreateUser, paramsRepo)
	if err != nil {
		return resp, errors.Wrapf(
			err,
			"userPGRepo.CreateUser.GetContext.queryCreateUser(email: %s)",
			paramsRepo.Email,
		)
	}

	return converter.ConvertCreateUserResponseFromRepoToService(&respRepo), nil
}

// GetUser retrieves a user from the database by their ID.
func (p *userPGRepo) GetUser(
	ctx context.Context,
	params *model.GetUserParams,
) (resp *model.GetUserResponse, err error) {
	paramsRepo := converter.ConvertGetUserParamsFromServiceToRepo(params)

	var respRepo modelRepo.GetUserResponse

	err = p.db.GetContext(ctx, &respRepo, queryGetUser, paramsRepo)
	if err != nil {
		return resp, errors.Wrapf(
			err,
			"userPGRepo.GetUser.GetContext.queryGetUser(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return converter.ConvertGetUserResponseFromRepoToService(&respRepo), nil
}

// UpdateUser modifies the details of an existing user in the database.
func (p *userPGRepo) UpdateUser(
	ctx context.Context,
	params *model.UpdateUserParams,
) (err error) {
	paramsRepo := converter.ConvertUpdateUserParamsFromServiceToRepo(params)

	_, err = p.db.ExecContext(ctx, queryUpdateUser, paramsRepo)
	if err != nil {
		return errors.Wrapf(
			err,
			"userPGRepo.UpdateUser.ExecContext.queryUpdateUser(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return nil
}

// DeleteUser removes a user from the database by their ID.
func (p *userPGRepo) DeleteUser(
	ctx context.Context,
	params *model.DeleteUserParams,
) (err error) {
	paramsRepo := converter.ConvertDeleteUserParamsFromServiceToRepo(params)

	_, err = p.db.ExecContext(ctx, queryDeleteUser, paramsRepo)
	if err != nil {
		return errors.Wrapf(
			err,
			"userPGRepo.DeleteUser.ExecContext.queryDeleteUser(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return nil
}
