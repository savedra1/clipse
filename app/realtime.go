package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"

	"github.com/savedra1/clipse/config"
	"github.com/savedra1/clipse/utils"
)

type ReRender struct{}

func (m Model) ListenRealTime(p *tea.Program) {
	historyPath := config.ClipseConfig.HistoryFilePath

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.LogERROR("Could not create file watcher: " + err.Error())
		return
	}
	defer watcher.Close()

	if err := watcher.Add(historyPath); err != nil {
		utils.LogERROR("Could not watch history file: " + err.Error())
		return
	}

	rr := ReRender{}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				p.Send(rr)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			utils.LogERROR("File watcher error: " + err.Error())
		}
	}
}
