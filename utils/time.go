package utils

import (
	"strings"
	"time"
)

func GetTime() string {
	return strings.TrimSpace(strings.Split(time.Now().UTC().String(), "+0000")[0])
}
