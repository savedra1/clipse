package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/atotto/clipboard"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

/*
runListener is essentially a while loop to be created as a system background process on boot.
	can be stopped at any time with:
		clipse -kill
		pkill -f clipse
		killall clipse
*/

var prevClipboardContent string // used to store clipboard content to avoid re-checking media data unnecessarily
var dataType string             // used to determine which poll interval to use based on current clipboard data format

func RunListener(displayServer string, imgEnabled bool) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// channel to pass clipboard events to
	clipboardData := make(chan string, 1)

	// Goroutine to monitor clipboard
	go func() {
		for {
			input, _ := clipboard.ReadAll() // ignoring err here to prevent system crash if input ever not recognized
			if input != prevClipboardContent {
				clipboardData <- input       // Pass clipboard data to main goroutine
				prevClipboardContent = input // update previous content
			}
			if dataType == "text" {
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
			dataType = utils.DataType(input)
			switch dataType {
			case "text":
				if input != "" && !config.Contains(input) {
					if err := config.AddClipboardItem(input, "null"); err != nil {
						utils.LogERROR(fmt.Sprintf("ERROR: failed to add new item `( %s )` | %s", input, err))
					}
				}
			case "png", "jpeg":
				if imgEnabled { // need to add something here to only check the same media image once to save CPU
					fileName := fmt.Sprintf("%s.%s", strconv.Itoa(len(input)), dataType)
					title := fmt.Sprintf("%s %s", imgIcon, fileName)
					if !config.Contains(title) {
						filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

						if err := shell.SaveImage(filePath, displayServer); err != nil {
							fmt.Println("")
							utils.LogERROR(fmt.Sprintf("ERROR: failed to save image | %s", err))
							break
						}
						if err := config.AddClipboardItem(title, filePath); err != nil {
							utils.LogERROR(fmt.Sprintf("ERROR: failed to save image | %s", err))
						}
					}
				}
			}
		case <-interrupt:
			break MainLoop // Exit main loop on interrupt signal
		}
	}

	return nil
}
