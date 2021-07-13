-- +goose Up

CREATE TABLE IF NOT EXISTS users (
	id SERIAL,
	username VARCHAR(64) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	source VARCHAR(64) /* source + username need to be unique */
	);

CREATE INDEX index_users ON users (username, source);

CREATE TABLE IF NOT EXISTS chat_message (
	id SERIAL,
	username VARCHAR(64) NOT NULL,
	message_body TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	source VARCHAR(64)
	);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS chat_message;
DROP INDEX IF EXISTS index_users;
