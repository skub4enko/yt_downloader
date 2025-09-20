# 🏗️ Инструкции по сборке YouTube Downloader

## 📋 Требования

### Для компиляции:
- **Go 1.24.2+** - [скачать с golang.org](https://golang.org/dl/)
- **Git** (опционально) - для клонирования репозитория

### Для работы программы:
- **Windows 10/11** - для portable версии
- **Интернет соединение** - для скачивания видео
- **~50MB свободного места** - для программы и зависимостей

## 🚀 Быстрая сборка (Windows)

### Вариант 1: Автоматическая сборка
```batch
# 1. Скачиваем зависимости
download_deps.bat

# 2. Собираем portable версию
build.bat
```

### Вариант 2: Ручная сборка
```batch
# 1. Создаем папки
mkdir bin assets

# 2. Скачиваем yt-dlp.exe в папку bin/
# Ссылка: https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe

# 3. Компилируем Go код
go mod tidy
go build -ldflags="-s -w" -o yt-downloader.exe .
```

## 🐧 Сборка для Linux/macOS

### Используя Makefile:
```bash
# Полная сборка
make all

# Только для Linux
make linux

# Только для Windows (кросс-компиляция)
make windows

# Скачать зависимости
make deps
make install-ytdlp
```

### Ручная сборка:
```bash
# 1. Загрузить зависимости
go mod tidy

# 2. Скачать yt-dlp
mkdir -p bin
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o bin/yt-dlp
chmod +x bin/yt-dlp

# 3. Компилировать
go build -ldflags="-s -w" -o yt-downloader .
```

## 📁 Структура готовой portable версии

```
YouTube-Downloader-Portable/
├── yt-downloader.exe          # Основная программа
├── bin/
│   └── yt-dlp.exe            # Утилита скачивания
├── assets/
│   └── beep_long.wav         # Звуковой сигнал
├── links.txt                 # Файл для URL списков
└── README.txt               # Инструкции
```

## ⚙️ Опции компиляции

### Размер файла:
- **Стандартная сборка**: ~15-20MB
- **С флагами оптимизации**: ~8-12MB
- **С UPX сжатием**: ~3-5MB

### Флаги сборки:
```bash
# Минимальный размер
go build -ldflags="-s -w" -o yt-downloader.exe .

# Статическая сборка (без внешних dll)
CGO_ENABLED=0 go build -ldflags="-s -w" -o yt-downloader.exe .

# С отладочной информацией
go build -o yt-downloader.exe .
```

### Кросс-компиляция:
```bash
# Для Windows из Linux
GOOS=windows GOARCH=amd64 go build -o yt-downloader.exe .

# Для Linux из Windows
set GOOS=linux
set GOARCH=amd64
go build -o yt-downloader .
```

## 🧪 Тестирование сборки

### Проверка зависимостей:
```bash
# Проверить что Go найден
go version

# Проверить модули
go mod verify

# Проверить что yt-dlp работает
bin/yt-dlp.exe --version    # Windows
bin/yt-dlp --version        # Linux
```

### Тест запуска:
```bash
# Запустить программу
./yt-downloader.exe    # Windows
./yt-downloader        # Linux

# Должно показать главное меню
```

## 📦 Создание installer (опционально)

### Для Windows (NSIS):
```nsis
# Можно создать installer используя NSIS
# Файл installer.nsi уже готов в репозитории
makensis installer.nsi
```

### Для Linux (AppImage):
```bash
# Создание AppImage
make linux
./create-appimage.sh dist/
```

## 🔧 Устранение проблем

### "go: command not found"
- Установите Go с официального сайта
- Добавьте Go в PATH
- Перезапустите командную строку

### "yt-dlp.exe not found"
- Запустите `download_deps.bat` 
- Или скачайте вручную в папку `bin/`

### Большой размер exe файла
- Используйте флаги `-ldflags="-s -w"`
- Установите UPX: `upx --best yt-downloader.exe`

### Ошибки при компиляции
- Проверьте версию Go: `go version`
- Обновите модули: `go mod tidy`
- Очистите кеш: `go clean -modcache`

## 📋 Checklist готовности

- ✅ Go 1.24.2+ установлен
- ✅ Проект скомпилирован без ошибок
- ✅ yt-dlp.exe находится в bin/
- ✅ beep_long.wav находится в assets/
- ✅ Программа запускается и показывает меню
- ✅ Можно скачать тестовое видео
- ✅ Все файлы упакованы в portable архив

## 🎯 Результат

После успешной сборки вы получите:
- **yt-downloader.exe** - готовая к работе программа
- **YouTube-Downloader-Portable.zip** - архив для распространения
- Программа работает на чистом Windows без установки дополнительных компонентов

Размер итогового архива: **~15-25MB**