package cloudwatch

import (
	"encoding/json"
	"fmt"
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
	prefix := "alarm:"
	if alarmNamePrefix != "" {
		prefix = "alarm:" + alarmNamePrefix
	}

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

// SetAlarmState updates the state of a CloudWatch alarm.
func (s *AlarmStore) SetAlarmState(name, state, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	alarm, err := s.GetAlarm(name)
	if err != nil {
		return err
	}

	alarm.State = state
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

// AddAlarmHistory adds a history entry for an alarm.
func (s *AlarmStore) AddAlarmHistory(entry *AlarmHistoryEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := fmt.Sprintf("history:%s:%d", entry.AlarmName, entry.Timestamp)
	return s.Put(key, entry)
}

// ListAlarmHistory retrieves alarm history entries, optionally filtered by alarm name and history item type.
func (s *AlarmStore) ListAlarmHistory(alarmName string, historyItemType string) ([]*AlarmHistoryEntry, error) {
	prefix := "history:"
	if alarmName != "" {
		prefix = "history:" + alarmName + ":"
	}

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
