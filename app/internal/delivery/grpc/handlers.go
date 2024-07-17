package grpc

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

type GRPCHandlers struct {
	pb.UnimplementedUserV1Server
}

func NewGRPCHandlers() pb.UserV1Server {
	return &GRPCHandlers{}
}

func (h *GRPCHandlers) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	log.Printf("rpc Create, request: %+v", req)

	return &pb.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (h *GRPCHandlers) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("rpc Get, request: %+v", req)

	return &pb.GetResponse{
		Id:        gofakeit.Int64(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      1,
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

func (h *GRPCHandlers) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Update, request: %+v", req)

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandlers) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("rpc Delete, request: %+v", req)

	return &emptypb.Empty{}, nil
}
