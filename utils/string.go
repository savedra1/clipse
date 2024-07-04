package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/* General purpose functions to be used by other modules
 */

func Shorten(s string) string {
	sl := strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(s, "\n", "\\n"),
			"\t", " ",
		),
	)
	if len(sl) <= maxChar {
		return strings.ReplaceAll(sl, "  ", " ")
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

func GetTime() string {
	return strings.TrimSpace(strings.Split(time.Now().String(), "+")[0])
}

func GetTimeStamp() string {
	return strings.Split(GetTime(), ".")[1]
}

// extact the byte data size from the image's clipboard title
func GetImgIdentifier(itemName string) string {
	itemName = strings.TrimSpace(itemName)
	if !strings.Contains(itemName, " ") || !strings.Contains(itemName, "-") {
		LogERROR(fmt.Sprintf("could not get img identifier due to irregular filename | '%s'", itemName))
		return ""
	}
	return strings.Split(strings.Split(itemName, "-")[0], " ")[1]
}

// Expands the path to include the home directory if the path is prefixed
// with `~`. If it isn't prefixed with `~`, the path is returned as-is.
func ExpandHome(relPath string) string {
	if len(relPath) == 0 {
		return relPath
	}

	if relPath[0] != '~' {
		// if not ~, it could be $HOME. Expand that.
		return os.ExpandEnv(relPath)
	}

	curUserHome, err := os.UserHomeDir()
	HandleError(err)

	return filepath.Join(curUserHome, relPath[1:])
}

func ExpandRel(relPath, absPath string) string {
	// Already absolute.
	if filepath.IsAbs(relPath) {
		return relPath
	}

	absRelPath, err := filepath.Abs(filepath.Join(absPath, relPath))
	if err != nil {
		fmt.Println("Current working directory is INVALID! How did you manage this?")
	}
	return absRelPath
}

/* NOT IN USE - Remove bad chars - can cause issues with fuzzy finder
func cleanString(s string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]+")
	sanitized := regex.ReplaceAllString(s, "")
	sl := strings.ReplaceAll(sanitized, "\n", "\\n")
	return strings.ReplaceAll(sl, "  ", " ")         // remove trailing space
}*/
