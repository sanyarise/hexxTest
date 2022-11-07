package server

/*import (
	"context"
	"log"
	"testing"

	"github.com/sanyarise/hezzl/internal/cash"
	"github.com/sanyarise/hezzl/internal/db"
	"github.com/sanyarise/hezzl/internal/pb"
	"github.com/sanyarise/hezzl/internal/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"
)

func TestServerCreateUser(t *testing.T) {
	t.Parallel()

	UserNoID := sample.NewUser()
	UserNoID.Id = ""

	UserInvalidID := sample.NewUser()
	UserInvalidID.Id = "invalid-uuid"

	userDuplicateID := sample.NewUser()
	storeDuplicateID, _ := db.NewUserPostgresStore("postgres://postgres:example@localhost:5432/postgres")
	ctx := context.Background()
	err := storeDuplicateID.SaveUser(ctx, userDuplicateID)
	require.Nil(t, err)
	store, _ := db.NewUserPostgresStore("postgres://postgres:example@localhost:5432/postgres")
	testCases := []struct {
		name  string
		user  *pb.User
		store db.UserStore
		code  codes.Code
	}{
		{
			name:  "success_with_id",
			user:  sample.NewUser(),
			store: store,
			code:  codes.OK,
		},
		{
			name:  "success_no_id",
			user:  UserNoID,
			store: store,
			code:  codes.OK,
		},
		{
			name:  "failure_invalid_id",
			user:  UserInvalidID,
			store: store,
			code:  codes.InvalidArgument,
		},
		{
			name:  "failure_duplicate_id",
			user:  userDuplicateID,
			store: storeDuplicateID,
			code:  codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreateUserRequest{
				User: tc.user,
			}
			red, err := cash.NewRedisClient("localhost", "6379", 1)
			if err != nil {
				log.Printf("error on new redis client: %v", err)
			}
			server := NewUserServer(tc.store, red)
			res, err := server.CreateUser(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.user.Id) > 0 {
					require.Equal(t, tc.user.Id, res.Id)
				}
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
*/