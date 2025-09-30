package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/hajimehoshi/oto/v2"
)

// =================== YouTube Downloader ===================

// getYTDLPBinary returns the path to yt-dlp binary based on the OS
func getYTDLPBinary() string {
	if runtime.GOOS == "windows" {
		return filepath.Join("bin", "yt-dlp.exe")
	}
	return filepath.Join("bin", "yt-dlp")
}

// GetVideoTitle extracts the video title using yt-dlp
func GetVideoTitle(url string) string {
	ytPath := getYTDLPBinary()

	// Ensure proper encoding flags
	cmd := exec.Command(ytPath, "--quiet", "--get-title", "--encoding", "utf-8", url)

	// Set Python encoding for Windows
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("⚠ Failed to get title:", err)
		return ""
	}

	title := strings.TrimSpace(string(output))

	// Debug: show what yt-dlp returned
	fmt.Printf("🔍 Title from yt-dlp: '%s'\n", title)
	fmt.Printf("📏 Length: %d bytes, UTF-8 valid: %t\n", len(title), utf8.ValidString(title))

	return title
}

// =================== File utilities ===================

// GenerateFallbackTitle creates a fallback filename
func GenerateFallbackTitle() string {
	return "video_" + fmt.Sprintf("%d", time.Now().Unix())
}

// SanitizeFileName removes unsafe characters but keeps non-Latin scripts
func SanitizeFileName(name string) string {
	if name == "" {
		return GenerateFallbackTitle()
	}

	original := name

	// Replace only characters disallowed by Windows filesystem
	dangerousChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = dangerousChars.ReplaceAllString(name, "_")

	// Remove control characters
	controlChars := regexp.MustCompile(`[\x00-\x1f\x7f]`)
	name = controlChars.ReplaceAllString(name, "")

	// Collapse multiple spaces
	multipleSpaces := regexp.MustCompile(`\s+`)
	name = multipleSpaces.ReplaceAllString(name, " ")

	// Trim spaces
	name = strings.TrimSpace(name)

	// Limit length (leave room for extension)
	if len(name) > 200 {
		name = truncateUTF8(name, 200)
	}

	// Handle Windows reserved names
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upper := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upper == reserved || strings.HasPrefix(upper, reserved+".") {
			name = "_" + name
			break
		}
	}

	// If name became empty after sanitization
	if name == "" {
		name = GenerateFallbackTitle()
	}

	// Debug: show what changed
	if original != name {
		fmt.Printf("📝 Sanitized: '%s' -> '%s'\n", original, name)
	} else {
		fmt.Printf("✅ Name kept: '%s'\n", name)
	}

	return name
}

// truncateUTF8 safely truncates a UTF-8 string
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

// ParseURLsFromFile extracts URLs from text content
func ParseURLsFromFile(content string) []string {
	lines := strings.Split(content, "\n")
	var urls []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Basic URL shape check
		if strings.Contains(line, "youtube.com") ||
			strings.Contains(line, "youtu.be") ||
			strings.Contains(line, "http://") ||
			strings.Contains(line, "https://") {
			urls = append(urls, line)
		}
	}

	return urls
}

// IsValidURL checks if string looks like a URL we accept
func IsValidURL(url string) bool {
	return strings.Contains(url, "youtube.com") ||
		strings.Contains(url, "youtu.be") ||
		strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://")
}

// =================== Audio player ===================

// intBufferToBytes converts *audio.IntBuffer to []byte (16-bit LE)
func intBufferToBytes(buf *audio.IntBuffer) []byte {
	out := make([]byte, len(buf.Data)*2)
	for i, v := range buf.Data {
		binary.LittleEndian.PutUint16(out[i*2:], uint16(int16(v)))
	}
	return out
}

// playWav plays a WAV file by path
func playWav(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)
	if !decoder.IsValidFile() {
		return fmt.Errorf("invalid WAV file: %s", path)
	}

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return fmt.Errorf("wav decode error: %v", err)
	}

	ctx, ready, err := oto.NewContext(int(buf.Format.SampleRate), buf.Format.NumChannels, 2)
	if err != nil {
		return err
	}
	<-ready

	player := ctx.NewPlayer(bytes.NewReader(intBufferToBytes(buf)))
	defer player.Close()

	player.Play()

	// Wait for playback to finish
	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// PlayBeepShort plays a short completion sound: assets/beep_short.wav
