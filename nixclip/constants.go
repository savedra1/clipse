package main

import "time"

/*
Global vars stored in separate module.
Any new additions to be added here.
*/

const (
	fileName      = "clipboard_history.json"
	configDirName = "clipboard_manager"
	listenCmd     = "--listen-shell"
	pollInterval  = 100 * time.Millisecond / 10
	maxLen        = 100
)
