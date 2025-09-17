package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"

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
	previewBack   key.Binding
}

func newKeyMap(config map[string]string) *keyMap {
	previewChar := config["preview"]
	if previewChar == " " {
		previewChar = spaceChar
	}

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
			key.WithHelp(previewChar, "preview"),
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
		previewBack: key.NewBinding(
			key.WithKeys(config["previewBack"]),
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

func newFilterKeymap(config map[string]string) *filterKeyMap {
	return &filterKeyMap{
		apply: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp(config["choose"], "apply"),
		),
		cancel: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(config["quit"], "cancel"),
		),
		yankMatches: key.NewBinding(
			key.WithKeys(config["yankFilter"]),
			key.WithHelp(config["yankFilter"], "yank matched"),
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

func newConfirmationKeymap(config map[string]string) *confirmationKeyMap {
	return &confirmationKeyMap{
		up: key.NewBinding(
			key.WithKeys(config["up"]),
			key.WithHelp(config["up"], "↑"),
		),
		down: key.NewBinding(
			key.WithKeys(config["down"]),
			key.WithHelp(config["down"], "↓"),
		),
		choose: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp(config["choose"], "choose"),
		),
	}
}

func (ck confirmationKeyMap) ConfirmationHelp() []key.Binding {
	return []key.Binding{
		ck.up, ck.down, ck.choose,
	}
}

type previewKeymap struct {
	up       key.Binding
	down     key.Binding
	back     key.Binding
	pageDown key.Binding
	pageUp   key.Binding
	choose   key.Binding
}

func newPreviewKeyMap() *previewKeymap {
	config := config.ClipseConfig.KeyBindings

	previewChar := config["preview"]
	if previewChar == " " {
		previewChar = spaceChar
	}

	return &previewKeymap{
		up: key.NewBinding(
			key.WithKeys(config["up"]),
			key.WithHelp(config["up"], "↑"),
		),
		down: key.NewBinding(
			key.WithKeys(config["down"]),
			key.WithHelp(config["down"], "↓"),
		),
		pageDown: key.NewBinding(
			key.WithKeys("PgDn"),
			key.WithHelp("PgDn", "page down"),
		),
		pageUp: key.NewBinding(
			key.WithKeys("PgUp"),
			key.WithHelp("PgUp", "page up"),
		),
		choose: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp("↵", "copy"),
		),
		back: key.NewBinding(
			key.WithKeys(config["preview"], config["previewBack"]),
			key.WithHelp(previewChar+" / "+config["previewBack"], "back"),
		),
	}
}

func (pk previewKeymap) PreviewHelp() []key.Binding {
	return []key.Binding{
		pk.up, pk.down, pk.pageDown, pk.pageUp, pk.back, pk.choose,
	}
}

func defaultOverrides(config map[string]string) list.KeyMap {
	return list.KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys(config["up"]),
			key.WithHelp(config["up"], "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys(config["down"]),
			key.WithHelp(config["down"], "down"),
		),
		NextPage: key.NewBinding(
			key.WithKeys(config["nextPage"]),
			key.WithHelp(config["nextPage"], "page down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys(config["prevPage"]),
			key.WithHelp(config["prevPage"], "page up"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys(config["home"]),
			key.WithHelp(config["home"], "start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys(config["end"]),
			key.WithHelp(config["end"], "end"),
		),
		Filter: key.NewBinding(
			key.WithKeys(config["filter"]),
			key.WithHelp(config["filter"], "filter"),
		),
		Quit: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(config["quit"], "quit"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "less"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithDisabled(),
		),
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithDisabled(),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithDisabled(),
		),
		ForceQuit: key.NewBinding(key.WithDisabled()),
	}
}
