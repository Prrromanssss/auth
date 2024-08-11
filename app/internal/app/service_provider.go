package app

import (
	"context"
	"log"

	redisCache "github.com/Prrromanssss/platform_common/pkg/cache"
	"github.com/Prrromanssss/platform_common/pkg/cache/redis"
	"github.com/Prrromanssss/platform_common/pkg/closer"
	"github.com/Prrromanssss/platform_common/pkg/db"
	"github.com/Prrromanssss/platform_common/pkg/db/pg"
	"github.com/Prrromanssss/platform_common/pkg/db/transaction"
	redigo "github.com/gomodule/redigo/redis"

	"github.com/Prrromanssss/auth/config"
	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/cache"
	userCache "github.com/Prrromanssss/auth/internal/cache/user"
	"github.com/Prrromanssss/auth/internal/repository"
	logRepository "github.com/Prrromanssss/auth/internal/repository/log"
	userRepository "github.com/Prrromanssss/auth/internal/repository/user"
	"github.com/Prrromanssss/auth/internal/service"
	userService "github.com/Prrromanssss/auth/internal/service/user"
)

type serviceProvider struct {
	cfg *config.Config

	db        db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient redisCache.RedisClient

	userRepository repository.UserRepository
	logRepository  repository.LogRepository

	userCache   cache.UserCache
	userService service.UserService
	userAPI     *userAPI.GRPCHandlers
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
			log.Panicf("failed to create db client: %v, %s", err, s.cfg.Postgres.DSN())
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Panicf("db ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.db = cl
	}

	return s.db
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.cfg.Redis.MaxIdle,
			IdleTimeout: s.cfg.Redis.IdleTimeout,
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.cfg.Redis.Address())
			},
		}

		closer.Add(s.redisPool.Close)
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient(ctx context.Context) redisCache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.cfg.Redis.ConnectionTimeout)

		err := s.redisClient.Ping(ctx)
		if err != nil {
			log.Panicf("redis ping error: %s", err.Error())
		}

	}

	return s.redisClient
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) UserCache(ctx context.Context) cache.UserCache {
	if s.userCache == nil {
		s.userCache = userCache.NewCache(s.RedisClient(ctx))
	}

	return s.userCache
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.LogRepository(ctx),
			s.UserCache(ctx),
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
