package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	tests := []struct {
		level    string
		expected logrus.Level
	}{
		{"debug", logrus.DebugLevel},
		{"info", logrus.InfoLevel},
		{"warn", logrus.WarnLevel},
		{"error", logrus.ErrorLevel},
		{"invalid", logrus.InfoLevel}, // Default fallback
		{"", logrus.InfoLevel},        // Default fallback
	}

	for _, test := range tests {
		logger := New(test.level)

		if logger == nil {
			t.Fatal("Expected logger to be created, got nil")
		}

		logrusLogger, ok := logger.(*LogrusLogger)
		if !ok {
			t.Fatal("Expected LogrusLogger type")
		}

		if logrusLogger.Logger.GetLevel() != test.expected {
			t.Errorf("Expected level %v for input %s, got %v", test.expected, test.level, logrusLogger.Logger.GetLevel())
		}
	}
}

func TestLogrusLogger_Interface(t *testing.T) {
	logger := New("debug")
	logrusLogger := logger.(*LogrusLogger)

	// Test that the logger implements the Logger interface
	var _ Logger = logger

	// Test all interface methods exist
	logger.Debug("test debug")
	logger.Debugf("test debug %s", "format")
	logger.Info("test info")
	logger.Infof("test info %s", "format")
	logger.Warn("test warn")
	logger.Warnf("test warn %s", "format")
	logger.Error("test error")
	logger.Errorf("test error %s", "format")

	// Test that the underlying logrus logger is accessible
	if logrusLogger.Logger == nil {
		t.Fatal("Expected underlying logrus logger to be set")
	}
}

func TestLogrusLogger_SetOutput(t *testing.T) {
	logger := New("debug")
	logrusLogger := logger.(*LogrusLogger)

	// Create a buffer to capture output
	var buf bytes.Buffer
	logrusLogger.SetOutput(&buf)

	// Log a message
	logger.Info("test message")

	// Check that the message was written to the buffer
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}

	if !strings.Contains(output, `"level":"info"`) {
		t.Errorf("Expected output to contain 'level':'info', got: %s", output)
	}
}

func TestLogrusLogger_JSONFormat(t *testing.T) {
	logger := New("debug")
	logrusLogger := logger.(*LogrusLogger)

	// Create a buffer to capture output
	var buf bytes.Buffer
	logrusLogger.SetOutput(&buf)

	// Log a message
	logger.Info("test json message")

	// Check that the output is in JSON format
	output := buf.String()
	if !strings.Contains(output, `"msg":"test json message"`) {
		t.Errorf("Expected JSON output with message, got: %s", output)
	}

	if !strings.Contains(output, `"level":"info"`) {
		t.Errorf("Expected JSON output with level, got: %s", output)
	}

	if !strings.Contains(output, `"time"`) {
		t.Errorf("Expected JSON output with timestamp, got: %s", output)
	}
}

func TestLogrusLogger_DifferentLevels(t *testing.T) {
	logger := New("debug")
	logrusLogger := logger.(*LogrusLogger)

	// Create a buffer to capture output
	var buf bytes.Buffer
	logrusLogger.SetOutput(&buf)

	// Test different log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	// Check that all levels are present
	levels := []string{"debug", "info", "warning", "error"} // Note: logrus uses "warning" not "warn"
	for _, level := range levels {
		if !strings.Contains(output, `"level":"`+level+`"`) {
			t.Errorf("Expected output to contain level %s, got: %s", level, output)
		}
	}
}

func TestLogrusLogger_SetOutput_Stdout(t *testing.T) {
	logger := New("debug")
	logrusLogger := logger.(*LogrusLogger)

	// Test setting output to stdout
	logrusLogger.SetOutput(os.Stdout)

	// This should not panic
	logger.Info("test stdout message")
}

func TestLogrusLogger_Formatting(t *testing.T) {
	logger := New("debug")
	logrusLogger := logger.(*LogrusLogger)

	// Create a buffer to capture output
	var buf bytes.Buffer
	logrusLogger.SetOutput(&buf)

	// Test formatted logging
	logger.Infof("User %s logged in with ID %d", "john", 123)

	output := buf.String()
	if !strings.Contains(output, "User john logged in with ID 123") {
		t.Errorf("Expected formatted message, got: %s", output)
	}
}

func TestLogrusLogger_MultipleMessages(t *testing.T) {
	logger := New("info")
	logrusLogger := logger.(*LogrusLogger)

	// Create a buffer to capture output
	var buf bytes.Buffer
	logrusLogger.SetOutput(&buf)

	// Log multiple messages
	logger.Info("message 1")
	logger.Info("message 2")
	logger.Info("message 3")

	output := buf.String()

	// Count the number of log entries (each should be on a separate line)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 log entries, got %d", len(lines))
	}

	// Check that each message is present
	for i := 1; i <= 3; i++ {
		expected := fmt.Sprintf("message %d", i)
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}
