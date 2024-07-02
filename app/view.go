package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	render := style.PaddingLeft(1).Render

	listView := m.list.View()
	helpView := style.PaddingLeft(2).Render(m.help.View(m.keys))

	switch {

	case m.showPreview:
		helpView = style.PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.previewKeys.PreviewHelp()))
		return fmt.Sprintf(
			"\n%s\n%s\n%s\n%s\n",
			m.previewHeaderView(), m.preview.View(), m.previewFooterView(), helpView,
		)

	case m.showConfirmation:
		listView = m.confirmationList.View()
		helpView = style.PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.confirmationKeys.ConfirmationHelp()))
		return render(listView + "\n" + helpView)

	case m.list.SettingFilter():
		return render(listView + "\n" + style.PaddingLeft(2).Render(
			m.list.Help.ShortHelpView(m.filterKeys.FilterHelp())),
		)

	case m.list.ShowHelp():
		return render(listView)

	default:
		return render(listView + "\n" + helpView)
	}
}

func (m *Model) previewHeaderView() string {
	title := previewTitleStyle.Render(previewHeader)
	line := strings.Repeat(borderMiddleChar, max(0, m.preview.Width-lipgloss.Width(title)))
	return m.styledPreviewHeader(lipgloss.JoinHorizontal(lipgloss.Center, title, line))
}

func (m *Model) previewFooterView() string {
	info := previewInfoStyle.Render(fmt.Sprintf("%3.f%%", m.preview.ScrollPercent()*100))
	line := strings.Repeat(borderMiddleChar, max(0, m.preview.Width-lipgloss.Width(info)))
	return m.styledPreviewFooter(lipgloss.JoinHorizontal(lipgloss.Center, line, info))
}
