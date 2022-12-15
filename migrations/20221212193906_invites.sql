-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invites (
    id           INTEGER   PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER   NOT NULL REFERENCES users(id),
    invite_hash  TEXT      NOT NULL UNIQUE,
    is_activated BOOLEAN   NOT NULL DEFAULT (false),
    created_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    updated_at   TIMESTAMP NOT NULL DEFAULT (datetime('now')),
    deleted_at   TIMESTAMP
);
CREATE UNIQUE INDEX invites_user_id_is_activated_uidx ON invites (user_id, is_activated) WHERE is_activated IS FALSE;
-- 00000000-0000-0000-0000-000000000000 invite id for the first user
INSERT INTO invites (id, user_id, invite_hash) VALUES (0, 0, 'Erk3fL5-XJTopw2dI5KVI9FK-pVHkxMPijlZx7hJrKg=');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE invites;
-- +goose StatementEnd
