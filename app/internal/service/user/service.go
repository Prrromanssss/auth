package user

import (
	"context"

	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/internal/service"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewService(
	userRepository repository.UserRepository,
) service.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(
	ctx context.Context,
	params *model.CreateUserParams,
) (resp *model.CreateUserResponse, err error) {
	resp, err = s.userRepository.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *userService) GetUser(
	ctx context.Context,
	params *model.GetUserParams,
) (resp *model.GetUserResponse, err error) {
	resp, err = s.userRepository.GetUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *userService) UpdateUser(
	ctx context.Context,
	params *model.UpdateUserParams,
) (err error) {
	err = s.userRepository.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(
	ctx context.Context,
	params *model.DeleteUserParams,
) (err error) {
	err = s.userRepository.DeleteUser(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
