package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/v-mars/zlog"
)

func main() {
	// Test the directory creation functionality by using a log file in a non-existent directory
	logDir := "./test/logs"
	logFile := filepath.Join(logDir, "test.log")

	fmt.Printf("Testing automatic directory creation for log file: %s\n", logFile)

	// Create rotation config with the log file in a subdirectory
	rotateConfig := zlog.GetDefaultRotateConfig(logFile,
		zlog.WithMaxSize(10),      // 10 MB
		zlog.WithMaxBackups(3),    // Keep 3 backups
		zlog.WithMaxAge(7),        // Keep logs for 7 days
	)

	// Create a rotating logger
	rotatingLogger := zlog.NewRotatingLoggerWithFormat(rotateConfig, zlog.ConsoleFormat)

	// Verify that the directory was created
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		fmt.Printf("ERROR: Directory %s was not created!\n", logDir)
		return
	}

	fmt.Printf("SUCCESS: Directory %s was created successfully.\n", logDir)

	// Test logging to the file
	rotatingLogger.Info("Test log message to verify logger works correctly.")
	rotatingLogger.Error("Another test log message.")

	fmt.Println("Log messages written successfully!")

	// Verify that the log file was created
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Printf("ERROR: Log file %s was not created!\n", logFile)
		return
	}

	fmt.Printf("SUCCESS: Log file %s was created and written to successfully.\n", logFile)

	// Clean up the test directory
	os.RemoveAll("./test")
	fmt.Println("Test completed and cleanup done.")
}