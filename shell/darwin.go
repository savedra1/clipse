package shell

import (
	"fmt"
	"os/exec"

	"github.com/savedra1/clipse/utils"
)

func RunDarwinListener() {
	cmd := exec.Command(ExeName, darwinListenCmd)
	runDetachedCmd(cmd)
}

func DarwinCopyImage(filePath string) {
	cmdFull := fmt.Sprintf(darwinCopyImgCmd, filePath)
	if err := exec.Command("sh", "-c", cmdFull).Run(); err != nil {
		utils.LogERROR(fmt.Sprintf("failed to copy image: %s", err))
	}
}

func DarwinActiveWindowTitle() string {
	output := execOutput("osascript", "-e", darwinGetWindowCmd)
	if output == "" {
		utils.LogWARN("Failed to get active window")
	}
	return output
}
