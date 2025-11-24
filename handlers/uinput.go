// handlers/uinput.go
//go:build wayland

/* uint lib used to handle automated paste following copy action.
Lib: https://github.com/bendahl/uinput/
Requires access to the /dev/uinput device
E.g. before running: sudo chmod +rwx /dev/uinput
*/

package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
	"gopkg.in/bendahl/uinput.v1"
)

var wlKeyboardDevice = "/dev/uinput"
var wlKeyboardDeviceName = "clipskb"

var uinputKeys = map[string]int{
	"a":          uinput.KeyA,
	"b":          uinput.KeyB,
	"c":          uinput.KeyC,
	"d":          uinput.KeyD,
	"e":          uinput.KeyE,
	"f":          uinput.KeyF,
	"g":          uinput.KeyG,
	"h":          uinput.KeyH,
	"i":          uinput.KeyI,
	"j":          uinput.KeyJ,
	"k":          uinput.KeyK,
	"l":          uinput.KeyL,
	"m":          uinput.KeyM,
	"n":          uinput.KeyN,
	"o":          uinput.KeyO,
	"p":          uinput.KeyP,
	"q":          uinput.KeyQ,
	"r":          uinput.KeyR,
	"s":          uinput.KeyS,
	"t":          uinput.KeyT,
	"u":          uinput.KeyU,
	"v":          uinput.KeyV,
	"w":          uinput.KeyW,
	"x":          uinput.KeyX,
	"y":          uinput.KeyY,
	"z":          uinput.KeyZ,
	"0":          uinput.Key0,
	"1":          uinput.Key1,
	"2":          uinput.Key2,
	"3":          uinput.Key3,
	"4":          uinput.Key4,
	"5":          uinput.Key5,
	"6":          uinput.Key6,
	"7":          uinput.Key7,
	"8":          uinput.Key8,
	"9":          uinput.Key9,
	"minus":      uinput.KeyMinus,
	"equal":      uinput.KeyEqual,
	"leftbrace":  uinput.KeyLeftbrace,
	"rightbrace": uinput.KeyRightbrace,
	"semicolon":  uinput.KeySemicolon,
	"backslash":  uinput.KeyBackslash,
	"comma":      uinput.KeyComma,
	"dot":        uinput.KeyDot,
	"slash":      uinput.KeySlash,
}

var uinputMods = map[string]int{
	"capslock": uinput.KeyCapslock,
	"ctrl":     uinput.KeyLeftctrl,
	"shift":    uinput.KeyLeftshift,
	"alt":      uinput.KeyLeftalt,
	"meta":     uinput.KeyLeftmeta,
	"insert":   uinput.KeyInsert,
}

func uinputPaste(keybind string) error {
	keyboard, err := uinput.CreateKeyboard(wlKeyboardDevice, []byte(wlKeyboardDeviceName))
	if err != nil {
		return fmt.Errorf("failed to create virtual keyboard: %w", err)
	}
	defer keyboard.Close()

	parts := strings.Split(strings.ToLower(keybind), "+")
	if len(parts) == 0 {
		return fmt.Errorf("invalid keybind: %s", keybind)
	}

	mainKey := parts[len(parts)-1]
	modifiers := parts[:len(parts)-1]

	mainKeyCode, ok := uinputKeys[mainKey]
	if !ok {
		// Check if it's a modifier being used as main key (uncommon but possible)
		mainKeyCode, ok = uinputMods[mainKey]
		if !ok {
			return fmt.Errorf("unknown key: %s", mainKey)
		}
	}

	// Press all modifiers
	var modCodes []int
	for _, mod := range modifiers {
		modCode, ok := uinputMods[mod]
		if !ok {
			return fmt.Errorf("unknown modifier: %s", mod)
		}
		if err := keyboard.KeyDown(modCode); err != nil {
			return fmt.Errorf("failed to press modifier %s: %w", mod, err)
		}
		modCodes = append(modCodes, modCode)
	}

	time.Sleep(time.Duration(config.ClipseConfig.AutoPaste.Buffer) * time.Millisecond)

	if err := keyboard.KeyPress(mainKeyCode); err != nil {
		return fmt.Errorf("failed to press key %s: %w", mainKey, err)
	}

	// Release modefiers
	for _, code := range modCodes {
		if err := keyboard.KeyUp(code); err != nil {
			return fmt.Errorf("failed to release modifier: %w", err)
		}
	}

	return nil
}

// Wrapper for compatibility with your existing error handling
func UinputPaste(keybind string) {
	utils.HandleError(uinputPaste(keybind))
}
