package app

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

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
