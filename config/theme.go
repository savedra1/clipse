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

// For now, reload each time the window is opened. Near future, on file change.
var themePaths []string

func GetTheme() CustomTheme {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	// Just choose the first theme in the list. Change to allow selecting
	// from multiple themes in the future maybe.
	fp := themePaths[0]

	file, err := os.OpenFile(fp, os.O_RDONLY, 0644)
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

func initTheme(fp string) error {
	/*
	  Creates custom_theme.json file is not found in path
	  and sets base config.
	*/
	_, err := os.Stat(fp)
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

		err = os.WriteFile(fp, jsonData, 0644)
		if err != nil {
			return err
		}

		return nil

	}
	return nil
}
