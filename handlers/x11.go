//go:build linux && !wayland
// +build linux,!wayland

package handlers

/*]
#cgo pkg-config: x11 xfixes
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

// Returns 1 if clipboard changed
// int hasClipboardChangedX11() {
//     init_x11();
//     if (!dpy) return 0;

//     XEvent ev;
//     int changed = 0;

//     while (XPending(dpy)) {
//         XNextEvent(dpy, &ev);
//         if (ev.type == XFixesSelectionNotify) {
//             long serial = ev.xfixesselection.selection_timestamp;
//             if (serial != last_serial) {
//                 last_serial = serial;
//                 changed = 1;
//             }
//         }
//     }
//     return changed;
// }

// // Sets text into clipboard
// void setClipboardTextX11(const char *text) {
//     init_x11();
//     if (!dpy) return;
//     if (!text) return;

//     XSetSelectionOwner(dpy, XA_CLIPBOARD, win, CurrentTime);
//     XFlush(dpy);

//     // respond to selection requests
//     XEvent ev;
//     for (;;) {
//         XNextEvent(dpy, &ev);
//         if (ev.type == SelectionRequest) {
//             XSelectionRequestEvent *req = &ev.xselectionrequest;

//             XEvent reply;
//             memset(&reply, 0, sizeof(reply));
//             reply.xselection.type = SelectionNotify;
//             reply.xselection.display = req->display;
//             reply.xselection.requestor = req->requestor;
//             reply.xselection.selection = req->selection;
//             reply.xselection.target = req->target;
//             reply.xselection.property = None;
//             reply.xselection.time = req->time;

//             if (req->target == XA_UTF8_STRING) {
//                 XChangeProperty(dpy, req->requestor, req->property,
//                                 XA_UTF8_STRING, 8, PropModeReplace,
//                                 (unsigned char*)text, strlen(text));
//                 reply.xselection.property = req->property;
//             }

//             XSendEvent(dpy, req->requestor, True, 0, &reply);
//             XFlush(dpy);
//             break;
//         }
//     }
// }
*/
import "C"
import "unsafe"

func GetClipboardText() string {
	cstr := C.getClipboardTextX11()
	if cstr == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

// func XHasClipboardChanged() bool {
// 	return C.hasClipboardChangedX11() == 1
// }

// func XSetClipboardText(s string) {
// 	cs := C.CString(s)
// 	defer C.free(unsafe.Pointer(cs))
// 	C.setClipboardTextX11(cs)
// }
