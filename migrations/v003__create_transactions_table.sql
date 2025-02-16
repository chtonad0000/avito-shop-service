CREATE TABLE IF NOT EXISTS coin_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    counterpart_username VARCHAR(255),
    amount INT NOT NULL,
    transaction_type VARCHAR(10),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (counterpart_username) REFERENCES users(username)
);


CREATE INDEX IF NOT EXISTS idx_user_id ON coin_transactions (user_id);