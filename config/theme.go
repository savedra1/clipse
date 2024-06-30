package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/savedra1/clipse/utils"
)

type CustomTheme struct {
	UseCustom          bool   `json:"useCustomTheme"`
	TitleFore          string `json:"TitleFore"`
	TitleBack          string `json:"TitleBack"`
	TitleInfo          string `json:"TitleInfo"`
	NormalTitle        string `json:"NormalTitle"`
	DimmedTitle        string `json:"DimmedTitle"`
	SelectedTitle      string `json:"SelectedTitle"`
	NormalDesc         string `json:"NormalDesc"`
	DimmedDesc         string `json:"DimmedDesc"`
	SelectedDesc       string `json:"SelectedDesc"`
	StatusMsg          string `json:"StatusMsg"`
	PinIndicatorColor  string `json:"PinIndicatorColor"`
	SelectedBorder     string `json:"SelectedBorder"`
	SelectedDescBorder string `json:"SelectedDescBorder"`
	FilteredMatch      string `json:"FilteredMatch"`
	FilterPrompt       string `json:"FilterPrompt"`
	FilterInfo         string `json:"FilterInfo"`
	FilterText         string `json:"FilterText"`
	FilterCursor       string `json:"FilterCursor"`
	HelpKey            string `json:"HelpKey"`
	HelpDesc           string `json:"HelpDesc"`
	PageActiveDot      string `json:"PageActiveDot"`
	PageInactiveDot    string `json:"PageInactiveDot"`
	DividerDot         string `json:"DividerDot"`
}

func GetTheme() CustomTheme {
	_, err := os.Stat(ClipseConfig.ThemeFilePath)
	if os.IsNotExist(err) {
		if err = initDefaultTheme(); err != nil {
			utils.LogERROR(fmt.Sprintf("could not initialize theme: %s", err))
		}
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
			TitleFore:          "#ffffff",
			TitleBack:          "#434C5E",
			TitleInfo:          "#ffffff",
			NormalTitle:        "#ffffff",
			DimmedTitle:        "#ffffff",
			SelectedTitle:      "#ffffff",
			NormalDesc:         "#ffffff",
			DimmedDesc:         "#ffffff",
			SelectedDesc:       "#ffffff",
			StatusMsg:          "#ffffff",
			PinIndicatorColor:  "#ff0000",
			SelectedBorder:     "#ffffff",
			SelectedDescBorder: "#ffffff",
			FilteredMatch:      "#ffffff",
			FilterPrompt:       "#ffffff",
			FilterInfo:         "#ffffff",
			FilterText:         "#ffffff",
			FilterCursor:       "#ffffff",
			HelpKey:            "#ffffff",
			HelpDesc:           "#ffffff",
			PageActiveDot:      "#ffffff",
			PageInactiveDot:    "#ffffff",
			DividerDot:         "#ffffff",
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
