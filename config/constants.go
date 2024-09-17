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
	defaultImageDisplay    = "basic"
	defaultImageScaleX     = 9
	defaultImageScaleY     = 9
	defaultHeightCut       = 2
	listenCmd              = "--listen-shell"
	maxChar                = 65
)

// Initialize default key bindings
func defaultKeyBindings() map[string]string {
	return map[string]string{
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
		"yankFilter":    "ctrl+s",
		"up":            "up",
		"down":          "down",
		"nextPage":      "right",
		"prevPage":      "left",
		"home":          "home",
		"end":           "end",
	}
}

// Because Go does not support constant Structs :(
func defaultConfig() Config {
	return Config{
		HistoryFilePath: defaultHistoryFile,
		MaxHistory:      defaultMaxHist,
		AllowDuplicates: defaultAllowDuplicates,
		TempDirPath:     defaultTempDir,
		LogFilePath:     defaultLogFile,
		ThemeFilePath:   defaultThemeFile,
		KeyBindings:     defaultKeyBindings(),
		ImageDisplay:    defaultImageDisplay,
		ImageScaleX:     defaultImageScaleX,
		ImageScaleY:     defaultImageScaleY,
		HeightCut:       defaultHeightCut,
	}
}
