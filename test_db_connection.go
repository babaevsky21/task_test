package main

import (
	"crypto-watcher/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Тестирование подключения к базе данных PostgreSQL...")

	cfg := config.Load()
	fmt.Printf("Попытка подключения к базе данных:\n")
	fmt.Printf("Host: %s\n", cfg.DBHost)
	fmt.Printf("Port: %s\n", cfg.DBPort)
	fmt.Printf("User: %s\n", cfg.DBUser)
	fmt.Printf("Database: %s\n", cfg.DBName)

	// Сначала подключимся к базе postgres для создания базы task
	postgresConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword)

	fmt.Println("\n1. Подключение к базе данных postgres...")
	postgresDB, err := sql.Open("postgres", postgresConnStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к postgres: %v", err)
	}
	defer postgresDB.Close()

	if err := postgresDB.Ping(); err != nil {
		log.Fatalf("Ошибка ping postgres: %v", err)
	}

	fmt.Println("✅ Успешное подключение к postgres!")

	// Проверим, существует ли база данных task
	var exists bool
	err = postgresDB.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", cfg.DBName).Scan(&exists)
	if err != nil {
		log.Fatalf("Ошибка проверки существования базы данных: %v", err)
	}

	if !exists {
		fmt.Printf("База данных '%s' не существует. Создаём...\n", cfg.DBName)
		_, err = postgresDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
		if err != nil {
			log.Fatalf("Ошибка создания базы данных: %v", err)
		}
		fmt.Printf("✅ База данных '%s' успешно создана!\n", cfg.DBName)
	} else {
		fmt.Printf("✅ База данных '%s' уже существует!\n", cfg.DBName)
	}

	// Теперь подключимся к созданной базе данных
	taskConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	fmt.Printf("\n2. Подключение к базе данных '%s'...\n", cfg.DBName)
	taskDB, err := sql.Open("postgres", taskConnStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных %s: %v", cfg.DBName, err)
	}
	defer taskDB.Close()

	if err := taskDB.Ping(); err != nil {
		log.Fatalf("Ошибка ping базы данных %s: %v", cfg.DBName, err)
	}

	fmt.Printf("✅ Успешное подключение к базе данных '%s'!\n", cfg.DBName)

	// Проверим, существуют ли наши таблицы
	var tableCount int
	err = taskDB.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name IN ('watchlist', 'price_history')").Scan(&tableCount)
	if err != nil {
		log.Printf("Ошибка проверки таблиц: %v", err)
	} else {
		fmt.Printf("Найдено таблиц в базе данных: %d из 2\n", tableCount)
		if tableCount == 0 {
			fmt.Println("\n3. Создание таблиц...")
			// Создаём таблицы
			_, err = taskDB.Exec(`
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
			`)
			if err != nil {
				log.Printf("Ошибка создания таблиц: %v", err)
			} else {
				fmt.Println("✅ Таблицы успешно созданы!")

				// Добавим тестовые данные
				fmt.Println("\n4. Добавление тестовых данных...")
				_, err = taskDB.Exec(`
					INSERT INTO watchlist (coin) VALUES ('BTC') ON CONFLICT (coin) DO NOTHING;
					INSERT INTO watchlist (coin) VALUES ('ETH') ON CONFLICT (coin) DO NOTHING;
					INSERT INTO price_history (coin, price, timestamp) VALUES 
						('BTC', 45000.50, 1736500490),
						('ETH', 3000.25, 1736500490);
				`)
				if err != nil {
					log.Printf("Ошибка добавления тестовых данных: %v", err)
				} else {
					fmt.Println("✅ Тестовые данные добавлены!")
				}
			}
		} else if tableCount == 2 {
			fmt.Println("✅ Все необходимые таблицы найдены!")
		} else {
			fmt.Println("⚠️  Найдены не все таблицы. Проверьте структуру базы данных.")
		}
	}

	fmt.Println("\n🎉 Инициализация базы данных завершена успешно!")
	fmt.Println("Теперь вы можете запустить приложение командой: go run cmd/main.go")
}
