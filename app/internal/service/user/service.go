package user

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"

	"github.com/Prrromanssss/auth/internal/client/db"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/internal/service"
)

type userService struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

// NewService creates a new instance of userService with the provided UserRepository.
func NewService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
) service.UserService {
	return &userService{
		userRepository: userRepository,
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

		requestData, txErr := json.Marshal(params)
		if txErr != nil {
			return txErr
		}

		responseData, txErr := json.Marshal(resp)
		if txErr != nil {
			return txErr
		}

		responseDataString := string(responseData)

		txErr = s.userRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Create",
			RequestData:  string(requestData),
			ResponseData: &responseDataString,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return model.CreateUserResponse{}, err
	}

	return resp, nil
}

// GetUser retrieves a user's information based on the provided GetUserParams.
func (s *userService) GetUser(
	ctx context.Context,
	params model.GetUserParams,
) (resp model.GetUserResponse, err error) {
	log.Infof("userService.GetUser, params: %+v", params)

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error

		resp, txErr = s.userRepository.GetUser(ctx, params)
		if err != nil {
			return txErr
		}

		requestData, txErr := json.Marshal(params)
		if txErr != nil {
			return txErr
		}

		responseData, txErr := json.Marshal(resp)
		if txErr != nil {
			return txErr
		}

		responseDataString := string(responseData)

		txErr = s.userRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Get",
			RequestData:  string(requestData),
			ResponseData: &responseDataString,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return model.GetUserResponse{}, err
	}

	return resp, nil
}

// UpdateUser updates an existing user's information based on the provided UpdateUserParams.
func (s *userService) UpdateUser(
	ctx context.Context,
	params model.UpdateUserParams,
) (err error) {
	log.Infof("userService.UpdateUser, params: %+v", params)

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var txErr error
		var responseData *string

		txErr = s.userRepository.UpdateUser(ctx, params)
		if txErr != nil {
			return txErr
		}

		requestData, txErr := json.Marshal(params)
		if txErr != nil {
			return txErr
		}

		txErr = s.userRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Update",
			RequestData:  string(requestData),
			ResponseData: responseData,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
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
		var responseData *string

		txErr = s.userRepository.DeleteUser(ctx, params)
		if txErr != nil {
			return txErr
		}

		requestData, txErr := json.Marshal(params)
		if txErr != nil {
			return txErr
		}

		txErr = s.userRepository.CreateAPILog(ctx, model.CreateAPILogParams{
			Method:       "Delete",
			RequestData:  string(requestData),
			ResponseData: responseData,
		})
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
