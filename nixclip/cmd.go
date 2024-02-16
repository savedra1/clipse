package main

import (
	"fmt"
	"os"
	"os/exec"
	//ps "github.com/mitchellh/go-ps"
)

/* CMD funcs
 */

/* NOT IN USE
func getPPIDs(process string) string {
	list, err := ps.Processes()
	if err != nil {
		panic(err)
	}

	results := ""

	for _, p := range list {
		if strings.Contains(p.Executable(), process) { //&& p.PPid() != 1 {
			results += fmt.Sprintf(
				"- Process %s with PID %d and PPID %d\n", p.Executable(), p.Pid(), p.PPid(),
			)
		}
	}
	return results
}*/

func clearShellOutput() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Start() // Not essential to wait for this process to complete
}

func killExistingProcess(bin string) {
	cmd := exec.Command("pkill", "-f", bin)
	err := cmd.Run() // Wait for this to finish before executing
	if err != nil {
		fmt.Printf("Failed to kill existing background processes for %s", bin)
		return
	}
	clearShellOutput()
}

func runNohupListener(cmdArg string) {
	//c := fmt.Sprintf("nohup %s %s >/dev/null 2>&1 &", os.Args[0], cmdArg)
	cmd := exec.Command("nohup", os.Args[0], cmdArg, ">/dev/null", "2>&1", "&")
	//cmd := exec.Command("zsh", "-c", c)
	err := cmd.Start()
	handleError(err)
	clearShellOutput()
}

func closeShell(ppid string) {
	cmd := exec.Command("kill", ppid)
	cmd.Run()
}
