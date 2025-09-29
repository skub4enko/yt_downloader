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

// =================== Audio bitrate selection ===================

var AudioBitrate string = "64" // default bitrate

func PromptAudioQuality() {
	fmt.Println("Select audio bitrate:")
	fmt.Println("0 - 32 kbps")
	fmt.Println("1 - 64 kbps (default)")
	fmt.Println("2 - 96 kbps")
	fmt.Println("3 - 128 kbps")
	fmt.Println("4 - 256 kbps")
	fmt.Println("5 - 320 kbps")
	fmt.Println("6 - 512 kbps")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "0":
		AudioBitrate = "32"
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

	fmt.Println("Selected bitrate:", AudioBitrate, "kbps")
}

// =================== YouTube ===================

func GetTitleFromURL(url string) string {
	title := utils.GetVideoTitle(url)
	if title == "" {
		title = utils.GenerateFallbackTitle()
	}
	return utils.SanitizeFileName(title)
}

// =================== Audio download ===================

func DownloadAudio(url, filename, folder string) {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	outPath := filepath.Join(folder, filename+".%(ext)s")

	args := []string{
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", AudioBitrate + "K",
		"--ffmpeg-location", "bin",
		"-o", outPath,
		url,
	}

	cmd := exec.Command(ytPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("⚠ Failed to run yt-dlp:", err)
		return
	}

	fmt.Println("\n✅ Audio download and extraction completed:", filename+".mp3")

	// secure call beep
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("⚠ Beep playback error:", r)
		}
	}()
	utils.PlayBeepShort()
}

// =================== Batch download ===================

func ProcessBatchFile(filePath, folder string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("⚠ Failed to open file:", filePath, "error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		fmt.Println("\n🎬 Downloading:", url)

		fileName := GetTitleFromURL(url)
		if fileName == "" {
			fmt.Println("⚠ Failed to get filename for:", url)
			continue
		}

		//secure start every download
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("⚠ Panic during download:", r)
				}
			}()
			DownloadAudio(url, fileName, folder)
		}()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("⚠ File read error:", err)
	}
}
