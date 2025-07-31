# Docker Deployment Guide

## 🐳 Запуск через Docker Compose

### Быстрый старт

```bash
# Клонирование репозитория
git clone <repository-url>
cd crypto-watcher

# Запуск всех сервисов
docker-compose up --build
```

### Проверка работы

После запуска сервисы будут доступны:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **PostgreSQL**: localhost:5432

## 📋 Доступные команды

### Основные команды

```bash
# Запуск в фоновом режиме
docker-compose up --build -d

# Просмотр логов
docker-compose logs -f app      # Логи приложения
docker-compose logs -f postgres # Логи базы данных

# Остановка сервисов
docker-compose down

# Полная очистка (удаление данных)
docker-compose down --volumes --remove-orphans
```

### Через Makefile

```bash
make docker-test        # Быстрое тестирование
make docker-up          # Запуск в фоне
make docker-stop        # Остановка
make docker-clean       # Полная очистка
```

### Автоматические скрипты

```bash
# Быстрая проверка Docker Compose
./quick_docker_test.bat

# Полное тестирование с API
./test_docker_compose.bat
```

## 🔧 Структура сервисов

### PostgreSQL
- **Образ**: postgres:15
- **База данных**: task
- **Порт**: 5432
- **Пользователь**: postgres
- **Пароль**: password

### Приложение
- **Сборка**: из локального Dockerfile
- **Порт**: 8080
- **Переменные окружения**: настроены для работы с PostgreSQL в Docker

## 🏗️ Архитектура контейнеров

```
┌─────────────────┐    ┌─────────────────┐
│   Application   │    │   PostgreSQL    │
│   (Go/Gin)      │────│   Database      │
│   Port: 8080    │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────────────────┘
               Docker Network
```

## 🔍 Диагностика

### Проверка статуса контейнеров

```bash
docker-compose ps
```

### Просмотр ресурсов

```bash
docker-compose top
```

### Подключение к контейнерам

```bash
# К приложению
docker-compose exec app sh

# К базе данных
docker-compose exec postgres psql -U postgres -d task
```

### Проверка сети

```bash
docker network ls
docker network inspect <network_name>
```

## 🐛 Решение проблем

### Проблема: Контейнер app падает с ошибкой подключения к БД

**Решение**: Добавлен healthcheck для PostgreSQL и depends_on с условием

### Проблема: Порт уже занят

```bash
# Найти процесс на порту 8080
netstat -ano | findstr :8080

# Остановить другие сервисы
docker-compose down
```

### Проблема: Данные не сохраняются

**Решение**: Используется именованный том `postgres_data`

### Проблема: Медленный запуск

**Решение**: Добавлены healthcheck'и для корректного ожидания готовности сервисов

## 📊 Мониторинг

### Просмотр логов в реальном времени

```bash
# Все сервисы
docker-compose logs -f

# Только приложение
docker-compose logs -f app

# Последние 50 строк
docker-compose logs --tail=50 app
```

### Метрики контейнеров

```bash
docker stats $(docker-compose ps -q)
```

## 🔄 Обновление

### Пересборка после изменений кода

```bash
docker-compose build --no-cache app
docker-compose up -d app
```

### Обновление базы данных

```bash
# Остановить приложение
docker-compose stop app

# Выполнить миграцию
docker-compose exec postgres psql -U postgres -d task -f /path/to/migration.sql

# Запустить приложение
docker-compose start app
```

## 🔒 Продакшн конфигурация

Для продакшна рекомендуется:

1. **Использовать переменные окружения из файла**:
```bash
# Создать .env файл с секретными данными
DB_PASSWORD=secure_password_here
```

2. **Настроить ограничения ресурсов**:
```yaml
services:
  app:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

3. **Добавить reverse proxy** (nginx)

4. **Настроить SSL/TLS**

5. **Настроить бэкапы БД**

## 📋 Checklist для развертывания

- [ ] Docker и Docker Compose установлены
- [ ] Порты 8080 и 5432 свободны
- [ ] Файлы docker-compose.yml, Dockerfile, init.sql присутствуют
- [ ] Переменные окружения настроены
- [ ] Запуск: `docker-compose up --build`
- [ ] Проверка: http://localhost:8080/swagger/index.html
- [ ] Тестирование API эндпоинтов
