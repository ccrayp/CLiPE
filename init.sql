-- Порядок заполнения таблиц
-- 1-я очередь (справочники):
--   actions
--   users
--   hosts
--   services
-- 2-я очередь (контроль доступа):
--   rules
--   policies
-- 3-я очередь (аудит):
--   requests
--   decisions


CREATE TABLE actions (
	action_id SERIAL PRIMARY KEY,
	action_name VARCHAR(50) NOT NULL
);

CREATE TABLE users (
	user_id SERIAL PRIMARY KEY,
	user_name VARCHAR(100),
	uid INT NOT NULL,
	gid INT
);

CREATE TABLE hosts (
	host_id SERIAL PRIMARY KEY,
	ip CHAR(15) NOT NULL
);

CREATE TABLE services (
	service_id SERIAL PRIMARY KEY,
	service_name VARCHAR(50) NOT NULL
);

CREATE TABLE rules (
	rule_id SERIAL PRIMARY KEY,
	rule_name VARCHAR(100) NOT NULL,
	condition JSON NOT NULL,
	effect BOOLEAN NOT NULL DEFAULT 'false'
);

CREATE TABLE policies (
	policy_id SERIAL PRIMARY KEY,
	policy_name VARCHAR(100),
	user_id INT REFERENCES users(user_id),
	host_id INT REFERENCES hosts(host_id),
	service_id INT REFERENCES services(service_id),
	action_id INT REFERENCES actions(action_id),
	rule_id INT REFERENCES rules(rule_id),
	status BOOLEAN NOT NULL DEFAULT 'false'
);

CREATE TABLE requests (
	request_id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(user_id) NOT NULL,
	host_id INT REFERENCES hosts(host_id) NOT NULL,
	service_id INT REFERENCES services(service_id) NOT NULL,
	action_id INT REFERENCES actions(action_id) NOT NULL
);

CREATE TABLE decisions (
	decision_id SERIAL PRIMARY KEY,
	request_id INT REFERENCES requests(request_id) NOT NULL,
	policy_id INT REFERENCES policies(policy_id) NOT NULL,
	result BOOLEAN NOT NULL
);
