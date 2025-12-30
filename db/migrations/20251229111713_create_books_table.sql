-- +goose Up
-- +goose StatementBegin
CREATE TABLE books
(
    id             INTEGER GENERATED ALWAYS AS IDENTITY,
    title          VARCHAR(100)                        NOT NULL,
    author         VARCHAR(100)                        NOT NULL,
    published_date DATE                                NOT NULL,
    language       VARCHAR(100)                        NOT NULL,
    added_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    added_by       INTEGER                             NOT NULL,
    deleted_at     TIMESTAMP,
    deleted_by     INTEGER,
    PRIMARY KEY (id),
    CONSTRAINT fk_books_added_by
        FOREIGN KEY (added_by) REFERENCES users(id),
    CONSTRAINT fk_books_deleted_by
        FOREIGN KEY (deleted_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS books;
-- +goose StatementEnd
