package utils

import (
	"fmt"

	"golang.org/x/sys/unix"
)

//path ex. config.ClipseConfig.HistoryFilePath, config.ClipseConfig.TempDirPath, defaults to "/"
func DiskspaceAvailable(bytes int, path ...string) bool {
    checkPath := "/"
    if len(path) > 0 {
        checkPath = path[0]
    }

    var stat unix.Statfs_t
    if err := unix.Statfs(checkPath, &stat); err != nil {
        LogERROR(fmt.Sprintf("failed to check disk space: %s", err))
        return true
    }
    bytefree := stat.Bavail * uint64(stat.Bsize)

    return bytefree > (uint64(bytes) * 2) // *2 for safety buffer
}