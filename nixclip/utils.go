package main

import (
	"fmt"
	"os"
	"strings"
)

/* General purpose functions to be used by other modules
 */

// Avoids repeat code by handling errors in a uniform way
func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// Contains checks if a string exists in a slice of strings
func contains(slice []ClipboardItem, str string) bool {
	for _, item := range slice {
		if item.Value == str {
			return true
		}
	}
	return false
}

// Shortens string val to show in list view
func shorten(s string) string {
	if len(s) <= maxLen { // maxLen defined in constants.go
		return strings.ReplaceAll(s, "\n", " ")
	}
	return strings.ReplaceAll(s[:maxLen-3], "\n", " ") + "..."
}
