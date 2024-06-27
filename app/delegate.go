package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/savedra1/clipse/config"
)

func (m *model) newItemDelegate() itemDelegate {
	return itemDelegate{
		theme: m.theme,
	}
}

// delegate used to override individual item appearance based on state

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

func (m *model) warningDelegate() secondaryItemDelegate {
	return secondaryItemDelegate{
		theme: m.theme,
	}
}

type secondaryItemDelegate struct {
	theme config.CustomTheme
	task  string
	items []item
}

func (sd secondaryItemDelegate) Height() int  { return 2 }
func (sd secondaryItemDelegate) Spacing() int { return 1 }
func (sd secondaryItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if msg.String() == "down" {
			m.CursorDown()
		}

	}
	return nil
}

func (sd secondaryItemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	listItem list.Item,
) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	renderStr := style.
		Foreground(lipgloss.Color(sd.theme.SelectedTitle)).
		PaddingLeft(1).
		BorderLeft(true).BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(sd.theme.SelectedDescBorder)).
		Render(i.title)

	if index != m.Index() {
		renderStr = style.
			Foreground(lipgloss.Color(sd.theme.NormalTitle)).
			PaddingLeft(1).
			Render(i.title)
	}

	fmt.Fprint(w, renderStr)
}
