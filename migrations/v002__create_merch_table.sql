CREATE TABLE IF NOT EXISTS merch (
    id SERIAL PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    price INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_item_name ON merch(item_name);

INSERT INTO merch (item_name, price) VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500);