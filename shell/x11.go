package shell

import (
	"os/exec"
	"strings"

	"github.com/savedra1/clipse/utils"
)

func RunX11Listener() {
	cmd := exec.Command(ExeName, x11ListenCmd)
	runDetachedCmd(cmd)
}

// getActiveWindowTitleX11 tries getting the window title using various X11 tools
func X11ActiveWindowTitle() string {
	if title := tryXdotool(); title != "" {
		return title
	}
	if title := tryXprop(); title != "" {
		return title
	}
	utils.LogWARN("Failed to get active window on X11: no suitable tool found (Xdotool, Xprop)")
	return ""
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

// tryXdotool tries getting the window title for X11 systems using xdotool - Command-line X11 automation tool
// xdotool is installed by default on some distros but it is rather uncommon
// Example output: Alacritty
func tryXdotool() string {
	return execOutput("xdotool", "getactivewindow", "getwindowname")
}
