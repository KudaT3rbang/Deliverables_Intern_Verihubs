-- name: CreateBorrowHistory :one
INSERT INTO borrow_history (book_id, user_id, borrowed_at, borrowed_until)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ListBorrowHistoryByBookID :many
SELECT id, book_id, user_id, borrowed_at, borrowed_until, returned_at
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
SELECT id, book_id, user_id, borrowed_at, borrowed_until, returned_at
FROM borrow_history
WHERE user_id = $1 AND book_id = $2 AND returned_at IS NULL;

-- name: UpdateBorrowHistoryReturnDate :exec
UPDATE borrow_history
SET returned_at = $1
WHERE id = $2;

-- name: GetOverdueBorrows :many
SELECT bh.id, bh.book_id, bh.user_id, bh. borrowed_at, bh.borrowed_until, bh.returned_at, b.title as book_title
FROM borrow_history bh
     JOIN books b ON bh. book_id = b.id
WHERE bh.returned_at IS NULL AND bh.borrowed_until < NOW();