package main

import (
	"crypto-watcher/internal/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("🔍 Быстрая проверка данных в БД")
	fmt.Println("==============================")
	
	cfg := config.Load()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()
	
	// Показываем watchlist
	fmt.Println("\n📋 Watchlist:")
	rows, err := db.Query("SELECT coin, created_at FROM watchlist ORDER BY created_at")
	if err != nil {
		log.Printf("Ошибка запроса watchlist: %v", err)
	} else {
		defer rows.Close()
		count := 0
		for rows.Next() {
			var coin string
			var createdAt time.Time
			rows.Scan(&coin, &createdAt)
			fmt.Printf("  %d. %s (добавлен: %s)\n", count+1, coin, createdAt.Format("15:04:05"))
			count++
		}
		if count == 0 {
			fmt.Println("  (пусто)")
		} else {
			fmt.Printf("Всего: %d монет\n", count)
		}
	}
	
	// Показываем последние цены
	fmt.Println("\n💰 Последние цены (топ 10):")
	rows, err = db.Query(`
		SELECT coin, price, timestamp, created_at 
		FROM price_history 
		ORDER BY created_at DESC 
		LIMIT 10
	`)
	if err != nil {
		log.Printf("Ошибка запроса price_history: %v", err)
	} else {
		defer rows.Close()
		count := 0
		for rows.Next() {
			var coin string
			var price float64
			var timestamp int64
			var createdAt time.Time
			rows.Scan(&coin, &price, &timestamp, &createdAt)
			fmt.Printf("  %d. %s: $%.2f (время: %s)\n", 
				count+1, coin, price, createdAt.Format("15:04:05"))
			count++
		}
		if count == 0 {
			fmt.Println("  (пусто)")
		}
		
		// Общее количество записей
		var total int
		db.QueryRow("SELECT COUNT(*) FROM price_history").Scan(&total)
		fmt.Printf("Всего записей в истории: %d\n", total)
	}
	
	// Статистика по монетам
	fmt.Println("\n📊 Статистика по монетам:")
	rows, err = db.Query(`
		SELECT coin, COUNT(*) as count, MIN(price) as min_price, MAX(price) as max_price, AVG(price) as avg_price
		FROM price_history 
		GROUP BY coin 
		ORDER BY count DESC
	`)
	if err != nil {
		log.Printf("Ошибка запроса статистики: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var coin string
			var count int
			var minPrice, maxPrice, avgPrice float64
			rows.Scan(&coin, &count, &minPrice, &maxPrice, &avgPrice)
			fmt.Printf("  %s: %d записей (мин: $%.2f, макс: $%.2f, сред: $%.2f)\n", 
				coin, count, minPrice, maxPrice, avgPrice)
		}
	}
}
