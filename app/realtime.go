package app

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

type ReRender struct{}

func (m Model) ListenRealTime(p *tea.Program) {
	historyPath := config.ClipseConfig.HistoryFilePath
	info, err := os.Stat(historyPath)
	if err != nil {
		utils.LogERROR("Could not get Modification time of history file, starting real time mode failed")
		return
	}
	m.lastUpdated = info.ModTime()

	rr := ReRender{}
	var currModTime time.Time
	for {
		historyFileInfo, err := os.Stat(historyPath)
		if err != nil {
			continue
		}
		currModTime = historyFileInfo.ModTime()

		if currModTime.After(m.lastUpdated) {
			m.lastUpdated = currModTime
			p.Send(rr)
		}
	}
}
