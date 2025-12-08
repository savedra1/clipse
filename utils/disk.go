package utils

import (
	"fmt"

	"golang.org/x/sys/unix"
)

func DiskspaceAvailable(bytes int) bool {
	var stat unix.Statfs_t
	if err := unix.Statfs("/", &stat); err != nil {
		LogERROR(fmt.Sprintf("failed to check disk space: %s", err))
		return true
	}

	bytefree := (stat.Bavail * uint64(stat.Bsize))

	return bytefree > (uint64(bytes) * 2) // *2 for safety buffer
}
