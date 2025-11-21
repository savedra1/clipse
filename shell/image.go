package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

var copyImgCmds = map[string]string{
	"darwin":  darwinCopyImgCmd,
	"wayland": wlCopyImgCmd,
}

var pasteImgCmds = map[string]string{
	"wayland": wlPasteImgCmd,
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
