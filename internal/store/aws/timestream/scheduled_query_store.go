package timestream

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_timestream"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ScheduledQueryStore manages Timestream scheduled queries.
type ScheduledQueryStore struct {
	*common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	createMu   sync.Mutex
}

// scheduledQueryRecordWriteMu serialises read-modify-write cycles on
// scheduled-query records across ScheduledQueryStore instances. The engine
// and the service each construct their own store over the same Pebble
// keyspace, so an instance-level lock cannot prevent one record writer's
// read-modify-write from losing another's fields (e.g. the engine's
// UpdateLastRun racing UpdateNextRunTime or a state update).
var scheduledQueryRecordWriteMu sync.Mutex

// NewScheduledQueryStore creates a new scheduled query store.
func NewScheduledQueryStore(store storage.BasicStorage, accountID, region string) *ScheduledQueryStore {
	return &ScheduledQueryStore{
		BaseStore:  common.NewBaseStore(store.Bucket(scheduledQueryBucketName(region)), "timestream-scheduled-query"),
		TagStore:   common.NewTagStoreWithRegion(store, "timestream-scheduled-query", region),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

// CreateScheduledQuery creates a new scheduled query.
func (s *ScheduledQueryStore) CreateScheduledQuery(name, queryString string, scheduleConfig *ScheduleConfiguration, notificationConfig *NotificationConfiguration, roleARN, kmsKeyID string, errorReportConfig *ErrorReportConfiguration, targetConfig *TargetConfiguration, clientToken string) (*ScheduledQuery, error) {
	s.createMu.Lock()
	defer s.createMu.Unlock()

	key := name
	if s.Exists(key) {
		return nil, ErrScheduledQueryAlreadyExists
	}

	now := time.Now().UTC()
	sq := &ScheduledQuery{
		ARN:                            s.arnBuilder.Timestream().ScheduledQuery(name),
		Name:                           name,
		QueryString:                    queryString,
		ScheduleConfiguration:          scheduleConfig,
		NotificationConfiguration:      notificationConfig,
		ScheduledQueryExecutionRoleARN: roleARN,
		KmsKeyID:                       kmsKeyID,
		ErrorReportConfiguration:       errorReportConfig,
		TargetConfiguration:            targetConfig,
		ClientToken:                    clientToken,
		ScheduledQueryStatus:           ScheduledQueryStatusEnabled,
		CreationTime:                   now,
	}

	if err := s.PutProto(key, ScheduledQueryToProto(sq)); err != nil {
		return nil, err
	}

	return sq, nil
}

// GetScheduledQuery retrieves a scheduled query by name.
func (s *ScheduledQueryStore) GetScheduledQuery(name string) (*ScheduledQuery, error) {
	var pbSq pb.ScheduledQuery
	if err := s.GetProto(name, &pbSq); err != nil {
		return nil, ErrScheduledQueryNotFound
	}
	return ProtoToScheduledQuery(&pbSq), nil
}

// UpdateScheduledQuery updates the state of an existing scheduled query.
// Per the Smithy UpdateScheduledQueryRequest, only the State field may
// be modified (ENABLED or DISABLED). All other fields are immutable
// after creation.
func (s *ScheduledQueryStore) UpdateScheduledQuery(name string, state ScheduledQueryStatus) (*ScheduledQuery, error) {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(name)
	if err != nil {
		return nil, err
	}

	sq.ScheduledQueryStatus = state

	if err := s.PutProto(name, ScheduledQueryToProto(sq)); err != nil {
		return nil, err
	}

	return sq, nil
}

// DeleteScheduledQuery deletes a scheduled query by name.
func (s *ScheduledQueryStore) DeleteScheduledQuery(name string) error {
	// The record lock closes the resurrection window: without it an
	// UpdateLastRun/UpdateNextRunTime cycle that read before the delete
	// could write the record back after it.
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	if !s.Exists(name) {
		return ErrScheduledQueryNotFound
	}
	if err := s.TagStore.Delete(s.arnBuilder.Timestream().ScheduledQuery(name)); err != nil {
		logs.Warn("Failed to delete tags for scheduled query", logs.String("query", name), logs.Err(err))
	}
	return s.BaseStore.Delete(name)
}

// ListScheduledQueries returns all scheduled queries.
func (s *ScheduledQueryStore) ListScheduledQueries() ([]*ScheduledQuery, error) {
	result, err := common.ListProto[*pb.ScheduledQuery](s.BaseStore, common.ListOptions{}, func() *pb.ScheduledQuery { return &pb.ScheduledQuery{} }, nil)
	if err != nil {
		return nil, err
	}

	queries := make([]*ScheduledQuery, 0, len(result.Items))
	for _, pbSq := range result.Items {
		queries = append(queries, ProtoToScheduledQuery(pbSq))
	}
	return queries, nil
}

// UpdateNextRunTime updates the next run time for a scheduled query.
func (s *ScheduledQueryStore) UpdateNextRunTime(name string, nextRunTime time.Time) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(name)
	if err != nil {
		return err
	}
	sq.NextRunTime = nextRunTime
	return s.PutProto(name, ScheduledQueryToProto(sq))
}

// UpdateLastRun updates the last run information for a scheduled query.
func (s *ScheduledQueryStore) UpdateLastRun(name string, status string, runTime time.Time) error {
	scheduledQueryRecordWriteMu.Lock()
	defer scheduledQueryRecordWriteMu.Unlock()
	sq, err := s.GetScheduledQuery(name)
	if err != nil {
		return err
	}
	sq.LastRunStatus = status
	sq.PreviousRunTime = runTime
	return s.PutProto(name, ScheduledQueryToProto(sq))
}

// ScheduledQueryRunStore manages scheduled query runs.
type ScheduledQueryRunStore struct {
	*common.BaseStore
}

// NewScheduledQueryRunStore creates a new ScheduledQueryRunStore.
func NewScheduledQueryRunStore(store storage.BasicStorage, region string) *ScheduledQueryRunStore {
	return &ScheduledQueryRunStore{
		BaseStore: common.NewBaseStore(store.Bucket(scheduledQueryRunBucketName(region)), "timestream-scheduled-query-runs"),
	}
}

// CreateRun creates a new scheduled query run.
func (s *ScheduledQueryRunStore) CreateRun(scheduledQueryARN string, invocationTime, triggerTime time.Time, triggerType string) (*ScheduledQueryRun, error) {
	runID := fmt.Sprintf("%d-%s", triggerTime.UnixNano(), strings.TrimPrefix(scheduledQueryARN, "arn:"))
	runARN := scheduledQueryARN + "/run/" + runID

	run := &ScheduledQueryRun{
		ARN:               runARN,
		ScheduledQueryARN: scheduledQueryARN,
		InvocationTime:    invocationTime,
		TriggerTime:       triggerTime,
		RunStatus:         ScheduleRunStatusPending,
		TriggerType:       triggerType,
	}

	if err := s.PutProto(runARN, ScheduledQueryRunToProto(run)); err != nil {
		return nil, err
	}

	return run, nil
}

// GetRun retrieves a scheduled query run by its ARN.
func (s *ScheduledQueryRunStore) GetRun(arn string) (*ScheduledQueryRun, error) {
	if arn == "" {
		return nil, ErrInvalidParameter
	}
	var pbRun pb.ScheduledQueryRun
	if err := s.GetProto(arn, &pbRun); err != nil {
		return nil, ErrScheduledQueryRunNotFound
	}
	return ProtoToScheduledQueryRun(&pbRun), nil
}

// UpdateRunStatus updates the status of a scheduled query run.
func (s *ScheduledQueryRunStore) UpdateRunStatus(arn string, status ScheduleRunStatus, errStr string, stats *ExecutionStats) error {
	run, err := s.GetRun(arn)
	if err != nil {
		return err
	}
	run.RunStatus = status
	if errStr != "" {
		run.Error = errStr
	}
	if stats != nil {
		run.ExecutionStats = stats
	}
	if status == ScheduleRunStatusSucceeded || status == ScheduleRunStatusFailed || status == ScheduleRunStatusCancelled {
		run.CompletionTime = time.Now().UTC()
	}
	return s.PutProto(arn, ScheduledQueryRunToProto(run))
}

// DeleteRun deletes a scheduled query run by ARN.
func (s *ScheduledQueryRunStore) DeleteRun(arn string) error {
	return s.BaseStore.Delete(arn)
}

// ListRuns lists scheduled query runs for a given scheduled query ARN.
func (s *ScheduledQueryRunStore) ListRuns(scheduledQueryARN string) ([]*ScheduledQueryRun, error) {
	prefix := ""
	if scheduledQueryARN != "" {
		prefix = scheduledQueryARN + "/run/"
	}

	result, err := common.ListProto[*pb.ScheduledQueryRun](s.BaseStore, common.ListOptions{Prefix: prefix}, func() *pb.ScheduledQueryRun { return &pb.ScheduledQueryRun{} }, nil)
	if err != nil {
		return nil, err
	}

	runs := make([]*ScheduledQueryRun, 0, len(result.Items))
	for _, pbRun := range result.Items {
		runs = append(runs, ProtoToScheduledQueryRun(pbRun))
	}
	return runs, nil
}
