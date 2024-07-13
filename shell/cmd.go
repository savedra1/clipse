package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	ps "github.com/mitchellh/go-ps"

	"github.com/savedra1/clipse/utils"
)

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
		if strings.Contains(os.Args[0], p.Executable()) {
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
	currentPS := strconv.Itoa(syscall.Getpid())
	cmd := exec.Command("sh", "-c", pgrepCmd)
	output, err := cmd.Output()
	/*
		EG Output returns as:
		156842 clipse --listen-shell >/dev/null 2>&1 &
		310228 clipse
	*/
	if err != nil {
		utils.LogWARN(fmt.Sprintf("failed to get processes | err msg: %s | output: %s", err, output))
		return
	}
	if output == nil {
		return // no clipse processes running
	}

	psList := strings.Split(string(output), "\n")
	for _, ps := range psList {
		if strings.Contains(ps, currentPS) || strings.Contains(ps, listenCmd) {
			continue
		}
		if ps != "" {
			KillProcess(strings.Split(ps, " ")[0])
		}
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

func RunNohupListener(displayServer string) {
	switch displayServer {
	case "wayland":
		// run optimized wl-clipboard listener
		utils.HandleError(nohupCmdWL("image").Start())
		utils.HandleError(nohupCmdWL("text").Start())

	default:
		// run default poll listener
		cmd := exec.Command("nohup", os.Args[0], listenCmd, ">/dev/null", "2>&1", "&")
		utils.HandleError(cmd.Start())
	}
}

func nohupCmdWL(dataType string) *exec.Cmd {
	cmd := exec.Command(
		"nohup",
		wlPasteHandler,
		wlTypeSpec,
		dataType,
		wlPasteWatcher,
		os.Args[0],
		wlStoreCmd,
		">/dev/null",
		"2>&1",
		"&",
	)
	return cmd
}

func KillProcess(ppid string) {
	cmd := exec.Command("kill", ppid)
	if err := cmd.Run(); err != nil {
		utils.LogERROR(fmt.Sprintf("failed to kill process: %s", err))
	}
}
