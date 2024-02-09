// TODO
// Get date added to json
// Get help menu updated :'((
// Decide on CLI vs GUI
// Multi-module structure format

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
}

func listKeyMap() *keyMap {
	return &keyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("⏎ return", "Copy to clipboard"),
		),
		Backspace: key.NewBinding(
			key.WithKeys("backspace", "delete"),
			key.WithHelp("⌫ backspace", "Delete from clipboard"),
		),
	}
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keys.Backspace, m.keys.Enter}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keys.Backspace, m.keys.Enter}, // First column
		// Second column
	}
}

// MAIN MODEL
type Model struct {
	list list.Model
	help help.Model
	keys *keyMap
	err  error
}

func New() *Model {
	return &Model{
		keys: listKeyMap(),
		help: help.New(),
	}
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

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			err := clipboard.WriteAll(itemName)
			if err != nil {
				panic(err)
			}
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + itemName))
			return m, statusCmd

		case "backspace", "delete":
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

	return m.list.View() //+ m.help.FullHelpView(m.FullHelp())
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
