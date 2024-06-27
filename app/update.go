package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

/*
	The main update function used to handle core TUI logic and update
	the model state.
*/

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if msg.String() == "y" {
			m.itemCach = m.list.Items()
			m.showConfirmation = true
			m.list.SetItems(confirmationItems())
			m.keys.remove.SetEnabled(false)
			m.list.SetFilteringEnabled(false)
			m.list.SetShowStatusBar(false)
			m.list.Title = "Delete pinned item(s)?"
			m.list.ResetSelected()
			//m.setConfirmationScreen()

			return m, tea.Batch(cmds...)

		}
		if key.Matches(msg, m.keys.choose) && m.showConfirmation {
			m.list.SetItems(m.itemCach)
			m.showConfirmation = false
			m.keys.remove.SetEnabled(true)
			m.list.SetFilteringEnabled(true)
			m.list.SetShowStatusBar(true)
			m.list.Title = clipboardTitle
			return m, tea.Batch(cmds...)

		}

		if key.Matches(msg, m.keys.filter) && m.list.ShowHelp() {
			m.list.Help.ShowAll = false // change default back to short help to keep in sync
			m.list.SetShowHelp(false)
			m.updatePaginator()
		}

		if m.list.SettingFilter() && key.Matches(msg, m.keys.yankFilter) {
			filterMatches := m.filterMatches()
			if len(filterMatches) >= 1 {
				err := clipboard.WriteAll(strings.Join(filterMatches, "\n"))
				if err == nil {
					return m, tea.Quit
				}
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("Failed to copy all selected items.")),
				)
			}
			return m, tea.Batch(cmds...)
		}

		// Don't match any of the keys below if we're actively filtering.
		if m.list.SettingFilter() {
			break
		}

		i, ok := m.list.SelectedItem().(item)
		if !ok {

			switch {
			case key.Matches(msg, m.keys.more):
				m.list.SetShowHelp(!m.list.ShowHelp())
				m.updatePaginator()
			}
			break
		}
		title := i.Title()
		fullValue := i.TitleFull()
		fp := i.FilePath()
		desc := i.TimeStamp()

		switch {

		case key.Matches(msg, m.keys.choose):
			selectedItems := m.selectedItems()

			if len(selectedItems) < 1 {
				switch {
				case fp != "null":
					ds := config.DisplayServer() // eg "wayland"
					utils.HandleError(shell.CopyImage(fp, ds))
					return m, tea.Quit

				case len(os.Args) > 2 && utils.IsInt(os.Args[2]):
					shell.KillProcess(os.Args[2])
					return m, tea.Quit

				case len(os.Args) > 1 && os.Args[1] == "keep":
					utils.HandleError(clipboard.WriteAll(fullValue))
					cmds = append(
						cmds,
						m.list.NewStatusMessage(statusMessageStyle("Copied to clipboard: "+title)),
					)
					return m, tea.Batch(cmds...)

				default:
					err := clipboard.WriteAll(fullValue)
					utils.HandleError(err)
					return m, tea.Quit
				}
			}

			yank := ""
			for _, item := range selectedItems {
				if fullValue != item.Value {
					yank += item.Value + "\n"
				}
			}
			yank += fullValue
			switch {

			case len(os.Args) > 2 && utils.IsInt(os.Args[2]):
				utils.HandleError(clipboard.WriteAll(yank))
				shell.KillProcess(os.Args[2])
				return m, tea.Quit

			case len(os.Args) > 1 && os.Args[1] == "keep":
				statusMsg := "Copied to clipboard: *selected items*"
				err := clipboard.WriteAll(yank)
				if err != nil {
					statusMsg = "Could not copy all selected items."
				}
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
				)
				return m, tea.Batch(cmds...)

			default:
				err := clipboard.WriteAll(yank)
				if err == nil {
					return m, tea.Quit
				}
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("Could not copy all selected items.")),
				)

			}

		case key.Matches(msg, m.keys.remove):
			selectedItems := m.selectedItems()
			currentIndex := m.list.Index()

			currentContent, _ := clipboard.ReadAll()

			for _, item := range selectedItems {
				if item.Value == currentContent {
					// clear clipboard to stop deleted content temp repopulating (temp solution)
					clipboard.WriteAll("")
				}
			}

			statusMsg := "Deleted: "

			if len(selectedItems) >= 1 {
				timeStamps := []string{}
				m.list.RemoveItem(currentIndex)
				m.removeMultiSelected()
				for _, item := range selectedItems {
					timeStamps = append(timeStamps, strings.Split(item.Description, "Date copied: ")[1])
				}

				timeStamps = append(timeStamps, desc)
				statusMsg += "*selected items*"
				config.DeleteItems(timeStamps)
			} else {
				m.list.RemoveItem(currentIndex)
				err := config.DeleteItems([]string{desc})
				utils.HandleError(err)
				statusMsg += title
			}

			if len(m.list.Items()) == 0 {
				m.keys.remove.SetEnabled(false)
				m.list.SetShowStatusBar(false)
			}

			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
			)

		case key.Matches(msg, m.keys.togglePin):
			if len(m.list.Items()) == 0 {
				m.keys.togglePin.SetEnabled(false)
			}
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
			m.list.Title = clipboardTitle
			if m.togglePinned {
				m.list.Title = "Pinned " + clipboardTitle
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
				for _, i := range filteredItems { // redraw all required items
					m.list.InsertItem(len(m.list.Items()), i)
				}
			}
		case key.Matches(msg, m.keys.selectDown):
			if m.list.IsFiltered() {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("cannot select with filter applied")),
				)
			} else {
				m.toggleSelected("down")
			}

		case key.Matches(msg, m.keys.selectUp):
			if m.list.IsFiltered() {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("cannot select with filter applied")),
				)
			} else {
				m.toggleSelected("up")
			}

		case key.Matches(msg, m.keys.selectSingle):
			if m.list.IsFiltered() {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("cannot select with filter applied")),
				)
			} else {
				m.toggleSelectedSingle()
			}

		case key.Matches(msg, m.keys.clearSelected), key.Matches(msg, m.keys.filter):
			m.resetSelected()

		case key.Matches(msg, m.keys.yankFilter):
			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle("no filtered items")),
			)

		case key.Matches(msg, m.keys.more):
			// switch to default help for full view (better rendering)
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

	item, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}

	index := m.list.Index()

	if item.selected {
		item.selected = false
	} else if m.prevDirection == direction && !item.selected {
		item.selected = true
	} else {
		m.prevDirection = ""
	}

	m.list.SetItem(index, item)

	switch direction {
	case "down":
		m.list.CursorDown()
	case "up":
		m.list.CursorUp()
	}
}

