package display

import (
	"fmt"

	"github.com/savedra1/clipse/handlers"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

type WaylandDS struct {
	runtime string
}

func (wds *WaylandDS) Runtime() string {
	return wds.runtime
}

func (wds *WaylandDS) CopyText(text string) {
	handlers.WaylandCopy(text)
}

func (wds *WaylandDS) CopyImage(filePath string) {
	shell.WLCopyImage(filePath)
}

func (wds *WaylandDS) ReadClipboard() string {
	wlContent, err := shell.GetWLClipBoard()
	utils.HandleError(err)
	return wlContent
}

func (wds *WaylandDS) RunListener() {
	fmt.Println("Wayland systems use the `wl-paste --watch` util. See https://github.com/bugaevc/wl-clipboard")
}

func (wds *WaylandDS) RunDetachedListener() {
	shell.RunWaylandListener()
}

func (wds *WaylandDS) Paste() {
	handlers.WaylandPaste()
}

func (wds WaylandDS) SendPasteKey(keybind string) {
	handlers.UinputPaste(keybind)
}
