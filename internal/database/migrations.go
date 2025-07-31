package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration представляет одну миграцию
type Migration struct {
	Version     int
	Description string
	SQL         string
}

// GetMigrations возвращает список всех миграций в порядке применения
func GetMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Create initial tables",
			SQL: `
				-- Создание таблицы для отслеживания миграций
				CREATE TABLE IF NOT EXISTS schema_migrations (
					version INTEGER PRIMARY KEY,
					applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Создание таблицы watchlist
				CREATE TABLE IF NOT EXISTS watchlist (
					id SERIAL PRIMARY KEY,
					coin VARCHAR(10) NOT NULL UNIQUE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Создание таблицы price_history
				CREATE TABLE IF NOT EXISTS price_history (
					id SERIAL PRIMARY KEY,
					coin VARCHAR(10) NOT NULL,
					price DECIMAL(20, 8) NOT NULL,
					timestamp BIGINT NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (coin) REFERENCES watchlist(coin) ON DELETE CASCADE
				);
			`,
		},
		{
			Version:     2,
			Description: "Add indexes for performance",
			SQL: `
				-- Создание индексов для оптимизации запросов
				CREATE INDEX IF NOT EXISTS idx_price_history_coin_timestamp ON price_history(coin, timestamp);
				CREATE INDEX IF NOT EXISTS idx_watchlist_coin ON watchlist(coin);
			`,
		},
	}
}

// RunMigrations выполняет все неприменённые миграции
func RunMigrations(db *sql.DB) error {
	log.Println("Starting database migrations...")

	// Создаём таблицу для отслеживания миграций если её нет
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %v", err)
	}

	// Получаем список уже применённых миграций
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	// Применяем новые миграции
	migrations := GetMigrations()
	for _, migration := range migrations {
		if !contains(appliedMigrations, migration.Version) {
			log.Printf("Applying migration %d: %s", migration.Version, migration.Description)

			// Начинаем транзакцию
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction for migration %d: %v", migration.Version, err)
			}

			// Выполняем SQL миграции
			_, err = tx.Exec(migration.SQL)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %d: %v", migration.Version, err)
			}

			// Записываем информацию о применённой миграции
			_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
			}

			// Коммитим транзакцию
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit migration %d: %v", migration.Version, err)
			}

			log.Printf("Migration %d applied successfully", migration.Version)
		} else {
			log.Printf("Migration %d already applied, skipping", migration.Version)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// getAppliedMigrations возвращает список версий уже применённых миграций
func getAppliedMigrations(db *sql.DB) ([]int, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}

// contains проверяет, содержит ли слайс определённое значение
func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
