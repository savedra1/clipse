//handlers/robotgoStub.go
//go:build wayland

/* Ignore robotgo import when building on Wayland; uinput
to be used instead. */

package handlers

func RobotPaste(_ string) {}
