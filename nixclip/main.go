package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// definitions for cmd flags and args
	listen := "listen"
	clear := "clear"
	listenStart := "listen-start-background-process-dev-null" // obscure arg to prevent accidental usage
	kill := "kill"

	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	// explicit path for config file is tested before program can continue
	fullPath, err := checkConfig()
	handleError(err)

	if *help {
		standardInfo := "| `clipboard` -> open clipboard history"
		clearInfo := "| `clipboard clear` -> truncate clipboard history"
		listenInfo := "| `clipboard listen` -> starts background process to listen for clipboard events"

		fmt.Printf(
			"Available commands:\n\n%s\n\n%s\n\n%s\n\n",
			standardInfo, clearInfo, listenInfo,
		)
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case listen:
			// Kill existing clipboard processes
			shellCmd := exec.Command("pkill", "-f", "clipboard")
			shellCmd.Run()
			shellCmd = exec.Command("nohup", "./clipboard", listenStart, ">/dev/null", "2>&1")
			err = shellCmd.Start()
			handleError(err)
			return

		case clear:
			// Remove contents of jsonFile.clipboardHistory array
			err = setBaseConfig(fullPath)
			handleError(err)
			fmt.Println("Cleared clipboard contents from system.")
			return

		case listenStart:
			//Hidden arg that starts listener as background process
			err := runListener(fullPath)
			handleError(err)
			return

		case kill:
			// End any existing background listener processes
			shellCmd := exec.Command("pkill", "-f", "clipboard")
			shellCmd.Run()
			fmt.Println("Stopped all clipboard listener processes. Use `clipboard listen` to resume.")
			return

		default:
			// Arg not recognised
			fmt.Println("Arg not recognised. Try `clipboard --help` for more details.")
			return
		}
	}

	// Open bubbletea app in terminal session
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error opening clipboard:\n", err)
		os.Exit(1)
	}
}
