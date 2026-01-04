package jobs

import (
	"context"
	"fmt"
	"lendbook/internal/entity"
	"lendbook/internal/repository"
	"time"

	"github.com/riverqueue/river"
)

type SendNotificationArgs struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

func (SendNotificationArgs) Kind() string { return "send_notification" }

type SendNotificationWorker struct {
	river.WorkerDefaults[SendNotificationArgs]
	notificationRepo repository.NotificationRepository
}

func NewSendNotificationWorker(notificationRepo repository.NotificationRepository) *SendNotificationWorker {
	return &SendNotificationWorker{
		notificationRepo: notificationRepo,
	}
}

func (w *SendNotificationWorker) Work(ctx context.Context, job *river.Job[SendNotificationArgs]) error {
	notification := &entity.Notification{
		UserID:  job.Args.UserID,
		Message: job.Args.Message,
		SentAt:  time.Now(),
		Status:  entity.NotificationStatusSent,
	}

	err := w.notificationRepo.Create(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}
