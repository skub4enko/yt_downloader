package history

import (
	"encoding/json"
	"fmt"
	"os"
	//"time"
)

const historyFile = "download_history.json"

// SaveToHistory appends a record about a downloaded file into history
func SaveToHistory(url, fileName, downloadTime string) {
	history := LoadHistory()
	history = append(history, map[string]string{
		"url":           url,
		"file_name":     fileName,
		"download_time": downloadTime,
	})
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		fmt.Println("⚠ Failed to save download history:", err)
		return
	}
	err = os.WriteFile(historyFile, data, 0644)
	if err != nil {
		fmt.Println("⚠ Failed to write history to file:", err)
	}
}

// LoadHistory loads download history from file
func LoadHistory() []map[string]string {
	var history []map[string]string
	data, err := os.ReadFile(historyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("⚠ Failed to read history file:", err)
		}
		return history
	}
	err = json.Unmarshal(data, &history)
	if err != nil {
		fmt.Println("⚠ Failed to parse history file:", err)
	}
	return history
}
