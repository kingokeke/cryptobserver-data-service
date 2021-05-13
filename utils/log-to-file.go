package utils

import (
	"log"
	"os"
)

func logToFile(filename string, logContent string) {
	logFile, e := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	CheckError(e)

	defer logFile.Close()

	log.SetOutput(logFile)
	log.Print(logContent)
}

func LogToGeneral(logContent string) {
	logToFile("raw-data-logs.log", logContent)
}