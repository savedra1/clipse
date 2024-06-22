/*
	TODO
	- fix for items not showing selected after selection was set
	- fix for selected item styles not updating during filter view
	- implement final behaviour for multi-selection:
		- directional unselect
		- selection of next item
	- update help menu
	- implement logic to copy all selected split by a configurable custom string val
*/

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
			//m.resetSelected() - need to implement
			break
		}
		switch msg.String() {
		case "p":
			m.togglePinUpdate()
		case "shift+down", "J":
			m.toggleSelected()
			m.list.CursorDown()
		case "shift+up", "K":
			m.toggleSelected()
			m.list.CursorUp()
		case "S":
			m.toggleSelected()
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
		item.description = fmt.Sprintf("Date copied: %s %s", item.timeStamp, styledPin(m.theme))
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
	item.selected = !item.selected
	item = updateSelectionStyle(item, m.theme)
	m.list.SetItem(index, item)
}

/*func (m *model) resetSelected() {
/*
	Need make selected items reset to dimmed when filtering
*/
//for i := 0; i < len(m.list.Items()); i++ {
//	item, ok := m.list.SelectedItem().(item)
//	if !ok {
//		continue
//	}
//	item.selected = false
//	m.list.SetItem(i, item)
//}
//m.list.Filter("", []string{""})
//}

/*func (m *model) addSelected() {



	index := m.list.Index()
	item, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}
	desc := item.descriptionBase
	if item.pinned {
		desc = desc + " " + styledPin(m.theme)
	}

	item.title = lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.SelectedTitle)).Render(item.titleBase)
	item.description = lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.SelectedTitle)).Render(desc)

	m.list.SetItem(index, item)
}*/
