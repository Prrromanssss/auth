package converter

import (
	"time"

	modelCache "github.com/Prrromanssss/auth/internal/cache/user/model"
	"github.com/Prrromanssss/auth/internal/model"
)

// ConvertUserFromServiceToCache converts a User model from the service layer to the cache layer.
func ConvertUserFromServiceToCache(params model.User) modelCache.User {
	return modelCache.User{
		UserID:    params.UserID,
		Name:      params.Name,
		Email:     params.Email,
		Role:      params.Role,
		CreatedAt: params.CreatedAt.Unix(),
		UpdatedAt: params.UpdatedAt.Unix(),
	}
}

// ConvertGetUserParamsFromServiceToCache converts GetUserParams from the service layer to the cache layer.
func ConvertGetUserParamsFromServiceToCache(params model.GetUserParams) modelCache.GetUserParams {
	return modelCache.GetUserParams(params)
}

// ConvertDeleteUserParamsFromServiceToCache converts DeleteUserParams from the service layer to the cache layer.
func ConvertDeleteUserParamsFromServiceToCache(params model.DeleteUserParams) modelCache.DeleteUserParams {
	return modelCache.DeleteUserParams(params)
}

// ConvertGetUserResponseFromCacheToService converts a User model from the cache layer to the service layer.
func ConvertGetUserResponseFromCacheToService(params modelCache.User) model.GetUserResponse {
	return model.GetUserResponse{
		User: model.User{
			UserID:    params.UserID,
			Name:      params.Name,
			Email:     params.Email,
			Role:      params.Role,
			CreatedAt: time.UnixMilli(params.CreatedAt),
			UpdatedAt: time.UnixMilli(params.UpdatedAt),
		},
	}
}
