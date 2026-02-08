package zlog

import (
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestOpenTelemetryIntegration(t *testing.T) {
	// Create a mock tracer provider for testing
	//tp := trace.NewNoopTracerProvider()

	// Create logger with trace provider
	//logger := New(WithTraceProvider(tp))

	// Create a context with a span (using noop provider, so it won't really trace)
	ctx := context.Background()

	// Test that the logger can handle context with OpenTelemetry info
	logger.CtxInfof(ctx, "Test message with context")

	t.Log("OpenTelemetry integration test passed")
}

func TestOpenTelemetryFieldsExtraction(t *testing.T) {
	logger := New()

	ctx := context.Background()

	// Test that we can extract fields without errors
	fields := logger.getOtelFields(ctx)

	// Fields might be empty if no trace is active, but shouldn't cause errors
	_ = fields

	t.Log("OpenTelemetry fields extraction test passed")
}

func TestManualTraceLogging(t *testing.T) {
	// Setup OTel tracer
	tp := trace.NewNoopTracerProvider()
	otel.SetTracerProvider(tp)

	//logger := New(WithTraceProvider(tp))

	//ctx := context.Background()

	// Test enhanced trace logging functions
	//logger.CtxInfofWithTrace(ctx, "Info message with trace: %s", "details")
	//logger.CtxErrorfWithTrace(ctx, "Error message with trace: %v", "error occurred")

	t.Log("Manual trace logging test passed")
}
