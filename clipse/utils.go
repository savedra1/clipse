package main

import (
	"bytes"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

/* General purpose functions to be used by other modules
 */

// Avoids repeat code by handling errors in a uniform way
func handleError(err error) {
	if err != nil {
		debug.PrintStack()
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

// Contains checks if a string exists in the most recent 3 items
func contains(str string) bool {
	data := getHistory()

	if len(data) > 3 {
		data = data[:3]
	}
	for _, item := range data {

		if item.Value == str {
			return true
		}
	}
	return false
}

// Shortens string val to show in list view
func shorten(s string) string {
	sl := strings.ReplaceAll(s, "\n", "\\n") // make single line
	if len(sl) <= maxChar {                  // maxChar defined in constants.go
		return strings.ReplaceAll(sl, "  ", " ") // remove double spaces
	}
	return strings.ReplaceAll(sl[:maxChar-3], "  ", " ") + "..."
}

/* NOT IN USE - Remove bad chars - can cause issues with fuzzy finder
func cleanString(s string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]+")
	sanitised := regex.ReplaceAllString(s, "")
	sl := strings.ReplaceAll(sanitised, "\n", "\\n")
	return strings.ReplaceAll(sl, "  ", " ")         // remove trailing space
}*/

func dataType(data string) string {
	/*
	   Confirms if clipboard data is currently folding a file vs a string
	*/
	dataBytes := []byte(data)
	reader := bytes.NewReader(dataBytes)

	_, err := png.Decode(reader)
	if err == nil {
		return "png"
	}
	_, err = jpeg.Decode(reader)
	if err == nil {
		return "jpg"
	}
	_, err = gif.Decode(reader)
	if err == nil {
		return "gif"
	}

	return "text"
}

func getTime() string {
	return strings.TrimSpace(strings.Split(time.Now().UTC().String(), "+0000")[0])
}
