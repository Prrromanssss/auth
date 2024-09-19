package converter

import "github.com/Prrromanssss/auth/internal/model"

// ConvertCreateRequestFromHandlerToService converts a CreateRequest to a CreateUserParams the api layer to the service layer.
func ConvertCreateUserKafkaParamsFromConsumerServiceToUserService(params model.CreateUserKafkaParams) model.CreateUserParams {
	return model.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.Password,
		Role:           int64(params.Role),
	}
}
