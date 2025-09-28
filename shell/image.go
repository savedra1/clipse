package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

var imgIsEnabledCmd = map[string]string{
	"darwin":  darwinVersionCmd,
	"wayland": wlVersionCmd,
	"x11":     xVersionCmd,
}

var copyImgCmds = map[string]string{
	"darwin":  darwinCopyImgCmd,
	"wayland": wlCopyImgCmd,
	"x11":     xCopyImgCmd,
}

var pasteImgCmds = map[string]string{
	"darwin":  darwinPasteImgCmd,
	"wayland": wlPasteImgCmd,
	"x11":     xPasteImgCmd,
}

func ImagesEnabled(displayServer string) bool {
	cmd, ok := imgIsEnabledCmd[displayServer]
	if !ok {
		utils.LogWARN(fmt.Sprintf("unknown display server: %s", displayServer))
		return false
	}
	execCmd := exec.Command("sh", "-c", cmd)
	if err := execCmd.Run(); err != nil {
		utils.LogERROR(fmt.Sprintf("%s system is missing image dependency", displayServer))
		return false
	}
	return true
}

func CopyImage(imagePath, displayServer string) error {
	cmd, ok := copyImgCmds[displayServer]
	if !ok {
		return fmt.Errorf("unknown display server: %s; could not copy image", displayServer)
	}
	cmdFull := fmt.Sprintf(cmd, imagePath)
	if err := exec.Command("sh", "-c", cmdFull).Run(); err != nil {
		return err
	}
	return nil
}

func SaveImage(imagePath, displayServer string) error {
	// imagePath string cannot contain space chars unless wrapped
	cmd, ok := pasteImgCmds[displayServer]
	if !ok {
		return fmt.Errorf("unknown display server: %s; could not save image", displayServer)
	}
	cmdFull := fmt.Sprintf(cmd, imagePath)
	if err := exec.Command("sh", "-c", cmdFull).Run(); err != nil {
		return err
	}
	return nil
}

func DeleteImage(imagePath string) error {
	if err := os.Remove(imagePath); err != nil {
		return err
	}
	return nil
}

func DeleteAllImages(imgDir string) error {
	files, err := os.ReadDir(imgDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := os.Remove(filepath.Join(imgDir, file.Name())); err != nil {
			utils.LogERROR(fmt.Sprintf("failed to delete file %s | %s", file.Name(), err))
		}
	}
	return nil
}

func DarwinImageDataPresent() (bool, []byte) {
	output, err := exec.Command("sh", "-c", darwinImgCheckCmd).Output()
	if err != nil {
		return false, nil
	}
	return true, output
}
