package app

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Prrromanssss/platform_common/pkg/closer"
	"github.com/gofiber/fiber/v2/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Prrromanssss/auth/config"
	"github.com/Prrromanssss/auth/internal/interceptor"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
	_ "github.com/Prrromanssss/auth/statik"
)

type App struct {
	cfg             *config.Config
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
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

	wg := &sync.WaitGroup{}

	wg.Add(4)

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

	// Starting Swagger server.
	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Panic(err)
		}
	}()

	// Starting Kafka consumer.
	go func() {
		defer wg.Done()

		err := a.serviceProvider.UserSaverConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			log.Panicf("failed to run consumer: %s", err.Error())
		}
	}()

	// Handle graceful shutdown.
	a.gracefulShutdown(ctx, cancel, wg)

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
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
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

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

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Addr:              a.cfg.HTTP.Address(),
		Handler:           corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(ctx context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Addr:              a.cfg.Swagger.Address(),
		Handler:           mux,
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

func (a *App) runSwaggerServer() error {
	log.Infof("Starting Swagger server on %s", a.cfg.Swagger.Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Infof("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			err = file.Close()
			if err != nil {
				log.Warnf("Cannot close file: %+v", err)
			}
		}()

		log.Infof("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Infof("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Infof("Served swagger file: %s", path)
	}
}

func (a *App) gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Info("terminating: context cancelled")
	case <-waitSignal():
		log.Info("terminating: via signal")
	}

	a.grpcServer.GracefulStop()
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		log.Panicf("cannot shutdown http server: %+v", err)
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}
