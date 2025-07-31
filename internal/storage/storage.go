package storage

import (
	"crypto-watcher/internal/models"
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddCoin(coin string) error {
	query := `INSERT INTO watchlist (coin) VALUES ($1) ON CONFLICT (coin) DO NOTHING`
	_, err := s.db.Exec(query, coin)
	return err
}

func (s *Storage) RemoveCoin(coin string) error {
	query := `DELETE FROM watchlist WHERE coin = $1`
	_, err := s.db.Exec(query, coin)
	return err
}

func (s *Storage) GetWatchlist() ([]models.Cryptocurrency, error) {
	query := `SELECT coin, created_at FROM watchlist ORDER BY created_at`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []models.Cryptocurrency
	for rows.Next() {
		var coin models.Cryptocurrency
		if err := rows.Scan(&coin.Coin, &coin.CreatedAt); err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

func (s *Storage) SavePrice(coin string, price float64, timestamp int64) error {
	query := `INSERT INTO price_history (coin, price, timestamp) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, coin, price, timestamp)
	return err
}

func (s *Storage) GetPrice(coin string, timestamp int64) (*models.PriceHistory, error) {
	query := `
		SELECT id, coin, price, timestamp, created_at 
		FROM price_history 
		WHERE coin = $1 AND timestamp <= $2 
		ORDER BY timestamp DESC 
		LIMIT 1`

	var price models.PriceHistory
	err := s.db.QueryRow(query, coin, timestamp).Scan(
		&price.ID, &price.Coin, &price.Price, &price.Timestamp, &price.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &price, nil
}

func (s *Storage) CoinExists(coin string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM watchlist WHERE coin = $1)`
	var exists bool
	err := s.db.QueryRow(query, coin).Scan(&exists)
	return exists, err
}
