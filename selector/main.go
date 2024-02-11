package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	maxLen := 30 // Define your max length here
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

type ClipboardEntry struct {
	Value    string `json:"value"`
	Recorded string `json:"recorded"`
}

type ClipboardHistory struct {
	ClipboardHistory []ClipboardEntry `json:"clipboardHistory"`
}

func getjsonData() []ClipboardEntry {
	file, err := os.Open("../history.json")
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
	filePath := "../history.json"
	fileContent, err := os.ReadFile(filePath)
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
	if err := os.WriteFile(filePath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func main() {
	// cmd args
	help := "--help"
	listen := "--listen"
	clear := "--clear"

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case help:
			fmt.Println("Place holder")
			os.Exit(0)
		case listen:
			fmt.Println("Place holder")
			os.Exit(0)
		case clear:
			fmt.Println("Place holder")
			os.Exit(0)

		default:
			fmt.Println("Arg not recognised. Try `clipboard --help` for more details.")
			os.Exit(0)
		}
	}

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error opening clipboard:", err)
		os.Exit(1)
	}
}
