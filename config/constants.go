package config

const (
	configFile             = "config.json"
	clipseDir              = "clipse"
	defaultAllowDuplicates = false
	defaultHistoryFile     = "clipboard_history.json"
	defaultMaxHist         = 100
	defaultDeleteAfter     = 0
	defaultLogFile         = "clipse.log"
	defaultTempDir         = "tmp_files"
	defaultThemeFile       = "custom_theme.json"
	listenShellCmd         = "--listen-shell"
	listenCmd              = "-listen"
	maxChar                = 65
)

// Initialize default key bindings
func defaultKeyBindings() map[string]string {
	return map[string]string{
		"filter":        "/",
		"quit":          "esc",
		"more":          "?",
		"choose":        "enter",
		"remove":        "backspace",
		"togglePin":     "p",
		"togglePinned":  "tab",
		"preview":       " ",
		"selectDown":    "shift+down",
		"selectUp":      "shift+up",
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

// Default list of applications to exclude from clipboard history
func defaultExcludedApps() []string {
	return []string{
		"1Password",
		"Bitwarden",
		"KeePassXC",
		"LastPass",
		"Dashlane",
		"Password Safe",
		"Keychain Access",
	}
}

// Because Go does not support constant Structs :(
func defaultConfig() Config {
	return Config{
		HistoryFilePath: defaultHistoryFile,
		MaxHistory:      defaultMaxHist,
		DeleteAfter:     defaultDeleteAfter,
		AllowDuplicates: defaultAllowDuplicates,
		TempDirPath:     defaultTempDir,
		LogFilePath:     defaultLogFile,
		ThemeFilePath:   defaultThemeFile,
		KeyBindings:     defaultKeyBindings(),
		ExcludedApps:    defaultExcludedApps(),
		ImageDisplay: ImageDisplay{
			Type:      "basic",
			ScaleX:    9,
			ScaleY:    9,
			HeightCut: 2,
		},
	}
}
