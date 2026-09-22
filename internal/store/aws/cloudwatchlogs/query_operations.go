package cloudwatchlogs

import (
	"strings"
	"sync"
	"time"
)

// --- Query Definition ---

func (s *Store) PutQueryDefinitionEntry(qd *QueryDefinition) error {
	// Seconds, not milliseconds: the member is a Timestamp shape whose
	// wire format is epoch seconds, matching the documented examples.
	qd.LastModified = time.Now().UTC().Unix()
	return s.Put(s.queryDefinitionKey(qd.QueryDefinitionId), qd)
}

func (s *Store) DeleteQueryDefinitionEntry(id string) error {
	key := s.queryDefinitionKey(id)
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

// GetQueryDefinitionEntry loads one saved query definition by its id, or
// the store's missing-resource sentinel when no record carries it.
func (s *Store) GetQueryDefinitionEntry(id string) (*QueryDefinition, error) {
	return getJSONRecord[QueryDefinition](s, s.queryDefinitionKey(id), ErrResourceNotFound)
}

func (s *Store) ListQueryDefinitions(namePrefix string) ([]*QueryDefinition, error) {
	return listJSONRecords[QueryDefinition](s, keyPrefixQueryDefinition, "query definition record", func(qd *QueryDefinition) bool {
		return namePrefix == "" || strings.HasPrefix(qd.Name, namePrefix)
	})
}

// --- Scheduled Query ---

func (s *Store) PutScheduledQuery(sq *ScheduledQuery) error {
	sq.LastUpdatedTime = time.Now().UTC().UnixMilli()
	return s.Put(s.scheduledQueryKey(sq.Id), sq)
}

func (s *Store) GetScheduledQuery(id string) (*ScheduledQuery, error) {
	return getJSONRecord[ScheduledQuery](s, s.scheduledQueryKey(id), ErrResourceNotFound)
}

// ScheduledQueryIdByName resolves a scheduled query's name to its id
// key, or "" when no record carries the name: the operations' documented
// Identifier form is ARN-or-name while the records are keyed by the
// opaque id alone.
func (s *Store) ScheduledQueryIdByName(name string) (string, error) {
	queries, err := s.ListScheduledQueries("")
	if err != nil {
		return "", err
	}
	for _, sq := range queries {
		if sq.Name == name {
			return sq.Id, nil
		}
	}
	return "", nil
}

// scheduledQueryRecordWriteMu serialises read-modify-write cycles on
// scheduled query records. The delivery worker and the API handlers
// write these records concurrently: a full-record put from either side
// must not interleave with the other's read-modify-write, and a write
// racing a delete must not resurrect the deleted record.
var scheduledQueryRecordWriteMu sync.Mutex

// MutateScheduledQuery loads the stored scheduled query, applies fn, and
// persists the result as one atomic read-modify-write. The put advances
// LastUpdatedTime, so this is the path for user-driven mutations.
func (s *Store) MutateScheduledQuery(id string, fn func(*ScheduledQuery) error) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(id)
	if err != nil {
		return err
	}
	if err := fn(sq); err != nil {
		return err
	}
	return s.PutScheduledQuery(sq)
}

// TouchScheduledQueryDelivery records the outcome of a scheduled query
// execution: the consumed schedule boundary (the internal
// deduplication marker, only ever advancing), the execution clock
// (surfaced as lastTriggeredTime, the timestamp the query was last
// executed), and the execution status. It never advances
// LastUpdatedTime — an execution is not an update.
func (s *Store) TouchScheduledQueryDelivery(id string, boundary, executedAt int64, status string) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(id)
	if err != nil {
		return err
	}
	if boundary > sq.LastExecutedBoundary {
		sq.LastExecutedBoundary = boundary
	}
	sq.LastTriggeredTime = executedAt
	sq.LastExecutionStatus = status
	return s.Put(s.scheduledQueryKey(sq.Id), sq)
}

