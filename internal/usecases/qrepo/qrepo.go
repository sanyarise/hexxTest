package qrepo

import "context"

type Queue interface {
	Enqueue(ctx context.Context, level string, msg string) error
}

type UserQueue struct {
	queue Queue
}

func NewUserQueue(queue Queue) *UserQueue {
	return &UserQueue{queue: queue}
}

func (uq *UserQueue) Enqueue(ctx context.Context, level string, msg string) error {
	return uq.queue.Enqueue(ctx, level, msg)
}
