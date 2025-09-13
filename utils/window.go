package utils

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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
		}
		return "x11"
	case "darwin":
		return "darwin"
	default:
		return "unknown"
	}
}

func GetActiveWindowTitle() string {
	switch runtime.GOOS {
	case "darwin":
		return getActiveWindowTitleMacOS()
	case "linux":
		return getActiveWindowTitleLinux()
	default:
		LogWARN("Unsupported platform for active window detection: " + runtime.GOOS)
		return ""
	}
}

func getActiveWindowTitleMacOS() string {
	output := execOutput("osascript", "-e", `tell application "System Events" to get name of first application process whose frontmost is true`)
	if output == "" {
		LogWARN("Failed to get active window on macOS")
	}
	return output
}

func getActiveWindowTitleLinux() string {
	if isWaylandSession() {
		if title := tryHyprctl(); title != "" {
			return title
		}
		LogWARN("Failed to get active window on Wayland: no suitable tool found (hyprctl)")
	} else {
		if title := tryXdotool(); title != "" {
			return title
		}
		LogWARN("Failed to get active window on X11: no suitable tool found (xdotool)")
	}
	return ""
}

func isWaylandSession() bool {
	return os.Getenv("XDG_SESSION_TYPE") == "wayland" || os.Getenv("WAYLAND_DISPLAY") != ""
}

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

func tryXdotool() string {
	return execOutput("xdotool", "getactivewindow", "getwindowname")
}

func IsAppExcluded(appName string, excludeList []string) bool {
	if appName == "" {
		return false
	}

	appNameLower := strings.ToLower(appName)

	for _, excluded := range excludeList {
		excludedLower := strings.ToLower(excluded)

		if excludedLower != "" && strings.Contains(appNameLower, excludedLower) {
			return true
		}
	}

	return false
}

func execOutput(name string, args ...string) string {
	if _, err := exec.LookPath(name); err != nil {
		return ""
	}

	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
