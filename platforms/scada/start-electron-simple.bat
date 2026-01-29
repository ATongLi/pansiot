@echo off
cd /d %~dp0platforms\scada\packages\desktop
set ELECTRON_MIRROR=https://npmmirror.com/mirrors/electron/
echo Starting Scada Electron Desktop...
pnpm dev
