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

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	//"github.com/charmbracelet/lipgloss"
	"github.com/atotto/clipboard"
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

// MAIN MODEL
type Model struct {
	list list.Model
	err  error
}

func New() *Model {
	return &Model{}
}

func (m *Model) initList(width, height int) { // window size
	var entryItems []list.Item
	clipboardHistory := getjsonData()
	for _, entry := range clipboardHistory {
		entryItems = append(entryItems, Entry{title: entry, description: "---"})
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

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			itemName := m.list.SelectedItem().FilterValue()
			err := clipboard.WriteAll(itemName)
			if err != nil {
				panic(err)
			}
			os.Exit(0)
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
	file, err := os.Open("history.json")
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

	// Print clipboard history items
	// fmt.Println("Clipboard History:")
	// for _, item := range clipboardHistory {
	// 	fmt.Println(item)
	// }
	return clipboardHistory

}

func main() {
	m := New()
	p := tea.NewProgram(m)
	err, _ := p.Run()
	if err != nil {
		os.Exit(1)

	}
}
