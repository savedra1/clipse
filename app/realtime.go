package app

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/config"
)

type ReRender struct{}

func (m Model) ListenRealTime(p *tea.Program) {
	historyPath := config.ClipseConfig.HistoryFilePath
	info, _ := os.Stat(historyPath)
	m.lastUpdated = info.ModTime()

	rr := ReRender{}
	var currModTime time.Time
	for {
		historyFileInfo, _ := os.Stat(historyPath)
		currModTime = historyFileInfo.ModTime()

		if currModTime.After(m.lastUpdated) {
			m.lastUpdated = currModTime
			p.Send(rr)
		}
	}
}
