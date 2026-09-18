package cloudtrail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/storage/storage_cloudtrail"

	"google.golang.org/protobuf/proto"
)

// Per-event-data-store event storage. Every recorded event is copied into
// each event data store whose advanced event selectors it matches (the
// selector-less default is management events, the documented CloudTrail
// default), making each store an independent data boundary with its own
// retention. The copies live in a per-EDS bucket keyed "<unixnano>#<eventID>"
// — the same ordered-key scheme the event history uses — and carry NO
// secondary indexes: Lake queries scan the bucket in key order and apply
// their predicates in the engine, so the write path adds exactly one entry
// per matching store instead of multiplying the history's index set.

// edsEventBucketName is the per-EDS event bucket.
func edsEventBucketName(region, edsID string) string {
	return "cloudtrail-eds-events-" + region + "-" + edsID
}

// DefaultEDSEventSelectors is the selector set a selector-less event data
// store ingests by default: management events. The AWS CLI reference
// documents that "You do not need to specify any advanced event selectors
// to include all management events", and its create-event-data-store
// example response shows the synthesised "Default management events"
// selector.
func DefaultEDSEventSelectors() []AdvancedEventSelector {
	return []AdvancedEventSelector{{
		Name: "Default management events",
		FieldSelectors: []AdvancedFieldSelector{{
			Field:  "eventCategory",
			Equals: []string{"Management"},
		}},
	}}
}

// EDSSelectorMatches reports whether an event falls inside an event data
// store's ingestion configuration: the store must have ingestion enabled
// and not be pending deletion, and the event must match the store's
// advanced event selectors (or, with none configured, be a management
// event — the documented default).
func EDSSelectorMatches(eds *EventDataStore, e *Event) bool {
	if !eds.IngestionEnabled || eds.Status == "PENDING_DELETION" {
		return false
	}
	selectors := eds.AdvancedEventSelectors
	if len(selectors) == 0 {
		return e.EventCategory == "Management"
	}
	for _, sel := range selectors {
		if AdvancedSelectorMatches(sel, e) {
			return true
		}
	}
	return false
}

// ingestEventIntoEDSs copies an already-serialised event record into every
// event data store in edsList whose ingestion configuration matches the
// event. The EDS list is read by the caller BEFORE the event's own
// transaction opens (store configuration changes are rare; a store created
// after the read misses this event, which is honest — ingestion begins at
// creation). A non-nil txn writes the copies inside it; nil writes them
// directly, mirroring the event write path's fallback structure.
func (s *CloudTrailStore) ingestEventIntoEDSs(txn storage.Transaction, edsList []*EventDataStore, event *Event, key string, eventData []byte) error {
	for _, eds := range edsList {
		if !EDSSelectorMatches(eds, event) {
			continue
		}
		bucketName := edsEventBucketName(s.region, eds.EventDataStoreID)
		if txn != nil {
			if err := txn.Bucket(bucketName).Put([]byte(key), eventData); err != nil {
				return err
			}
			continue
		}
		if err := s.indexer.storage.Bucket(bucketName).Put([]byte(key), eventData); err != nil {
			return err
		}
	}
	return nil
}

// EDSQuery bounds a per-EDS event walk. Both time bounds are optional;
// NextToken is the opaque cursor LookupEDSEvents issued.
type EDSQuery struct {
	StartTime  *time.Time
	EndTime    *time.Time
	MaxResults int
	NextToken  string
}

