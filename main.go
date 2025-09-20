package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"yt_downloader/audio"
	"yt_downloader/video"
	"yt_downloader/subtitles" // новый модуль субтитров
	"yt_downloader/utils"
)

func main() {
	fmt.Println("🎬 YouTube Downloader v2.0")
	fmt.Println("==========================")
	
	// Проверка и автообновление yt-dlp
	utils.CheckUpdateYtDlp()
	
	// Выбор типа контента
	var contentType string
	fmt.Println("\n📋 Что будем скачивать?")
	fmt.Println("1 - Аудио (MP3)")
	fmt.Println("2 - Видео (MP4/WebM)")
	fmt.Print("Ваш выбор: ")
	fmt.Scanln(&contentType)
	
	switch contentType {
	case "1":
		fmt.Println("\n🎵 === РЕЖИМ СКАЧИВАНИЯ АУДИО ===")
		handleAudioDownload()
	case "2":
		fmt.Println("\n🎬 === РЕЖИМ СКАЧИВАНИЯ ВИДЕО ===")
		handleVideoDownload()
	default:
		fmt.Println("⚠ Неверный выбор. Завершаем программу.")
		return
	}
	
	fmt.Println("\n🎉 Программа завершена!")
}

// handleAudioDownload обрабатывает скачивание аудио
func handleAudioDownload() {
	// Выбор качества аудио
	audio.PromptAudioQuality()
	
	// Выбор режима
	var mode string
	fmt.Println("\n📥 Выберите режим скачивания:")
	fmt.Println("1 - Одиночная загрузка")
	fmt.Println("2 - Пакетная загрузка из файла")
	fmt.Print("Ваш выбор: ")
	fmt.Scanln(&mode)
	
	switch mode {
	case "1":
		// Одиночная загрузка аудио
		fmt.Print("\n🔗 Введите URL видео: ")
		var url string
		fmt.Scanln(&url)
		
		if !utils.IsValidURL(url) {
			fmt.Println("⚠ Неверный формат URL")
			return
		}
		
		fmt.Println("\n🔍 Получаем информацию о видео...")
		fileName := audio.GetTitleFromURL(url)
		folder, _ := os.Getwd()
		
		fmt.Printf("📁 Файл будет сохранен: %s.mp3\n", fileName)
		audio.DownloadAudio(url, fileName, folder)
		
	case "2":
		// Пакетная загрузка аудио
		fmt.Print("\n📄 Используется файл: links.txt")
		
		// Проверяем существование файла
		if _, err := os.Stat("links.txt"); os.IsNotExist(err) {
			fmt.Println("\n⚠ Файл links.txt не найден!")
			fmt.Println("💡 Создайте файл links.txt и добавьте в него URL видео (по одному на строку)")
			return
		}
		
		audio.ProcessBatchFile("links.txt")
		
	default:
		fmt.Println("⚠ Неверный выбор режима.")
	}
}

// handleVideoDownload обрабатывает скачивание видео
func handleVideoDownload() {
	// Выбор качества видео
	video.PromptVideoQuality()
	
	// Настройки субтитров
	subOptions := subtitles.PromptSubtitleOptions()
	
	// Выбор режима
	var mode string
	fmt.Println("\n📥 Выберите режим скачивания:")
	fmt.Println("1 - Одиночная загрузка")
	fmt.Println("2 - Пакетная загрузка из файла")
	fmt.Println("3 - Проверить доступные субтитры (без скачивания)")
	fmt.Print("Ваш выбор: ")
	fmt.Scanln(&mode)
	
	switch mode {
	case "1":
		// Одиночная загрузка видео
		fmt.Print("\n🔗 Введите URL видео: ")
		var url string
		fmt.Scanln(&url)
		
		if !utils.IsValidURL(url) {
			fmt.Println("⚠ Неверный формат URL")
			return
		}
		
		// Показываем доступные субтитры если нужно
		if subOptions.DownloadSubtitles {
			subtitles.ShowAvailableSubtitles(url)
		}
		
		fmt.Println("\n🔍 Получаем информацию о видео...")
		fileName := video.GetVideoTitle(url)
		folder, _ := os.Getwd()
		
		fmt.Printf("📁 Файл будет сохранен: %s\n", fileName)
		video.DownloadVideoWithSubtitles(url, fileName, folder, subOptions)
		
	case "2":
		// Пакетная загрузка видео
		fmt.Print("\n📄 Используется файл: links.txt")
		
		// Проверяем существование файла
		if _, err := os.Stat("links.txt"); os.IsNotExist(err) {
			fmt.Println("\n⚠ Файл links.txt не найден!")
			fmt.Println("💡 Создайте файл links.txt и добавьте в него URL видео (по одному на строку)")
			return
		}
		
		// Для пакетной загрузки используем те же настройки субтитров
		processVideoBatchFileWithSubtitles("links.txt", subOptions)
		
	case "3":
		// Просто проверить субтитры
		fmt.Print("\n🔗 Введите URL видео: ")
		var url string
		fmt.Scanln(&url)
		
		if !utils.IsValidURL(url) {
			fmt.Println("⚠ Неверный формат URL")
			return
		}
		
		subtitles.ShowAvailableSubtitles(url)
		
	default:
		fmt.Println("⚠ Неверный выбор режима.")
	}
}

// processVideoBatchFileWithSubtitles обрабатывает пакетную загрузку с субтитрами
func processVideoBatchFileWithSubtitles(filePath string, subOptions subtitles.SubtitleOptions) {
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
		
		fileName := video.GetVideoTitle(url)
		video.DownloadVideoWithSubtitles(url, fileName, folder, subOptions)
		
		// Пауза между скачиваниями (опционально)
		if i < len(urls)-1 {
			fmt.Println("⏳ Пауза 2 секунды перед следующим видео...")
			// time.Sleep(2 * time.Second) // раскомментируйте если нужна пауза
		}
	}
	
	fmt.Printf("\n🎉 Пакетное скачивание завершено! Обработано видео: %d\n", len(urls))
}