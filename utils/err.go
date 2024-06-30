package utils

import (
	"fmt"
	"os"
	"runtime/debug"
)

func HandleError(err error) {
	if err != nil {
		debug.PrintStack()
		LogERROR(fmt.Sprint(err))
		os.Exit(1)
	}
}
