@echo off
echo ===============================================
echo Starting PanTools Scada Desktop Application
echo ===============================================
echo.

echo Step 1: Setting Electron mirror...
set ELECTRON_MIRROR=https://npmmirror.com/mirrors/electron/

echo Step 2: Starting Electron...
echo.
pnpm dev

pause
