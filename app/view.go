package app

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	render := lipgloss.NewStyle().PaddingLeft(1).Render

	listView := m.list.View()
	helpView := lipgloss.NewStyle().PaddingLeft(2).Render(m.help.View(m.keys))

	switch {
	case m.showConfirmation:
		listView = m.confirmationList.View()
		helpView = lipgloss.NewStyle().PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.confirmationKeys.ConfirmationHelp()))
		return render(listView + "\n" + helpView)

	case m.list.SettingFilter():
		return render(listView + "\n" + lipgloss.NewStyle().PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.filterKeys.FilterHelp())),
		)

	case m.list.ShowHelp():
		return render(listView)

	default:
		return render(listView + "\n" + helpView)
	}
}
