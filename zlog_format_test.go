package zlog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	hertzlog "github.com/cloudwego/hertz/pkg/common/hlog"
)

func TestFormatTypes(t *testing.T) {
	// Test Console format
	consoleBuf := &bytes.Buffer{}
	consoleLogger := New(
		WithFormat(ConsoleFormat),
		WithOutput(consoleBuf),
		WithLevel(hertzlog.LevelInfo),
	)

	consoleLogger.Info("Test console message")
	consoleOutput := consoleBuf.String()

	// Console output should contain human-readable format, not pure JSON
	assert.Contains(t, consoleOutput, "Test console message")
	// Should not be pure JSON format (no curly braces at beginning)
	assert.NotRegexp(t, `^\s*\{`, consoleOutput)

	// Test JSON format
	jsonBuf := &bytes.Buffer{}
	jsonLogger := New(
		WithFormat(JSONFormat),
		WithOutput(jsonBuf),
		WithLevel(hertzlog.LevelInfo),
	)

	jsonLogger.Info("Test JSON message")
	jsonOutput := jsonBuf.String()

	// JSON output should be in JSON format
	assert.Contains(t, jsonOutput, `"message":"Test JSON message"`)
	assert.Regexp(t, `^\s*\{.*\}\s*$`, strings.TrimSpace(jsonOutput))

	// Test that format is correctly stored
	assert.Equal(t, ConsoleFormat, consoleLogger.format)
	assert.Equal(t, JSONFormat, jsonLogger.format)
}

func TestSetOutputPreservesFormat(t *testing.T) {
	buf := &bytes.Buffer{}

	// Create JSON logger
	jsonLogger := New(
		WithFormat(JSONFormat),
		WithOutput(buf),
		WithLevel(hertzlog.LevelInfo),
	)

	// Change output
	newBuf := &bytes.Buffer{}
	jsonLogger.SetOutput(newBuf)

	// Log message
	jsonLogger.Info("Test after output change")
	output := newBuf.String()

	// Should still be in JSON format
	assert.Contains(t, output, `"message":"Test after output change"`)
	assert.Regexp(t, `^\s*\{.*\}\s*$`, strings.TrimSpace(output))

	// Create console logger
	consoleBuf := &bytes.Buffer{}
	consoleLogger := New(
		WithFormat(ConsoleFormat),
		WithOutput(consoleBuf),
		WithLevel(hertzlog.LevelInfo),
	)

	// Change output
	newConsoleBuf := &bytes.Buffer{}
	consoleLogger.SetOutput(newConsoleBuf)

	// Log message
	consoleLogger.Info("Test after output change")
	consoleOutput := newConsoleBuf.String()

	// Should still be in console format
	assert.Contains(t, consoleOutput, "Test after output change")
	assert.NotRegexp(t, `^\s*\{`, consoleOutput)
}