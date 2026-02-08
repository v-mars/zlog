// Package zlog is a flexible and high-performance logging library for Go applications
// that supports integration with Hertz's hlog and provides log rotation capabilities.
package zlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	LogIDKey = "request_id"
	ReqIDKey = "X-Request-ID"
)

// Logger defines the core logging interface with basic log levels
type Logger interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Notice(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

// FormatLogger extends Logger with formatted logging methods
type FormatLogger interface {
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Noticef(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

// CtxLogger provides logging methods that accept a context
type CtxLogger interface {
	CtxTracef(ctx context.Context, format string, v ...interface{})
	CtxDebugf(ctx context.Context, format string, v ...interface{})
	CtxInfof(ctx context.Context, format string, v ...interface{})
	CtxNoticef(ctx context.Context, format string, v ...interface{})
	CtxWarnf(ctx context.Context, format string, v ...interface{})
	CtxErrorf(ctx context.Context, format string, v ...interface{})
	CtxFatalf(ctx context.Context, format string, v ...interface{})
}

// Control provides methods to configure the logger
type Control interface {
	SetLevel(level hertzlog.Level)
	SetOutput(w io.Writer)
}

// FullLogger combines all logging interfaces
type FullLogger interface {
	Logger
	FormatLogger
	CtxLogger
	Control
}

// FormatType defines the output format for the logger
type FormatType int

const (
	// ConsoleFormat outputs logs in human-readable console format
	ConsoleFormat FormatType = iota
	// JSONFormat outputs logs in JSON format
	JSONFormat
)

// ZLogger implements the FullLogger interface using zerolog
type ZLogger struct {
	logger zerolog.Logger
	level  hertzlog.Level
	tp     trace.TracerProvider
	format FormatType
}

// Ensure ZLogger implements FullLogger interface
var _ FullLogger = (*ZLogger)(nil)

// New creates a new ZLogger instance
func New(options ...Option) *ZLogger {
	cfg := &config{
		output:          os.Stdout,
		level:           hertzlog.LevelInfo,
		tp:              trace.NewNoopTracerProvider(),
		format:          ConsoleFormat, // Default to console format
		loggerEnrichers: []func(zerolog.Logger) zerolog.Logger{},
	}

	for _, opt := range options {
		opt(cfg)
	}

	var zlogger zerolog.Logger
	switch cfg.format {
	case JSONFormat:
		// JSON format - default zerolog behavior with caller info
		zlogger = zerolog.New(cfg.output).Level(toZerologLevel(cfg.level)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	case ConsoleFormat:
		// Console format - human readable with RFC3339 time format, caller info and custom format
		consoleWriter := &zerolog.ConsoleWriter{
			Out:        cfg.output,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				// Ensure full level name is shown instead of 3-letter abbreviation
				if ll, ok := i.(string); ok {
					return fmt.Sprintf("%-6s", strings.ToUpper(ll))
				}
				return fmt.Sprintf("%-6s", strings.ToUpper(fmt.Sprintf("%s", i)))
			},
		}
		zlogger = zerolog.New(consoleWriter).Level(toZerologLevel(cfg.level)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	default:
		// Default to console format with customization
		consoleWriter := &zerolog.ConsoleWriter{
			Out: cfg.output,
			//TimeFormat: time.RFC3339,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				// Ensure full level name is shown instead of 3-letter abbreviation
				if ll, ok := i.(string); ok {
					return fmt.Sprintf("%-6s", strings.ToUpper(ll))
				}
				return fmt.Sprintf("%-6s", strings.ToUpper(fmt.Sprintf("%s", i)))
			},
		}
		zlogger = zerolog.New(consoleWriter).Level(toZerologLevel(cfg.level)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	}

	// Apply any additional logger enrichments
	for _, enricher := range cfg.loggerEnrichers {
		zlogger = enricher(zlogger)
	}

	return &ZLogger{
		logger: zlogger,
		level:  cfg.level,
		tp:     cfg.tp,
		format: cfg.format,
	}
}

// Option configures the logger
type Option func(*config)

// config holds the configuration for the logger
type config struct {
	output io.Writer
	level  hertzlog.Level
	tp     trace.TracerProvider
	format FormatType
	// Functions to customize the base logger after initial setup
	loggerEnrichers []func(zerolog.Logger) zerolog.Logger
}

// WithOutput sets the output writer for the logger
func WithOutput(output io.Writer) Option {
	return func(c *config) {
		c.output = output
	}
}

// WithLevel sets the log level for the logger
func WithLevel(level hertzlog.Level) Option {
	return func(c *config) {
		c.level = level
	}
}

// WithTraceProvider sets the OpenTelemetry trace provider for the logger
func WithTraceProvider(tp trace.TracerProvider) Option {
	return func(c *config) {
		if tp != nil {
			c.tp = tp
		} else {
			c.tp = trace.NewNoopTracerProvider()
		}
	}
}

// WithFormat sets the output format for the logger (Console or JSON)
func WithFormat(format FormatType) Option {
	return func(c *config) {
		c.format = format
	}
}

// WithZerologOptions sets additional zerolog options using enricher functions
func WithZerologOptions(enrichers ...func(zerolog.Logger) zerolog.Logger) Option {
	return func(c *config) {
		c.loggerEnrichers = append(c.loggerEnrichers, enrichers...)
	}
}

// Helper function to convert hertz log level to zerolog level
func toZerologLevel(level hertzlog.Level) zerolog.Level {
	switch level {
	case hertzlog.LevelTrace:
		return zerolog.TraceLevel
	case hertzlog.LevelDebug:
		return zerolog.DebugLevel
	case hertzlog.LevelInfo:
		return zerolog.InfoLevel
	case hertzlog.LevelNotice, hertzlog.LevelWarn:
		return zerolog.WarnLevel
	case hertzlog.LevelError:
		return zerolog.ErrorLevel
	case hertzlog.LevelFatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// Helper function to convert zerolog level to hertz log level
func fromZerologLevel(level zerolog.Level) hertzlog.Level {
	switch level {
	case zerolog.TraceLevel:
		return hertzlog.LevelTrace
	case zerolog.DebugLevel:
		return hertzlog.LevelDebug
	case zerolog.InfoLevel:
		return hertzlog.LevelInfo
	case zerolog.WarnLevel:
		return hertzlog.LevelNotice
	case zerolog.ErrorLevel:
		return hertzlog.LevelError
	case zerolog.FatalLevel:
		return hertzlog.LevelFatal
	default:
		return hertzlog.LevelInfo
	}
}

// Implementation of Logger interface methods
func (zl *ZLogger) Trace(v ...interface{}) {
	zl.logger.Trace().Msg(fmt.Sprint(v...))
}

func (zl *ZLogger) Debug(v ...interface{}) {
	zl.logger.Debug().Msg(fmt.Sprint(v...))
}

func (zl *ZLogger) Info(v ...interface{}) {
	zl.logger.Info().Msg(fmt.Sprint(v...))
}

func (zl *ZLogger) Notice(v ...interface{}) {
	zl.logger.Warn().Msg(fmt.Sprint(v...)) // Map Notice to Warn level
}

func (zl *ZLogger) Warn(v ...interface{}) {
	zl.logger.Warn().Msg(fmt.Sprint(v...))
}

func (zl *ZLogger) Error(v ...interface{}) {
	zl.logger.Error().Msg(fmt.Sprint(v...))
}

func (zl *ZLogger) Fatal(v ...interface{}) {
	zl.logger.Fatal().Msg(fmt.Sprint(v...))
}

// Implementation of FormatLogger interface methods
func (zl *ZLogger) Tracef(format string, v ...interface{}) {
	zl.logger.Trace().Msgf(format, v...)
}

func (zl *ZLogger) Debugf(format string, v ...interface{}) {
	zl.logger.Debug().Msgf(format, v...)
}

func (zl *ZLogger) Infof(format string, v ...interface{}) {
	zl.logger.Info().Msgf(format, v...)
}

func (zl *ZLogger) Noticef(format string, v ...interface{}) {
	zl.logger.Warn().Msgf(format, v...) // Map Noticef to Warnf level
}

func (zl *ZLogger) Warnf(format string, v ...interface{}) {
	zl.logger.Warn().Msgf(format, v...)
}

func (zl *ZLogger) Errorf(format string, v ...interface{}) {
	zl.logger.Error().Msgf(format, v...)
}

func (zl *ZLogger) Fatalf(format string, v ...interface{}) {
	zl.logger.Fatal().Msgf(format, v...)
}

// Implementation of CtxLogger interface methods
func (zl *ZLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Trace()
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	logEvt.Msgf(format, v...)

	// Add as event to the current span if it exists
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", fmt.Sprintf(format, v...)),
			attribute.String("level", "trace"),
		))
	}
}

func (zl *ZLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Debug()
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	logEvt.Msgf(format, v...)

	// Add as event to the current span if it exists
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", fmt.Sprintf(format, v...)),
			attribute.String("level", "debug"),
		))
	}
}

