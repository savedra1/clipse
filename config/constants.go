package config

const (
	defaultConfigFile        = "config.json"
	defaultClipseDir         = "clipse"
	defaultAllowDuplicates   = false
	defaultHistoryFile       = "clipboard_history.json"
	defaultMaxHist           = 100
	defaultDeleteAfter       = 0
	defaultLogFile           = "clipse.log"
	defaultPollInterval      = 50
	defaultMaxEntryLength    = 65
	defaultTempDir           = "tmp_files"
	defaultThemeFile         = "custom_theme.json"
	defaultEnableAutoPaste   = false
	defaultAutoPasteKeyBind  = "ctrl+v"
	defaultAutoPasteBuffer   = 10
	defaultEnableMouse       = true
	defaultEnableDescription = true
)

// Initialize default key bindings
func defaultKeyBindings() map[string]string {
	return map[string]string{
		"filter":        "/",
		"quit":          "esc",
		"forceQuit":     "Q",
		"more":          "?",
		"choose":        "enter",
		"remove":        "backspace",
		"togglePin":     "p",
		"togglePinned":  "tab",
		"preview":       "space",
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
		HistoryFilePath:   defaultHistoryFile,
		MaxHistory:        defaultMaxHist,
		DeleteAfter:       defaultDeleteAfter,
		AllowDuplicates:   defaultAllowDuplicates,
		TempDirPath:       defaultTempDir,
		LogFilePath:       defaultLogFile,
		PollInterval:      defaultPollInterval,
		MaxEntryLength:    defaultMaxEntryLength,
		ThemeFilePath:     defaultThemeFile,
		KeyBindings:       defaultKeyBindings(),
		ExcludedApps:      defaultExcludedApps(),
		EnableMouse:       defaultEnableMouse,
		EnableDescription: defaultEnableDescription,
		ImageDisplay: ImageDisplay{
			Type:      "basic",
			ScaleX:    9,
			ScaleY:    9,
			HeightCut: 2,
		},
		AutoPaste: AutoPaste{
			Enabled: defaultEnableAutoPaste,
			Keybind: defaultAutoPasteKeyBind,
			Buffer:  defaultAutoPasteBuffer,
		},
	}
}
