package main

import (
	"os"
	"os/signal"
	"strings"
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

	for {
		// Get the current clipboard content
		text, err := clipboard.ReadAll()
		handleError(err)

		// If clipboard content is not empty and not already in the list, add it
		if text != "" && !contains(data.ClipboardHistory, text) {
			// If the length exceeds 50, remove the oldest item
			if len(data.ClipboardHistory) >= 50 {
				lastIndex := len(data.ClipboardHistory) - 1
				data.ClipboardHistory = data.ClipboardHistory[:lastIndex] // Remove the oldest item
			}

			// yyyy-mm-dd hh-mm-s.msmsms Time format
			timeNow := strings.Split(time.Now().UTC().String(), "+0000")[0]

			// {"value": "copied_strig", "recorded": "2024-01-02 12:34:78743687"}
			item := ClipboardItem{Value: text, Recorded: timeNow}

			data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)

			// Save updated data to JSON file
			err = saveDataToFile(fullPath, data)
			handleError(err)

			// Check for updates every 0.1 second
			time.Sleep(pollInterval) // pollInterval defined in constants.go
		}

		// Wait for SIGINT or SIGTERM signal
		<-interrupt
		return nil
	}

}
