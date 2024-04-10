package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/savedra1/clipse/app"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/handlers"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"

	"github.com/atotto/clipboard"
)

var (
	version     = "v0.0.7"
	help        = flag.Bool("help", false, "Show help message.")
	v           = flag.Bool("v", false, "Show app version.")
	add         = flag.Bool("a", false, "Add the following arg to the clipboard history.")
	copy        = flag.Bool("c", false, "Copy the input to your systems clipboard.")
	paste       = flag.Bool("p", false, "Prints the current clipboard content.")
	listen      = flag.Bool("listen", false, "Start background process for monitoring clipboard activity.")
	listenShell = flag.Bool("listen-shell", false, "Starts a clipboard monitor process in the current shell.")
	kill        = flag.Bool("kill", false, "Kill any existing background processes.")
	clear       = flag.Bool("clear", false, "Remove all contents from the clipboard's history.")
	forceClose  = flag.Bool("fc", false, "Forces the terminal session to quick by taking the $PPID var as an arg. EG `clipse -fc $PPID`")
)

func main() {
	flag.Parse()
	historyFilePath, clipseDir, imgDir, displayServer, imgEnabled, err := config.Init()
	utils.HandleError(err)

	switch {

	case flag.NFlag() == 0:
		if len(os.Args) > 2 {
			fmt.Println("Too many args provided. See usage:")
			flag.PrintDefaults()
			return
		}
		launchTUI()

	case flag.NFlag() > 1:
		fmt.Printf("Too many flags provided. Use %s --help for more info.", os.Args[0])

	case *help:
		flag.PrintDefaults()

	case *v:
		fmt.Println(os.Args[0], version)

	case *add:
		handleAdd(historyFilePath)

	case *copy:
		handleCopy()

	case *paste:
		handlePaste()

	case *listen:
		handleListen()

	case *listenShell:
		handleListenShell(historyFilePath, clipseDir, displayServer, imgEnabled)

	case *kill:
		handleKill()

	case *clear:
		handleClear(historyFilePath, imgDir)

	case *forceClose:
		handleForceClose()

	default:
		fmt.Printf("Command not recognized. See %s --help for usage instructions.", os.Args[0])
	}
}

func launchTUI() {
	_ = shell.KillExistingFG() // err ignored to mitigate panic when no existinmg clipse ps
	_, err := tea.NewProgram(app.NewModel()).Run()
	utils.HandleError(err)
}

func handleAdd(historyFilePath string) {
	var input string
	if len(os.Args) < 3 {
		input = utils.GetStdin()
	} else {
		input = os.Args[2]
	}

	err := config.AddClipboardItem(historyFilePath, input, "null")
	utils.HandleError(err)
}

func handleListen() {
	shell.KillExisting()
	shell.RunNohupListener() // hardcoded as const
}

func handleListenShell(historyFilePath, clipseDir, displayServer string, imgEnabled bool) {
	err := handlers.RunListener(historyFilePath, clipseDir, displayServer, imgEnabled)
	utils.HandleError(err)
}

func handleKill() {
	shell.KillAll(os.Args[0])
}

func handleClear(historyFilePath, imgDir string) {
	clipboard.WriteAll("")
	err := config.ClearHistory(historyFilePath, imgDir)
	utils.HandleError(err)
}

func handleCopy() {
	var input string
	if len(os.Args) < 3 {
		input = utils.GetStdin()
	} else {
		input = os.Args[2]
	}
	err := clipboard.WriteAll(input)
	utils.HandleError(err)
}

func handlePaste() {
	currentItem, err := clipboard.ReadAll()
	utils.HandleError(err)
	if currentItem != "" {
		fmt.Println(currentItem)
	}
}

func handleForceClose() {
	if len(os.Args) < 3 {
		fmt.Printf("No PPID provided. Usage: %s' -fc $PPID'", os.Args[0])
		return
	} else if len(os.Args) > 3 {
		fmt.Printf("Too many args. Usage: %s' -fc $PPID'", os.Args[0])
		return
	}

	if !utils.IsInt(os.Args[2]) {
		fmt.Printf("Invalid PPID supplied: %s\nPPID must be integer. use var `$PPID` as the arg.", os.Args[2])
		return
	}

	launchTUI()
}
