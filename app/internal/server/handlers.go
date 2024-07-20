package server

import (
	"context"

	deliveryGRPC "github.com/Prrromanssss/auth/internal/delivery/grpc"
	"github.com/Prrromanssss/auth/internal/repository/user"

	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func (s *Server) MapHandlers(ctx context.Context) error {
	// repos

	userRepo := user.NewPGRepo(s.pgDB)
	// handlers
	GRPCHandlers := deliveryGRPC.NewGRPCHandlers(userRepo)

	pb.RegisterUserV1Server(s.grpc, GRPCHandlers)

	return nil
}
