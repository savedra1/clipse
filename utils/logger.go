package utils

import (
	"log"
	"os"
)

var logger *log.Logger
var debugging = false

func SetUpLogger(logFilePath string, debug bool) {
	debugging = debug

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func LogERROR(message string) {
	logger.Printf("ERROR: %s", message)
}

func LogINFO(message string) {
	logger.Printf("INFO: %s", message)
}

func LogWARN(message string) {
	logger.Printf("WARN: %s", message)
}

func LogDEBUG(message string) {
	if debugging {
		logger.Printf("DEBUG: %s", message)
	}
}
