# Crypto Watcher

Микросервис для мониторинга и хранения стоимости криптовалют.

## Функциональность

- Добавление криптовалют в список наблюдения (`/currency/add`)
- Удаление криптовалют из списка наблюдения (`/currency/remove`)
- Получение цены криптовалюты на определенный момент времени (`/currency/price`)

## Запуск

### С помощью Docker Compose

1. Убедитесь, что у вас установлены Docker и Docker Compose
2. Клонируйте репозиторий
3. Запустите проект:

```bash
docker-compose up --build
```

Сервис будет доступен по адресу `http://localhost:8080`

### Локальный запуск

1. Установите PostgreSQL и создайте базу данных
2. Настройте переменные окружения в файле `.env`
3. Установите зависимости:

```bash
go mod tidy
```

4. Запустите сервис:

```bash
go run cmd/main.go
```

## API Документация

После запуска сервиса документация Swagger будет доступна по адресу:
`http://localhost:8080/swagger/index.html`

## Примеры использования

### Добавление криптовалюты в watchlist

```bash
curl -X POST http://localhost:8080/currency/add \
  -H "Content-Type: application/json" \
  -d '{"coin": "BTC"}'
```

### Удаление криптовалюты из watchlist

```bash
curl -X DELETE http://localhost:8080/currency/remove \
  -H "Content-Type: application/json" \
  -d '{"coin": "BTC"}'
```

### Получение цены

```bash
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=1736500490"
```

## Архитектура

Проект построен с использованием чистой архитектуры:

- `cmd/` - точка входа в приложение
- `internal/api/` - HTTP обработчики
- `internal/service/` - бизнес-логика
- `internal/storage/` - работа с базой данных
- `internal/models/` - модели данных
- `internal/config/` - конфигурация
- `internal/database/` - подключение к БД

## База данных

Используется PostgreSQL с двумя таблицами:
- `watchlist` - список отслеживаемых криптовалют
- `price_history` - история цен

## Технологии

- Go 1.21
- Gin (HTTP framework)
- PostgreSQL
- Docker & Docker Compose
- Swagger для документации API
