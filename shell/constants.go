package shell

const (
	listenCmd         = "-listen"
	listenShellCmd    = "--listen-shell" // internal
	pgrepCmd          = "pgrep 'clipse'"
	psCmd             = "ps -o command"
	wlVersionCmd      = "wl-copy -v"
	wlPasteHandler    = "wl-paste"
	wlPasteWatcher    = "--watch"
	wlCopyImgCmd      = "wl-copy -t image/png < %s"
	wlPasteImgCmd     = "wl-paste -t image/png > %s"
	wlStoreCmd        = "--wl-store" // internal
	wlTypeSpec        = "--type"
	xVersionCmd       = "xclip -version"
	xCopyImgCmd       = "xclip -selection clipboard -t image/png -i %s"
	xPasteImgCmd      = "xclip -selection clipboard -t image/png -o > %s"
	darwinVersionCmd  = "pngpaste -v"
	darwinCopyImgCmd  = "osascript -e 'set the clipboard to (read (POSIX file \"%s\") as «class PNGf»)'"
	darwinPasteImgCmd = "pngpaste %s"
	darwinImgCheckCmd = "pngpaste -"
)
