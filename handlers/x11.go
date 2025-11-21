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

// Initialize X11 resources once
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

// Clipboard data holder
static unsigned char *clipboard_data = NULL;
static int clipboard_data_len = 0;
static Atom clipboard_data_type;

// Selection handler - responds to requests for our clipboard data
static Bool handleSelectionRequest(XEvent *ev) {
    XSelectionRequestEvent *req = &ev->xselectionrequest;
    XSelectionEvent notify;

    notify.type = SelectionNotify;
    notify.requestor = req->requestor;
    notify.selection = req->selection;
    notify.target = req->target;
    notify.time = req->time;
    notify.property = None;

    Atom TARGETS = XInternAtom(dpy, "TARGETS", False);

    // Handle TARGETS request - tell what formats we support
    if (req->target == TARGETS) {
        Atom supported[] = { clipboard_data_type, TARGETS };
        XChangeProperty(dpy, req->requestor, req->property,
                       XA_ATOM, 32, PropModeReplace,
                       (unsigned char*)supported, 2);
        notify.property = req->property;
    }
    // Handle request for our actual data
    else if (req->target == clipboard_data_type && clipboard_data != NULL) {
        XChangeProperty(dpy, req->requestor, req->property,
                       clipboard_data_type, 8, PropModeReplace,
                       clipboard_data, clipboard_data_len);
        notify.property = req->property;
    }

    XSendEvent(dpy, req->requestor, False, 0, (XEvent*)&notify);
    XFlush(dpy);
    return True;
}

// Set text to clipboard
int setClipboardTextX11(const char *text) {
    init_x11();
    if (!dpy || !text) return 0;

    // Free old data
    if (clipboard_data) {
        free(clipboard_data);
        clipboard_data = NULL;
    }

    // Copy text data
    clipboard_data_len = strlen(text);
    clipboard_data = malloc(clipboard_data_len);
    memcpy(clipboard_data, text, clipboard_data_len);
    clipboard_data_type = XA_UTF8_STRING;

    // Take ownership of clipboard
    XSetSelectionOwner(dpy, XA_CLIPBOARD, win, CurrentTime);
    XFlush(dpy);

    // Verify we own it
    if (XGetSelectionOwner(dpy, XA_CLIPBOARD) != win) {
        free(clipboard_data);
        clipboard_data = NULL;
        return 0;
    }

    // Process selection requests for a bit to let other apps grab the data
    // This is a simple approach - a proper implementation would do this in the background
    time_t start = time(NULL);
    while (time(NULL) - start < 1) {
        while (XPending(dpy)) {
            XEvent ev;
            XNextEvent(dpy, &ev);
            if (ev.type == SelectionRequest) {
                handleSelectionRequest(&ev);
            }
        }
        usleep(10000); // 10ms
    }

    return 1;
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
			// Error
			fmt.Println("Error waiting for clipboard change")
			time.Sleep(1 * time.Second)
		}
	}
}

// Alternative: polling approach (less efficient but simpler)
func RunX11ListenerPolling() {
	fmt.Println("Starting X11 clipboard monitor (polling)...")

	for {
		if X11ClipboardChanged() {
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
		}

		time.Sleep(100 * time.Millisecond)
	}
}

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

func X11SetClipboardText(text string) error {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	if C.setClipboardTextX11(cstr) == 0 {
		return fmt.Errorf("failed to set clipboard text")
	}
	return nil
}
