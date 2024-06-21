package app

import (
	"fmt"
	"strings"

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
		case "shift+down":
			m.toggleSelected()
			m.list.CursorDown()
		case "shift+up":
			m.toggleSelected()
			m.list.CursorUp()

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

// This updates the TUI when an item is pinned/unpinned
func (m *model) togglePinUpdate() {
	index := m.list.Index()
	item, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}
	item.description = fmt.Sprintf("Date copied: %s", item.timeStamp)
	if !item.pinned {
		item.description = fmt.Sprintf("Date copied: %s %s", item.timeStamp, styledPin())
	}

	item.pinned = !item.pinned
	m.list.SetItem(index, item)
	if m.list.IsFiltered() {
		m.list.ResetFilter() // move selected pinned item to front
	}
}

func (m *model) updatePaginator() {
	pagStyle := lipgloss.NewStyle().MarginBottom(1).MarginLeft(2)
	if m.list.ShowHelp() {
		pagStyle = lipgloss.NewStyle().MarginBottom(0).MarginLeft(2)
	}
	m.list.Styles.PaginationStyle = pagStyle
}

func (m *model) toggleSelected() {
	index := m.list.Index()
	item, ok := m.list.SelectedItem().(item)

	if !ok {
		return
	}
	selectedChar := "➤➤➤ "
	item.selected = !item.selected
	if !item.selected {
		item.title = strings.Replace(item.title, selectedChar, "", 1)
	} else if string(item.title[0]) != selectedChar {
		item.title = selectedChar + item.title
	}
	m.list.SetItem(index, item)

}
