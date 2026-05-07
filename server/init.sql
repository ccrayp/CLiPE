-- Порядок заполнения таблиц
-- 1-я очередь (справочники):
--   users
--   hosts
-- 2-я очередь (контроль доступа):
--   rules
--   policies
-- 3-я очередь (аудит):
--   requests
--   decisions


-- Создание структуры

CREATE TABLE hosts (
	host_id SERIAL PRIMARY KEY,
	ip VARCHAR(15) NOT NULL UNIQUE
);

CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	user_name VARCHAR(100) NOT NULL UNIQUE, 
	uid INT NOT NULL,
	gid INT
	host_id INT NOT NULL REFERENCES hosts(host_id) ON DELETE CASCADE
);

CREATE TABLE rules (
	rule_id SERIAL PRIMARY KEY,
	rule_name VARCHAR(100) NOT NULL UNIQUE,
	condition JSON NOT NULL,
	effect BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE policies (
	policy_id SERIAL PRIMARY KEY,
	policy_name VARCHAR(100) NOT NULL UNIQUE,
	user_id INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	rule_id INT NOT NULL REFERENCES rules(rule_id) ON DELETE CASCADE,
	status BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE requests (
	request_id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
	context JSONB NOT NULL DEFAULT '{"data": null}',
	timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE decisions (
	decision_id SERIAL PRIMARY KEY,
	request_id INT NOT NULL REFERENCES requests(request_id) ON DELETE CASCADE,
	policy_id INT REFERENCES policies(policy_id) ON DELETE CASCADE,
	result BOOLEAN NOT NULL
	timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    username TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TABLE sys_users (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL,
	password_hash CHAR(60) NOT NULL
)

INSERT INTO sys_users VALUES (DEFAULT, admin, $2a$10$sTQ8Zx0/KbfuvEn2JhbLjOPn5XEO5iHgxnQpKBE11RQaAvUkkiW9a)