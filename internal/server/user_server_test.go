package server_test

import (
	"context"
	"testing"

	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/server"
	"github.com/sanyarise/hezzl/internal/usecases/cashrepo"
	"github.com/sanyarise/hezzl/internal/usecases/qrepo"
	"github.com/sanyarise/hezzl/internal/usecases/userrepo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"
)

func TestServerCreateUser(t *testing.T) {
	t.Parallel()

	userNoID := pb.User{Id: "", Name: "TestName"}

	userInvalidID := pb.User{Id: "InvalidUUID", Name: "TestName"}

	userStore := userrepo.NewCustomMockStore()
	userCash := cashrepo.NewCustomMockCash()
	userQueue := qrepo.NewCustomMockQueue()
	store := userrepo.NewUserStorage(userStore)
	cash := cashrepo.NewCashStore(userCash)
	queue := qrepo.NewUserQueue(userQueue)

	testCases := []struct {
		name  string
		user  *pb.User
		store userrepo.UserStore
		cash  cashrepo.CashStore
		queue qrepo.UserQueue
		code  codes.Code
	}{
		{
			name:  "success_with_id",
			user:  &pb.User{Id: "d8ea77b6-ca89-478d-8e19-cd1f9666f209", Name: "John Dou"},
			store: *store,
			cash:  *cash,
			queue: *queue,
			code:  codes.OK,
		},
		{
			name:  "success_no_id",
			user:  &userNoID,
			store: *store,
			cash:  *cash,
			queue: *queue,
			code:  codes.OK,
		},
		{
			name:  "failure_invalid_id",
			user:  &userInvalidID,
			store: *store,
			cash:  *cash,
			queue: *queue,
			code:  codes.InvalidArgument,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := &pb.CreateUserRequest{
				User: tc.user,
			}

			server := server.NewUserServer(tc.store, &tc.cash, &tc.queue)
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

func TestServerDeleteUser(t *testing.T) {
	t.Parallel()

	InvalidID := "TestName"
	ValidID := "d8ea77b6-ca89-478d-8e19-cd1f9666f209"

	userStore := userrepo.NewCustomMockStore()
	userCash := cashrepo.NewCustomMockCash()
	userQueue := qrepo.NewCustomMockQueue()
	store := userrepo.NewUserStorage(userStore)
	cash := cashrepo.NewCashStore(userCash)
	queue := qrepo.NewUserQueue(userQueue)

	testCases := []struct {
		name  string
		id    string
		store userrepo.UserStore
		cash  cashrepo.CashStore
		queue qrepo.UserQueue
		code  codes.Code
	}{
		{
			name:  "success_with_id",
			id:    ValidID,
			store: *store,
			cash:  *cash,
			queue: *queue,
			code:  codes.OK,
		},
		{
			name:  "failure_invalid_id",
			id:    InvalidID,
			store: *store,
			cash:  *cash,
			queue: *queue,
			code:  codes.InvalidArgument,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := &pb.DeleteUserRequest{
				Id: tc.id,
			}

			server := server.NewUserServer(tc.store, &tc.cash, &tc.queue)
			res, err := server.DeleteUser(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
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

func TestServerGetAllUsers(t *testing.T) {
	userStore := userrepo.NewCustomMockStore()
	userCash := cashrepo.NewCustomMockCash()
	userQueue := qrepo.NewCustomMockQueue()
	store := userrepo.NewUserStorage(userStore)
	cash := cashrepo.NewCashStore(userCash)
	queue := qrepo.NewUserQueue(userQueue)
	server := server.NewUserServer(store, cash, queue)
	stream := newStreamMock()
	req := &pb.AllUsersRequest{}
	err := server.GetAllUsers(req, stream)
	require.Nil(t, err)
}

func newStreamMock() *StreamMock {
	return &StreamMock{
		ctx:            context.Background(),
		sentFromServer: make(chan pb.AllUsersResponse, 10),
	}
}

type StreamMock struct {
	grpc.ServerStream
	ctx            context.Context
	sentFromServer chan pb.AllUsersResponse
}

func (m *StreamMock) Send(resp *pb.AllUsersResponse) error {
	m.sentFromServer <- *resp
	return nil
}

func (m *StreamMock) Context() context.Context {
	return m.ctx
}
