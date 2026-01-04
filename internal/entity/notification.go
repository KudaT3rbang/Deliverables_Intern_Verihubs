package entity

import "time"

type Notification struct {
	ID      int       `json:"id"`
	UserID  int       `json:"user_id"`
	Message string    `json:"message"`
	SentAt  time.Time `json:"sent_at"`
	Status  string    `json:"status"`
}

const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)
