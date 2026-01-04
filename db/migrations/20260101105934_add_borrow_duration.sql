-- +goose Up
-- +goose StatementBegin
ALTER TABLE books
    ADD COLUMN max_borrow_days INTEGER NOT NULL DEFAULT 14;

ALTER TABLE borrow_history
    ADD COLUMN borrowed_until TIMESTAMP;

UPDATE borrow_history
SET borrowed_until = borrowed_at + INTERVAL '14 days'
WHERE borrowed_until IS NULL;

ALTER TABLE borrow_history
    ALTER COLUMN borrowed_until SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE borrow_history DROP COLUMN borrowed_until;
ALTER TABLE books DROP COLUMN max_borrow_days;
-- +goose StatementEnd
