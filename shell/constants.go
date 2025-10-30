package shell

const (
	listenCmd        = "-listen"
	listenShellCmd   = "--listen-shell"
	pgrepCmd         = "pgrep 'clipse'"
	psCmd            = "ps -o command"
	wlVersionCmd     = "wl-copy -v"
	wlCopyHandler    = "wl-copy"
	wlPasteHandler   = "wl-paste"
	wlPasteWatcher   = "--watch"
	wlCopyImgCmd     = "wl-copy -t image/png < %s"
	wlPasteImgCmd    = "wl-paste -t image/png > %s"
	wlStoreCmd       = "--wl-store"
	wlTypeSpec       = "--type"
	xVersionCmd      = "xclip -v"
	xCopyImgCmd      = "xclip -selection clipboard -t image/png -i %s"
	xPasteImgCmd     = "xclip -selection clipboard -t image/png -o > %s"
	darwinCopyImgCmd = "osascript -e 'set the clipboard to (read (POSIX file \"%s\") as «class PNGf»)'"
	darwinListenCmd  = "--listen-darwin"
)
