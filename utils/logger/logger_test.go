package logger_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/mannyfresh11/cloudflare-ddns/utils/logger"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := logger.Logger{
		Log: slog.New(handler),
	}
	logger.Info("Test message", slog.String("key", "value"))

	logOutput := buf.String()

	if !bytes.Contains([]byte(logOutput), []byte("INFO")) {
		t.Errorf("Expected log level 'INFO' not found in log output: %s", logOutput)
	}
	if !bytes.Contains([]byte(logOutput), []byte("Test message")) {
		t.Errorf("Expected log message 'Test message' not found in log output: %s", logOutput)
	}
	if !bytes.Contains([]byte(logOutput), []byte("key=value")) {
		t.Errorf("Expected key-value pair 'key=value' not found in log output: %s", logOutput)
	}
}
