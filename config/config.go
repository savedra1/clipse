package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/savedra1/clipse/utils"
)

var clipseConfig = Config {}

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
	if err = json.Unmarshal(confData, &clipseConfig); err != nil {
		fmt.Println("Failed to read config. Fallback to default.\nErr: %w", err)
		clipseConfig = defaultConfig()
	}
}
