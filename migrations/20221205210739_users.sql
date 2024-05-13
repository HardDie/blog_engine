-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER   PRIMARY KEY AUTOINCREMENT,
    username        TEXT      NOT NULL UNIQUE,
    displayed_name  TEXT      NOT NULL,
    email           TEXT,
    invited_by_user INTEGER   NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at      TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at      TIMESTAMP
);
INSERT INTO users (id, username, displayed_name, invited_by_user) VALUES (0, 'root', 'root', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
