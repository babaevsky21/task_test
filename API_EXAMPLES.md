# API Examples - Crypto Watcher

Примеры использования всех API эндпоинтов микросервиса Crypto Watcher.

**Базовый URL:** `http://localhost:8080`

## 📋 Содержание

1. [Добавление криптовалюты (POST)](#добавление-криптовалюты-post)
2. [Удаление криптовалюты (DELETE)](#удаление-криптовалюты-delete)
3. [Получение цены (GET)](#получение-цены-get)
4. [Swagger документация](#swagger-документация)
5. [Примеры ответов](#примеры-ответов)

---

## 1. Добавление криптовалюты (POST)

### Эндпоинт: `POST /currency/add`

Добавляет криптовалюту в список наблюдения (watchlist).

### cURL пример:
```bash
curl -X POST "http://localhost:8080/currency/add" \
  -H "Content-Type: application/json" \
  -d '{"coin": "BTC"}'
```

### PowerShell пример:
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/currency/add" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"coin": "BTC"}'
```

### JavaScript (fetch) пример:
```javascript
fetch('http://localhost:8080/currency/add', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    coin: 'BTC'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

### Python пример:
```python
import requests

url = "http://localhost:8080/currency/add"
payload = {"coin": "BTC"}
headers = {"Content-Type": "application/json"}

response = requests.post(url, json=payload, headers=headers)
print(response.json())
```

### Успешный ответ:
```json
{
  "message": "Coin added successfully"
}
```

### Возможные ошибки:
```json
{
  "error": "Key: 'AddCoinRequest.Coin' Error:Field validation for 'Coin' failed on the 'required' tag"
}
```

---

## 2. Удаление криптовалюты (DELETE)

### Эндпоинт: `DELETE /currency/remove`

Удаляет криптовалюту из списка наблюдения.

### cURL пример:
```bash
curl -X DELETE "http://localhost:8080/currency/remove" \
  -H "Content-Type: application/json" \
  -d '{"coin": "BTC"}'
```

### PowerShell пример:
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/currency/remove" `
  -Method DELETE `
  -ContentType "application/json" `
  -Body '{"coin": "BTC"}'
```

### JavaScript (fetch) пример:
```javascript
fetch('http://localhost:8080/currency/remove', {
  method: 'DELETE',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    coin: 'BTC'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

### Python пример:
```python
import requests

url = "http://localhost:8080/currency/remove"
payload = {"coin": "BTC"}
headers = {"Content-Type": "application/json"}

response = requests.delete(url, json=payload, headers=headers)
print(response.json())
```

### Успешный ответ:
```json
{
  "message": "Coin removed successfully"
}
```

### Возможные ошибки:
```json
{
  "error": "coin not found in watchlist"
}
```

---

## 3. Получение цены (GET)

### Эндпоинт: `GET /currency/price`

Получает цену криптовалюты на определенный момент времени.

### Параметры:
- `coin` - название криптовалюты (обязательный)
- `timestamp` - Unix timestamp (обязательный)

### cURL примеры:
```bash
# Получить цену BTC
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=1753954000"

# Получить цену ETH
curl "http://localhost:8080/currency/price?coin=ETH&timestamp=1753954000"

# Текущее время (Linux/Mac)
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=$(date +%s)"
```

### PowerShell пример:
```powershell
# Получить цену BTC
Invoke-RestMethod -Uri "http://localhost:8080/currency/price?coin=BTC&timestamp=1753954000"

# С текущим временем
$timestamp = [int][double]::Parse((Get-Date -UFormat %s))
Invoke-RestMethod -Uri "http://localhost:8080/currency/price?coin=BTC&timestamp=$timestamp"
```

### JavaScript (fetch) пример:
```javascript
// С определенным timestamp
fetch('http://localhost:8080/currency/price?coin=BTC&timestamp=1753954000')
  .then(response => response.json())
  .then(data => console.log(data));

// С текущим временем
const timestamp = Math.floor(Date.now() / 1000);
fetch(`http://localhost:8080/currency/price?coin=BTC&timestamp=${timestamp}`)
  .then(response => response.json())
  .then(data => console.log(data));
```

### Python пример:
```python
import requests
import time

# С определенным timestamp
url = "http://localhost:8080/currency/price"
params = {
    "coin": "BTC",
    "timestamp": 1753954000
}

response = requests.get(url, params=params)
print(response.json())

# С текущим временем
params = {
    "coin": "BTC",
    "timestamp": int(time.time())
}

response = requests.get(url, params=params)
print(response.json())
```

### Успешный ответ:
```json
{
  "id": 1,
  "coin": "BTC",
  "price": 45994.0,
  "timestamp": 1753953994,
  "created_at": "2025-07-31T12:26:34.835323Z"
}
```

### Возможные ошибки:
```json
{
  "error": "coin not found in watchlist"
}
```

```json
{
  "error": "price not found for the given timestamp"
}
```

```json
{
  "error": "coin and timestamp parameters are required"
}
```

---

## 4. Swagger документация

### URL: `http://localhost:8080/swagger/index.html`

Интерактивная документация API с возможностью тестирования прямо в браузере.

---

## 5. Примеры ответов

### Структура успешного ответа для цены:
```json
{
  "id": 1,              // ID записи в БД
  "coin": "BTC",        // Название криптовалюты
  "price": 45994.0,     // Цена в USD
  "timestamp": 1753953994, // Unix timestamp
  "created_at": "2025-07-31T12:26:34.835323Z" // Время записи в БД
}
```

### Структура ошибки:
```json
{
  "error": "описание ошибки"
}
```

---

## 🧪 Быстрое тестирование

### Полный цикл тестирования:
```bash
# 1. Добавляем BTC
curl -X POST "http://localhost:8080/currency/add" \
  -H "Content-Type: application/json" \
  -d '{"coin": "BTC"}'

# 2. Ждем несколько секунд для автоматического сбора цен
sleep 15

# 3. Получаем цену
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=$(date +%s)"

# 4. Удаляем BTC
curl -X DELETE "http://localhost:8080/currency/remove" \
  -H "Content-Type: application/json" \
  -d '{"coin": "BTC"}'
```

### Использование готовых скриптов:
```bash
# Автоматическое тестирование всех эндпоинтов
./test_data_insertion.bat

# Или через Makefile
make test-api
```

---

## 📝 Примечания

1. **Автоматический сбор цен:** После добавления монеты в watchlist, сервис автоматически начинает собирать цены каждые 10 секунд.

2. **Временные метки:** Используйте Unix timestamp в секундах. Сервис вернет цену, ближайшую к указанному времени.

3. **Поддерживаемые монеты:** Любые строки до 10 символов. Примеры: BTC, ETH, DOGE, ADA, LTC.

4. **CORS:** В текущей версии CORS не настроен. Для фронтенд-приложений может потребоваться настройка.

5. **Лимиты:** В текущей версии лимиты запросов не установлены.
