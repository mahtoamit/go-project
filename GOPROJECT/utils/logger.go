package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var Logger *log.Logger
var once sync.Once
var lastRotation time.Time
var lastLogEntryTime time.Time

func InitLogger() {
	once.Do(func() {
		logFile, err := os.OpenFile("record_log/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Error creating log file:", err)
			return
		}
		Logger = log.New(logFile, "", 0) // Set flags to 0 to remove the default timestamp
		lastRotation = time.Now()
		lastLogEntryTime = time.Now()

		// Start a goroutine to periodically check for log rotation
		go performLogRotation()
	})
}

// Log logs the message with the specified log level, API name, endpoint, function name, user's id, and message.
func Log(level, apiName, endpoint, funcName, userId, message string, times ...time.Time) {
	timestamp := time.Now().Format("2006-01-02T15:04:05.000")
	fileName := "app-" + time.Now().Format("2006-01-02-15-04") + ".log"
	logFilePath := filepath.Join("record_log", fileName)

	rotateLogFile(logFilePath)

	responseTime := ""
	if len(times) >= 2 {
		responseTime = times[1].Sub(times[0]).String()
	}

	logEntry := fmt.Sprintf("%s %s %s %s %s %s %s ResponseTime:%s", timestamp, level, apiName, endpoint, funcName, userId, message, responseTime)
	Logger.Println(logEntry)

	lastLogEntryTime = time.Now()
}

// performLogRotation checks if a log rotation is required every minute
func performLogRotation() {
	for {
		time.Sleep(time.Minute)

		lastRotation = time.Now()
	}
}

// rotateLogFile creates a new log file and updates the Logger to write to the new file
func rotateLogFile(logFilePath string) {
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}

	Logger.SetOutput(logFile) // Update the log file for rotation
}
