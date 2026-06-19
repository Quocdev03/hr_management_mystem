@echo off
setlocal enabledelayedexpansion

set COMPOSE=docker compose -f docker-compose.dev.yml

echo.
echo   HR Management System - DEVELOPMENT
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
    echo [SUCCESS] .env created.
)

REM Start services
echo [INFO]    Starting services...
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

echo [INFO]    Running seed data...
%COMPOSE% run --rm seed
if %ERRORLEVEL% neq 0 (
    echo [ERROR]   Seed failed.
    exit /b 1
)
echo [SUCCESS] Seed completed.
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
echo [SUCCESS] Dev stack is running!
echo.
echo   Frontend : http://localhost:5173
echo   Backend  : http://localhost:8080/api/v1/health
echo.
echo   First time running? Run: start-dev.bat setup
echo.
echo   Stop   : start-dev.bat down
echo   Logs   : start-dev.bat logs
echo   Status : start-dev.bat status
echo ============================================================
exit /b 0

:END_SUCCESS_SETUP
echo.
echo ============================================================
echo [SUCCESS] Dev stack is running!
echo.
echo   Frontend : http://localhost:5173
echo   Backend  : http://localhost:8080/api/v1/health
echo.
echo   Stop   : start-dev.bat down
echo   Logs   : start-dev.bat logs
echo   Status : start-dev.bat status
echo ============================================================
exit /b 0
