package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

type Model struct {
	list             list.Model          // list items
	keys             *keyMap             // keybindings
	filterKeys       *filterKeyMap       // keybindings for filter view
	confirmationKeys *confirmationKeyMap // keybindings for the confirmation view
	help             help.Model          // custom help menu
	togglePinned     bool                // pinned view indicator
	theme            config.CustomTheme  // colors scheme to uses
	prevDirection    string              // prev direction used to track selections
	confirmationList list.Model          // secondary list Model used for confirmation screen
	showConfirmation bool                // whether to show confirmation screen
	itemCache        []SelectedItem      // easy access for related items following confirmation screen
	preview          viewport.Model      // viewport model used for displaying previews
	originalHeight   int                 // for restore height of preview viewport in sixel mode
	previewReady     bool                // viewport needs to wait for the initial window size message
	showPreview      bool                // whether the viewport preview should be displayed
	previewKeys      *previewKeymap      // keybindings for the viewport model
	lastUpdated      time.Time
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

type SelectedItem struct {
	Index     int    // list index needed for deletion
	TimeStamp string // timestamp needed for deletion
	Value     string // full val needed for copy
	Pinned    bool   // pinned val needed to determine whether confirmation screen is needed
}

func (i item) Title() string       { return i.title }
func (i item) TitleFull() string   { return i.titleFull }
func (i item) TimeStamp() string   { return i.timeStamp }
func (i item) Description() string { return i.description }
func (i item) FilePath() string    { return i.filePath }
func (i item) FilterValue() string { return i.title }

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func NewModel() Model {
	var (
		keyConfig        = config.ClipseConfig.KeyBindings
		listKeys         = newKeyMap(keyConfig)
		filterKeys       = newFilterKeymap(keyConfig)
		confirmationKeys = newConfirmationKeymap(keyConfig)
	)

	clipboardItems := config.GetHistory()

	theme := config.GetTheme()

	m := Model{
		keys:             listKeys,
		filterKeys:       filterKeys,
		confirmationKeys: confirmationKeys,
		help:             help.New(),
		togglePinned:     false,
		theme:            theme,
		prevDirection:    "",
		showConfirmation: false,
		preview:          NewPreview(),
		showPreview:      false,
		previewKeys:      newPreviewKeyMap(),
	}

	entryItems := filterItems(clipboardItems, false, m.theme)

	del := m.newItemDelegate()

	clipboardList := list.New(entryItems, del, 0, 0)
	clipboardList.KeyMap = defaultOverrides(config.ClipseConfig.KeyBindings)   // override default list keys with custom values
	clipboardList.Title = clipboardTitle                                       // set hardcoded title
	clipboardList.SetShowHelp(false)                                           // override with custom
	clipboardList.Styles.PaginationStyle = style.MarginBottom(1).MarginLeft(2) // set custom pagination spacing
	//clipboardList.StatusMessageLifetime = time.Second // can override this if necessary
	clipboardList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.preview,
			listKeys.selectDown,
			listKeys.selectSingle,
			listKeys.clearSelected,
		}
	}

	confirmationList := newConfirmationList(del)

	if len(clipboardItems) < 1 {
		clipboardList.SetShowStatusBar(false) // remove duplicate "No items"
	}

	statusMessageStyle = styledStatusMessage(theme)
	m.help = styledHelp(m.help, theme)
	m.list = styledList(clipboardList, theme)
	m.confirmationList = styledList(confirmationList, theme)
	m.enableConfirmationKeys(false)

	return m
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

func newConfirmationList(del itemDelegate) list.Model {
	items := confirmationItems()
	l := list.New(items, del, 0, 10)
	l.Title = confirmationTitle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.KeyMap.Quit.SetEnabled(false)
	return l
}

func confirmationItems() []list.Item {
	return []list.Item{
		item{
			title:           "No",
			titleBase:       "No",
			descriptionBase: "go back",
		},
		item{
			title:           "Yes",
			titleBase:       "Yes",
			descriptionBase: "delete the item(s)",
		},
	}
}
