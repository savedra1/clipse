package app

import (
	"strings"
	
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

var charSub = map[string]string{
	"enter":     enterChar,
	" ":         spaceChar,
	"up":        upChar,
	"down":      downChar,
	"right":     rightChar,
	"left":      leftChar,
	"backspace": backspaceChar,
}

func getHelpChar(configChar string) string {
	if val, ok := charSub[configChar]; ok {
		return val
	}
	return configChar
}

func dedupeNonEmptyKeys(keys ...string) []string {
	seen := make(map[string]struct{}, len(keys))
	var result []string
	for _, k := range keys {
		if k == "" {
			continue
		}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, k)
	}
	return result
}

func formatHelpKeys(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	help := make([]string, len(keys))
	for i, k := range keys {
		help[i] = getHelpChar(k)
	}
	return strings.Join(help, " / ")
}

func newKeyMap(config map[string]string) *keyMap {
	return &keyMap{
		filter: key.NewBinding(
			key.WithKeys(config["filter"]),
			key.WithHelp(getHelpChar(config["filter"]), "filter"),
		),
		quit: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(getHelpChar(config["quit"]), "quit"),
		),
		more: key.NewBinding(
			key.WithKeys(config["more"]),
			key.WithHelp(getHelpChar(config["more"]), "more"),
		),
		choose: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp(getHelpChar(config["choose"]), "copy"),
		),
		remove: key.NewBinding(
			key.WithKeys(config["remove"]),
			key.WithHelp(getHelpChar(config["remove"]), "delete"),
		),
		togglePin: key.NewBinding(
			key.WithKeys(config["togglePin"]),
			key.WithHelp(getHelpChar(config["togglePin"]), "pin/unpin"),
		),
		togglePinned: key.NewBinding(
			key.WithKeys(config["togglePinned"]),
			key.WithHelp(getHelpChar(config["togglePinned"]), "show pinned"),
		),
		preview: key.NewBinding(
			key.WithKeys(config["preview"]),
			key.WithHelp(getHelpChar(config["preview"]), "preview"),
		),
		selectDown: key.NewBinding(
			key.WithKeys(config["selectDown"]),
			key.WithHelp(getHelpChar(config["selectDown"]), "select down"),
		),
		selectUp: key.NewBinding(
			key.WithKeys(config["selectUp"]),
			key.WithHelp(getHelpChar(config["selectUp"]), "select up"),
		),
		selectSingle: key.NewBinding(
			key.WithKeys(config["selectSingle"]),
			key.WithHelp(getHelpChar(config["selectSingle"]), "select single"),
		),
		clearSelected: key.NewBinding(
			key.WithKeys(config["clearSelected"]),
			key.WithHelp(getHelpChar(config["clearSelected"]), "clear selected"),
		),
		yankFilter: key.NewBinding(
			key.WithKeys(config["yankFilter"]),
			key.WithHelp(getHelpChar(config["yankFilter"]), "yank filter results"),
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
		{k.selectDown, k.selectUp, k.selectSingle, k.clearSelected},
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
			key.WithHelp(getHelpChar(config["choose"]), "apply"),
		),
		cancel: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(getHelpChar(config["quit"]), "cancel"),
		),
		yankMatches: key.NewBinding(
			key.WithKeys(config["yankFilter"]),
			key.WithHelp(getHelpChar(config["yankFilter"]), "yank matched"),
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
	back   key.Binding
}

func newConfirmationKeymap(config map[string]string) *confirmationKeyMap {
	return &confirmationKeyMap{
		up: key.NewBinding(
			key.WithKeys(config["up"]),
			key.WithHelp(getHelpChar(config["up"]), "up"),
		),
		down: key.NewBinding(
			key.WithKeys(config["down"]),
			key.WithHelp(getHelpChar(config["down"]), "down"),
		),
		choose: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp(getHelpChar(config["choose"]), "choose"),
		),
		back: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(getHelpChar(config["quit"]), "back"),
		),
	}
}

func (ck confirmationKeyMap) ConfirmationHelp() []key.Binding {
	return []key.Binding{
		ck.up, ck.down, ck.choose, ck.back,
	}
}

type previewKeymap struct {
	up       key.Binding
	down     key.Binding
	pageDown key.Binding
	pageUp   key.Binding
	back     key.Binding
	choose   key.Binding
	quit     key.Binding
}

func newPreviewKeyMap(config map[string]string) *previewKeymap {
	backKeys := dedupeNonEmptyKeys(config["preview"], config["previewBack"], config["quit"])
	quitKeys := dedupeNonEmptyKeys(config["previewQuit"])

	backBinding := key.NewBinding(key.WithDisabled())
	if len(backKeys) > 0 {
		backBinding = key.NewBinding(
			key.WithKeys(backKeys...),
			key.WithHelp(formatHelpKeys(backKeys), "back"),
		)
	}

	quitBinding := key.NewBinding(key.WithDisabled())
	if len(quitKeys) > 0 {
		quitBinding = key.NewBinding(
			key.WithKeys(quitKeys...),
			key.WithHelp(formatHelpKeys(quitKeys), "quit"),
		)
	}

	return &previewKeymap{
		up: key.NewBinding(
			key.WithKeys(config["up"]),
			key.WithHelp(getHelpChar(config["up"]), "up"),
		),
		down: key.NewBinding(
			key.WithKeys(config["down"]),
			key.WithHelp(getHelpChar(config["down"]), "down"),
		),
		pageDown: key.NewBinding(
			key.WithKeys(config["nextPage"]),
			key.WithHelp(getHelpChar(config["nextPage"]), "page down"),
		),
		pageUp: key.NewBinding(
			key.WithKeys(config["prevPage"]),
			key.WithHelp(getHelpChar(config["prevPage"]), "page up"),
		),
		back: backBinding,
		choose: key.NewBinding(
			key.WithKeys(config["choose"]),
			key.WithHelp(getHelpChar(config["choose"]), "copy"),
		),
		quit: quitBinding,
	}
}

func (pk previewKeymap) PreviewHelp() []key.Binding {
	bindings := []key.Binding{pk.up, pk.down, pk.pageDown, pk.pageUp, pk.choose}
	if len(pk.back.Keys()) > 0 {
		bindings = append(bindings, pk.back)
	}
	if len(pk.quit.Keys()) > 0 {
		bindings = append(bindings, pk.quit)
	}
	return bindings
}

func defaultOverrides(config map[string]string) list.KeyMap {
	return list.KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys(config["up"]),
			key.WithHelp(getHelpChar(config["up"]), "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys(config["down"]),
			key.WithHelp(getHelpChar(config["down"]), "down"),
		),
		NextPage: key.NewBinding(
			key.WithKeys(config["nextPage"]),
			key.WithHelp(getHelpChar(config["nextPage"]), "page down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys(config["prevPage"]),
			key.WithHelp(getHelpChar(config["prevPage"]), "page up"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys(config["home"]),
			key.WithHelp(getHelpChar(config["home"]), "start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys(config["end"]),
			key.WithHelp(getHelpChar(config["end"]), "end"),
		),
		Filter: key.NewBinding(
			key.WithKeys(config["filter"]),
			key.WithHelp(getHelpChar(config["filter"]), "filter"),
		),
		Quit: key.NewBinding(
			key.WithKeys(config["quit"]),
			key.WithHelp(getHelpChar(config["quit"]), "quit"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys(config["more"]),
			key.WithHelp(getHelpChar(config["more"]), "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys(config["more"]),
			key.WithHelp(getHelpChar(config["more"]), "less"),
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
