package cashrepo

import (
	"context"

	"github.com/sanyarise/hezzl/internal/pb"
)

/*type Cash interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res chan *pb.User, key string) error
	GetCash(key string) ([]*pb.User, error)
}*/

type CustomMockCash struct {
}

func NewCustomMockCash() *CustomMockCash {
	return &CustomMockCash{}
}

func (cmc *CustomMockCash) CheckCash(key string) bool {
	return false
}

func (cmc *CustomMockCash) CreateCash(ctx context.Context, res chan *pb.User, key string) error {
	return nil
}

func (cmc *CustomMockCash) GetCash(key string) ([]*pb.User, error) {
	out := make([]*pb.User,0, 1)
	out = append(out, &pb.User{Id: "TestID", Name: "TestName"})
	return out, nil
}