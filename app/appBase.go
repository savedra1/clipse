package app

import (
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* This is where the base level configuration for the
   bubbletea CLI app is defined.

   Base level config includes:
   - Color scheme of text
   - Font
   - Key bindings
   - List sructure
   - Help menu
   - Defualt actions
*/

var (
	// base styling config using lipgloss
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#434C5E")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	// (Each Row in clipboard view)
	title       string
	titleFull   string
	description string
	filePath    string
}

func (i item) Title() string       { return i.title }
func (i item) TitleFull() string   { return i.titleFull }
func (i item) Description() string { return i.description }
func (i item) FilePath() string    { return i.filePath }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	// default keybind definitions
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{

		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	// model pulls all relevant elems together for rendering
	list         list.Model      // list items
	keys         *listKeyMap     // keybindings
	delegateKeys *delegateKeyMap // custom key bindings
}

func NewModel() model {
	// new model needs raising to render additional custom keys
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Make initial list of items
	clipboardItems := config.GetHistory()
	var entryItems []list.Item
	for _, entry := range clipboardItems {
		shortenedVal := utils.Shorten(entry.Value)
		item := item{
			title:       shortenedVal,
			titleFull:   entry.Value,
			description: "Date copied: " + entry.Recorded,
			filePath:    entry.FilePath,
		}
		entryItems = append(entryItems, item)
	}

	// Setup list

	del := newItemDelegate(delegateKeys)
	ct := config.GetTheme()
	if ct.UseCustom {
		del.Styles.DimmedDesc = del.Styles.DimmedDesc.
			Foreground(lipgloss.Color(ct.DimmedDesc))

		del.Styles.DimmedTitle = del.Styles.DimmedTitle.
			Foreground(lipgloss.Color(ct.DimmedTitle))

		del.Styles.FilterMatch = del.Styles.FilterMatch.
			Foreground(lipgloss.Color(ct.FilteredMatch))

		del.Styles.NormalDesc = del.Styles.NormalDesc.
			Foreground(lipgloss.Color(ct.NormalDesc))

		del.Styles.NormalTitle = del.Styles.NormalTitle.
			Foreground(lipgloss.Color(ct.NormalTitle))

		del.Styles.SelectedDesc = del.Styles.SelectedDesc.
			Foreground(lipgloss.Color(ct.SelectedDesc)).
			BorderForeground(lipgloss.Color(ct.SelectedDescBorder))

		del.Styles.SelectedTitle = del.Styles.SelectedTitle.
			Foreground(lipgloss.Color(ct.SelectedTitle)).
			BorderForeground(lipgloss.Color(ct.SelectedBorder))

		titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ct.TitleFore)).
			Background(lipgloss.Color(ct.TitleBack)).
			Padding(0, 1)

		statusMessageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: ct.StatusMsg, Dark: ct.StatusMsg}).
			Render
	}

	//c := lipgloss.Color("#6f03fc")
	//delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(c)

	clipboardList := list.New(entryItems, del, 0, 0)
	clipboardList.Title = "Clipboard History"
	clipboardList.Styles.Title = titleStyle
	clipboardList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         clipboardList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd { // initialise app
	return tea.EnterAltScreen
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
		//txtlog(fmt.Sprintf("%s: key state = %s", time.Now(), m.list.FilterState()))
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string { // Render app in terminal using client libs
	return appStyle.Render(m.list.View())
}
