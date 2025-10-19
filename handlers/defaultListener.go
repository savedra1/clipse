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
			input, err := getClipboardData(displayServer)
			if err != nil {
				time.Sleep(time.Duration(1 * time.Second)) // wait for boot
			}
			if input != prevClipboardContent {
				clipboardData <- input       // Pass clipboard data to main goroutine
				prevClipboardContent = input // update previous content
			}
			if dataType == Text {
				time.Sleep(time.Duration(config.ClipseConfig.PollInterval))
				continue
			}
			time.Sleep(time.Duration(config.ClipseConfig.PollInterval * 10))
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
			case Text:
				if err := config.AddClipboardItem(input, "null"); err != nil {
					utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", input, err))
				}
			case PNG, JPEG:
				if imgEnabled {
					fileName := fmt.Sprintf("%s-%s.%s", strconv.Itoa(len(input)), utils.GetTimeStamp(), dataType)
					itemTitle := fmt.Sprintf("%s %s", imgIcon, fileName)
					filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

					if err := shell.SaveImage(utils.CleanPath(filePath), displayServer); err != nil {
						utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
						break
					}
					if err := config.AddClipboardItem(itemTitle, filePath); err != nil {
						utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
					}
				}
			}
		case <-interrupt:
			break MainLoop
		}
	}

	return nil
}

// helper func to handler darwin image data
func getClipboardData(ds string) (string, error) {
	if ds != "darwin" {
		return clipboard.ReadAll()
	}

	imgData := shell.DarwinImageDataPresent()
	if imgData == nil {
		return clipboard.ReadAll()
	}
	return string(imgData), nil
}
