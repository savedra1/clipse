package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/display"
	"github.com/savedra1/clipse/utils"
)

/*
	The main update function used to handle core TUI logic and update
	the Model state.
*/

var KeepEnabled = keepEnabled()

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ReRender:
		clipboardItems := config.GetHistory()
		entryItems := filterItems(clipboardItems, false)
		cmds = append(cmds, m.list.SetItems(entryItems))
		return m, tea.Batch(cmds...)
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.confirmationList.SetSize(msg.Width-h, msg.Height-v)

		headerHeight := lipgloss.Height(m.previewHeaderView())
		footerHeight := lipgloss.Height(m.previewFooterView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.previewReady {
			m.preview = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.previewReady = true
			m.preview.YPosition = headerHeight // '+ 1' needed for high performance rendering only
			break
		}
		m.preview.Width = msg.Width
		m.preview.Height = msg.Height - verticalMarginHeight

	case tea.KeyMsg:
		if key.Matches(msg, m.keys.filter) && m.list.ShowHelp() {
			m.list.Help.ShowAll = false // change default back to short help to keep in sync
			m.list.SetShowHelp(false)
			m.updatePaginator()
		}

		if m.list.SettingFilter() && key.Matches(msg, m.keys.yankFilter) {
			filterMatches := m.filterMatches()
			if len(filterMatches) >= 1 {
				display.DisplayServer.CopyText(strings.Join(filterMatches, "\n"))
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("Failed to copy all selected items.")),
				)
			}
			return m, tea.Batch(cmds...)
		}

		// Don't match any of the keys below if we're actively filtering.
		if m.list.SettingFilter() {
			m.setQuitEnabled(false) // disable main list quit to allow filter cancel
			break
		}

		if m.showPreview {
			switch {
			case key.Matches(msg, m.keys.nextPage):
				m.preview.ViewDown()
				return m, nil
			case key.Matches(msg, m.keys.prevPage):
				m.preview.ViewUp()
				return m, nil

			case key.Matches(msg, m.keys.down):
				m.preview.LineDown(1)
				return m, nil

			case key.Matches(msg, m.keys.up):
				m.preview.LineUp(1)
				return m, nil
			}
		}

		i, ok := m.list.SelectedItem().(item)
		if !ok { // no keys in current list; handle available options
			switch {
			case key.Matches(msg, m.keys.more):
				m.list.SetShowHelp(!m.list.ShowHelp())
				m.updatePaginator()

			case key.Matches(msg, m.keys.quit):
				if m.togglePinned {
					statusMsg := m.togglePinView()
					cmds = append(cmds, m.list.NewStatusMessage(statusMsg))
					return m, tea.Batch(cmds...)
				}
				m.ExitCode = 1
				return m, tea.Quit

			case key.Matches(msg, m.keys.forceQuit):
				m.ExitCode = 1
				return m, tea.Quit

			case key.Matches(msg, m.keys.togglePinned):
				statusMsg := m.togglePinView()
				cmds = append(cmds, m.list.NewStatusMessage(statusMsg))
				return m, tea.Batch(cmds...)
			}

			break
		}

		switch {

		case key.Matches(msg, m.keys.choose):
			if m.showConfirmation {
				cmds = m.handleConfirmationSelection(cmds)
				break
			}

			m, updatedCmds, toQuit := m.handleChooseOperation(i, cmds)
			if !toQuit {
				cmds = append(cmds, updatedCmds...)
				return m, tea.Batch(cmds...)
			}
			return m, tea.Quit

		case key.Matches(msg, m.keys.remove):
			selectedItems := m.selectedItems()
			var pinnedItemSelected bool
			m.itemCache = append(
				m.itemCache,
				SelectedItem{
					Index:     m.list.Index(),
					TimeStamp: i.timeStamp,
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
				m.confirmationList.Select(0)
				m.setConfirmationKeys(true)
				m.enableConfirmationKeys(true)
				m.showConfirmation = true
				break
			}

			currentIndex := m.list.Index()
			statusMsg := "Deleted: "

			if len(selectedItems) >= 1 {
				timeStamps := []string{}
				m.list.RemoveItem(currentIndex)
				m.removeMultiSelected()
				for _, item := range selectedItems {
					timeStamps = append(timeStamps, item.TimeStamp)
				}

				timeStamps = append(timeStamps, i.timeStamp)
				statusMsg += "*selected items*"
				if err := config.DeleteItems(timeStamps); err != nil {
					utils.LogERROR(fmt.Sprintf("failed to delete all items from history file: %s", err))
				}
			} else {
				m.list.RemoveItem(currentIndex)
				utils.HandleError(config.DeleteItems([]string{i.timeStamp}))
				statusMsg += i.title
			}

			if len(m.list.Items()) == 0 {
				m.list.SetShowStatusBar(false)
			}

			m.itemCache = []SelectedItem{}
			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
			)

		case key.Matches(msg, m.keys.togglePin):
			isPinned, err := config.TogglePinClipboardItem(i.timeStamp)
			utils.HandleError(err)
			m.togglePinUpdate()

			pinEvent := "Pinned"
			if isPinned {
				pinEvent = "Unpinned"
			}
			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("%s: %s", pinEvent, i.title))),
			)

		case key.Matches(msg, m.keys.togglePinned):
			newStatusMsg := m.togglePinView()
			cmds = append(
				cmds,
				m.list.NewStatusMessage(newStatusMsg),
			)

		case key.Matches(msg, m.keys.selectDown):
			if m.list.IsFiltered() {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("cannot select items with filter applied")),
				)
				break
			}
			m.toggleSelected("down")

		case key.Matches(msg, m.keys.selectUp):
			if m.list.IsFiltered() {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("cannot select items with filter applied")),
				)
				break
			}
			m.toggleSelected("up")

		case key.Matches(msg, m.keys.selectSingle):
			if m.list.IsFiltered() {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(statusMessageStyle("cannot select items with filter applied")),
				)
				break
			}
			m.toggleSelectedSingle()

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

		case key.Matches(msg, m.keys.preview):
			m.showPreview = !m.showPreview
			if config.ClipseConfig.ImageDisplay.Type == "kitty" {
				fmt.Print("\x1B_Ga=d\x1B\\")
			}
			if m.showPreview {
				content := m.styledPreviewContent(i.titleFull)
				if i.filePath != "null" {
					content = getImgPreview(i.filePath, m.preview.Width, m.preview.Height)
					if config.ClipseConfig.ImageDisplay.Type != "basic" {
						m.originalHeight = m.preview.Height
						m.preview.Height /= config.ClipseConfig.ImageDisplay.HeightCut
					}
				}
				m.preview.SetContent(content)
				var cmd tea.Cmd
				m.preview, cmd = m.preview.Update(msg)
				cmds = append(cmds, cmd)
				cmds = append(cmds, tea.ClearScreen)
				m.preview.GotoTop()
				m.setPreviewKeys(true)

				return m, tea.Batch(cmds...)
			}
			if i.filePath != "null" && config.ClipseConfig.ImageDisplay.Type != "basic" {
				m.preview.Height = m.originalHeight
			}
			m.setPreviewKeys(false)

		case key.Matches(msg, m.keys.quit):
			switch {
			case m.showPreview:
				if config.ClipseConfig.ImageDisplay.Type == "kitty" {
					fmt.Print("\x1B_Ga=d\x1B\\")
				}
				m.showPreview = !m.showPreview
				var cmd tea.Cmd
				m.preview, cmd = m.preview.Update(msg)
				cmds = append(cmds, cmd)
				m.preview.GotoTop()
				m.setPreviewKeys(false)
				if i.filePath != "null" && config.ClipseConfig.ImageDisplay.Type != "basic" {
					m.preview.Height = m.originalHeight
				}
				return m, tea.Batch(cmds...)

			case m.list.IsFiltered():
				m.list.ResetFilter()

			case m.showConfirmation:
				m.confirmationList.Select(0)
				m.setConfirmationKeys(false)
				m.enableConfirmationKeys(false)
				m.showConfirmation = false
				return m, tea.Batch(cmds...)

			default:
				m.ExitCode = 1
				return m, tea.Quit
			}

		case key.Matches(msg, m.keys.forceQuit):
			m.ExitCode = 1
			return m, tea.Quit
		}

	case tea.MouseMsg:
		if !config.ClipseConfig.EnableMouse {
			break
		}
		if m.showPreview {
			break
		}
		i, ok := m.list.SelectedItem().(item)
		if !ok {
			break
		}

		switch msg.Action {
		case tea.MouseActionMotion:
			listItemsStart := 5 // first 5 lines take up title and info
			if m.showConfirmation {
				listItemsStart = listItemsStart - 1 // confirmation list starts from line 4
			}
			linesPerItem := 3 // may need to make configurable if desc is disabled

			if msg.Y >= listItemsStart {
				relativeY := msg.Y - listItemsStart
				hoveredIndex := relativeY / linesPerItem

				if hoveredIndex < len(m.list.Items()) {
					if m.showConfirmation {
						m.confirmationList.Select(hoveredIndex)
						break
					}
					m.list.Select(hoveredIndex)
				}
			}

		case tea.MouseActionPress:
			switch msg.Button {
			case tea.MouseButtonLeft:
				if m.showConfirmation {
					cmds = m.handleConfirmationSelection(cmds)
					break
				}

				m, updatedCmds, toQuit := m.handleChooseOperation(i, cmds)
				if !toQuit {
					cmds = append(cmds, updatedCmds...)
					return m, tea.Batch(cmds...)
				}
				return m, tea.Quit

			case tea.MouseButtonRight:
				if m.showPreview || m.showConfirmation {
					break
				}
				i.selected = !i.selected
				m.list.SetItem(m.list.Index(), i)

			case tea.MouseButtonWheelUp:
				if m.list.Index() > 0 {
					m.list.Select(m.list.Index() - 1)
				}

			case tea.MouseButtonWheelDown:
				if m.list.Index() < len(m.list.Items())-1 {
					m.list.Select(m.list.Index() + 1)
				}
			}
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	m.confirmationList, cmd = m.confirmationList.Update(msg)
	cmds = append(cmds, cmd)

	m.preview, cmd = m.preview.Update(msg)
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

	switch {
	case item.selected:
		item.selected = false
	case m.prevDirection == direction && !item.selected:
		item.selected = true
	default:
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

func (m Model) handleChooseOperation(i item, cmds []tea.Cmd) (Model, []tea.Cmd, bool) {
	selectedItems := m.selectedItems()
	if len(selectedItems) < 1 {
		switch {
		case i.filePath != "null":
			display.DisplayServer.CopyImage(i.filePath)

		default:
			display.DisplayServer.CopyText(i.titleFull)
		}

		if KeepEnabled {
			cmds = append(
				cmds,
				m.list.NewStatusMessage(statusMessageStyle("Copied to clipboard: "+i.title)),
			)
			return m, cmds, false
		}
		return m, cmds, true
	}

	yank := ""
	for _, item := range selectedItems {
		if i.titleFull != item.Value {
			yank += item.Value + "\n"
		}
	}
	yank += i.titleFull

	display.DisplayServer.CopyText(yank)

	if KeepEnabled {
		statusMsg := "Copied to clipboard: *selected items*"
		display.DisplayServer.CopyText(yank)
		cmds = append(
			cmds,
			m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
		)
		return m, cmds, false
	}
	display.DisplayServer.CopyText(yank)
	return m, cmds, true
}

func (m *Model) deleteCachedItems(cmds []tea.Cmd) []tea.Cmd {
	timeStamps := []string{}
	for _, item := range m.itemCache {
		timeStamps = append(timeStamps, item.TimeStamp)
		m.removeCachedItem(item.TimeStamp)
	}

	statusMsg := "Deleted: *selected items*"
	if len(m.itemCache) == 1 {
		statusMsg += strings.Replace(statusMsg, "*selected items*", m.itemCache[0].Value, 1)
	}

	if err := config.DeleteItems(timeStamps); err != nil {
		utils.LogERROR(fmt.Sprintf("could not delete all items from history: %s", err))
	}

	cmds = append(
		cmds,
		m.list.NewStatusMessage(statusMessageStyle(statusMsg)),
	)
	m.itemCache = []SelectedItem{}

	if len(m.list.Items()) == 0 {
		m.list.SetShowStatusBar(false)
	}

	return cmds
}

func (m *Model) removeCachedItem(ts string) {
	// helper func for deleteCachedItems
	items := m.list.Items()
	for i := len(items) - 1; i >= 0; i-- {
		if item, ok := items[i].(item); ok && item.timeStamp == ts {
			m.list.RemoveItem(i)
		}
	}
}

func (m *Model) handleConfirmationSelection(cmds []tea.Cmd) []tea.Cmd {
	if m.showConfirmation && m.confirmationList.Index() == 0 { // No
		m.itemCache = []SelectedItem{}
		m.showConfirmation = false
		m.setPreviewKeys(false)
		m.enableConfirmationKeys(false)
		m.setConfirmationKeys(false)
		return cmds

	}
	// yes selected
	m.showConfirmation = false
	m.setPreviewKeys(false)
	m.enableConfirmationKeys(false)
	m.setConfirmationKeys(false)
	cmds = m.deleteCachedItems(cmds)
	return cmds
}

// enable/disable the main keys relevant to the preview view
func (m *Model) setPreviewKeys(v bool) {
	m.list.KeyMap.CursorUp.SetEnabled(!v)
	m.list.KeyMap.CursorDown.SetEnabled(!v)
	m.list.KeyMap.Filter.SetEnabled(!v)
	m.list.KeyMap.GoToEnd.SetEnabled(!v)
	m.list.KeyMap.GoToStart.SetEnabled(!v)
	m.list.KeyMap.Quit.SetEnabled(!v)
	m.list.KeyMap.NextPage.SetEnabled(!v)
	m.list.KeyMap.PrevPage.SetEnabled(!v)
	m.list.KeyMap.ShowFullHelp.SetEnabled(!v)

	m.keys.remove.SetEnabled(!v)
	m.keys.togglePin.SetEnabled(!v)
	m.keys.togglePinned.SetEnabled(!v)
	m.keys.selectDown.SetEnabled(!v)
	m.keys.selectUp.SetEnabled(!v)
	m.keys.selectSingle.SetEnabled(!v)
	m.keys.clearSelected.SetEnabled(!v)
}

// enable/disable the main keys relevant to the confirmation view
func (m *Model) setConfirmationKeys(v bool) {
	m.list.KeyMap.CursorUp.SetEnabled(!v)
	m.list.KeyMap.CursorDown.SetEnabled(!v)
	m.list.KeyMap.Filter.SetEnabled(!v)
	m.list.KeyMap.GoToEnd.SetEnabled(!v)
	m.list.KeyMap.GoToStart.SetEnabled(!v)
	m.list.KeyMap.Quit.SetEnabled(!v)
	m.list.KeyMap.NextPage.SetEnabled(!v)
	m.list.KeyMap.PrevPage.SetEnabled(!v)
	m.list.KeyMap.ShowFullHelp.SetEnabled(!v)

	m.keys.remove.SetEnabled(!v)
	m.keys.togglePin.SetEnabled(!v)
	m.keys.togglePinned.SetEnabled(!v)
	m.keys.selectDown.SetEnabled(!v)
	m.keys.selectUp.SetEnabled(!v)
	m.keys.selectSingle.SetEnabled(!v)
	m.keys.clearSelected.SetEnabled(!v)
	m.keys.preview.SetEnabled(!v)
}

// enable/disable the navigation for the confirmation list
func (m *Model) enableConfirmationKeys(v bool) {
	m.confirmationList.KeyMap.CursorUp.SetEnabled(v)
	m.confirmationList.KeyMap.CursorDown.SetEnabled(v)
}

func (m *Model) setQuitEnabled(v bool) {
	m.list.KeyMap.Quit.SetEnabled(v)
}

func (m *Model) togglePinView() string {
	m.togglePinned = !m.togglePinned
	m.list.Title = clipboardTitle

	clipboardItems := config.GetHistory()
	filteredItems := filterItems(clipboardItems, m.togglePinned)

	if len(filteredItems) == 0 {
		if m.togglePinned {
			m.togglePinned = !m.togglePinned
			return statusMessageStyle("No pinned items")
		}
		return ""
	}

	if m.togglePinned {
		m.list.Title = "Pinned " + m.list.Title
	}

	for i := len(m.list.Items()) - 1; i >= 0; i-- { // clear all items
		m.list.RemoveItem(i)
	}
	for _, i := range filteredItems { // redraw all required items
		m.list.InsertItem(len(m.list.Items()), i)
	}

	return ""
}

func keepEnabled() bool {
	if len(os.Args) > 1 && os.Args[1] == "keep" {
		return true
	}
	return false
}
