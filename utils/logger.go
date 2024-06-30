package utils

import (
	"log"
	"os"
)

var logger *log.Logger

func SetUpLogger(logFilePath string) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func LogERROR(message string) {
	logger.Printf("ERROR: " + message)
}

func LogINFO(message string) {
	logger.Printf("INFO: " + message)
}

func LogWARN(message string) {
	logger.Printf("WARN: " + message)
}
