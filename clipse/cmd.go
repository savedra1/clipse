package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	ps "github.com/mitchellh/go-ps"
)

/*
CMD funcs
*/

func killExisting() error {
	/*
		Kills any existing clipse processes but keeps current ps live
	*/
	currentPS := syscall.Getpid()
	//fmt.Println("current:", currentPS)
	psList, err := ps.Processes()
	if err != nil {
		return err
	}

	for _, p := range psList {
		if strings.Contains(os.Args[0], p.Executable()) {
			//fmt.Println("Process:", p.Pid())
			if p.Pid() != currentPS {
				killProcess(strconv.Itoa(p.Pid()))
			}
		}
	}
	return nil
}

func killExistingFG() error {
	//currentPS := syscall.Getpid()
	//fmt.Println("current:", currentPS)
	psList, err := ps.Processes()
	if err != nil {
		return err
	}

	for _, p := range psList {
		if strings.Contains(os.Args[0], p.Executable()) {
			fmt.Println("Process:", p.Executable())
			stingPID := strconv.Itoa(p.Pid())
			cmd := exec.Command("ps", "-p", stingPID, "-o", "args=")
			output, _ := cmd.Output()

			fmt.Println(output)
			//if p.Pid() != currentPS {
			//	killProcess(strconv.Itoa(p.Pid()))
			//}
		}
	}
	return nil

}

/* Not currently used
func clearShellOutput() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Start() // Not essential to wait for this process to complete
}
*/

func killAll(bin string) {
	cmd := exec.Command("pkill", "-f", bin)
	err := cmd.Run() // Wait for this to finish before executing
	if err != nil {
		fmt.Printf("Failed to kill all existing processes for %s.", bin)
		return
	}
	//clearShellOutput()
}

func runNohupListener(cmdArg string) {
	//c := fmt.Sprintf("nohup %s %s >/dev/null 2>&1 &", os.Args[0], cmdArg)
	cmd := exec.Command("nohup", os.Args[0], cmdArg, ">/dev/null", "2>&1", "&")
	//cmd := exec.Command("zsh", "-c", c)
	err := cmd.Start()
	handleError(err)
	//clearShellOutput()
}

func killProcess(ppid string) {
	cmd := exec.Command("kill", ppid)
	cmd.Run()
}

func imagesEnabled(displayServer string) bool {
	var cmd *exec.Cmd
	switch displayServer {
	case "wayland":
		cmd = exec.Command("sh", "-c", "wl-copy -v")
	case "x11", "darwin":
		cmd = exec.Command("sh", "-c", "xclip -v")
	default:
		return false
	}
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func copyImage(imagePath, displayServer string) error {
	var cmd string
	if displayServer == "wayland" {
		cmd = fmt.Sprintf("wl-copy -t image/png < %s", imagePath)
	} else {
		cmd = fmt.Sprintf("xclip -selection clipboard -t image/png -i %s", imagePath)
	}
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func saveImage(imagePath, displayServer string) error {
	var cmd string
	if displayServer == "wayland" {
		cmd = fmt.Sprintf("wl-paste -t image/png > %s", imagePath)
	} else {
		cmd = fmt.Sprintf("xclip -selection clipboard -t image/png -o > %s", imagePath)
	}

	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func deleteImage(imagePath string) error {
	cmd := fmt.Sprintf("rm %s", imagePath)
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		return err
	}
	return nil

}