func (s *Store) DeleteScheduledQuery(id string) error {
	key := s.scheduledQueryKey(id)
	// The lock closes the resurrection window: a delivery touch or a
	// mutation that read the record before the delete must not write it
	// back afterwards.
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	if !s.Exists(key) {
		return ErrResourceNotFound
	}
	return s.Delete(key)
}

func (s *Store) ListScheduledQueries(state string) ([]*ScheduledQuery, error) {
	return listJSONRecords[ScheduledQuery](s, keyPrefixScheduledQuery, "scheduled query record", func(sq *ScheduledQuery) bool {
		return state == "" || sq.State == state
	})
}

// --- Scheduled Query Execution History ---

func (s *Store) PutScheduledQueryExecution(exec *ScheduledQueryExecution) error {
	return s.Put(s.scheduledQueryExecutionKey(exec.ScheduledQueryId, exec.TriggerTime), exec)
}

// FinaliseScheduledQueryExecution writes the worker's terminal outcome
// unless the persisted record already reached a terminal status: a
// StopQuery that landed between the worker's cancellation checkpoint
// and this write keeps its CANCELLED status — the terminal-wins rule
// the ad-hoc query plane enforces, applied at the record plane.
func (s *Store) FinaliseScheduledQueryExecution(exec *ScheduledQueryExecution) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	key := s.scheduledQueryExecutionKey(exec.ScheduledQueryId, exec.TriggerTime)
	var current ScheduledQueryExecution
	if err := s.Get(key, &current); err == nil && IsTerminalScheduledExecutionStatus(current.Status) {
		return nil
	}
	return s.Put(key, exec)
}

// CancelScheduledQueryExecutionIfRunning flips the persisted RUNNING
// execution to CANCELLED as one atomic read-modify-write. cancelled
// reports whether the flip happened; an execution already past RUNNING
// keeps its status.
func (s *Store) CancelScheduledQueryExecutionIfRunning(exec *ScheduledQueryExecution) (bool, error) {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	key := s.scheduledQueryExecutionKey(exec.ScheduledQueryId, exec.TriggerTime)
	var current ScheduledQueryExecution
	if err := s.Get(key, &current); err != nil {
		return false, err
	}
	if current.Status != ScheduledExecutionStatusRunning {
		return false, nil
	}
	current.Status = ScheduledExecutionStatusCancelled
	return true, s.Put(key, &current)
}

func (s *Store) ListScheduledQueryExecutions(sqId string, startTime, endTime int64) ([]*ScheduledQueryExecution, error) {
	prefix := keyPrefixScheduledQueryExecution + escapePath(sqId) + ":"
	return listJSONRecords[ScheduledQueryExecution](s, prefix, "scheduled query execution record", func(exec *ScheduledQueryExecution) bool {
		if startTime > 0 && exec.TriggerTime < startTime {
			return false
		}
		return endTime <= 0 || exec.TriggerTime <= endTime
	})
}

// --- Insights Query State ---

func (s *Store) PutQueryRecord(rec *QueryRecord) error {
	return s.Put(s.queryRecordKey(rec.QueryId), rec)
}

// GetQueryRecord reads one query record by id; an absent record is the
// not-found sentinel (the caller surfaces it as the operation's
// ResourceNotFoundException).
func (s *Store) GetQueryRecord(queryId string) (*QueryRecord, error) {
	return getJSONRecord[QueryRecord](s, s.queryRecordKey(queryId), ErrResourceNotFound)
}

func (s *Store) DeleteQueryRecord(queryId string) error {
	return s.Delete(s.queryRecordKey(queryId))
}

func (s *Store) ListQueryRecords() ([]*QueryRecord, error) {
	return listJSONRecords[QueryRecord](s, keyPrefixQueryRecord, "query record", nil)
}

