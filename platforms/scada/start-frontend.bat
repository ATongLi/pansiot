@echo off
cd /d %~dp0platforms\scada\packages\renderer
echo Starting Scada Frontend...
pnpm dev
