package main

import (
	"crypto-watcher/internal/config"
	"crypto-watcher/internal/database"
	"database/sql"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к базе данных (без автоматических миграций)
	db, err := database.ConnectWithoutMigrations(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	switch command {
	case "migrate":
		log.Println("Running migrations...")
		if err := database.RunMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations completed successfully!")

	case "status":
		log.Println("Checking migration status...")
		if err := showMigrationStatus(db); err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

	case "rollback":
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s rollback <version>", os.Args[0])
		}
		// Функционал отката пока не реализован
		log.Println("Rollback functionality is not implemented yet")

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`Migration CLI Tool

Usage:
    %s migrate     - Run all pending migrations
    %s status      - Show migration status
    %s rollback <version> - Rollback to specific version (not implemented)

Examples:
    %s migrate
    %s status
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func showMigrationStatus(db *sql.DB) error {
	// Получаем список всех миграций
	allMigrations := database.GetMigrations()

	// Получаем примененные миграции
	appliedVersions, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %v", err)
	}

	appliedMap := make(map[int]bool)
	for _, version := range appliedVersions {
		appliedMap[version] = true
	}

	fmt.Println("Migration Status:")
	fmt.Println("=================")

	for _, migration := range allMigrations {
		status := "PENDING"
		if appliedMap[migration.Version] {
			status = "APPLIED"
		}
		fmt.Printf("Version %d: %s [%s]\n", migration.Version, migration.Description, status)
	}

	return nil
}

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
