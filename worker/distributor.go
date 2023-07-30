package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

// Constructor for RedisTaskDistributor
func NewRedisTaskDistributor(opt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(opt)
	return &RedisTaskDistributor{
		client: client,
	}
}
