package user

import (
	"context"

	"github.com/Prrromanssss/platform_common/pkg/db"
	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"

	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/internal/repository/user/converter"
	modelRepo "github.com/Prrromanssss/auth/internal/repository/user/model"
)

type userPGRepo struct {
	db db.Client
}

// NewRepository creates a new instance of userPGRepo with the provided database connection.
func NewRepository(db db.Client) repository.UserRepository {
	return &userPGRepo{db: db}
}

// CreateUser inserts a new user into the database with the provided parameters.
func (p *userPGRepo) CreateUser(
	ctx context.Context,
	params model.CreateUserParams,
) (resp model.CreateUserResponse, err error) {
	log.Infof("userPGRepo.CreateUser, params: %+v", params)

	paramsRepo := converter.ConvertCreateUserParamsFromServiceToRepo(params)

	var respRepo modelRepo.CreateUserResponse

	q := db.Query{
		Name:     "userPGRepo.CreateUser",
		QueryRaw: queryCreateUser,
	}

	err = p.db.DB().ScanOneContext(ctx, &respRepo, q, paramsRepo.Name, paramsRepo.Email, params.HashedPassword, params.Role)
	if err != nil {
		return resp, errors.Wrapf(
			err,
			"Cannot create user(email: %s)",
			paramsRepo.Email,
		)
	}

	return converter.ConvertCreateUserResponseFromRepoToService(respRepo), nil
}

// GetUser retrieves a user from the database by their ID.
func (p *userPGRepo) GetUser(
	ctx context.Context,
	params model.GetUserParams,
) (resp model.GetUserResponse, err error) {
	log.Infof("userPGRepo.GetUser, params: %+v", params)

	paramsRepo := converter.ConvertGetUserParamsFromServiceToRepo(params)

	var respRepo modelRepo.GetUserResponse

	q := db.Query{
		Name:     "userPGRepo.GetUser",
		QueryRaw: queryGetUser,
	}

	err = p.db.DB().ScanOneContext(ctx, &respRepo, q, paramsRepo.UserID)
	if err != nil {
		return resp, errors.Wrapf(
			err,
			"Cannot get user(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return converter.ConvertGetUserResponseFromRepoToService(respRepo), nil
}

// UpdateUser modifies the details of an existing user in the database.
func (p *userPGRepo) UpdateUser(
	ctx context.Context,
	params model.UpdateUserParams,
) (resp model.UpdateUserResponse, err error) {
	log.Infof("userPGRepo.UpdateUser, params: %+v", params)

	paramsRepo := converter.ConvertUpdateUserParamsFromServiceToRepo(params)

	q := db.Query{
		Name:     "userPGRepo.UpdateUser",
		QueryRaw: queryUpdateUser,
	}

	var respRepo modelRepo.UpdateUserResponse

	err = p.db.DB().ScanOneContext(ctx, &respRepo, q, paramsRepo.UserID, paramsRepo.Name, paramsRepo.Role)
	if err != nil {
		return resp, errors.Wrapf(
			err,
			"Cannot update user(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return converter.ConvertUpdateUserResponseFromRepoToService(respRepo), nil
}

// DeleteUser removes a user from the database by their ID.
func (p *userPGRepo) DeleteUser(
	ctx context.Context,
	params model.DeleteUserParams,
) (err error) {
	log.Infof("userPGRepo.DeleteUser, params: %+v", params)

	paramsRepo := converter.ConvertDeleteUserParamsFromServiceToRepo(params)

	q := db.Query{
		Name:     "userPGRepo.DeleteUser",
		QueryRaw: queryDeleteUser,
	}

	_, err = p.db.DB().ExecContext(ctx, q, paramsRepo.UserID)
	if err != nil {
		return errors.Wrapf(
			err,
			"Cannot delete user(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return nil
}
