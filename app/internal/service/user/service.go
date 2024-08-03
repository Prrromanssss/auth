package user

import (
	"context"

	"github.com/gofiber/fiber/v2/log"

	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/internal/service"
)

type userService struct {
	userRepository repository.UserRepository
}

// NewService creates a new instance of userService with the provided UserRepository.
func NewService(userRepository repository.UserRepository) service.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

// CreateUser creates a new user using the provided CreateUserParams.
func (s *userService) CreateUser(
	ctx context.Context,
	params model.CreateUserParams,
) (resp model.CreateUserResponse, err error) {
	log.Infof("userService.CreateUser, params: %+v", params)

	resp, err = s.userRepository.CreateUser(ctx, params)
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

	resp, err = s.userRepository.GetUser(ctx, params)
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

	err = s.userRepository.UpdateUser(ctx, params)
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

	err = s.userRepository.DeleteUser(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
