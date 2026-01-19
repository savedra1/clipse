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
	forceQuit     key.Binding
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
	"space":     spaceChar,
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

func parseKeys(input string) []string {
	keys := strings.Split(input, ",")
	for i, key := range keys {
		k := strings.TrimSpace(key)
		if k == "space" {
			k = " "
		}
		keys[i] = k
	}
	return keys
}

func newKeyMap(config map[string]string) *keyMap {
	return &keyMap{
		filter: key.NewBinding(
			key.WithKeys(parseKeys(config["filter"])...),
			key.WithHelp(getHelpChar(config["filter"]), "filter"),
		),
		quit: key.NewBinding(
			key.WithKeys(parseKeys(config["quit"])...),
			key.WithHelp(getHelpChar(config["quit"]), "quit"),
		),
		forceQuit: key.NewBinding(
			key.WithKeys(parseKeys(config["forceQuit"])...),
			key.WithHelp(getHelpChar(config["forceQuit"]), "force quit"),
		),
		more: key.NewBinding(
			key.WithKeys(parseKeys(config["more"])...),
			key.WithHelp(getHelpChar(config["more"]), "more"),
		),
		choose: key.NewBinding(
			key.WithKeys(parseKeys(config["choose"])...),
			key.WithHelp(getHelpChar(config["choose"]), "copy"),
		),
		remove: key.NewBinding(
			key.WithKeys(parseKeys(config["remove"])...),
			key.WithHelp(getHelpChar(config["remove"]), "delete"),
		),
		togglePin: key.NewBinding(
			key.WithKeys(parseKeys(config["togglePin"])...),
			key.WithHelp(getHelpChar(config["togglePin"]), "pin/unpin"),
		),
		togglePinned: key.NewBinding(
			key.WithKeys(parseKeys(config["togglePinned"])...),
			key.WithHelp(getHelpChar(config["togglePinned"]), "show pinned"),
		),
		preview: key.NewBinding(
			key.WithKeys(parseKeys(config["preview"])...),
			key.WithHelp(getHelpChar(config["preview"]), "preview"),
		),
		selectDown: key.NewBinding(
			key.WithKeys(parseKeys(config["selectDown"])...),
			key.WithHelp(getHelpChar(config["selectDown"]), "select down"),
		),
		selectUp: key.NewBinding(
			key.WithKeys(parseKeys(config["selectUp"])...),
			key.WithHelp(getHelpChar(config["selectUp"]), "select up"),
		),
		selectSingle: key.NewBinding(
			key.WithKeys(parseKeys(config["selectSingle"])...),
			key.WithHelp(getHelpChar(config["selectSingle"]), "select single"),
		),
		clearSelected: key.NewBinding(
			key.WithKeys(parseKeys(config["clearSelected"])...),
			key.WithHelp(getHelpChar(config["clearSelected"]), "clear selected"),
		),
		yankFilter: key.NewBinding(
			key.WithKeys(parseKeys(config["yankFilter"])...),
			key.WithHelp(getHelpChar(config["yankFilter"]), "yank filter results"),
		),
		up: key.NewBinding(
			key.WithKeys(parseKeys(config["up"])...),
		),
		down: key.NewBinding(
			key.WithKeys(parseKeys(config["down"])...),
		),
		nextPage: key.NewBinding(
			key.WithKeys(parseKeys(config["nextPage"])...),
		),
		prevPage: key.NewBinding(
			key.WithKeys(parseKeys(config["prevPage"])...),
		),
		home: key.NewBinding(
			key.WithKeys(parseKeys(config["home"])...),
		),
		end: key.NewBinding(
			key.WithKeys(parseKeys(config["end"])...),
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
		{k.filter, k.quit, k.forceQuit},
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
			key.WithKeys("enter"), // hardcoded enter prevents `choose` key from interrupting filter
			key.WithHelp(enterChar, "apply"),
		),
		cancel: key.NewBinding(
			key.WithKeys("esc"), // hardcoded esc prevents custom `quit` key from interrupting filter
			key.WithHelp(getHelpChar("esc"), "cancel"),
		),
		yankMatches: key.NewBinding(
			key.WithKeys(parseKeys(config["yankFilter"])...),
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
			key.WithKeys(parseKeys(config["up"])...),
			key.WithHelp(getHelpChar(config["up"]), "up"),
		),
		down: key.NewBinding(
			key.WithKeys(parseKeys(config["down"])...),
			key.WithHelp(getHelpChar(config["down"]), "down"),
		),
		choose: key.NewBinding(
			key.WithKeys(parseKeys(config["choose"])...),
			key.WithHelp(getHelpChar(config["choose"]), "choose"),
		),
		back: key.NewBinding(
			key.WithKeys(parseKeys(config["quit"])...),
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
}

func newPreviewKeyMap(config map[string]string) *previewKeymap {
	return &previewKeymap{
		up: key.NewBinding(
			key.WithKeys(parseKeys(config["up"])...),
			key.WithHelp(getHelpChar(config["up"]), "up"),
		),
		down: key.NewBinding(
			key.WithKeys(parseKeys(config["down"])...),
			key.WithHelp(getHelpChar(config["down"]), "down"),
		),
		pageDown: key.NewBinding(
			key.WithKeys(parseKeys(config["nextPage"])...),
			key.WithHelp(getHelpChar(config["nextPage"]), "page down"),
		),
		pageUp: key.NewBinding(
			key.WithKeys(parseKeys(config["prevPage"])...),
			key.WithHelp(getHelpChar(config["prevPage"]), "page up"),
		),
		back: key.NewBinding(
			key.WithKeys(append(parseKeys(config["preview"]), parseKeys(config["quit"])...)...),
			key.WithHelp(getHelpChar(config["preview"])+" / "+getHelpChar(config["quit"]), "back"),
		),
		choose: key.NewBinding(
			key.WithKeys(parseKeys(config["choose"])...),
			key.WithHelp(getHelpChar(config["choose"]), "copy"),
		),
	}
}

func (pk previewKeymap) PreviewHelp() []key.Binding {
	return []key.Binding{
		pk.up, pk.down,
		pk.pageDown, pk.pageUp,
		pk.choose,
	}
}

// keys defined here do not need to be handled via the update func unless purposefully disabled
func defaultOverrides(config map[string]string) list.KeyMap {
	return list.KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys(parseKeys(config["up"])...),
			key.WithHelp(getHelpChar(config["up"]), "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys(parseKeys(config["down"])...),
			key.WithHelp(getHelpChar(config["down"]), "down"),
		),
		NextPage: key.NewBinding(
			key.WithKeys(parseKeys(config["nextPage"])...),
			key.WithHelp(getHelpChar(config["nextPage"]), "page down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys(parseKeys(config["prevPage"])...),
			key.WithHelp(getHelpChar(config["prevPage"]), "page up"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys(parseKeys(config["home"])...),
			key.WithHelp(getHelpChar(config["home"]), "start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys(parseKeys(config["end"])...),
			key.WithHelp(getHelpChar(config["end"]), "end"),
		),
		Filter: key.NewBinding(
			key.WithKeys(parseKeys(config["filter"])...),
			key.WithHelp(getHelpChar(config["filter"]), "filter"),
		),
		Quit: key.NewBinding(
			key.WithDisabled(), // quit keys handles by update func
			key.WithHelp(getHelpChar(config["quit"]), "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithDisabled(),
			key.WithHelp("ctrl+c / "+getHelpChar(config["forceQuit"]), "force quit"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys(parseKeys(config["more"])...),
			key.WithHelp(getHelpChar(config["more"]), "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys(parseKeys(config["more"])...),
			key.WithHelp(getHelpChar(config["more"]), "less"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter"), // hardcoded enter prevents `choose` key from interrupting filter
			key.WithDisabled(),
		),
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"), // hardcoded esc prevents custom `quit` key from interrupting filter
			key.WithDisabled(),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys(parseKeys(config["quit"])...),
			key.WithDisabled(),
		),
	}
}
