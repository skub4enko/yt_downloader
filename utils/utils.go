//go:build windows

package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/oto/v2"
)

// =================== YouTube Downloader ===================

// GetVideoTitle извлекает название видео с YouTube
func GetVideoTitle(url string) string {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	
	// ИСПРАВЛЕНИЕ: добавляем флаги для правильной кодировки
	cmd := exec.Command(ytPath, "--quiet", "--get-title", "--encoding", "utf-8", url)
	
	// Устанавливаем кодировку для Windows
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("⚠ Failed to get title:", err)
		return ""
	}
	
	title := strings.TrimSpace(string(output))
	
	// ОТЛАДКА: показываем что получили от yt-dlp
	fmt.Printf("🔍 Название от yt-dlp: '%s'\n", title)
	fmt.Printf("📏 Длина: %d байт, UTF-8 валидно: %t\n", len(title), utf8.ValidString(title))
	
	return title
}

// =================== Файловые утилиты ===================

// GenerateFallbackTitle генерирует fallback-название
func GenerateFallbackTitle() string {
	return "video_" + fmt.Sprintf("%d", time.Now().Unix())
}

// SanitizeFileName очищает название от нежелательных символов
// ИСПРАВЛЕНИЕ: НЕ удаляем кириллицу, только действительно опасные символы!
func SanitizeFileName(name string) string {
	if name == "" {
		return GenerateFallbackTitle()
	}
	
	original := name
	
	// СТАРАЯ ВЕРСИЯ (НЕПРАВИЛЬНО):
	// re := regexp.MustCompile("[^a-zA-Z0-9-_.]+")  // <-- ЭТО УБИВАЛО КИРИЛЛИЦУ!
	// return re.ReplaceAllString(name, "_")
	
	// НОВАЯ ВЕРСИЯ (ПРАВИЛЬНО):
	// Заменяем только символы, которые нельзя использовать в именах файлов
	dangerousChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = dangerousChars.ReplaceAllString(name, "_")
	
	// Удаляем управляющие символы (но не обычные!)
	controlChars := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	name = controlChars.ReplaceAllString(name, "")
	
	// Заменяем множественные пробелы на один
	multipleSpaces := regexp.MustCompile(`\s+`)
	name = multipleSpaces.ReplaceAllString(name, " ")
	
	// Убираем пробелы по краям
	name = strings.TrimSpace(name)
	
	// Ограничиваем длину (оставляем место для расширения)
	if len(name) > 200 {
		name = truncateUTF8(name, 200)
	}
	
	// Обрабатываем зарезервированные имена Windows
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upper := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upper == reserved || strings.HasPrefix(upper, reserved+".") {
			name = "_" + name
			break
		}
	}
	
	// Если название стало пустым после очистки
	if name == "" {
		name = GenerateFallbackTitle()
	}
	
	// ОТЛАДКА: показываем что изменилось
	if original != name {
		fmt.Printf("📝 Санитизация: '%s' -> '%s'\n", original, name)
	} else {
		fmt.Printf("✅ Название сохранено: '%s'\n", name)
	}
	
	return name
}

// truncateUTF8 корректно обрезает UTF-8 строку
func truncateUTF8(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	
	for i := maxLen; i >= 0; i-- {
		if utf8.ValidString(s[:i]) {
			return s[:i]
		}
	}
	
	return s[:maxLen]
}

// ParseURLsFromFile парсит URL из текстового содержимого файла
func ParseURLsFromFile(content string) []string {
	lines := strings.Split(content, "\n")
	var urls []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		
		// Базовая проверка что это похоже на URL
		if strings.Contains(line, "youtube.com") || 
		   strings.Contains(line, "youtu.be") ||
		   strings.Contains(line, "http://") ||
		   strings.Contains(line, "https://") {
			urls = append(urls, line)
		}
	}
	
	return urls
}

// IsValidURL проверяет является ли строка валидным URL
func IsValidURL(url string) bool {
	return strings.Contains(url, "youtube.com") || 
		   strings.Contains(url, "youtu.be") ||
		   strings.HasPrefix(url, "http://") ||
		   strings.HasPrefix(url, "https://")
}

// =================== WAV-плеер ===================

// intBufferToBytes конвертирует *audio.IntBuffer в []byte (16-bit little endian)
func intBufferToBytes(buf *audio.IntBuffer) []byte {
	out := make([]byte, len(buf.Data)*2)
	for i, v := range buf.Data {
		binary.LittleEndian.PutUint16(out[i*2:], uint16(int16(v)))
	}
	return out
}

// PlayBeep воспроизводит звук из файла assets/beep_long.wav
func PlayBeep() {
	f, err := os.Open("assets/beep_long.wav")
	if err != nil {
		fmt.Println("⚠ Не удалось открыть beep_long.wav:", err)
		return
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		fmt.Println("⚠ Неверный WAV файл")
		return
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		fmt.Println("⚠ Ошибка декодирования WAV:", err)
		return
	}

	ctx, ready, err := oto.NewContext(int(buf.Format.SampleRate), buf.Format.NumChannels, 2)
	if err != nil {
		fmt.Println("⚠ Ошибка инициализации аудио:", err)
		return
	}
	<-ready

	player := ctx.NewPlayer(bytes.NewReader(intBufferToBytes(buf)))
	defer player.Close()

	player.Play()

	// Ждем окончания воспроизведения
	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}
}

// =================== Автообновление yt-dlp ===================

// UpdateYtDlp обновляет yt-dlp.exe в папке bin
func UpdateYtDlp() {
	ytPath := filepath.Join("bin", "yt-dlp.exe")

	cmd := exec.Command(ytPath, "-U")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("⚠ Ошибка обновления yt-dlp:", err)
		fmt.Println(string(output))
		return
	}
	fmt.Println("✅ yt-dlp обновлён")
	fmt.Println(string(output))
}

// CheckUpdateYtDlp проверяет наличие новой версии yt-dlp
func CheckUpdateYtDlp() {
	ytPath := filepath.Join("bin", "yt-dlp.exe")

	// Проверка, существует ли локальный бинарь
	if _, err := os.Stat(ytPath); os.IsNotExist(err) {
		fmt.Println("⚠ bin/yt-dlp.exe не найден, скачиваю последнюю версию...")
		downloadYtDlp(ytPath)
		return
	}

	// Получаем локальную версию
	cmd := exec.Command(ytPath, "--version")
	currentVerBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("⚠ Не удалось получить локальную версию yt-dlp:", err)
		return
	}
	currentVer := strings.TrimSpace(string(currentVerBytes))

	// Получаем последнюю версию с GitHub
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		fmt.Println("⚠ Не удалось проверить последнюю версию:", err)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	latestVer, _ := data["tag_name"].(string)

	if latestVer != "" && latestVer != currentVer {
		fmt.Println("⬆ Доступен новый yt-dlp:", latestVer, "текущая:", currentVer)
		UpdateYtDlp()
	} else {
		fmt.Println("✅ yt-dlp актуален:", currentVer)
	}
}

// downloadYtDlp скачивает yt-dlp.exe в bin/
func downloadYtDlp(ytPath string) {
	url := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("⚠ Ошибка загрузки yt-dlp:", err)
		return
	}
	defer resp.Body.Close()

	os.MkdirAll(filepath.Dir(ytPath), os.ModePerm)

	out, err := os.Create(ytPath)
	if err != nil {
		fmt.Println("⚠ Ошибка создания файла yt-dlp.exe:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("⚠ Ошибка записи yt-dlp.exe:", err)
		return
	}

	fmt.Println("✅ yt-dlp.exe скачан в bin/")
}