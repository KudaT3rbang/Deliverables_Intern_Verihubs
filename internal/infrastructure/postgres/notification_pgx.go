package postgres

import (
	"context"
	generated "lendbook/internal/db"
	"lendbook/internal/entity"
	"lendbook/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type notificationRepository struct {
	queries *generated.Queries
}

func NewNotificationRepository(db *pgxpool.Pool) repository.NotificationRepository {
	return &notificationRepository{
		queries: generated.New(db),
	}
}

func (r *notificationRepository) Create(ctx context.Context, notification *entity.Notification) error {
	params := generated.CreateNotificationParams{
		UserID:  int32(notification.UserID),
		Message: notification.Message,
		SentAt:  pgtype.Timestamp{Time: notification.SentAt, Valid: true},
		Status:  notification.Status,
	}

	id, err := r.queries.CreateNotification(ctx, params)
	if err != nil {
		return err
	}

	notification.ID = int(id)
	return nil
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID int) ([]entity.Notification, error) {
	rows, err := r.queries.GetNotificationsByUserID(ctx, int32(userID))
	if err != nil {
		return nil, err
	}

	notifications := make([]entity.Notification, 0, len(rows))
	for _, row := range rows {
		n := entity.Notification{
			ID:      int(row.ID),
			UserID:  int(row.UserID),
			Message: row.Message,
			SentAt:  row.SentAt.Time,
			Status:  row.Status,
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	params := generated.UpdateNotificationStatusParams{
		Status: status,
		ID:     int32(id),
	}
	return r.queries.UpdateNotificationStatus(ctx, params)
}
