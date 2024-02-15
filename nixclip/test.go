package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	ps "github.com/mitchellh/go-ps"
)

/* This is where the base level configuration for the
   bubbletea CLI app is defined.

   Base level config includes:
   - Color scheme of text
   - Font
   - Key bindings
   - List sructure
   - Help menu
   - Defualt actions
*/

var (
	// base styling config using lipgloss
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#434C5E")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	// (Each Row in clipboard view)
	title       string
	titleFull   string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) TitleFull() string   { return i.titleFull }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	// default keybind definitions
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{

		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	// model pulls all relevant elems together for rendering
	list         list.Model      // list items
	keys         *listKeyMap     // keybindings
	delegateKeys *delegateKeyMap // custom key bindings
}

func newModel() model {
	// new model needs raising to render additional custom keys
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Make initial list of items
	clipboardItems := getjsonData()
	var entryItems []list.Item
	for _, entry := range clipboardItems {
		shortenedVal := shorten(entry.Value)
		item := item{
			title:       shortenedVal,
			titleFull:   entry.Value,
			description: "Copied to clipboard: " + entry.Recorded,
		}
		entryItems = append(entryItems, item)
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	clipboardList := list.New(entryItems, delegate, 0, 0)
	clipboardList.Title = "Clipboard History"
	clipboardList.Styles.Title = titleStyle
	clipboardList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         clipboardList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd { // initialise app
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	/* this is where the base logic is held for what action to take from
	   the predefined key bindings
	*/
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string { // Render app in terminal using client libs
	return appStyle.Render(m.list.View())
}

/* This is were we define additional config to add to our
base-level bubbletea app. Here including keybinds only.
*/

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

/*
Additional short/full help entries. This satisfies the help.KeyMap interface and

	is entirely optional.
*/
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

// final config map for new keys
func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "copy"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	/* This is where the additional keybinding actions are defined:
	   - enter/reurn: copies selected item to the clipboard and adds a status message
	   - backspace/delete: removes item from list view and json file

	*/
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var fullValue string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
			fullValue = i.TitleFull()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				err := clipboard.WriteAll(fullValue)
				if err != nil {
					panic(err)
				}
				cmd := exec.Command("kill", os.Args[2])
				err = cmd.Run()
				handleError(err)
				return m.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				fullPath := getFullPath()
				err := deleteJsonItem(fullPath, fullValue)
				handleError(err)
				return m.NewStatusMessage(statusMessageStyle("Deleted: " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

/* File contains logic for parseing the cilpboard data and
   general config.
   - fileName defined in constants.go
   - dirName defined in constants.go
*/

// ClipboardItem struct for individual clipboardHistor array item
type ClipboardItem struct {
	// EG: {"value": "copied_string", "recorded": "datetime"}
	Value    string `json:"value"`
	Recorded string `json:"recorded"`
}

type ClipboardHistory struct {
	ClipboardHistory []ClipboardItem `json:"clipboardHistory"`
}

// saveDataToFile saves data to a JSON file
func saveDataToFile(fullPath string, data ClipboardHistory) error {
	/* Triggered from the system copy action:
	   Adds the copied string to the clipboard_history.json file.
	*/
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func getjsonData() []ClipboardItem {
	/* returns the clipboardHistory array from the
	   clipboard_history.json file
	*/
	fullPath := getFullPath()
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("error opening file:", err)
		file.Close()
	}

	var data ClipboardHistory
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	// Extract clipboard history items
	return data.ClipboardHistory

}

func deleteJsonItem(fullPath, item string) error {
	/* Accessed by bubbletea method on backspace keybinding:
	   Deletes selected item from json file.
	*/
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var data ClipboardHistory
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var updatedClipboardHistory []ClipboardItem
	for _, entry := range data.ClipboardHistory {
		if entry.Value != item {
			updatedClipboardHistory = append(updatedClipboardHistory, entry)
		}
	}

	updatedData := ClipboardHistory{
		ClipboardHistory: updatedClipboardHistory,
	}
	updatedJSON, err := json.Marshal(updatedData)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Write the updated JSON back to the file
	if err := os.WriteFile(fullPath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func createConfigDir(configDir string) error {
	/* Used to create the ~/.config/clipboard_manager dir
	   in relative path.
	*/
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		os.Exit(1)
	}
	return nil
}

func createHistoryFile(fullPath string) error {
	/* Used to create the clipboard_history.json file
	   in relative path.
	*/
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = setBaseConfig(fullPath)
	if err != nil {
		return err
	}
	return nil
}

func getFullPath() string {
	/* Returns full path string for clipboard file.
	   useful when needing to be accessed form a
	   bubbletea method.
	*/
	currentUser, err := user.Current()
	handleError(err)
	// Construct the path to the config directory
	configDir := filepath.Join(currentUser.HomeDir, ".config", configDirName)
	fullPath := filepath.Join(configDir, fileName)
	return fullPath
}

func checkConfig() (string, error) {
	/* Ensure $HOME/.config/clipboard_manager/clipboard_history.json
	   exists and create the path if not. Full path returned as string
	   when successful
	*/
	currentUser, err := user.Current()
	handleError(err)

	// Construct the path to the config directory
	configDir := filepath.Join(currentUser.HomeDir, ".config", configDirName)
	fullPath := filepath.Join(configDir, fileName)

	_, err = os.Stat(fullPath) // File already exist?
	if os.IsNotExist(err) {

		_, err = os.Stat(configDir) // Config dir at least exists?
		if os.IsNotExist(err) {
			err = createConfigDir(configDir)
			if err != nil {
				fmt.Println("Failed to create config dir. Please create:", configDir)
				os.Exit(1)
			}
		}

		_, err = os.Stat(fullPath) // Attempts creation of full path now that relative path exists on system
		if os.IsNotExist(err) {
			err = createHistoryFile(fullPath)
			if err != nil {
				fmt.Println("Failed to create", fullPath)
				os.Exit(1)
			}

		}

	} else if err != nil {
		fmt.Println("Unable to check if config file exists. Please update binary permisisons.")
		os.Exit(1)
	}
	return fullPath, nil
}

func setBaseConfig(fullPath string) error {
	/*
		 Sets clipboard_history.json file to:
			{
				"clipboardHistory": []
			}
	*/
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0644) // Permisisons specified for file to allow write
	if err != nil {
		return err
	}
	defer file.Close()

	// Truncate the file to zero length
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// Rewind the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	baseConfig := ClipboardHistory{
		ClipboardHistory: []ClipboardItem{},
	}

	// Encode initial history to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(baseConfig); err != nil {
		return err
	}

	return nil
}

/*
Global vars stored in separate module.
Any new additions to be added here.
*/

const (
	fileName      = "clipboard_history.json"
	configDirName = "clipboard_manager"
	pollInterval  = 100 * time.Millisecond / 10
	maxLen        = 50
)

/* runListener is essentially a while loop to be created as a system background process on boot.
   can be stopped at any time with:
   	clipboard kill
   	pkill -f clipboard
   	killall clipboard
*/

func runListener(fullPath string) error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Load existing data from file, if any
	var data ClipboardHistory

	go func() { // go routine necessary to acheive desired CTRL+C behavior
		for {

			// Get the current clipboard content
			text, err := clipboard.ReadAll()
			handleError(err)

			// If clipboard content is not empty and not already in the list, add it
			if text != "" && !contains(data.ClipboardHistory, text) {
				// If the length exceeds 50, remove the oldest item
				if len(data.ClipboardHistory) >= 50 {
					lastIndex := len(data.ClipboardHistory) - 1
					data.ClipboardHistory = data.ClipboardHistory[:lastIndex] // Remove the oldest item
				}

				// yyyy-mm-dd hh-mm-s.msmsms Time format
				timeNow := strings.Split(time.Now().UTC().String(), "+0000")[0]

				// {"value": "copied_strig", "recorded": "2024-01-02 12:34:78743687"}
				item := ClipboardItem{Value: text, Recorded: timeNow}

				data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)

				// Save updated data to JSON file
				err = saveDataToFile(fullPath, data)
				handleError(err)

			}

			time.Sleep(pollInterval) // pollInterval defined in constants.go

		}

	}()
	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	return nil
}

func main() {
	// definitions for cmd flags and args
	listen := "listen"
	clear := "clear"
	listenStart := "listen-start-background-process" // obscure arg to prevent accidental usage
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
	bin := "test.go"
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case listen:
			// Kill existing clipboard processes
			shellCmd := exec.Command("pkill", "-f", bin)
			err = shellCmd.Run()
			handleError(err)

			shellCmd = exec.Command("nohup", "go", "run", bin, listenStart, ">/dev/null", "2>&1", "&")
			err = shellCmd.Run()
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
			shellCmd := exec.Command("pkill", "-f", bin)
			shellCmd.Run()
			fmt.Println("Stopped all clipboard listener processes. Use `clipboard listen` to resume.")
			return

		case "test":
			//htop()
			c := exec.Command("kill", "3407899")
			out, err := c.Output()
			handleError(err)
			fmt.Println(out)

			return

		case "htop":
			htop()
			return

		case "open":

			// Open bubbletea app in terminal session
			if _, err := tea.NewProgram(newModel()).Run(); err != nil {
				fmt.Println("Error opening clipboard:\n", err)
				os.Exit(1)
			}

		default:
			// Arg not recognised
			fmt.Println("Arg not recognised. Try `clipboard --help` for more details.")
			return
		}
	}

}

/* General purpose functions to be used by other modules
 */

func htop() {
	list, err := ps.Processes()
	if err != nil {
		panic(err)
	}
	for _, p := range list {
		if strings.Contains(p.Executable(), "zsh") { //&& p.PPid() != 1 {
			fmt.Printf("- Process %s with PID %d and PPID %d\n", p.Executable(), p.Pid(), p.PPid())
		}

	}
}

// Avoids repeat code by handling errors in a uniform way
func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Contains checks if a string exists in a slice of strings
func contains(slice []ClipboardItem, str string) bool {
	for _, item := range slice {
		if item.Value == str {
			return true
		}
	}
	return false
}

// Shortens string val to show in list view
func shorten(s string) string {
	if len(s) <= maxLen { // maxLen defined in constants.go
		return strings.ReplaceAll(s, "\n", " ")
	}
	return strings.ReplaceAll(s[:maxLen-3], "\n", " ") + "..."
}
