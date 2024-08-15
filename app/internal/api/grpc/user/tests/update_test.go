package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/service"
	serviceMocks "github.com/Prrromanssss/auth/internal/service/mocks"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id   = gofakeit.Int64()
		name = gofakeit.Name()
		role = pb.Role_USER

		ErrService = errors.New("service error")

		serviceParamsWithName = model.UpdateUserParams{
			UserID: id,
			Name:   &name,
			Role:   int64(role),
		}

		serviceParamsWithoutName = model.UpdateUserParams{
			UserID: id,
			Name:   nil,
			Role:   int64(role),
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case with name",
			args: args{
				ctx: ctx,
				req: &pb.UpdateRequest{
					Id:   id,
					Name: &wrapperspb.StringValue{Value: name},
					Role: role,
				},
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, serviceParamsWithName).Return(nil)
				return mock
			},
		},
		{
			name: "success case without name",
			args: args{
				ctx: ctx,
				req: &pb.UpdateRequest{
					Id:   id,
					Role: role,
				},
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, serviceParamsWithoutName).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: &pb.UpdateRequest{
					Id:   id,
					Name: &wrapperspb.StringValue{Value: name},
					Role: role,
				},
			},
			want: nil,
			err:  ErrService,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateUserMock.Expect(ctx, serviceParamsWithName).Return(ErrService)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewGRPCHandlers(userServiceMock)

			resp, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
