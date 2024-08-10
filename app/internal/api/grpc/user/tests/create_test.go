package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/service"
	serviceMocks "github.com/Prrromanssss/auth/internal/service/mocks"
	"github.com/Prrromanssss/auth/pkg/crypto"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.CreateRequest
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

		ErrInvalidPassword = errors.New("passwords don't match")
		ErrService         = errors.New("service error")

		serviceParams = model.CreateUserParams{
			Name:           name,
			Email:          email,
			HashedPassword: hashedPassword,
			Role:           int64(role),
		}

		resp = &pb.CreateResponse{
			Id: id,
		}

		serviceResp = model.CreateUserResponse{
			UserID: id,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *pb.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &pb.CreateRequest{
					Name:            name,
					Email:           email,
					Password:        password,
					PasswordConfirm: password,
					Role:            role,
				},
			},
			want: resp,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, serviceParams).Return(serviceResp, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: &pb.CreateRequest{
					Name:            name,
					Email:           email,
					Password:        password,
					PasswordConfirm: password,
					Role:            role,
				},
			},
			want: nil,
			err:  ErrService,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateUserMock.Expect(ctx, serviceParams).Return(model.CreateUserResponse{}, ErrService)
				return mock
			},
		},
		{
			name: "different passwords",
			args: args{
				ctx: ctx,
				req: &pb.CreateRequest{
					Name:            name,
					Email:           email,
					Password:        password,
					PasswordConfirm: gofakeit.Password(false, true, false, true, true, 9),
					Role:            role,
				},
			},
			want: nil,
			err:  ErrInvalidPassword,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewGRPCHandlers(userServiceMock)

			resp, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
