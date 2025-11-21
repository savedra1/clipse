//go:build linux && !wayland
// +build linux,!wayland

package handlers

/*
#cgo pkg-config: x11 xfixes
#include <sys/types.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/Xfixes.h>
#include <stdlib.h>
#include <string.h>

static Display *dpy = NULL;
static Window win;
static Atom XA_CLIPBOARD;
static Atom XA_UTF8_STRING;
static long last_serial = -1;
static int xfixes_event_base = 0;
static unsigned char *clipboard_data = NULL;

// Initialize X11 resources
static void init_x11() {
    if (dpy != NULL) return;

    dpy = XOpenDisplay(NULL);
    if (!dpy) return;

    int xfixes_error_base;
    if (!XFixesQueryExtension(dpy, &xfixes_event_base, &xfixes_error_base)) {
        XCloseDisplay(dpy);
        dpy = NULL;
        return;
    }

    win = XCreateSimpleWindow(dpy, DefaultRootWindow(dpy), 0, 0, 1, 1, 0, 0, 0);

    XA_CLIPBOARD = XInternAtom(dpy, "CLIPBOARD", False);
    XA_UTF8_STRING = XInternAtom(dpy, "UTF8_STRING", False);

    XFixesSelectSelectionInput(dpy, win, XA_CLIPBOARD, XFixesSetSelectionOwnerNotifyMask);

    // Flush to ensure the request is sent
    XFlush(dpy);
}

// Returns the X11 connection file descriptor for select/poll
int getX11ConnectionFd() {
    init_x11();
    if (!dpy) return -1;
    return ConnectionNumber(dpy);
}

// Returns 1 if clipboard has changed since last check, 0 otherwise
int hasClipboardChangedX11() {
    init_x11();
    if (!dpy) return 0;

    XEvent ev;
    int changed = 0;

    // Process all pending events
    while (XPending(dpy)) {
        XNextEvent(dpy, &ev);

        // XFixes events start at xfixes_event_base
        if (ev.type == xfixes_event_base + XFixesSelectionNotify) {
            XFixesSelectionNotifyEvent *xfe = (XFixesSelectionNotifyEvent *)&ev;
            if (xfe->selection == XA_CLIPBOARD) {
                long serial = xfe->selection_timestamp;
                if (serial != last_serial) {
                    last_serial = serial;
                    changed = 1;
                }
            }
        }
    }

    return changed;
}

// Blocking wait for clipboard change (with timeout in milliseconds)
// Returns 1 if changed, 0 if timeout, -1 on error
int waitForClipboardChange(int timeout_ms) {
    init_x11();
    if (!dpy) return -1;

    int fd = ConnectionNumber(dpy);
    fd_set fds;
    struct timeval tv;

    tv.tv_sec = timeout_ms / 1000;
    tv.tv_usec = (timeout_ms % 1000) * 1000;

    FD_ZERO(&fds);
    FD_SET(fd, &fds);

    int ret = select(fd + 1, &fds, NULL, NULL, &tv);
    if (ret > 0) {
        return hasClipboardChangedX11();
    }

    return ret; // 0 = timeout, -1 = error
}

// Returns clipboard text (UTF-8) or NULL
char* getClipboardTextX11() {
    init_x11();
    if (!dpy) return NULL;

    Atom sel = XA_CLIPBOARD;
    Atom target = XA_UTF8_STRING;

    XConvertSelection(dpy, sel, target, target, win, CurrentTime);
    XFlush(dpy);

    XEvent ev;
    XNextEvent(dpy, &ev);

    if (ev.type != SelectionNotify) return NULL;
    if (ev.xselection.property == None) return NULL;

    Atom type;
    int format;
    unsigned long len, bytes_left;
    unsigned char *data = NULL;

    XGetWindowProperty(dpy, win, target, 0, ~0, False,
                       AnyPropertyType, &type, &format,
                       &len, &bytes_left, &data);

    if (!data) return NULL;

    char *out = strdup((char*)data);
    XFree(data);
    return out;
}

unsigned char* getClipboardImageX11(int *out_len) {
    init_x11();
    if (!dpy) return NULL;

    *out_len = 0;

    Atom sel = XA_CLIPBOARD;

    // preferred MIME targets (BMP removed)
    Atom PNG  = XInternAtom(dpy, "image/png", False);
    Atom JPEG = XInternAtom(dpy, "image/jpeg", False);

    Atom targets[] = { PNG, JPEG };
    const int ntargets = sizeof(targets) / sizeof(targets[0]);

    for (int i = 0; i < ntargets; i++) {
        Atom target = targets[i];

        // Ask clipboard owner to convert to requested type
        XConvertSelection(dpy, sel, target, target, win, CurrentTime);
        XFlush(dpy);

        // Wait for the SelectionNotify event
        XEvent ev;
        XNextEvent(dpy, &ev);

        if (ev.type != SelectionNotify)
            continue;

        if (ev.xselection.property == None)
            continue;

        Atom type;
        int format;
        unsigned long len, bytes_left;
        unsigned char *data = NULL;

        if (XGetWindowProperty(dpy, win, target, 0, ~0, False,
                               AnyPropertyType, &type, &format,
                               &len, &bytes_left, &data) != Success) {
            continue;
        }

        if (!data || len == 0) {
            if (data) XFree(data);
            continue;
        }

        // Copy result to malloc'd buffer (Go will free this)
        unsigned char *copy = malloc(len);
        memcpy(copy, data, len);
        XFree(data);

        *out_len = (int)len;
        return copy;
    }

    return NULL; // neither PNG nor JPEG available
}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

var clipboardContents string

func X11GetClipboardText() string {
	cstr := C.getClipboardTextX11()
	if cstr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

func X11ClipboardChanged() bool {
	return C.hasClipboardChangedX11() != 0
}

// Efficient listener using blocking waits
func RunX11Listner() {
	fmt.Println("Starting X11 clipboard monitor...")

	for {
		// Wait up to 1 second for clipboard change
		result := int(C.waitForClipboardChange(1000))

		if result > 0 {
			// Clipboard changed
			imgContents, err := GetClipboardImage()
			if err != nil {
				fmt.Printf("Error getting clipboard image: %v\n", err)
			}

			if imgContents != nil {
				fmt.Printf("Clipboard changed - Image detected (%d bytes)\n", len(imgContents))
			} else {
				textContents := X11GetClipboardText()
				fmt.Printf("Clipboard changed - Text: %s\n", textContents)
			}
		} else if result == 0 {
			// Timeout - no change, this is normal
		} else {
			fmt.Println("Error waiting for clipboard change")
			time.Sleep(1 * time.Second)
		}
	}
}

// Alternative: polling approach (less efficient but simpler)
// func RunX11ListenerPolling() {
// 	fmt.Println("Starting X11 clipboard monitor (polling)...")

// 	for {
// 		if X11ClipboardChanged() {
// 			imgContents, err := GetClipboardImage()
// 			if err != nil {
// 				fmt.Printf("Error getting clipboard image: %v\n", err)
// 			}

// 			if imgContents != nil {
// 				fmt.Printf("Clipboard changed - Image detected (%d bytes)\n", len(imgContents))
// 			} else {
// 				textContents := X11GetClipboardText()
// 				fmt.Printf("Clipboard changed - Text: %s\n", textContents)
// 			}
// 		}

// 		time.Sleep(100 * time.Millisecond)
// 	}
// }

func GetClipboardImage() ([]byte, error) {
	var outLen C.int

	ptr := C.getClipboardImageX11(&outLen)
	if ptr == nil || outLen == 0 {
		return nil, nil
	}

	buf := C.GoBytes(unsafe.Pointer(ptr), outLen)
	C.free(unsafe.Pointer(ptr))

	return buf, nil
}

func X11Paste() error {
	imgContents, err := GetClipboardImage()
	if err != nil {
		return err
	}
	if imgContents != nil {
		fmt.Println(string(imgContents))
		return nil
	}

	textContents := X11GetClipboardText()
	fmt.Println(textContents)

	return nil
}

// func X11CopyText(s string) error {
// 	cstr := C.CString(s)
// 	defer C.free(unsafe.Pointer(cstr))

// 	if C.setClipboardTextX11(cstr) == 0 {
// 		return fmt.Errorf("failed to set clipboard text")
// 	}
// 	return nil
// }
