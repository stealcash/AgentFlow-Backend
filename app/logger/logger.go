package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile       *os.File
	currentLogDay string
)

// InitLogger creates the daily log file in storage/logs
func InitLogger() {
	currentLogDay = time.Now().Format("2006-01-02")
	logFilePath := filepath.Join("storage", "logs")

	// Ensure the folder exists
	if err := os.MkdirAll(logFilePath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log folder: %v", err)
	}

	openLogFile()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(logFile)
}

// openLogFile opens the log file for the current day
func openLogFile() {
	fileName := fmt.Sprintf("%s.log", currentLogDay)
	fullPath := filepath.Join("storage", "logs", fileName)

	var err error
	logFile, err = os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
}

// Close closes the log file
func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

// checkRotate checks if the date has changed and rotates the file if needed
func checkRotate() {
	today := time.Now().Format("2006-01-02")
	if today != currentLogDay {
		Close()
		currentLogDay = today
		openLogFile()
		log.SetOutput(logFile)
	}
}

// Info logs an info message with format + args
func Info(format string, args ...interface{}) {
	checkRotate()
	msg := fmt.Sprintf(format, args...)
	log.Printf("[INFO] %s\n", msg)
}

// Error logs an error message with format + args
func Error(format string, args ...interface{}) {
	checkRotate()
	msg := fmt.Sprintf(format, args...)
	log.Printf("[ERROR] %s\n", msg)
}

// Track logs a tracker message with format + args
func Track(format string, args ...interface{}) {
	checkRotate()
	msg := fmt.Sprintf(format, args...)
	log.Printf("[TRACK] %s\n", msg)
}

// InfoArgs logs with auto-join like fmt.Println (no format string needed)
func InfoArgs(args ...interface{}) {
	checkRotate()
	msg := fmt.Sprintln(args...)
	log.Printf("[INFO] %s", msg)
}

// ErrorArgs logs with auto-join like fmt.Println (no format string needed)
func ErrorArgs(args ...interface{}) {
	checkRotate()
	msg := fmt.Sprintln(args...)
	log.Printf("[ERROR] %s", msg)
}

// TrackArgs logs with auto-join like fmt.Println (no format string needed)
func TrackArgs(args ...interface{}) {
	checkRotate()
	msg := fmt.Sprintln(args...)
	log.Printf("[TRACK] %s", msg)
}
