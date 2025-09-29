# 🏗️ Build Instructions — YouTube Downloader

## 📋 Requirements

### To build:
- **Go 1.24.2+** — download from https://golang.org/dl/
- **Git** (optional) — to clone the repo

### To run:
- **Windows 10/11** (portable build)
- **Internet connection** (to download videos)
- **~50 MB free space** for app and dependencies

## 🚀 Quick build (Windows)

### Option 1: Automatic
```batch
REM 1) Download dependencies (yt-dlp, sample assets, links.txt)
download_deps.bat

REM 2) Build portable bundle into dist\ and create ZIP
build.bat
```

### Option 2: Manual
```batch
REM 1) Create folders
mkdir bin assets

REM 2) Download yt-dlp.exe to bin/
REM https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe

REM 3) Build Go binary
go mod tidy
go build -ldflags="-s -w" -o yt-downloader.exe .
```

Notes:
- We no longer use Windows resources (no resource.syso, no icon at runtime).
- Completion sounds are synthesized in code; WAV files are not required to play beeps.
- FFmpeg is required for merging streams and MP3 extraction. The project can auto-download FFmpeg to bin/.

## 🐧 Building on Linux/macOS (optional)

Using Makefile:
```bash
make all         # full build
make linux       # linux build
make windows     # cross-compile for Windows
make deps        # helper targets
```

Manual (example):
```bash
go mod tidy
mkdir -p bin
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o bin/yt-dlp
chmod +x bin/yt-dlp
go build -ldflags="-s -w" -o yt-downloader .
```

## 📁 Portable output layout

```
YouTube-Downloader-Portable/
├── yt-downloader.exe        # main app
├── bin/
│   └── yt-dlp.exe          # downloader utility
├── assets/                 # optional assets folder
├── links.txt               # batch URL list
└── README.txt              # usage (EN/ES)
```

## ⚙️ Build options

### Binary size (rough):
- Standard build: ~15–20 MB
- With -ldflags -s -w: ~8–12 MB
- With UPX: ~3–5 MB

### Useful flags
```bash
go build -ldflags="-s -w" -o yt-downloader.exe .     # smaller size
CGO_ENABLED=0 go build -ldflags="-s -w" -o yt-downloader.exe .  # fully static
go build -o yt-downloader.exe .                      # with debug info
```

### Cross-compilation
```bash
# From Linux to Windows
GOOS=windows GOARCH=amd64 go build -o yt-downloader.exe .

# From Windows to Linux
set GOOS=linux
set GOARCH=amd64
go build -o yt-downloader .
```

## 🧪 Verifying the build

### Check tools
```bash
go version
go mod verify
bin/yt-dlp.exe --version   # Windows
bin/yt-dlp --version       # Linux
bin/ffmpeg.exe -version    # Windows (if auto-downloaded)
```

### Run the app
```bash
./yt-downloader.exe   # Windows
./yt-downloader       # Linux
```
You should see the main menu with options:
- Audio (MP3)
- Video (MP4/WebM) including 2160p, 1440p, 1080p, 720p, 480p, 360p.

## 🔧 Troubleshooting

### "go: command not found"
- Install Go and add it to PATH; restart your terminal

### "yt-dlp.exe not found"
- Run `download_deps.bat` or manually place yt-dlp.exe into `bin/`

### "FFmpeg not found"
- Run `download_deps.bat` (it downloads ffmpeg.exe and ffprobe.exe into `bin/`)
- Or install FFmpeg manually and add to PATH, or place binaries in `bin/`
- The app passes `--ffmpeg-location bin` to yt-dlp

### Large .exe size
- Use `-ldflags="-s -w"` and optionally UPX: `upx --best yt-downloader.exe`

### Build errors
- Check Go version: `go version`
- Refresh modules: `go mod tidy`
- Clean cache: `go clean -modcache`

## 📋 Release checklist

- ✅ Go 1.24.2+ installed
- ✅ Project builds without errors
- ✅ `bin/yt-dlp.exe` present
- ✅ `bin/ffmpeg.exe` and `bin/ffprobe.exe` present (or FFmpeg in PATH)
- ✅ App launches and shows the menu
- ✅ Test video downloads work
- ✅ Portable archive created in project root

## 🎯 Result

You will get:
- **yt-downloader.exe** — the portable app
- **YouTube-Downloader-Portable.zip** — ready-to-share archive

No Windows resource embedding is used; sounds are synthesized at runtime.

---

# 🇷🇺 Инструкции по сборке (кратко)

## Требования
- Go 1.24.2+
- Windows 10/11 (для portable)
- Интернет (для скачивания видео)

## Быстрая сборка (Windows)
```
download_deps.bat   # загрузит yt-dlp и подготовит структуру
build.bat           # соберёт dist\ и ZIP
```

## Ручная сборка
```
mkdir bin assets
# скачайте yt-dlp.exe в bin/
go mod tidy
go build -ldflags="-s -w" -o yt-downloader.exe .
```

## Portable структура
```
YouTube-Downloader-Portable/
├─ yt-downloader.exe
├─ bin/yt-dlp.exe
├─ assets/            (опционально)
├─ links.txt
└─ README.txt (EN/ES)
```

## Запуск и опции
- Главное меню: Аудио (MP3) / Видео (MP4/WebM)
- Качество видео: 2160p, 1440p, 1080p, 720p, 480p, 360p
- Звуки завершения синтезируются кодом (файлы WAV не обязательны)

## Частые проблемы
- "go: command not found" — установите Go и добавьте в PATH
- "yt-dlp.exe not found" — запустите `download_deps.bat` или положите файл в `bin/`
- Большой размер exe — используйте `-ldflags="-s -w"`, при желании UPX

---

# 🇪🇸 Instrucciones de compilación (resumen)

## Requisitos
- Go 1.24.2+
- Windows 10/11 (portable)
- Conexión a internet (para descargar videos)

## Compilación rápida (Windows)
```
download_deps.bat   # descarga yt-dlp y prepara carpetas
build.bat           # crea dist\ y el ZIP portable
```

## Compilación manual
```
mkdir bin assets
# descargue yt-dlp.exe a bin/
go mod tidy
go build -ldflags="-s - w" -o yt-downloader.exe .
```

## Estructura portable
```
YouTube-Downloader-Portable/
├─ yt-downloader.exe
├─ bin/yt-dlp.exe
├─ assets/            (opcional)
├─ links.txt
└─ README.txt (EN/ES)
```

## Ejecución y opciones
- Menú principal: Audio (MP3) / Video (MP4/WebM)
- Calidad de video: 2160p, 1440p, 1080p, 720p, 480p, 360p
- Los sonidos de finalización se sintetizan en tiempo de ejecución (no requiere WAV)

## Problemas comunes
- "go: command not found": instale Go y añádalo al PATH
- "yt-dlp.exe not found": ejecute `download_deps.bat` o coloque el archivo en `bin/`
- Tamaño grande del .exe: use `-ldflags="-s - w"`, opcionalmente UPX