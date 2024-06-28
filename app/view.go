package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	render := lipgloss.NewStyle().PaddingLeft(1).Render

	listView := m.list.View()
	helpView := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(m.keys))

	if m.showConfirmation {
		listView = m.confirmationList.View()
		helpView = lipgloss.NewStyle().PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.confirmationKeys.ConfirmationHelp()))
		return render(listView + strings.Repeat("\n", 22) + helpView)
	}

	if m.list.SettingFilter() {
		return render(listView + "\n" + lipgloss.NewStyle().PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.filterKeys.FilterHelp())),
		)
	}

	if m.list.ShowHelp() {
		return render(listView)
	}
	return render(listView + "\n" + helpView)
}
