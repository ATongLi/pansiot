@echo off
echo ===============================================
echo PanTools Scada - Complete Development Environment
echo ===============================================
echo.

echo This script will open 3 separate windows:
echo   1. Backend Server (Go)
echo   2. Frontend Dev Server (Vite)
echo   3. Electron Desktop App
echo.

echo Press any key to start...
pause

echo.
echo [1/3] Starting Backend Server...
start "Scada Backend" cmd /k "cd /d %~dp0platforms\scada\backend && echo Running Go backend... && go run main.go"

echo.
echo [2/3] Starting Frontend Dev Server...
timeout /t 2 /nobreak >nul
start "Scada Frontend" cmd /k "cd /d %~dp0platforms\scada\packages\renderer && echo Starting Vite dev server... && pnpm dev"

echo.
echo [3/3] Starting Electron Desktop App...
timeout /t 5 /nobreak >nul
start "Scada Desktop" cmd /k "cd /d %~dp0platforms\scada\packages\desktop && set ELECTRON_MIRROR=https://npmmirror.com/mirrors/electron/ && echo Starting Electron... && pnpm dev"

echo.
echo ===============================================
echo All services started!
echo.
echo Backend:  http://localhost:3000
echo Frontend: http://localhost:5173
echo.
echo Press any key to close this window...
pause
