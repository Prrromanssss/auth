package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Prrromanssss/auth/internal/models"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/pkg/crypto"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

// GRPCHandlers implements the gRPC server for user operations using a UserRepository.
type GRPCHandlers struct {
	pb.UnimplementedUserV1Server                           // Embeds the unimplemented server for backward compatibility.
	repo                         repository.UserRepository // User repository for data operations.
}

// NewGRPCHandlers creates a new instance of GRPCHandlers with the provided user repository.
func NewGRPCHandlers(repo repository.UserRepository) pb.UserV1Server {
	return &GRPCHandlers{repo: repo}
}

// Create handles the creation of a new user. It checks if passwords match,
// hashes the password, and saves the user to the repository.
func (h *GRPCHandlers) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Printf("rpc Create, request: %+v", req)

	// Check if the provided passwords match
	if req.Password != req.PasswordConfirm {
		return nil, errors.New("passwords don't match")
	}

	// Hash the password
	hashPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create a new user in the repository
	userID, err := h.repo.CreateUser(ctx, models.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		Role:           int64(req.Role),
		HashedPassword: hashPassword,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Id: userID,
	}, nil
}

// Get retrieves a user by their ID and returns their details.
func (h *GRPCHandlers) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("rpc Get, request: %+v", req)

	// Fetch user details from the repository
	resp, err := h.repo.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		Id:        resp.UserID,
		Name:      resp.Name,
		Email:     resp.Email,
		Role:      pb.Role(resp.Role),
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
	}, nil
}

// Update modifies the details of an existing user.
func (h *GRPCHandlers) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Update, request: %+v", req)

	// Update user details in the repository
	err := h.repo.UpdateUser(ctx, models.UpdateUserParams{
		UserID: req.Id,
		Name:   req.GetName().GetValue(),
		Role:   int64(req.Role),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// Delete removes a user from the repository by their ID.
func (h *GRPCHandlers) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Delete, request: %+v", req)

	// Delete user from the repository
	err := h.repo.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
