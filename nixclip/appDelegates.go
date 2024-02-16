package main

import (
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

/* This is were we define additional config to add to our
base-level bubbletea app. Here including keybinds only.
*/

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

/*
Additional short/full help entries. This satisfies the help.KeyMap interface and

	is entirely optional.
*/
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

// final config map for new keys
func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "copy"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	/* This is where the additional keybinding actions are defined:
	   - enter/reurn: copies selected item to the clipboard and adds a status message
	   - backspace/delete: removes item from list view and json file

	*/
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string
		var fullValue string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
			fullValue = i.TitleFull()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				err := clipboard.WriteAll(fullValue)
				if err != nil {
					panic(err)
				}
				if len(os.Args) > 1 {
					closeShell(os.Args[1])
				}
				return m.NewStatusMessage(statusMessageStyle("Copied to clipboard: " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				fullPath := getFullPath()
				err := deleteJsonItem(fullPath, fullValue)
				handleError(err)
				return m.NewStatusMessage(statusMessageStyle("Deleted: " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
