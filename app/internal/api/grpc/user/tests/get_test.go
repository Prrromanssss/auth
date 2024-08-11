package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/service"
	serviceMocks "github.com/Prrromanssss/auth/internal/service/mocks"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.GetRequest
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

		ErrService = errors.New("service error")

		req = &pb.GetRequest{
			Id: id,
		}

		resp = &pb.GetResponse{
			Id:        id,
			Name:      name,
			Email:     email,
			Role:      role,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}

		serviceParams = model.GetUserParams{
			UserID: id,
		}

		serviceResp = model.GetUserResponse{
			User: model.User{
				UserID:    id,
				Name:      name,
				Email:     email,
				Role:      int64(role),
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *pb.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: resp,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, serviceParams).Return(serviceResp, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  ErrService,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetUserMock.Expect(ctx, serviceParams).Return(model.GetUserResponse{}, ErrService)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewGRPCHandlers(userServiceMock)

			resp, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