// QueryParameter is one saved-query parameter: the placeholder name a
// query string references with the {{parameterName}} syntax plus the
// default value and description presented alongside it.
type QueryParameter struct {
	Name         string `json:"name"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Description  string `json:"description,omitempty"`
}

// QueryDefinition represents a saved CloudWatch Logs Insights query definition.
type QueryDefinition struct {
	QueryDefinitionId string           `json:"queryDefinitionId"`
	Name              string           `json:"name"`
	QueryString       string           `json:"queryString"`
	LogGroupNames     []string         `json:"logGroupNames,omitempty"`
	QueryLanguage     string           `json:"queryLanguage,omitempty"`
	Parameters        []QueryParameter `json:"parameters,omitempty"`
	// LastModified is epoch seconds — the Timestamp shape's wire format,
	// matching the documented response examples ("lastModified":
	// 1549321515).
	LastModified int64 `json:"lastModified"`
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
	ScheduledQueryStatusTimeout  = "Timeout"
)

// Scheduled query State enum values.
const (
	ScheduledQueryStateEnabled  = "ENABLED"
	ScheduledQueryStateDisabled = "DISABLED"
)

// Scheduled-execution record statuses: the internal vocabulary of the
// execution-history records and the DescribeQueries Scheduled surface.
const (
	ScheduledExecutionStatusRunning   = "RUNNING"
	ScheduledExecutionStatusSuccess   = "SUCCESS"
	ScheduledExecutionStatusFailed    = "FAILED"
	ScheduledExecutionStatusCancelled = "CANCELLED"
	ScheduledExecutionStatusTimeout   = "TIMEOUT"
)

// IsTerminalScheduledExecutionStatus reports whether the execution
// record's status is final — a terminal record keeps its status against
// every later write.
func IsTerminalScheduledExecutionStatus(status string) bool {
	return status == ScheduledExecutionStatusSuccess ||
		status == ScheduledExecutionStatusFailed ||
		status == ScheduledExecutionStatusCancelled ||
		status == ScheduledExecutionStatusTimeout
}

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

// QueryResultRow persists one Insights query result row: the column order
// the result presented plus its field values. Both halves are needed to
// round-trip a row — the columns carry the presentation order that a bare
// field map loses.
type QueryResultRow struct {
	Columns []string          `json:"columns,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// QueryRecord is the persisted form of one Insights query's lifecycle
// state — everything DescribeQueries and GetQueryResults serve for the
// retention window, so a restart is not the availability boundary for
// query results. Status carries the service's QueryStatus vocabulary
// (Running/Complete/Failed/Cancelled); a non-terminal record found at
// construction has no worker behind it and is reconciled by the service.
type QueryRecord struct {
	QueryId             string   `json:"queryId"`
	Region              string   `json:"region"`
	LogGroupNames       []string `json:"logGroupNames,omitempty"`
	LogGroupIdentifiers []string `json:"logGroupIdentifiers,omitempty"`
	// ScannedGroups is the group set the query actually analysed — the
	// explicit list for identifier-selected queries, the SOURCE-resolved
	// set once execution has fixed it. ListLogGroupsForQuery serves it.
	ScannedGroups     []string         `json:"scannedGroups,omitempty"`
	StartTime         int64            `json:"startTime"`
	EndTime           int64            `json:"endTime"`
	QueryString       string           `json:"queryString"`
	QueryLanguage     string           `json:"queryLanguage"`
	Status            string           `json:"status"`
	CreatedAtUnixNano int64            `json:"createdAtUnixNano"`
	Results           []QueryResultRow `json:"results,omitempty"`
	RecordsScanned    int64            `json:"recordsScanned"`
	RecordsMatched    int64            `json:"recordsMatched"`
	BytesScanned      int64            `json:"bytesScanned"`
	// UserIdentity is the ARN of the principal that started the query —
	// DescribeQueries reports it back; empty for anonymous callers and
	// internally driven executions. QueryDurationMs is the wall-clock
	// duration the terminal execution took; zero while the query runs.
	UserIdentity    string `json:"userIdentity,omitempty"`
	QueryDurationMs int64  `json:"queryDurationMs,omitempty"`
}
