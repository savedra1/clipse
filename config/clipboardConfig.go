package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

/* File contains logic for parseing the cilpboard data and
general config.
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

func Init() (string, string, string, bool, error) {
	/* Ensure $HOME/.config/clipboard_manager/clipboard_history.json
	exists and create the path if not. Full path returned as string
	when successful
	*/
	currentUser, err := user.Current()
	utils.HandleError(err)

	// Construct the path to the config directory
	clipseDir := filepath.Join(currentUser.HomeDir, ".config", clipseDirName) // the ~/.config/clipboard_manager dir
	historyFilePath := filepath.Join(clipseDir, historyFileName)              // the path to the clipboard_history.json file
	tmpFileDir := filepath.Join(clipseDir, tmpDir)                            // where tmporary image files are stored
	themePath := filepath.Join(clipseDir, themeFile)                          // explicit path to theme.json file

	initTheme(themePath)

	_, err = os.Stat(historyFilePath) // File already exist?
	if os.IsNotExist(err) {

		_, err = os.Stat(clipseDir) // Config dir at least exists?
		if os.IsNotExist(err) {
			err = createDir(clipseDir)
			utils.HandleError(err)
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

	_, err = os.Stat(tmpFileDir)
	if os.IsNotExist(err) { // create temp files dir within main config
		err = createDir(tmpFileDir)
		utils.HandleError(err)
	}

	ds := DisplayServer()
	var ie bool // imagesEnabled?
	if ds == "unknown" {
		ie = false
	} else {
		ie = shell.ImagesEnabled(ds)
	}

	return historyFilePath, clipseDir, ds, ie, nil
}

func GetHistory() []ClipboardItem {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	historyFilePath, _ := Paths()
	file, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0644)
	utils.HandleError(err)

	var data ClipboardHistory

	err = json.NewDecoder(file).Decode(&data)
	utils.HandleError(err)

	// Extract clipboard history items
	return data.ClipboardHistory
}

func DeleteJsonItem(historyFilePath, item string) error {
	/* Accessed by bubbletea method on backspace keybinding:
	Deletes selected item from json file.
	*/
	//fileMutex.Lock()
	//defer fileMutex.Unlock()

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
		} else {
			if entry.FilePath != "null" {
				err = shell.DeleteImage(entry.FilePath)
				utils.HandleError(err)
			}
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

func createDir(dirPath string) error {
	/* Used to create the ~/.config/clipboard_manager dir
	in relative path.
	*/
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
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

func Paths() (string, string) {
	/* Returns full path string for clipboard file.
	useful when needing to be accessed form a
	bubbletea method.
	*/
	currentUser, err := user.Current()
	utils.HandleError(err)
	// Construct the path to the config directory
	clipseDir := filepath.Join(currentUser.HomeDir, ".config", clipseDirName)
	historyFilePath := filepath.Join(clipseDir, historyFileName)

	return historyFilePath, clipseDir
}

func ClearHistory(historyFilePath string) error {
	/* Sets clipboard_history.json file to:
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

func AddClipboardItem(configFile, text, imgPath string) error {
	var data ClipboardHistory

	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return err
	}

	item := ClipboardItem{
		Value:    text,
		Recorded: utils.GetTime(),
		FilePath: imgPath,
		Pinned:   false,
	}

	// Append the new item to the beginning of the array to appear at top of list
	data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)

	if len(data.ClipboardHistory) > maxLen {
		for i := len(data.ClipboardHistory) - 1; i >= 0; i-- { // remove the first unpinned entry starting with the oldest
			if !data.ClipboardHistory[i].Pinned {
				data.ClipboardHistory = append(data.ClipboardHistory[:i], data.ClipboardHistory[i+1:]...)
				break
			}
		}
	}

	if err = saveDataToFile(configFile, data); err != nil {
		return err
	}
	return nil
}

// This pins and unpins an item in the clipboard
func TogglePinClipboardItem(configFile string, timeStamp string) (bool, error) {
	var data ClipboardHistory
	var pinned bool // gets the pinned state of the iteem

	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return pinned, err
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return pinned, err
	}

	for i := range data.ClipboardHistory {
		if data.ClipboardHistory[i].Recorded == timeStamp {
			pinned = data.ClipboardHistory[i].Pinned
			// Toggle the pinned state
			data.ClipboardHistory[i].Pinned = !data.ClipboardHistory[i].Pinned
			break
		}
	}

	if err = saveDataToFile(configFile, data); err != nil {
		return pinned, err
	}

	return pinned, nil
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
