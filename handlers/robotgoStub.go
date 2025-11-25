//handlers/robotgoStub.go
//go:build wayland || ci

/* Ignore robotgo import when building on Wayland; uinput
to be used instead. */

package handlers

func RobotPaste(_ string) {}
