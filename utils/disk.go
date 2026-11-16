package utils

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func DiskspaceAvailable(path string, bytes int) bool {
	dir := filepath.Dir(path)
	var stat unix.Statfs_t
	if err := unix.Statfs(dir, &stat); err != nil {
		LogERROR(fmt.Sprintf("failed to check disk space at '%s': %s", path, err))
		return true
	}

	bytefree := (stat.Bavail * uint64(stat.Bsize))

	return bytefree > (uint64(bytes) * 2) // *2 for safety buffer
}
