
-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(75) UNIQUE NOT NULL,
    user_id INT
    );

INSERT INTO roles (name, user_id) VALUES
    ('Admin', 1),
    ('User', 1),
    ('Guest', 1);

-- +migrate Down
DROP TABLE IF EXISTS roles;