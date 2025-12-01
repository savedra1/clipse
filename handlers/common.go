// handlers/common.go
package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

func SaveImageCommon(imgData []byte) error {
	byteLength := strconv.Itoa(len(string(imgData)))
	fileName := fmt.Sprintf("%s-%s.png", byteLength, utils.GetTimeStamp())
	itemTitle := fmt.Sprintf("%s %s", imgIcon, fileName)
	filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

	if err := os.WriteFile(filePath, imgData, 0644); err != nil {
		return err
	}

	if err := config.AddClipboardItem(itemTitle, filePath); err != nil {
		return err
	}
	return nil
}

func SaveTextCommon(textData string) error {
	if err := config.AddClipboardItem(textData, "null"); err != nil {
		return err
	}
	return nil
}

func isAppExcluded(appName string, excludedList []string) bool {
	if appName == "" {
		return false
	}

	appNameLower := strings.ToLower(appName)

	for _, excluded := range excludedList {
		excludedLower := strings.ToLower(excluded)

		if excludedLower != "" && strings.Contains(appNameLower, excludedLower) {
			return true
		}
	}

	return false
}
