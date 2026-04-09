// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

// Logger defines the interface for logging.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
	WithError(err error) Logger
	WithSource(source string) Logger
}

// Config holds configuration for the logger.
type Config struct {
	Level      Level
	TenantID   string
	Region     string
	Output     io.Writer
	Store      *LogStore
	ShowSource bool
}

// logger is the internal implementation of the Logger interface.
type logger struct {
	mu       sync.Mutex
	config   *Config
	fields   []Field
	source   string
	logCh    chan *LogEntry
	logWg    sync.WaitGroup
	stopOnce sync.Once
}

// NewLogger creates a new logger with the given configuration.
func NewLogger(cfg *Config) Logger {
	if cfg == nil {
		cfg = &Config{Level: LevelInfo}
	}
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}
	l := &logger{
		config: cfg,
		logCh:  make(chan *LogEntry, 10000),
	}
	if cfg.Store != nil {
		l.logWg.Add(1)
		go l.logWorker()
	}
	return l
}

// logWorker asynchronously writes log entries to the persistent store.
func (l *logger) logWorker() {
	defer l.logWg.Done()
	for entry := range l.logCh {
		if err := l.config.Store.Write(entry); err != nil {
			log.Printf("Failed to write log entry to store: %v", err)
		}
	}
}

// Close stops the log worker and waits for pending entries to be flushed.
func (l *logger) Close() {
	l.stopOnce.Do(func() {
		close(l.logCh)
	})
	l.logWg.Wait()
}

// log writes a log entry at the given level if it meets the configured threshold.
func (l *logger) log(level Level, msg string, fields ...Field) {
	if level < l.config.Level {
		return
	}

	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		TenantID:  l.config.TenantID,
		Region:    l.config.Region,
	}

	if l.config.ShowSource {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			entry.Source = fmt.Sprintf("%s:%d", file, line)
		}
	}

	if l.source != "" {
		entry.Source = l.source
	}

	allFields := make([]Field, 0, len(l.fields)+len(fields))
	allFields = append(allFields, l.fields...)
	allFields = append(allFields, fields...)
	entry.Fields = allFields

	l.mu.Lock()
	l.writeToOutput(entry)
	l.mu.Unlock()

	if l.config.Store != nil {
		select {
		case l.logCh <- entry:
		default:
			log.Printf("Log store channel full, dropping entry")
		}
	}
}

// writeToOutput formats and writes a log entry to the configured output writer.
func (l *logger) writeToOutput(entry *LogEntry) {
	fmt.Fprintf(l.config.Output, "[%s] %s %s",
		entry.Timestamp.Format("2006-01-02 15:04:05.000"),
		entry.Level.String(),
		entry.Message,
	)

	for _, f := range entry.Fields {
		fmt.Fprintf(l.config.Output, " %s=%v", f.Key, f.Value)
	}

	if entry.Source != "" {
		fmt.Fprintf(l.config.Output, " source=%s", entry.Source)
	}

	fmt.Fprintln(l.config.Output)
}

// Debug logs a message at debug level.
//
// Parameters:
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (l *logger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

// Info logs a message at info level.
//
// Parameters:
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (l *logger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

// Warn logs a message at warn level.
//
// Parameters:
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (l *logger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

// Error logs a message at error level.
//
// Parameters:
//   - msg: The message to log
//   - fields: Optional fields to add to the log entry
func (l *logger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

// With returns a new logger with the given fields added.
//
// Parameters:
//   - fields: Fields to add to the logger
//
// Returns:
//   - Logger: A new logger with the added fields
func (l *logger) With(fields ...Field) Logger {
	return &logger{
		config: l.config,
		fields: append(l.fields, fields...),
		source: l.source,
	}
}

// WithError returns a new logger with the given error added as a field.
//
// Parameters:
//   - err: The error to add (can be nil)
//
// Returns:
//   - Logger: A new logger with the error field
func (l *logger) WithError(err error) Logger {
	if err == nil {
		return l.With(Field{Key: "error", Value: nil})
	}
	return l.With(Field{Key: "error", Value: err.Error()})
}

// WithSource returns a new logger with the given source added.
//
// Parameters:
//   - source: The source to add
//
// Returns:
//   - Logger: A new logger with the source field
func (l *logger) WithSource(source string) Logger {
	return &logger{
		config: l.config,
		fields: l.fields,
		source: source,
	}
}

// String creates a string field for logging.
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Int creates an integer field for logging.
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 creates a 64-bit integer field for logging.
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 creates a floating point field for logging.
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool creates a boolean field for logging.
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Err creates an error field for logging.
func Err(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}

// Any creates a field with any value for logging.
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

var (
	defaultLogger Logger
	defaultOnce   sync.Once
)

// InitDefault initialises the default logger with the given configuration.
func InitDefault(cfg *Config) {
	defaultOnce.Do(func() {
		defaultLogger = NewLogger(cfg)
	})
}

// Default returns the default logger, creating one if it does not exist.
func Default() Logger {
	defaultOnce.Do(func() {
		defaultLogger = NewLogger(nil)
	})
	return defaultLogger
}

// Debug logs a message at debug level using the default logger.
func Debug(msg string, fields ...Field) { Default().Debug(msg, fields...) }

// Info logs a message at info level using the default logger.
func Info(msg string, fields ...Field) { Default().Info(msg, fields...) }

// Warn logs a message at warning level using the default logger.
func Warn(msg string, fields ...Field) { Default().Warn(msg, fields...) }

// Error logs a message at error level using the default logger.
func Error(msg string, fields ...Field) { Default().Error(msg, fields...) }
