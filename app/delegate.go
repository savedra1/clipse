package app

import (
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/shell"
	"github.com/savedra1/clipse/utils"
)

func (parentModel *model) newItemDelegate(keys *keyMap) list.DefaultDelegate {
	/* This is where the additional keybinding actions are defined:
	   - enter/reurn: copies selected item to the clipboard and adds a status message
	   - backspace/delete: removes item from list view and json file
	   - p: pins/unpins an item
	   - tab: toggles pinned items
	*/
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var fullValue string
		var fp string
		var desc string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
			fullValue = i.TitleFull()
			fp = i.FilePath()
			desc = i.TimeStamp()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
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
						return m.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + title))
					}
				}
				return tea.Quit

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
					m.SetShowStatusBar(false)
				}
				go func() { // stop cached clipboard item repopulating
					currentContent, _ := clipboard.ReadAll()
					if currentContent == fullValue {
						clipboard.WriteAll("")
					}
					err := config.DeleteJsonItem(desc) // This func will also delete the temoraily stored image if filepath present
					utils.HandleError(err)
				}()

				return m.NewStatusMessage(statusMessageStyle("Deleted: " + title))

			case key.Matches(msg, keys.togglePin):
				if len(m.Items()) == 0 {
					keys.togglePin.SetEnabled(false)
				}

				// update pinned status in history file
				isPinned, err := config.TogglePinClipboardItem(desc)
				utils.HandleError(err)

				if isPinned {
					return m.NewStatusMessage(statusMessageStyle("UnPinned: " + title))
				}
				return m.NewStatusMessage(statusMessageStyle("Pinned: " + title))

			case key.Matches(msg, keys.togglePinned):
				if len(m.Items()) == 0 {
					keys.togglePinned.SetEnabled(false)
				}

				parentModel.togglePinned = !parentModel.togglePinned

				if parentModel.togglePinned {
					m.Title = "Pinned " + clipboardTitle
				} else {
					m.Title = clipboardTitle
				}

				clipboardItems := config.GetHistory()
				filteredItems := filterItems(clipboardItems, parentModel.togglePinned)

				if len(filteredItems) == 0 {
					m.Title = clipboardTitle
					return m.NewStatusMessage(statusMessageStyle("No pinned items"))
				}

				for i := len(m.Items()) - 1; i >= 0; i-- { // clear all items
					m.RemoveItem(i)
				}

				for _, item := range filteredItems { // adds all required items
					m.InsertItem(len(m.Items()), item)
				}
			}
		}
		return nil
	}

	// add custom options to the default full help view
	help := []key.Binding{
		keys.choose,
		keys.remove,
		keys.togglePin,
		keys.togglePinned,
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
