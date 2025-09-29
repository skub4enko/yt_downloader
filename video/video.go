package video

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"yt_downloader/subtitles"
	"yt_downloader/utils"
)

// VideoQuality describes video quality option
type VideoQuality struct {
	Format      string
	Resolution  string
	Description string
	YtDlpFormat string
}

// Available quality options
var VideoQualities = []VideoQuality{
	{"mp4", "720p", "720p MP4 (recommended)", "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720][ext=mp4]/best[height<=720]"},
	{"mp4", "1080p", "1080p MP4 (Full HD)", "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[height<=1080][ext=mp4]/best[height<=1080]"},
	{"mp4", "1440p", "1440p MP4 (2K)", "bestvideo[height<=1440][ext=mp4]+bestaudio[ext=m4a]/best[height<=1440][ext=mp4]/best[height<=1440]"},
	{"mp4", "2160p", "2160p MP4 (4K)", "bestvideo[height<=2160][ext=mp4]+bestaudio[ext=m4a]/best[height<=2160][ext=mp4]/best[height<=2160]"},
	{"mp4", "480p", "480p MP4", "bestvideo[height<=480][ext=mp4]+bestaudio[ext=m4a]/best[height<=480][ext=mp4]/best[height<=480]"},
	{"mp4", "360p", "360p MP4", "bestvideo[height<=360][ext=mp4]+bestaudio[ext=m4a]/best[height<=360][ext=mp4]/best[height<=360]"},
	{"webm", "720p", "720p WebM", "bestvideo[height<=720][ext=webm]+bestaudio[ext=webm]/best[height<=720][ext=webm]/best[height<=720]"},
	{"webm", "1080p", "1080p WebM", "bestvideo[height<=1080][ext=webm]+bestaudio[ext=webm]/best[height<=1080][ext=webm]/best[height<=1080]"},
	{"mp4", "best", "Best MP4", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"},
	{"any", "best", "Best quality (any format)", "bestvideo+bestaudio/best"},
}

// Selected quality (default 720p MP4)
var SelectedVideoQuality = VideoQualities[0]

// PromptVideoQuality lets user choose video quality
func PromptVideoQuality() {
	fmt.Println("Select video quality:")
	for i, quality := range VideoQualities {
		fmt.Printf("%d - %s\n", i+1, quality.Description)
	}

	var choice string
	fmt.Print("Your choice (1-", len(VideoQualities), "): ")
	fmt.Scanln(&choice)

	switch choice {
	case "2":
		SelectedVideoQuality = VideoQualities[1] // 1080p MP4
	case "3":
		SelectedVideoQuality = VideoQualities[2] // 1440p MP4
	case "4":
		SelectedVideoQuality = VideoQualities[3] // 2160p MP4
	case "5":
		SelectedVideoQuality = VideoQualities[4] // 480p MP4
	case "6":
		SelectedVideoQuality = VideoQualities[5] // 360p MP4
	case "7":
		SelectedVideoQuality = VideoQualities[6] // 720p WebM
	case "8":
		SelectedVideoQuality = VideoQualities[7] // 1080p WebM
	case "9":
		SelectedVideoQuality = VideoQualities[8] // Best MP4
	case "10":
		SelectedVideoQuality = VideoQualities[9] // Best quality
	default:
		SelectedVideoQuality = VideoQualities[0] // 720p MP4 by default
	}

	fmt.Printf("✅ Selected quality: %s\n", SelectedVideoQuality.Description)
}

// GetVideoTitle grabs a safe video title from URL
func GetVideoTitle(url string) string {
	title := utils.GetVideoTitle(url)
	if title == "" {
		title = utils.GenerateFallbackTitle()
	}
	return utils.SanitizeFileName(title)
}

// DownloadVideo downloads a video with default subtitle options
func DownloadVideo(url string, filename string, folder string) {
	DownloadVideoWithOptions(url, filename, folder, subtitles.DefaultSubtitleOptions)
}

// DownloadVideoWithSubtitles downloads a video with subtitle options
func DownloadVideoWithSubtitles(url string, filename string, folder string, subOptions subtitles.SubtitleOptions) {
	DownloadVideoWithOptions(url, filename, folder, subOptions)
}

// DownloadVideoWithOptions downloads with fully specified options
func DownloadVideoWithOptions(url string, filename string, folder string, subOptions subtitles.SubtitleOptions) {
	// If subtitles requested, use subtitle pipeline
	if subOptions.DownloadSubtitles {
		err := subtitles.DownloadWithSubtitles(url, filename, folder, SelectedVideoQuality.YtDlpFormat, subOptions)
		if err != nil {
			fmt.Printf("⚠ Error: %v\n", err)
		}
		return
	}

	// Regular download without subtitles
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	outPath := filepath.Join(folder, filename+".%(ext)s")

	fmt.Printf("🎬 Downloading video: %s\n", filename)
	fmt.Printf("📁 Saving to: %s\n", folder)
	fmt.Printf("🎯 Quality: %s\n", SelectedVideoQuality.Description)

	// yt-dlp arguments
	args := []string{
		"-f", SelectedVideoQuality.YtDlpFormat, // quality format
		"-o", outPath, // output path
		"--no-warnings",   // warnings off
		"--console-title", // show process in title
		"--ffmpeg-location", "bin",
		url,
	}

	// Extra stability options
	args = append(args,
		"--retries", "3", // repeate 3 times in case of error
		"--fragment-retries", "3", // repeate fragments 3 times
	)

	cmd := exec.Command(ytPath, args...)

	// Set encoding for proper console output
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	// Forward output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("🚀 Starting download...")

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠ Video download error: %v\n", err)
		return
	}

	fmt.Printf("✅ Video downloaded successfully: %s\n", filename)
	utils.PlayBeepShort() // short completion beep
}

// ProcessVideoBatchFile processes a file with video URLs
func ProcessVideoBatchFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("⚠ can't open file: %s\n", filePath)
		return
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)

	// Read URLs from file
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" || strings.HasPrefix(url, "#") {
			continue // skip empty lines and comments
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
	folder, _ := os.Getwd()

	// Process each URL
	for i, url := range urls {
		fmt.Printf("\n🎬 Processing %d/%d: %s\n", i+1, len(urls), url)

		fileName := GetVideoTitle(url)
		DownloadVideo(url, fileName, folder)

		// Optional pause between downloads
		if i < len(urls)-1 {
			fmt.Println("⏳ Pause 2 seconds before next video...")
			// time.Sleep(2 * time.Second) // uncomment if pause needs
		}
	}

	fmt.Printf("\n🎉 Batch download completed! Processed: %d\n", len(urls))
}
