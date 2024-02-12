package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ##################### LISTENER SECTION ####################### //
// Data struct for storing clipboard strings
type Data struct {
	ClipboardHistory []ClipboardItem `json:"clipboardHistory"`
}

// ClipboardItem struct for individual clipboard history item
type ClipboardItem struct {
	Value    string `json:"value"`
	Recorded string `json:"recorded"`
}

func runListener() error {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Load existing data from file, if any
	var data Data
	err := loadDataFromFile("../history.json", &data)
	if err != nil {
		fmt.Println("Error loading data from file:", err)
	}

	// Start a goroutine to continuously monitor clipboard changes
	go func() {
		for {
			// Get the current clipboard content
			text, err := clipboard.ReadAll()
			if err != nil {
				fmt.Println("Error reading clipboard:", err)
			}

			// If clipboard content is not empty and not already in the list, add it
			if text != "" && !contains(data.ClipboardHistory, text) {
				// If the length exceeds 50, remove the oldest item
				if len(data.ClipboardHistory) >= 50 {
					lastIndex := len(data.ClipboardHistory) - 1
					data.ClipboardHistory = data.ClipboardHistory[:lastIndex] // Remove the oldest item
				}
				timeNow := time.Now().UTC().String()
				item := ClipboardItem{Value: text, Recorded: timeNow}
				data.ClipboardHistory = append([]ClipboardItem{item}, data.ClipboardHistory...)
				//fmt.Println("Added to clipboard history:", text)

				// Save data to file
				err := saveDataToFile("../history.json", data)
				if err != nil {
					fmt.Println("Error saving data to file:", err)
				}
			}

			// Check for updates every 0.1 second
			duration := 100 * time.Millisecond / 10
			time.Sleep(duration)
		}
	}()

	fmt.Println("Clipboard history listener running... Press Ctrl+C to exit.")

	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	fmt.Println("Exiting...")
	return nil
}

// contains checks if a string exists in a slice of strings
func contains(slice []ClipboardItem, str string) bool {
	for _, item := range slice {
		if item.Value == str {
			return true
		}
	}
	return false
}

// loadDataFromFile loads data from a JSON file
func loadDataFromFile(filename string, data *Data) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// saveDataToFile saves data to a JSON file
func saveDataToFile(filename string, data Data) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

// ##################### LISTENER SECTION ####################### //

var (
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
	title       string
	titleFull   string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) TitleFull() string   { return i.titleFull }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
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
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func newModel() model {
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

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func shorten(s string) string {
	maxLen := 50 // Define your max length here
	if len(s) <= maxLen {
		return strings.ReplaceAll(s, "\n", " ")
	}
	return strings.ReplaceAll(s[:maxLen-3], "\n", " ") + "..."
}

// NEW ITEM DELEGATE SECTION
func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
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
				return m.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				err := deleteJsonItem(fullValue)
				if err != nil {
					os.Exit(1)
				}
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

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

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

type jsonFile struct {
}

type ClipboardEntry struct {
	Value    string `json:"value"`
	Recorded string `json:"recorded"`
}

type ClipboardHistory struct {
	ClipboardHistory []ClipboardEntry `json:"clipboardHistory"`
}

func getjsonData() []ClipboardEntry {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error opening file:", err)
		file.Close()
	}

	// Decode JSON from the file
	var data ClipboardHistory
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	// Extract clipboard history items
	clipboardHistory := data.ClipboardHistory

	return clipboardHistory

}

func deleteJsonItem(item string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var data ClipboardHistory
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var updatedClipboardHistory []ClipboardEntry
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
	if err := os.WriteFile(fileName, updatedJSON, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func checkConfig() error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// File does not exist, create it with default values

		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		defer file.Close()

		baseConfig := ClipboardHistory{
			ClipboardHistory: []ClipboardEntry{},
		}

		err = setBaseConfig(baseConfig)
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(baseConfig); err != nil {
			return err
		}

	} else if err != nil {
		return err
	}
	return nil
}

func setBaseConfig(baseConfig ClipboardHistory) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
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

	// Encode initial history to JSON and write to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(baseConfig); err != nil {
		return err
	}

	return nil
}

const (
	fileName = "../clipboard_history.json"
)

func main() {
	// cmd flags and args
	listen := "listen"
	clear := "clear"
	listenStart := "listen-start-background-process-0088" // obscure string to prevent accidental usage
	help := flag.Bool("help", false, "Show help message")

	err := checkConfig()
	if err != nil {
		fmt.Println("No clipboard_history.json file found in path. Failed to create:", err)
		return
	}

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
			cmd := exec.Command("pkill", "main.go")
			cmd.Run() // Kill existing clipboard processes
			cmd = exec.Command("nohup", "go", "run", "main.go", listenStart, ">/dev/null", "2>&1", "&")
			//cmd = exec.Command("nohup", "sh", "-c", "go run "+os.Args[0]+" --start >/dev/null 2>&1 &")

			if err := cmd.Start(); err != nil {
				fmt.Println("Error starting clipboard listener:", err)
				os.Exit(1)
			}
			fmt.Println("Starting clipboard listener...")
			return
		case clear:
			fmt.Println("Place holder")
			return
		case listenStart:
			err := runListener()
			if err != nil {
				fmt.Println(err)
			}
			return
		default:
			fmt.Println("Arg not recognised. Try `clipboard --help` for more details.")
			return
		}
	}

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error opening clipboard:", err)
		os.Exit(1)
	}
}
