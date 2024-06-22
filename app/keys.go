package app

import (
	"github.com/charmbracelet/bubbles/key"
)

// default keybind definitions
type keyMap struct {
	filter       key.Binding
	quit         key.Binding
	more         key.Binding
	choose       key.Binding
	remove       key.Binding
	togglePin    key.Binding
	togglePinned key.Binding
	selectDown   key.Binding
	selectUp     key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q/esc", "quit"),
		),
		more: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "copy"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x/⌫", "delete"),
		),
		togglePin: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "pin/unpin"),
		),
		togglePinned: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("↹", "show pinned"),
		),
		selectDown: key.NewBinding(
			key.WithKeys("shift+down"),
			key.WithHelp("⇧+↓/↑", "select"),
		),
		selectUp: key.NewBinding(
			key.WithKeys("shift+up"),
			key.WithHelp("⇧+↓/↑", "select"),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.choose, k.remove, k.filter, k.togglePin, k.togglePinned, k.more,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.choose, k.remove},
		{k.togglePin, k.togglePinned},
		{k.filter, k.quit},
	}
}

type filterKeyMap struct {
	apply  key.Binding
	cancel key.Binding
}

func newFilterKeymap() *filterKeyMap {
	return &filterKeyMap{
		apply: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "apply"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}

func (fk filterKeyMap) FilterHelp() []key.Binding {
	return []key.Binding{
		fk.apply, fk.cancel,
	}
}
