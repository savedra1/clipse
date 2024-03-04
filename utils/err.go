package utils

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func HandleError(err error) {
	if err != nil {
		debug.PrintStack()
		log.Fatalln(err)
		errLog(fmt.Sprintf("%s", err))
		os.Exit(1)
	}
}

func errLog(msg string) {
	file, err := os.OpenFile("~/.config/clipse/errLog.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to open ~/.config/clipse/errLog.txt: %s", err))
	}
	defer file.Close()

	if _, err := file.WriteString(msg); err != nil {
		log.Fatalln(err)
	}

}
