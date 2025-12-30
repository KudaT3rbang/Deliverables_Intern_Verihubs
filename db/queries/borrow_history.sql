-- name: CreateBorrowHistory :one
INSERT INTO borrow_history (book_id, user_id, borrowed_at)
VALUES ($1, $2, $3)
RETURNING id;

-- name: ListBorrowHistoryByBookID :many
SELECT id, book_id, user_id, borrowed_at, returned_at
FROM borrow_history
WHERE book_id = $1
ORDER BY borrowed_at DESC;

-- name: IsBookBorrowed :one
SELECT EXISTS(
    SELECT 1
    FROM borrow_history
    WHERE book_id = $1 AND returned_at IS NULL
);

-- name: GetActiveBorrowHistory :one
SELECT id, book_id, user_id, borrowed_at, returned_at
FROM borrow_history
WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL;

-- name: UpdateBorrowHistoryReturnDate :exec
UPDATE borrow_history
SET returned_at = $1
WHERE id = $2;