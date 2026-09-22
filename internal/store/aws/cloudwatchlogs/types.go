// Package logs provides CloudWatch Logs storage functionality for vorpalstacks.
package cloudwatchlogs

import (
	"time"
)

// The log-plane substrate records: the group and stream records, the ingested
// event forms and the chunk bookkeeping, with their constructors and the
// retention vocabulary. Every feature family's record types ride their
// operations files; the bounds register lives in limits.go.
// validRetentionDays is the set of retention values accepted by AWS
// CloudWatch Logs PutRetentionPolicy. Any value outside this set is
// rejected with InvalidParameterException.
var validRetentionDays = map[int32]bool{
	1: true, 3: true, 5: true, 7: true, 14: true, 30: true,
	60: true, 90: true, 120: true, 150: true, 180: true,
	365: true, 400: true, 545: true, 731: true, 1096: true,
	1827: true, 2192: true, 2557: true, 2922: true, 3288: true,
	3653: true,
}

// IsValidRetentionDays returns true if the given value is one of the
// allowed retention periods per the AWS CloudWatch Logs specification.
func IsValidRetentionDays(days int32) bool {
	return validRetentionDays[days]
}

// LogGroup represents a CloudWatch Logs log group.
type LogGroup struct {
	Name                      string    `json:"name"`
	ARN                       string    `json:"arn"`
	Region                    string    `json:"region"`
	AccountID                 string    `json:"accountId"`
	CreatedAt                 time.Time `json:"createdAt"`
	RetentionInDays           int32     `json:"retentionInDays,omitempty"`
	MetricFilterCount         int32     `json:"metricFilterCount"`
	StoredBytes               int64     `json:"storedBytes"`
	LogGroupClass             string    `json:"logGroupClass,omitempty"`
	KmsKeyId                  string    `json:"kmsKeyId,omitempty"`
	DeletionProtectionEnabled bool      `json:"deletionProtectionEnabled"`
	// DataProtectionStatus is the DescribeLogGroups display member
	// ("Displays whether this log group has a protection policy, or
	// whether it had one in the past"): empty until the first policy
	// write, ACTIVATED while one is stored, DELETED after its deletion.
	// The archive and disable transitions have no platform path and are
	// never produced.
	DataProtectionStatus string `json:"dataProtectionStatus,omitempty"`
	// BearerTokenAuthenticationEnabled mirrors the group's bearer token
	// switch on the group record itself, so the read surfaces (the
	// DescribeLogGroups member and the HTTP ingestion gate) observe it
	// with the record they already hold.
	BearerTokenAuthenticationEnabled bool              `json:"bearerTokenAuthenticationEnabled"`
	Tags                             map[string]string `json:"tags,omitempty"`
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
	Timestamp     int64  `json:"timestamp"`
	Message       string `json:"message"`
	IngestionTime int64  `json:"ingestionTime,omitempty"`
}

// OutputLogEvent represents an output log event.
type OutputLogEvent struct {
	Timestamp     int64  `json:"timestamp"`
	Message       string `json:"message"`
	IngestionTime int64  `json:"ingestionTime"`
	LogStreamName string `json:"logStreamName,omitempty"`
	// Ordinal is the read engine's per-event disambiguator (the chunk
	// identity plus the event's position inside the chunk). It carries
	// no wire meaning of its own: readers use it to mint members that
	// must stay unique across byte-identical duplicate events, whose
	// content-derived identity alone would collide.
	Ordinal string `json:"-"`
}

// ChunkMeta represents metadata for a log chunk.
type ChunkMeta struct {
	ChunkID      string `json:"chunkId"`
	LogGroupName string `json:"logGroupName"`
	LogStream    string `json:"logStream"`
	MinTs        int64  `json:"minTs"`
	MaxTs        int64  `json:"maxTs"`
	// MaxIngestionTs is the greatest ingestion time among the chunk's
	// entries — the index-level bound the delivery engine's late window
	// selects chunks by (events a backdated or mid-pass put landed below
	// a cursor that had already passed their timestamps). A record from
	// before the member existed decodes as zero and never qualifies,
	// which is the correct answer: its entries predate every delivery
	// cursor movement that matters.
	MaxIngestionTs int64  `json:"maxIngestionTs,omitempty"`
	EntryCount     int    `json:"entryCount"`
	ChunkPath      string `json:"chunkPath"`
	// ByteSize is the sum of the entry message lengths ingested into the
	// chunk — the exact quantity PutLogEvents added to the parent
	// LogGroup's StoredBytes — so removal paths decrement on the same
	// basis ingestion incremented. The chunk file on disk is compressed;
	// its size cannot stand in for this number.
	ByteSize int64 `json:"byteSize"`
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
