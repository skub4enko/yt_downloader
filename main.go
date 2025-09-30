package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"yt_downloader/audio"
	"yt_downloader/subtitles"
	"yt_downloader/utils"
	"yt_downloader/video"
)

func main() {
	fmt.Println("🎬 YouTube Downloader v2.0")
	fmt.Println("==========================")

	// Check and auto-update yt-dlp
	utils.CheckUpdateYtDlp()

	// Ensure links.txt exists in the project root
	linksFile := filepath.Join(".", "links.txt") // Корень проекта
	if _, err := os.Stat(linksFile); os.IsNotExist(err) {
		err := os.WriteFile(linksFile, []byte("# Enter YouTube URLs here, one per line\n# Example: https://www.youtube.com/watch?v=example\n"), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠ Failed to create links.txt: %v\n", err)
			return
		}
	}

	// Choose content type
	var contentType string
	fmt.Println("\n📋 What do you want to download?")
	fmt.Println("1 - Audio (MP3)")
	fmt.Println("2 - Video (MP4/WebM)")
	fmt.Print("Your choice: ")
	fmt.Scanln(&contentType)

	switch contentType {
	case "1":
		fmt.Println("\n🎵 === AUDIO DOWNLOAD MODE ===")
		handleAudioDownload()
	case "2":
		fmt.Println("\n🎬 === VIDEO DOWNLOAD MODE ===")
		handleVideoDownload()
	default:
		fmt.Println("⚠ Invalid choice. Exiting.")
		return
	}

	fmt.Println("\n🎉 Done!")
	fmt.Println("Press Enter to close")
	fmt.Scanln()
}

// handleAudioDownload handles audio download flow
func handleAudioDownload() {
	audio.PromptAudioQuality()

	var mode string
	fmt.Println("\n📥 Select download mode:")
	fmt.Println("1 - Single URL")
	fmt.Println("2 - Batch from file")
	fmt.Print("Your choice: ")
	fmt.Scanln(&mode)

	switch mode {
	case "1":
		fmt.Print("\n🔗 Enter video URL: ")
		var url string
		fmt.Scanln(&url)

		if !utils.IsValidURL(url) {
			fmt.Println("⚠ Invalid URL format")
			return
		}

		folder := chooseDownloadFolder()
		fmt.Println("\n🔍 Fetching video info...")
		fileName := audio.GetTitleFromURL(url)
		fmt.Printf("📁 Output file: %s.mp3\n", fileName)
		audio.DownloadAudio(url, fileName, folder)

	case "2":
		folder := chooseDownloadFolder()
		audio.ProcessBatchFile("links.txt", folder)

	default:
		fmt.Println("⚠ Invalid mode selection.")
	}
}

// handleVideoDownload handles video download flow
func handleVideoDownload() {
	video.PromptVideoQuality()
	subOptions := subtitles.PromptSubtitleOptions()

	var mode string
	fmt.Println("\n📥 Select download mode:")
	fmt.Println("1 - Single URL")
	fmt.Println("2 - Batch from file")
	fmt.Println("3 - List available subtitles (no download)")
	fmt.Print("Your choice: ")
	fmt.Scanln(&mode)

	switch mode {
	case "1":
		fmt.Print("\n🔗 Enter video URL: ")
		var url string
		fmt.Scanln(&url)

		if !utils.IsValidURL(url) {
			fmt.Println("⚠ Invalid URL format")
			return
		}

		if subOptions.DownloadSubtitles {
			subtitles.ShowAvailableSubtitles(url)
		}

		folder := chooseDownloadFolder()
		fmt.Println("\n🔍 Fetching video info...")
		fileName := video.GetVideoTitle(url)
		fmt.Printf("📁 Output file: %s\n", fileName)
		video.DownloadVideoWithSubtitles(url, fileName, folder, subOptions)

	case "2":
		folder := chooseDownloadFolder()
		processVideoBatchFileWithSubtitles("links.txt", folder, subOptions)

	case "3":
		fmt.Print("\n🔗 Enter video URL: ")
		var url string
		fmt.Scanln(&url)

		if !utils.IsValidURL(url) {
			fmt.Println("⚠ Invalid URL format")
			return
		}

		subtitles.ShowAvailableSubtitles(url)

	default:
		fmt.Println("⚠ Invalid mode selection.")
	}
}

// chooseDownloadFolder asks where to save files
func chooseDownloadFolder() string {
	var choice string
	fmt.Println("\n📂 Where to save files?")
	fmt.Println("1 - Current program folder (default)")
	fmt.Println("2 - Enter a custom folder path")
	fmt.Print("Your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case "2":
		fmt.Print("Enter full folder path: ")
		var path string
		fmt.Scanln(&path)
		return path
	default:
		folder, _ := os.Getwd()
		return folder
	}
}

// processVideoBatchFileWithSubtitles handles batch video downloads
func processVideoBatchFileWithSubtitles(filePath string, folder string, subOptions subtitles.SubtitleOptions) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("⚠ Failed to open file: %s\n", filePath)
		return
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" || strings.HasPrefix(url, "#") {
			continue
		}
		urls = append(urls, url)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("⚠ File read error: %v\n", err)
		return
	}

	if len(urls) == 0 {
		fmt.Println("⚠ No valid URLs found in file")
		return
	}

	fmt.Printf("📋 Found %d videos to download\n", len(urls))

	for i, url := range urls {
		fmt.Printf("\n🎬 Processing %d/%d: %s\n", i+1, len(urls), url)
		fileName := video.GetVideoTitle(url)
		video.DownloadVideoWithSubtitles(url, fileName, folder, subOptions)
		if i < len(urls)-1 {
			fmt.Println("⏳ Pause 2 seconds before next video...")
			time.Sleep(2 * time.Second)
		}
	}

	fmt.Printf("\n🎉 Batch download completed! Processed: %d\n", len(urls))
	utils.PlayBeepLong()
}
