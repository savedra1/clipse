package config

import (
	"os/user"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

const (
	baseDir            = ".config"
	defaultHistoryFile = "clipboard_history.json"
	defaultThemeFile   = "custom_theme.json"
	configFile         = "config.json"
	clipseDir          = "clipse"
	defaultTempDir     = "tmp_files"
	listenCmd          = "--listen-shell"
	defaultMaxHist     = 100
	maxChar            = 65
)

// Because Go does not support constant Structs :(
func defaultConfig() Config {
	currentUser, err := user.Current()
	utils.HandleError(err)

	return Config{
		Sources:         []string{filepath.Join(currentUser.HomeDir, baseDir, clipseDir, defaultThemeFile)},
		MaxHistory:      defaultMaxHist,
		HistoryFilePath: defaultHistoryFile,
		TempDirPath:     defaultTempDir,
	}
}
