package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

type Config struct {
	AllowDuplicates   bool              `json:"allowDuplicates"`
	HistoryFilePath   string            `json:"historyFile"`
	MaxHistory        int               `json:"maxHistory"`
	DeleteAfter       int               `json:"deleteAfter"`
	LogFilePath       string            `json:"logFile"`
	PollInterval      int               `json:"pollInterval"`
	MaxEntryLength    int               `json:"maxEntryLength"`
	ThemeFilePath     string            `json:"themeFile"`
	TempDirPath       string            `json:"tempDir"`
	KeyBindings       map[string]string `json:"keyBindings"`
	ImageDisplay      ImageDisplay      `json:"imageDisplay"`
	ExcludedApps      []string          `json:"excludedApps"`
	AutoPaste         AutoPaste         `json:"autoPaste"`
	EnableMouse       bool              `json:"enableMouse"`
	EnableDescription bool              `json:"enableDescription"`
	Search            SearchConfig      `json:"search"`
}

type SearchConfig struct {
	Engine          string       `json:"engine"`
	Algo            string       `json:"algo"`
	MatchMode       string       `json:"matchMode"`
	CaseSensitivity string       `json:"caseSensitivity"`
	Normalize       bool         `json:"normalize"`
	Tiebreak        TiebreakList `json:"tiebreak"`
}

type TiebreakEntry struct {
	Key    string `json:"key"`
	Bucket string `json:"bucket,omitempty"`
}

type TiebreakList []TiebreakEntry

func (t *TiebreakList) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	out := make(TiebreakList, 0, len(raw))
	for _, r := range raw {
		trimmed := len(r) > 0 && r[0] == '"'
		if trimmed {
			var s string
			if err := json.Unmarshal(r, &s); err != nil {
				return err
			}
			out = append(out, TiebreakEntry{Key: s})
			continue
		}
		var e TiebreakEntry
		if err := json.Unmarshal(r, &e); err != nil {
			return err
		}
		out = append(out, e)
	}
	*t = out
	return nil
}

func (t TiebreakList) MarshalJSON() ([]byte, error) {
	raw := make([]interface{}, len(t))
	for i, e := range t {
		if e.Bucket == "" {
			raw[i] = e.Key
		} else {
			raw[i] = e
		}
	}
	return json.Marshal(raw)
}

type AutoPaste struct {
	Enabled bool   `json:"enabled"`
	Keybind string `json:"keybind"`
	Buffer  int    `json:"buffer"`
}

type ImageDisplay struct {
	Type      string `json:"type"`
	ScaleX    int    `json:"scaleX"`
	ScaleY    int    `json:"scaleY"`
	HeightCut int    `json:"heightCut"`
}

// Global config object, accessed and used when any configuration is needed.
var ClipseConfig = defaultConfig()

func Init() error {
	/* Ensure $HOME/.config/clipse/clipboard_history.json OR $XDG_CONFIG_HOME
	exists and create the path if not.
	*/

	// returns $HOME/.config || $XDG_CONFIG_HOME
	userHome, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to read home dir.\nerror: %s", err)
	}

	// Construct the path to the config directory
	clipseDir := filepath.Join(userHome, defaultClipseDir)    // the ~/.config/clipse dir
	configPath := filepath.Join(clipseDir, defaultConfigFile) // the path to the config.json file

	// Does Config dir exist, if no make it.
	_, err = os.Stat(clipseDir)
	if os.IsNotExist(err) {
		utils.HandleError(os.MkdirAll(clipseDir, 0755))
	}

	// load the config from file into ClipseConfig struct
	loadConfig(configPath)

	// The history path is absolute at this point. Create it if it does not exist
	utils.HandleError(initHistoryFile())

	// Create TempDir for images if it does not exist.
	_, err = os.Stat(ClipseConfig.TempDirPath)
	if os.IsNotExist(err) {
		utils.HandleError(os.MkdirAll(ClipseConfig.TempDirPath, 0755))
	}

	return nil
}

func loadConfig(configPath string) {
	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		baseConfig := defaultConfig()
		jsonData, err := json.MarshalIndent(baseConfig, "", "    ")
		utils.HandleError(err)
		utils.HandleError(os.WriteFile(configPath, jsonData, 0644))
	}

	configDir := filepath.Dir(configPath)
	confData, err := os.ReadFile(configPath)
	utils.HandleError(err)

	if err = json.Unmarshal(confData, &ClipseConfig); err != nil {
		fmt.Println("Failed to read config. Skipping.\nErr: %w", err)
		utils.LogERROR(fmt.Sprintf("failed to read config. Skipping.\nsrr: %s", err))
	}

	// Expand HistoryFile, ThemeFile, LogFile and TempDir paths
	ClipseConfig.HistoryFilePath = utils.ExpandRel(utils.ExpandHome(ClipseConfig.HistoryFilePath), configDir)
	ClipseConfig.TempDirPath = utils.ExpandRel(utils.ExpandHome(ClipseConfig.TempDirPath), configDir)
	ClipseConfig.ThemeFilePath = utils.ExpandRel(utils.ExpandHome(ClipseConfig.ThemeFilePath), configDir)
	ClipseConfig.LogFilePath = utils.ExpandRel(utils.ExpandHome(ClipseConfig.LogFilePath), configDir)
}
