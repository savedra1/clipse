//handlers/robotgo.go
//go:build !wayland && !ci

package handlers

import (
	"fmt"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/savedra1/clipse/utils"
)

func RobotPaste(keybind string) {
	parts := strings.Split(keybind, "+")

	if len(parts) == 0 {
		utils.LogERROR(fmt.Sprintf("invalid keybind: %s", keybind))
		return
	}

	// Last element is the key, everything else is modifiers
	key := parts[len(parts)-1]
	mods := parts[:len(parts)-1]

	utils.HandleError(robotgo.KeyTap(key, mods))
}
