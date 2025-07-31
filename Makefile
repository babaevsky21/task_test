.PHONY: build run test clean docs docker-build docker-run migrate migrate-status

# Переменные
APP_NAME := crypto-watcher
BUILD_DIR := bin
DOCKER_IMAGE := $(APP_NAME):latest

# Сборка приложения
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Сборка CLI утилиты миграций
build-migrate:
	@echo "Building migration CLI..."
	@go build -o $(BUILD_DIR)/migrate cmd/migrate/main.go
	@echo "Migration CLI built: $(BUILD_DIR)/migrate"

# Запуск приложения
run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/main.go

# Установка зависимостей
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Генерация Swagger документации
docs:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/main.go

# Выполнение миграций базы данных
migrate:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go migrate

# Проверка статуса миграций
migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go status

# Запуск тестов
test:
	@echo "Running tests..."
	@go test -v ./...

# Инициализация базы данных (устаревший метод)
init-db:
	@echo "Initializing database..."
	@psql -h localhost -p 5432 -U postgres -d postgres -f setup_database.sql

# Очистка build директории
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

# Сборка Docker образа
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Запуск с Docker Compose
docker-run:
	@echo "Starting with Docker Compose..."
	@docker-compose up --build

# Запуск в фоновом режиме
docker-up:
	@echo "Starting Docker Compose in background..."
	@docker-compose up --build -d

# Остановка Docker Compose
docker-stop:
	@echo "Stopping Docker Compose..."
	@docker-compose down

# Полная очистка Docker
docker-clean:
	@echo "Cleaning up Docker..."
	@docker-compose down --rmi all --volumes --remove-orphans

# Полная настройка проекта
setup: deps docs build build-migrate
	@echo "Setup complete!"

# По умолчанию
all: setup
