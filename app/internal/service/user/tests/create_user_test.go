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
	"github.com/Prrromanssss/auth/pkg/crypto"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type (
		userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
		txManagerMockFunc      func(f func(context.Context) error, mc *minimock.Controller) db.TxManager
	)

	type args struct {
		ctx context.Context
		req model.CreateUserParams
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id             = gofakeit.Int64()
		name           = gofakeit.Name()
		email          = gofakeit.Email()
		role           = pb.Role_ADMIN
		password       = gofakeit.Password(true, true, true, true, true, 10)
		hashedPassword = crypto.HashPassword(password)

		ErrRepository    = errors.New("repository error")
		ErrRepositoryLog = errors.New("repository error in log")

		req = model.CreateUserParams{
			Name:           name,
			Email:          email,
			Role:           int64(role),
			HashedPassword: hashedPassword,
		}

		resp = model.CreateUserResponse{
			UserID: id,
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
		Method:       "Create",
		RequestData:  string(requestData),
		ResponseData: &responseDataString,
	}

	tests := []struct {
		name               string
		args               args
		want               model.CreateUserResponse
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
				mock.CreateUserMock.Expect(ctx, req).Return(resp, nil)
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
			want: model.CreateUserResponse{},
			err:  ErrRepository,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(model.CreateUserResponse{}, ErrRepository)

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
			want: model.CreateUserResponse{},
			err:  ErrRepositoryLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(mc)
				mock.CreateUserMock.Expect(ctx, req).Return(resp, nil)
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
				resp, txErr := userRepositoryMock.CreateUser(ctx, req)
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

			resp, err := service.CreateUser(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
