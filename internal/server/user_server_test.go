package server_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/server"
	"github.com/sanyarise/hezzl/internal/usecases/userrepo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"
)

func TestServerCreateLaptop(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userNoID := pb.User{Id: "", Name: "TestName"}

	userInvalidID := pb.User{Id: "InvalidUUID", Name: "TestName"}

	testCases := []struct {
		name  string
		user  *pb.User
		store userrepo.UserStore
		code  codes.Code
	}{
		{
			name:  "success_with_id",
			user:  &pb.User{Id: "d8ea77b6-ca89-478d-8e19-cd1f9666f209", Name: "John Dou"},
			store: userrepo.NewMockUserStore(ctrl),
			code:  codes.OK,
		},
		{
			name:   "success_no_id",
			user: &userNoID,
			store:  userrepo.NewMockUserStore(ctrl),
			code:   codes.OK,
		},
		{
			name:   "failure_invalid_id",
			user: &userInvalidID,
			store:  userrepo.NewMockUserStore(ctrl),
			code:   codes.InvalidArgument,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreateUserRequest{
				User: tc.user,
			}

			server := server.NewUserServer(tc.store, nil, nil)
			res, err := server.CreateUser(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}
}
