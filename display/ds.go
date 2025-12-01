package display

import (
	"fmt"
	"os"
	"runtime"

	"github.com/savedra1/clipse/utils"
)

var DisplayServer = GetDisplayServer()

type DS interface {
	Runtime() string
	ReadClipboard() string
	CopyText(string)
	CopyImage(string)
	Paste()
	RunListener()
	RunDetachedListener()
	SendPasteKey(string)
}

func GetDisplayServer() DS {
	osName := runtime.GOOS
	switch osName {
	case "linux":
		waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
		if waylandDisplay != "" {
			return &WaylandDS{
				runtime: "wayland",
			}
		}
		return &XDS{
			runtime: "x11",
		}
	case "darwin":
		return &DarwinDS{
			runtime: "darwin",
		}
	default:
		utils.LogERROR(fmt.Sprintf("display server not recognized: %s", osName))
		os.Exit(1)
	}
	return nil
}
