CREATE TABLE IF NOT EXISTS watchlist (
    id SERIAL PRIMARY KEY,
    coin VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    price DECIMAL(20, 8) NOT NULL,
    timestamp BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coin) REFERENCES watchlist(coin) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_price_history_coin_timestamp ON price_history(coin, timestamp);
CREATE INDEX IF NOT EXISTS idx_watchlist_coin ON watchlist(coin);
