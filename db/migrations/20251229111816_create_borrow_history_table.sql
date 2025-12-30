-- +goose Up
-- +goose StatementBegin
CREATE TABLE borrow_history
(
    id          INTEGER GENERATED ALWAYS AS IDENTITY,
    book_id     INTEGER                             NOT NULL,
    user_id     INTEGER                             NOT NULL,
    borrowed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    returned_at TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_borrow_book
        FOREIGN KEY (book_id) REFERENCES books(id),
    CONSTRAINT fk_borrow_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS borrow_history;
-- +goose StatementEnd
