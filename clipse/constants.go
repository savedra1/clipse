package main

import "time"

/*
Global vars stored in separate module.
Any new additions to be added here.
*/

const (
	baseDir       = ".config"
	fileName      = "clipboard_history.json"
	configDirName = "clipse"
	fileDir       = "tmp_files"
	listenCmd     = "--listen-shell"
	pollInterval  = 100 * time.Millisecond / 10
	maxLen        = 100
	maxChar       = 125
)
