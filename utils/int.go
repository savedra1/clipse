package utils

import "strconv"

func IsInt(arg string) bool {
	if _, err := strconv.Atoi(arg); err != nil {
		return false
	}
	return true
}
