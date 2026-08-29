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
	// MaxLookupTables is the documented per-account, per-Region quota of
	// lookup tables.
	MaxLookupTables = 100
	// DefaultDescribeLookupTablesResults is the documented default of the
	// maxResults parameter of DescribeLookupTables.
	DefaultDescribeLookupTablesResults = 50
	// MaxDescribeLookupTablesResults is the documented maximum of the
	// maxResults parameter of DescribeLookupTables.
	MaxDescribeLookupTablesResults = 100
	// MaxLookupTableBodyBytes is the documented size ceiling of a lookup
	// table's CSV content (10 MB).
	MaxLookupTableBodyBytes = 10485760
	// MaxLookupTableNameLength is the documented maximum length of a lookup
	// table name.
	MaxLookupTableNameLength = 256
	// MaxLookupTableDescriptionLength is the documented maximum length of a
	// lookup table description.
	MaxLookupTableDescriptionLength = 1024
	// MaxKmsKeyIdLength is the documented maximum length of the kmsKeyId
	// parameter of lookup table and scheduled-query destination operations
	// (the shared KmsKeyId shape).
	MaxKmsKeyIdLength = 256
	// MaxLookupTableTags is the documented maximum number of tags attached
	// to one lookup table resource.
	MaxLookupTableTags = 50
)

// LookupTableNamePattern is the documented character set of lookup table
// names: alphanumeric characters and underscores.
const LookupTableNamePattern = `^[a-zA-Z0-9_]+$`

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
	Name                      string            `json:"name"`
	ARN                       string            `json:"arn"`
	Region                    string            `json:"region"`
	AccountID                 string            `json:"accountId"`
	CreatedAt                 time.Time         `json:"createdAt"`
	RetentionInDays           int32             `json:"retentionInDays,omitempty"`
	MetricFilterCount         int32             `json:"metricFilterCount"`
	StoredBytes               int64             `json:"storedBytes"`
	LogGroupClass             string            `json:"logGroupClass,omitempty"`
	KmsKeyId                  string            `json:"kmsKeyId,omitempty"`
	DeletionProtectionEnabled bool              `json:"deletionProtectionEnabled"`
	Tags                      map[string]string `json:"tags,omitempty"`
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
	Name                      string                 `json:"name"`
	LogGroupName              string                 `json:"logGroupName"`
	FilterPattern             string                 `json:"filterPattern"`
	MetricTransformations     []MetricTransformation `json:"metricTransformations"`
	ApplyOnTransformedLogs    bool                   `json:"applyOnTransformedLogs,omitempty"`
	FieldSelectionCriteria    string                 `json:"fieldSelectionCriteria,omitempty"`
	EmitSystemFieldDimensions []string               `json:"emitSystemFieldDimensions,omitempty"`
	CreatedAt                 time.Time              `json:"createdAt"`
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
	LogGroupName           string    `json:"logGroupName"`
	FilterName             string    `json:"filterName"`
	FilterPattern          string    `json:"filterPattern"`
	DestinationArn         string    `json:"destinationArn"`
	RoleArn                string    `json:"roleArn"`
	Distribution           string    `json:"distribution"`
	ApplyOnTransformedLogs bool      `json:"applyOnTransformedLogs,omitempty"`
	FieldSelectionCriteria string    `json:"fieldSelectionCriteria,omitempty"`
	EmitSystemFields       []string  `json:"emitSystemFields,omitempty"`
	CreationTime           time.Time `json:"creationTime"`
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

// ResourcePolicy represents a CloudWatch Logs resource policy.
type ResourcePolicy struct {
	PolicyName      string `json:"policyName"`
	PolicyDocument  string `json:"policyDocument"`
	ResourceArn     string `json:"resourceArn,omitempty"`
	PolicyScope     string `json:"policyScope,omitempty"`
	RevisionId      string `json:"revisionId,omitempty"`
	LastUpdatedTime int64  `json:"lastUpdatedTime"`
}

// AccountPolicy represents a CloudWatch Logs account-level policy.
type AccountPolicy struct {
	PolicyName        string `json:"policyName"`
	PolicyDocument    string `json:"policyDocument"`
	PolicyType        string `json:"policyType"`
	Scope             string `json:"scope,omitempty"`
	SelectionCriteria string `json:"selectionCriteria,omitempty"`
	AccountId         string `json:"accountId,omitempty"`
	LastUpdatedTime   int64  `json:"lastUpdatedTime"`
}

// DataProtectionPolicy represents a CloudWatch Logs data protection policy.
type DataProtectionPolicy struct {
	LogGroupIdentifier string `json:"logGroupIdentifier"`
	PolicyDocument     string `json:"policyDocument"`
	LastUpdatedTime    int64  `json:"lastUpdatedTime"`
}

// QueryDefinition represents a saved CloudWatch Logs Insights query definition.
type QueryDefinition struct {
	QueryDefinitionId string                 `json:"queryDefinitionId"`
	Name              string                 `json:"name"`
	QueryString       string                 `json:"queryString"`
	LogGroupNames     []string               `json:"logGroupNames,omitempty"`
	QueryLanguage     string                 `json:"queryLanguage,omitempty"`
	Parameters        map[string]interface{} `json:"parameters,omitempty"`
	LastModified      int64                  `json:"lastModified"`
}

