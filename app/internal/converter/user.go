package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Prrromanssss/auth/internal/model"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

// ConvertCreateRequestFromHandlerToService converts a gRPC CreateRequest to a CreateUserParams model used by the service layer.
func ConvertCreateRequestFromHandlerToService(params *pb.CreateRequest) *model.CreateUserParams {
	return &model.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: params.Password,
		Role:           int64(params.Role),
	}
}

// ConvertCreateUserResponseFromServiceToHandler converts a CreateUserResponse model from the service layer to a gRPC CreateResponse.
func ConvertCreateUserResponseFromServiceToHandler(params *model.CreateUserResponse) *pb.CreateResponse {
	return &pb.CreateResponse{
		Id: params.UserID,
	}
}

// ConvertGetRequestFromHandlerToService converts a gRPC GetRequest to a GetUserParams model used by the service layer.
func ConvertGetRequestFromHandlerToService(params *pb.GetRequest) *model.GetUserParams {
	return &model.GetUserParams{
		UserID: params.Id,
	}
}

// ConvertGetUserResponseFromHandlerToService converts a GetUserResponse model from the service layer to a gRPC GetResponse.
func ConvertGetUserResponseFromHandlerToService(params *model.GetUserResponse) *pb.GetResponse {
	return &pb.GetResponse{
		Id:        params.UserID,
		Name:      params.Name,
		Email:     params.Email,
		Role:      pb.Role(params.Role),
		CreatedAt: timestamppb.New(params.CreatedAt),
		UpdatedAt: timestamppb.New(params.UpdatedAt),
	}
}

// ConvertUpdateRequestFromHandlerToService converts a gRPC UpdateRequest to an UpdateUserParams model used by the service layer.
func ConvertUpdateRequestFromHandlerToService(params *pb.UpdateRequest) *model.UpdateUserParams {
	var name *string

	if params.Name != nil {
		name = &params.Name.Value
	}

	return &model.UpdateUserParams{
		UserID: params.Id,
		Name:   name,
		Role:   int64(params.Role),
	}
}

// ConvertDeleteRequestFromHandlerToService converts a gRPC DeleteRequest to a DeleteUserParams model used by the service layer.
func ConvertDeleteRequestFromHandlerToService(params *pb.DeleteRequest) *model.DeleteUserParams {
	return &model.DeleteUserParams{
		UserID: params.Id,
	}
}
