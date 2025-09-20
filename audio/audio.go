package audio

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"yt_downloader/utils"
)

// =================== Выбор битрейта ===================

var AudioBitrate string = "64" // дефолтный битрейт теперь 64 kbps

func PromptAudioQuality() {
	fmt.Println("Выберите битрейт аудио:")
	fmt.Println("1 - 64 kbps (по умолчанию)")
	fmt.Println("2 - 96 kbps")
	fmt.Println("3 - 128 kbps")
	fmt.Println("4 - 256 kbps")
	fmt.Println("5 - 320 kbps")
	fmt.Println("6 - 512 kbps")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "2":
		AudioBitrate = "96"
	case "3":
		AudioBitrate = "128"
	case "4":
		AudioBitrate = "256"
	case "5":
		AudioBitrate = "320"
	case "6":
		AudioBitrate = "512"
	default:
		AudioBitrate = "64"
	}

	fmt.Println("Выбран битрейт:", AudioBitrate, "kbps")
}

// =================== YouTube ===================

// GetTitleFromURL получает название видео с YouTube
func GetTitleFromURL(url string) string {
	title := utils.GetVideoTitle(url)
	if title == "" {
		title = utils.GenerateFallbackTitle()
	}
	return utils.SanitizeFileName(title)
}

// =================== Загрузка и извлечение аудио ===================

func DownloadAudio(url string, filename string, folder string) {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	outPath := filepath.Join(folder, filename+".%(ext)s")

	args := []string{
		"-x",                    // извлечение аудио
		"--audio-format", "mp3", // формат mp3
		"--audio-quality", AudioBitrate + "K",
		"-o", outPath,
		url,
	}

	cmd := exec.Command(ytPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("⚠ Ошибка при запуске yt-dlp:", err)
		return
	}

	fmt.Println("\n✅ Загрузка и извлечение аудио завершены:", filename+".mp3")
	utils.PlayBeep() // сигнал окончания
}

// =================== Пакетная обработка ===================

func ProcessBatchFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("⚠ Не удалось открыть файл:", filePath)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		fmt.Println("\n🎬 Загружаем:", url)
		fileName := GetTitleFromURL(url)
		folder, _ := os.Getwd()
		DownloadAudio(url, fileName, folder) // beep внутри функции
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("⚠ Ошибка чтения файла:", err)
	}
}
