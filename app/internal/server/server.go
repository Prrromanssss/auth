package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/jmoiron/sqlx"

	"github.com/Prrromanssss/auth/config"
)

// Server holds the configuration, gRPC server instance, and PostgreSQL database connection.
type Server struct {
	cfg  *config.Config // Configuration settings
	grpc *grpc.Server   // gRPC server instance
	pgDB *sqlx.DB       // PostgreSQL database connection
}

// NewServer creates a new Server instance with the provided configuration and database connection.
func NewServer(
	cfg *config.Config, // Configuration for the server
	database *sqlx.DB, // Database connection
) *Server {
	return &Server{
		cfg:  cfg,
		grpc: grpc.NewServer(), // Initialize gRPC server
		pgDB: database,
	}
}

// Run starts the gRPC server, sets up handlers, and manages graceful shutdown on context cancellation or termination signals.
func (s *Server) Run(ctx context.Context, cancel context.CancelFunc) error {
	// Map gRPC handlers to the server
	if err := s.MapHandlers(ctx); err != nil {
		log.Fatalf("Cannot map handlers: %v", err)
	}

	go func() {
		// Start listening for incoming gRPC requests
		listener, err := net.Listen(
			"tcp",
			fmt.Sprintf("%s:%s", s.cfg.GRPC.Host, s.cfg.GRPC.Port),
		)
		if err != nil {
			log.Fatalf("Error starting listener: %s", err.Error())
		}

		log.Printf("Starting gRPC server on %s:%s", s.cfg.GRPC.Host, s.cfg.GRPC.Port)
		if err := s.grpc.Serve(listener); err != nil {
			log.Fatalf("Error starting gRPC server: %s", err.Error())
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		log.Println("Context cancelled, initiating graceful shutdown...")
		s.grpc.GracefulStop()
	case <-quit:
		log.Println("Received termination signal, initiating graceful shutdown...")
		s.grpc.GracefulStop()
	}

	log.Println("gRPC server shut down gracefully")
	cancel()

	return nil
}
