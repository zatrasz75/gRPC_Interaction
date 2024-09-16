
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(125) NOT NULL,
    password VARCHAR(255) NOT NULL,
    data  TIMESTAMP WITHOUT TIME ZONE
    );


-- +migrate Down
DROP TABLE IF EXISTS users;