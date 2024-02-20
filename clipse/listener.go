package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

/* runListener is essentially a while loop to be created as a system background process on boot.
   can be stopped at any time with:
   	clipse -kill
   	pkill -f clipse
   	killall clipse
*/

func runListener(fullPath string) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {
			// Get the current clipboard content
			text, err := clipboard.ReadAll()
			handleError(err)
			dataType := checkDataType(text)

			if dataType == "text" {
				// If clipboard content is not empty and not already in the list, add it
				err := addClipboardItem(fullPath, text, "")
				handleError(err)
			} else {
				if imagesEnabled() && (dataType == "png" || dataType == "JPG") {
					randoTimeStamp := fmt.Sprintf("%d", time.Now().UnixNano())
					fileName := fmt.Sprintf("%s/%s.%s", fileDir, randoTimeStamp, dataType) // fileDir defined in constants.go
					saveImage(fileName)
					displayName := fmt.Sprintf("<BINARY FILE> %s", fileName)
					addClipboardItem(fullPath, displayName, fileName)
				}
			}
			time.Sleep(pollInterval) // pollInterval defined in constants.go
		}
	}()
	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	return nil
}
