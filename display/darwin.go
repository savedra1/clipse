package display

import (
	"github.com/savedra1/clipse/handlers"
	"github.com/savedra1/clipse/shell"
)

type DarwinDS struct {
	runtime string
}

func (dds *DarwinDS) Runtime() string {
	return dds.runtime
}

func (dds *DarwinDS) CopyText(text string) {
	handlers.DarwinCopyText(text)
}

func (dds *DarwinDS) CopyImage(filePath string) {
	shell.DarwinCopyImage(filePath)
}

func (dds *DarwinDS) ReadClipboard() string {
	return handlers.DarwinGetClipboardText()
}

func (dds *DarwinDS) RunListener() {
	handlers.RunDarwinListener()
}

func (dds *DarwinDS) RunDetachedListener() {
	shell.RunDarwinListener()
}

func (dds *DarwinDS) Paste() {
	handlers.DarwinPaste()
}

func (dds *DarwinDS) SendPasteKey(keybind string) {
	handlers.RobotPaste(keybind)
}
