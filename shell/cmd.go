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
			// fmt.Println("Process:", p.Pid())
			if p.Pid() != currentPS {
				KillProcess(strconv.Itoa(p.Pid()))
			}
		}
	}
	return nil
}

func KillExistingFG() error {
	/*
		Only kill other clipboard GUI windows to prevent
		file conflicts.
	*/

	currentPS := strconv.Itoa(syscall.Getpid())
	// fmt.Println("current:", currentPS)
	cmd := exec.Command("sh", "-c", "pgrep -a clipse")
	output, err := cmd.Output()
	if err != nil || output == nil { // allows local usage when no clipse ps
		return fmt.Errorf("no clipse processes are running")
	}
	/*
		EG Output returns as:
		156842 ./clipse --listen-shell >/dev/null 2>&1 &
		310228 ./clipse
	*/

	psList := strings.Split(string(output), "\n")
	for _, ps := range psList {
		if !strings.Contains(ps, currentPS) && !strings.Contains(ps, listenCmd) {
			KillProcess(strings.Split(ps, " ")[0])
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

func KillAll(bin string) {
	cmd := exec.Command("pkill", "-f", bin)
	err := cmd.Run() // Wait for this to finish before executing
	if err != nil {
		fmt.Printf("Failed to kill all existing processes for %s.", bin)
		return
	}
	// clearShellOutput()
}

func RunNohupListener() {
	cmd := exec.Command("nohup", os.Args[0], listenCmd, ">/dev/null", "2>&1", "&")
	utils.HandleError(cmd.Start())
}

func KillProcess(ppid string) {
	cmd := exec.Command("kill", ppid)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to kill process: %s", err)
	}
}
