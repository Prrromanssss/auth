package grpc

import (
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

type GRPCHandlers struct {
	pb.UnimplementedUserV1Server
}

func NewGRPCHandlers() pb.UserV1Server {
	return &GRPCHandlers{}
}
