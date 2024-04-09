package config

const (
	baseDir            = ".config"
	defaultHistoryFile = "clipboard_history.json"
	defaultThemeFile   = "custom_theme.json"
	configFile         = "config.json"
	clipseDirName      = "clipse"
	tmpDir             = "tmp_files"
	listenCmd          = "--listen-shell"
	defaultMaxHist     = 100
	maxChar            = 65
)

// Because Go does not support constant Structs :(
func defaultConfig() Config {
	return Config {
		Sources: []string {""},
		MaxHistory: defaultMaxHist,
		HistoryFile: defaultHistoryFile,
	}
}
