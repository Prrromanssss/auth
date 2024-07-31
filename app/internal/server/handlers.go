package server

import (
	"context"

	apiGRPC "github.com/Prrromanssss/auth/internal/api/grpc"
	"github.com/Prrromanssss/auth/internal/repository/user"

	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

// MapHandlers initializes the repository and registers gRPC handlers with the server.
func (s *Server) MapHandlers(ctx context.Context) error {
	// Initialize repository
	userRepo := user.NewPGRepo(s.pgDB)

	// Create and register gRPC handlers
	GRPCHandlers := apiGRPC.NewGRPCHandlers(userRepo)
	pb.RegisterUserV1Server(s.grpc, GRPCHandlers)

	return nil
}
