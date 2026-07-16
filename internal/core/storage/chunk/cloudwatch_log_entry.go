package chunk

import (
	"encoding/json"
	"time"
)

// CloudWatchLogEntry represents a log entry in CloudWatch Logs format.
type CloudWatchLogEntry struct {
	LogGroupName  string `json:"logGroupName"`
	LogStreamName string `json:"logStreamName"`
	Ts            int64  `json:"timestamp"`
	Msg           string `json:"message"`
	IngestionTs   int64  `json:"ingestionTime"`
	EventID       string `json:"eventId"`
}

// Timestamp returns the timestamp of the log entry in nanoseconds.
func (e *CloudWatchLogEntry) Timestamp() int64 {
	return e.Ts * 1e6
}

// Message returns the log entry message as JSON.
func (e *CloudWatchLogEntry) Message() []byte {
	data, _ := json.Marshal(e)
	return data
}

// Time returns the timestamp of the log entry as a time.Time value.
func (e *CloudWatchLogEntry) Time() time.Time {
	return time.UnixMilli(e.Ts)
}

// IngestionTime returns the ingestion time of the log entry as a time.Time value.
func (e *CloudWatchLogEntry) IngestionTime() time.Time {
	return time.UnixMilli(e.IngestionTs)
}

// Tags returns the tags associated with the log entry.
func (e *CloudWatchLogEntry) Tags() map[string]string {
	return map[string]string{
		"log_group":  e.LogGroupName,
		"log_stream": e.LogStreamName,
	}
}
