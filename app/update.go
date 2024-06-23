/*
	TODO
	- implement final behaviour for multi-selection:
		- directional unselect
		- selection of next item
	- implement logic to copy all selected split by a configurable custom string val
*/

package app

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
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
		item, ok := m.list.SelectedItem().(item)
		if !ok {
			return nil, nil
		}
		title := item.Title()
		fullValue := item.TitleFull()
		fp := item.FilePath()
		desc := item.TimeStamp()
		switch {
		case key.Matches(msg, m.keys.choose):
			if fp != "null" {
				ds := config.DisplayServer() // eg "wayland"
				err := shell.CopyImage(fp, ds)
				utils.HandleError(err)
			} else {
				err := clipboard.WriteAll(fullValue)
				utils.HandleError(err)
			}

			if len(os.Args) > 2 {
				if utils.IsInt(os.Args[2]) {
					shell.KillProcess(os.Args[2])
				}
			} else if len(os.Args) > 1 {
				if os.Args[1] == "keep" {
					cmds = append(
						cmds,
						m.list.NewStatusMessage(statusMessageStyle("Copied to clipboard: "+title)),
					)
				}
			} else {
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.remove):
			index := m.list.Index()
			m.list.RemoveItem(index)
			if len(m.list.Items()) == 0 {
				m.keys.remove.SetEnabled(false)
				m.list.SetShowStatusBar(false)
			}
			go func() { // stop cached clipboard item repopulating
				currentContent, _ := clipboard.ReadAll()
				if currentContent == fullValue {
					clipboard.WriteAll("")
				}
				err := config.DeleteJsonItem(desc)
				utils.HandleError(err)
			}()
			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle("Deleted: "+title)),
			)

		case key.Matches(msg, m.keys.togglePin):
			if len(m.list.Items()) == 0 {
				m.keys.togglePin.SetEnabled(false)
			}
			// update pinned status in history file
			isPinned, err := config.TogglePinClipboardItem(desc)
			utils.HandleError(err)
			m.togglePinUpdate()

			if isPinned {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("UnPinned: "+title)),
				)
			} else {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("Pinned: "+title)),
				)
			}
		case key.Matches(msg, m.keys.togglePinned):
			if len(m.list.Items()) == 0 {
				m.keys.togglePinned.SetEnabled(false)
			}
			m.togglePinned = !m.togglePinned
			if m.togglePinned {
				m.list.Title = "Pinned " + clipboardTitle
			} else {
				m.list.Title = clipboardTitle
			}
			clipboardItems := config.GetHistory()
			filteredItems := filterItems(clipboardItems, m.togglePinned, m.theme)

			if len(filteredItems) == 0 {
				m.list.Title = clipboardTitle
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("No pinned items")),
				)
			} else {
				for i := len(m.list.Items()) - 1; i >= 0; i-- { // clear all items
					m.list.RemoveItem(i)
				}
				for _, i := range filteredItems { // adds all required items
					m.list.InsertItem(len(m.list.Items()), i)
				}
			}
		case key.Matches(msg, m.keys.selectDown):
			m.toggleSelected("down")

		case key.Matches(msg, m.keys.selectUp):
			m.toggleSelected("up")

		case key.Matches(msg, m.keys.selectSingle):
			m.toggleSelectedSingle()

		case key.Matches(msg, m.keys.more):
			// swap custom help menu for default list.Model help view when expanding
			// the menu. doing this because the custom help menu causing rendering
			// conflits with the list view
			m.list.SetShowHelp(!m.list.ShowHelp())
			m.updatePaginator()

		case key.Matches(msg, m.keys.up),
			key.Matches(msg, m.keys.down),
			key.Matches(msg, m.keys.nextPage),
			key.Matches(msg, m.keys.prevPage),
			key.Matches(msg, m.keys.home),
			key.Matches(msg, m.keys.end):
			m.prevDirection = ""
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

func (m *model) toggleSelectedSingle() {
	m.prevDirection = ""
	index := m.list.Index()
	item, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}
	item.selected = !item.selected
	m.list.SetItem(index, item)
}

func (m *model) toggleSelected(direction string) {
	if m.prevDirection == "" {
		m.prevDirection = direction
	}

	index := m.list.Index()
	item, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}

	if item.selected {
		item.selected = false
	} else if m.prevDirection == direction && !item.selected {
		item.selected = true
	}

	m.list.SetItem(index, item)

	switch direction {
	case "down":
		m.list.CursorDown()
	case "up":
		m.list.CursorUp()
	}

}
