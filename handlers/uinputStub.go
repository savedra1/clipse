// handlers/uinputStub.go
//go:build !wayland

/* Ignore the uinput code when not built on Wayland; robotgo to be used
instead.
*/

package handlers

func UinputPaste(_ string) {}
