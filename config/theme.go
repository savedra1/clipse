package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type CustomTheme struct {
	UseCustom          bool   `json:"useCustomTheme"`
	DimmedDesc         string `json:"DimmedDesc"`
	DimmedTitle        string `json:"DimmedTitle"`
	FilteredMatch      string `json:"FilteredMatch"`
	NormalDesc         string `json:"NormalDesc"`
	NormalTitle        string `json:"NormalTitle"`
	SelectedDesc       string `json:"SelectedDesc"`
	SelectedTitle      string `json:"SelectedTitle"`
	SelectedBorder     string `json:"SelectedBorder"`
	SelectedDescBorder string `json:"SelectedDescBorder"`
	TitleFore          string `json:"TitleFore"`
	TitleBack          string `json:"Titleback"`
	StatusMsg          string `json:"StatusMsg"`
	PinIndicatorColor  string `json:"PinIndicatorColor"`
}

func GetTheme() CustomTheme {
	_, err := os.Stat(ClipseConfig.ThemeFilePath)
	if os.IsNotExist(err) {
		initDefaultTheme()
	}

	file, err := os.OpenFile(ClipseConfig.ThemeFilePath, os.O_RDONLY, 0644)
	if err != nil {
		file.Close()
	}

	var theme CustomTheme
	if err := json.NewDecoder(file).Decode(&theme); err != nil {
		fmt.Println("Error decoding JSON for custom_theme.json. Try creating this file manually instead. Err:", err)
		// handleError(err) // No need to terminate program here as default theme can be kept
	}

	// Extract clipboard history items
	return theme
}

func initDefaultTheme() error {
	/*
	  Creates custom_theme.json file is not found in path
	  and sets base config.
	*/
	_, err := os.Stat(ClipseConfig.ThemeFilePath)
	if os.IsNotExist(err) {

		baseConfig := CustomTheme{
			UseCustom:          false,
			DimmedDesc:         "#ffffff",
			DimmedTitle:        "#ffffff",
			FilteredMatch:      "#ffffff",
			NormalDesc:         "#ffffff",
			NormalTitle:        "#ffffff",
			SelectedDesc:       "#ffffff",
			SelectedTitle:      "#ffffff",
			SelectedBorder:     "#ffffff",
			SelectedDescBorder: "#ffffff",
			TitleFore:          "#ffffff",
			TitleBack:          "#434C5E",
			StatusMsg:          "#ffffff",
			PinIndicatorColor:  "#ff0000",
		}

		jsonData, err := json.MarshalIndent(baseConfig, "", "    ")
		if err != nil {
			return err
		}

		err = os.WriteFile(ClipseConfig.ThemeFilePath, jsonData, 0644)
		if err != nil {
			return err
		}

		return nil

	}
	return nil
}
