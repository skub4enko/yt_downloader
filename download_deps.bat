@echo off
echo ======================================
echo   Скачивание зависимостей
echo ======================================

REM Создаем папки если их нет
if not exist "bin" mkdir bin
if not exist "assets" mkdir assets

REM Проверяем и скачиваем yt-dlp.exe
if not exist "bin\yt-dlp.exe" (
    echo [1/2] Скачиваем yt-dlp.exe...
    powershell -Command "& {Invoke-WebRequest -Uri 'https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe' -OutFile 'bin\yt-dlp.exe'}"
    if %errorlevel% equ 0 (
        echo ✅ yt-dlp.exe скачан успешно
    ) else (
        echo ❌ Ошибка скачивания yt-dlp.exe
    )
) else (
    echo ✅ yt-dlp.exe уже существует
)

REM Проверяем и создаем простой beep файл если его нет
if not exist "assets\beep_long.wav" (
    echo [2/2] Создаем звуковой файл beep_long.wav...
    REM Создаем минимальный WAV файл (тишина 0.5 сек)
    powershell -Command "& {[System.IO.File]::WriteAllBytes('assets\beep_long.wav', @(0x52,0x49,0x46,0x46,0x24,0x08,0x00,0x00,0x57,0x41,0x56,0x45,0x66,0x6D,0x74,0x20,0x10,0x00,0x00,0x00,0x01,0x00,0x02,0x00,0x44,0xAC,0x00,0x00,0x10,0xB1,0x02,0x00,0x04,0x00,0x10,0x00,0x64,0x61,0x74,0x61,0x00,0x08,0x00,0x00) + @(0x00) * 2048)}"
    echo ✅ beep_long.wav создан
) else (
    echo ✅ beep_long.wav уже существует
)

REM Создаем пустой links.txt если его нет
if not exist "links.txt" (
    echo # Добавьте URL YouTube видео, по одному на строку > links.txt
    echo # Пример: >> links.txt
    echo # https://www.youtube.com/watch?v=dQw4w9WgXcQ >> links.txt
    echo ✅ links.txt создан
) else (
    echo ✅ links.txt уже существует
)

echo.
echo ======================================
echo    Все зависимости готовы!
echo ======================================
echo.
echo Структура проекта:
echo - bin\yt-dlp.exe (скачивание видео)
echo - assets\beep_long.wav (звуковой сигнал)  
echo - links.txt (для пакетной загрузки)
echo.
echo Теперь можно запустить build.bat для сборки
pause