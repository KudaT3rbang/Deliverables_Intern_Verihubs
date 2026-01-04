package repository

import (
	"context"
	"lendbook/internal/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetByUserID(ctx context.Context, userID int) ([]entity.Notification, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}
