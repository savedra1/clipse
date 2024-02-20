package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
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

func runListener(fullPath, fileDir string) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {
			// Get the current clipboard content
			input, err := clipboard.ReadAll()
			handleError(err)

			dataType := checkDataType(input)

			switch dataType {
			case "text":
				if input != "" && !contains(input) {
					err := addClipboardItem(fullPath, input, "null")
					handleError(err)
				}
			case "png", "jpeg":
				if imagesEnabled() {
					file := fmt.Sprintf("%s.%s", strconv.Itoa(len(input)), dataType) // fileDir defined in constants.go
					filePath := filepath.Join(fileDir, file)
					title := fmt.Sprintf("<BINARY FILE> %s", file)
					if !contains(title) {
						err = saveImage(filePath)
						handleError(err)
						err = addClipboardItem(fullPath, title, filePath)
						handleError(err)
					}
				}

			}
			time.Sleep(pollInterval) // pollInterval defined in constants.go
		}

	}()

	//time.Sleep(pollInterval) // pollInterval defined in constants.go
	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	return nil
}

/*
func runListener(fullPath, fp string) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any

	fmt.Println("DEBUG 1")
	time.Sleep(2 * time.Second)

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {
			// Get the current clipboard content
			text, err := clipboard.ReadAll()
			fmt.Println("Input:", text)
			fmt.Println("DEBUG 2")
			time.Sleep(2 * time.Second)

			handleError(err)

			// If clipboard content is not empty and not already in the list, add it
			if text != "" && !contains(text) {
				fmt.Println("DEBUG 3")
				time.Sleep(2 * time.Second)
				err := addClipboardItem(fullPath, text, "null")
				handleError(err)
				fmt.Println("DEBUG 4")
				time.Sleep(2 * time.Second)
			}
			time.Sleep(pollInterval) // pollInterval defined in constants.go

		}
	}()
	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	return nil
}
*/
