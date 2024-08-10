package tests

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Prrromanssss/platform_common/pkg/db"
	dbMocks "github.com/Prrromanssss/platform_common/pkg/db/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/repository"
	repositoryMocks "github.com/Prrromanssss/auth/internal/repository/mocks"
	userService "github.com/Prrromanssss/auth/internal/service/user"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type (
		userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
		txManagerMockFunc      func(f func(context.Context) error, mc *minimock.Controller) db.TxManager
	)

	type args struct {
		ctx context.Context
		req model.UpdateUserParams
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		role = pb.Role_ADMIN
		name = gofakeit.Name()

		ErrRepository    = errors.New("repository error")
		ErrRepositoryLog = errors.New("repository error in log")

		req = model.UpdateUserParams{
			UserID: id,
			Role:   int64(role),
			Name:   &name,
		}
	)

	requestData, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	logApiReq := model.CreateAPILogParams{
		Method:      "Update",
		RequestData: string(requestData),
	}

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
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
				mock.UpdateUserMock.Expect(ctx, req).Return(nil)
				mock.CreateAPILogMock.Expect(ctx, logApiReq).Return(nil)

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
			name: "repository error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrRepository,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(ErrRepository)

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
			name: "repository error in log",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: ErrRepositoryLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.UpdateUserMock.Expect(ctx, req).Return(nil)
				mock.CreateAPILogMock.Expect(ctx, logApiReq).Return(ErrRepositoryLog)

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
			txManagerMock := tt.txManagerMock(func(ctx context.Context) error {
				txErr := userRepositoryMock.UpdateUser(ctx, req)
				if txErr != nil {
					return txErr
				}

				requestData, err := json.Marshal(req)
				if err != nil {
					t.Error(err)
				}

				txErr = userRepositoryMock.CreateAPILog(ctx, model.CreateAPILogParams{
					Method:      "Update",
					RequestData: string(requestData),
				})

				if txErr != nil {
					return txErr
				}

				return nil
			}, mc)
			service := userService.NewService(userRepositoryMock, txManagerMock)

			err := service.UpdateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
