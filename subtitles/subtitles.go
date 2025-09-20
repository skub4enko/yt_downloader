package subtitles

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"yt_downloader/utils"
)

// SubtitleInfo представляет информацию о субтитрах
type SubtitleInfo struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	URL      string `json:"url"`
}

// AudioTrackInfo представляет информацию о звуковой дорожке
type AudioTrackInfo struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	Quality  string `json:"quality"`
}

// VideoMetadata содержит всю информацию о видео
type VideoMetadata struct {
	Title             string            `json:"title"`
	AvailableSubtitles []SubtitleInfo   `json:"subtitles"`
	AudioTracks       []AudioTrackInfo `json:"formats"`
}

// SubtitleOptions настройки для скачивания субтитров
type SubtitleOptions struct {
	DownloadSubtitles bool
	SubtitleFormat    string   // srt, vtt, ass
	Languages         []string // языки для скачивания
	DownloadAll       bool     // скачать все доступные
}

// AudioTrackOptions настройки для звуковых дорожек
type AudioTrackOptions struct {
	PreferredLanguages []string // предпочитаемые языки
	DownloadMultiple   bool     // скачать несколько дорожек
}

// По умолчанию
var (
	DefaultSubtitleOptions = SubtitleOptions{
		DownloadSubtitles: false,
		SubtitleFormat:    "srt",
		Languages:         []string{"ru", "en", "uk"},
		DownloadAll:       false,
	}
	
	DefaultAudioOptions = AudioTrackOptions{
		PreferredLanguages: []string{"ru", "en", "uk", "orig"},
		DownloadMultiple:   false,
	}
)

// GetVideoMetadata получает полную информацию о видео
func GetVideoMetadata(url string) (*VideoMetadata, error) {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	
	// Получаем JSON с полной информацией о видео
	cmd := exec.Command(ytPath, 
		"--dump-json",           // выводить JSON
		"--no-warnings",         // без предупреждений
		"--encoding", "utf-8",   // кодировка
		url,
	)
	
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения метаданных: %v", err)
	}
	
	var metadata VideoMetadata
	if err := json.Unmarshal(output, &metadata); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}
	
	return &metadata, nil
}

// GetAvailableSubtitles получает список доступных субтитров
func GetAvailableSubtitles(url string) ([]SubtitleInfo, error) {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	
	cmd := exec.Command(ytPath,
		"--list-subs",           // список субтитров
		"--no-warnings",
		url,
	)
	
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка субтитров: %v", err)
	}
	
	// Парсим вывод --list-subs
	return parseSubtitlesList(string(output)), nil
}

// parseSubtitlesList парсит вывод команды --list-subs
func parseSubtitlesList(output string) []SubtitleInfo {
	var subtitles []SubtitleInfo
	lines := strings.Split(output, "\n")
	
	var inSubtitlesSection bool
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Ищем секцию с субтитрами
		if strings.Contains(line, "Available subtitles") {
			inSubtitlesSection = true
			continue
		}
		
		// Пропускаем автоматические субтитры
		if strings.Contains(line, "Available automatic captions") {
			inSubtitlesSection = false
			continue
		}
		
		if !inSubtitlesSection || line == "" {
			continue
		}
		
		// Парсим строку с субтитрами (формат: "ru vtt")
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			subtitle := SubtitleInfo{
				Language: parts[0],
				Name:     getLanguageName(parts[0]),
				Ext:      parts[1],
			}
			subtitles = append(subtitles, subtitle)
		}
	}
	
	return subtitles
}

// getLanguageName возвращает полное название языка
func getLanguageName(code string) string {
	languages := map[string]string{
		"ru":    "Русский",
		"en":    "English",
		"uk":    "Українська",
		"de":    "Deutsch",
		"fr":    "Français",
		"es":    "Español",
		"it":    "Italiano",
		"pt":    "Português",
		"ja":    "日本語",
		"ko":    "한국어",
		"zh":    "中文",
		"ar":    "العربية",
		"hi":    "हिन्दी",
		"pl":    "Polski",
		"tr":    "Türkçe",
		"nl":    "Nederlands",
		"sv":    "Svenska",
		"no":    "Norsk",
		"da":    "Dansk",
		"fi":    "Suomi",
		"orig":  "Оригинал",
	}
	
	if name, exists := languages[code]; exists {
		return name
	}
	return strings.ToUpper(code) // fallback
}

