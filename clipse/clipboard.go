package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
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

func getjsonData() []ClipboardItem {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	fullPath := getFullPath()
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
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

func deleteJsonItem(fullPath, item string) error {
	/* Accessed by bubbletea method on backspace keybinding:
	Deletes selected item from json file.
	*/
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var data ClipboardHistory
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var updatedClipboardHistory []ClipboardItem
	for _, entry := range data.ClipboardHistory {
		if entry.Value != item {
			updatedClipboardHistory = append(updatedClipboardHistory, entry)
		}
		if entry.Value == item && entry.FilePath != "null" {
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
	if err := os.WriteFile(fullPath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func createConfigDir(configDir string) error {
	/* Used to create the ~/.config/clipboard_manager dir
	in relative path.
	*/
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}
	return nil
}

func createHistoryFile(fullPath string) error {
	/* Used to create the clipboard_history.json file
	in relative path.
	*/
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = setBaseConfig(fullPath)
	if err != nil {
		return err
	}
	return nil
}

func getFullPath() string {
	/* Returns full path string for clipboard file.
	useful when needing to be accessed form a
	bubbletea method.
	*/
	currentUser, err := user.Current()
	handleError(err)
	// Construct the path to the config directory
	configDir := filepath.Join(currentUser.HomeDir, ".config", configDirName)
	fullPath := filepath.Join(configDir, fileName)
	return fullPath
}

func Init() (string, string, error) {
	/* Ensure $HOME/.config/clipboard_manager/clipboard_history.json
	exists and create the path if not. Full path returned as string
	when successful
	*/
	currentUser, err := user.Current()
	handleError(err)

	// Construct the path to the config directory
	configDir := filepath.Join(currentUser.HomeDir, ".config", configDirName) // the ~/.config/clipboard_manager dir
	fullPath := filepath.Join(configDir, fileName)                            // the path to the clipboard_history.json file
	filePath := filepath.Join(configDir, fileDir)                             // where tmporary image files are stored

	_, err = os.Stat(fullPath) // File already exist?
	if os.IsNotExist(err) {

		_, err = os.Stat(configDir) // Config dir at least exists?
		if os.IsNotExist(err) {
			err = createConfigDir(configDir)
			if err != nil {
				fmt.Println("Failed to create config dir:", configDir)
				os.Exit(1)
			}
		}

		err = createHistoryFile(fullPath) // Attempts creation of file now that dir path exists
		if err != nil {
			fmt.Println("Failed to create:", fullPath)
			os.Exit(1)
		}

	} else if err != nil {
		fmt.Println("Unable to check if config file exists. Please update binary permisisons.")
		os.Exit(1)
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		err = createConfigDir(filePath)
		if err != nil {
			fmt.Errorf("Failed to initialise config directory for storing image files.")
		}
	}

	return fullPath, filePath, nil
}

func setBaseConfig(fullPath string) error {
	/*
		  Sets clipboard_history.json file to:
			 {
				 "clipboardHistory": []
			 }
	*/
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644) // Permisisons specified for file to allow write
	if err != nil {
		return err
	}
	defer file.Close()

	// Truncate the file to zero length
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// Rewind the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

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
	// Read data from JSON file
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

	item := ClipboardItem{Value: text, Recorded: timeNow, FilePath: imgPath}

	// Append the new item to the beginning of the array to appear at top of list
	data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)

	err = saveDataToFile(configFile, data)
	if err != nil {
		return err
	}

	return nil
}

// saveDataToFile saves data to a JSON file
func saveDataToFile(fullPath string, data ClipboardHistory) error {
	/* Triggered from the system copy action:
	Adds the copied string to the clipboard_history.json file.
	*/
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

/*
func Init() error {
	configFile, err := checkConfig()
}*/
