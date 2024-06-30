package config

const (
	defaultHistoryFile = "clipboard_history.json"
	defaultThemeFile   = "custom_theme.json"
	defaultLogFile     = "clipse.log"
	configFile         = "config.json"
	clipseDir          = "clipse"
	defaultTempDir     = "tmp_files"
	listenCmd          = "--listen-shell"
	defaultMaxHist     = 100
	maxChar            = 65
)

// Because Go does not support constant Structs :(
func defaultConfig() Config {
	return Config{
		HistoryFilePath: defaultHistoryFile,
		MaxHistory:      defaultMaxHist,
		TempDirPath:     defaultTempDir,
		LogFilePath:     defaultLogFile,
		ThemeFilePath:   defaultThemeFile,
	}
}
