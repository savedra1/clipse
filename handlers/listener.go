package handlers

import (
	"clipse/config"
	"clipse/shell"
	"clipse/utils"

	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

/*
runListener is essentially a while loop to be created as a system background process on boot.

	can be stopped at any time with:
		clipse -kill
		pkill -f clipse
		killall clipse
*/

func bootLoaded() bool {
	/*
		System fails to read clipboard data when run on boot.
		Needs a buffer period to continue.
	*/
	var loaded bool
	startTime := time.Now()

	for {
		if time.Since(startTime) >= 60*time.Second {
			loaded = false
			break
		}

		_, err := clipboard.ReadAll()
		if err == nil {
			loaded = true
			break
		}

		time.Sleep(time.Second)
	}
	return loaded
}

func RunListener(historyFilePath, clipsDir, displayServer string, imgEnabled bool) error {
	if !bootLoaded() {
		time.Sleep(30 * time.Second) // Account for extra slow boot loaders
	}
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {
			// Get the current clipboard content
			input, err := clipboard.ReadAll()
			utils.HandleError(err)
			dt := utils.DataType(input)

			switch dt {
			case "text":
				if input != "" && !config.Contains(input) {
					err = config.AddClipboardItem(historyFilePath, input, "null")
					utils.HandleError(err)
				}
			case "png", "jpeg":
				if imgEnabled {
					file := fmt.Sprintf("%s.%s", strconv.Itoa(len(input)), dt)
					filePath := filepath.Join(clipsDir, tmpDir, file)
					title := fmt.Sprintf("<BINARY FILE> %s", file)
					if !config.Contains(title) {
						err = shell.SaveImage(filePath, displayServer)
						utils.HandleError(err)

						err = config.AddClipboardItem(historyFilePath, title, filePath)
						utils.HandleError(err)
					}
				}
			}

			time.Sleep(pollInterval) // pollInterval defined in constants.go
		}
	}()

	<-interrupt
	return nil
}
