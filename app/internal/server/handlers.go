package server

import (
	"context"

	deliveryGRPC "github.com/Prrromanssss/auth/internal/delivery/grpc"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func (s *Server) MapHandlers(ctx context.Context) error {
	GRPCHandlers := deliveryGRPC.NewGRPCHandlers()
	pb.RegisterUserV1Server(s.grpc, GRPCHandlers)

	return nil
}
