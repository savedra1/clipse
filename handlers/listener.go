package handlers

import (
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"

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

func RunListener(historyFilePath, clipsDir, displayServer string, imgEnabled bool) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// channel to pass clipboard events to
	clipboardData := make(chan string, 1)

	// Goroutine to monitor clipboard
	go func() {
		for {
			input, _ := clipboard.ReadAll() // ignoring err here to prevent system crash if input ever not recognised
			clipboardData <- input          // Pass clipboard data to main goroutine
			time.Sleep(pollInterval)        // pollInterval defined in constants.go
		}
	}()

MainLoop:
	for {
		select {
		case input := <-clipboardData:
			dt := utils.DataType(input)

			switch dt {
			case "text":
				if input != "" && !config.Contains(input) {
					err := config.AddClipboardItem(historyFilePath, input, "null")
					utils.HandleError(err)
				}
			case "png", "jpeg":
				if imgEnabled {
					file := fmt.Sprintf("%s.%s", strconv.Itoa(len(input)), dt)
					filePath := filepath.Join(clipsDir, tmpDir, file)
					title := fmt.Sprintf("%s %s", imgIcon, file)
					if !config.Contains(title) {
						err := shell.SaveImage(filePath, displayServer)
						utils.HandleError(err)

						err = config.AddClipboardItem(historyFilePath, title, filePath)
						utils.HandleError(err)
					}
				}
			}
		case <-interrupt:
			break MainLoop // Exit main loop on interrupt signal
		}
	}

	return nil
}

// Function to explicity await boot is no longer required as err returned
// from clipboard read operation can be ignored in Mainloop
/*
func bootLoaded() bool {
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
*/
