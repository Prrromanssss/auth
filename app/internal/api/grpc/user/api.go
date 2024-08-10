package user

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2/log"
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
func NewGRPCHandlers(userService service.UserService) *GRPCHandlers {
	return &GRPCHandlers{
		userService: userService,
	}
}

// Create handles the request for creating a new user.
func (h *GRPCHandlers) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Infof("rpc Create, request: %+v", req)

	log.Info(h.userService)

	if req.Password != req.PasswordConfirm {
		return nil, errors.New("passwords don't match")
	}

	hashedPassword := crypto.HashPassword(req.Password)

	req.Password = hashedPassword

	resp, err := h.userService.CreateUser(ctx, converter.ConvertCreateRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return converter.ConvertCreateUserResponseFromServiceToHandler(resp), nil
}

// Get handles the request for retrieving user data.
func (h *GRPCHandlers) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Infof("rpc Get, request: %+v", req)

	resp, err := h.userService.GetUser(ctx, converter.ConvertGetRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return converter.ConvertGetUserResponseFromHandlerToService(resp), nil
}

// Update handles the request for updating user data.
func (h *GRPCHandlers) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Infof("rpc Update, request: %+v", req)

	err := h.userService.UpdateUser(ctx, converter.ConvertUpdateRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Delete handles the request for deleting a user.
func (h *GRPCHandlers) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Infof("rpc Delete, request: %+v", req)

	err := h.userService.DeleteUser(ctx, converter.ConvertDeleteRequestFromHandlerToService(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
