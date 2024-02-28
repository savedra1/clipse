package main

import (
	"clipse/app"
	"clipse/config"
	"clipse/handlers"
	"clipse/shell"
	"clipse/utils"

	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	version     = "v1.0.0"
	help        = flag.Bool("help", false, "Show help message.")
	v           = flag.Bool("v", false, "Show app version.")
	add         = flag.Bool("a", false, "Add the following arg to the clipboard history.")
	listen      = flag.Bool("listen", false, "Start background process for monitoring clipboard activity.")
	listenShell = flag.Bool("listen-shell", false, "Starts a clipboard monitor process in the current shell.")
	kill        = flag.Bool("kill", false, "Kill any existing background processes.")
	clear       = flag.Bool("clear", false, "Remove all contents from the clipboard's history.")
)

func main() {
	//time.Sleep(10000 * time.Second)
	flag.Parse()
	historyFilePath, clipseDir, displayServer, imgEnabled, err := config.Init()
	utils.HandleError(err)

	if flag.NFlag() == 0 {
		shell.KillExistingFG()
		if len(os.Args) > 1 {
			_, err := strconv.Atoi(os.Args[1]) // check for valid PPID by attempting conversion to an int
			// above line causes canic so cannot catch this error effictively
			if err != nil {
				fmt.Printf("Invalid PPID supplied: %s\nPPID must be integer. use var `$PPID`", os.Args[1])
				return
			}
		}
		_, err := tea.NewProgram(app.NewModel()).Run()
		utils.HandleError(err)
		return

	} else if flag.NFlag() > 1 {
		fmt.Printf("Too many flags provided. Use %s --help for more info.", os.Args[0])
		return
	}

	if *help {
		flag.PrintDefaults()
		return
	}

	if *v {
		fmt.Println(os.Args[0], version)
		return
	}

	if *add {
		var input string
		if len(os.Args) < 3 {
			input = utils.GetStdin()
		} else {
			input = os.Args[2]
		}

		err := config.AddClipboardItem(historyFilePath, input, "null")
		utils.HandleError(err)
		fmt.Println("added the following val to clipboard:", input)

		return
	}

	if *listen {
		shell.KillExisting()
		shell.RunNohupListener() // hardcoded as const
		return
	}

	if *listenShell {
		err = handlers.RunListener(historyFilePath, clipseDir, displayServer, imgEnabled)
		utils.HandleError(err)
		return
	}

	if *kill {
		shell.KillAll(os.Args[0])
		utils.HandleError(err)
		return
	}

	if *clear {
		clipboard.WriteAll("")
		err = config.ClearHistory(historyFilePath)
		utils.HandleError(err)
		fmt.Println("Removed clipboard contents from system.")
		return
	}

	fmt.Printf("Command not recognised. See %s --help for usage instructions.", os.Args[0])

}
