package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"

	userAPI "github.com/Prrromanssss/auth/internal/api/grpc/user"
	"github.com/Prrromanssss/auth/internal/model"
	"github.com/Prrromanssss/auth/internal/service"
	serviceMocks "github.com/Prrromanssss/auth/internal/service/mocks"
	pb "github.com/Prrromanssss/auth/pkg/user_v1"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *pb.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = gofakeit.Int64()

		ErrService = errors.New("service error")

		serviceParams = model.DeleteUserParams{
			UserID: id,
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
			name: "success case",
			args: args{
				ctx: ctx,
				req: &pb.DeleteRequest{
					Id: id,
				},
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteUserMock.Expect(ctx, serviceParams).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: &pb.DeleteRequest{
					Id: id,
				},
			},
			want: nil,
			err:  ErrService,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteUserMock.Expect(ctx, serviceParams).Return(ErrService)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userServiceMock := tt.userServiceMock(mc)
			api := userAPI.NewGRPCHandlers(userServiceMock)

			resp, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