// LookupEDSEvents walks one event data store's events in ascending event
// time order, applying the time bounds to the ordered keys. The returned
// token names the issuing bucket beside the last served key; an empty
// token means the walk reached the end. A token that does not decode as
// this bucket's cursor — another event data store's, the keys all share
// one shape — is rejected as invalid.
func (s *CloudTrailStore) LookupEDSEvents(edsID string, query EDSQuery) ([]*Event, string, error) {
	if query.MaxResults <= 0 {
		query.MaxResults = DefaultLookupEventsResults
	}
	bucketName := edsEventBucketName(s.region, edsID)
	bucket := s.indexer.storage.Bucket(bucketName)

	var start []byte
	if query.NextToken != "" {
		resume, err := base64.StdEncoding.DecodeString(query.NextToken)
		if err != nil {
			return nil, "", fmt.Errorf("%w: invalid event data store cursor", ErrInvalidNextToken)
		}
		issuer, lastKey, ok := strings.Cut(string(resume), "#")
		if !ok || issuer != bucketName || !strings.Contains(lastKey, "#") {
			return nil, "", fmt.Errorf("%w: invalid event data store cursor", ErrInvalidNextToken)
		}
		// Resume strictly after the last served key.
		start = append([]byte(lastKey), 0x00)
	}
	var end []byte
	if query.EndTime != nil {
		// The scan's end is exclusive and the bound must include an event
		// at exactly EndTime: one nanosecond past it covers every key of
		// that instant (keys are equal-width decimal nanos, so lexical
		// order is numeric order) and nothing after — channel-ingested
		// events carry partner-supplied times with full nanosecond
		// precision, so a millisecond pad would serve sub-millisecond-late
		// events.
		end = []byte(fmt.Sprintf("%d#", query.EndTime.Add(time.Nanosecond).UnixNano()))
	}
	if query.StartTime != nil {
		lower := []byte(fmt.Sprintf("%d#", query.StartTime.UnixNano()))
		if len(start) == 0 || string(lower) > string(start) {
			start = lower
		}
	}

	iter := bucket.ScanRange(start, end)
	defer iter.Close()

	events := make([]*Event, 0, query.MaxResults)
	lastKey := ""
	for iter.Next() {
		key := string(iter.Key())
		nanos, eventID, ok := strings.Cut(key, "#")
		if !ok {
			return nil, "", fmt.Errorf("malformed event data store event key %q", key)
		}
		ts, err := strconv.ParseInt(nanos, 10, 64)
		if err != nil {
			return nil, "", fmt.Errorf("malformed event data store event key %q: %w", key, err)
		}
		var p pb.Event
		if err := proto.Unmarshal(iter.Value(), &p); err != nil {
			return nil, "", fmt.Errorf("corrupt event data store event record %q: %w", key, err)
		}
		event := ProtoToEvent(&p)
		event.EventTime = time.Unix(0, ts).UTC()
		event.EventID = eventID
		events = append(events, event)
		lastKey = key
		if len(events) >= query.MaxResults {
			break
		}
	}
	if err := iter.Error(); err != nil {
		return nil, "", err
	}
	if lastKey == "" {
		return events, "", nil
	}
	return events, base64.StdEncoding.EncodeToString([]byte(bucketName + "#" + lastKey)), nil
}

// RecordEventFromPayload parses an AWS record-format JSON payload into an
// Event whose CloudTrailEvent carries the payload verbatim. The two callers
// that share it differ in identity and category: the channel path delivers
// the ID CloudTrail assigned and forces the ActivityAuditLog category, so it
// passes both; the import path keeps the record's own eventID and category,
// so it passes empty overrides — the record's eventCategory wins and a
// record without one is a management event (the record-format vocabulary).
// A payload that does not parse as JSON fails; an eventTime that is absent
// or unparseable falls back to the ingestion time, since the documented
// per-event integrity gate is the checksum, not the timestamp.
func RecordEventFromPayload(payload, overrideID, categoryOverride string) (*Event, error) {
	var record cloudTrailRecord
	if err := json.Unmarshal([]byte(payload), &record); err != nil {
		return nil, err
	}
	eventTime := time.Now().UTC()
	if record.EventTime != "" {
		if parsed, err := time.Parse(time.RFC3339, record.EventTime); err == nil {
			eventTime = parsed.UTC()
		}
	}
	readOnly := "false"
	if record.ReadOnly {
		readOnly = "true"
	}
	eventID := overrideID
	if eventID == "" {
		eventID = record.EventID
	}
	category := categoryOverride
	if category == "" {
		category = record.EventCategory
		if category == "" {
			category = "Management"
		}
	}
	event := &Event{
		EventID:         eventID,
		EventName:       record.EventName,
		EventSource:     record.EventSource,
		EventTime:       eventTime,
		EventType:       record.EventType,
		EventVersion:    record.EventVersion,
		AwsRegion:       record.AWSRegion,
		SourceIPAddress: record.SourceIPAddress,
		UserAgent:       record.UserAgent,
		ReadOnly:        readOnly,
		UserIdentity:    record.UserIdentity,
		CloudTrailEvent: payload,
		EventCategory:   category,
	}
	for _, r := range record.Resources {
		event.Resources = append(event.Resources, Resource{ResourceName: r.ARN, ResourceType: r.Type})
	}
	return event, nil
}

