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

var (
	appStyle           = lipgloss.NewStyle().Padding(1, 2) // default padding
	statusMessageStyle = lipgloss.NewStyle().Foreground(
		lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"},
	).Render // default styling func used to render message
)

type model struct {
	list          list.Model         // list items
	keys          *keyMap            // keybindings
	filterKeys    *filterKeyMap      // keybindings for filter view
	help          help.Model         // custom help menu
	togglePinned  bool               // pinned view indicator
	theme         config.CustomTheme // colors scheme to uses
	prevDirection string             // prev direction used to track selections
}

type item struct {
	title           string // display title in list
	titleBase       string // unstyled string used for rendering
	titleFull       string // full value stored in history file
	timeStamp       string // local date and time of copy event
	description     string // displayed description in list
	descriptionBase string // unstyled string used for rendering
	filePath        string // "path/to/file" | "null"
	pinned          bool   // pinned status
	selected        bool   // selected status
}

func (i item) Title() string       { return i.title }
func (i item) TitleFull() string   { return i.titleFull }
func (i item) TimeStamp() string   { return i.timeStamp }
func (i item) Description() string { return i.description }
func (i item) FilePath() string    { return i.filePath }
func (i item) FilterValue() string { return i.title }

func NewModel() model {
	var (
		listKeys   = newKeyMap()
		filterKeys = newFilterKeymap()
	)

	clipboardItems := config.GetHistory()

	ct := config.GetTheme()

	m := model{
		keys:          listKeys,
		filterKeys:    filterKeys,
		help:          help.New(),
		togglePinned:  false,
		theme:         ct,
		prevDirection: "",
	}

	entryItems := filterItems(clipboardItems, false, m.theme)

	del := m.newItemDelegate()

	clipboardList := list.New(entryItems, del, 0, 0)
	clipboardList.Title = clipboardTitle // set hardcoded title
	clipboardList.SetShowTitle(true)
	clipboardList.SetShowStatusBar(true)
	clipboardList.SetShowHelp(false)                                           // override with custom
	clipboardList.Styles.PaginationStyle = style.MarginBottom(1).MarginLeft(2) // set custom pagination spacing
	//clipboardList.StatusMessageLifetime = time.Second // can override this if necessary
	clipboardList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.selectDown,
			listKeys.selectSingle,
			listKeys.clearSelected,
		}
	}

	if len(clipboardItems) < 1 {
		clipboardList.SetShowStatusBar(false) // remove duplicate "No items"
	}

	if !ct.UseCustom {
		m.list = setDefaultStyling(clipboardList)
		return m
	}

	statusMessageStyle = styledStatusMessage(ct)
	m.help = styledHelp(m.help, ct)
	m.list = styledList(clipboardList, ct)
	return m
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// if isPinned is true, returns only an array of pinned items, otherwise all
func filterItems(clipboardItems []config.ClipboardItem, isPinned bool, theme config.CustomTheme) []list.Item {
	var filteredItems []list.Item

	for _, entry := range clipboardItems {
		shortenedVal := utils.Shorten(entry.Value)
		item := item{
			title:           shortenedVal,
			titleBase:       shortenedVal,
			titleFull:       entry.Value,
			description:     "Date copied: " + entry.Recorded,
			descriptionBase: "Date copied: " + entry.Recorded,
			filePath:        entry.FilePath,
			pinned:          entry.Pinned,
			timeStamp:       entry.Recorded,
			selected:        false,
		}

		if entry.Pinned {
			item.description = fmt.Sprintf("Date copied: %s %s", entry.Recorded, styledPin(theme))
		}

		if !isPinned || entry.Pinned {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}
