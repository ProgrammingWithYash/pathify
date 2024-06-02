@echo off
for /f "delims=" %%i in ('pathify.exe %*') do set newdir=%%i
cd /d "%newdir%"
