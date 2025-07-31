-- Подключитесь к PostgreSQL как суперпользователь и выполните эти команды

-- Создание базы данных (если не существует)
CREATE DATABASE task;

-- Подключение к базе данных task
\c task;

-- Создание таблиц
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

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_price_history_coin_timestamp ON price_history(coin, timestamp);
CREATE INDEX IF NOT EXISTS idx_watchlist_coin ON watchlist(coin);

-- Проверка созданных таблиц
\dt

-- Вставка тестовых данных для проверки
INSERT INTO watchlist (coin) VALUES ('BTC') ON CONFLICT (coin) DO NOTHING;
INSERT INTO watchlist (coin) VALUES ('ETH') ON CONFLICT (coin) DO NOTHING;

-- Вставка тестовых цен
INSERT INTO price_history (coin, price, timestamp) VALUES 
    ('BTC', 45000.50, 1736500490),
    ('ETH', 3000.25, 1736500490);

-- Проверка данных
SELECT * FROM watchlist;
SELECT * FROM price_history;
