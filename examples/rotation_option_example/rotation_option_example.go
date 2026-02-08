package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/v-mars/zlog"
)

func main() {
	// Example 1: Using GetDefaultRotateConfig with Option pattern
	// This demonstrates the flexibility of the configuration
	rotateConfig := zlog.GetDefaultRotateConfig("app.log",
		zlog.WithMaxSize(100),       // Increase max size to 100 MB
		zlog.WithMaxBackups(10),     // Keep 10 backups instead of 5
		zlog.WithMaxAge(30),         // Keep logs for 30 days instead of 10
		zlog.WithCompress(false),    // Disable compression
		zlog.WithLocalTime(true),    // Use local time (default is true anyway)
	)

	consoleRotatingLogger := zlog.NewRotatingLoggerWithFormat(rotateConfig, zlog.ConsoleFormat)
	consoleRotatingLogger.Info("This is a console format message in a rotating log")
	consoleRotatingLogger.Errorf("Error in console rotating log: %v", "console rotating error")

	// Example 2: Custom configuration with minimal changes
	customConfig := zlog.GetDefaultRotateConfig("custom.log",
		zlog.WithMaxSize(50),        // Only change max size
	)

	jsonRotatingLogger := zlog.NewRotatingLoggerWithFormat(customConfig, zlog.JSONFormat)
	jsonRotatingLogger.Info("This is a JSON format message in a rotating log")
	jsonRotatingLogger.Errorf("Error in rotating log: %v", "rotating error")

	// Example 3: Using with zerolog options as well
	enrichedLogger := zlog.New(
		zlog.WithFormat(zlog.JSONFormat),
		zlog.WithLevel(hertzlog.LevelInfo),
		// Add global fields using enricher
		zlog.WithZerologOptions(
			func(logger zerolog.Logger) zerolog.Logger {
				return logger.With().
					Str("service", "rotation-example").
					Str("version", "1.0.0").
					Int("pid", os.Getpid()).Logger()
			},
		),
		zlog.WithRotation(customConfig), // Use rotation with our custom config
	)

	enrichedLogger.Info("Message with global fields and rotation")
	enrichedLogger.Error("Error with global fields and rotation")

	// Example 4: Using with rotation and format options together
	multiOptionLogger := zlog.New(
		zlog.WithFormat(zlog.ConsoleFormat),
		zlog.WithLevel(hertzlog.LevelDebug),
		zlog.WithRotationAndFormat(
			zlog.GetDefaultRotateConfig("multi-option.log",
				zlog.WithMaxSize(25),    // 25 MB max
				zlog.WithMaxBackups(3),  // 3 backups
			),
			zlog.ConsoleFormat,
		),
	)

	multiOptionLogger.Info("Message with multiple options combined")
	multiOptionLogger.Warn("Warning with multiple options combined")

	time.Sleep(time.Millisecond * 100) // Small delay to ensure all logs are written
}