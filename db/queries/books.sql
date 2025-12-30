-- name: CreateBook :one
INSERT INTO books (title, author, published_date, language, added_at, added_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: GetBookByID :one
SELECT id, title, author, published_date, language, added_at, added_by, deleted_at, deleted_by
FROM books
WHERE id = $1;

-- name: UpdateBook :exec
UPDATE books
SET title=$1, author=$2, published_date=$3, language=$4, deleted_at=$5, deleted_by=$6
WHERE id=$7;

-- name: ListAvailableBooks :many
SELECT id, title, author, published_date, language, added_at, added_by
FROM books
WHERE deleted_at IS NULL;