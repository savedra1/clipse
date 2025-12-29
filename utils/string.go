// General purpose functions to be used by other modules

package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// any chars that cause the fuzzy find to crash can be appended here
var badChars = []string{
	"\x00", // \u0000
}

func Shorten(s string, maxChar int) string {
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
	return time.Now().Format(DateLayout)
}

func GetTimeStamp() string {
	return strings.Split(GetTime(), ".")[1]
}

func GetImgIdentifier(filename string) string {
	parts := strings.SplitN(filename, " ", 2)
	if len(parts) < 2 {
		LogERROR(
			fmt.Sprintf(
				"could not get img identifier due to missing space in filename | '%s'",
				filename,
			),
		)
		return ""
	}
	filename = parts[1]
	fileNamePattern := regexp.MustCompile(imgNameRegEx)
	matches := fileNamePattern.FindStringSubmatch(filename)
	if matches == nil {
		LogERROR(
			fmt.Sprintf(
				"could not get img identifier due to irregular filename | '%s'",
				filename,
			),
		)
		return ""
	}
	return matches[1]
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

func CleanPath(fp string) string {
	if strings.Contains(fp, " ") {
		return fmt.Sprintf("'%s'", fp)
	}
	return fp
}

func ParseDuration(s string) (*time.Duration, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return nil, err
	}
	return &duration, nil
}

func SanitizeChars(s string) string {
	for _, char := range badChars {
		s = strings.ReplaceAll(s, char, "")
	}

	return s
}
