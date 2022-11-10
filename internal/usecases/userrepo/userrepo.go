package userrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/sanyarise/hezzl/internal/pb"
)

// UserStore is an interface to store user
type UserStore interface {
	SaveUser(ctx context.Context, user *pb.User) error
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) (chan *pb.User, error)
}

type UserStorage struct {
	userStore UserStore
}

func NewUserStorage(userStore UserStore) *UserStorage {
	log.Println("Enter in NewUserStorage")
	return &UserStorage{
		userStore: userStore,
	}
}

func (us UserStorage) SaveUser(ctx context.Context, user *pb.User) error {
	log.Println("enter in repo CreateUser")
	err := us.userStore.SaveUser(ctx, user)
	if err != nil {
		return fmt.Errorf("error on create new user: %w", err)
	}
	return nil
}

func (us UserStorage) DeleteUser(ctx context.Context, id string) error {
	log.Println("Enter in repo DeleteUser")
	err := us.userStore.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("error on delete user by id: %w", err)
	}
	return nil
}

func (us UserStorage) GetAllUsers(ctx context.Context) (chan *pb.User, error) {
	log.Println("Enter in repo GetAllUsers")
	res, err := us.userStore.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on GetAllUsers: %W", err)
	}
	return res, nil
}
