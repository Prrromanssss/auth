package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Prrromanssss/platform_common/pkg/closer"
	"github.com/gofiber/fiber/v2/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Prrromanssss/auth/config"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

type App struct {
	cfg             *config.Config
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

// NewApp creates a new instance of App.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run runs the App.
func (a *App) Run(ctx context.Context, cancel context.CancelFunc) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}

	wg.Add(2)

	// Starting gRPC server.
	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Panic(err)
		}
	}()

	// Starting HTTP server.
	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Panic(err)
		}
	}()

	// Handle graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-ctx.Done():
		log.Info("Context cancelled, initiating graceful shutdown...")
		a.grpcServer.GracefulStop()
		a.httpServer.Shutdown(ctx)
	case <-quit:
		log.Info("Received termination signal, initiating graceful shutdown...")
		a.grpcServer.GracefulStop()
		a.httpServer.Shutdown(ctx)
	}

	log.Info("App shut down gracefully")
	cancel()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	log.Info("Config loaded")

	a.cfg = cfg

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.cfg)

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	pb.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserAPI(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := pb.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.cfg.GRPC.Address(), opts)
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.cfg.HTTP.Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	listener, err := net.Listen(
		"tcp",
		a.cfg.GRPC.Address(),
	)
	if err != nil {
		return errors.Wrapf(err, "Error starting listener")
	}

	log.Infof("Starting gRPC server on %s", a.cfg.GRPC.Address())

	if err := a.grpcServer.Serve(listener); err != nil {
		return errors.Wrapf(err, "Error starting gRPC server")
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Infof("Starting HTTP server on %s", a.cfg.HTTP.Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
