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
	the Model state.
*/

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.confirmationList.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.filter) && m.list.ShowHelp() {
			m.list.Help.ShowAll = false // change default back to short help to keep in sync
			m.list.SetShowHelp(false)
			m.updatePaginator()
		}

		if m.list.SettingFilter() && key.Matches(msg, m.keys.yankFilter) {
			filterMatches := m.filterMatches()
			if len(filterMatches) >= 1 {
				if err := clipboard.WriteAll(strings.Join(filterMatches, "\n")); err == nil {
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
		timestamp := i.TimeStamp()

		switch {

		case key.Matches(msg, m.keys.choose):
			if m.showConfirmation && m.confirmationList.Index() == 0 { // No
				m.itemCache = []SelectedItem{}
				m.showConfirmation = false
				break

			} else if m.showConfirmation && m.confirmationList.Index() == 1 { // Yes
				m.showConfirmation = false
				currentContent, _ := clipboard.ReadAll()
				timeStamps := []string{}
				for _, item := range m.itemCache {
					if item.Value == currentContent {
						if err := clipboard.WriteAll(""); err != nil {
							utils.LogERROR(fmt.Sprintf("ERROR: could not delete all items from history: %s", err))

						}
					}
					timeStamps = append(timeStamps, item.TimeStamp)
					m.removeCachedItem(item.TimeStamp)
				}

				statusMsg := "Deleted: "
				if len(m.itemCache) == 1 {
					statusMsg += m.itemCache[0].Value
				} else {
					statusMsg += "*selected items*"
				}

				if err := config.DeleteItems(timeStamps); err != nil {
					utils.LogERROR(fmt.Sprintf("ERROR: could not delete all items from history: %s", err))
				}

				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
				)
				m.itemCache = []SelectedItem{}

				if len(m.list.Items()) == 0 {
					m.keys.remove.SetEnabled(false)
					m.list.SetShowStatusBar(false)
				}
				break
			}

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
					utils.HandleError(clipboard.WriteAll(fullValue))
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
			var pinnedItemSelected bool

			m.itemCache = append(
				m.itemCache,
				SelectedItem{
					Index:     m.list.Index(),
					TimeStamp: timestamp,
					Value:     i.titleFull,
					Pinned:    i.pinned,
				},
			)

			if i.pinned {
				pinnedItemSelected = true
			}

			for _, selectedItem := range selectedItems {
				if selectedItem.Pinned {
					pinnedItemSelected = true
				}
				m.itemCache = append(
					m.itemCache,
					selectedItem,
				)
			}

			if pinnedItemSelected {
				m.showConfirmation = true
				break
			}

			currentIndex := m.list.Index()
			currentContent, _ := clipboard.ReadAll()
			statusMsg := "Deleted: "

			if len(selectedItems) >= 1 {
				for _, item := range selectedItems {
					if item.Value == currentContent {
						if err := clipboard.WriteAll(""); err != nil {
							utils.LogERROR(fmt.Sprintf("ERROR: failed to reset clipboard buffer value: %s", err))
						}
					}
				}
				timeStamps := []string{}
				m.list.RemoveItem(currentIndex)
				m.removeMultiSelected()
				for _, item := range selectedItems {
					timeStamps = append(timeStamps, item.TimeStamp)
				}

				timeStamps = append(timeStamps, timestamp)
				statusMsg += "*selected items*"
				if err := config.DeleteItems(timeStamps); err != nil {
					utils.LogERROR(fmt.Sprintf("ERROR: failed to delete all items from history file: %s", err))
				}
			} else {
				m.list.RemoveItem(currentIndex)
				utils.HandleError(config.DeleteItems([]string{timestamp}))
				statusMsg += title
			}

			if len(m.list.Items()) == 0 {
				m.keys.remove.SetEnabled(false)
				m.list.SetShowStatusBar(false)
			}

			m.itemCache = []SelectedItem{}
			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
			)

		case key.Matches(msg, m.keys.togglePin):
			if len(m.list.Items()) == 0 {
				m.keys.togglePin.SetEnabled(false)
			}
			isPinned, err := config.TogglePinClipboardItem(timestamp)
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

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)
	m.confirmationList, cmd = m.confirmationList.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

/*
	HELPER FUNCS
*/

func (m *Model) togglePinUpdate() {
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

func (m *Model) updatePaginator() {
	pagStyle := lipgloss.NewStyle().MarginBottom(1).MarginLeft(2)
	if m.list.ShowHelp() {
		pagStyle = lipgloss.NewStyle().MarginBottom(0).MarginLeft(2)
	}
	m.list.Styles.PaginationStyle = pagStyle
}

func (m *Model) toggleSelectedSingle() {
	m.prevDirection = ""
	index := m.list.Index()
	item, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}
	item.selected = !item.selected
	m.list.SetItem(index, item)
}

func (m *Model) toggleSelected(direction string) {
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

func (m *Model) selectedItems() []SelectedItem {
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
					Index:     index,
					TimeStamp: item.TimeStamp(),
					Value:     item.titleFull,
					Pinned:    item.pinned,
				},
			)
		}

	}
	return selectedItems
}

func (m *Model) removeMultiSelected() {
	items := m.list.Items()
	for i := len(items) - 1; i >= 0; i-- {
		if item, ok := items[i].(item); ok && item.selected {
			m.list.RemoveItem(i)
		}
	}
}

func (m *Model) resetSelected() {
	items := m.list.Items()
	for i := len(items) - 1; i >= 0; i-- {
		if item, ok := items[i].(item); ok && item.selected {
			item.selected = false
			m.list.SetItem(i, item)
		}
	}
}

func (m *Model) filterMatches() []string {
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

func (m *Model) removeCachedItem(ts string) {
	items := m.list.Items()
	for i := len(items) - 1; i >= 0; i-- {
		if item, ok := items[i].(item); ok && item.timeStamp == ts {
			m.list.RemoveItem(i)
		}
	}
}
