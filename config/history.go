package config

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

/* File contains logic for parsing the clipboard history.
- fileName defined in constants.go
- dirName defined in constants.go
*/

type ClipboardItem struct {
	Value    string `json:"value"`
	Recorded string `json:"recorded"`
	FilePath string `json:"filePath"`
	Pinned   bool   `json:"pinned"`
}

type ClipboardHistory struct {
	ClipboardHistory []ClipboardItem `json:"clipboardHistory"`
}

func initHistoryFile() error {
	/* Used to create the clipboard_history.json file
	in relative path.
	*/
	_, err := os.Stat(ClipseConfig.HistoryFilePath) // File already exist?
	if os.IsNotExist(err) {
		baseConfig := ClipboardHistory{
			ClipboardHistory: []ClipboardItem{},
		}

		jsonData, err := json.MarshalIndent(baseConfig, "", "    ")
		if err != nil {
			return err
		}
		err = os.WriteFile(ClipseConfig.HistoryFilePath, jsonData, 0644)

		if err != nil {
			fmt.Println("Failed to create:", ClipseConfig.HistoryFilePath)
			os.Exit(1)
		}

		// fmt.Println("Created history file:", ClipseConfig.HistoryFilePath)

	} else if err != nil {
		fmt.Println("Unable to check if history file exists. Please update binary permissions.")
		os.Exit(1)
	}

	return nil
}

func DisplayServer() string {
	/* Determine runtime and return appropriate window server.
	used to determine which dependency is required for handling
	image files.
	*/
	osName := runtime.GOOS
	switch osName {
	case "linux":
		waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
		if waylandDisplay != "" {
			return "wayland"
		} else {
			return "x11"
		}
	case "darwin":
		return "darwin"
	default:
		return "unknown"
	}
}

func GetHistory() []ClipboardItem {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	file, err := os.OpenFile(ClipseConfig.HistoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	utils.HandleError(err)

	var data ClipboardHistory

	err = json.NewDecoder(file).Decode(&data)
	utils.HandleError(err)

	// Extract clipboard history items
	return data.ClipboardHistory
}

func fileContents() ClipboardHistory {
	file, err := os.OpenFile(ClipseConfig.HistoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	utils.HandleError(err)

	var data ClipboardHistory

	err = json.NewDecoder(file).Decode(&data)
	utils.HandleError(err)

	return data
}

func WriteUpdate(data ClipboardHistory) error {
	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(ClipseConfig.HistoryFilePath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func DeleteJsonItem(item string) error {
	/* Accessed by bubbletea method on backspace keybinding:
	Deletes selected item from json file.
	*/
	data := fileContents()
	var updatedClipboardHistory []ClipboardItem

	for _, entry := range data.ClipboardHistory {
		if entry.Recorded != item {
			updatedClipboardHistory = append(updatedClipboardHistory, entry)
		} else {
			if entry.FilePath != "null" {
				err := shell.DeleteImage(entry.FilePath)
				utils.HandleError(err)
			}
		}
	}
	updatedData := ClipboardHistory{
		ClipboardHistory: updatedClipboardHistory,
	}
	err := WriteUpdate(updatedData)
	if err != nil {
		return nil
	}
	return nil
}

func createDir(dirPath string) error {
	/* Used to create the ~/.config/clipboard_manager dir
	in relative path. Takes arg to allow other dirs to be created also.
	*/
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		os.Exit(1)
	}
	return nil
}

func ClearHistory(clearType string) error {
	var data ClipboardHistory
	switch clearType {
	case "all":
		data = ClipboardHistory{
			ClipboardHistory: []ClipboardItem{},
		}
		shell.DeleteAllImages(ClipseConfig.TempDirPath)
	case "images":
		data = ClipboardHistory{
			ClipboardHistory: textItems(),
		}
		shell.DeleteAllImages(ClipseConfig.TempDirPath)
	case "text":
		data = ClipboardHistory{
			ClipboardHistory: imageItems(),
		}
	default:
		data = ClipboardHistory{
			ClipboardHistory: pinnedItems(),
		}
	}
	err := WriteUpdate(data)
	if err != nil {
		return err
	}
	return nil
}

func pinnedItems() []ClipboardItem {
	pinnedItems := []ClipboardItem{}
	history := GetHistory()
	for _, item := range history {
		if item.Pinned {
			pinnedItems = append(pinnedItems, item)
		}
	}
	return pinnedItems
}

func imageItems() []ClipboardItem {
	images := []ClipboardItem{}
	history := GetHistory()
	for _, item := range history {
		if item.FilePath != "null" {
			images = append(images, item)
		}
	}
	return images
}

func textItems() []ClipboardItem {
	textItems := []ClipboardItem{}
	history := GetHistory()
	for _, item := range history {
		if item.FilePath == "null" {
			textItems = append(textItems, item)
		}
	}
	return textItems
}
func AddClipboardItem(text, fp string) error {
	data := fileContents()
	item := ClipboardItem{
		Value:    text,
		Recorded: utils.GetTime(),
		FilePath: fp,
		Pinned:   false,
	}

	// Append the new item to the beginning of the array to appear at top of list
	data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)

	if len(data.ClipboardHistory) > ClipseConfig.MaxHistory {
		for i := len(data.ClipboardHistory) - 1; i >= 0; i-- {
			// remove the first unpinned entry starting with the oldest
			if !data.ClipboardHistory[i].Pinned {
				data.ClipboardHistory = append(data.ClipboardHistory[:i], data.ClipboardHistory[i+1:]...)
				break
			}
		}
	}

	if err := WriteUpdate(data); err != nil {
		return err
	}
	return nil
}

// This pins and unpins an item in the clipboard
func TogglePinClipboardItem(timeStamp string) (bool, error) {
	data := fileContents()
	var pinned bool

	for i, item := range data.ClipboardHistory {
		if item.Recorded == timeStamp {
			pinned = item.Pinned
			// Toggle the pinned state
			data.ClipboardHistory[i].Pinned = !item.Pinned
			break
		}
	}

	if err := WriteUpdate(data); err != nil {
		return pinned, err
	}
	return pinned, nil
}

// Contains checks if a string exists in the most recent 3 items
func Contains(str string) bool {
	data := GetHistory()
	if len(data) > 3 {
		data = data[:3]
	}
	for _, item := range data {
		if item.Value == str {
			return true
		}
	}
	return false
}
