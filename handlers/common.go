// handlers/common.go
package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

func CopyText(text, ds string) {
	switch ds {
	case "darwin":
		DarwinCopyText(text)
	case "wayland":
		WaylandCopy(text)
	case "x11":
		X11SetClipboardText(text)
	}
}

func ReadClipboard(ds string) string {
	switch ds {
	case "darwin":
		return DarwinGetClipboardText()
	case "wayland":
		wlContent, err := shell.GetWLClipBoard()
		utils.HandleError(err)
		return wlContent
	case "x11":
		return X11GetClipboardText()
	}
	return ""
}

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

// run the listener is the current shell
func RunListener(displayServer string) {
	switch displayServer {
	case "darwin":
		RunDarwinListener()
	case "wayland":
		fmt.Println("Wayland systems use the wl-paste --watch util. See https://github.com/bugaevc/wl-clipboard")
	case "x11":
		RunX11Listener()
	}
}

func SendPaste(keybind, displayServer string) {
	switch displayServer {
	case "wayland":
		utils.LogERROR("auto paste is not yet available for wayland")
	default:
		parts := strings.Split(keybind, "+")
		utils.HandleError(robotgo.KeyTap(parts[1], parts[0]))
	}
}
