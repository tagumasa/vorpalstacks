// Package logs provides CloudWatch Logs storage functionality for vorpalstacks.
package cloudwatchlogs

import (
	"time"
)

const (
	// MaxChunkSize is the maximum number of log entries per chunk.
	MaxChunkSize = 10000
	// MaxRetentionDays is the maximum retention period in days.
	MaxRetentionDays = 3653
	// DefaultRetentionDays is the default retention period in days.
	DefaultRetentionDays = 30
)

// LogGroup represents a CloudWatch Logs log group.
type LogGroup struct {
	Name              string            `json:"name"`
	ARN               string            `json:"arn"`
	Region            string            `json:"region"`
	AccountID         string            `json:"accountId"`
	CreatedAt         time.Time         `json:"createdAt"`
	RetentionInDays   int32             `json:"retentionInDays,omitempty"`
	MetricFilterCount int32             `json:"metricFilterCount"`
	StoredBytes       int64             `json:"storedBytes"`
	Tags              map[string]string `json:"tags,omitempty"`
}

// LogStream represents a CloudWatch Logs log stream.
type LogStream struct {
	Name                string    `json:"name"`
	LogGroupName        string    `json:"logGroupName"`
	ARN                 string    `json:"arn"`
	CreatedAt           time.Time `json:"createdAt"`
	FirstEventTs        int64     `json:"firstEventTs,omitempty"`
	LastEventTs         int64     `json:"lastEventTs,omitempty"`
	LastIngestionTs     int64     `json:"lastIngestionTs,omitempty"`
	UploadSequenceToken string    `json:"uploadSequenceToken,omitempty"`
}

// LogEntry represents a single log event entry.
type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// OutputLogEvent represents an output log event.
type OutputLogEvent struct {
	Timestamp     int64  `json:"timestamp"`
	Message       string `json:"message"`
	IngestionTime int64  `json:"ingestionTime"`
}

// ChunkMeta represents metadata for a log chunk.
type ChunkMeta struct {
	ChunkID    string `json:"chunkId"`
	LogGroup   string `json:"logGroup"`
	LogStream  string `json:"logStream"`
	MinTs      int64  `json:"minTs"`
	MaxTs      int64  `json:"maxTs"`
	EntryCount int    `json:"entryCount"`
	ChunkPath  string `json:"chunkPath"`
}

// MetricFilter represents a CloudWatch Logs metric filter.
type MetricFilter struct {
	Name                  string                 `json:"name"`
	LogGroupName          string                 `json:"logGroupName"`
	FilterPattern         string                 `json:"filterPattern"`
	MetricTransformations []MetricTransformation `json:"metricTransformations"`
	CreatedAt             time.Time              `json:"createdAt"`
}

// MetricTransformation represents a metric transformation for a metric filter.
type MetricTransformation struct {
	MetricName      string  `json:"metricName"`
	MetricNamespace string  `json:"metricNamespace"`
	MetricValue     string  `json:"metricValue"`
	DefaultValue    float64 `json:"defaultValue,omitempty"`
	DefaultValueSet bool    `json:"defaultValueSet,omitempty"`
}

// SubscriptionFilter represents a CloudWatch Logs subscription filter.
type SubscriptionFilter struct {
	LogGroupName   string    `json:"logGroupName"`
	FilterName     string    `json:"filterName"`
	FilterPattern  string    `json:"filterPattern"`
	DestinationArn string    `json:"destinationArn"`
	RoleArn        string    `json:"roleArn"`
	Distribution   string    `json:"distribution"`
	CreationTime   time.Time `json:"creationTime"`
}

// Destination represents a CloudWatch Logs destination (cross-account).
type Destination struct {
	Name         string            `json:"name"`
	ARN          string            `json:"arn"`
	RoleArn      string            `json:"roleArn"`
	TargetArn    string            `json:"targetArn"`
	AccessPolicy string            `json:"accessPolicy"`
	CreationTime int64             `json:"creationTime"`
	Tags         map[string]string `json:"tags,omitempty"`
}

// NewLogGroup creates a new CloudWatch Logs log group.
func NewLogGroup(name, region, accountID string) *LogGroup {
	return &LogGroup{
		Name:            name,
		Region:          region,
		AccountID:       accountID,
		CreatedAt:       time.Now().UTC(),
		RetentionInDays: 0,
		Tags:            make(map[string]string),
	}
}

// NewLogStream creates a new CloudWatch Logs log stream.
func NewLogStream(name, logGroupName string) *LogStream {
	return &LogStream{
		Name:         name,
		LogGroupName: logGroupName,
		CreatedAt:    time.Now().UTC(),
	}
}

// SetRetention sets the retention period for the log group in days.
func (lg *LogGroup) SetRetention(days int32) {
	if days == 0 {
		lg.RetentionInDays = 0
	} else if days > 0 && days <= MaxRetentionDays {
		lg.RetentionInDays = days
	}
}

// UpdateEventTimestamps updates the first and last event timestamps for the log stream.
func (cs *LogStream) UpdateEventTimestamps(firstTs, lastTs int64) {
	if cs.FirstEventTs == 0 || firstTs < cs.FirstEventTs {
		cs.FirstEventTs = firstTs
	}
	if lastTs > cs.LastEventTs {
		cs.LastEventTs = lastTs
	}
}
