package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/savedra1/clipse/utils"
)

type Config struct {
	Sources     []string `json:"sources"`
	MaxHistory  int      `json:"maxHistory"`
	HistoryFile string   `json:"historyFile"`
}

type source struct {
	SourceType string `json:"sourceType"`
}

// Global config object, accessed and used when any configuration is needed.
var ClipseConfig = Config{}

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

	confData, err := os.ReadFile(configPath)
	if err = json.Unmarshal(confData, &tempConfig); err != nil {
		fmt.Println("Failed to read config. Fallback to default.\nErr: %w", err)
		ClipseConfig = defaultConfig()
	}
	// fmt.Println("WE AT LEAST GET TO HERE!")

	for _, src := range tempConfig.Sources {
		loadSource(src)
	}
	
	mergeStructs(&ClipseConfig, &tempConfig)
	fmt.Printf("%+v\n", ClipseConfig)
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
	case "history":
	case "":
		fmt.Printf("Error: \"sourceType\" tag not found in source file: %s. File skipped.\n", path)
		fmt.Println("Possible values for sourceType:\n\t- config\n\t- theme\n\t- history")
	default:
		fmt.Printf("Error: Invalid value \"%s\" in \"sourceType\" tag for source file: %s\n",
			src.SourceType, path)
		fmt.Println("Possible values for sourceType:\n\t- config\n\t- theme\n\t- history")
	}
}

/*Merges 2 structs together. It has one fatal flaw that I am unsure how to fix:
It does not let you overrite the destination (dest) with a newer value from 
source (src) if the src value is the zero value of a primitive. I see the main
issue arising if the config has a bool option in the future, which will never
be able to overwrite a true value with a false value.

I can't find a way around this because Go does not have optional types, nor are
primitives nilable, so there is no explicit "not a value" that we can check.

One fix would be for the Config struct's members to be pointers to primitives,
or be interfaces, but I am not too sure about this method.
*/
func mergeStructs(dest, src interface{}) {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		sField := srcVal.Field(i)
		dField := destVal.Field(i)

		switch {
			// If field is a slice, merge them (eg: For getting a merged list of sources)
		case sField.Kind() == reflect.Slice:
			// If the slice in the destination is nil, set it. Else, append.
			if dField.IsNil() {
				dField.Set(sField)
			} else {
				for j := 0; j < sField.Len(); j++ {
					dField.Set(reflect.Append(dField, sField.Index(j)))
				}
			}

			// if source is nil, we don't copy anything. Move on to next field.
		case sField.IsZero():
			continue

			// By default, overwrite the dest value with the src value.
		default:
			dField.Set(sField)
		}
	}
}
