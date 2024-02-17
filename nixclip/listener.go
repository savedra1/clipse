package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

/* runListener is essentially a while loop to be created as a system background process on boot.
   can be stopped at any time with:
   	clipboard kill
   	pkill -f clipboard
   	killall clipboard
*/

func runListener(fullPath string) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any
	var data ClipboardHistory

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {
			// Get the current clipboard content
			text, err := clipboard.ReadAll()
			handleError(err)
			if !isFile(text) {
				// If clipboard content is not empty and not already in the list, add it
				if text != "" && !contains(data.ClipboardHistory, text) {
					err := addClipboardItem(fullPath, text)
					handleError(err)
				}
				time.Sleep(pollInterval) // pollInterval defined in constants.go

			}

		}
	}()
	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	return nil
}
