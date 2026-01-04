-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications
(
    id         INTEGER GENERATED ALWAYS AS IDENTITY,
    user_id    INTEGER                             NOT NULL,
    message    TEXT                                NOT NULL,
    sent_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status     VARCHAR(50)                         NOT NULL DEFAULT 'pending',
    PRIMARY KEY (id),
    CONSTRAINT fk_notifications_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
