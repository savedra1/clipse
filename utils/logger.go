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
		return
	}
	logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func LogERROR(message string) {
	if logger != nil {
		logger.Printf("ERROR: %s", message)
		return
	}
	log.Fatalf("ERROR: %s", message)
}

func LogINFO(message string) {
	if logger != nil {
		logger.Printf("INFO: %s", message)
		return
	}
	log.Fatalf("INFO: %s", message)
}

func LogWARN(message string) {
	if logger != nil {
		logger.Printf("WARN: %s", message)
		return
	}
	log.Fatalf("WARN: %s", message)
}
