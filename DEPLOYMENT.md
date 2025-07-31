# Инструкция по развертыванию Crypto Watcher

## Быстрый старт

### 1. Запуск с Docker Compose (Рекомендуется)

```bash
# Клонирование репозитория
git clone <repository-url>
cd crypto-watcher

# Запуск всех сервисов
docker-compose up --build
```

Сервис будет доступен по адресу: `http://localhost:8080`
Swagger документация: `http://localhost:8080/swagger/index.html`

### 2. Локальная разработка

#### Требования
- Go 1.21+
- PostgreSQL 15+
- Git

#### Шаги установки

1. **Установка зависимостей**
   ```bash
   go mod tidy
   ```

2. **Установка swag для документации**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

3. **Генерация Swagger документации**
   ```bash
   swag init -g cmd/main.go
   ```

4. **Настройка базы данных**
   - Создайте базу данных PostgreSQL с именем `crypto_watcher`
   - Выполните скрипт инициализации: `init.sql`

5. **Настройка переменных окружения**
   ```bash
   cp .env.production .env
   # Отредактируйте .env файл под ваши настройки
   ```

6. **Запуск приложения**
   ```bash
   go run cmd/main.go
   ```

#### Или используя Makefile

```bash
# Полная настройка
make setup

# Запуск
make run

# Сборка
make build

# Тесты
make test
```

## API Эндпоинты

### Добавление криптовалюты
```bash
POST /currency/add
Content-Type: application/json

{
  "coin": "BTC"
}
```

### Удаление криптовалюты
```bash
DELETE /currency/remove
Content-Type: application/json

{
  "coin": "BTC"
}
```

### Получение цены
```bash
GET /currency/price?coin=BTC&timestamp=1736500490
```

## Архитектура проекта

```
crypto-watcher/
├── cmd/               # Точка входа в приложение
│   └── main.go
├── internal/          # Внутренние пакеты
│   ├── api/          # HTTP обработчики
│   ├── config/       # Конфигурация
│   ├── database/     # Подключение к БД
│   ├── models/       # Модели данных
│   ├── service/      # Бизнес-логика
│   └── storage/      # Работа с хранилищем
├── docs/             # Swagger документация
├── docker-compose.yml # Docker Compose конфигурация
├── Dockerfile        # Docker образ
├── init.sql          # Инициализация БД
└── README.md         # Документация
```

## База данных

### Схема

**watchlist** - список отслеживаемых криптовалют
- `id` (SERIAL PRIMARY KEY)
- `coin` (VARCHAR(10) UNIQUE)
- `created_at` (TIMESTAMP)

**price_history** - история цен
- `id` (SERIAL PRIMARY KEY)
- `coin` (VARCHAR(10) FOREIGN KEY)
- `price` (DECIMAL(20, 8))
- `timestamp` (BIGINT)
- `created_at` (TIMESTAMP)

## Мониторинг и логирование

Приложение выводит логи в консоль. В продакшене рекомендуется:
- Настроить централизованное логирование
- Добавить метрики для мониторинга
- Настроить алерты

## Безопасность

- Используйте сильные пароли для базы данных
- В продакшене настройте HTTPS
- Ограничьте доступ к базе данных
- Регулярно обновляйте зависимости

## Масштабирование

Для увеличения нагрузки можно:
- Добавить кэширование (Redis)
- Настроить репликацию базы данных
- Использовать load balancer
- Добавить горизонтальное масштабирование

## Разработка

### Добавление новых функций

1. Добавьте модели в `internal/models/`
2. Обновите storage слой в `internal/storage/`
3. Добавьте бизнес-логику в `internal/service/`
4. Создайте API обработчики в `internal/api/`
5. Обновите документацию и тесты

### Тестирование

```bash
# Запуск всех тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Тесты конкретного пакета
go test ./internal/storage/
```
