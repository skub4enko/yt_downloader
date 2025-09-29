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

// SubtitleInfo represents subtitle info
type SubtitleInfo struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	URL      string `json:"url"`
}

// AudioTrackInfo represents audio track info
type AudioTrackInfo struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	Ext      string `json:"ext"`
	Quality  string `json:"quality"`
}

// VideoMetadata contains video metadata
type VideoMetadata struct {
	Title              string           `json:"title"`
	AvailableSubtitles []SubtitleInfo   `json:"subtitles"`
	AudioTracks        []AudioTrackInfo `json:"formats"`
}

// SubtitleOptions controls subtitle download
type SubtitleOptions struct {
	DownloadSubtitles bool
	SubtitleFormat    string   // srt, vtt, ass
	Languages         []string // языки для скачивания
	DownloadAll       bool     // скачать все доступные
}

// AudioTrackOptions controls audio track preferences
type AudioTrackOptions struct {
	PreferredLanguages []string // предпочитаемые языки
	DownloadMultiple   bool     // скачать несколько дорожек
}

// Defaults
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

// GetVideoMetadata retrieves complete metadata for a video
func GetVideoMetadata(url string) (*VideoMetadata, error) {
	ytPath := filepath.Join("bin", "yt-dlp.exe")

	// Ask yt-dlp for JSON metadata
	cmd := exec.Command(ytPath,
		"--dump-json",         // выводить JSON
		"--no-warnings",       // без предупреждений
		"--encoding", "utf-8", // кодировка
		url,
	)

	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("metadata retrieval error: %v", err)
	}

	var metadata VideoMetadata
	if err := json.Unmarshal(output, &metadata); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return &metadata, nil
}

// GetAvailableSubtitles returns available subtitles for a video
func GetAvailableSubtitles(url string) ([]SubtitleInfo, error) {
	ytPath := filepath.Join("bin", "yt-dlp.exe")

	cmd := exec.Command(ytPath,
		"--list-subs", // список субтитров
		"--no-warnings",
		url,
	)

	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("subtitles list retrieval error: %v", err)
	}

	// Parse --list-subs output
	return parseSubtitlesList(string(output)), nil
}

