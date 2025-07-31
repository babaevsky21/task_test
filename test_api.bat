@echo off
echo Тестирование API Crypto Watcher
echo ================================

echo.
echo 1. Добавление BTC в watchlist:
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"BTC\"}"

echo.
echo.
echo 2. Добавление ETH в watchlist:
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"ETH\"}"

echo.
echo.
echo 3. Получение цены BTC по timestamp:
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=1736500490"

echo.
echo.
echo 4. Получение цены ETH по timestamp:
curl "http://localhost:8080/currency/price?coin=ETH&timestamp=1736500490"

echo.
echo.
echo 5. Попытка получить цену несуществующей монеты:
curl "http://localhost:8080/currency/price?coin=DOGE&timestamp=1736500490"

echo.
echo.
echo 6. Swagger документация доступна по адресу:
echo http://localhost:8080/swagger/index.html

echo.
pause
