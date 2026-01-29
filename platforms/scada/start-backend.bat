@echo off
cd /d %~dp0platforms\scada\backend
echo Starting Scada Backend...
go run main.go
