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

type Server struct {
	cfg  *config.Config
	grpc *grpc.Server
	pgDB *sqlx.DB
}

func NewServer(
	cfg *config.Config,
	database *sqlx.DB,
) *Server {
	return &Server{
		cfg:  cfg,
		grpc: grpc.NewServer(),
		pgDB: database,
	}
}

func (s *Server) Run(ctx context.Context, cancel context.CancelFunc) error {
	err := s.MapHandlers(ctx)
	if err != nil {
		log.Fatalf("Cannot map handlers: %v", err)
	}

	go func() {
		listener, err := net.Listen(
			"tcp",
			fmt.Sprintf("%s:%s", s.cfg.GRPC.Host, s.cfg.GRPC.Port),
		)
		if err != nil {
			log.Fatalf("Error start listener: %s", err.Error())
		}

		log.Printf("Start GRPC server on port: %s:%s", s.cfg.GRPC.Host, s.cfg.GRPC.Port)
		if err := s.grpc.Serve(listener); err != nil {
			log.Fatalf("Error starting GRPC Server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		log.Println("Context cancelled, start graceful shutdown...")
		s.grpc.GracefulStop()
	case <-quit:
		log.Println("Received signal, start graceful shutdown...")
		s.grpc.GracefulStop()
	}

	log.Println("gRPC server exited properly")
	cancel()

	return nil
}
