// Package zlog provides adapter for hertz hlog compatibility
package zlog

import (
	"context"
	"io"

	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
)

// HlogAdapter adapts ZLogger to be compatible with hlog interface
type HlogAdapter struct {
	logger *ZLogger
}

// NewHlogAdapter creates a new adapter that makes ZLogger compatible with hlog
func NewHlogAdapter(zlogger *ZLogger) *HlogAdapter {
	return &HlogAdapter{
		logger: zlogger,
	}
}

// Ensure HlogAdapter implements hlog.FullLogger interface
var _ hertzlog.FullLogger = (*HlogAdapter)(nil)

// Logger interface methods (hlog compatibility)
func (h *HlogAdapter) Trace(v ...interface{}) {
	h.logger.Trace(v...)
}

func (h *HlogAdapter) Debug(v ...interface{}) {
	h.logger.Debug(v...)
}

func (h *HlogAdapter) Info(v ...interface{}) {
	h.logger.Info(v...)
}

func (h *HlogAdapter) Notice(v ...interface{}) {
	h.logger.Notice(v...)
}

func (h *HlogAdapter) Warn(v ...interface{}) {
	h.logger.Warn(v...)
}

func (h *HlogAdapter) Error(v ...interface{}) {
	h.logger.Error(v...)
}

func (h *HlogAdapter) Fatal(v ...interface{}) {
	h.logger.Fatal(v...)
}

// FormatLogger interface methods (hlog compatibility)
func (h *HlogAdapter) Tracef(format string, v ...interface{}) {
	h.logger.Tracef(format, v...)
}

func (h *HlogAdapter) Debugf(format string, v ...interface{}) {
	h.logger.Debugf(format, v...)
}

func (h *HlogAdapter) Infof(format string, v ...interface{}) {
	h.logger.Infof(format, v...)
}

func (h *HlogAdapter) Noticef(format string, v ...interface{}) {
	h.logger.Noticef(format, v...)
}

func (h *HlogAdapter) Warnf(format string, v ...interface{}) {
	h.logger.Warnf(format, v...)
}

func (h *HlogAdapter) Errorf(format string, v ...interface{}) {
	h.logger.Errorf(format, v...)
}

func (h *HlogAdapter) Fatalf(format string, v ...interface{}) {
	h.logger.Fatalf(format, v...)
}

// CtxLogger interface methods (hlog compatibility)
func (h *HlogAdapter) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxTracef(ctx, format, v...)
}

func (h *HlogAdapter) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxDebugf(ctx, format, v...)
}

func (h *HlogAdapter) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxInfof(ctx, format, v...)
}

func (h *HlogAdapter) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxNoticef(ctx, format, v...)
}

func (h *HlogAdapter) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxWarnf(ctx, format, v...)
}

func (h *HlogAdapter) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxErrorf(ctx, format, v...)
}

func (h *HlogAdapter) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	h.logger.CtxFatalf(ctx, format, v...)
}

// Control interface methods (hlog compatibility)
func (h *HlogAdapter) SetLevel(level hertzlog.Level) {
	h.logger.SetLevel(level)
}

func (h *HlogAdapter) SetOutput(w io.Writer) {
	h.logger.SetOutput(w)
}

// Convenience function to set hlog's default logger to use our ZLogger
func SetAsHlogDefault(zlogger *ZLogger) {
	adapter := NewHlogAdapter(zlogger)
	hertzlog.SetLogger(adapter)
}

// Convenience function to set hlog's system logger to use our ZLogger
func SetAsHlogSystem(zlogger *ZLogger) {
	adapter := NewHlogAdapter(zlogger)
	hertzlog.SetSystemLogger(adapter)
}

// GetDefaultLogger returns the default Hertz logger
func GetDefaultLogger() hertzlog.FullLogger {
	return hertzlog.DefaultLogger()
}

// GetSystemLogger returns the system Hertz logger
func GetSystemLogger() hertzlog.FullLogger {
	return hertzlog.SystemLogger()
}