// PromptSubtitleOptions позволяет настроить параметры субтитров
func PromptSubtitleOptions() SubtitleOptions {
	var options SubtitleOptions
	
	fmt.Println("\n📝 === НАСТРОЙКА СУБТИТРОВ ===")
	
	var choice string
	fmt.Println("Скачивать субтитры?")
	fmt.Println("1 - Да")
	fmt.Println("2 - Нет")
	fmt.Print("Ваш выбор: ")
	fmt.Scanln(&choice)
	
	if choice != "1" {
		options.DownloadSubtitles = false
		return options
	}
	
	options.DownloadSubtitles = true
	
	// Выбор формата субтитров
	fmt.Println("\nВыберите формат субтитров:")
	fmt.Println("1 - SRT (рекомендуется)")
	fmt.Println("2 - VTT (WebVTT)")
	fmt.Println("3 - ASS (Advanced SubStation)")
	fmt.Print("Ваш выбор: ")
	fmt.Scanln(&choice)
	
	switch choice {
	case "2":
		options.SubtitleFormat = "vtt"
	case "3":
		options.SubtitleFormat = "ass"
	default:
		options.SubtitleFormat = "srt"
	}
	
	// Выбор языков
	fmt.Println("\nВыберите языки субтитров:")
	fmt.Println("1 - Русский и английский")
	fmt.Println("2 - Все доступные")
	fmt.Println("3 - Только русский")
	fmt.Println("4 - Только английский")
	fmt.Println("5 - Пользовательский выбор")
	fmt.Print("Ваш выбор: ")
	fmt.Scanln(&choice)
	
	switch choice {
	case "2":
		options.DownloadAll = true
	case "3":
		options.Languages = []string{"ru"}
	case "4":
		options.Languages = []string{"en"}
	case "5":
		fmt.Print("Введите коды языков через запятую (например: ru,en,de,fr): ")
		var langInput string
		fmt.Scanln(&langInput)
		options.Languages = strings.Split(strings.ReplaceAll(langInput, " ", ""), ",")
	default:
		options.Languages = []string{"ru", "en"}
	}
	
	fmt.Printf("✅ Настройки субтитров: формат %s\n", options.SubtitleFormat)
	if options.DownloadAll {
		fmt.Println("📝 Языки: ВСЕ ДОСТУПНЫЕ")
	} else {
		fmt.Printf("📝 Языки: %v\n", options.Languages)
	}
	
	return options
}

// ShowAvailableSubtitles показывает доступные субтитры для видео
func ShowAvailableSubtitles(url string) {
	fmt.Println("\n🔍 Проверяем доступные субтитры...")
	
	subtitles, err := GetAvailableSubtitles(url)
	if err != nil {
		fmt.Printf("⚠ Ошибка: %v\n", err)
		return
	}
	
	if len(subtitles) == 0 {
		fmt.Println("❌ Субтитры не найдены")
		return
	}
	
	fmt.Printf("✅ Найдено субтитров: %d\n", len(subtitles))
	fmt.Println("📝 Доступные языки:")
	
	for _, sub := range subtitles {
		fmt.Printf("   • %s (%s) - формат: %s\n", 
			sub.Name, sub.Language, strings.ToUpper(sub.Ext))
	}
}

// BuildSubtitleArgs создает аргументы для yt-dlp с субтитрами
// ИСПРАВЛЕНИЕ: правильные флаги для субтитров
func BuildSubtitleArgs(options SubtitleOptions) []string {
	if !options.DownloadSubtitles {
		return []string{}
	}
	
	args := []string{
		"--write-subs",        // скачивать обычные субтитры
		"--sub-format", options.SubtitleFormat,
	}
	
	if options.DownloadAll {
		// ИСПРАВЛЕНИЕ: для всех субтитров используем --all-subs
		args = append(args, "--all-subs")
	} else {
		// ИСПРАВЛЕНИЕ: правильный формат языков
		langs := strings.Join(options.Languages, ",")
		args = append(args, "--sub-langs", langs)
	}
	
	// Добавляем отладочную информацию
	fmt.Printf("🔧 Аргументы субтитров: %v\n", args)
	
	return args
}

// DownloadWithSubtitles скачивает видео с субтитрами
func DownloadWithSubtitles(url, filename, folder string, videoFormat string, subOptions SubtitleOptions) error {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	outPath := filepath.Join(folder, filename+".%(ext)s")
	
	// ИСПРАВЛЕНИЕ: более точный селектор формата для 1080p
	finalVideoFormat := videoFormat
	if videoFormat == "best[height<=1080][ext=mp4]" {
		// Попробуем более точный селектор
		finalVideoFormat = "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[height<=1080][ext=mp4]/best[height<=1080]"
	} else if videoFormat == "best[height<=720][ext=mp4]" {
		finalVideoFormat = "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720][ext=mp4]/best[height<=720]"
	}
	
	// Базовые аргументы для видео
	args := []string{
		"-f", finalVideoFormat,
		"-o", outPath,
		"--no-warnings",
		"--console-title",        // показать прогресс в заголовке
	}
	
	// Добавляем аргументы для субтитров
	subArgs := BuildSubtitleArgs(subOptions)
	args = append(args, subArgs...)
	
	// Дополнительные опции для стабильности
	args = append(args, 
		"--retries", "3",
		"--fragment-retries", "3",
	)
	
	args = append(args, url)
	
	fmt.Printf("🎬 Скачиваем видео с субтитрами: %s\n", filename)
	fmt.Printf("🎯 Формат видео: %s\n", finalVideoFormat)
	if subOptions.DownloadSubtitles {
		if subOptions.DownloadAll {
			fmt.Printf("📝 Субтитры: %s (ВСЕ ЯЗЫКИ)\n", subOptions.SubtitleFormat)
		} else {
			fmt.Printf("📝 Субтитры: %s, языки: %v\n", 
				subOptions.SubtitleFormat, subOptions.Languages)
		}
	}
	
	// Показываем полную команду для отладки
	fmt.Printf("🔧 Команда: %s %s\n", ytPath, strings.Join(args, " "))
	
	cmd := exec.Command(ytPath, args...)
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ошибка скачивания: %v", err)
	}
	
	fmt.Printf("✅ Готово! Проверьте папку на наличие файлов субтитров (.%s)\n", 
		subOptions.SubtitleFormat)
	
	utils.PlayBeep()
	return nil
}