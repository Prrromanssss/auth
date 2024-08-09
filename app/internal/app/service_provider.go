package app

import (
	"context"
	"log"

	"github.com/Prrromanssss/platform_common/pkg/closer"

	"github.com/Prrromanssss/platform_common/pkg/db"
	"github.com/Prrromanssss/platform_common/pkg/db/pg"
	"github.com/Prrromanssss/platform_common/pkg/db/transaction"

	"github.com/Prrromanssss/auth/config"
	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/repository"
	userRepository "github.com/Prrromanssss/auth/internal/repository/user"
	"github.com/Prrromanssss/auth/internal/service"
	userService "github.com/Prrromanssss/auth/internal/service/user"
)

type serviceProvider struct {
	cfg *config.Config

	db        db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	userService    service.UserService
	userAPI        *userAPI.GRPCHandlers
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{
		cfg: cfg,
	}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.db == nil {
		cl, err := pg.New(ctx, s.cfg.Postgres.DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.db = cl
	}

	return s.db
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserAPI(ctx context.Context) *userAPI.GRPCHandlers {
	if s.userAPI == nil {
		s.userAPI = userAPI.NewGRPCHandlers(s.UserService(ctx))
	}

	return s.userAPI
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}
