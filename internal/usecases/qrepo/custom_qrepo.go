package qrepo

import "context"

type CustomMockQueue struct {
}

func NewCustomMockQueue() *CustomMockQueue {
	return &CustomMockQueue{}
}

func (cmq *CustomMockQueue) Enqueue(ctx context.Context, level string, msg string) error {
	return nil
}
