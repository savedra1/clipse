package shell

const (
	listenCmd          = "-listen"
	listenShellCmd     = "--listen-shell"
	pgrepCmd           = "pgrep 'clipse'"
	psCmd              = "ps -o command"
	wlCopyHandler      = "wl-copy"
	wlPasteHandler     = "wl-paste"
	wlPasteWatcher     = "--watch"
	wlCopyImgCmd       = "wl-copy -t image/png < %s"
	wlPasteImgCmd      = "wl-paste -t image/png > %s"
	wlStoreCmd         = "--wl-store"
	wlTypeSpec         = "--type"
	darwinCopyImgCmd   = "osascript -e 'set the clipboard to (read (POSIX file \"%s\") as «class PNGf»)'"
	darwinGetWindowCmd = `tell application "System Events" to get name of first application process whose frontmost is true`
	darwinListenCmd    = "--listen-darwin"
	x11ListenCmd       = "--listen-x11"
)
