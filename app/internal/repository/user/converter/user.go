package converter

import (
	"database/sql"

	"github.com/Prrromanssss/auth/internal/model"
	modelRepo "github.com/Prrromanssss/auth/internal/repository/user/model"
)

// CreateUserParamsFromServiceToRepo converts CreateUserParams from the service layer to the repository layer.
func CreateUserParamsFromServiceToRepo(params *model.CreateUserParams) *modelRepo.CreateUserParams {
	return &modelRepo.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
		Role:           params.Role,
	}
}

// CreateUserResponseFromRepoToService converts CreateUserResponse from the repository layer to the service layer.
func CreateUserResponseFromRepoToService(params *modelRepo.CreateUserResponse) *model.CreateUserResponse {
	return &model.CreateUserResponse{
		UserID: params.UserID,
	}
}

// GetUserParamsFromServiceToRepo converts GetUserParams from the service layer to the repository layer.
func GetUserParamsFromServiceToRepo(params *model.GetUserParams) *modelRepo.GetUserParams {
	return &modelRepo.GetUserParams{
		UserID: params.UserID,
	}
}

// GetUserResponseFromRepoToService converts GetUserResponse from the repository layer to the service layer.
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

// UpdateUserParamsFromServiceToRepo converts UpdateUserParams from the service layer to the repository layer.
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

// DeleteUserParamsFromServiceToRepo converts DeleteUserParams from the service layer to the repository layer.
func DeleteUserParamsFromServiceToRepo(params *model.DeleteUserParams) *modelRepo.DeleteUserParams {
	return &modelRepo.DeleteUserParams{
		UserID: params.UserID,
	}
}
