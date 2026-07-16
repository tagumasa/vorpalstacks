// Package logs provides logging functionality for vorpalstacks.
package logs

import (
	"encoding/json"
	"time"
)

// Level represents the severity of a log message.
type Level int

// Level constants represent the severity levels for log entries.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the string representation of a Level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel converts a string representation to a Level.
func ParseLevel(s string) Level {
	switch s {
	case "DEBUG", "debug":
		return LevelDebug
	case "INFO", "info":
		return LevelInfo
	case "WARN", "warn", "WARNING", "warning":
		return LevelWarn
	case "ERROR", "error":
		return LevelError
	default:
		return LevelInfo
	}
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// LogEntry represents a single log entry.
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     Level                  `json:"level"`
	Message   string                 `json:"message"`
	Fields    []Field                `json:"fields,omitempty"`
	Source    string                 `json:"source,omitempty"`
	TenantID  string                 `json:"tenant_id,omitempty"`
	Region    string                 `json:"region,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Encode serialises the LogEntry to JSON.
func (e *LogEntry) Encode() ([]byte, error) {
	return json.Marshal(e)
}

// DecodeLogEntry deserialises a LogEntry from JSON data.
func DecodeLogEntry(data []byte) (*LogEntry, error) {
	var entry LogEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetField retrieves a field value by its key.
func (e *LogEntry) GetField(key string) (interface{}, bool) {
	for _, f := range e.Fields {
		if f.Key == key {
			return f.Value, true
		}
	}
	return nil, false
}

// WithFields returns a new LogEntry with the given fields appended.
func (e *LogEntry) WithFields(fields ...Field) *LogEntry {
	newEntry := *e
	newEntry.Fields = make([]Field, len(e.Fields), len(e.Fields)+len(fields))
	copy(newEntry.Fields, e.Fields)
	newEntry.Fields = append(newEntry.Fields, fields...)
	return &newEntry
}
