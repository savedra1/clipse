package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	//time.Sleep(100 * time.Second)

	help := flag.Bool("help", false, "Show help message.")
	version := flag.Bool("version", false, "Show app version.")
	add := flag.Bool("a", false, "Add the following arg to the clipboard history.")
	listen := flag.Bool("listen", false, "Start background process for monitoring clipboard activity.")
	listenShell := flag.Bool("listen-shell", false, "Starts a clipboard monitor process in the current shell.")
	kill := flag.Bool("kill", false, "Kill any existing background processes.")
	clear := flag.Bool("clear", false, "Remove all contents from the clipboard's history.")

	flag.Parse()

	// explicit path for config file is tested before program can continue
	fullPath, err := checkConfig()
	handleError(err)

	if flag.NFlag() == 0 {
		if len(os.Args) > 1 {
			_, err := strconv.Atoi(os.Args[1]) // check for valid PPID by attempting conversion to an int
			// above line causes canic so cannot catch this error effictively
			if err != nil {
				fmt.Printf("Invalid PPID supplied: %s\nPPID must be integer. use var `$PPID`", os.Args[2])
				return
			}
		}
		_, err := tea.NewProgram(newModel()).Run()
		handleError(err)
		return
	} else if flag.NFlag() > 1 {
		fmt.Printf("Too many flags provided. Use %s --help for more info.", os.Args[0])
		return
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	if *version {
		fmt.Println(os.Args[0], "1.00")
		return
	}

	if *add {
		if len(os.Args) < 3 {
			fmt.Printf("Nothing to add. %s -a requires a following arg. Use --help of more info.", os.Args[0])
			return
		}
		err = addClipboardItem(fullPath, os.Args[2], "")
		handleError(err)
		fmt.Printf("Added %s to clipboard", os.Args[2])
		return
	}

	if *listen {
		err = killExisting()
		handleError(err)
		runNohupListener(listenCmd) // hardcoded as const
		return
	}

	if *listenShell {
		err = runListener(fullPath)
		handleError(err)
		return
	}

	if *kill {
		killAll(os.Args[0])
		handleError(err)
		return
	}

	if *clear {
		err = setBaseConfig(fullPath)
		handleError(err)
		fmt.Println("Removed clipboard contents from system.")
		return
	}

	fmt.Printf("Command not recognised. See %s --help for usage instructions.", os.Args[0])

}
