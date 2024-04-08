package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/savedra1/clipse/utils"
)

type Config struct {
	SourcePaths []string `json:"sourcePaths"`
	MaxHist     int      `json:"maxList"`
	HistoryFile string   `json:"historyFile"`
}

// Global config object, accessed and used when any configuration is needed.
var ClipseConfig = Config {}

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
	
	confData, err := os.ReadFile(configPath)
	if err = json.Unmarshal(confData, &ClipseConfig); err != nil {
		fmt.Println("Failed to read config. Fallback to default.\nErr: %w", err)
		ClipseConfig = defaultConfig()
	} else {
		fmt.Println("CONFIG SUCCESSFULLY LOADED!")
	}
}
