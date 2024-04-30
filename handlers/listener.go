package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"

	"github.com/atotto/clipboard"
)

/*
runListener is essentially a while loop to be created as a system background process on boot.
	can be stopped at any time with:
		clipse -kill
		pkill -f clipse
		killall clipse
*/

var prevClipboardContent string // used to store clipboard content to avoid re-checking media data unnecessarily
var currentDataType string      // used to determine which poll interval to use based on current clipboard data format

func RunListener(clipsDir, displayServer string, imgEnabled bool) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// channel to pass clipboard events to
	clipboardData := make(chan string, 1)

	// Goroutine to monitor clipboard
	go func() {
		for {
			input, _ := clipboard.ReadAll() // ignoring err here to prevent system crash if input ever not recognised
			if input != prevClipboardContent {
				clipboardData <- input // Pass clipboard data to main goroutine
			}
			if currentDataType == "text" {
				time.Sleep(defaultPollInterval)
			} else {
				time.Sleep(mediaPollInterval)
			}
		}
	}()

MainLoop:
	for {
		select {
		case input := <-clipboardData:
			if input == "" {
				continue
			}
			dt := utils.DataType(input)
			currentDataType = dt // update var to determine interval time
			switch dt {
			case "text":
				if input != "" && !config.Contains(input) {
					err := config.AddClipboardItem(input, "null")
					utils.HandleError(err)
				}
			case "png", "jpeg":
				if imgEnabled { // need to add something here to only check the same media image once to save CPU
					fileName := fmt.Sprintf("%s.%s", strconv.Itoa(len(input)), dt)
					title := fmt.Sprintf("%s %s", imgIcon, fileName)
					if !config.Contains(title) {
						filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)
						err := shell.SaveImage(filePath, displayServer)
						if err != nil {
							fmt.Println("failed to save media data to tmp dir")
						}
						err = config.AddClipboardItem(title, filePath)
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

/*
Function to explicity await boot is no longer required as err returned
from clipboard read operation can be ignored in Mainloop

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
