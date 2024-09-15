package app

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/savedra1/clipse/config"
)

// default keybind definitions
type keyMap struct {
	filter        key.Binding
	quit          key.Binding
	more          key.Binding
	choose        key.Binding
	remove        key.Binding
	togglePin     key.Binding
	togglePinned  key.Binding
	preview       key.Binding
	selectDown    key.Binding
	selectUp      key.Binding
	selectSingle  key.Binding
	clearSelected key.Binding
	yankFilter    key.Binding
	up            key.Binding
	down          key.Binding
	nextPage      key.Binding
	prevPage      key.Binding
	home          key.Binding
	end           key.Binding
}

func newKeyMap() *keyMap {
	config := config.ClipseConfig.KeyBindings

	return &keyMap{
		filter: key.NewBinding(
			key.WithKeys(config["filter"]),
			key.WithHelp(config["filter"], "filter"),
		),
		quit: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(config["quit"], "quit"),
		),
		more: key.NewBinding(
			key.WithKeys(config["more"]),
			key.WithHelp(config["more"], "more"),
		),
		choose: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp("↵", "copy"),
		),
		remove: key.NewBinding(
			key.WithKeys(config["remove"]),
			key.WithHelp(config["remove"], "delete"),
		),
		togglePin: key.NewBinding(
			key.WithKeys(config["togglePin"]),
			key.WithHelp(config["togglePin"], "pin/unpin"),
		),
		togglePinned: key.NewBinding(
			key.WithKeys(config["togglePinned"]),
			key.WithHelp(config["togglePinned"], "show pinned"),
		),
		preview: key.NewBinding(
			key.WithKeys(config["preview"]),
			key.WithHelp(config["preview"], "preview"),
		),
		selectDown: key.NewBinding(
			key.WithKeys(config["selectDown"]),
			key.WithHelp(config["selectDown"], "select"),
		),
		selectUp: key.NewBinding(
			key.WithKeys(config["selectUp"]),
			key.WithHelp(config["selectUp"], "select"),
		),
		selectSingle: key.NewBinding(
			key.WithKeys(config["selectSingle"]),
			key.WithHelp(config["selectSingle"], "select single"),
		),
		clearSelected: key.NewBinding(
			key.WithKeys(config["clearSelected"]),
			key.WithHelp(config["clearSelected"], "clear selected"),
		),
		yankFilter: key.NewBinding(
			key.WithKeys(config["yankFilter"]),
			key.WithHelp(config["yankFilter"], "yank filter results"),
		),
		up: key.NewBinding(
			key.WithKeys(config["up"]),
		),
		down: key.NewBinding(
			key.WithKeys(config["down"]),
		),
		nextPage: key.NewBinding(
			key.WithKeys(config["nextPage"]),
		),
		prevPage: key.NewBinding(
			key.WithKeys(config["prevPage"]),
		),
		home: key.NewBinding(
			key.WithKeys(config["home"]),
		),
		end: key.NewBinding(
			key.WithKeys(config["end"]),
		),
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.choose, k.remove, k.togglePin, k.togglePinned, k.more,
	}
}

// not currently in use as intentionally being overridden by the default
// full help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down, k.home, k.end},
		{k.choose, k.remove},
		{k.togglePin, k.togglePinned},
		{k.selectDown, k.selectSingle, k.yankFilter},
		{k.filter, k.quit},
	}
}

// used only for the default filter input view
type filterKeyMap struct {
	apply       key.Binding
	cancel      key.Binding
	yankMatches key.Binding
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
		yankMatches: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "yank matched"),
		),
	}
}

func (fk filterKeyMap) FilterHelp() []key.Binding {
	return []key.Binding{
		fk.apply, fk.cancel, fk.yankMatches,
	}
}

type confirmationKeyMap struct {
	up     key.Binding
	down   key.Binding
	choose key.Binding
}

func newConfirmationKeymap() *confirmationKeyMap {
	return &confirmationKeyMap{
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↓/k", "down"),
		),
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "choose"),
		),
	}
}

func (ck confirmationKeyMap) ConfirmationHelp() []key.Binding {
	return []key.Binding{
		ck.up, ck.down, ck.choose,
	}
}

// Not currently used

type previewKeymap struct {
	up       key.Binding
	down     key.Binding
	back     key.Binding
	pageDown key.Binding
	pageUp   key.Binding
}

func newPreviewKeyMap() *previewKeymap {
	return &previewKeymap{
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↓/j", "down"),
		),
		pageDown: key.NewBinding(
			key.WithKeys("PgDn"),
			key.WithHelp("PgDn", "page down"),
		),
		pageUp: key.NewBinding(
			key.WithKeys("PgUp"),
			key.WithHelp("PgUp", "page up"),
		),
		back: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("␣", "back"),
		),
	}
}

func (pk previewKeymap) PreviewHelp() []key.Binding {
	return []key.Binding{
		pk.up, pk.down, pk.pageDown, pk.pageUp, pk.back,
	}
}
