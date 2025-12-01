package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

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
			LogERROR(fmt.Sprintf("failed to delete file %s | %s", file.Name(), err))
		}
	}
	return nil
}
