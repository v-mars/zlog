// Package zlog provides log rotation functionality
package zlog

import (
	"context"
	"fmt"
	"io"

	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// RotatingLogger provides log rotation functionality
type RotatingLogger struct {
	baseLogger *ZLogger
	writer     io.Writer
	config     *RotateConfig
}

// RotateConfig holds the configuration for log rotation
type RotateConfig struct {
	Filename   string // Filename is the file to write logs to
	MaxSize    int    // MaxSize is the maximum size in megabytes of the log file before rotation
	MaxBackups int    // MaxBackups is the maximum number of old log files to retain
	MaxAge     int    // MaxAge is the maximum number of days to retain old log files
	Compress   bool   // Compress determines if the rotated log files should be compressed
	LocalTime  bool   // LocalTime determines if the time used for formatting the timestamps in backup files is the computer's local time
}

// NewRotatingLogger creates a new logger with rotation capabilities
func NewRotatingLogger(config *RotateConfig) *RotatingLogger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  config.LocalTime,
	}

	// Create a new ZLogger with lumberjack writer using console format by default
	zLogger := New(WithOutput(lumberjackLogger), WithFormat(ConsoleFormat))

	return &RotatingLogger{
		baseLogger: zLogger,
		writer:     lumberjackLogger,
		config:     config,
	}
}

// NewRotatingLoggerWithFormat creates a new logger with rotation capabilities and specified format
func NewRotatingLoggerWithFormat(config *RotateConfig, format FormatType) *RotatingLogger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  config.LocalTime,
	}

	// Create a new ZLogger with lumberjack writer and specified format
	zLogger := New(WithOutput(lumberjackLogger), WithFormat(format))

	return &RotatingLogger{
		baseLogger: zLogger,
		writer:     lumberjackLogger,
		config:     config,
	}
}

// WithRotation is an option function that configures the logger with rotation
func WithRotation(config *RotateConfig) Option {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime:  config.LocalTime,
	}

	return WithOutput(lumberjackLogger)
}

// WithRotationAndFormat is an option function that configures the logger with rotation and format
func WithRotationAndFormat(rotationConfig *RotateConfig, format FormatType) Option {
	return func(c *config) {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   rotationConfig.Filename,
			MaxSize:    rotationConfig.MaxSize,
			MaxBackups: rotationConfig.MaxBackups,
			MaxAge:     rotationConfig.MaxAge,
			Compress:   rotationConfig.Compress,
			LocalTime:  rotationConfig.LocalTime,
		}

		c.output = lumberjackLogger
		c.format = format
	}
}

// RotateConfigOption configures the RotateConfig
type RotateConfigOption func(*RotateConfig)

// WithMaxSize sets the maximum size in megabytes of the log file before rotation
func WithMaxSize(maxSize int) RotateConfigOption {
	return func(c *RotateConfig) {
		c.MaxSize = maxSize
	}
}

// WithMaxBackups sets the maximum number of old log files to retain
func WithMaxBackups(maxBackups int) RotateConfigOption {
	return func(c *RotateConfig) {
		c.MaxBackups = maxBackups
	}
}

// WithMaxAge sets the maximum number of days to retain old log files
func WithMaxAge(maxAge int) RotateConfigOption {
	return func(c *RotateConfig) {
		c.MaxAge = maxAge
	}
}

// WithCompress sets whether the rotated log files should be compressed
func WithCompress(compress bool) RotateConfigOption {
	return func(c *RotateConfig) {
		c.Compress = compress
	}
}

// WithLocalTime sets whether to use local time for formatting timestamps in backup files
func WithLocalTime(localTime bool) RotateConfigOption {
	return func(c *RotateConfig) {
		c.LocalTime = localTime
	}
}

// GetDefaultRotateConfig returns a default rotation configuration with optional configurations
func GetDefaultRotateConfig(filename string, opts ...RotateConfigOption) *RotateConfig {
	c := &RotateConfig{
		Filename:   filename,
		MaxSize:    20,   // 20 MB
		MaxBackups: 5,    // Keep 5 backups
		MaxAge:     10,   // Keep logs for 10 days
		Compress:   true, // Compress rotated files
		LocalTime:  true, // Use local time for filenames
	}

	// Apply any provided options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Rotate manually rotates the log file
func (rl *RotatingLogger) Rotate() error {
	if lj, ok := rl.writer.(*lumberjack.Logger); ok {
		return lj.Rotate()
	}
	return fmt.Errorf("unable to rotate: writer is not a lumberjack.Logger")
}

// GetRotatingWriter returns the underlying lumberjack writer for direct access
func (rl *RotatingLogger) GetRotatingWriter() *lumberjack.Logger {
	if lj, ok := rl.writer.(*lumberjack.Logger); ok {
		return lj
	}
	return nil
}

// SetLevel implements the Control interface for RotatingLogger
func (rl *RotatingLogger) SetLevel(level hertzlog.Level) {
	rl.baseLogger.SetLevel(level)
}

// SetOutput implements the Control interface for RotatingLogger
func (rl *RotatingLogger) SetOutput(w io.Writer) {
	// We typically don't want to change the output for rotating loggers
	// but if we do, we'd need to handle this differently
	// For now, just update the underlying logger with the new writer
	newZLogger := New(WithOutput(w), WithLevel(rl.baseLogger.level))
	rl.baseLogger = newZLogger
}

// Implement all the logging methods by delegating to the base logger

// Logger methods
func (rl *RotatingLogger) Trace(v ...interface{})  { rl.baseLogger.Trace(v...) }
func (rl *RotatingLogger) Debug(v ...interface{})  { rl.baseLogger.Debug(v...) }
func (rl *RotatingLogger) Info(v ...interface{})   { rl.baseLogger.Info(v...) }
func (rl *RotatingLogger) Notice(v ...interface{}) { rl.baseLogger.Notice(v...) }
func (rl *RotatingLogger) Warn(v ...interface{})   { rl.baseLogger.Warn(v...) }
func (rl *RotatingLogger) Error(v ...interface{})  { rl.baseLogger.Error(v...) }
func (rl *RotatingLogger) Fatal(v ...interface{})  { rl.baseLogger.Fatal(v...) }

// FormatLogger methods
func (rl *RotatingLogger) Tracef(format string, v ...interface{}) { rl.baseLogger.Tracef(format, v...) }
func (rl *RotatingLogger) Debugf(format string, v ...interface{}) { rl.baseLogger.Debugf(format, v...) }
func (rl *RotatingLogger) Infof(format string, v ...interface{})  { rl.baseLogger.Infof(format, v...) }
func (rl *RotatingLogger) Noticef(format string, v ...interface{}) {
	rl.baseLogger.Noticef(format, v...)
}
func (rl *RotatingLogger) Warnf(format string, v ...interface{})  { rl.baseLogger.Warnf(format, v...) }
func (rl *RotatingLogger) Errorf(format string, v ...interface{}) { rl.baseLogger.Errorf(format, v...) }
func (rl *RotatingLogger) Fatalf(format string, v ...interface{}) { rl.baseLogger.Fatalf(format, v...) }

// CtxLogger methods
func (rl *RotatingLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxTracef(ctx, format, v...)
}
func (rl *RotatingLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxDebugf(ctx, format, v...)
}
func (rl *RotatingLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxInfof(ctx, format, v...)
}
func (rl *RotatingLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxNoticef(ctx, format, v...)
}
func (rl *RotatingLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxWarnf(ctx, format, v...)
}
func (rl *RotatingLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxErrorf(ctx, format, v...)
}
func (rl *RotatingLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	rl.baseLogger.CtxFatalf(ctx, format, v...)
}
