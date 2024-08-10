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

func TestGet(t *testing.T) {
	t.Parallel()

	type (
		userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
		txManagerMockFunc      func(f func(context.Context) error, mc *minimock.Controller) db.TxManager
	)

	type args struct {
		ctx context.Context
		req model.GetUserParams
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		role      = pb.Role_ADMIN
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		ErrRepository    = errors.New("repository error")
		ErrRepositoryLog = errors.New("repository error in log")

		req = model.GetUserParams{
			UserID: id,
		}

		resp = model.GetUserResponse{
			UserID:    id,
			Name:      name,
			Email:     email,
			Role:      int64(role),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	)

	requestData, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}

	responseData, err := json.Marshal(resp)
	if err != nil {
		t.Error(err)
	}

	responseDataString := string(responseData)

	logApiReq := model.CreateAPILogParams{
		Method:       "Get",
		RequestData:  string(requestData),
		ResponseData: &responseDataString,
	}

	tests := []struct {
		name               string
		args               args
		want               model.GetUserResponse
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
			want: resp,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(resp, nil)
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
			want: model.GetUserResponse{},
			err:  ErrRepository,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(model.GetUserResponse{}, ErrRepository)

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
			want: model.GetUserResponse{},
			err:  ErrRepositoryLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.GetUserMock.Expect(ctx, req).Return(resp, nil)
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
				resp, txErr := userRepositoryMock.GetUser(ctx, req)
				if txErr != nil {
					return txErr
				}

				requestData, txErr := json.Marshal(req)
				if txErr != nil {
					return txErr
				}

				responseData, txErr := json.Marshal(resp)
				if txErr != nil {
					return txErr
				}

				responseDataString := string(responseData)

				txErr = userRepositoryMock.CreateAPILog(ctx, model.CreateAPILogParams{
					Method:       "Create",
					RequestData:  string(requestData),
					ResponseData: &responseDataString,
				})

				if txErr != nil {
					return txErr
				}

				return nil
			}, mc)
			service := userService.NewService(userRepositoryMock, txManagerMock)

			resp, err := service.GetUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
