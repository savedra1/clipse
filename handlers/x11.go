//go:build linux && !wayland
// +build linux,!wayland

package handlers

/*
#cgo pkg-config: x11 xfixes
#include <stdlib.h>
#include <string.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/extensions/Xfixes.h>

static Display *dpy = NULL;
static Window win;

// Atoms
static Atom XA_CLIPBOARD;
static Atom XA_UTF8_STRING;
static Atom XA_STRING;
static Atom TARGETS;
static Atom PNG;
static Atom JPEG;

static long last_serial = -1;   // For XFixes change detection
static int last_type = 0;       // Cached clipboard type

// Clipboard types:
#define CLIP_NONE   0
#define CLIP_TEXT   1
#define CLIP_IMAGE  2
#define CLIP_OTHER  3

// -----------------------------------------------------------
// Init X11
// -----------------------------------------------------------
static void init_x11() {
    if (dpy) return;

    dpy = XOpenDisplay(NULL);
    if (!dpy) return;

    win = XCreateSimpleWindow(dpy, DefaultRootWindow(dpy), 0, 0, 1, 1, 0, 0, 0);

    XA_CLIPBOARD   = XInternAtom(dpy, "CLIPBOARD", False);
    XA_UTF8_STRING = XInternAtom(dpy, "UTF8_STRING", False);
    XA_STRING      = XInternAtom(dpy, "STRING", False);
    TARGETS        = XInternAtom(dpy, "TARGETS", False);
    PNG            = XInternAtom(dpy, "image/png", False);
    JPEG           = XInternAtom(dpy, "image/jpeg", False);

    // Listen for clipboard owner changes
    XFixesSelectSelectionInput(
        dpy, win, XA_CLIPBOARD, XFixesSetSelectionOwnerNotifyMask
    );
}

// -----------------------------------------------------------
// Check if clipboard changed (XFixes)
// -----------------------------------------------------------
int hasClipboardChangedX11() {
    init_x11();
    if (!dpy) return 0;

    XEvent ev;

    int changed = 0;

    while (XPending(dpy)) {
        XNextEvent(dpy, &ev);

        if (ev.type == XFixesSelectionNotify) {
            XFixesSelectionNotifyEvent *xfe =
                (XFixesSelectionNotifyEvent *)&ev;

            long serial = xfe->selection_timestamp;

            if (serial != last_serial) {
                last_serial = serial;
                changed = 1;
            }
        }
    }

    return changed;
}

// -----------------------------------------------------------
// Determine clipboard type: TEXT, IMAGE, OTHER
// -----------------------------------------------------------
int getClipboardTypeX11() {
    init_x11();
    if (!dpy) return CLIP_NONE;

    // Request TARGETS
    XConvertSelection(dpy, XA_CLIPBOARD, TARGETS, TARGETS, win, CurrentTime);
    XFlush(dpy);

    XEvent ev;
    XNextEvent(dpy, &ev);
    if (ev.type != SelectionNotify || ev.xselection.property == None)
        return CLIP_NONE;

    Atom type;
    int format;
    unsigned long len, left;
    Atom *list = NULL;

    XGetWindowProperty(
        dpy, win, TARGETS, 0, ~0, False, XA_ATOM,
        &type, &format, &len, &left, (unsigned char**)&list
    );

    if (!list || len == 0) {
        if (list) XFree(list);
        return CLIP_NONE;
    }

    int found_text = 0;
    int found_png  = 0;
    int found_jpeg = 0;

    for (unsigned long i = 0; i < len; i++) {
        if (list[i] == XA_UTF8_STRING || list[i] == XA_STRING)
            found_text = 1;
        if (list[i] == PNG)
            found_png = 1;
        if (list[i] == JPEG)
            found_jpeg = 1;
    }

    XFree(list);

    if (found_png || found_jpeg)
        return CLIP_IMAGE;

    if (found_text)
        return CLIP_TEXT;

    return CLIP_OTHER;
}

// -----------------------------------------------------------
// Get TEXT
// -----------------------------------------------------------
char* getClipboardTextX11() {
    init_x11();
    if (!dpy) return NULL;

    XConvertSelection(dpy, XA_CLIPBOARD, XA_UTF8_STRING,
                      XA_UTF8_STRING, win, CurrentTime);
    XFlush(dpy);

    XEvent ev;
    XNextEvent(dpy, &ev);

    if (ev.type != SelectionNotify || ev.xselection.property == None)
        return NULL;

    Atom type;
    int format;
    unsigned long len, left;
    unsigned char *data = NULL;

    XGetWindowProperty(
        dpy, win, XA_UTF8_STRING, 0, ~0, False,
        AnyPropertyType, &type, &format, &len, &left, &data
    );

    if (!data || len == 0) {
        if (data) XFree(data);
        return NULL;
    }

    char *out = strndup((char*)data, len);
    XFree(data);
    return out;
}

// -----------------------------------------------------------
// Get IMAGE (PNG or JPEG)
// -----------------------------------------------------------
unsigned char* getClipboardImageX11(int *out_len) {
    init_x11();
    if (!dpy) return NULL;

    *out_len = 0;

    Atom sel = XA_CLIPBOARD;
    Atom targets[] = { PNG, JPEG };
    const int ntargets = sizeof(targets) / sizeof(targets[0]);

    for (int i = 0; i < ntargets; i++) {
        Atom target = targets[i];

        XConvertSelection(dpy, sel, target, target, win, CurrentTime);
        XFlush(dpy);

        XEvent ev;
        XNextEvent(dpy, &ev);

        if (ev.type != SelectionNotify || ev.xselection.property == None)
            continue;

        Atom type;
        int format;
        unsigned long len, left;
        unsigned char *data = NULL;

        XGetWindowProperty(
            dpy, win, target, 0, ~0, False,
            AnyPropertyType, &type, &format, &len, &left, &data
        );

        if (!data || len == 0) {
            if (data) XFree(data);
            continue;
        }

        unsigned char *copy = malloc(len);
        memcpy(copy, data, len);
        XFree(data);

        *out_len = (int)len;
        return copy;
    }

    return NULL; // No image formats available
}
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type ClipboardType int

const (
	ClipNone  ClipboardType = 0
	ClipText  ClipboardType = 1
	ClipImage ClipboardType = 2
	ClipOther ClipboardType = 3
)

func X11ClipboardChanged() bool {
	return C.hasClipboardChangedX11() != 0
}

func X11ClipboardType() ClipboardType {
	return ClipboardType(C.getClipboardTypeX11())
}

func X11ClipboardText() (string, bool) {
	ptr := C.getClipboardTextX11()
	if ptr == nil {
		return "", false
	}
	defer C.free(unsafe.Pointer(ptr))
	return C.GoString(ptr), true
}

func X11ClipboardImage() ([]byte, bool) {
	var outLen C.int
	ptr := C.getClipboardImageX11(&outLen)
	if ptr == nil || outLen == 0 {
		return nil, false
	}
	defer C.free(unsafe.Pointer(ptr))
	return C.GoBytes(unsafe.Pointer(ptr), outLen), true
}

// ----------------------------------

func X11GetClipboardText() string {
	cstr := C.getClipboardTextX11()
	if cstr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

func RunX11Listner() {
	for {

		if X11ClipboardChanged() {
			fmt.Println("cliboard changed!")
			fmt.Println(X11ClipboardType())
		}
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
