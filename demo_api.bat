@echo off
echo ======================================
echo       DEMO API CRYPTO WATCHER
echo ======================================
echo.

:menu
echo Выберите действие:
echo 1. Добавить монету в watchlist
echo 2. Удалить монету из watchlist  
echo 3. Получить цену монеты
echo 4. Показать все примеры
echo 5. Выход
echo.
set /p choice="Введите номер (1-5): "

if "%choice%"=="1" goto add_coin
if "%choice%"=="2" goto remove_coin
if "%choice%"=="3" goto get_price
if "%choice%"=="4" goto show_examples
if "%choice%"=="5" goto exit
goto menu

:add_coin
echo.
set /p coin="Введите название монеты (например BTC, ETH, DOGE): "
echo Добавляем %coin%...
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"%coin%\"}"
echo.
echo.
pause
goto menu

:remove_coin
echo.
set /p coin="Введите название монеты для удаления: "
echo Удаляем %coin%...
curl -X DELETE "http://localhost:8080/currency/remove" -H "Content-Type: application/json" -d "{\"coin\": \"%coin%\"}"
echo.
echo.
pause
goto menu

:get_price
echo.
set /p coin="Введите название монеты: "
echo Получаем текущую цену для %coin%...
echo (используется текущий timestamp)
curl "http://localhost:8080/currency/price?coin=%coin%&timestamp=1753954400"
echo.
echo.
pause
goto menu

:show_examples
echo.
echo ======================================
echo            ПРИМЕРЫ ЗАПРОСОВ
echo ======================================
echo.
echo 1. POST - Добавление BTC:
echo curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"BTC\"}"
echo.
echo 2. DELETE - Удаление BTC:
echo curl -X DELETE "http://localhost:8080/currency/remove" -H "Content-Type: application/json" -d "{\"coin\": \"BTC\"}"
echo.
echo 3. GET - Получение цены:
echo curl "http://localhost:8080/currency/price?coin=BTC&timestamp=1753954400"
echo.
echo 4. Swagger UI:
echo http://localhost:8080/swagger/index.html
echo.
echo ======================================
pause
goto menu

:exit
echo До свидания!
exit
