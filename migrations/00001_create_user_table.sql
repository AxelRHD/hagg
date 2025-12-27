-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
uid TEXT UNIQUE NOT NULL, -- random secret, handed to user
    display_name TEXT UNIQUE NOT NULL,
    last_name TEXT NOT NULL DEFAULT '',
    first_name TEXT NOT NULL DEFAULT '',
    created_at TEXT DEFAULT (datetime('now', 'localtime')) NOT NULL,
    updated_at TEXT DEFAULT (datetime('now', 'localtime')) NOT NULL
);

CREATE TRIGGER IF NOT EXISTS on_update_ts_users
AFTER UPDATE ON users
WHEN OLD.updated_at <> (datetime('now', 'localtime'))
BEGIN
    UPDATE users
    SET updated_at = (datetime('now', 'localtime'))
    WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
