// TODO
// Decide on CLI vs GUI
// Clear clipboard method
// Multi-module structure format
//

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	statusMessageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
		Render
)

// Create format for each history item
type Entry struct {
	title       string
	description string
}

// Function to effect the Item struct directly
func (e Entry) FilterValue() string {
	return e.title
}

func (e Entry) Title() string {
	return e.title
}

func (e Entry) Description() string {
	return e.description
}

type keyMap struct {
	Enter     key.Binding
	Backspace key.Binding
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Help      key.Binding
	Quit      key.Binding
}

var keys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("⏎ return", "Copy to clipboard"),
	),
	Backspace: key.NewBinding(
		key.WithKeys("backspace", "delete"),
		key.WithHelp("⌫ backspace", "Delete from clipboard"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit},                // second column
	}
}

// MAIN MODEL
type Model struct {
	list list.Model
	help help.Model
	keys keyMap
	err  error
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) { // window size
	var entryItems []list.Item
	clipboardHistory := getjsonData()
	for _, entry := range clipboardHistory {
		entryItems = append(entryItems, Entry{title: entry, description: "Added to clipboard 22:05|08/02/24"})
	}

	m.list = list.New(entryItems, list.NewDefaultDelegate(), width, height)
	m.list.Title = "Clipboard History"

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
		m.help.Width = msg.Width

	case tea.KeyMsg:
		itemName := m.list.SelectedItem().FilterValue()

		//switch msg.String() {
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Enter):
			err := clipboard.WriteAll(itemName)
			if err != nil {
				panic(err)
			}
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + itemName))
			return m, statusCmd

		case key.Matches(msg, m.keys.Backspace):
			index := m.list.Index()
			m.list.RemoveItem(index)

			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Deleted: " + itemName))
			deleteJsonItem(itemName)
			return m, statusCmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}

type ClipboardData struct {
	ClipboardHistory []string `json:"clipboardHistory"`
}

func getjsonData() []string {
	file, err := os.Open("../history.json")
	if err != nil {
		fmt.Println("error opening file:", err)
		file.Close()
	}

	// Decode JSON from the file
	var data ClipboardData
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

	var data ClipboardData
	if err := json.Unmarshal(fileContent, &data); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	var updatedClipboardHistory []string
	for _, entry := range data.ClipboardHistory {
		if entry != item {
			updatedClipboardHistory = append(updatedClipboardHistory, entry)
		}
	}

	updatedData := ClipboardData{
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

// MENU

func main() {
	m := New()
	p := tea.NewProgram(m)

	err, _ := p.Run()
	if err != nil {
		os.Exit(1)

	}
}
