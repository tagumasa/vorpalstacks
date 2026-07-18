// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package iot

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
)

// Alarm state names as defined by AWS IoT Events.
const (
	AlarmStateNormal       = "NORMAL"
	AlarmStateActive       = "ACTIVE"
	AlarmStateAcknowledged = "ACKNOWLEDGED"
	AlarmStateSnooze       = "SNOOZE"
	AlarmStateLatched      = "LATCHED"
	AlarmStateDisabled     = "DISABLED"
)

// alarmInstanceState represents the runtime state of a single alarm instance,
// keyed by (alarmModelName, keyValue).
type alarmInstanceState struct {
	stateName    string
	creationTime time.Time
	lastUpdate   time.Time
}

// alarmStateRecord is the JSON-serialisable form of alarmInstanceState
// used for Pebble persistence (H-SM1).
type alarmStateRecord struct {
	StateName    string `json:"stateName"`
	CreationTime int64  `json:"creationTime"`
	LastUpdate   int64  `json:"lastUpdate"`
}

// AlarmStateMachine manages runtime state for IoT Events alarm instances.
// State is persisted to Pebble via the embedded BaseStore so that alarm
// instances survive server restarts (H-SM1). If store is nil (e.g. in
// unit tests) the machine operates in-memory only.
//
// Alarms are created explicitly via EnableAlarm and queried via
// DescribeAlarm/ListAlarms. Non-existent alarms return a sentinel so that
// the service layer can respond with ResourceNotFoundException (matching
// AWS behaviour).
type AlarmStateMachine struct {
	mu         sync.RWMutex
	store      *common.BaseStore
	alarms     map[string]map[string]*alarmInstanceState
	modelDefs  map[string]bool              // tracks known alarm model names
	alarmRules map[string]*loadedAlarmModel // alarm models with parsed SimpleRules
}

// NewAlarmStateMachine creates a new AlarmStateMachine. If store is
// non-nil, existing alarm state is loaded from Pebble on construction
// so that alarm instances survive server restarts (H-SM1).
func NewAlarmStateMachine(store *common.BaseStore) *AlarmStateMachine {
	sm := &AlarmStateMachine{
		store:      store,
		alarms:     make(map[string]map[string]*alarmInstanceState),
		modelDefs:  make(map[string]bool),
		alarmRules: make(map[string]*loadedAlarmModel),
	}
	if store != nil {
		sm.loadFromStore()
	}
	return sm
}

// alarmStateKey builds the Pebble key for an alarm instance.
func alarmStateKey(modelName, keyValue string) string {
	return modelName + "#" + keyValue
}

// loadFromStore rebuilds the in-memory map from persisted state. Called
// once during construction before any API traffic, so no lock is needed.
func (sm *AlarmStateMachine) loadFromStore() {
	_ = sm.store.ScanPrefix("", func(key string, value []byte) error {
		var rec alarmStateRecord
		if err := json.Unmarshal(value, &rec); err != nil {
			logs.Warn("Failed to unmarshal persisted alarm state; skipping",
				logs.String("key", key), logs.Err(err))
			return nil
		}
		idx := strings.Index(key, "#")
		if idx < 0 {
			return nil
		}
		modelName := key[:idx]
		keyValue := key[idx+1:]
		if sm.alarms[modelName] == nil {
			sm.alarms[modelName] = make(map[string]*alarmInstanceState)
		}
		sm.alarms[modelName][keyValue] = &alarmInstanceState{
			stateName:    rec.StateName,
			creationTime: time.Unix(0, rec.CreationTime).UTC(),
			lastUpdate:   time.Unix(0, rec.LastUpdate).UTC(),
		}
		return nil
	})
}

// persist writes the alarm instance state to Pebble. Called while the
// caller holds the write lock, so no additional locking is needed.
// Persistence failures are logged as warnings; the in-memory state is
// still updated so the current request succeeds, and the next state
// transition will attempt to persist again.
func (sm *AlarmStateMachine) persist(modelName, keyValue string, alarm *alarmInstanceState) {
	if sm.store == nil {
		return
	}
	rec := alarmStateRecord{
		StateName:    alarm.stateName,
		CreationTime: alarm.creationTime.UnixNano(),
		LastUpdate:   alarm.lastUpdate.UnixNano(),
	}
	if err := sm.store.Put(alarmStateKey(modelName, keyValue), rec); err != nil {
		logs.Warn("Failed to persist alarm state to Pebble",
			logs.String("alarmModel", modelName),
			logs.String("keyValue", keyValue),
			logs.Err(err))
	}
}

