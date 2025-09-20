@echo off
echo ======================================
echo      YouTube Downloader Builder
echo ======================================

REM Проверяем наличие Go
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go не установлен или не найден в PATH
    echo Скачайте Go с https://golang.org/dl/
    pause
    exit /b 1
)

echo [1/6] Очистка старых файлов...
if exist "dist" rmdir /s /q dist
if exist "yt-downloader.exe" del "yt-downloader.exe"

echo [2/6] Создание структуры папок...
mkdir dist
mkdir dist\bin
mkdir dist\assets

echo [3/6] Компиляция Go проекта...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\yt-downloader.exe .
if %errorlevel% neq 0 (
    echo ERROR: Ошибка компиляции!
    pause
    exit /b 1
)

echo [4/6] Копирование внешних файлов...
REM Копируем yt-dlp.exe
if exist "bin\yt-dlp.exe" (
    copy "bin\yt-dlp.exe" "dist\bin\" >nul
    echo     - yt-dlp.exe скопирован
) else (
    echo WARNING: bin\yt-dlp.exe не найден, будет скачан при первом запуске
)

REM Копируем звуковой файл
if exist "assets\beep_long.wav" (
    copy "assets\beep_long.wav" "dist\assets\" >nul
    echo     - beep_long.wav скопирован
) else (
    echo WARNING: assets\beep_long.wav не найден
)

REM Создаем пустой links.txt
echo. > dist\links.txt
echo     - links.txt создан

echo [5/6] Создание README файла...
echo # YouTube Downloader Portable > dist\README.txt
echo. >> dist\README.txt
echo Portable версия YouTube Downloader >> dist\README.txt
echo. >> dist\README.txt
echo Файлы: >> dist\README.txt
echo - yt-downloader.exe - основная программа >> dist\README.txt
echo - bin\yt-dlp.exe - утилита для скачивания >> dist\README.txt
echo - assets\beep_long.wav - звуковой сигнал >> dist\README.txt
echo - links.txt - файл для пакетной загрузки >> dist\README.txt
echo. >> dist\README.txt
echo Использование: >> dist\README.txt
echo 1. Запустите yt-downloader.exe >> dist\README.txt
echo 2. Для пакетной загрузки добавьте URL в links.txt >> dist\README.txt

echo [6/6] Создание архива...
powershell Compress-Archive -Path "dist\*" -DestinationPath "YouTube-Downloader-Portable.zip" -Force
if %errorlevel% equ 0 (
    echo ✅ Архив создан: YouTube-Downloader-Portable.zip
) else (
    echo WARNING: Не удалось создать архив, но файлы готовы в папке dist\
)

echo.
echo ======================================
echo           СБОРКА ЗАВЕРШЕНА!
echo ======================================
echo.
echo Готовые файлы в папке: dist\
echo Portable архив: YouTube-Downloader-Portable.zip
echo.
echo Размер exe файла:
for %%A in (dist\yt-downloader.exe) do echo - yt-downloader.exe: %%~zA bytes
echo.
pause