package shell

const (
	listenCmd     = "--listen-shell"
	pgrepCmd      = "pgrep -a clipse"
	wlVersionCmd  = "wl-copy"
	wlCopyImgCmd  = "wl-copy -t image/png <"
	wlPasteImgCmd = "wl-paste -t image/png >"
	xVersionCmd   = "xclip -v"
	xCopyImgCmd   = "xclip -selection clipboard -t image/png -i"
	xPasteImgCmd  = "xclip -selection clipboard -t image/png -o >"
)
