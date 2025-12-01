package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	ps "github.com/mitchellh/go-ps"

	"github.com/savedra1/clipse/utils"
)

var ExeName = exeName()

func IsListenerRunning() (bool, error) {
	/*
		Check if clipse is running by checking if the process is in the process list.
	*/
	psList, err := ps.Processes()
	if err != nil {
		return false, fmt.Errorf("failed to get processes: %w", err)
	}

	for _, p := range psList {
		if strings.Contains(os.Args[0], p.Executable()) || strings.Contains(wlPasteHandler, p.Executable()) {
			return true, nil
		}
	}
	return false, nil
}

func KillExisting() error {
	/*
		Kills any existing clipse processes but keeps current ps live
	*/
	currentPS := syscall.Getpid()
	psList, err := ps.Processes()
	if err != nil {
		return err
	}

	for _, p := range psList {
		if strings.Contains(os.Args[0], p.Executable()) || strings.Contains(wlPasteHandler, p.Executable()) {
			if p.Pid() != currentPS {
				KillProcess(strconv.Itoa(p.Pid()))
			}
		}
	}
	return nil
}

func KillExistingFG() {
	/*
		Only kill other clipboard TUI windows to prevent
		file conflicts.
	*/
	currentPid := strconv.Itoa(syscall.Getpid())

	// Get PIDs of process names (as opposed to full commands) containing "clipse"
	cmd := exec.Command("sh", "-c", pgrepCmd)
	pgrepOutput, err := cmd.Output() // returns something like "1234\n5678\n\n"
	if err != nil {
		utils.LogWARN(fmt.Sprintf("failed to get processes | err msg: %s | output: %s", err, pgrepOutput))
		return
	}
	pidList := strings.Fields(string(pgrepOutput))

	for _, pid := range pidList {
		if pid == currentPid {
			continue
		}

		// Get full command for given PID, eg. "./clipse --listen-shell >/dev/null 2>&1 &"
		cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", psCmd, pid))
		psOutput, err := cmd.Output()
		if err != nil {
			utils.LogWARN(fmt.Sprintf("failed to get pid's command | pid: %s | err msg: %s | output: %s", pid, err, psOutput))
			continue
		}
		pidCmd := strings.Split(string(psOutput), "\n")[1] // skip headers (macOS's ps doesn't support --no-headers)
		if strings.Contains(pidCmd, listenShellCmd) ||
			strings.Contains(pidCmd, wlStoreCmd) ||
			strings.Contains(pidCmd, darwinListenCmd) ||
			strings.Contains(pidCmd, x11ListenCmd) {
			continue
		}

		utils.LogINFO(fmt.Sprintf("Killing pid %s, cmd %s", pid, pidCmd))
		KillProcess(pid)
	}
}

func KillAll(bin string) {
	cmd := exec.Command("pkill", "-f", bin)
	err := cmd.Run() // Wait for this to finish before executing
	if err != nil {
		utils.LogERROR(fmt.Sprintf("Failed to kill all existing processes for %s.", bin))
		return
	}
}

func RunListenerAfterDelay(delay *time.Duration) {
	if delay == nil {
		utils.LogERROR("Delay cannot be nil")
		return
	}

	cmd := exec.Command("sleep", strconv.Itoa(int(delay.Seconds())), "&&", ExeName, listenCmd)
	runDetachedCmd(cmd)
}

func RunAutoPaste() {
	cmd := exec.Command(ExeName, "--auto-paste")
	runDetachedCmd(cmd)
}

func exeName() string {
	exe, err := os.Executable()
	utils.HandleError(err)
	return exe
}

func runDetachedCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	cmd.Stdout = devNull
	cmd.Stderr = devNull
	cmd.Stdin = nil

	utils.HandleError(cmd.Start())
}

func KillProcess(ppid string) {
	cmd := exec.Command("kill", ppid)
	if err := cmd.Run(); err != nil {
		utils.LogERROR(fmt.Sprintf("failed to kill process: %s", err))
	}
}

func execOutput(name string, args ...string) string {
	if _, err := exec.LookPath(name); err != nil {
		return ""
	}

	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
