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

CREATE TABLE services (
	service_id SERIAL PRIMARY KEY,
	service_name TEXT NOT NULL UNIQUE
)

CREATE TABLE rules (
	rule_id SERIAL PRIMARY KEY,
	rule_name VARCHAR(100) NOT NULL UNIQUE,
	condition JSON NOT NULL,
	effect BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE policies (
	policy_id SERIAL PRIMARY KEY,
	policy_name VARCHAR(100) NOT NULL UNIQUE,
	user_id INT NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
	status BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE policy_contents (
	policy_id INT NOT NULL REFERENCES policies(policy_id) ON DELETE CASCADE,
    service_id INT NOT NULL REFERENCES services(service_id) ON DELETE CASCADE,
    rule_id INT NOT NULL REFERENCES rules(rule_id) ON DELETE CASCADE,
    PRIMARY KEY (policy_id, service_id)
)

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


INSERT INTO sys_users VALUES (DEFAULT, 'admin', '$2a$10$sTQ8Zx0/KbfuvEn2JhbLjOPn5XEO5iHgxnQpKBE11RQaAvUkkiW9a')


CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_rules_rule_name_trgm ON rules USING GIN (rule_name gin_trgm_ops);
CREATE INDEX idx_users_user_name_trgm ON users USING GIN (user_name gin_trgm_ops);
CREATE INDEX idx_services_service_name_trgm ON services USING GIN (service_name gin_trgm_ops);
CREATE INDEX idx_policies_policy_name_trgm ON policies USING GIN (policy_name gin_trgm_ops);

-- SELECT 
--     *
-- FROM rules
-- WHERE similarity(rule_name, 'wa') > 0.1