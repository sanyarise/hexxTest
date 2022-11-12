package cashrepo

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sanyarise/hezzl/internal/pb"
)

/*type Cash interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res chan *pb.User, key string) error
	GetCash(key string) ([]*pb.User, error)
}*/

func TestCheckCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cash := NewMockCash(ctrl)
	cashStore := NewCashStore(cash)

	key := "test key"

	cash.EXPECT().CheckCash(key).Return(false)

	ok := cashStore.CheckCash(key)

	if ok {
		t.Errorf("failed to check cash: expected %t, got %t", false, ok)
	}

	cash.EXPECT().CheckCash(key).Return(true)

	ok = cashStore.CheckCash(key)

	if !ok {
		t.Errorf("failed to check cash: expected %t, got %t", true, ok)
	}
}

func TestCreateCash(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cash := NewMockCash(ctrl)
	cashStore := NewCashStore(cash)

	key := "test key"
	testChan := make(chan *pb.User, 100)
	cash.EXPECT().CreateCash(ctx, testChan, key).Return(nil)

	err := cashStore.CreateCash(ctx, testChan, key)
	if err != nil {
		t.Errorf("failed to create cash: expected nil, got %s", err)
	}
}

func TestGetCash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cash := NewMockCash(ctrl)
	cashStore := NewCashStore(cash)

	key := "test key"
	testSlice := make([]*pb.User, 100)

	cash.EXPECT().GetCash(key).Return(testSlice, nil)

	_, err := cashStore.GetCash(key)
	if err != nil {
		t.Errorf("failed to get cash: expected err = nil, got %s", err)
	}

}