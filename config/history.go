package config

import (
	"encoding/json"
	"fmt"
	"os"

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
		if err = os.WriteFile(ClipseConfig.HistoryFilePath, jsonData, 0644); err != nil {
			utils.LogERROR(fmt.Sprintf("Failed to create %s", ClipseConfig.HistoryFilePath))
			return err
		}
		return nil
	}

	if err != nil {
		utils.LogERROR("Unable to check if history file exists. Please update binary permissions.")
		return err
	}
	return nil
}

func GetHistory() []ClipboardItem {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	file, err := os.OpenFile(ClipseConfig.HistoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	utils.HandleError(err)

	var data ClipboardHistory

	utils.HandleError(json.NewDecoder(file).Decode(&data))

	return data.ClipboardHistory
}

func fileContents() ClipboardHistory {
	file, err := os.OpenFile(ClipseConfig.HistoryFilePath, os.O_RDWR|os.O_CREATE, 0644)
	utils.HandleError(err)

	var data ClipboardHistory

	utils.HandleError(json.NewDecoder(file).Decode(&data))

	return data
}

func WriteUpdate(data ClipboardHistory) error {
	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(ClipseConfig.HistoryFilePath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}

	return nil
}

func DeleteItems(timeStamps []string) error {
	data := fileContents()
	updatedData := []ClipboardItem{}

	toDelete := make(map[string]bool)
	for _, ts := range timeStamps {
		toDelete[ts] = true
	}
	for _, item := range data.ClipboardHistory {
		if toDelete[item.Recorded] {
			if item.FilePath == "null" {
				continue
			}
			if err := shell.DeleteImage(item.FilePath); err != nil {
				utils.LogERROR(fmt.Sprintf("failed to delete image file | %s", item.FilePath))
			}
			continue
		}
		updatedData = append(updatedData, item)
	}
	updatedFile := ClipboardHistory{
		ClipboardHistory: updatedData,
	}
	return WriteUpdate(updatedFile)

}

func ClearHistory(clearType string) error {
	var data ClipboardHistory
	switch clearType {
	case "all":
		data = ClipboardHistory{
			ClipboardHistory: []ClipboardItem{},
		}
		if err := shell.DeleteAllImages(ClipseConfig.TempDirPath); err != nil {
			utils.LogERROR(fmt.Sprintf("could not delete all images: %s", err))
		}
	case "images":
		data = ClipboardHistory{
			ClipboardHistory: TextItems(),
		}
		if err := shell.DeleteAllImages(ClipseConfig.TempDirPath); err != nil {
			utils.LogERROR(fmt.Sprintf("could not read file dir: %s", err))
		}
	case "text":
		data = ClipboardHistory{
			ClipboardHistory: imageItems(),
		}
	default:
		data = ClipboardHistory{
			ClipboardHistory: pinnedItems(),
		}
	}
	return WriteUpdate(data)

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

func TextItems() []ClipboardItem {
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

	if !ClipseConfig.AllowDuplicates {
		duplicates, isPinned := duplicateItems(data.ClipboardHistory, item)
		data.ClipboardHistory = removeDuplicates(data.ClipboardHistory, duplicates)
		item.Pinned = isPinned
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
	return WriteUpdate(data)
}

func duplicateItems(currentHistory []ClipboardItem, newItem ClipboardItem) ([]string, bool) {
	isPinned := false
	timestamps := []string{}

	for _, item := range currentHistory {
		if isItemDuplicate(item, newItem) {
			timestamps = append(timestamps, item.Recorded)
			if item.Pinned {
				isPinned = true
			}
		}
	}

	return timestamps, isPinned
}

func isItemDuplicate(item, newItem ClipboardItem) bool {
	if item.FilePath == "null" && newItem.FilePath == "null" {
		return item.Value == newItem.Value
	}
	if item.FilePath != "null" && newItem.FilePath != "null" {
		return utils.GetImgIdentifier(item.Value) == utils.GetImgIdentifier(newItem.Value)
	}
	return false
}

func removeDuplicates(clipboardHistory []ClipboardItem, duplicates []string) []ClipboardItem {
	toDelete := make(map[string]bool)
	for _, ts := range duplicates {
		toDelete[ts] = true
	}
	updatedHistory := []ClipboardItem{}
	for _, item := range clipboardHistory {
		if !toDelete[item.Recorded] {
			updatedHistory = append(updatedHistory, item)
			continue
		}

		if item.FilePath != "null" {
			if err := shell.DeleteImage(item.FilePath); err != nil {
				utils.LogERROR(fmt.Sprintf("failed to delete image | %s | %s", item.FilePath, err))
			}
		}
	}
	return updatedHistory
}

// This pins and unpins an item in the clipboard
func TogglePinClipboardItem(timeStamp string) (bool, error) {
	data := fileContents()
	var pinned bool

	for i, item := range data.ClipboardHistory {
		if item.Recorded == timeStamp {
			// Toggle the pinned state
			data.ClipboardHistory[i].Pinned = !item.Pinned
			pinned = item.Pinned
			break
		}
	}

	if err := WriteUpdate(data); err != nil {
		return pinned, err
	}
	return pinned, nil
}
