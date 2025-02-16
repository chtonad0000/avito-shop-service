CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    coins INT DEFAULT 1000
);

CREATE INDEX IF NOT EXISTS idx_username ON users(username);
