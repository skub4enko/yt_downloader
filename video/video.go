package video

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"yt_downloader/utils"
	"yt_downloader/subtitles"
)

// VideoQuality представляет качество видео
type VideoQuality struct {
	Format      string
	Resolution  string
	Description string
	YtDlpFormat string
}

// Доступные варианты качества видео
var VideoQualities = []VideoQuality{
	{"mp4", "720p", "720p MP4 (рекомендуется)", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720][ext=mp4]/best[height<=720]"},
	{"mp4", "1080p", "1080p MP4 (Full HD)", "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[height<=1080][ext=mp4]/best[height<=1080]"},
	{"mp4", "1440p", "1440p MP4 (2K)", "bestvideo[height<=1440][ext=mp4]+bestaudio[ext=m4a]/best[height<=1440][ext=mp4]/best[height<=1440]"},
	{"mp4", "2160p", "2160p MP4 (4K)", "bestvideo[height<=2160][ext=mp4]+bestaudio[ext=m4a]/best[height<=2160][ext=mp4]/best[height<=2160]"},
	{"webm", "720p", "720p WebM", "bestvideo[height<=720][ext=webm]+bestaudio[ext=webm]/best[height<=720][ext=webm]/best[height<=720]"},
	{"webm", "1080p", "1080p WebM", "bestvideo[height<=1080][ext=webm]+bestaudio[ext=webm]/best[height<=1080][ext=webm]/best[height<=1080]"},
	{"mp4", "best", "Лучшее MP4", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"},
	{"any", "best", "Лучшее качество (любой формат)", "bestvideo+bestaudio/best"},
}

// Выбранное качество видео (по умолчанию 720p MP4)
var SelectedVideoQuality = VideoQualities[0]

// PromptVideoQuality позволяет выбрать качество видео
func PromptVideoQuality() {
	fmt.Println("Выберите качество видео:")
	for i, quality := range VideoQualities {
		fmt.Printf("%d - %s\n", i+1, quality.Description)
	}
	
	var choice string
	fmt.Print("Ваш выбор (1-", len(VideoQualities), "): ")
	fmt.Scanln(&choice)
	
	switch choice {
	case "2":
		SelectedVideoQuality = VideoQualities[1] // 1080p MP4
	case "3":
		SelectedVideoQuality = VideoQualities[2] // 1440p MP4
	case "4":
		SelectedVideoQuality = VideoQualities[3] // 2160p MP4
	case "5":
		SelectedVideoQuality = VideoQualities[4] // 720p WebM
	case "6":
		SelectedVideoQuality = VideoQualities[5] // 1080p WebM
	case "7":
		SelectedVideoQuality = VideoQualities[6] // Лучшее MP4
	case "8":
		SelectedVideoQuality = VideoQualities[7] // Лучшее качество
	default:
		SelectedVideoQuality = VideoQualities[0] // 720p MP4 по умолчанию
	}
	
	fmt.Printf("✅ Выбрано качество: %s\n", SelectedVideoQuality.Description)
}

// GetVideoTitle получает название видео с YouTube
func GetVideoTitle(url string) string {
	title := utils.GetVideoTitle(url)
	if title == "" {
		title = utils.GenerateFallbackTitle()
	}
	return utils.SanitizeFileName(title)
}

// DownloadVideo скачивает видео в выбранном качестве
func DownloadVideo(url string, filename string, folder string) {
	DownloadVideoWithOptions(url, filename, folder, subtitles.DefaultSubtitleOptions)
}

// DownloadVideoWithSubtitles скачивает видео с настройками субтитров
func DownloadVideoWithSubtitles(url string, filename string, folder string, subOptions subtitles.SubtitleOptions) {
	DownloadVideoWithOptions(url, filename, folder, subOptions)
}

// DownloadVideoWithOptions скачивает видео с полными настройками
func DownloadVideoWithOptions(url string, filename string, folder string, subOptions subtitles.SubtitleOptions) {
	// Если нужны субтитры, используем специальную функцию
	if subOptions.DownloadSubtitles {
		err := subtitles.DownloadWithSubtitles(url, filename, folder, SelectedVideoQuality.YtDlpFormat, subOptions)
		if err != nil {
			fmt.Printf("⚠ Ошибка: %v\n", err)
		}
		return
	}
	
	// Обычное скачивание видео без субтитров
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	outPath := filepath.Join(folder, filename+".%(ext)s")
	
	fmt.Printf("🎬 Скачиваем видео: %s\n", filename)
	fmt.Printf("📁 Сохраняем в: %s\n", folder)
	fmt.Printf("🎯 Качество: %s\n", SelectedVideoQuality.Description)
	
	// Аргументы для yt-dlp
	args := []string{
		"-f", SelectedVideoQuality.YtDlpFormat, // формат качества
		"-o", outPath,                          // путь вывода
		"--no-warnings",                        // отключить предупреждения
		"--console-title",                      // показать прогресс в заголовке
		url,
	}
	
	// Дополнительные опции для стабильности
	args = append(args, 
		"--retries", "3",           // повторить 3 раза при ошибке
		"--fragment-retries", "3",  // повторить фрагменты 3 раза
	)
	
	cmd := exec.Command(ytPath, args...)
	
	// Устанавливаем кодировку для корректного отображения
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	
	// Перенаправляем вывод
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	fmt.Println("🚀 Начинаем скачивание...")
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠ Ошибка при скачивании видео: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Видео скачано успешно: %s\n", filename)
	utils.PlayBeep() // звуковой сигнал окончания
}

// ProcessVideoBatchFile обрабатывает файл со списком URL для видео
func ProcessVideoBatchFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("⚠ Не удалось открыть файл: %s\n", filePath)
		return
	}
	defer file.Close()
	
	var urls []string
	scanner := bufio.NewScanner(file)
	
	// Читаем все URL из файла
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" || strings.HasPrefix(url, "#") {
			continue // пропускаем пустые строки и комментарии
		}
		urls = append(urls, url)
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("⚠ Ошибка чтения файла: %v\n", err)
		return
	}
	
	if len(urls) == 0 {
		fmt.Println("⚠ В файле не найдено валидных URL")
		return
	}
	
	fmt.Printf("📋 Найдено %d видео для скачивания\n", len(urls))
	folder, _ := os.Getwd()
	
	// Обрабатываем каждый URL
	for i, url := range urls {
		fmt.Printf("\n🎬 Обрабатываем %d/%d: %s\n", i+1, len(urls), url)
		
		fileName := GetVideoTitle(url)
		DownloadVideo(url, fileName, folder)
		
		// Пауза между скачиваниями (опционально)
		if i < len(urls)-1 {
			fmt.Println("⏳ Пауза 2 секунды перед следующим видео...")
			// time.Sleep(2 * time.Second) // раскомментируйте если нужна пауза
		}
	}
	
	fmt.Printf("\n🎉 Пакетное скачивание завершено! Обработано видео: %d\n", len(urls))
}