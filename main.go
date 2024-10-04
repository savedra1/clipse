package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/app"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/handlers"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

var (
	version     = "v1.1.0"
	help        = flag.Bool("help", false, "Show help message.")
	v           = flag.Bool("v", false, "Show app version.")
	add         = flag.Bool("a", false, "Add the following arg to the clipboard history.")
	copyInput   = flag.Bool("c", false, "Copy the input to your systems clipboard.")
	paste       = flag.Bool("p", false, "Prints the current clipboard content.")
	listen      = flag.Bool("listen", false, "Start background process for monitoring clipboard activity on wayland/x11/macOS.")
	listenShell = flag.Bool("listen-shell", false, "Starts a clipboard monitor process in the current shell.")
	kill        = flag.Bool("kill", false, "Kill any existing background processes.")
	clear       = flag.Bool("clear", false, "Remove all contents from the clipboard history except for pinned items.")
	clearAll    = flag.Bool("clear-all", false, "Remove all contents the clipboard history including pinned items.")
	clearImages = flag.Bool("clear-images", false, "Removes all images from the clipboard history including pinned images.")
	clearText   = flag.Bool("clear-text", false, "Removes all text from the clipboard history including pinned text entries.")
	forceClose  = flag.Bool("fc", false, "Forces the terminal session to quick by taking the $PPID var as an arg. EG `clipse -fc $PPID`")
	wlStore     = flag.Bool("wl-store", false, "Store data from the stdin directly using the wl-clipboard API.")
)

func main() {
	flag.Parse()
	logPath, displayServer, imgEnabled, err := config.Init()
	utils.HandleError(err)
	utils.SetUpLogger(logPath)

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
		handleAdd()

	case *paste:
		handlePaste()

	case *copyInput:
		handleCopy()

	case *listen:
		handleListen(displayServer)

	case *listenShell:
		handleListenShell(displayServer, imgEnabled)

	case *kill:
		handleKill()

	case *clear, *clearAll, *clearImages, *clearText:
		handleClear()

	case *forceClose:
		handleForceClose()

	case *wlStore:
		handlers.StoreWLData()

	default:
		fmt.Printf("Command not recognized. See %s --help for usage instructions.", os.Args[0])
	}
}

func launchTUI() {
	shell.KillExistingFG()
	_, err := tea.NewProgram(app.NewModel()).Run()
	utils.HandleError(err)
}

func handleAdd() {
	var input string
	switch {
	case len(os.Args) < 3:
		input = utils.GetStdin()
	default:
		input = os.Args[2]
	}
	utils.HandleError(config.AddClipboardItem(input, "null"))
}

func handleListen(displayServer string) {
	if err := shell.KillExisting(); err != nil {
		fmt.Printf("ERROR: failed to kill existing listener process: %s", err)
		utils.LogERROR(fmt.Sprintf("failed to kill existing listener process: %s", err))
	}
	shell.RunNohupListener(displayServer)
}

func handleListenShell(displayServer string, imgEnabled bool) {
	utils.HandleError(handlers.RunListener(displayServer, imgEnabled))
}

func handleKill() {
	shell.KillAll(os.Args[0])
}

func handleClear() {
	if err := clipboard.WriteAll(""); err != nil {
		utils.LogERROR(fmt.Sprintf("failed to reset clipboard buffer value: %s", err))
	}

	var clearType string

	switch {
	case *clearImages:
		clearType = "images"
	case *clearAll:
		clearType = "all"
	case *clearText:
		clearType = "text"
	default:
		clearType = "default"
	}

	utils.HandleError(config.ClearHistory(clearType))
}

func handleCopy() {
	var input string
	switch {
	case len(os.Args) < 3:
		input = utils.GetStdin()
	default:
		input = os.Args[2]
	}
	if input != "" {
		fmt.Println(input)
		utils.HandleError(clipboard.WriteAll(input))
	}
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
	}

	if len(os.Args) > 3 {
		fmt.Printf("Too many args. Usage: %s' -fc $PPID'", os.Args[0])
		return
	}

	if !utils.IsInt(os.Args[2]) {
		fmt.Printf("Invalid PPID supplied: %s\nPPID must be integer. use var `$PPID` as the arg.", os.Args[2])
		return
	}

	launchTUI()
}
