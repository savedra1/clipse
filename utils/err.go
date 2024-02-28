package utils

import (
	"fmt"
	"os"
	"runtime/debug"
)

func HandleError(err error) {
	if err != nil {
		debug.PrintStack()
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
