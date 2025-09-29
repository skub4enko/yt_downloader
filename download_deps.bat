@echo off
echo ======================================
echo   Downloading dependencies
echo ======================================

REM Create folders if missing
if not exist "bin" mkdir bin
if not exist "assets" mkdir assets

REM Check and download yt-dlp.exe
if not exist "bin\yt-dlp.exe" (
    echo [1/2] Downloading yt-dlp.exe...
    powershell -Command "& {Invoke-WebRequest -Uri 'https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe' -OutFile 'bin\yt-dlp.exe'}"
    if %errorlevel% equ 0 (
        echo ✅ yt-dlp.exe downloaded successfully
    ) else (
        echo ❌ Failed to download yt-dlp.exe
    )
) else (
    echo ✅ yt-dlp.exe already exists
)

REM Check and download FFmpeg (ffmpeg.exe and ffprobe.exe)
if not exist "bin\ffmpeg.exe" (
    echo [FFmpeg] Downloading static build...
    powershell -Command "& { $ProgressPreference='SilentlyContinue'; $tmp=New-Item -ItemType Directory -Path $env:TEMP\ffmpeg -Force; Invoke-WebRequest -Uri 'https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip' -OutFile $tmp\ffmpeg.zip; Expand-Archive -Path $tmp\ffmpeg.zip -DestinationPath $tmp -Force; $dir=(Get-ChildItem $tmp | Where-Object {$_.PSIsContainer -and $_.Name -like 'ffmpeg-*'} | Select-Object -First 1).FullName; Copy-Item "$dir\bin\ffmpeg.exe" bin\ -Force; Copy-Item "$dir\bin\ffprobe.exe" bin\ -Force }"
    if %errorlevel% equ 0 (
        echo ✅ FFmpeg downloaded to bin\
    ) else (
        echo ❌ Failed to download FFmpeg. Please install manually and ensure ffmpeg.exe is in PATH or bin\
    )
) else (
    echo ✅ FFmpeg already exists
)

REM Create simple beep files if missing
if not exist "assets\beep_long.wav" (
    echo [2/2] Creating beep_long.wav...
    REM Minimal WAV (0.5s silence)
    powershell -Command "& {[System.IO.File]::WriteAllBytes('assets\beep_long.wav', @(0x52,0x49,0x46,0x46,0x24,0x08,0x00,0x00,0x57,0x41,0x56,0x45,0x66,0x6D,0x74,0x20,0x10,0x00,0x00,0x00,0x01,0x00,0x02,0x00,0x44,0xAC,0x00,0x00,0x10,0xB1,0x02,0x00,0x04,0x00,0x10,0x00,0x64,0x61,0x74,0x61,0x00,0x08,0x00,0x00) + @(0x00) * 2048)}"
    echo ✅ beep_long.wav created
) else (
    echo ✅ beep_long.wav already exists
)

REM Create short beep file if missing
if not exist "assets\beep_short.wav" (
    echo [2/2] Creating beep_short.wav...
    powershell -Command "& {[System.IO.File]::WriteAllBytes('assets\beep_short.wav', @(0x52,0x49,0x46,0x46,0x24,0x04,0x00,0x00,0x57,0x41,0x56,0x45,0x66,0x6D,0x74,0x20,0x10,0x00,0x00,0x00,0x01,0x00,0x02,0x00,0x44,0xAC,0x00,0x00,0x10,0xB1,0x02,0x00,0x04,0x00,0x10,0x00,0x64,0x61,0x74,0x61,0x00,0x04,0x00,0x00) + @(0x00) * 1024)}"
    echo ✅ beep_short.wav created
) else (
    echo ✅ beep_short.wav already exists
)

REM Create empty links.txt if missing
if not exist "links.txt" (
    echo # Add YouTube URLs, one per line > links.txt
    echo # Example: >> links.txt
    echo # https://www.youtube.com/watch?v=dQw4w9WgXcQ >> links.txt
    echo ✅ links.txt created
) else (
    echo ✅ links.txt already exists
)

echo.
echo ======================================
echo    All dependencies are ready!
echo ======================================
echo.
echo Project structure:
echo - bin\yt-dlp.exe (downloading utility)
echo - assets\beep_long.wav (beep)  
echo - links.txt (batch list)
echo.
echo You can now run build.bat to build
pause