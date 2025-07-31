package storage

import (
	"crypto-watcher/internal/models"
	"testing"
)

// Это базовый пример теста. В реальном проекте здесь были бы более полные тесты с подключением к тестовой БД

func TestStorageInterface(t *testing.T) {
	// Проверяем, что структура Storage имеет все необходимые методы
	var s *Storage

	// Эти методы должны существовать (компиляция не пройдет, если их нет)
	_ = s.AddCoin
	_ = s.RemoveCoin
	_ = s.GetWatchlist
	_ = s.SavePrice
	_ = s.GetPrice
	_ = s.CoinExists

	t.Log("All Storage methods are present")
}

func TestModelsStructure(t *testing.T) {
	// Проверяем, что структуры моделей имеют необходимые поля
	coin := models.Cryptocurrency{}
	if coin.Coin == "" {
		t.Log("Coin field exists in Cryptocurrency model")
	}

	price := models.PriceHistory{}
	if price.Price == 0 {
		t.Log("Price field exists in PriceHistory model")
	}

	req := models.AddCoinRequest{}
	if req.Coin == "" {
		t.Log("Coin field exists in AddCoinRequest model")
	}
}
