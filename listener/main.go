package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
)

// Data struct for storing clipboard strings
type Data struct {
	ClipboardHistory []string `json:"clipboardHistory"`
}

func main() {
	// Listen for SIGINT (Ctrl+C) and SIGTERM signals to properly close the program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Load existing data from file, if any
	var data Data
	err := loadDataFromFile("../history.json", &data)
	if err != nil {
		fmt.Println("Error loading data from file:", err)
	}

	// Start a goroutine to continuously monitor clipboard changes
	go func() {
		for {
			// Get the current clipboard content
			text, err := clipboard.ReadAll()
			if err != nil {
				fmt.Println("Error reading clipboard:", err)
			}

			// If clipboard content is not empty and not already in the list, add it
			if text != "" && !contains(data.ClipboardHistory, text) {
				// If the length exceeds 20, remove the oldest item
				if len(data.ClipboardHistory) >= 20 {
					data.ClipboardHistory = data.ClipboardHistory[1:] // Remove the oldest item (first element)
				}
				data.ClipboardHistory = append(data.ClipboardHistory, text)
				fmt.Println("Added to clipboard history:", text)

				// Save data to file
				err := saveDataToFile("../history.json", data)
				if err != nil {
					fmt.Println("Error saving data to file:", err)
				}
			}

			// Check for updates every second
			time.Sleep(time.Second)
		}
	}()

	fmt.Println("Clipboard history listener running... Press Ctrl+C to exit.")

	// Wait for SIGINT or SIGTERM signal
	<-interrupt
	fmt.Println("Exiting...")
}

// loadDataFromFile loads data from a JSON file.
func loadDataFromFile(filename string, data *Data) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(data)
}

// saveDataToFile saves data to a JSON file.
func saveDataToFile(filename string, data Data) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(data)
}

// contains checks if a string is present in a slice of strings.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
