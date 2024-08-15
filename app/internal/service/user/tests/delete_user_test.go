package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/Prrromanssss/platform_common/pkg/db"
	dbMocks "github.com/Prrromanssss/platform_common/pkg/db/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Prrromanssss/auth/internal/cache"
	cacheMocks "github.com/Prrromanssss/auth/internal/cache/mocks"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	repositoryMocks "github.com/Prrromanssss/auth/internal/repository/mocks"
	userService "github.com/Prrromanssss/auth/internal/service/user"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type (
		userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
		logRepositoryMockFunc  func(mc *minimock.Controller) repository.LogRepository
		txManagerMockFunc      func(f func(context.Context) error, mc *minimock.Controller) db.TxManager
		cacheMockFunc          func(mc *minimock.Controller) cache.UserCache
	)

	type args struct {
		ctx context.Context
		req model.DeleteUserParams
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		ErrUserRepository = errors.New("user repository error")
		ErrLogRepository  = errors.New("log repository error")
		ErrCache          = errors.New("cache error")

		req = model.DeleteUserParams{
			UserID: id,
		}

		logApiReq = model.CreateAPILogParams{
			Method:      "Delete",
			RequestData: req,
		}
	)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
		cacheMock          cacheMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(nil)

				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)
				mock.CreateAPILogMock.Expect(ctx, logApiReq).Return(nil)

				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.DeleteMock.Expect(ctx, req).Return(nil)

				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Optional().Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return mock
			},
		},
		{
			name: "user repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrUserRepository,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(ErrUserRepository)

				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)

				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)

				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Optional().Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return mock
			},
		},
		{
			name: "log repository error",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrLogRepository,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(nil)

				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)
				mock.CreateAPILogMock.Expect(ctx, logApiReq).Return(ErrLogRepository)

				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)

				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Optional().Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return mock
			},
		},
		{
			name: "cache error",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.DeleteUserMock.Expect(ctx, req).Return(nil)

				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) repository.LogRepository {
				mock := repositoryMocks.NewLogRepositoryMock(mc)
				mock.CreateAPILogMock.Expect(ctx, logApiReq).Return(nil)

				return mock
			},
			cacheMock: func(mc *minimock.Controller) cache.UserCache {
				mock := cacheMocks.NewUserCacheMock(mc)
				mock.DeleteMock.Expect(ctx, req).Return(ErrCache)

				return mock
			},
			txManagerMock: func(f func(context.Context) error, mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Optional().Set(func(ctx context.Context, f db.Handler) (err error) {
					return f(ctx)
				})

				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepositoryMock := tt.userRepositoryMock(mc)
			logRepositoryMock := tt.logRepositoryMock(mc)
			cacheMock := tt.cacheMock(mc)
			txManagerMock := tt.txManagerMock(func(ctx context.Context) error {
				txErr := userRepositoryMock.DeleteUser(ctx, req)
				if txErr != nil {
					return txErr
				}

				txErr = logRepositoryMock.CreateAPILog(ctx, model.CreateAPILogParams{
					Method:      "Delete",
					RequestData: req,
				})
				if txErr != nil {
					return txErr
				}

				return nil
			}, mc)

			service := userService.NewService(userRepositoryMock, logRepositoryMock, cacheMock, txManagerMock)

			err := service.DeleteUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
