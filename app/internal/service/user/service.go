package user

import (
	"context"
	"errors"

	"github.com/Prrromanssss/platform_common/pkg/db"
	"github.com/gofiber/fiber/v2/log"

	modelCache "github.com/Prrromanssss/auth/internal/cache/user/model"

	"github.com/Prrromanssss/auth/internal/cache"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/internal/service"
)

type userService struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	cacheClient    cache.UserCache
	txManager      db.TxManager
}

// NewService creates a new instance of userService with the provided UserRepository.
func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	cacheClient cache.UserCache,
	txManager db.TxManager,
) service.UserService {
	return &userService{
		userRepository: userRepository,
		logRepository:  logRepository,
		cacheClient:    cacheClient,
		txManager:      txManager,
	}
}

// CreateUser creates a new user using the provided CreateUserParams.
func (s *userService) CreateUser(
	ctx context.Context,
	params model.CreateUserParams,
) (resp model.CreateUserResponse, err error) {
	log.Infof("userService.CreateUser, params: %+v", params)

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error

		resp, txErr = s.userRepository.CreateUser(ctx, params)
		if txErr != nil {
			return txErr
		}

		txErr = s.logRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Create",
			RequestData:  params,
			ResponseData: resp,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return model.CreateUserResponse{}, err
	}

	cacheErr := s.cacheClient.Create(ctx, resp.User)
	if cacheErr != nil {
		log.Warnf("Failed to create user in cache, params: %+v, err: %+v", params, cacheErr)
	}

	return resp, nil
}

// GetUser retrieves a user's information based on the provided GetUserParams.
func (s *userService) GetUser(
	ctx context.Context,
	params model.GetUserParams,
) (resp model.GetUserResponse, err error) {
	log.Infof("userService.GetUser, params: %+v", params)

	user, cacheErr := s.cacheClient.Get(ctx, params)
	if cacheErr != nil && !errors.Is(cacheErr, modelCache.ErrUserNotFound) {
		log.Warnf("Failed to get user from cache, params: %+v, err: %+v", params, cacheErr)
	} else if cacheErr == nil {
		log.Infof("User retrieved from cache, userID: %d", params.UserID)
		return user, nil
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error

		resp, txErr = s.userRepository.GetUser(ctx, params)
		if txErr != nil {
			return txErr
		}

		txErr = s.logRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Get",
			RequestData:  params,
			ResponseData: resp,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return model.GetUserResponse{}, err
	}

	if errors.Is(cacheErr, modelCache.ErrUserNotFound) {
		cacheErr = s.cacheClient.Create(ctx, resp.User)
		if cacheErr != nil {
			log.Warnf("Failed to create user in cache, params: %+v, err: %+v", resp, cacheErr)
		}
	}

	return resp, nil
}

// UpdateUser updates an existing user's information based on the provided UpdateUserParams.
func (s *userService) UpdateUser(
	ctx context.Context,
	params model.UpdateUserParams,
) (err error) {
	log.Infof("userService.UpdateUser, params: %+v", params)

	var resp model.UpdateUserResponse

	cacheErr := s.cacheClient.Delete(ctx, model.DeleteUserParams{UserID: params.UserID})
	if cacheErr != nil {
		log.Warnf("Failed to delete user from cache before update, UserID: %d, err: %v", params.UserID, cacheErr)
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error

		resp, txErr = s.userRepository.UpdateUser(ctx, params)
		if txErr != nil {
			return txErr
		}

		txErr = s.logRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Update",
			RequestData:  params,
			ResponseData: resp,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	cacheErr = s.cacheClient.Create(ctx, resp.User)
	if cacheErr != nil {
		log.Warnf("Failed to create user in cache, params: %+v, err: %+v", resp, cacheErr)
	}

	return nil
}

// DeleteUser deletes a user based on the provided DeleteUserParams.
func (s *userService) DeleteUser(
	ctx context.Context,
	params model.DeleteUserParams,
) (err error) {
	log.Infof("userService.DeleteUser, params: %+v", params)

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error

		txErr = s.userRepository.DeleteUser(ctx, params)
		if txErr != nil {
			return txErr
		}

		txErr = s.logRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:      "Delete",
			RequestData: params,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	cacheErr := s.cacheClient.Delete(ctx, params)
	if cacheErr != nil {
		log.Warnf("Failed to delete user from cache, params: %+v, err: %v", params, cacheErr)
	}

	return nil
}
