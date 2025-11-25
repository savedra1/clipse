package shell

import (
	"os/exec"
)

func GetWLClipBoard() (string, error) {
	cmd := exec.Command(wlPasteHandler)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func UpdateWLClipboard(s string) error {
	cmd := exec.Command(wlCopyHandler, s)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func WLDependencyCheck() error {
	cmd := exec.Command("which", wlCopyHandler)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
