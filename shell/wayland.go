package shell

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/savedra1/clipse/utils"
)

func GetWLClipBoard() (string, error) {
	cmd := exec.Command(wlPasteHandler)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func UpdateWLClipboard(s string) error {
	cmd := exec.Command(wlCopyHandler, "--", s)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func WLDependencyCheck() error {
	cmd := exec.Command("which", wlCopyHandler)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func RunWaylandListener() {
	for _, i := range []string{"image/png", "text"} {
		cmd := exec.Command(wlPasteHandler, wlTypeSpec, i, wlPasteWatcher, ExeName, wlStoreCmd)
		runDetachedCmd(cmd)
	}
}

func WLCopyImage(filePath string) {
	cmdFull := fmt.Sprintf(wlCopyImgCmd, filePath)
	if err := exec.Command("sh", "-c", cmdFull).Run(); err != nil {
		utils.LogERROR(fmt.Sprintf("failed to copy image: %s", err))
	}
}

func WLSaveImage(imagePath string) error {
	cmdFull := fmt.Sprintf(wlPasteImgCmd, imagePath)
	if err := exec.Command("sh", "-c", cmdFull).Run(); err != nil {
		return err
	}
	return nil
}

// getActiveWindowTitleWayland tries getting the window title using various Wayland tools
func WLActiveWindowTitle() string {
	if title := tryWlrctl(); title != "" {
		return title
	}
	if title := tryHyprctl(); title != "" {
		return title
	}
	utils.LogWARN("Failed to get active window on Wayland: no suitable tool found (Hyprctl, Wlrctl)")
	return ""
}

// tryWlrctl tries getting the window title for wl-roots compositors using wlrctl - Utility for miscellaneous wlroots extensions
// wlrctl can be installed on any wl-roots system but may not be installed by default
// Example output: Alacritty: user@arch:~
func tryWlrctl() string {
	output := execOutput("wlrctl", "toplevel", "list", "state:focused")
	lines := strings.Split(output, "\n")

	// Find first non-empty line and process it
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		// Get first word and remove trailing colon
		appName := strings.TrimSuffix(fields[0], ":")
		return appName
	}

	return ""
}

// tryHyprctl tries getting the window title for Hyprland using hyprctl - Utility for controlling parts of Hyprland from a CLI or a script
// hyprctl is typically installed by default on Hyprland distros
// Example output: { "class": "Alacritty", "title": "user@arch:~", ... }
func tryHyprctl() string {
	output := execOutput("hyprctl", "activewindow", "-j")
	if output == "" {
		return ""
	}

	var windowInfo struct {
		Class string `json:"class"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal([]byte(output), &windowInfo); err != nil {
		return ""
	}

	if windowInfo.Class != "" {
		return windowInfo.Class
	}
	return windowInfo.Title
}
