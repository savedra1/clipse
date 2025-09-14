package utils

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DisplayServer determines runtime and returns appropriate window server.
// Used to determine the active window and which dependency is required for handling image files.
func DisplayServer() string {
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

// GetActiveWindowTitle returns the title of the currently active window.
func GetActiveWindowTitle() string {
	displayServer := DisplayServer()
	switch displayServer {
	case "darwin":
		return getActiveWindowTitleMacOS()
	case "wayland":
		return getActiveWindowTitleWayland()
	case "x11":
		return getActiveWindowTitleX11()
	default:
		LogWARN("Unsupported display server for active window detection: " + displayServer)
		return ""
	}
}

// IsAppExcluded checks if an application name matches any in the excluded list.
func IsAppExcluded(appName string, excludedList []string) bool {
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

// getActiveWindowTitleMacOS returns the window title on macOS. Should work on any recent system
func getActiveWindowTitleMacOS() string {
	output := execOutput("osascript", "-e", `tell application "System Events" to get name of first application process whose frontmost is true`)
	if output == "" {
		LogWARN("Failed to get active window on macOS")
	}
	return output
}

// getActiveWindowTitleWayland tries getting the window title using various Wayland tools
func getActiveWindowTitleWayland() string {
	if title := tryHyprctl(); title != "" {
		return title
	}
	LogWARN("Failed to get active window on Wayland: no suitable tool found")
	return ""
}

// getActiveWindowTitleX11 tries getting the window title using various X11 tools
func getActiveWindowTitleX11() string {
	if title := tryXdotool(); title != "" {
		return title
	}
	LogWARN("Failed to get active window on X11: no suitable tool found")
	return ""
}

// tryHyprctl tries getting the window title using hyprctl - Utility for controlling parts of Hyprland from a CLI or a script
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

// tryXdotool tries getting the window title using xdotool - Command-line X11 automation tool
func tryXdotool() string {
	return execOutput("xdotool", "getactivewindow", "getwindowname")
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
