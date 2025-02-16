CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    user_id INT,
    item_id INT,
    quantity INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (item_id) REFERENCES merch(id),
    UNIQUE (user_id, item_id)
);

CREATE INDEX IF NOT EXISTS idx_user_id ON inventory(user_id);