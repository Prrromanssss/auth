package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/guregu/null"

	"github.com/Prrromanssss/auth/internal/models"
	"github.com/Prrromanssss/auth/internal/repository"
	"github.com/Prrromanssss/auth/pkg/crypto"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

type GRPCHandlers struct {
	pb.UnimplementedUserV1Server
	repo repository.UserRepository
}

func NewGRPCHandlers(repo repository.UserRepository) pb.UserV1Server {
	return &GRPCHandlers{repo: repo}
}

func (h *GRPCHandlers) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Printf("rpc Create, request: %+v", req)

	if req.Password != req.PasswordConfirm {
		return nil, errors.New("passwords don't match")
	}

	hashPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

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

func (h *GRPCHandlers) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("rpc Get, request: %+v", req)

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

func (h *GRPCHandlers) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Update, request: %+v", req)

	err := h.repo.UpdateUser(ctx, models.UpdateUserParams{
		UserID: req.Id,
		Name:   null.StringFrom(req.Name.Value),
		Role:   int64(req.Role),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandlers) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Delete, request: %+v", req)

	err := h.repo.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