// parseSubtitlesList parses --list-subs output
func parseSubtitlesList(output string) []SubtitleInfo {
	var subtitles []SubtitleInfo
	lines := strings.Split(output, "\n")

	var inSubtitlesSection bool
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Seek the subtitles section
		if strings.Contains(line, "Available subtitles") {
			inSubtitlesSection = true
			continue
		}

		// Skip automatic captions section
		if strings.Contains(line, "Available automatic captions") {
			inSubtitlesSection = false
			continue
		}

		if !inSubtitlesSection || line == "" {
			continue
		}

		// Parse subtitle row (format: "ru vtt")
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

// getLanguageName returns language display name
func getLanguageName(code string) string {
	languages := map[string]string{
		"ru":   "Russian",
		"en":   "English",
		"uk":   "Ukrainian",
		"de":   "Deutsch",
		"fr":   "Français",
		"es":   "Español",
		"it":   "Italiano",
		"pt":   "Português",
		"ja":   "日本語",
		"ko":   "한국어",
		"zh":   "中文",
		"ar":   "العربية",
		"hi":   "हिन्दी",
		"pl":   "Polski",
		"tr":   "Türkçe",
		"nl":   "Nederlands",
		"sv":   "Svenska",
		"no":   "Norsk",
		"da":   "Dansk",
		"fi":   "Suomi",
		"orig": "Original",
	}

	if name, exists := languages[code]; exists {
		return name
	}
	return strings.ToUpper(code) // fallback
}

// PromptSubtitleOptions prompts user for subtitle options
func PromptSubtitleOptions() SubtitleOptions {
	var options SubtitleOptions

	fmt.Println("\n📝 === SUBTITLES SETTINGS ===")

	var choice string
	fmt.Println("Download subtitles?")
	fmt.Println("1 - Yes")
	fmt.Println("2 - No")
	fmt.Print("Your choice: ")
	fmt.Scanln(&choice)

	if choice != "1" {
		options.DownloadSubtitles = false
		return options
	}

	options.DownloadSubtitles = true

	// Choose subtitle format
	fmt.Println("\nChoose subtitle format:")
	fmt.Println("1 - SRT (recommended)")
	fmt.Println("2 - VTT (WebVTT)")
	fmt.Println("3 - ASS (Advanced SubStation)")
	fmt.Print("Your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case "2":
		options.SubtitleFormat = "vtt"
	case "3":
		options.SubtitleFormat = "ass"
	default:
		options.SubtitleFormat = "srt"
	}

	// Choose languages
	fmt.Println("\nChoose subtitle languages:")
	fmt.Println("1 - Russian and English")
	fmt.Println("2 - All available")
	fmt.Println("3 - Russian only")
	fmt.Println("4 - English only")
	fmt.Println("5 - Custom list")
	fmt.Print("Your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case "2":
		options.DownloadAll = true
	case "3":
		options.Languages = []string{"ru"}
	case "4":
		options.Languages = []string{"en"}
	case "5":
		fmt.Print("Enter language codes comma-separated (e.g.: ru,en,de,fr): ")
		var langInput string
		fmt.Scanln(&langInput)
		options.Languages = strings.Split(strings.ReplaceAll(langInput, " ", ""), ",")
	default:
		options.Languages = []string{"ru", "en"}
	}

	fmt.Printf("✅ Subtitles: format %s\n", options.SubtitleFormat)
	if options.DownloadAll {
		fmt.Println("📝 Languages: ALL AVAILABLE")
	} else {
		fmt.Printf("📝 Languages: %v\n", options.Languages)
	}

	return options
}

// ShowAvailableSubtitles prints available subtitles for a video
func ShowAvailableSubtitles(url string) {
	fmt.Println("\n🔍 Checking available subtitles...")

	subtitles, err := GetAvailableSubtitles(url)
	if err != nil {
		fmt.Printf("⚠ Error: %v\n", err)
		return
	}

	if len(subtitles) == 0 {
		fmt.Println("❌ No subtitles found")
		return
	}

	fmt.Printf("✅ Subtitles found: %d\n", len(subtitles))
	fmt.Println("📝 Available languages:")

	for _, sub := range subtitles {
		fmt.Printf("   • %s (%s) - формат: %s\n",
			sub.Name, sub.Language, strings.ToUpper(sub.Ext))
	}
}

// BuildSubtitleArgs builds yt-dlp flags for subtitles
func BuildSubtitleArgs(options SubtitleOptions) []string {
	if !options.DownloadSubtitles {
		return []string{}
	}

	args := []string{
		"--write-subs", // regular subtitles
		"--sub-format", options.SubtitleFormat,
	}

	if options.DownloadAll {
		// use --all-subs to fetch all languages
		args = append(args, "--all-subs")
	} else {
		// languages list format
		langs := strings.Join(options.Languages, ",")
		args = append(args, "--sub-langs", langs)
	}

	// Debug output
	fmt.Printf("🔧 Subtitle args: %v\n", args)

	return args
}

// DownloadWithSubtitles downloads a video with subtitles
func DownloadWithSubtitles(url, filename, folder string, videoFormat string, subOptions SubtitleOptions) error {
	ytPath := filepath.Join("bin", "yt-dlp.exe")
	outPath := filepath.Join(folder, filename+".%(ext)s")

	// More precise 1080p selector
	finalVideoFormat := videoFormat
	if videoFormat == "best[height<=1080][ext=mp4]" {
		// Попробуем более точный селектор
		finalVideoFormat = "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[height<=1080][ext=mp4]/best[height<=1080]"
	} else if videoFormat == "best[height<=720][ext=mp4]" {
		finalVideoFormat = "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720][ext=mp4]/best[height<=720]"
	}

	// Base video args
	args := []string{
		"-f", finalVideoFormat,
		"-o", outPath,
		"--no-warnings",
		"--console-title", // show progress in console title
		"--ffmpeg-location", "bin",
	}

	// Add subtitle args
	subArgs := BuildSubtitleArgs(subOptions)
	args = append(args, subArgs...)

	// Stability options
	args = append(args,
		"--retries", "3",
		"--fragment-retries", "3",
	)

	args = append(args, url)

	fmt.Printf("🎬 Downloading with subtitles: %s\n", filename)
	fmt.Printf("🎯 Video format: %s\n", finalVideoFormat)
	if subOptions.DownloadSubtitles {
		if subOptions.DownloadAll {
			fmt.Printf("📝 Subtitles: %s (ALL LANGUAGES)\n", subOptions.SubtitleFormat)
		} else {
			fmt.Printf("📝 Subtitles: %s, languages: %v\n",
				subOptions.SubtitleFormat, subOptions.Languages)
		}
	}

	// Show full command for debugging
	fmt.Printf("🔧 Command: %s %s\n", ytPath, strings.Join(args, " "))

	cmd := exec.Command(ytPath, args...)
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("download error: %v", err)
	}

	fmt.Printf("✅ Done! Check output folder for subtitle files (.%s)\n",
		subOptions.SubtitleFormat)

	utils.PlayBeepShort()
	return nil
}
