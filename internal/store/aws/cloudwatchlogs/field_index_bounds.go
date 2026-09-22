package cloudwatchlogs

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
)

// The field-index bounds record family: the per-group matching-event time
// bounds the ingestion scan maintains. CloudWatch Logs builds field
// indexes when log events are ingested — FieldIndex.lastScanTime is "The
// most recent time that CloudWatch Logs scanned ingested log events to
// search for this field index" — so the bounds accumulate at ingestion
// and the DescribeFieldIndexes listing reads this record instead of
// walking the group's events. A field's bounds count only events
// ingested while a policy watches it: firstEventTime is "the earliest
// log event that matches this field index, after the index policy that
// contains it was created", so pre-creation events never enter the
// record. Bounds freeze when the policy drops the field (the INACTIVE
// category reports the frozen range) and resume if the field returns.

// IndexIdentityFields names the DEFAULT and generated fields every
// ingested event carries by construction (its stream, its ingestion
// stamp, the account and region of its group): each such field matches
// every event in a scan, so its DEFAULT bounds are the group's event
// bounds. The remaining @-prefixed members of the documented default
// listing describe vended-logs record metadata the platform's log plane
// never stamps onto events, so they match no event.
var IndexIdentityFields = map[string]bool{
	"@logStream": true, "@ingestionTime": true,
	"@aws.region": true, "@aws.account": true,
}

// FieldScanMessage is one ingested event as the field-index scan sees
// it: the event timestamp and the message form the query plane serves
// (the transformed copy where one exists).
type FieldScanMessage struct {
	Timestamp int64
	Message   string
}

// FieldIndexEventRange is one field's matching-event time bounds.
type FieldIndexEventRange struct {
	First int64 `json:"first"`
	Last  int64 `json:"last"`
}

// FieldIndexBounds is one log group's accumulated field-index state: the
// earliest and latest matching-event timestamp per scanned field and the
// most recent ingestion scan's stamp.
type FieldIndexBounds struct {
	LogGroupName string                          `json:"logGroupName"`
	Fields       map[string]FieldIndexEventRange `json:"fields,omitempty"`
	LastScanTs   int64                           `json:"lastScanTs"`
}

// MergeFieldIndexBounds folds one ingested batch's scan into the group's
// record: every watched field a message carries (a JSON path lookup)
// extends its bounds, and the scan stamp advances whether or not
// anything matched. The merge is best-effort — the events themselves are
// already durable when this runs, so a bounds failure never fails the
// ingestion that produced it.
func (s *Store) MergeFieldIndexBounds(logGroupName string, fields []string, messages []FieldScanMessage) {
	if len(messages) == 0 {
		return
	}
	parsed := make([]map[string]interface{}, len(messages))
	for i, m := range messages {
		var data map[string]interface{}
		if json.Unmarshal([]byte(m.Message), &data) == nil {
			parsed[i] = data
		}
	}
	matched := make(map[string]FieldIndexEventRange, len(fields))
	for _, name := range fields {
		for i, m := range messages {
			if !IndexIdentityFields[name] && (parsed[i] == nil || !jsonPathPresent(parsed[i], name)) {
				continue
			}
			r, ok := matched[name]
			if !ok {
				r = FieldIndexEventRange{First: m.Timestamp, Last: m.Timestamp}
			} else {
				if m.Timestamp < r.First {
					r.First = m.Timestamp
				}
				if m.Timestamp > r.Last {
					r.Last = m.Timestamp
				}
			}
			matched[name] = r
		}
	}

	// The record's read-modify-write and the group-existence check run
	// under the group's write lock — the same lock the whole teardown
	// holds — so a merge in flight while the group is deleted either
	// completes before the teardown (and its record dies with the group)
	// or observes the removed group and writes nothing.
	gLock := s.groupLock(logGroupName)
	gLock.Lock()
	defer gLock.Unlock()
	// A group deleted while the batch was in flight must not resurrect
	// records.
	if _, err := s.GetLogGroup(logGroupName); err != nil {
		return
	}
	var rec FieldIndexBounds
	if err := s.Get(s.fieldIndexBoundsKey(logGroupName), &rec); err != nil {
		rec = FieldIndexBounds{LogGroupName: logGroupName}
	}
	if rec.Fields == nil {
		rec.Fields = make(map[string]FieldIndexEventRange, len(matched))
	}
	for name, r := range matched {
		cur, ok := rec.Fields[name]
		if !ok {
			rec.Fields[name] = r
			continue
		}
		if r.First < cur.First {
			cur.First = r.First
		}
		if r.Last > cur.Last {
			cur.Last = r.Last
		}
		rec.Fields[name] = cur
	}
	rec.LastScanTs = time.Now().UTC().UnixMilli()
	if err := s.Put(s.fieldIndexBoundsKey(logGroupName), &rec); err != nil {
		logs.Error("Failed to persist field index bounds",
			logs.String("logGroupName", logGroupName), logs.Err(err))
	}
}

// ReadFieldIndexBounds returns the group's accumulated bounds; a group
// whose events predate the ingestion scan (or that never ingested) has
// no record and reports the empty bounds.
func (s *Store) ReadFieldIndexBounds(logGroupName string) *FieldIndexBounds {
	var rec FieldIndexBounds
	if err := s.Get(s.fieldIndexBoundsKey(logGroupName), &rec); err != nil {
		return &FieldIndexBounds{LogGroupName: logGroupName}
	}
	return &rec
}

// DeleteFieldIndexBounds removes the record (the group teardown path:
// a same-named recreation starts unscanned).
func (s *Store) DeleteFieldIndexBounds(logGroupName string) error {
	return s.Delete(s.fieldIndexBoundsKey(logGroupName))
}

// jsonPathPresent walks a dotted field path into a JSON log event,
// numeric segments indexing arrays ("requestParameters.instancesSet.
// items.0.instanceId" — "The 0 refers to that field's place in the
// array").
func jsonPathPresent(data map[string]interface{}, path string) bool {
	var current interface{} = data
	for _, segment := range strings.Split(path, ".") {
		switch node := current.(type) {
		case map[string]interface{}:
			value, ok := node[segment]
			if !ok {
				return false
			}
			current = value
		case []interface{}:
			index, err := strconv.Atoi(segment)
			if err != nil || index < 0 || index >= len(node) {
				return false
			}
			current = node[index]
		default:
			return false
		}
	}
	return true
}
