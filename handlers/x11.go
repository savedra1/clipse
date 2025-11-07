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
	"unsafe"
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