// ExportTask represents a CloudWatch Logs export-to-S3 task.
type ExportTask struct {
	TaskId              string                 `json:"taskId"`
	TaskName            string                 `json:"taskName"`
	LogGroupName        string                 `json:"logGroupName"`
	LogStreamNamePrefix string                 `json:"logStreamNamePrefix,omitempty"`
	From                int64                  `json:"from"`
	To                  int64                  `json:"to"`
	Destination         string                 `json:"destination"`
	DestinationPrefix   string                 `json:"destinationPrefix,omitempty"`
	Status              string                 `json:"status"`
	StatusMessage       string                 `json:"statusMessage,omitempty"`
	ExecutionInfo       map[string]interface{} `json:"executionInfo,omitempty"`
	CreationTime        int64                  `json:"creationTime"`
}

// ImportTask represents a CloudWatch Logs import-from-S3 task.
type ImportTask struct {
	ImportId             string                 `json:"importId"`
	ImportSourceArn      string                 `json:"importSourceArn"`
	ImportRoleArn        string                 `json:"importRoleArn,omitempty"`
	LogGroupName         string                 `json:"logGroupName"`
	ImportStatus         string                 `json:"importStatus"`
	ImportDestinationArn string                 `json:"importDestinationArn,omitempty"`
	ImportStatistics     map[string]interface{} `json:"importStatistics,omitempty"`
	ImportFilter         map[string]interface{} `json:"importFilter,omitempty"`
	ErrorMessage         string                 `json:"errorMessage,omitempty"`
	CreationTime         int64                  `json:"creationTime"`
	LastUpdatedTime      int64                  `json:"lastUpdatedTime"`
}

// ScheduledQuery represents a scheduled CloudWatch Logs Insights query.
type ScheduledQuery struct {
	Id                       string                 `json:"id"`
	Name                     string                 `json:"name"`
	Description              string                 `json:"description,omitempty"`
	QueryString              string                 `json:"queryString"`
	QueryLanguage            string                 `json:"queryLanguage,omitempty"`
	LogGroupIdentifiers      []string               `json:"logGroupIdentifiers,omitempty"`
	ScheduleExpression       string                 `json:"scheduleExpression"`
	ScheduleType             string                 `json:"scheduleType,omitempty"`
	State                    string                 `json:"state"`
	ExecutionRoleArn         string                 `json:"executionRoleArn,omitempty"`
	Timezone                 string                 `json:"timezone,omitempty"`
	StartTimeOffset          int64                  `json:"startTimeOffset,omitempty"`
	EndTimeOffset            int64                  `json:"endTimeOffset,omitempty"`
	ScheduleStartTime        int64                  `json:"scheduleStartTime,omitempty"`
	ScheduleEndTime          int64                  `json:"scheduleEndTime,omitempty"`
	DestinationConfiguration map[string]interface{} `json:"destinationConfiguration,omitempty"`
	// LastExecutionStatus carries the outcome of the most recent
	// execution on the wire (Running, InvalidQuery, Complete, Failed,
	// Timeout per the service model).
	LastExecutionStatus string `json:"lastExecutionStatus,omitempty"`
	CreationTime        int64  `json:"creationTime"`
	LastUpdatedTime     int64  `json:"lastUpdatedTime"`
	LastTriggeredTime   int64  `json:"lastTriggeredTime,omitempty"`
	// LastExecutedBoundary is an internal marker holding the schedule
	// boundary of the most recent executed occurrence. It is the
	// deduplication truth across restarts and never surfaces on the
	// wire; lastTriggeredTime remains the execution clock.
	LastExecutedBoundary int64             `json:"lastExecutedBoundary,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty"`
}

// Wire values of the ExecutionStatus enum (Running, InvalidQuery,
// Complete, Failed, Timeout) carried by the lastExecutionStatus member
// of the scheduled query shapes.
const (
	ScheduledQueryStatusComplete = "Complete"
	ScheduledQueryStatusFailed   = "Failed"
)

// ScheduledQueryDestination records the delivery outcome of one destination
// of a scheduled query execution, reported through GetScheduledQueryHistory.
type ScheduledQueryDestination struct {
	DestinationType       string `json:"destinationType"`
	DestinationIdentifier string `json:"destinationIdentifier"`
	Status                string `json:"status"`
	ProcessedIdentifier   string `json:"processedIdentifier,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
}

// ScheduledQueryExecution represents a single execution of a scheduled query.
type ScheduledQueryExecution struct {
	ScheduledQueryId string                       `json:"scheduledQueryId"`
	QueryId          string                       `json:"queryId"`
	Destinations     []*ScheduledQueryDestination `json:"destinations,omitempty"`
	TriggerTime      int64                        `json:"triggerTime"`
	Status           string                       `json:"status"`
	ErrorMessage     string                       `json:"errorMessage,omitempty"`
	RecordsScanned   int64                        `json:"recordsScanned"`
	RecordsMatched   int64                        `json:"recordsMatched"`
}

// LookupTable stores the reference data the lookup and cidrlookup query
// commands enrich events with. TableBody holds the CSV content including
// the header row; TableFields mirrors the header and RecordsCount counts
// the data rows. When the table is encrypted with a customer-managed KMS
// key, TableBody is empty and EncryptedBody, EncryptedDataKey and
// ContentNonce hold the envelope-encrypted content instead.
type LookupTable struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	TableBody        string            `json:"tableBody,omitempty"`
	TableFields      []string          `json:"tableFields,omitempty"`
	RecordsCount     int64             `json:"recordsCount"`
	SizeBytes        int64             `json:"sizeBytes"`
	KmsKeyId         string            `json:"kmsKeyId,omitempty"`
	EncryptedBody    []byte            `json:"encryptedBody,omitempty"`
	EncryptedDataKey []byte            `json:"encryptedDataKey,omitempty"`
	ContentNonce     []byte            `json:"contentNonce,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	CreationTime     int64             `json:"creationTime"`
	LastUpdatedTime  int64             `json:"lastUpdatedTime"`
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
