package user

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/pkg/errors"

	"github.com/Prrromanssss/auth/internal/client/db"
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
			"userPGRepo.CreateUser.DB.ScanOneContext.queryCreateUser(email: %s)",
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
			"userPGRepo.GetUser.DB.ScanOneContext.queryGetUser(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return converter.ConvertGetUserResponseFromRepoToService(respRepo), nil
}

// UpdateUser modifies the details of an existing user in the database.
func (p *userPGRepo) UpdateUser(
	ctx context.Context,
	params model.UpdateUserParams,
) (err error) {
	log.Infof("userPGRepo.UpdateUser, params: %+v", params)

	paramsRepo := converter.ConvertUpdateUserParamsFromServiceToRepo(params)

	q := db.Query{
		Name:     "userPGRepo.UpdateUser",
		QueryRaw: queryUpdateUser,
	}

	_, err = p.db.DB().ExecContext(ctx, q, paramsRepo.UserID, paramsRepo.Name, paramsRepo.Role)
	if err != nil {
		return errors.Wrapf(
			err,
			"userPGRepo.UpdateUser.DB.ExecContext.queryUpdateUser(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return nil
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
			"userPGRepo.DeleteUser.DB.ExecContext.queryDeleteUser(userID: %d)",
			paramsRepo.UserID,
		)
	}

	return nil
}

// CreateAPILog creates log in database of every api action.
func (p *userPGRepo) CreateAPILog(
	ctx context.Context,
	params model.CreateAPILogParams,
) (err error) {
	log.Infof("userPGRepo.CreateAPILog, params: %+v", params)

	paramsRepo := converter.ConvertCreateAPILogParamsFromServiceToRepo(params)

	q := db.Query{
		Name:     "userPGRepo.CreateAPILog",
		QueryRaw: queryCreateAPILog,
	}

	_, err = p.db.DB().ExecContext(ctx, q, paramsRepo.Method, paramsRepo.RequestData, paramsRepo.ResponseData)
	if err != nil {
		return errors.Wrapf(
			err,
			"userPGRepo.CreateAPILog.DB.ExecContext.queryCreateAPILog",
		)
	}

	return nil
}
