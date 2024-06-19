package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/savedra1/clipse/config"
)

func setDefaultStyling(clipboardList list.Model) list.Model {
	// align list elements
	clipboardList.FilterInput.PromptStyle = lipgloss.NewStyle().PaddingTop(1)
	clipboardList.Styles.Title = lipgloss.NewStyle().MarginTop(1)
	clipboardList.Styles.StatusBar = lipgloss.NewStyle().MarginBottom(1).MarginLeft(2)
	clipboardList.Styles.DividerDot = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
	clipboardList.Help.FullSeparator = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Render("•")
	clipboardList.Help.ShortSeparator = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Render("•")
	clipboardList.Styles.NoItems = lipgloss.NewStyle().PaddingBottom(1).PaddingLeft(2)
	return clipboardList
}

func styledDelegate(del list.DefaultDelegate, ct config.CustomTheme) list.DefaultDelegate {
	del.Styles.DimmedDesc = del.Styles.DimmedDesc.
		Foreground(lipgloss.Color(ct.DimmedDesc))
	del.Styles.DimmedTitle = del.Styles.DimmedTitle.
		Foreground(lipgloss.Color(ct.DimmedTitle))
	del.Styles.FilterMatch = del.Styles.FilterMatch.
		Foreground(lipgloss.Color(ct.FilteredMatch))
	del.Styles.NormalDesc = del.Styles.NormalDesc.
		Foreground(lipgloss.Color(ct.NormalDesc))
	del.Styles.NormalTitle = del.Styles.NormalTitle.
		Foreground(lipgloss.Color(ct.NormalTitle))
	del.Styles.SelectedDesc = del.Styles.SelectedDesc.
		Foreground(lipgloss.Color(ct.SelectedDesc)).
		BorderForeground(lipgloss.Color(ct.SelectedDescBorder))
	del.Styles.SelectedTitle = del.Styles.SelectedTitle.
		Foreground(lipgloss.Color(ct.SelectedTitle)).
		BorderForeground(lipgloss.Color(ct.SelectedBorder))

	return del
}

func styledList(clipboardList list.Model, ct config.CustomTheme) list.Model {
	clipboardList.FilterInput.PromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.FilterPrompt)).PaddingTop(1)
	clipboardList.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.FilterText))
	clipboardList.Styles.StatusBarFilterCount = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.FilterInfo))
	clipboardList.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.FilterCursor))
	clipboardList.Styles.StatusEmpty = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.FilterInfo))
	clipboardList.Help.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpKey))
	clipboardList.Help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpDesc))
	clipboardList.Help.Styles.FullKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpKey))
	clipboardList.Help.Styles.FullDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpDesc))
	clipboardList.Paginator.ActiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.PageActiveDot)).Render("•")
	clipboardList.Paginator.InactiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.PageInactiveDot)).Render("•")
	clipboardList.Styles.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.TitleInfo)).MarginBottom(1).MarginLeft(2)
	clipboardList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.TitleFore)).Background(lipgloss.Color(ct.TitleBack)).MarginTop(1).
		Align(lipgloss.Position(1))
	clipboardList.Styles.DividerDot = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.DividerDot)).SetString("•").PaddingLeft(1).PaddingRight(1)
	clipboardList.Help.FullSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.DividerDot)).PaddingLeft(1).PaddingRight(1).Render("•")
	clipboardList.Help.ShortSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.DividerDot)).PaddingLeft(1).PaddingRight(1).Render("•")
	clipboardList.Styles.NoItems = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ct.TitleInfo)).PaddingBottom(1).PaddingLeft(2)

	return clipboardList
}

func styledHelp(help help.Model, ct config.CustomTheme) help.Model {
	help.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpKey))
	help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpDesc))
	help.Styles.FullKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpKey))
	help.Styles.FullDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.HelpDesc))
	help.FullSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.DividerDot)).PaddingLeft(1).PaddingRight(1).Render("•")
	help.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color(ct.DividerDot)).PaddingLeft(1).PaddingRight(1).Render("•")
	return help
}

func styledStatusMessage(ct config.CustomTheme) func(strs ...string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: ct.StatusMsg, Dark: ct.StatusMsg}).
		Render
}

func pinnedStyle() string {
	color := "#FF0000"
	pinChar := " "
	config := config.GetTheme()

	if config.UseCustom {
		color = config.PinIndicatorColor
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).SetString(pinChar).Render()
}