// PutEventIntoEDS writes an event directly into one event data store's
// bucket, bypassing the ingestion selectors — the channel ingestion path:
// events delivered through a channel are copied to the channel's
// destination store as delivered (their eventCategory is ActivityAuditLog,
// which no management selector matches). The destination record must still
// exist: a write into a store whose record is gone (the restore window's
// hard delete already dropped it) would repopulate a bucket no record
// references and no sweeper revisits — the refusal surfaces as the
// caller's per-entry failure.
func (s *CloudTrailStore) PutEventIntoEDS(edsID string, event *Event) error {
	if _, err := s.GetEventDataStore(edsID); err != nil {
		return err
	}
	if event.EventID == "" {
		event.EventID = generateEventID()
	}
	if event.EventTime.IsZero() {
		event.EventTime = time.Now().UTC()
	}
	key := fmt.Sprintf("%d#%s", event.EventTime.UnixNano(), event.EventID)
	eventData, err := proto.Marshal(EventToProto(event))
	if err != nil {
		return err
	}
	return s.indexer.storage.Bucket(edsEventBucketName(s.region, edsID)).Put([]byte(key), eventData)
}

// PutEventIntoEDSs writes an event into every named event data store's
// bucket as one atomic fan-out, bypassing the ingestion selectors — the
// channel ingestion path: a failure discovered at a later destination must
// not leave the event copied into an earlier one only, whose retry carries
// a fresh event record and would duplicate the surviving copies. Every
// copy joins the event write path's single transaction; each destination
// record must still exist at its write — a missing record (the restore
// window's hard delete) aborts the whole batch.
func (s *CloudTrailStore) PutEventIntoEDSs(edsIDs []string, event *Event) error {
	if event.EventID == "" {
		event.EventID = generateEventID()
	}
	if event.EventTime.IsZero() {
		event.EventTime = time.Now().UTC()
	}
	key := fmt.Sprintf("%d#%s", event.EventTime.UnixNano(), event.EventID)
	eventData, err := proto.Marshal(EventToProto(event))
	if err != nil {
		return err
	}
	return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
		for _, edsID := range edsIDs {
			if _, err := s.GetEventDataStore(edsID); err != nil {
				return err
			}
			if err := txn.Bucket(edsEventBucketName(s.region, edsID)).Put([]byte(key), eventData); err != nil {
				return err
			}
		}
		return nil
	})
}

// PurgeEDSEventsBefore deletes every event of one event data store whose
// event time is before cutoff, in batched transactions, mirroring the
// event-history purge. The returned count is the number of events removed;
// the operation is idempotent.
func (s *CloudTrailStore) PurgeEDSEventsBefore(edsID string, cutoff time.Time, batchSize int) (int, error) {
	if batchSize <= 0 {
		batchSize = defaultPurgeBatchSize
	}
	bucketName := edsEventBucketName(s.region, edsID)
	bound := []byte(fmt.Sprintf("%d#", cutoff.UnixNano()))
	removed := 0
	var resumeAfter string
	for {
		var keys []string
		var lastKey string
		var start []byte
		if resumeAfter != "" {
			start = append([]byte(resumeAfter), 0x00)
		}
		iter := s.indexer.storage.Bucket(bucketName).ScanRange(start, bound)
		for iter.Next() {
			keys = append(keys, string(iter.Key()))
			lastKey = string(iter.Key())
			if len(keys) >= batchSize {
				break
			}
		}
		if err := iter.Error(); err != nil {
			iter.Close()
			return removed, err
		}
		iter.Close()
		if len(keys) == 0 {
			return removed, nil
		}
		if err := s.deleteEDSEventKeys(bucketName, keys); err != nil {
			return removed, err
		}
		removed += len(keys)
		resumeAfter = lastKey
	}
}

// deleteEDSEventKeys removes one batch of per-EDS event keys in a single
// transaction when the backing storage provides them.
func (s *CloudTrailStore) deleteEDSEventKeys(bucketName string, keys []string) error {
	if s.storage != nil {
		return s.storage.Update(context.Background(), func(txn storage.Transaction) error {
			for _, k := range keys {
				if err := txn.Bucket(bucketName).Delete([]byte(k)); err != nil {
					return err
				}
			}
			return nil
		})
	}
	for _, k := range keys {
		if err := s.indexer.storage.Bucket(bucketName).Delete([]byte(k)); err != nil {
			return err
		}
	}
	return nil
}

// DropEDSEvents removes every event of one event data store — the hard
// delete path after the restore window closes.
func (s *CloudTrailStore) DropEDSEvents(edsID string) error {
	bucketName := edsEventBucketName(s.region, edsID)
	bucket := s.indexer.storage.Bucket(bucketName)
	var keys []string
	iter := bucket.ScanPrefix(nil)
	for iter.Next() {
		keys = append(keys, string(iter.Key()))
	}
	if err := iter.Error(); err != nil {
		iter.Close()
		return err
	}
	iter.Close()
	return s.deleteEDSEventKeys(bucketName, keys)
}
