package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

func ImagesEnabled(displayServer string) bool {
	var cmd *exec.Cmd
	switch displayServer {
	case "wayland":
		cmd = exec.Command("sh", "-c", wlVersionCmd)
	case "x11", "darwin":
		cmd = exec.Command("sh", "-c", xVersionCmd)
	default:
		return false
	}
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func CopyImage(imagePath, displayServer string) error {
	cmd := fmt.Sprintf("%s %s", xCopyImgCmd, imagePath)
	if displayServer == "wayland" {
		cmd = fmt.Sprintf("%s %s", wlCopyImgCmd, imagePath)
	}
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		return err
	}
	return nil
}

func SaveImage(imagePath, displayServer string) error {
	cmd := fmt.Sprintf("%s %s", xPasteImgCmd, imagePath)
	if displayServer == "wayland" {
		cmd = fmt.Sprintf("%s %s", wlPasteImgCmd, imagePath)
	}
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
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
