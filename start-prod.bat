@echo off
setlocal enabledelayedexpansion

set COMPOSE=docker compose -f docker-compose.prod.yml

echo.
echo   HR Management System - PRODUCTION
echo.

if /I "%1"=="down" goto DOWN
if /I "%1"=="-down" goto DOWN
if /I "%1"=="/down" goto DOWN

if /I "%1"=="logs" goto LOGS
if /I "%1"=="-logs" goto LOGS
if /I "%1"=="/logs" goto LOGS

if /I "%1"=="status" goto STATUS
if /I "%1"=="-status" goto STATUS
if /I "%1"=="/status" goto STATUS

REM Check Docker
docker info > nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR]   Docker is not running. Please start Docker Desktop first.
    exit /b 1
)

REM Create .env if it does not exist
if not exist .env (
    echo [INFO]    Copying .env.example to .env
    copy .env.example .env > nul
    echo.
    echo   [!] Please open .env file and fill in production values before continuing:
    echo       DB_PASSWORD, MYSQL_ROOT_PASSWORD, JWT_SECRET
    echo.
    exit /b 0
)

REM Check SSL certificate
if not exist certs\cert.pem (
    echo.
    echo   [!] SSL certificate not found.
    echo       Please run: bash scripts/generate-certs.sh
    echo.
)

REM Start services
echo [INFO]    Starting production stack...
%COMPOSE% up -d
if %ERRORLEVEL% neq 0 (
    echo [ERROR]   Failed to start services.
    exit /b 1
)
echo [SUCCESS] Services are running.

if /I "%1"=="setup" goto SETUP
if /I "%1"=="-setup" goto SETUP
if /I "%1"=="/setup" goto SETUP

goto END_SUCCESS

:SETUP
echo [INFO]    Waiting for MySQL to be ready (10s)...
timeout /t 10 /nobreak > nul

echo [INFO]    Running migrations...
%COMPOSE% run --rm migrate
if %ERRORLEVEL% neq 0 (
    echo [ERROR]   Migration failed.
    exit /b 1
)
echo [SUCCESS] Migration completed.

echo.
echo   [i] Seed is NOT run automatically in production.
echo       Please seed manually if needed.
goto END_SUCCESS_SETUP

:DOWN
%COMPOSE% down
echo [SUCCESS] Stopped.
exit /b 0

:LOGS
%COMPOSE% logs -f
exit /b 0

:STATUS
%COMPOSE% ps
exit /b 0

:END_SUCCESS
echo.
echo ============================================================
echo [SUCCESS] Production stack is running!
echo.
echo   App : http://localhost  (port 80)
echo         https://localhost (port 443)
echo.
echo   First time deploy? Run: start-prod.bat setup
echo.
echo   Stop   : start-prod.bat down
echo   Logs   : start-prod.bat logs
echo   Status : start-prod.bat status
echo ============================================================
exit /b 0

:END_SUCCESS_SETUP
echo.
echo ============================================================
echo [SUCCESS] Production stack is running!
echo.
echo   App : http://localhost  (port 80)
echo         https://localhost (port 443)
echo.
echo   Stop   : start-prod.bat down
echo   Logs   : start-prod.bat logs
echo   Status : start-prod.bat status
echo ============================================================
exit /b 0
