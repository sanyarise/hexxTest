package cashrepo

import (
	"context"

	"github.com/sanyarise/hezzl/internal/pb"
)

type Cash interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res chan *pb.User, key string) error
	GetCash(key string) ([]*pb.User, error)
}

type CashStore struct {
	cash Cash
}

func NewCashStore(cash Cash) *CashStore {
	return &CashStore{cash: cash}
}

func (cs *CashStore) CheckCash(key string) bool {
	return cs.cash.CheckCash(key)
}

func (cs *CashStore) CreateCash(ctx context.Context, res chan *pb.User, key string) error {
	return cs.cash.CreateCash(ctx, res, key)
}

func (cs *CashStore) GetCash(key string) ([]*pb.User, error) {
	res, err := cs.cash.GetCash(key)
	if err != nil {
		return nil, err
	}
	return res, nil
}
