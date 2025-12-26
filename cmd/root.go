package cmd

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/app"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/display"
	"github.com/savedra1/clipse/handlers"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

var (
	version       = "v1.1.5"
	help          = flag.Bool("help", false, "Show help message.")
	v             = flag.Bool("v", false, "Show app version.")
	add           = flag.Bool("a", false, "Add the following arg to the clipboard history.")
	copyInput     = flag.Bool("c", false, "Copy the input to your systems clipboard.")
	paste         = flag.Bool("p", false, "Prints the current clipboard content.")
	listen        = flag.Bool("listen", false, "Start background process for monitoring clipboard activity on wayland/x11/macOS.")
	listenShell   = flag.Bool("listen-shell", false, "Starts a clipboard monitor process in the current shell.")
	listenDarwin  = flag.Bool("listen-darwin", false, "Starts a clipboard monitor process in the current shell for Darwin systems.")
	listenX11     = flag.Bool("listen-x11", false, "Starts a clipboard monitor process in the current shell for X11 systems.")
	kill          = flag.Bool("kill", false, "Kill any existing background processes.")
	clearUnpinned = flag.Bool("clear", false, "Remove all contents from the clipboard history except for pinned items.")
	clearAll      = flag.Bool("clear-all", false, "Remove all contents the clipboard history including pinned items.")
	clearImages   = flag.Bool("clear-images", false, "Removes all images from the clipboard history including pinned images.")
	clearText     = flag.Bool("clear-text", false, "Removes all text from the clipboard history including pinned text entries.")
	wlStore       = flag.Bool("wl-store", false, "Store data from the stdin directly using the wl-clipboard API.")
	realTime      = flag.Bool("enable-real-time", false, "Enable real time updates to the TUI")
	outputAll     = flag.String("output-all", "", "Print clipboard text content to stdout, each entry separated by a newline, possible values: (raw, unescaped)")
	autoPaste     = flag.Bool("auto-paste", false, "send key event to paste")
	pause         = flag.String("pause", "", "Pause clipboard monitoring for a specified duration. Example: `clipse -pause 5m` pauses for 5 minutes.")
)

func Main() int {
	flag.Parse()

	utils.HandleError(config.Init())
	utils.SetUpLogger(config.ClipseConfig.LogFilePath)

	app.ClipseTheme = config.GetTheme()

	switch {

	case flag.NFlag() == 0:
		if len(os.Args) > 2 {
			fmt.Println("Too many args provided. See usage:")
			flag.PrintDefaults()
			return 1
		}
		launchTUI()

	case flag.NFlag() > 1:
		fmt.Printf("Too many flags provided. Use %s --help for more info.", os.Args[0])
		return 1

	case *help:
		flag.PrintDefaults()

	case *v:
		fmt.Println(os.Args[0], version)

	case *add:
		handleAdd()

	case *paste:
		display.DisplayServer.Paste()

	case *copyInput:
		handleCopy()

	case *listen:
		handleListen()

	case *listenShell:
		display.DisplayServer.RunListener()

	case *listenDarwin:
		handlers.RunDarwinListener()

	case *listenX11:
		handlers.RunX11Listener()

	case *kill:
		handleKill()

	case *clearUnpinned, *clearAll, *clearImages, *clearText:
		handleClear()

	case *wlStore:
		handlers.StoreWLData()

	case *realTime:
		launchTUI()

	case *outputAll != "":
		handleOutputAll(*outputAll)

	case *pause != "":
		handlePause(*pause)

	case *autoPaste:
		handleAutoPaste()

	default:
		fmt.Printf("Command not recognized. See %s --help for usage instructions.", os.Args[0])
		return 1
	}
	return 0
}

func handlePause(s string) {
	ok, err := shell.IsListenerRunning()
	if err != nil {
		fmt.Printf("Error checking for active clipboard monitoring process: %s\n", err)
		os.Exit(1)
	}
	if !ok {
		fmt.Println("No active clipboard monitoring process found. Cannot pause.")
		os.Exit(1)
	}
	usageMsg := fmt.Sprintf("Usage: %s -pause <duration>\nWhere duration is in seconds, minutes, or hours. Example: %s -pause 5m pauses for 5 minutes.", os.Args[0], os.Args[0])
	if s == "0" {
		fmt.Println("Invalid duration. Use a positive duration to pause clipboard monitoring.")
		fmt.Println(usageMsg)
		os.Exit(1)
	}
	duration, err := utils.ParseDuration(s)
	if err != nil {
		fmt.Printf("Invalid duration format: %s\n", err)
		fmt.Println(usageMsg)
		os.Exit(1)
	}
	if err := shell.KillExisting(); err != nil {
		fmt.Printf("ERROR: failed to kill existing listener process: %s", err)
		utils.LogERROR(fmt.Sprintf("failed to kill existing listener process: %s", err))
		os.Exit(1)
	}
	fmt.Printf("Pausing clipboard monitoring for %s...\n", duration)
	shell.RunListenerAfterDelay(duration)
}

func launchTUI() {
	shell.KillExistingFG()
	newModel := app.NewModel()
	p := tea.NewProgram(
		newModel,
		tea.WithMouseCellMotion(),
		tea.WithMouseAllMotion(),
	)
	if *realTime {
		go newModel.ListenRealTime(p)
	}
	finalModel, err := p.Run()
	utils.HandleError(err)

	if m, ok := finalModel.(app.Model); ok {
		if m.ExitCode != 0 {
			os.Exit(m.ExitCode)
		}
		if config.ClipseConfig.AutoPaste.Enabled {
			shell.RunAutoPaste()
		}
	}
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

func handleListen() {
	if err := shell.KillExisting(); err != nil {
		fmt.Printf("ERROR: failed to kill existing listener process: %s", err)
		utils.LogERROR(fmt.Sprintf("failed to kill existing listener process: %s", err))
		os.Exit(1)
	}
	display.DisplayServer.RunDetachedListener()
}

func handleKill() {
	shell.KillAll(os.Args[0])
}

func handleClear() {
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
	display.DisplayServer.CopyText(input)
}

func handleOutputAll(format string) {
	items := config.TextItems()

	switch format {
	case "raw":
		for _, v := range items {
			fmt.Printf("%q\n", v.Value)
		}
	case "unescaped":
		for _, v := range items {
			fmt.Println(v.Value)
		}
	default:
		fmt.Printf("Invalid argument to -output-all\nSee %s --help for usage", os.Args[0])
		os.Exit(1)
	}
}

func handleAutoPaste() {
	time.Sleep(time.Duration(config.ClipseConfig.AutoPaste.Buffer) * time.Microsecond)
	display.DisplayServer.SendPasteKey(config.ClipseConfig.AutoPaste.Keybind)
}
