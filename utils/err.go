package utils

import (
	"fmt"
	"os"
	"runtime/debug"
)

func HandleError(err error) {
	if err != nil {
		debug.PrintStack()
		if logger != nil {
			LogERROR(fmt.Sprint(err))
		}
		os.Exit(1)
	}
}
