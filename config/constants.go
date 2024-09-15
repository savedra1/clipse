package config

const (
	configFile             = "config.json"
	clipseDir              = "clipse"
	defaultAllowDuplicates = false
	defaultHistoryFile     = "clipboard_history.json"
	defaultMaxHist         = 100
	defaultLogFile         = "clipse.log"
	defaultTempDir         = "tmp_files"
	defaultThemeFile       = "custom_theme.json"
	listenCmd              = "--listen-shell"
	maxChar                = 65
)

// Because Go does not support constant Structs :(
func defaultConfig() Config {
	return Config{
		AllowDuplicates: true,
		HistoryFilePath: "clipboard_history.json",
		MaxHistory:      1000,
		LogFilePath:     "log.txt",
		ThemeFilePath:   "theme.json",
		TempDirPath:     "temp",
		KeyBindings: map[string]string{
			"filter":        "/",
			"quit":          "q",
			"more":          "?",
			"choose":        "enter",
			"remove":        "x",
			"togglePin":     "p",
			"togglePinned":  "tab",
			"preview":       " ",
			"selectDown":    "ctrl+down",
			"selectUp":      "ctrl+up",
			"selectSingle":  "s",
			"clearSelected": "S",
			"fuzzySelect":   "F",
			"yankFilter":    "ctrl+s",
			"up":            "up",
			"down":          "down",
			"nextPage":      "right",
			"prevPage":      "left",
			"home":          "home",
			"end":           "end",
		},
	}
}
