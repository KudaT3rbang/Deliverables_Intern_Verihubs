-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id       INTEGER GENERATED ALWAYS AS IDENTITY,
    email    VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
