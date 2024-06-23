package app

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/config"
)

// Delegate used to override individual item appearance based on state.

type itemDelegate struct {
	theme config.CustomTheme
}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	listItem list.Item,
) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	var renderStr string

	switch {
	case m.FilterState() == list.Filtering:
		renderStr = d.itemFilterStyle(i, m)
	case i.selected, index == m.Index():
		renderStr = d.itemSelectedStyle(i, m, index)
	default:
		renderStr = d.itemNormalStyle(i)
	}

	fmt.Fprint(w, renderStr)
}

func (m *model) newItemDelegate() itemDelegate {
	return itemDelegate{
		theme: m.theme,
	}
}
