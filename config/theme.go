package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/savedra1/clipse/utils"
)

type CustomTheme struct {
	UseCustom          bool   `json:"useCustom"`
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
	PreviewedText      string `json:"PreviewedText"`
	PreviewBorder      string `json:"PreviewBorder"`
}

func GetTheme() CustomTheme {
	_, err := os.Stat(ClipseConfig.ThemeFilePath)
	if os.IsNotExist(err) {
		if err = initDefaultTheme(); err != nil {
			utils.LogERROR(fmt.Sprintf("could not initialize theme: %s", err))
			return defaultTheme()
		}
	}

	file, err := os.OpenFile(ClipseConfig.ThemeFilePath, os.O_RDONLY, 0644)
	if err != nil {
		file.Close()
	}

	var theme CustomTheme

	if err := json.NewDecoder(file).Decode(&theme); err != nil {
		utils.LogERROR(
			fmt.Sprintf(
				"Error decoding JSON for custom_theme.json. Try creating this file manually instead. Err: %s",
				err,
			),
		)
	}
	if !theme.UseCustom {
		return defaultTheme()
	}
	return theme
}

func initDefaultTheme() error {
	/*
	  Creates custom_theme.json file is not found in path
	  and sets base config.
	*/
	_, err := os.Stat(ClipseConfig.ThemeFilePath)
	if os.IsNotExist(err) {

		baseConfig := defaultTheme()

		jsonData, err := json.MarshalIndent(baseConfig, "", "    ")
		if err != nil {
			return err
		}

		if err = os.WriteFile(ClipseConfig.ThemeFilePath, jsonData, 0644); err != nil {
			return err
		}

		return nil
	}
	return nil
}

// hardcoded default theme when UseCustom set to false
func defaultTheme() CustomTheme {
	return CustomTheme{
		UseCustom:          false,
		TitleFore:          "#ffffff",
		TitleBack:          "#6F4CBC",
		TitleInfo:          "#3498db",
		NormalTitle:        "#ffffff",
		DimmedTitle:        "#808080",
		SelectedTitle:      "#FF69B4",
		NormalDesc:         "#808080",
		DimmedDesc:         "#808080",
		SelectedDesc:       "#FF69B4",
		StatusMsg:          "#2ecc71",
		PinIndicatorColor:  "#FFD700",
		SelectedBorder:     "#3498db",
		SelectedDescBorder: "#3498db",
		FilteredMatch:      "#ffffff",
		FilterPrompt:       "#2ecc71",
		FilterInfo:         "#3498db",
		FilterText:         "#ffffff",
		FilterCursor:       "#FFD700",
		HelpKey:            "#999999",
		HelpDesc:           "#808080",
		PageActiveDot:      "#3498db",
		PageInactiveDot:    "#808080",
		DividerDot:         "#3498db",
		PreviewedText:      "#ffffff",
		PreviewBorder:      "#3498db",
	}
}