func PlayBeepShort() { _ = playTone(1000, 300) }

// PlayBeepLong plays a long completion sound: assets/beep_long.wav
func PlayBeepLong() { _ = playTone(800, 3000) }

// playTone synthesizes and plays a simple sine-wave tone using oto.
// frequencyHz: tone frequency; durationMs: duration in milliseconds.
func playTone(frequencyHz int, durationMs int) error {
	const sampleRate = 44100
	const channels = 2
	const bytesPerSample = 2 // 16-bit

	totalSamples := sampleRate * durationMs / 1000
	// stereo interleaved
	data := make([]byte, totalSamples*channels*bytesPerSample)
	amplitude := 0.25 // 25% of max to avoid clipping
	twoPiF := 2.0 * math.Pi * float64(frequencyHz)
	for i := 0; i < totalSamples; i++ {
		t := float64(i) / float64(sampleRate)
		s := math.Sin(twoPiF * t)
		v := int16(s * amplitude * 32767)
		// write to both channels
		idx := i * channels * bytesPerSample
		binary.LittleEndian.PutUint16(data[idx:], uint16(v))
		binary.LittleEndian.PutUint16(data[idx+2:], uint16(v))
	}

	ctx, ready, err := oto.NewContext(sampleRate, channels, bytesPerSample)
	if err != nil {
		return err
	}
	<-ready
	player := ctx.NewPlayer(bytes.NewReader(data))
	defer player.Close()
	player.Play()
	for player.IsPlaying() {
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

// =================== yt-dlp auto update ===================

// UpdateYtDlp updates yt-dlp binary in bin
func UpdateYtDlp() {
	ytPath := getYTDLPBinary()

	cmd := exec.Command(ytPath, "-U")
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("⚠ yt-dlp update error:", err)
		fmt.Println(string(output))
		return
	}
	fmt.Println("✅ yt-dlp updated")
	fmt.Println(string(output))
}

// CheckUpdateYtDlp checks for new yt-dlp version
func CheckUpdateYtDlp() {
	ytPath := getYTDLPBinary()

	// Check local binary exists
	if _, err := os.Stat(ytPath); os.IsNotExist(err) {
		fmt.Println("⚠ yt-dlp binary not found, downloading latest...")
		downloadYtDlp(ytPath)
		return
	}

	// Get local version
	cmd := exec.Command(ytPath, "--version")
	cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	currentVerBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("⚠ Failed to get local yt-dlp version:", err)
		return
	}
	currentVer := strings.TrimSpace(string(currentVerBytes))

	// Get latest version from GitHub
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		fmt.Println("⚠ Failed to check latest version:", err)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	latestVer, _ := data["tag_name"].(string)

	if latestVer != "" && latestVer != currentVer {
		fmt.Println("⬆ New yt-dlp available:", latestVer, "current:", currentVer)
		downloadYtDlp(ytPath) // Download new version instead of running -U
	} else {
		fmt.Println("✅ yt-dlp is up to date:", currentVer)
	}
}

// downloadYtDlp downloads the appropriate yt-dlp binary into bin/
func downloadYtDlp(ytPath string) {
    url := "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"
    if runtime.GOOS != "windows" {
        url = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux" // Используем yt-dlp_linux
    }

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("⚠ yt-dlp download error:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("⚠ Failed to download yt-dlp: HTTP status", resp.Status)
        return
    }

    os.MkdirAll(filepath.Dir(ytPath), os.ModePerm)

    out, err := os.Create(ytPath)
    if err != nil {
        fmt.Println("⚠ Failed to create yt-dlp binary:", err)
        return
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        fmt.Println("⚠ Failed to write yt-dlp binary:", err)
        return
    }

    // Set executable permissions on Linux
    if runtime.GOOS != "windows" {
        err = os.Chmod(ytPath, 0755)
        if err != nil {
            fmt.Println("⚠ Failed to set executable permissions for yt-dlp:", err)
            return
        }
    }

    fmt.Println("✅ yt-dlp binary downloaded to", ytPath)
}
