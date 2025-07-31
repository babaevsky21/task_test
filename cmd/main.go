package main

import (
	_ "crypto-watcher/docs"
	"crypto-watcher/internal/api"
	"crypto-watcher/internal/config"
	"crypto-watcher/internal/database"
	"crypto-watcher/internal/service"
	"crypto-watcher/internal/storage"
	"log"
)

// @title Crypto Watcher API
// @version 1.0
// @description Микросервис для мониторинга криптовалют
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	storage := storage.New(db)
	cryptoService := service.NewCryptoService(storage)

	router := api.SetupRouter(cryptoService)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
