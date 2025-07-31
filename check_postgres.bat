@echo off
echo Проверка подключения к PostgreSQL на разных портах...

echo.
echo Попытка подключения к порту 5432 (стандартный):
psql -h localhost -p 5432 -U postgres -d postgres -c "SELECT version();"

echo.
echo Попытка подключения к порту 5436:
psql -h localhost -p 5436 -U postgres -d postgres -c "SELECT version();"

echo.
echo Попытка подключения к порту 5433:
psql -h localhost -p 5433 -U postgres -d postgres -c "SELECT version();"

echo.
echo Если подключение успешно, выполните:
echo psql -h localhost -p [ПОРТ] -U postgres -d postgres -f setup_database.sql

pause
