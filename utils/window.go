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
	if title := tryWlrctl(); title != "" {
		return title
	}
	LogWARN("Failed to get active window on Wayland: no suitable tool found (Hyprctl, Wlrctl)")
	return ""
}

// getActiveWindowTitleX11 tries getting the window title using various X11 tools
func getActiveWindowTitleX11() string {
	if title := tryXdotool(); title != "" {
		return title
	}
	if title := tryXprop(); title != "" {
		return title
	}
	LogWARN("Failed to get active window on X11: no suitable tool found (Xdotool, Xprop)")
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

// tryXdotool tries getting the window title for X11 systems using xdotool - Command-line X11 automation tool
// xdotool is installed by default on some distros but it is rather uncommon
// Example output: Alacritty
func tryXdotool() string {
	return execOutput("xdotool", "getactivewindow", "getwindowname")
}

// tryXprop tries getting the window title for X11 systems using xprop - property displayer for X
// xprop is widely available on X11 desktop environments
// Example output: _NET_ACTIVE_WINDOW(WINDOW): window id # 0x1a00005 (then) WM_NAME(STRING) = "Alacritty"
func tryXprop() string {
	activeWindowOutput := execOutput("xprop", "-root", "_NET_ACTIVE_WINDOW")

	// Find the "#" and extract the first hex ID after it (there can be more than one, space separated)
	hashIndex := strings.Index(activeWindowOutput, "#")
	if hashIndex == -1 {
		return ""
	}
	hexIds := strings.Fields(activeWindowOutput[hashIndex+1:])
	if len(hexIds) == 0 {
		return ""
	}
	windowID := hexIds[0]

	wmNameOutput := execOutput("xprop", "-id", windowID, "WM_NAME")

	// Return the text between the first and last double-quotes
	firstQuote := strings.Index(wmNameOutput, `"`)
	if firstQuote == -1 {
		return ""
	}
	lastQuote := strings.LastIndex(wmNameOutput, `"`)
	if lastQuote > firstQuote {
		return wmNameOutput[firstQuote+1 : lastQuote]
	}

	return ""
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
