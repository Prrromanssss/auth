package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Prrromanssss/auth/internal/converter"
	"github.com/Prrromanssss/auth/internal/service"
	"github.com/Prrromanssss/auth/pkg/crypto"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

// GRPCHandlers represents the gRPC handlers that implement the UserV1Server interface
// and use the UserService for business logic operations.
type GRPCHandlers struct {
	pb.UnimplementedUserV1Server
	userService service.UserService
}

// NewGRPCHandlers creates a new instance of GRPCHandlers with the provided UserService.
func NewGRPCHandlers(userService service.UserService) pb.UserV1Server {
	return &GRPCHandlers{userService: userService}
}

// Create handles the request for creating a new user.
func (h *GRPCHandlers) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Printf("rpc Create, request: %+v", req)

	if req.Password != req.PasswordConfirm {
		return nil, errors.New("passwords don't match")
	}

	hashPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashPassword

	resp, err := h.userService.CreateUser(ctx, converter.ConvertCreateRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return converter.ConvertCreateUserResponseFromServiceToHandler(resp), nil
}

// Get handles the request for retrieving user data.
func (h *GRPCHandlers) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("rpc Get, request: %+v", req)

	resp, err := h.userService.GetUser(ctx, converter.ConvertGetRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return converter.ConvertGetUserResponseFromHandlerToService(resp), nil
}

// Update handles the request for updating user data.
func (h *GRPCHandlers) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Update, request: %+v", req)

	err := h.userService.UpdateUser(ctx, converter.ConvertUpdateRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Delete handles the request for deleting a user.
func (h *GRPCHandlers) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Delete, request: %+v", req)

	err := h.userService.DeleteUser(ctx, converter.ConvertDeleteRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
