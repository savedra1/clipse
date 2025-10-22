// handlers/darwin.go
//go:build darwin && cgo

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

unsigned char* getClipboardImage(int *outLen) {
    NSPasteboard *pb = [NSPasteboard generalPasteboard];
    NSData *data = [pb dataForType:NSPasteboardTypePNG];
    if (data == nil) {
        return NULL;
    }
    *outLen = (int)[data length];
    unsigned char *buffer = (unsigned char *)malloc(*outLen);
    if (buffer == NULL) {
        *outLen = 0;
        return NULL;
    }
    memcpy(buffer, [data bytes], *outLen);
    return buffer;
}

// returns 1 if text, 2 if image, 0 if unknown/empty
int getClipboardType() {
    NSPasteboard *pb = [NSPasteboard generalPasteboard];
    NSArray *types = [pb types];
    if ([types containsObject:NSPasteboardTypeString]) {
        return 1;
    }
    if ([types containsObject:NSPasteboardTypePNG] || [types containsObject:NSPasteboardTypeTIFF]) {
        return 2;
    }
    return 0;
}

// set clipboard content
void setClipboardText(const char* text) {
    if (text == NULL) {
        return;
    }

    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    [pasteboard clearContents];

    NSString *string = [NSString stringWithUTF8String:text];
    if (string != nil) {
        [pasteboard setString:string forType:NSPasteboardTypeString];
    }
}
*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unsafe"

	"github.com/savedra1/clipse/config"
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

func readClipboardImage() []byte {
	var length C.int
	data := C.getClipboardImage(&length)
	if data == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(data))
	return C.GoBytes(unsafe.Pointer(data), length)
}

func saveDarwinImage(imgData []byte) error {
	byteLength := strconv.Itoa(len(string(imgData)))
	fileName := fmt.Sprintf("%s-%s.png", byteLength, utils.GetTimeStamp())
	itemTitle := fmt.Sprintf("%s %s", imgIcon, fileName)
	filePath := filepath.Join(config.ClipseConfig.TempDirPath, fileName)

	if err := os.WriteFile(filePath, imgData, 0644); err != nil {
		return err
	}

	if err := config.AddClipboardItem(itemTitle, filePath); err != nil {
		return err
	}
	return nil
}

func saveDarwinText(textData string) error {
	if err := config.AddClipboardItem(textData, "null"); err != nil {
		return err
	}
	return nil
}

func RunDarwinListener(displayServer string, imgEnabled bool) error {
	var prevText string
	var prevImg []byte
	for {
		if HasClipboardChanged() {
			// Check if the clipboard content should be excluded based on source application
			activeWindow := utils.GetActiveWindowTitle()
			if utils.IsAppExcluded(activeWindow, config.ClipseConfig.ExcludedApps) {
				utils.LogINFO(fmt.Sprintf("Skipping clipboard content from excluded app: %s", activeWindow))
				continue
			}

			clipboardType := C.getClipboardType()

			switch clipboardType {
			case 1: // text
				text := GetClipboardText()
				if text == prevText {
					break
				}
				prevText = text
				if err := saveDarwinText(text); err != nil {
					utils.LogERROR(fmt.Sprintf("failed to add new item `( %s )` | %s", text, err))
				}

			case 2: // image
				img := readClipboardImage()
				if string(img) == string(prevImg) {
					break
				}
				prevImg = img
				if err := saveDarwinImage(img); err != nil {
					utils.LogERROR(fmt.Sprintf("failed to save image | %s", err))
				}

			default:
				utils.LogWARN("Unknown data type found in darwin clipboard")
			}
		}
		time.Sleep(time.Duration(config.ClipseConfig.PollInterval) * time.Millisecond)
	}
}

func DarwinPaste() error {
	clipboardType := C.getClipboardType()

	switch clipboardType {
	case 1: // text
		_, err := fmt.Println(GetClipboardText())
		utils.HandleError(err)

	case 2: // image
		img := readClipboardImage()
		_, err := fmt.Println(string(img))
		utils.HandleError(err)
	}

	return nil
}

func DarwinCopyText(s string) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	C.setClipboardText(cstr)
}
