package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"log"

	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task %w", err)
	}

	log.Printf("RedisTaskDistributor Type %v and task Payload: %v", task.Type(), string(info.Payload))
	log.Printf("RedisTaskDistributor Queue %v and info MaxRetry: %v", info.Queue, info.MaxRetry)

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal task payload %w", err)
	}

	user, err := processor.store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with email %v doesnt exist", payload.Email)
		}
		return fmt.Errorf("failed to get user %w", err)
	}

	// Send Email to user
	log.Printf("RedisTaskProcessor Type %v and task payload: %v for user: %v", task.Type(), string(task.Payload()), user.Email)

	return nil
}