// EnableAlarm creates or resets the alarm instance for the given model/key.
// A newly created alarm starts in NORMAL state. If the alarm already exists,
// its state transitions to NORMAL (matching AWS BatchEnableAlarm semantics).
func (sm *AlarmStateMachine) EnableAlarm(modelName, keyValue string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := time.Now().UTC()
	if sm.alarms[modelName] == nil {
		sm.alarms[modelName] = make(map[string]*alarmInstanceState)
	}
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		alarm = &alarmInstanceState{
			stateName:    AlarmStateNormal,
			creationTime: now,
			lastUpdate:   now,
		}
		sm.alarms[modelName][keyValue] = alarm
		sm.persist(modelName, keyValue, alarm)
		return
	}
	alarm.stateName = AlarmStateNormal
	alarm.lastUpdate = now
	sm.persist(modelName, keyValue, alarm)
}

// DisableAlarm transitions the alarm to DISABLED state. If the alarm does
// not exist, it is created in DISABLED state (AWS silently accepts
// BatchDisableAlarm for non-existent alarms).
func (sm *AlarmStateMachine) DisableAlarm(modelName, keyValue string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := time.Now().UTC()
	if sm.alarms[modelName] == nil {
		sm.alarms[modelName] = make(map[string]*alarmInstanceState)
	}
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		alarm = &alarmInstanceState{
			stateName:    AlarmStateDisabled,
			creationTime: now,
			lastUpdate:   now,
		}
		sm.alarms[modelName][keyValue] = alarm
		sm.persist(modelName, keyValue, alarm)
		return
	}
	alarm.stateName = AlarmStateDisabled
	alarm.lastUpdate = now
	sm.persist(modelName, keyValue, alarm)
}

// AcknowledgeAlarm transitions an ACTIVE or LATCHED alarm to ACKNOWLEDGED.
// Returns true if the transition was applied, false if the alarm does not
// exist or is not in a state that accepts acknowledgement.
func (sm *AlarmStateMachine) AcknowledgeAlarm(modelName, keyValue string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.alarms[modelName] == nil {
		return false
	}
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		return false
	}
	if alarm.stateName == AlarmStateActive || alarm.stateName == AlarmStateLatched {
		alarm.stateName = AlarmStateAcknowledged
		alarm.lastUpdate = time.Now().UTC()
		sm.persist(modelName, keyValue, alarm)
		return true
	}
	return false
}

// ResetAlarm transitions the alarm to NORMAL state. Returns true if the
// alarm existed and was reset.
func (sm *AlarmStateMachine) ResetAlarm(modelName, keyValue string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.alarms[modelName] == nil {
		return false
	}
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		return false
	}
	alarm.stateName = AlarmStateNormal
	alarm.lastUpdate = time.Now().UTC()
	sm.persist(modelName, keyValue, alarm)
	return true
}

// SnoozeAlarm transitions an ACTIVE alarm to SNOOZE state. Returns true
// if the alarm existed and was snoozed. AWS only allows snoozing alarms
// that are in the ACTIVE state (M-S2).
func (sm *AlarmStateMachine) SnoozeAlarm(modelName, keyValue string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.alarms[modelName] == nil {
		return false
	}
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		return false
	}
	if alarm.stateName != AlarmStateActive {
		return false
	}
	alarm.stateName = AlarmStateSnooze
	alarm.lastUpdate = time.Now().UTC()
	sm.persist(modelName, keyValue, alarm)
	return true
}

// DescribeAlarm returns the state of the named alarm instance. Returns
// (state, creationTime, lastUpdate, true) if the alarm exists, or
// ("", zero, zero, false) if it does not.
func (sm *AlarmStateMachine) DescribeAlarm(modelName, keyValue string) (stateName string, creationTime, lastUpdate time.Time, ok bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if sm.alarms[modelName] == nil {
		return "", time.Time{}, time.Time{}, false
	}
	alarm, exists := sm.alarms[modelName][keyValue]
	if !exists {
		return "", time.Time{}, time.Time{}, false
	}
	return alarm.stateName, alarm.creationTime, alarm.lastUpdate, true
}

// ListAlarms returns all alarm instances for the given alarm model name.
// Each entry is (keyValue, stateName, creationTime, lastUpdate).
func (sm *AlarmStateMachine) ListAlarms(modelName string) []AlarmInstanceSummary {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	modelAlarms, exists := sm.alarms[modelName]
	if !exists {
		return nil
	}
	result := make([]AlarmInstanceSummary, 0, len(modelAlarms))
	for key, alarm := range modelAlarms {
		result = append(result, AlarmInstanceSummary{
			KeyValue:     key,
			StateName:    alarm.stateName,
			CreationTime: alarm.creationTime,
			LastUpdate:   alarm.lastUpdate,
		})
	}
	return result
}

// AlarmInstanceSummary is a read-only snapshot of an alarm instance.
type AlarmInstanceSummary struct {
	KeyValue     string
	StateName    string
	CreationTime time.Time
	LastUpdate   time.Time
}
