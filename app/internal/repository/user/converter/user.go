package converter

import (
	"database/sql"

	"github.com/Prrromanssss/auth/internal/model"
	modelRepo "github.com/Prrromanssss/auth/internal/repository/user/model"
)

// ConvertCreateUserParamsFromServiceToRepo converts CreateUserParams from the service layer to the repository layer.
func ConvertCreateUserParamsFromServiceToRepo(params model.CreateUserParams) modelRepo.CreateUserParams {
	return modelRepo.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.HashedPassword,
		Role:           params.Role,
	}
}

// ConvertCreateUserResponseFromRepoToService converts CreateUserResponse from the repository layer to the service layer.
func ConvertCreateUserResponseFromRepoToService(params modelRepo.CreateUserResponse) model.CreateUserResponse {
	return model.CreateUserResponse{
		User: model.User{
			UserID:    params.UserID,
			Name:      params.Name,
			Email:     params.Email,
			Role:      params.Role,
			CreatedAt: params.CreatedAt,
			UpdatedAt: params.UpdatedAt,
		},
	}
}

// ConvertGetUserParamsFromServiceToRepo converts GetUserParams from the service layer to the repository layer.
func ConvertGetUserParamsFromServiceToRepo(params model.GetUserParams) modelRepo.GetUserParams {
	return modelRepo.GetUserParams{
		UserID: params.UserID,
	}
}

// ConvertGetUserResponseFromRepoToService converts GetUserResponse from the repository layer to the service layer.
func ConvertGetUserResponseFromRepoToService(params modelRepo.GetUserResponse) model.GetUserResponse {
	return model.GetUserResponse{
		User: model.User{
			UserID:    params.UserID,
			Name:      params.Name,
			Email:     params.Email,
			Role:      params.Role,
			CreatedAt: params.CreatedAt,
			UpdatedAt: params.UpdatedAt,
		},
	}
}

// ConvertUpdateUserParamsFromServiceToRepo converts UpdateUserParams from the service layer to the repository layer.
func ConvertUpdateUserParamsFromServiceToRepo(params model.UpdateUserParams) modelRepo.UpdateUserParams {
	var name sql.NullString

	if params.Name != nil {
		name = sql.NullString{String: *params.Name, Valid: true}
	} else {
		name = sql.NullString{String: "", Valid: false}
	}

	return modelRepo.UpdateUserParams{
		UserID: params.UserID,
		Name:   name,
		Role:   params.Role,
	}
}

// ConvertUpdateUserResponseFromRepoToService converts UpdateUserResponse from the repository layer to the service layer.
func ConvertUpdateUserResponseFromRepoToService(params modelRepo.UpdateUserResponse) model.UpdateUserResponse {
	return model.UpdateUserResponse{
		User: model.User{
			UserID:    params.UserID,
			Name:      params.Name,
			Email:     params.Email,
			Role:      params.Role,
			CreatedAt: params.CreatedAt,
			UpdatedAt: params.UpdatedAt,
		},
	}
}

// ConvertDeleteUserParamsFromServiceToRepo converts DeleteUserParams from the service layer to the repository layer.
func ConvertDeleteUserParamsFromServiceToRepo(params model.DeleteUserParams) modelRepo.DeleteUserParams {
	return modelRepo.DeleteUserParams{
		UserID: params.UserID,
	}
}