func (zl *ZLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Info()
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	logEvt.Msgf(format, v...)

	// Add as event to the current span if it exists
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", fmt.Sprintf(format, v...)),
			attribute.String("level", "info"),
		))
	}
}

func (zl *ZLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Warn() // Map Notice to Warn level
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	logEvt.Msgf(format, v...)

	// Add as event to the current span if it exists
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", fmt.Sprintf(format, v...)),
			attribute.String("level", "notice"),
		))
	}
}

func (zl *ZLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Warn()
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	logEvt.Msgf(format, v...)

	// Add as event to the current span if it exists
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", fmt.Sprintf(format, v...)),
			attribute.String("level", "warn"),
		))
	}
}

func (zl *ZLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Error()
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	msg := fmt.Sprintf(format, v...)
	logEvt.Msg(msg)

	// Add as event to the current span if it exists and mark span as error
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", msg),
			attribute.String("level", "error"),
		))

		// Mark the span as error
		span.SetStatus(codes.Error, msg)
	}
}

func (zl *ZLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	// Let the otel.go implementation handle this to avoid duplication
	// For now, we'll call the basic logger with OTel fields
	fields := zl.getOtelFields(ctx)

	logEvt := zl.logger.Fatal()
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	msg := fmt.Sprintf(format, v...)
	logEvt.Msg(msg)

	// Add as event to the current span if it exists and mark span as error
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", msg),
			attribute.String("level", "fatal"),
		))

		// Mark the span as error
		span.SetStatus(codes.Error, msg)
	}
}

