@echo off
echo Проверка данных в базе данных PostgreSQL
echo =========================================

echo.
echo 1. Проверяем watchlist:
psql -h localhost -p 5432 -U postgres -d task -c "SELECT coin, created_at FROM watchlist ORDER BY created_at;"

echo.
echo 2. Проверяем price_history (последние 10 записей):
psql -h localhost -p 5432 -U postgres -d task -c "SELECT coin, price, to_timestamp(timestamp), created_at FROM price_history ORDER BY created_at DESC LIMIT 10;"

echo.
echo 3. Статистика по монетам:
psql -h localhost -p 5432 -U postgres -d task -c "SELECT coin, COUNT(*) as records, MIN(price) as min_price, MAX(price) as max_price, ROUND(AVG(price)::numeric, 2) as avg_price FROM price_history GROUP BY coin ORDER BY records DESC;"

echo.
echo 4. Общее количество записей:
psql -h localhost -p 5432 -U postgres -d task -c "SELECT 'watchlist' as table_name, COUNT(*) as count FROM watchlist UNION ALL SELECT 'price_history' as table_name, COUNT(*) as count FROM price_history;"

pause
