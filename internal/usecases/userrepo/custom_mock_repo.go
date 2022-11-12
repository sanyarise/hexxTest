package userrepo

import (
	"context"

	"github.com/sanyarise/hezzl/internal/pb"
)

type CustomMockStore struct {
}

/*type UserStore interface {
	SaveUser(ctx context.Context, user *pb.User) error
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) (chan *pb.User, error)
}*/

func NewCustomMockStore() *CustomMockStore {
	return &CustomMockStore{}
}

func (cms *CustomMockStore) SaveUser(ctx context.Context, user *pb.User) error {
	return nil
}

func (cms *CustomMockStore) DeleteUser(ctx context.Context, id string) error {
	return nil
}

func (cms *CustomMockStore) GetAllUsers(ctx context.Context) (chan *pb.User, error) {
	out := make(chan *pb.User, 100)
	user := pb.User{Id: "TestId", Name: "TestName"}
	out <- &user
	return out, nil
}
