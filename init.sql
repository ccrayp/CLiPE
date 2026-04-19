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


INSERT INTO actions (action_name) VALUES
('read'),
('write'),
('execute'),
('delete');

INSERT INTO users (user_name, uid, gid) VALUES
('alice', 1001, 100),
('bob', 1002, 100),
('charlie', 1003, 200),
('dave', 1004, 200);

INSERT INTO hosts (ip) VALUES
('192.168.1.10'),
('192.168.1.20'),
('10.0.0.5'),
('10.0.0.6');

INSERT INTO services (service_name) VALUES
('ssh'),
('http'),
('database'),
('ftp');

INSERT INTO rules (rule_name, condition, effect) VALUES
('Allow local SSH', '{"ip_range": "192.168.1.0/24", "service": "ssh"}', true),
('Deny external DB', '{"ip_range": "0.0.0.0/0", "service": "database"}', false),
('Allow HTTP', '{"service": "http"}', true),
('Deny FTP for group 200', '{"gid": 200, "service": "ftp"}', false);

INSERT INTO policies (policy_name, user_id, host_id, service_id, action_id, rule_id, status) VALUES
('Alice SSH local', 1, 1, 1, 3, 1, true),
('Bob HTTP access', 2, 2, 2, 1, 3, true),
('Charlie DB denied', 3, 3, 3, 1, 2, false),
('Dave FTP denied', 4, 4, 4, 2, 4, false);

INSERT INTO requests (user_id, host_id, service_id, action_id) VALUES
(1, 1, 1, 3),
(2, 2, 2, 1),
(3, 3, 3, 1),
(4, 4, 4, 2),
(1, 2, 2, 1),
(1, 3, 3, 1),
(2, 1, 1, 3),
(2, 4, 4, 2),
(3, 2, 2, 1),
(3, 1, 1, 3),
(4, 3, 3, 1),
(4, 2, 2, 1),
(1, 4, 4, 2),
(2, 3, 3, 1),
(3, 4, 4, 2),
(4, 1, 1, 3);

INSERT INTO decisions (request_id, policy_id, result) VALUES
(1, 1, true),
(2, 2, true),
(3, 3, false),
(4, 4, false),
(5, 2, true),
(6, 3, false),
(7, 1, true),
(8, 4, false),
(9, 2, true),
(10, 1, true),
(11, 3, false),
(12, 2, true),
(13, 4, false),
(14, 3, false),
(15, 4, false),
(16, 1, true);