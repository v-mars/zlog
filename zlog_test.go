package zlog

import (
	"context"
	"testing"
)

func TestZLogger(t *testing.T) {
	// Test basic logging
	logger := New()

	logger.Info("Test basic logging")
	logger.Infof("Test formatted logging: %s", "info")

	logger.Debug("Test debug logging")
	logger.Debugf("Test formatted debug: %d", 42)

	logger.Warn("Test warning logging")
	logger.Warnf("Test formatted warning: %t", true)

	// Test context logging
	ctx := context.WithValue(context.Background(), "request_id", "12345")
	logger.CtxInfof(ctx, "Test context logging: %s", "with context")

	t.Log("Basic logger functionality test passed")
}

func TestRotatingLogger(t *testing.T) {
	// Test rotating logger
	config := GetDefaultRotateConfig("test.log")
	rotatingLogger := NewRotatingLogger(config)

	rotatingLogger.Info("Test rotating logger")
	rotatingLogger.Infof("Test formatted rotating logger: %s", "info")

	// Manually rotate
	err := rotatingLogger.Rotate()
	if err != nil {
		t.Errorf("Failed to rotate log: %v", err)
	}

	t.Log("Rotating logger functionality test passed")
}

func TestHlogAdapter(t *testing.T) {
	// Test hlog adapter
	zlogger := New()
	adapter := NewHlogAdapter(zlogger)

	// Test that adapter methods work
	adapter.Info("Test hlog adapter")
	adapter.Infof("Test formatted hlog adapter: %s", "info")

	// Test context logging
	ctx := context.Background()
	adapter.CtxInfof(ctx, "Test context with adapter: %s", "context")

	t.Log("Hlog adapter functionality test passed")
}

func BenchmarkBasicLogging(b *testing.B) {
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message")
	}
}

func BenchmarkContextLogging(b *testing.B) {
	ctx := context.Background()
	logger := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.CtxInfof(ctx, "Benchmark context message: %d", i)
	}
}