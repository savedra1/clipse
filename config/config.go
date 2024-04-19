package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/savedra1/clipse/utils"
)

type Config struct {
	Sources         []string `json:"sources"`
	MaxHistory      int      `json:"maxHistory"`
	HistoryFilePath string   `json:"historyFile"`
	TempDirPath     string   `json:"tempDir"`
}

type source struct {
	SourceType string `json:"sourceType"`
}

// Global config object, accessed and used when any configuration is needed.
var ClipseConfig = defaultConfig()

func configInit(path string) {
	loadConfig(path)
}

func loadConfig(configPath string) {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		baseConfig := defaultConfig()

		jsonData, err := json.MarshalIndent(baseConfig, "", "    ")
		utils.HandleError(err)

		err = os.WriteFile(configPath, jsonData, 0644)
		if err != nil {
			fmt.Println("Failed to create:", configPath)
		}
	}

	// When recursively calling the sources, we want this source to not be
	// overwritten. So we store it in this, and at the end set ClipseConfig.
	//
	// This means that the last instance is the most signigicant.
	var tempConfig Config

	configDir := filepath.Dir(configPath)

	confData, err := os.ReadFile(configPath)
	if err = json.Unmarshal(confData, &tempConfig); err != nil {
		fmt.Println("Failed to read config. Skipping.\nErr: %w", err)
	}

	for i := range tempConfig.Sources {
		// Expand all cases of `~` in source, and call the loadSource func.
		src := &tempConfig.Sources[i]
		*src = utils.ExpandRel(utils.ExpandHome(*src), configDir)

		loadSource(*src)
	}

	// ClipseConfig contains all the settings from all sources. Store sources in temp var.
	tempConfig.Sources = append(tempConfig.Sources, ClipseConfig.Sources...)

	// All other configs have loaded, load this one.
	if err = json.Unmarshal(confData, &ClipseConfig); err != nil {
		fmt.Println("Failed to read config. Skipping.\nErr: %w", err)
	}

	// Recover source files list.
	ClipseConfig.Sources = tempConfig.Sources

	// Expand HistoryFile and TempDir paths
	ClipseConfig.HistoryFilePath = utils.ExpandRel(utils.ExpandHome(ClipseConfig.HistoryFilePath), configDir)
	ClipseConfig.TempDirPath = utils.ExpandRel(utils.ExpandHome(ClipseConfig.TempDirPath), configDir)
}

func loadSource(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Linked source file not found at:", path)
		return
	}

	var src source

	data, err := os.ReadFile(path)
	if err = json.Unmarshal(data, &src); err != nil {
		fmt.Printf("Failed to read source at %s. Incorrectly formatted json!\n", path)
	}

	switch src.SourceType {
	case "config":
		loadConfig(path)
	case "theme":
		themePaths = append(themePaths, path)
	case "":
		fmt.Printf("Error: \"sourceType\" tag not found in source file: %s. File skipped.\n", path)
		fmt.Println("Possible values for sourceType:\n\t- config\n\t- theme\n\t- history")
	default:
		fmt.Printf("Error: Invalid value \"%s\" in \"sourceType\" tag for source file: %s\n",
			src.SourceType, path)
		fmt.Println("Possible values for sourceType:\n\t- config\n\t- theme\n\t- history")
	}
}
