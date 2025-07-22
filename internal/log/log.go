// Package log provides a production-ready, dependency-free structured logger.
package log

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

// Level type for log levels.
type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelDebug Level = "debug"
)

// LevelPriority maps log levels to their priority for filtering.
var LevelPriority = map[Level]int{
	LevelDebug: 1,
	LevelInfo:  2,
	LevelWarn:  3,
	LevelError: 4,
}

// Logger defines the interface for structured logging with generic fields.
type Logger[F any] interface {
	Info(msg string, fields F)
	Warn(msg string, fields F)
	Error(msg string, fields F)
	Debug(msg string, fields F)
	WithFields(fields F) Logger[F]
}

// logEntry is the structure for each log message.
type logEntry[F any] struct {
	Timestamp string `json:"ts"`
	Level     Level  `json:"level"`
	Message   string `json:"msg"`
	Fields    F      `json:"fields,omitempty"`
}

// StdLogger is a production-ready, dependency-free implementation of Logger.
type StdLogger[F any] struct {
	mu       sync.Mutex // protects fields and config
	fields   F
	merge    func(a, b F) F // merge function for fields
	minLevel Level
	writer   io.Writer
}

// New creates a new StdLogger instance with a merge function for fields, minimum log level, and output writer.
func New[F any](zero F, merge func(a, b F) F, opts ...Option[F]) *StdLogger[F] {
	logger := &StdLogger[F]{
		fields:   zero,
		merge:    merge,
		minLevel: LevelDebug, // default: log everything
		writer:   os.Stdout,  // default: stdout
	}
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

// Option configures a StdLogger.
type Option[F any] func(*StdLogger[F])

// WithMinLevel sets the minimum log level for the logger.
func WithMinLevel[F any](level Level) Option[F] {
	return func(l *StdLogger[F]) {
		l.minLevel = level
	}
}

// WithWriter sets the output writer for the logger.
func WithWriter[F any](w io.Writer) Option[F] {
	return func(l *StdLogger[F]) {
		if w != nil {
			l.writer = w
		}
	}
}

// WithFields returns a new Logger with additional context fields.
func (l *StdLogger[F]) WithFields(fields F) Logger[F] {
	l.mu.Lock()
	defer l.mu.Unlock()
	merged := l.merge(l.fields, fields)
	return &StdLogger[F]{
		fields:   merged,
		merge:    l.merge,
		minLevel: l.minLevel,
		writer:   l.writer,
	}
}

func (l *StdLogger[F]) logf(level Level, msg string, fields F) {
	l.mu.Lock()
	minLevel := l.minLevel
	writer := l.writer
	l.mu.Unlock()
	if LevelPriority[level] < LevelPriority[minLevel] {
		return // filtered out
	}
	allFields := l.merge(l.fields, fields)
	entry := logEntry[F]{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     level,
		Message:   msg,
		Fields:    allFields,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		if _, serr := os.Stderr.Write([]byte(`{"level":"error","msg":"failed to marshal log entry","err":"` + err.Error() + `"}\n`)); serr != nil {
			_ = serr // intentionally ignored: cannot log further if stderr write fails
		}
		return
	}
	if _, werr := writer.Write(append([]byte(string(b)), '\n')); werr != nil {
		if _, serr := os.Stderr.Write([]byte(`{"level":"error","msg":"failed to write log entry","err":"` + werr.Error() + `"}\n`)); serr != nil {
			_ = serr // intentionally ignored: cannot log further if stderr write fails
		}
	}
}

// Info logs an info message with optional fields.
func (l *StdLogger[F]) Info(msg string, fields F) {
	l.logf(LevelInfo, msg, fields)
}

// Warn logs a warning message with optional fields.
func (l *StdLogger[F]) Warn(msg string, fields F) {
	l.logf(LevelWarn, msg, fields)
}

// Error logs an error message with optional fields.
func (l *StdLogger[F]) Error(msg string, fields F) {
	l.logf(LevelError, msg, fields)
}

// Debug logs a debug message with optional fields.
func (l *StdLogger[F]) Debug(msg string, fields F) {
	l.logf(LevelDebug, msg, fields)
}

// Fields is a type alias for map[string]interface{} used for structured log fields.
type Fields = map[string]interface{}

var defaultMerge = func(a, b Fields) Fields {
	merged := make(Fields, len(a)+len(b))
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

var Default = New(make(Fields), defaultMerge)

// Info logs an info message using the default logger.
func Info(msg string, fields Fields) {
	Default.Info(msg, fields)
}

// Warn logs a warning message using the default logger.
func Warn(msg string, fields Fields) {
	Default.Warn(msg, fields)
}

// Error logs an error message using the default logger.
func Error(msg string, fields Fields) {
	Default.Error(msg, fields)
}

// Debug logs a debug message using the default logger.
func Debug(msg string, fields Fields) {
	Default.Debug(msg, fields)
}

// TODO(context: log package, priority: low, effort: 2h): Add support for log level filtering and output destination configuration.