// used to retrieve selected items with their list view indexes
type SelectedItem struct {
	Index       int
	Description string
	Value       string
}

// index values, descriptions and full values of all selected items
func (m *model) selectedItems() []SelectedItem {
	selectedItems := []SelectedItem{}
	for index, i := range m.list.Items() {
		item, ok := i.(item)
		if !ok {
			continue
		}
		if item.selected {
			selectedItems = append(
				selectedItems,
				SelectedItem{
					Index:       index,
					Description: item.descriptionBase,
					Value:       item.titleFull,
				},
			)
		}
	}
	return selectedItems
}

// iterate over the list items backwards so the indexes are not affected
func (m *model) removeMultiSelected() {
	items := m.list.Items()
	for i := len(items) - 1; i >= 0; i-- {
		if item, ok := items[i].(item); ok && item.selected {
			m.list.RemoveItem(i)
		}
	}
}

// remove selected state from all items
func (m *model) resetSelected() {
	items := m.list.Items()
	for i := len(items) - 1; i >= 0; i-- {
		if item, ok := items[i].(item); ok && item.selected {
			item.selected = false
			m.list.SetItem(i, item)
		}
	}
}

// return a list of all current filter matches
func (m *model) filterMatches() []string {
	filteredItems := []string{}
	for _, i := range m.list.Items() {
		item, ok := i.(item)
		if !ok {
			continue
		}
		if strings.Contains(
			strings.ToLower(item.titleFull),
			strings.ToLower(m.list.FilterValue()),
		) {
			filteredItems = append(filteredItems, item.titleFull)
		}
	}

	return filteredItems
}

func (m model) setConfirmationScreen() {
	m.list.SetItems(confirmationItems())
	m.keys.remove.SetEnabled(false)
	m.list.SetFilteringEnabled(false)
	m.list.SetShowFilter(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowHelp(false)
	m.list.Title = "Delete pinned item(s)?"

}
