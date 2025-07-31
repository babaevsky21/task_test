@echo off
echo Тестирование API и проверка данных
echo ==================================

echo.
echo 1. Добавляем DOGE в watchlist:
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"DOGE\"}"

echo.
echo.
echo 2. Добавляем LTC в watchlist:
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"LTC\"}"

echo.
echo.
echo 3. Пытаемся получить цену BTC (должна быть в БД):
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=1753954000"

echo.
echo.
echo 4. Пытаемся получить цену ETH (должна быть в БД):
curl "http://localhost:8080/currency/price?coin=ETH&timestamp=1753954000"

echo.
echo.
echo 5. Пытаемся получить цену DOGE (только что добавили, но цены может не быть):
curl "http://localhost:8080/currency/price?coin=DOGE&timestamp=1753954000"

echo.
echo.
echo 6. Пытаемся получить цену несуществующей монеты:
curl "http://localhost:8080/currency/price?coin=NONEXISTENT&timestamp=1753954000"

echo.
echo.
echo ✅ Тестирование завершено!
echo Если видите успешные ответы, значит данные корректно добавляются в БД.

pause
