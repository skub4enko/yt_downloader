@echo off
echo ======================================
echo      YouTube Downloader Builder
echo ======================================

REM Check Go toolchain
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Download Go at https://golang.org/dl/
    pause
    exit /b 1
)

echo [1/6] Cleaning previous artifacts...
if exist "dist" rmdir /s /q dist
if exist "yt-downloader.exe" del "yt-downloader.exe"

echo [2/6] Creating folder structure...
mkdir dist
mkdir dist\bin
mkdir dist\assets

echo [3/6] Building Go project...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dist\yt-downloader.exe .
if %errorlevel% neq 0 (
    echo ERROR: Build failed!
    pause
    exit /b 1
)

echo [4/6] Copying external files...
REM Copy yt-dlp.exe
if exist "bin\yt-dlp.exe" (
    copy "bin\yt-dlp.exe" "dist\bin\" >nul
    echo     - yt-dlp.exe copied
) else (
    echo WARNING: bin\yt-dlp.exe not found, will be downloaded on first run
)

REM Copy beep file
if exist "assets\beep_long.wav" (
    copy "assets\beep_long.wav" "dist\assets\" >nul
    echo     - beep_long.wav copied
) else (
    echo WARNING: assets\beep_long.wav not found
)

REM Create empty links.txt
echo. > dist\links.txt
echo     - links.txt created

echo [5/6] Creating README file (EN and ES)...
echo # YouTube Downloader Portable > dist\README.txt
echo. >> dist\README.txt
echo EN: Portable version of YouTube Downloader >> dist\README.txt
echo. >> dist\README.txt
echo Files: >> dist\README.txt
echo - yt-downloader.exe - main application >> dist\README.txt
echo - bin\yt-dlp.exe - downloading utility >> dist\README.txt
echo - assets\beep_long.wav - completion sound >> dist\README.txt
echo - links.txt - batch download list >> dist\README.txt
echo. >> dist\README.txt
echo Usage: >> dist\README.txt
echo 1. Run yt-downloader.exe >> dist\README.txt
echo 2. For batch downloads, add URLs to links.txt >> dist\README.txt
echo. >> dist\README.txt
echo ES: Version portable del descargador de YouTube >> dist\README.txt
echo Archivos: >> dist\README.txt
echo - yt-downloader.exe - aplicacion principal >> dist\README.txt
echo - bin\yt-dlp.exe - utilidad de descarga >> dist\README.txt
echo - assets\beep_long.wav - sonido de finalizacion >> dist\README.txt
echo - links.txt - lista para descargas por lotes >> dist\README.txt
echo Uso: >> dist\README.txt
echo 1. Ejecute yt-downloader.exe >> dist\README.txt
echo 2. Para descargas por lotes, agregue URL en links.txt >> dist\README.txt

echo [6/6] Creating archive...
powershell Compress-Archive -Path "dist\*" -DestinationPath "YouTube-Downloader-Portable.zip" -Force
if %errorlevel% equ 0 (
    echo ✅ Archive created: YouTube-Downloader-Portable.zip
) else (
    echo WARNING: Failed to create archive, but files are ready in dist\
)

echo.
echo ======================================
echo           BUILD COMPLETED!
echo ======================================
echo.
echo Artifacts folder: dist\
echo Portable archive: YouTube-Downloader-Portable.zip
echo.
echo EXE file size:
for %%A in (dist\yt-downloader.exe) do echo - yt-downloader.exe: %%~zA bytes
echo.
pause