package logger

import (
	"context"
	"log/slog"
	"testing"
)

func TestNewLogger_Levels(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected slog.Level
	}{
		{"DebugLevel", "debug", slog.LevelDebug},
		{"InfoLevel", "info", slog.LevelInfo},
		{"WarnLevel", "warn", slog.LevelWarn},
		{"ErrorLevel", "error", slog.LevelError},
		{"DefaultLevel", "unknown", slog.LevelInfo},
		{"EmptyLevel", "", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.level)
			if logger.Handler().Enabled(context.TODO(), tt.expected) != true {
				t.Errorf("Expected level %v, but logger is not enabled for this level", tt.expected)
			}
		})
	}
}
