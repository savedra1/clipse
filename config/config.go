package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

type Config struct {
	AllowDuplicates bool              `json:"allowDuplicates"`
	HistoryFilePath string            `json:"historyFile"`
	MaxHistory      int               `json:"maxHistory"`
	DeleteAfter     int               `json:"deleteAfter"`
	LogFilePath     string            `json:"logFile"`
	ThemeFilePath   string            `json:"themeFile"`
	TempDirPath     string            `json:"tempDir"`
	KeyBindings     map[string]string `json:"keyBindings"`
	ImageDisplay    ImageDisplay      `json:"imageDisplay"`
	ExcludedApps    []string          `json:"excludedApps"`
}
type ImageDisplay struct {
	Type      string `json:"type"`
	ScaleX    int    `json:"scaleX"`
	ScaleY    int    `json:"scaleY"`
	HeightCut int    `json:"heightCut"`
}

// Global config object, accessed and used when any configuration is needed.
var ClipseConfig = defaultConfig()

func Init() (string, string, error) {
	/* Ensure $HOME/.config/clipse/clipboard_history.json OR $XDG_CONFIG_HOME
	exists and create the path if not.
	*/

	// returns $HOME/.config || $XDG_CONFIG_HOME
	userHome, err := os.UserConfigDir()
	if err != nil {
		return "", "", fmt.Errorf("failed to read home dir.\nerror: %s", err)
	}

	// Construct the path to the config directory
	clipseDir := filepath.Join(userHome, clipseDir)    // the ~/.config/clipse dir
	configPath := filepath.Join(clipseDir, configFile) // the path to the config.json file

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

	return ClipseConfig.LogFilePath, utils.DisplayServer(), nil
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
