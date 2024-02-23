package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

/* File contains logic for parseing the cilpboard data and
general config.
- fileName defined in constants.go
- dirName defined in constants.go
*/

// ClipboardItem struct for individual clipboardHistor array item
type ClipboardItem struct {
	// EG: {"value": "copied_string", "recorded": "datetime"}
	Value    string `json:"value"`
	Recorded string `json:"recorded"`
	FilePath string `json:"filePath"`
}

type ClipboardHistory struct {
	ClipboardHistory []ClipboardItem `json:"clipboardHistory"`
}

func displayServer() string {
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

func Init() (string, string, string, bool, error) {
	/* Ensure $HOME/.config/clipboard_manager/clipboard_history.json
	exists and create the path if not. Full path returned as string
	when successful
	*/
	currentUser, err := user.Current()
	handleError(err)

	// Construct the path to the config directory
	clipseDir := filepath.Join(currentUser.HomeDir, ".config", clipseDirName) // the ~/.config/clipboard_manager dir
	historyFilePath := filepath.Join(clipseDir, historyFileName)              // the path to the clipboard_history.json file
	themePath := filepath.Join(clipseDir, themeFile)                          // where tmporary image files are stored

	_, err = os.Stat(historyFilePath) // File already exist?
	if os.IsNotExist(err) {

		_, err = os.Stat(clipseDir) // Config dir at least exists?
		if os.IsNotExist(err) {
			err = createConfigDir(clipseDir)
			if err != nil {
				fmt.Println("Failed to create config dir:", clipseDir)
				os.Exit(1)
			}
		}

		err = createHistoryFile(historyFilePath) // Attempts creation of file now that dir path exists
		if err != nil {
			fmt.Println("Failed to create:", historyFilePath)
			os.Exit(1)
		}

	} else if err != nil {
		fmt.Println("Unable to check if config file exists. Please update binary permisisons.")
		os.Exit(1)
	}

	initTheme(themePath)

	ds := displayServer()
	var ie bool // imagesEnabled?
	if ds == "unknown" {
		ie = false
	} else {
		ie = imagesEnabled(ds)
	}

	return historyFilePath, clipseDir, ds, ie, nil
}

func getHistory() []ClipboardItem {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	historyFilePath, _ := paths()
	file, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("error opening file:", err)
		file.Close()
	}

	var data ClipboardHistory
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	// Extract clipboard history items
	return data.ClipboardHistory

}

func deleteJsonItem(historyFilePath, item string) error {
	/* Accessed by bubbletea method on backspace keybinding:
	Deletes selected item from json file.
	*/
	fileContent, err := os.ReadFile(historyFilePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var data ClipboardHistory
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var updatedClipboardHistory []ClipboardItem
	for _, entry := range data.ClipboardHistory {
		if entry.Recorded != item {
			updatedClipboardHistory = append(updatedClipboardHistory, entry)
		} else if entry.FilePath != "null" {
			err = deleteImage(entry.FilePath)
			handleError(err)
		}
	}
	updatedData := ClipboardHistory{
		ClipboardHistory: updatedClipboardHistory,
	}
	updatedJSON, err := json.Marshal(updatedData)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Write the updated JSON back to the file
	if err := os.WriteFile(historyFilePath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func createConfigDir(clipseDir string) error {
	/* Used to create the ~/.config/clipboard_manager dir
	in relative path.
	*/
	if err := os.MkdirAll(clipseDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}
	return nil
}

func createHistoryFile(historyFilePath string) error {
	/* Used to create the clipboard_history.json file
	in relative path.
	*/

	baseConfig := ClipboardHistory{
		ClipboardHistory: []ClipboardItem{},
	}

	jsonData, err := json.MarshalIndent(baseConfig, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(historyFilePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func paths() (string, string) {
	/* Returns full path string for clipboard file.
	useful when needing to be accessed form a
	bubbletea method.
	*/
	currentUser, err := user.Current()
	handleError(err)
	// Construct the path to the config directory
	clipseDir := filepath.Join(currentUser.HomeDir, ".config", clipseDirName)
	historyFilePath := filepath.Join(clipseDir, historyFileName)

	return historyFilePath, clipseDir
}

func clearHistory(historyFilePath string) error {
	/*
		  Sets clipboard_history.json file to:
			 {
				 "clipboardHistory": []
			 }
	*/
	file, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0644) // Permisisons specified for file to allow write
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	baseConfig := ClipboardHistory{
		ClipboardHistory: []ClipboardItem{},
	}

	// Encode initial history to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(baseConfig); err != nil {
		return err
	}

	return nil
}

func addClipboardItem(configFile, text, imgPath string) error {
	var data ClipboardHistory

	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return err
	}

	// If the length exceeds maxLen, remove the oldest item
	if len(data.ClipboardHistory) >= maxLen {
		data.ClipboardHistory = data.ClipboardHistory[:1]
	}

	// yyyy-mm-dd hh-mm-s.msmsms Time format
	timeNow := strings.Split(time.Now().UTC().String(), "+0000")[0]

	item := ClipboardItem{
		Value:    text,
		Recorded: timeNow,
		FilePath: imgPath,
	}

	// Append the new item to the beginning of the array to appear at top of list
	data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)

	if err = saveDataToFile(configFile, data); err != nil {
		return err
	}
	return nil
}

// saveDataToFile saves data to a JSON file
func saveDataToFile(historyFilePath string, data ClipboardHistory) error {
	/* Triggered from the system copy action:
	Adds the copied string to the clipboard_history.json file.
	*/
	file, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Truncate(0)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
