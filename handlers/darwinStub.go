// handlers/darwinStub.go
//go:build !darwin || !cgo

/* This file is a stub used for CI tests. The cgo lib has introduced
complexity with cross-platform builds, as ubuntu-latest cannot be used for testing
objective-c(++) deps, and ssing macos-latest would fail when testing XWayland deps.
Any cgo files will be omitted from the linter until a better approach is found.

When setting CGO_ENABLED="0" in `go-test.yml`, any files that import "C" are ignored, so we need a stub
to expose the global functions.
*/

package handlers

import "errors"

var darwinErrString = "macOS-only feature"
var errUnsupported = errors.New(darwinErrString)

func RunDarwinListener(_ string, _ bool) error { return errUnsupported }
func DarwinCopyText(_ string)                  {}
func DarwinPaste()                             {}
func GetClipboardText() string                 { return darwinErrString }
func HasClipboardChanged() bool                { return false }
