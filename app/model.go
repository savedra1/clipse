package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

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
	title       string
	titleFull   string
	timeStamp   string
	description string
	filePath    string
	pinned      bool
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

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
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
