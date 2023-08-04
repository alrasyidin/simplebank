package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send-verify-email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal json payload %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueu task %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", taskInfo.Queue).Int("max_retry", taskInfo.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload *PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed unmarshal task payload %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user doens't exists %w", asynq.SkipRetry)
		// }

		return fmt.Errorf("failed to get user %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}

	subject := "Welcome to Simple Bank"
	content := fmt.Sprintf(`Hello %s, Thank you for registering us, use this <a href="http://localhost:8081/v1/verify?email_id=%d&secret_code=%s">Verify</a> to verify email address.`, user.FullName, verifyEmail.ID, verifyEmail.SecretCode)
	err = processor.mailer.SendEmail(subject, content, []string{user.Email}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email verification: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("processed task")
	return nil
}
