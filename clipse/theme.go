package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
}

func getTheme() CustomTheme {
	/* returns the clipboardHistory array from the
	clipboard_history.json file
	*/
	_, configDir := getFullPath()
	fp := filepath.Join(configDir, themeFile)

	file, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		file.Close()
	}

	var theme CustomTheme
	if err := json.NewDecoder(file).Decode(&theme); err != nil {
		fmt.Println("Error decoding JSON:", err)
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
