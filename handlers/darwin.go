package handlers

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit
#import <AppKit/AppKit.h>
#include <stdlib.h>

static long lastChangeCount = -1;

char* getClipboardText() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    NSString *string = [pasteboard stringForType:NSPasteboardTypeString];
    if (string == nil) {
        return NULL;
    }
    const char* utf8String = [string UTF8String];
    if (utf8String == NULL) {
        return NULL;
    }
    return strdup(utf8String);
}

long getClipboardChangeCount() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    return [pasteboard changeCount];
}

// Returns 1 if changed, 0 if not
int hasClipboardChanged() {
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    long currentCount = [pasteboard changeCount];

    if (currentCount != lastChangeCount) {
        lastChangeCount = currentCount;
        return 1;
    }
    return 0;
}
*/
import "C"
import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

func GetClipboardText() string {
	cstr := C.getClipboardText()
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

func HasClipboardChanged() bool {
	return C.hasClipboardChanged() == 1
}

func RunDarwinListner(displayServer string, imgEnabled bool) error {
	for {
		if HasClipboardChanged() {
			text := GetClipboardText()
			if text != "" {
				if err := config.AddClipboardItem(text, "null"); err != nil {
					utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", text, err))
					return err
				}
				continue
			}

			if !imgEnabled {
				continue
			}

			imgDataPresent, data := shell.DarwinImageDataPresent()
			if !imgDataPresent {
				continue
			}
			// TODO: create func to avoid repeat code
			fileName := fmt.Sprintf("%s-%s.%s", strconv.Itoa(len(data)), utils.GetTimeStamp(), dataType)
			itemTitle := fmt.Sprintf("%s %s", imgIcon, fileName)
			filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

			if err := shell.SaveImage(utils.CleanPath(filePath), displayServer); err != nil {
				utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
				return err
			}
			if err := config.AddClipboardItem(itemTitle, filePath); err != nil {
				utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
				return err
			}

		}

		time.Sleep(defaultPollInterval * time.Millisecond)
	}
}
