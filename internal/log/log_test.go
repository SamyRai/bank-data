package log

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestStdLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := New(make(Fields), defaultMerge, WithWriter[Fields](&buf), WithMinLevel[Fields](LevelDebug))

	fields := Fields{"user_id": 123}
	logger.Info("test message", fields)

	var entry logEntry[Fields]
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if entry.Message != "test message" {
		t.Errorf("expected message 'test message', got '%s'", entry.Message)
	}
	if entry.Level != LevelInfo {
		t.Errorf("expected level 'info', got '%s'", entry.Level)
	}
	if entry.Fields["user_id"].(float64) != 123 {
		t.Errorf("expected field user_id 123, got %v", entry.Fields["user_id"])
	}
}

func TestStdLogger_MinLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(make(Fields), defaultMerge, WithWriter[Fields](&buf), WithMinLevel[Fields](LevelWarn))

	logger.Info("this should be filtered", nil)
	if buf.Len() > 0 {
		t.Errorf("expected no log entry, got %d bytes", buf.Len())
	}

	logger.Warn("this should be logged", nil)
	if buf.Len() == 0 {
		t.Error("expected log entry, got nothing")
	}
}

func TestStdLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger := New(make(Fields), defaultMerge, WithWriter[Fields](&buf))

	child := logger.WithFields(Fields{"component": "test"})
	child.Info("hello", Fields{"extra": "data"})

	var entry logEntry[Fields]
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to unmarshal log entry: %v", err)
	}

	if entry.Fields["component"] != "test" {
		t.Errorf("expected component field 'test', got '%v'", entry.Fields["component"])
	}
	if entry.Fields["extra"] != "data" {
		t.Errorf("expected extra field 'data', got '%v'", entry.Fields["extra"])
	}
}
