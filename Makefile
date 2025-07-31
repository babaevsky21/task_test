.PHONY: build run test clean docs docker-build docker-run

# Переменные
APP_NAME := crypto-watcher
BUILD_DIR := bin
DOCKER_IMAGE := $(APP_NAME):latest

# Сборка приложения
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/main.go
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

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

# Запуск тестов
test:
	@echo "Running tests..."
	@go test -v ./...

# Тестирование подключения к БД
test-db:
	@echo "Testing database connection..."
	@go run test_db_connection.go

# Инициализация базы данных
init-db:
	@echo "Initializing database..."
	@psql -h localhost -p 5436 -U postgres -d postgres -f setup_database.sql

# Проверка PostgreSQL на разных портах
check-postgres:
	@echo "Checking PostgreSQL on different ports..."
	@./check_postgres.bat

# Тестирование API
test-api:
	@echo "Testing API endpoints..."
	@./test_data_insertion.bat

# Интерактивная демонстрация API
demo:
	@echo "Starting interactive API demo..."
	@./demo_api.bat

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

# Быстрое тестирование Docker
docker-test:
	@echo "Quick Docker Compose test..."
	@./quick_docker_test.bat

# Полное тестирование Docker
docker-test-full:
	@echo "Full Docker Compose test..."
	@./test_docker_compose.bat

# Полная настройка проекта
setup: deps docs build
	@echo "Setup complete!"

# По умолчанию
all: setup
