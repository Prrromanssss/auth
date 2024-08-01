package converter

import (
	"database/sql"

	"github.com/Prrromanssss/auth/internal/model"
	modelRepo "github.com/Prrromanssss/auth/internal/repository/user/model"
)

func CreateUserParamsFromServiceToRepo(params *model.CreateUserParams) *modelRepo.CreateUserParams {
	return &modelRepo.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
		Role:           params.Role,
	}
}

func CreateUserResponseFromRepoToService(params *modelRepo.CreateUserResponse) *model.CreateUserResponse {
	return &model.CreateUserResponse{
		UserID: params.UserID,
	}
}

func GetUserParamsFromServiceToRepo(params *model.GetUserParams) *modelRepo.GetUserParams {
	return &modelRepo.GetUserParams{
		UserID: params.UserID,
	}
}

func GetUserResponseFromRepoToService(params *modelRepo.GetUserResponse) *model.GetUserResponse {
	return &model.GetUserResponse{
		UserID:    params.UserID,
		Name:      params.Name,
		Email:     params.Email,
		Role:      params.Role,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	}
}

func UpdateUserParamsFromServiceToRepo(params *model.UpdateUserParams) *modelRepo.UpdateUserParams {
	var name sql.NullString
	if params.Name != nil {
		name.String = *params.Name
		name.Valid = true
	} else {
		name.String = ""
		name.Valid = false
	}

	return &modelRepo.UpdateUserParams{
		UserID: params.UserID,
		Name:   name,
		Role:   params.Role,
	}
}

func DeleteUserParamsFromServiceToRepo(params *model.DeleteUserParams) *modelRepo.DeleteUserParams {
	return &modelRepo.DeleteUserParams{
		UserID: params.UserID,
	}
}
