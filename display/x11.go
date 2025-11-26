package display

import (
	"github.com/savedra1/clipse/handlers"
)

type XDS struct {
	runtime string
}

func (xds *XDS) Runtime() string {
	return xds.runtime
}

func (xds *XDS) CopyText(text string) {
	handlers.X11SetClipboardText(text)
}

func (xds *XDS) ReadClipboard() string {
	return handlers.X11GetClipboardText()
}

func (xds *XDS) RunListener() {
	handlers.RunX11Listener()
}

func (xds *XDS) Paste() {
	handlers.X11Paste()
}

func (xds *XDS) SendPasteKey(keybind string) {
	handlers.RobotPaste(keybind)
}
