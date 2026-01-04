-- name: CreateNotification :one
INSERT INTO notifications (user_id, message, sent_at, status)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetNotificationsByUserID :many
SELECT id, user_id, message, sent_at, status
FROM notifications
WHERE user_id = $1
ORDER BY sent_at DESC;

-- name: UpdateNotificationStatus :exec
UPDATE notifications
SET status = $1
WHERE id = $2;