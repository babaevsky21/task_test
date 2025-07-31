@echo off
echo ========================================
echo      ТЕСТИРОВАНИЕ DOCKER COMPOSE
echo ========================================

echo.
echo 1. Проверяем наличие Docker...
docker --version
if %errorlevel% neq 0 (
    echo ❌ Docker не найден! Установите Docker Desktop.
    pause
    exit /b 1
)

echo.
echo 2. Проверяем наличие Docker Compose...
docker-compose --version
if %errorlevel% neq 0 (
    echo ❌ Docker Compose не найден!
    pause
    exit /b 1
)

echo.
echo 3. Останавливаем существующие контейнеры...
docker-compose down

echo.
echo 4. Удаляем старые образы (опционально)...
docker-compose down --rmi all --volumes --remove-orphans

echo.
echo 5. Собираем и запускаем сервисы...
docker-compose up --build -d

echo.
echo 6. Проверяем статус контейнеров...
docker-compose ps

echo.
echo 7. Ожидание запуска сервисов (30 секунд)...
timeout /t 30 /nobreak

echo.
echo 8. Проверяем логи приложения...
docker-compose logs app

echo.
echo 9. Тестируем API...
echo.
echo 9.1 Добавляем BTC:
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"BTC\"}"

echo.
echo.
echo 9.2 Добавляем ETH:
curl -X POST "http://localhost:8080/currency/add" -H "Content-Type: application/json" -d "{\"coin\": \"ETH\"}"

echo.
echo.
echo 9.3 Ждем сбор цен (15 секунд)...
timeout /t 15 /nobreak

echo.
echo 9.4 Получаем цену BTC:
curl "http://localhost:8080/currency/price?coin=BTC&timestamp=1753954000"

echo.
echo.
echo 10. Проверяем логи базы данных...
docker-compose logs postgres

echo.
echo.
echo ========================================
echo           ТЕСТИРОВАНИЕ ЗАВЕРШЕНО
echo ========================================
echo.
echo Доступные команды:
echo - docker-compose logs app      # Логи приложения
echo - docker-compose logs postgres # Логи БД
echo - docker-compose down         # Остановить все
echo - http://localhost:8080/swagger/index.html # Swagger UI
echo.

pause
