
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    surname VARCHAR(125) DEFAULT '',
    patronymic VARCHAR(125) DEFAULT '',
    email VARCHAR(125) NOT NULL,
    user_id INT
    );

-- +migrate Down
DROP TABLE IF EXISTS users;