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

// Initialize X11 resources once
static void init_x11() {
    if (dpy != NULL) return;

    dpy = XOpenDisplay(NULL);
    if (!dpy) return;

    win = XCreateSimpleWindow(dpy, DefaultRootWindow(dpy), 0, 0, 1, 1, 0, 0, 0);

    XA_CLIPBOARD = XInternAtom(dpy, "CLIPBOARD", False);
    XA_UTF8_STRING = XInternAtom(dpy, "UTF8_STRING", False);

    XFixesSelectSelectionInput(dpy, win, XA_CLIPBOARD, XFixesSetSelectionOwnerNotifyMask);
}

int hasClipboardChangedX11() {
    init_x11();
    if (!dpy) return 0;

    XEvent ev;

    while (1) {
        XNextEvent(dpy, &ev); // blocks until an event arrives

        if (ev.type == XFixesSelectionNotify) {
            XFixesSelectionNotifyEvent *xfe = (XFixesSelectionNotifyEvent *)&ev;
            long serial = xfe->selection_timestamp;

            if (serial != last_serial) {
                last_serial = serial;
                return 1; // clipboard changed
            }
        }
        // ignore other events
    }

    return 0; // never reached
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

func RunX11Listner() {
	for {
		fmt.Println(C.hasClipboardChangedX11())
		// imgContents, err := GetClipboardImage()
		// textContents := X11GetClipboardText()
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// if imgContents == nil {
		// 	fmt.Println(textContents)
		// } else {
		// 	fmt.Println("<img data>")
		// }
		// if C.hasClipboardChangedX11() == 1 {
		// 	fmt.Printf("Clipborad changed. New value: %s", X11GetClipboardText())
		// } else {
		// 	fmt.Printf("Cliboard contents: %s", X11GetClipboardText())
		// }
		time.Sleep(time.Duration(2000) * time.Millisecond)
	}
}

func GetClipboardImage() ([]byte, error) {
	var outLen C.int

	// Call your C function
	ptr := C.getClipboardImageX11(&outLen)
	if ptr == nil || outLen == 0 {
		return nil, nil // no image
	}

	// Copy into Go memory
	buf := C.GoBytes(unsafe.Pointer(ptr), outLen)

	// Free C memory
	C.free(unsafe.Pointer(ptr))

	return buf, nil
}

// Optional functions
// func XDarwinHasClipboardChanged() bool {
// 	return C.hasClipboardChangedX11() == 1
// }

// func XSetClipboardText(s string) {
// 	cs := C.CString(s)
// 	defer C.free(unsafe.Pointer(cs))
// 	C.setClipboardTextX11(cs)
// }
