package utils

import "strconv"

func IsInt(arg string) bool {
	_, err := strconv.Atoi(arg)
	if err != nil {
		return false
	}
	return true
}
