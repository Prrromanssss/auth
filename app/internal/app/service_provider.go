package app

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/Prrromanssss/auth/config"
	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/repository"
	userRepository "github.com/Prrromanssss/auth/internal/repository/user"
	"github.com/Prrromanssss/auth/internal/service"
	userService "github.com/Prrromanssss/auth/internal/service/user"
	"github.com/Prrromanssss/auth/pkg/closer"
)

type serviceProvider struct {
	cfg *config.Config

	db *sqlx.DB

	userRepository repository.UserRepository
	userService    service.UserService
	userAPI        *userAPI.GRPCHandlers
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{
		cfg: cfg,
	}
}

func (s *serviceProvider) DBClient() *sqlx.DB {
	if s.db == nil {
		db, err := sqlx.Connect("postgres", s.cfg.Postgres.DSN())
		if err != nil {
			log.Panic(err)
		}

		closer.Add(db.Close)

		s.db = db
	}

	return s.db
}

func (s *serviceProvider) UserRepository() repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient())
	}

	return s.userRepository
}

func (s *serviceProvider) UserService() service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(s.UserRepository())
	}

	return s.userService
}

func (s *serviceProvider) UserAPI() *userAPI.GRPCHandlers {
	if s.userAPI == nil {
		s.userAPI = userAPI.NewGRPCHandlers(s.UserService())
	}

	return s.userAPI
}
