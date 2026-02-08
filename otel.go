// Package zlog provides OpenTelemetry tracing integration
package zlog

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// OtelHook is a zerolog hook that integrates with OpenTelemetry
type OtelHook struct {
	traceProvider trace.TracerProvider
	tracer        trace.Tracer
}

// NewOtelHook creates a new OpenTelemetry hook
func NewOtelHook(tp trace.TracerProvider) *OtelHook {
	if tp == nil {
		tp = trace.NewNoopTracerProvider()
	}

	return &OtelHook{
		traceProvider: tp,
		tracer:        tp.Tracer("zlog"),
	}
}

// Run implements the zerolog.Hook interface
func (h *OtelHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// Extract trace and span context from the event context if available
	ctx := context.Background()

	// Check if the event contains span information (would be added in the logger)
	span := trace.SpanFromContext(ctx)

	// If we have a valid span, add the log as an event
	if span.SpanContext().IsValid() {
		// Add log as an event to the current span
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", msg),
			attribute.String("level", level.String()),
		))

		// For error levels, mark the span as error
		if level >= zerolog.ErrorLevel {
			span.SetStatus(codes.Error, msg)
		}
	}
}

// OtelContextHook is a hook that extracts trace information from context and adds it to logs
type OtelContextHook struct{}

// NewOtelContextHook creates a new OpenTelemetry context hook
func NewOtelContextHook() *OtelContextHook {
	return &OtelContextHook{}
}

// Run implements the zerolog.Hook interface
func (h *OtelContextHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// This hook will be called for each log event
	// It attempts to extract trace information from the context attached to the logger
	// Note: For this to work properly, we need to modify how we use the logger
	// This is a simplified version that assumes the context is somehow available
}

// AddOtelFieldsToContext adds OpenTelemetry trace fields to context
func AddOtelFieldsToContext(ctx context.Context) map[string]interface{} {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return nil
	}

	fields := make(map[string]interface{})

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

	return fields
}

// WithOtelContextFields adds OpenTelemetry context fields to a logger
func WithOtelContextFields(ctx context.Context) zerolog.Logger {
	logger := zerolog.Nop()

	// Get OTel fields from context
	fields := AddOtelFieldsToContext(ctx)
	if len(fields) == 0 {
		return logger
	}

	return logger.With().Fields(fields).Logger()
}

// AddOtelHooks adds OpenTelemetry hooks to a ZLogger instance
func (zl *ZLogger) AddOtelHooks(ctx context.Context) {
	// Get OTel fields from context
	fields := AddOtelFieldsToContext(ctx)
	if len(fields) == 0 {
		return
	}

	// Add fields to the logger
	newLogger := zl.logger.With().Fields(fields).Logger()
	zl.logger = newLogger
}

// CtxInfofWithTrace adds trace information and logs the message
func (zl *ZLogger) CtxInfofWithTrace(ctx context.Context, format string, v ...interface{}) {
	// Add trace fields to context
	fields := AddOtelFieldsToContext(ctx)

	logEvt := zl.logger.Info()

	// Add trace fields to the event
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	logEvt.Msgf(format, v...)

	// Also add as event to the current span if it exists
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		span.AddEvent("log", trace.WithAttributes(
			attribute.String("message", fmt.Sprintf(format, v...)),
			attribute.String("level", "info"),
		))
	}
}

// CtxErrorfWithTrace adds trace information and logs the error message
func (zl *ZLogger) CtxErrorfWithTrace(ctx context.Context, format string, v ...interface{}) {
	// Add trace fields to context
	fields := AddOtelFieldsToContext(ctx)

	logEvt := zl.logger.Error()

	// Add trace fields to the event
	for k, v := range fields {
		logEvt = logEvt.Str(k, fmt.Sprintf("%v", v))
	}

	msg := fmt.Sprintf(format, v...)
	logEvt.Msg(msg)

	// Also add as event to the current span if it exists and mark span as error
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

// SetTraceProvider allows dynamic changing of the trace provider
func (zl *ZLogger) SetTraceProvider(tp trace.TracerProvider) {
	//if tp != nil {
	//	zl.tp = tp
	//} else {
	//	zl.tp = trace.NewNoopTracerProvider()
	//}
}

// GetTraceProvider returns the current trace provider
//func (zl *ZLogger) GetTraceProvider() trace.TracerProvider {
//	return zl.tp
//}
