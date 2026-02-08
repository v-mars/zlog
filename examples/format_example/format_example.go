package main

import (
	"os"
	"time"

	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/v-mars/zlog"
)

func main() {
	// Example 1: Console format (default)
	consoleLogger := zlog.New(
		zlog.WithFormat(zlog.ConsoleFormat),
		zlog.WithLevel(hertzlog.LevelDebug),
	)

	consoleLogger.Info("This is a console format log message")
	consoleLogger.Debugf("Debug message: %s", "console format")
	consoleLogger.Warn("Warning message")
	consoleLogger.Errorf("Error occurred: %v", "sample error")

	time.Sleep(time.Millisecond * 100) // Small delay to separate logs

	// Example 2: JSON format
	jsonLogger := zlog.New(
		zlog.WithFormat(zlog.JSONFormat),
		zlog.WithLevel(hertzlog.LevelDebug),
	)

	jsonLogger.Info("This is a JSON format log message")
	jsonLogger.Debugf("Debug message: %s", "JSON format")
	jsonLogger.Warn("Warning message")
	jsonLogger.Errorf("Error occurred: %v", "sample error")

	time.Sleep(time.Millisecond * 100) // Small delay to separate logs

	// Example 3: JSON format with custom output file
	file, _ := os.Create("example.log")
	defer file.Close()

	fileLogger := zlog.New(
		zlog.WithFormat(zlog.JSONFormat),
		zlog.WithOutput(file),
		zlog.WithLevel(hertzlog.LevelInfo),
	)

	fileLogger.Info("This is a JSON format log message written to a file")
	fileLogger.Errorf("Error occurred in file: %v", "file error")
	fileLogger.Notice("Notice message in file")

	// Example 4: Using rotating logger with JSON format
	rotateConfig := &zlog.RotateConfig{
		Filename:   "rotated.log",
		MaxSize:    1, // 1 MB
		MaxBackups: 3,
		MaxAge:     7, // 7 days
		Compress:   true,
		LocalTime:  true,
	}

	rotatingLogger := zlog.NewRotatingLoggerWithFormat(rotateConfig, zlog.JSONFormat)
	rotatingLogger.Info("This is a JSON format message in a rotating log")
	rotatingLogger.Errorf("Error in rotating log: %v", "rotating error")

	// Example 5: Console format with rotating logger
	consoleRotatingLogger := zlog.NewRotatingLoggerWithFormat(rotateConfig, zlog.ConsoleFormat)
	consoleRotatingLogger.Info("This is a console format message in a rotating log")
	consoleRotatingLogger.Errorf("Error in console rotating log: %v", "console rotating error")
}