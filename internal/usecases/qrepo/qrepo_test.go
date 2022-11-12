package qrepo

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sanyarise/hezzl/internal/model"
)

func TestEnqueue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	queue := NewMockQueue(ctrl)
	userQueue := NewUserQueue(queue)

	message := model.Log{
		Time:    "test time",
		Level:   "test level",
		Message: "test message",
	}

	queue.EXPECT().Enqueue(ctx, message.Level, message.Message).Return(nil)

	err := userQueue.Enqueue(ctx, message.Level, message.Message)

	if err != nil {
		t.Errorf("failed to enqueue: %v", err)
	}
}