// getOtelFields extracts OpenTelemetry fields from context
func (zl *ZLogger) getOtelFields(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	fields := make(map[string]interface{})

	// Add request ID if present
	if reqID := ctx.Value(ReqIDKey); reqID != nil {
		fields[LogIDKey] = reqID
	}

	// Add OpenTelemetry trace information if available
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID()
		spanID := span.SpanContext().SpanID()

		if traceID.IsValid() {
			fields["trace_id"] = traceID.String()
		}
		if spanID.IsValid() {
			fields["span_id"] = spanID.String()
		}

		// Add trace flags if needed
		traceFlags := span.SpanContext().TraceFlags()
		if traceFlags != 0 {
			fields["trace_flags"] = fmt.Sprintf("%02x", uint8(traceFlags))
		}
	}

	return fields
}

// getContextFields extracts fields from context for logging
func getContextFields(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	fields := make(map[string]interface{})

	// Add request ID if present
	if reqID := ctx.Value("request_id"); reqID != nil {
		fields["request_id"] = reqID
	}

	// Add OpenTelemetry trace information if available
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID()
		spanID := span.SpanContext().SpanID()

		if traceID.IsValid() {
			fields["trace_id"] = traceID.String()
		}
		if spanID.IsValid() {
			fields["span_id"] = spanID.String()
		}

		// Add trace flags if needed
		traceFlags := span.SpanContext().TraceFlags()
		if traceFlags != 0 {
			fields["trace_flags"] = fmt.Sprintf("%02x", uint8(traceFlags))
		}
	}

	return fields
}

// Implementation of Control interface methods
func (zl *ZLogger) SetLevel(level hertzlog.Level) {
	zl.level = level
	zl.logger = zl.logger.Level(toZerologLevel(level))
}

func (zl *ZLogger) SetOutput(w io.Writer) {
	// Rebuild logger with the same configuration but new output
	switch zl.format {
	case JSONFormat:
		// JSON format - default zerolog behavior with caller info
		zl.logger = zerolog.New(w).Level(toZerologLevel(zl.level)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	case ConsoleFormat:
		// Console format - human readable with RFC3339 time format, caller info and custom format
		consoleWriter := &zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.DateTime,
			FormatLevel: func(i interface{}) string {
				// Ensure full level name is shown instead of 3-letter abbreviation
				if ll, ok := i.(string); ok {
					return fmt.Sprintf("%-6s", strings.ToUpper(ll))
				}
				return fmt.Sprintf("%-6s", strings.ToUpper(fmt.Sprintf("%s", i)))
			},
		}
		zl.logger = zerolog.New(consoleWriter).Level(toZerologLevel(zl.level)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	default:
		// Default to console format
		consoleWriter := &zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
			FormatLevel: func(i interface{}) string {
				// Ensure full level name is shown instead of 3-letter abbreviation
				if ll, ok := i.(string); ok {
					return fmt.Sprintf("%-6s", strings.ToUpper(ll))
				}
				return fmt.Sprintf("%-6s", strings.ToUpper(fmt.Sprintf("%s", i)))
			},
		}
		zl.logger = zerolog.New(consoleWriter).Level(toZerologLevel(zl.level)).With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	}
}
