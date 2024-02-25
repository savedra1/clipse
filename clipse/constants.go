package main

import "time"

/*
Global vars stored in separate module.
Any new additions to be added here.
*/

const (
	baseDir         = ".config"
	historyFileName = "clipboard_history.json"
	themeFile       = "custom_theme.json"
	clipseDirName   = "clipse"
	tmpDir          = "tmp_files"
	listenCmd       = "--listen-shell"
	pollInterval    = 100 * time.Millisecond / 10
	maxLen          = 100
	maxChar         = 100
)
