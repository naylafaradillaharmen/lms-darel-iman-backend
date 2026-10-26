@echo off
title PKBM Digital Dar El Iman - Backend API (Go Fiber)
cd /d "%~dp0"
echo ========================================================
echo   PKBM DIGITAL DAR EL IMAN - BACKEND API (GO FIBER v2)
echo   Listening on http://localhost:8080
echo ========================================================
echo.
if exist "bin\server.exe" (
    echo Menjalankan server Go Fiber dari bin\server.exe...
    bin\server.exe
) else (
    echo Menjalankan server Go Fiber via go run cmd/server/main.go...
    go run cmd/server/main.go
)
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo Terjadi kendala saat menjalankan backend.
    echo Coba jalankan manual: go run cmd/server/main.go
)
pause
