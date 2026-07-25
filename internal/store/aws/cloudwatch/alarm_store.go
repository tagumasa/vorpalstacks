package cloudwatch

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func alarmBucketName(region string) string {
	return "cw_alarms-" + region
}

// AlarmStore provides CloudWatch alarm storage operations.
type AlarmStore struct {
	*common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.Mutex
}

// NewAlarmStore creates a new CloudWatch alarm store.
func NewAlarmStore(store storage.BasicStorage, accountID, region string) *AlarmStore {
	return &AlarmStore{
		BaseStore:  common.NewBaseStore(store.Bucket(alarmBucketName(region)), "cloudwatch"),
		TagStore:   common.NewTagStoreWithRegion(store, "cloudwatch", region),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

func (s *AlarmStore) buildAlarmArn(name string) string {
	return s.arnBuilder.Build("cloudwatch", "alarm:"+name)
}

func (s *AlarmStore) buildAlarmKey(name string) string {
	return "alarm:" + name
}

// CreateAlarm creates a new CloudWatch alarm.
func (s *AlarmStore) CreateAlarm(alarm *Alarm) (*Alarm, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if alarm.Name == "" {
		return nil, ErrInvalidParameter
	}

	key := s.buildAlarmKey(alarm.Name)
	if s.Exists(key) {
		return nil, ErrAlarmAlreadyExists
	}

	alarm.ARN = s.buildAlarmArn(alarm.Name)
	if alarm.CreatedAt.IsZero() {
		alarm.CreatedAt = time.Now().UTC()
	}
	if alarm.StateUpdatedTimestamp.IsZero() {
		alarm.StateUpdatedTimestamp = alarm.CreatedAt
	}
	if alarm.State == "" {
		alarm.State = "INSUFFICIENT_DATA"
	}
	if alarm.Tags == nil {
		alarm.Tags = make(map[string]string)
	}

	if err := s.Put(key, alarm); err != nil {
		return nil, err
	}

	return alarm, nil
}

// GetAlarm retrieves a CloudWatch alarm by name.
func (s *AlarmStore) GetAlarm(name string) (*Alarm, error) {
	if name == "" {
		return nil, ErrInvalidParameter
	}

	key := s.buildAlarmKey(name)
	var alarm Alarm
	if err := s.BaseStore.Get(key, &alarm); err != nil {
		return nil, ErrAlarmNotFound
	}
	return &alarm, nil
}

// UpdateAlarm modifies an existing CloudWatch alarm.
func (s *AlarmStore) UpdateAlarm(alarm *Alarm) error {
	if alarm.Name == "" {
		return ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildAlarmKey(alarm.Name)
	if !s.Exists(key) {
		return ErrAlarmNotFound
	}

	return s.Put(key, alarm)
}

// DeleteAlarm removes a CloudWatch alarm by name.
func (s *AlarmStore) DeleteAlarm(name string) error {
	if name == "" {
		return ErrInvalidParameter
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.buildAlarmKey(name)
	if !s.Exists(key) {
		return ErrAlarmNotFound
	}

	arn := s.buildAlarmArn(name)
	if err := s.TagStore.Delete(arn); err != nil {
		return err
	}

	return s.BaseStore.Delete(key)
}

// ListAlarms returns a list of CloudWatch alarms, optionally filtered by name prefix.
func (s *AlarmStore) ListAlarms(alarmNamePrefix string) ([]*Alarm, error) {
	prefix := s.buildAlarmKey(alarmNamePrefix)

	var alarms []*Alarm
	err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var alarm Alarm
		if err := json.Unmarshal(value, &alarm); err != nil {
			return err
		}
		alarms = append(alarms, &alarm)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return alarms, nil
}

// ListAlarmsPaginated returns a paginated list of CloudWatch alarms with optional filter.
func (s *AlarmStore) ListAlarmsPaginated(alarmNamePrefix string, opts common.ListOptions, filter func(*Alarm) bool) (*common.ListResult[Alarm], error) {
	opts.Prefix = s.buildAlarmKey(alarmNamePrefix)
	return common.List[Alarm](s.BaseStore, opts, filter)
}

// SetAlarmState updates the state of a CloudWatch alarm.
func (s *AlarmStore) SetAlarmState(name, state, reason, reasonData string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	alarm, err := s.GetAlarm(name)
	if err != nil {
		return err
	}

	alarm.State = state
	alarm.StateReason = reason
	alarm.StateReasonData = reasonData
	alarm.StateUpdatedTimestamp = time.Now().UTC()

	key := s.buildAlarmKey(alarm.Name)
	return s.Put(key, alarm)
}

// SetAlarmActionsEnabled enables or disables actions for a CloudWatch alarm.
func (s *AlarmStore) SetAlarmActionsEnabled(name string, enabled bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	alarm, err := s.GetAlarm(name)
	if err != nil {
		return err
	}

	alarm.ActionsEnabled = enabled

	key := s.buildAlarmKey(name)
	return s.Put(key, alarm)
}

// SetAlarmActionsSuppressed records the suppression reason on the alarm
// when a suppressor alarm causes actions to be withheld.
func (s *AlarmStore) SetAlarmActionsSuppressed(name string, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	alarm, err := s.GetAlarm(name)
	if err != nil {
		return err
	}

	alarm.ActionsSuppressedBy = reason
	alarm.ActionsSuppressedReason = fmt.Sprintf("Actions suppressed by %s", reason)

	key := s.buildAlarmKey(name)
	return s.Put(key, alarm)
}

// AddAlarmHistory adds a history entry for an alarm.
func (s *AlarmStore) AddAlarmHistory(entry *AlarmHistoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := alarmHistoryKey(entry.AlarmName, entry.Timestamp)
	return s.Put(key, entry)
}

// ListAlarmHistory retrieves alarm history entries, optionally filtered by alarm name and history item type.
func (s *AlarmStore) ListAlarmHistory(alarmName string, historyItemType string) ([]*AlarmHistoryEntry, error) {
	prefix := alarmHistoryPrefix(alarmName)

	var entries []*AlarmHistoryEntry
	err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var entry AlarmHistoryEntry
		if err := json.Unmarshal(value, &entry); err != nil {
			return err
		}
		if historyItemType != "" && entry.HistoryItemType != historyItemType {
			return nil
		}
		entries = append(entries, &entry)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return entries, nil
}

// ListAlarmHistoryPaginated returns a paginated list of alarm history entries.
// ListAlarmHistoryOpts holds all filter and pagination options for
// listing alarm history entries.
type ListAlarmHistoryOpts struct {
	AlarmName       string
	HistoryItemType string
	AlarmTypes      map[string]bool
	StartDate       int64 // UnixMilli; 0 = no filter
	EndDate         int64 // UnixMilli; 0 = no filter
	ScanBy          string
	ListOpts        common.ListOptions
}

func (s *AlarmStore) ListAlarmHistoryPaginated(opts ListAlarmHistoryOpts) (*common.ListResult[AlarmHistoryEntry], error) {
	prefix := alarmHistoryPrefix(opts.AlarmName)

	// Build the combined filter for HistoryItemType, AlarmTypes,
	// StartDate and EndDate — applied during iteration so that
	// pagination counts only matching entries.
	var filter func(*AlarmHistoryEntry) bool
	if opts.HistoryItemType != "" || len(opts.AlarmTypes) > 0 || opts.StartDate > 0 || opts.EndDate > 0 {
		filter = func(e *AlarmHistoryEntry) bool {
			if opts.HistoryItemType != "" && e.HistoryItemType != opts.HistoryItemType {
				return false
			}
			if len(opts.AlarmTypes) > 0 && !opts.AlarmTypes[e.AlarmType] {
				return false
			}
			if opts.StartDate > 0 && e.Timestamp < opts.StartDate {
				return false
			}
			if opts.EndDate > 0 && e.Timestamp > opts.EndDate {
				return false
			}
			return true
		}
	}

	// Fetch ALL matching items via ListMatching (unlimited iteration)
	// so that ScanBy sort is global, not per-page. Using common.List
	// with MaxItems=0 would silently clamp to 100 items.
	allItems, err := common.ListMatching[AlarmHistoryEntry](s.BaseStore, prefix, filter)
	if err != nil {
		return nil, err
	}

	// Sort: Pebble keys are ascending by timestamp (oldest first).
	// AWS default ScanBy is TimestampDescending (newest first), so
	// reverse unless explicitly TimestampAscending.
	if opts.ScanBy != "TimestampAscending" {
		for i, j := 0, len(allItems)-1; i < j; i, j = i+1, j-1 {
			allItems[i], allItems[j] = allItems[j], allItems[i]
		}
	}

	// Apply manual offset-based pagination on the sorted list.
	startIdx := 0
	maxItems := opts.ListOpts.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}
	if opts.ListOpts.Marker != "" {
		if idx, err := strconv.Atoi(opts.ListOpts.Marker); err == nil && idx >= 0 && idx < len(allItems) {
			startIdx = idx
		}
	}
	endIdx := startIdx + maxItems
	if endIdx > len(allItems) {
		endIdx = len(allItems)
	}

	pageItems := allItems[startIdx:endIdx]
	nextMarker := ""
	if endIdx < len(allItems) {
		nextMarker = strconv.Itoa(endIdx)
	}

	return &common.ListResult[AlarmHistoryEntry]{
		Items:       pageItems,
		NextMarker:  nextMarker,
		IsTruncated: endIdx < len(allItems),
	}, nil
}

func alarmHistoryKey(alarmName string, timestamp int64) string {
	return fmt.Sprintf("history:%s:%d", alarmName, timestamp)
}

func alarmHistoryPrefix(alarmName string) string {
	if alarmName == "" {
		return "history:"
	}
	return "history:" + alarmName + ":"
}
