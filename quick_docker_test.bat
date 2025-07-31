@echo off
echo ========================================
echo    БЫСТРАЯ ПРОВЕРКА DOCKER COMPOSE
echo ========================================

echo.
echo Проверяем Docker...
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker не запущен или не установлен
    echo 💡 Запустите Docker Desktop и повторите попытку
    pause
    exit /b 1
)

echo ✅ Docker найден

echo.
echo Проверяем файлы конфигурации...
if not exist docker-compose.yml (
    echo ❌ Файл docker-compose.yml не найден
    pause
    exit /b 1
)

if not exist Dockerfile (
    echo ❌ Файл Dockerfile не найден
    pause
    exit /b 1
)

if not exist init.sql (
    echo ❌ Файл init.sql не найден
    pause
    exit /b 1
)

echo ✅ Все файлы конфигурации найдены

echo.
echo Останавливаем существующие контейнеры...
docker-compose down --remove-orphans >nul 2>&1

echo.
echo Собираем образы...
docker-compose build --no-cache
if %errorlevel% neq 0 (
    echo ❌ Ошибка сборки образов
    pause
    exit /b 1
)

echo ✅ Образы собраны успешно

echo.
echo Запускаем сервисы...
docker-compose up -d
if %errorlevel% neq 0 (
    echo ❌ Ошибка запуска сервисов
    pause
    exit /b 1
)

echo ✅ Сервисы запущены

echo.
echo Ожидание готовности сервисов...
timeout /t 20 /nobreak >nul

echo.
echo Проверяем статус контейнеров...
docker-compose ps

echo.
echo Тестируем подключение к API...
curl -f http://localhost:8080/swagger/index.html >nul 2>&1
if %errorlevel% neq 0 (
    echo ⚠️  API пока не отвечает, проверяем логи...
    echo.
    echo === ЛОГИ ПРИЛОЖЕНИЯ ===
    docker-compose logs --tail=10 app
    echo.
    echo === ЛОГИ БАЗЫ ДАННЫХ ===
    docker-compose logs --tail=10 postgres
) else (
    echo ✅ API работает!
)

echo.
echo ========================================
echo             РЕЗУЛЬТАТ
echo ========================================
echo.
echo 🌐 Swagger UI: http://localhost:8080/swagger/index.html
echo 📊 API Base URL: http://localhost:8080
echo.
echo 📋 Полезные команды:
echo   docker-compose logs app      - логи приложения
echo   docker-compose logs postgres - логи БД
echo   docker-compose down          - остановить все
echo   docker-compose restart       - перезапустить
echo.

pause
