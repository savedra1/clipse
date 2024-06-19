package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

/* This is where the base level configuration for the
   bubbletea CLI app is defined.

   Base level config includes:
   - Color scheme of text
   - Font
   - Key bindings
   - List structure
   - Help menu
   - Default actions
*/

var (
	// base styling config using lipgloss
	appStyle           = lipgloss.NewStyle().Padding(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type model struct {
	// model pulls all relevant elems together for rendering
	list         list.Model    // list items
	keys         *keyMap       // keybindings
	filterKeys   *filterKeyMap // keybindings for filter view
	help         help.Model
	togglePinned bool // pinned indicator
	showFullHelp bool // whether full help menu is shown
}

type item struct {
	// (Each Row in clipboard view)
	title       string
	titleFull   string
	timeStamp   string
	description string
	filePath    string
	pinned      bool
}

func (i item) Title() string       { return i.title }
func (i item) TitleFull() string   { return i.titleFull }
func (i item) TimeStamp() string   { return i.timeStamp }
func (i item) Description() string { return i.description }
func (i item) FilePath() string    { return i.filePath }
func (i item) FilterValue() string { return i.title }

// default keybind definitions
type keyMap struct {
	filter       key.Binding
	quit         key.Binding
	more         key.Binding
	choose       key.Binding
	remove       key.Binding
	togglePin    key.Binding
	togglePinned key.Binding
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

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
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

func NewModel() model {
	var (
		listKeys   = newKeyMap()
		filterKeys = newFilterKeymap()
	)

	// Make initial list of items
	clipboardItems := config.GetHistory()
	entryItems := filterItems(clipboardItems, false)

	m := model{
		keys:         listKeys,
		filterKeys:   filterKeys,
		help:         help.New(),
		togglePinned: false,
		showFullHelp: false,
	}
	del := m.newItemDelegate(listKeys)
	clipboardList := list.New(entryItems, del, 0, 0)
	clipboardList.Title = "Clipboard History" // set hardcoded title
	clipboardList.SetShowHelp(false)          // override with custom
	clipboardList.Styles.PaginationStyle = lipgloss.NewStyle().
		MarginBottom(1).MarginLeft(2) // set custom pagination spacing

	if len(clipboardItems) < 1 {
		clipboardList.SetShowStatusBar(false) // remove duplicate "No items"
	}

	ct := config.GetTheme()
	if !ct.UseCustom {
		m.list = setDefaultStyling(clipboardList)
		return m
	}

	statusMessageStyle = styledStatusMessage(ct)
	m.help = styledHelp(m.help, ct)
	clipboardList.SetDelegate(styledDelegate(del, ct))
	m.list = styledList(clipboardList, ct)
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	/* this is where the base logic is held for what action to take from
	   the predefined key bindings
	*/
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		if key.Matches(msg, m.keys.more) {
			m.list.SetShowHelp(!m.list.ShowHelp())
			m.updatePaginator()
		}
		switch msg.String() {
		case "p":
			m.togglePinUpdate()
		}
	}

	// this will also call our delegate's update function
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) updatePaginator() {
	pagStyle := lipgloss.NewStyle().MarginBottom(1).MarginLeft(2)
	if m.list.ShowHelp() {
		pagStyle = lipgloss.NewStyle().MarginBottom(0).MarginLeft(2)
	}
	m.list.Styles.PaginationStyle = pagStyle
}

func (m model) View() string {
	listView := m.list.View()
	helpView := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(m.keys))
	render := lipgloss.NewStyle().PaddingLeft(1).Render

	if m.list.FilterState() == list.Filtering {
		return render(listView + "\n" + lipgloss.NewStyle().PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.filterKeys.FilterHelp())))
	}
	if m.list.ShowHelp() {
		return render(listView) // default full view used as replacement
	}
	return render(listView + "\n" + helpView)
}

func filterItems(clipboardItems []config.ClipboardItem, isPinned bool) []list.Item {
	var filteredItems []list.Item

	for _, entry := range clipboardItems {
		shortenedVal := utils.Shorten(entry.Value)
		item := item{
			title:       shortenedVal,
			titleFull:   entry.Value,
			description: "Date copied: " + entry.Recorded,
			filePath:    entry.FilePath,
			pinned:      entry.Pinned,
			timeStamp:   entry.Recorded,
		}

		if entry.Pinned {
			item.description = fmt.Sprintf("Date copied: %s %s", entry.Recorded, pinnedStyle())
		}

		if !isPinned || entry.Pinned {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

// This updates the TUI when an item is pinned/unpinned
func (m *model) togglePinUpdate() {
	index := m.list.Index()
	i, ok := m.list.SelectedItem().(item)
	if !ok {
		return
	}

	i.description = fmt.Sprintf("Date copied: %s", i.timeStamp)
	if !i.pinned {
		i.description = fmt.Sprintf("Date copied: %s %s", i.timeStamp, pinnedStyle())
	}

	i.pinned = !i.pinned
	m.list.SetItem(index, i)

}
