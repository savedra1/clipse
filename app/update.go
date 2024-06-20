package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
		switch msg.String() {
		case "p":
			m.togglePinUpdate()
		case "?":
			// swap custom help menu for default list.Model help view when expanding
			// the menu. doing this because the custom help menu causing rendering
			// conflits with the list view
			m.list.SetShowHelp(!m.list.ShowHelp())
			m.updatePaginator()
		}
	}

	// this will also call our delegate's update function
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) updatePaginator() {
	pagStyle := lipgloss.NewStyle().MarginBottom(1).MarginLeft(2)
	if m.list.ShowHelp() {
		pagStyle = lipgloss.NewStyle().MarginBottom(0).MarginLeft(2)
	}
	m.list.Styles.PaginationStyle = pagStyle
}

// This updates the TUI when an item is pinned/unpinned
func (m *model) togglePinUpdate() {
	index := m.list.Index()
	i, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}
	i.description = fmt.Sprintf("Date copied: %s", i.timeStamp)
	if !i.pinned {
		i.description = fmt.Sprintf("Date copied: %s %s", i.timeStamp, pinnedStyle())
	}

	i.pinned = !i.pinned
	m.list.SetItem(index, i)

}
