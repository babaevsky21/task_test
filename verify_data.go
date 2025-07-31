package main

import (
	"crypto-watcher/internal/config"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🔍 Проверка добавления данных в базу данных")
	fmt.Println("==========================================")

	// Подключение к базе данных
	cfg := config.Load()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка ping БД: %v", err)
	}

	fmt.Println("✅ Подключение к БД успешно")

	// Проверяем текущее состояние таблиц
	fmt.Println("\n📊 Текущее состояние watchlist:")
	showWatchlist(db)

	fmt.Println("\n📊 Текущее состояние price_history:")
	showPriceHistory(db)

	// Ждем несколько секунд, чтобы сервер запустился
	fmt.Println("\n⏳ Ожидание запуска сервера...")
	time.Sleep(3 * time.Second)

	// Тестируем API эндпоинты
	fmt.Println("\n🧪 Тестирование API:")

	// 1. Добавляем новую монету
	fmt.Println("\n1. Добавление DOGE в watchlist...")
	err = addCoin("DOGE")
	if err != nil {
		fmt.Printf("❌ Ошибка добавления DOGE: %v\n", err)
	} else {
		fmt.Println("✅ DOGE успешно добавлен")
	}

	// 2. Добавляем еще одну монету
	fmt.Println("\n2. Добавление ADA в watchlist...")
	err = addCoin("ADA")
	if err != nil {
		fmt.Printf("❌ Ошибка добавления ADA: %v\n", err)
	} else {
		fmt.Println("✅ ADA успешно добавлен")
	}

	// Ждем немного для обработки
	time.Sleep(2 * time.Second)

	// Проверяем изменения в БД
	fmt.Println("\n📊 Состояние watchlist после добавления:")
	showWatchlist(db)

	// Ждем немного для автоматического сбора цен
	fmt.Println("\n⏳ Ожидание автоматического сбора цен (15 секунд)...")
	time.Sleep(15 * time.Second)

	fmt.Println("\n📊 Состояние price_history после сбора цен:")
	showPriceHistory(db)

	// Тестируем получение цены
	fmt.Println("\n3. Тестирование получения цены...")
	timestamp := time.Now().Unix()
	price, err := getPrice("DOGE", timestamp)
	if err != nil {
		fmt.Printf("❌ Ошибка получения цены DOGE: %v\n", err)
	} else {
		fmt.Printf("✅ Цена DOGE на timestamp %d: $%.2f\n", timestamp, price)
	}

	fmt.Println("\n🎉 Тестирование завершено!")
}

func showWatchlist(db *sql.DB) {
	rows, err := db.Query("SELECT coin, created_at FROM watchlist ORDER BY created_at")
	if err != nil {
		fmt.Printf("❌ Ошибка запроса watchlist: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var coin string
		var createdAt time.Time
		if err := rows.Scan(&coin, &createdAt); err != nil {
			fmt.Printf("❌ Ошибка сканирования: %v\n", err)
			continue
		}
		fmt.Printf("  - %s (добавлен: %s)\n", coin, createdAt.Format("2006-01-02 15:04:05"))
		count++
	}

	if count == 0 {
		fmt.Println("  (пусто)")
	} else {
		fmt.Printf("Всего монет в watchlist: %d\n", count)
	}
}

func showPriceHistory(db *sql.DB) {
	rows, err := db.Query("SELECT coin, price, timestamp, created_at FROM price_history ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		fmt.Printf("❌ Ошибка запроса price_history: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var coin string
		var price float64
		var timestamp int64
		var createdAt time.Time
		if err := rows.Scan(&coin, &price, &timestamp, &createdAt); err != nil {
			fmt.Printf("❌ Ошибка сканирования: %v\n", err)
			continue
		}
		fmt.Printf("  - %s: $%.2f (timestamp: %d, записано: %s)\n",
			coin, price, timestamp, createdAt.Format("2006-01-02 15:04:05"))
		count++
	}

	if count == 0 {
		fmt.Println("  (пусто)")
	}

	// Показываем общее количество записей
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM price_history").Scan(&total)
	if err == nil {
		fmt.Printf("Всего записей в price_history: %d\n", total)
	}
}

func addCoin(coin string) error {
	payload := fmt.Sprintf(`{"coin": "%s"}`, coin)
	resp, err := http.Post("http://localhost:8080/currency/add",
		"application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return nil
}

func getPrice(coin string, timestamp int64) (float64, error) {
	url := fmt.Sprintf("http://localhost:8080/currency/price?coin=%s&timestamp=%d", coin, timestamp)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result struct {
		Price float64 `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Price, nil
}
