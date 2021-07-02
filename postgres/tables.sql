CREATE TABLE IF NOT EXISTS users (
	id INT PRIMARY KEY AUTO INCREMENT,
	username VARCHAR(64) NOT NULL,
	created_at TIMESTAMPZ DEFAULT NOW(),
	source VARCHAR(64) /* source + username need to be unique */
	);

CREATE INDEX index_users ON users (username, source);

CREATE TABLE IF NOT EXISTS chat_message (
	id INT PRIMARY KEY AUTO INCREMENT,
	username VARCHAR(64) NOT NULL,
	message_body TEXT,
	created_at TIMESTAMPZ DEFAULT NOW(),
	source VARCHAR(64)
	);
