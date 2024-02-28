package utils

import (
	"io"
	"os"
	"strings"
)

/* General purpose functions to be used by other modules
 */

func Shorten(s string) string {
	sl := strings.ReplaceAll(s, "\n", "\\n") // make single line
	if len(sl) <= maxChar {                  // maxChar defined in constants.go
		return strings.ReplaceAll(sl, "  ", " ") // remove double spaces
	}
	return strings.ReplaceAll(sl[:maxChar-3], "  ", " ") + "..."
}

func GetStdin() string {
	/*
		Gets piped input from the terminal when n
		no additional arg provided
	*/
	buffer := make([]byte, 1024)
	n, err := os.Stdin.Read(buffer)
	if err != nil && err != io.EOF {
		return "Error reading Stdin"
	}
	return string(buffer[:n])

}

/* NOT IN USE - Remove bad chars - can cause issues with fuzzy finder
func cleanString(s string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]+")
	sanitised := regex.ReplaceAllString(s, "")
	sl := strings.ReplaceAll(sanitised, "\n", "\\n")
	return strings.ReplaceAll(sl, "  ", " ")         // remove trailing space
}*/
