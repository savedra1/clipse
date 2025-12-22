package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/config"
)

func (m *Model) newItemDelegate() itemDelegate {
	return itemDelegate{
		theme: m.theme,
	}
}

type itemDelegate struct {
	theme config.CustomTheme
}

func (d itemDelegate) Height() int {
	if !config.ClipseConfig.EnableDescription {
		return 1
	}
	return 2
}

func (d itemDelegate) Spacing() int {
	if !config.ClipseConfig.EnableDescription {
		return 2 // extra space needed when no description to keep mouse in sync
	}
	return 1
}

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

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

	case m.SettingFilter():
		if strings.Contains(
			strings.ToLower(i.titleFull),
			strings.ToLower(m.FilterValue()),
		) && m.FilterValue() != "" {
			renderStr = d.itemSelectedStyle(i)
		} else {
			renderStr = d.itemFilterStyle(i)
		}

	case index == m.Index():
		renderStr = d.itemChosenStyle(i)

	case i.selected:
		renderStr = d.itemSelectedStyle(i)

	default:
		renderStr = d.itemNormalStyle(i)
	}

	fmt.Fprint(w, renderStr)
}
