package service

import (
	"crypto-watcher/internal/models"
	"crypto-watcher/internal/storage"
	"errors"
	"fmt"
	"time"
)

type CryptoService struct {
	storage *storage.Storage
}

func NewCryptoService(storage *storage.Storage) *CryptoService {
	return &CryptoService{storage: storage}
}

func (s *CryptoService) AddCoin(coin string) error {
	return s.storage.AddCoin(coin)
}

func (s *CryptoService) RemoveCoin(coin string) error {
	exists, err := s.storage.CoinExists(coin)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("coin not found in watchlist")
	}
	return s.storage.RemoveCoin(coin)
}

func (s *CryptoService) GetWatchlist() ([]models.Cryptocurrency, error) {
	return s.storage.GetWatchlist()
}

func (s *CryptoService) GetPrice(coin string, timestamp int64) (*models.PriceHistory, error) {
	// Проверяем, что монета есть в watchlist
	exists, err := s.storage.CoinExists(coin)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("coin not found in watchlist")
	}

	// Получаем цену из базы данных
	price, err := s.storage.GetPrice(coin, timestamp)
	if err != nil {
		return nil, err
	}

	if price == nil {
		return nil, errors.New("price not found for the given timestamp")
	}

	return price, nil
}

// SavePrice сохраняет цену криптовалюты (для внутреннего использования или тестирования)
func (s *CryptoService) SavePrice(coin string, price float64, timestamp int64) error {
	// Проверяем, что монета есть в watchlist
	exists, err := s.storage.CoinExists(coin)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("coin not found in watchlist")
	}

	return s.storage.SavePrice(coin, price, timestamp)
}

// SimulatePriceCollection имитирует сбор цен каждые N секунд
func (s *CryptoService) SimulatePriceCollection() {
	ticker := time.NewTicker(10 * time.Second) // каждые 10 секунд для демонстрации
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			coins, err := s.storage.GetWatchlist()
			if err != nil {
				fmt.Printf("Error getting watchlist: %v\n", err)
				continue
			}

			for _, coin := range coins {
				// Генерируем случайную цену для демонстрации
				// В реальном приложении здесь был бы вызов к внешнему API
				mockPrice := generateMockPrice(coin.Coin)
				timestamp := time.Now().Unix()

				if err := s.storage.SavePrice(coin.Coin, mockPrice, timestamp); err != nil {
					fmt.Printf("Error saving price for %s: %v\n", coin.Coin, err)
				} else {
					fmt.Printf("Saved price for %s: $%.2f at %d\n", coin.Coin, mockPrice, timestamp)
				}
			}
		}
	}
}

// generateMockPrice генерирует мок-цену для демонстрации
func generateMockPrice(coin string) float64 {
	prices := map[string]float64{
		"BTC": 45000.0 + float64(time.Now().Unix()%1000),
		"ETH": 3000.0 + float64(time.Now().Unix()%500),
		"ADA": 1.0 + float64(time.Now().Unix()%100)/100,
		"DOT": 20.0 + float64(time.Now().Unix()%50),
	}

	if price, exists := prices[coin]; exists {
		return price
	}
	return 100.0 + float64(time.Now().Unix()%1000) // дефолтная цена
}
