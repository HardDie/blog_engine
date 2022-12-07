-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER   PRIMARY KEY AUTOINCREMENT,
    username        TEXT      NOT NULL UNIQUE,
    displayed_name  TEXT,
    email           TEXT,
    invited_by_user INTEGER,
    created_at      TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at      TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at      TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
