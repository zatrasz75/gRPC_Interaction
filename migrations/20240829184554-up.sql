
-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(75) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    surname VARCHAR(125) DEFAULT '',
    patronymic VARCHAR(125) DEFAULT '',
    email VARCHAR(125) NOT NULL,
    password VARCHAR(255) NOT NULL,
    data  TIMESTAMP WITHOUT TIME ZONE
    );

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
    );

INSERT INTO roles (name) VALUES ('Admin'), ('User'), ('Guest');

-- +migrate Down
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;