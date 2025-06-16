package shell

const (
	listenCmd      = "-listen"
	listenShellCmd = "--listen-shell" // internal
	pgrepCmd       = "pgrep -a clipse"
	wlVersionCmd   = "wl-copy -v"
	wlPasteHandler = "wl-paste"
	wlPasteWatcher = "--watch"
	wlCopyImgCmd   = "wl-copy -t image/png <"
	wlPasteImgCmd  = "wl-paste -t image/png >"
	wlStoreCmd     = "--wl-store" // internal
	wlTypeSpec     = "--type"
	xVersionCmd    = "xclip -v"
	xCopyImgCmd    = "xclip -selection clipboard -t image/png -i"
	xPasteImgCmd   = "xclip -selection clipboard -t image/png -o >"
)
