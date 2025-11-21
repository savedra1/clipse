// handlers/x11Stub.go
//go:build !linux || wayland || !cgo

/* This file is a stub used for CI tests. The cgo lib has introduced
complexity with cross-platform builds, as ubuntu-latest cannot be used for testing
c/xorg deps, and using x11 runners would fail when testing XWayland deps.
Any cgo files will be omitted from the linter until a better approach is found.

When setting CGO_ENABLED="0" in `go-test.yml`, any files that import "C" are ignored, so we need a stub
to expose the global functions.
*/

package handlers

import "errors"

var x11ErrString = "X11-only feature"
var errX11Unsupported = errors.New(x11ErrString)

func X11GetClipboardText() string               { return x11ErrString }
func X11ClipboardChanged() bool                 { return false }
func RunX11Listener()                           {}
func GetClipboardImage() ([]byte, error)        { return []byte{}, errX11Unsupported }
func X11Paste()                                 {}
func X11SetClipboardText(_ string) error        { return errX11Unsupported }
func X11SetClipboardImage([]byte, string) error { return errX11Unsupported }
