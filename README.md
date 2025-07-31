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

**Ответ:** `{"message":"Coin added successfully"}`

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

**Ответ:**
```json
{
  "id": 1,
  "coin": "BTC",
  "price": 45000.50,
  "timestamp": 1736500490,
  "created_at": "2025-07-30T20:00:00Z"
}
```

### Тестирование API

Примеры тестирования через PowerShell:

```powershell
# Добавление монеты
Invoke-RestMethod -Uri "http://localhost:8080/currency/add" -Method POST -ContentType "application/json" -Body '{"coin": "BTC"}'

# Получение цены
$timestamp = [int][double]::Parse((Get-Date -UFormat %s))
Invoke-RestMethod -Uri "http://localhost:8080/currency/price?coin=BTC&timestamp=$timestamp" -Method GET

# Удаление монеты
Invoke-RestMethod -Uri "http://localhost:8080/currency/remove" -Method DELETE -ContentType "application/json" -Body '{"coin": "BTC"}'
```

## Архитектура

Проект построен с использованием чистой архитектуры:

- `cmd/` - точка входа в приложение
  - `cmd/main.go` - основное приложение
  - `cmd/migrate/` - CLI утилита миграций
- `internal/api/` - HTTP обработчики
- `internal/service/` - бизнес-логика
- `internal/storage/` - работа с базой данных
- `internal/models/` - модели данных
- `internal/config/` - конфигурация
- `internal/database/` - подключение к БД и миграции

## База данных

Используется PostgreSQL с автоматической системой миграций:

### Таблицы:
- `watchlist` - список отслеживаемых криптовалют
- `price_history` - история цен
- `schema_migrations` - отслеживание примененных миграций

### Автоматические миграции:
При запуске приложения автоматически:
- Создаются необходимые таблицы
- Применяются новые миграции
- Создаются индексы для оптимизации

### Управление миграциями:
```bash
# Проверка статуса миграций
go run cmd/migrate/main.go status

# Принудительное выполнение миграций
go run cmd/migrate/main.go migrate
```

## Технологии

- Go 1.21
- Gin (HTTP framework)
- PostgreSQL
- Docker & Docker Compose
- Swagger для документации API
