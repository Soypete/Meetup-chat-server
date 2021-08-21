-- +goose Up

CREATE TABLE IF NOT EXISTS chat_message (
	id SERIAL,
	username VARCHAR(64) NOT NULL,
	message_body TEXT,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	source VARCHAR(64)
	);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL,
	username VARCHAR(64) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	source VARCHAR(64) /* source + username need to be unique */
	);


-- +goose Down

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS chat_message;
DROP INDEX IF EXISTS index_users;
