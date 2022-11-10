package userrepo

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanyarise/hezzl/internal/pb"
)

func TestSaveUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockUserStore(ctrl)
	srv := NewUserStorage(repo)

	user := pb.User{
		Id:   "e428ac20-c89b-459a-b43b-83563836a536",
		Name: "John Dow",
	}
	repo.EXPECT().SaveUser(ctx, &user).Return(nil)

	err := srv.SaveUser(ctx, &user)

	if err != nil {
		t.Errorf("failed to save user: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockUserStore(ctrl)
	srv := NewUserStorage(repo)

	id := "e428ac20-c89b-459a-b43b-83563836a536"
	repo.EXPECT().DeleteUser(ctx, id).Return(nil)

	err := srv.DeleteUser(ctx, id)

	if err != nil {
		t.Errorf("failed to save user: %v", err)
	}
}

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockUserStore(ctrl)
	srv := NewUserStorage(repo)

	res := make(chan *pb.User, 100)

	repo.EXPECT().GetAllUsers(ctx).Return(res, nil)

	_, err := srv.GetAllUsers(ctx)

	if err != nil {
		t.Errorf("failed to save user: %v", err)
	}
}
