package app

import (
	"fmt"
	"strings"

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
	list         list.Model // list items
	keys         *keyMap    // keybindings
	togglePinned bool       // pinned indicator
	showFullHelp bool
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

type keyMap struct {
	// default keybind definitions
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

// ShortHelp returns the key bindings for the short help screen.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.choose, k.remove, k.filter, k.togglePin, k.togglePinned, k.more,
	}
}

// FullHelp returns the key bindings for the full help screen.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.choose, k.remove, k.togglePin, k.togglePinned},
		{k.filter, k.quit},
	}
}

func (m model) Init() tea.Cmd { // initialize app
	return tea.EnterAltScreen
}

func NewModel() model {
	var listKeys = newKeyMap()

	// Make initial list of items
	clipboardItems := config.GetHistory()
	entryItems := filterItemsByPinned(clipboardItems, false)

	// Setup list
	m := model{}
	del := m.newItemDelegate(listKeys)
	ct := config.GetTheme()
	if ct.UseCustom {
		del = styledDelegate(del, ct)
		statusMessageStyle = styledStatusMessage(ct)
	}

	clipboardList := list.New(entryItems, del, 0, 0)
	clipboardList.Title = "Clipboard History"
	clipboardList.SetShowHelp(false) // override with only custom options

	if ct.UseCustom { // add additional customizations after delegate created
		clipboardList = styledList(clipboardList, ct)
	}

	if len(clipboardItems) < 1 {
		clipboardList.SetShowStatusBar(false)
	}

	return model{
		list:         clipboardList,
		keys:         listKeys,
		togglePinned: false,
		showFullHelp: false,
	}
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
			m.showFullHelp = !m.showFullHelp
			m.list.SetShowHelp(!m.list.ShowHelp())
		}
		switch msg.String() {
		case "p":
			m.togglePinUpdate()
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	listView := m.list.View()
	parts := strings.Split(listView, "\n")

	var helpView string

	if !m.showFullHelp {
		helpView = "  " + m.list.Help.ShortHelpView(m.keys.ShortHelp())
	}

	// Find the index of the pagination line
	paginationIndex := -1
	for i, line := range parts {
		if strings.Contains(line, "items") {
			paginationIndex = i
			break
		}
	}

	// Remove the line before the pagination line
	if paginationIndex > 0 {
		parts = append(parts[:paginationIndex-1], parts[paginationIndex:]...)
	}

	parts = append(parts, helpView)

	return appStyle.Render(strings.Join(parts, "\n"))
}

func filterItemsByPinned(clipboardItems []config.ClipboardItem, isPinned bool) []list.Item {
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

func pinnedStyle() string {
	color := "#FF0000"
	pinChar := " "
	config := config.GetTheme()

	if config.UseCustom {
		color = config.PinIndicatorColor
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).SetString(pinChar).Render()
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
