package cloudtrail

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// awsErrorCode reads the wire error code a Core refusal carries.
func awsErrorCode(err error) string {
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return awsErr.GetCode()
	}
	return ""
}

// auditPayload builds a record-format channel payload and its checksum.
func auditPayload(id, eventName string) (string, string) {
	payload, err := json.Marshal(map[string]interface{}{
		"eventVersion": "1.08",
		"eventTime":    time.Now().UTC().Format(time.RFC3339),
		"eventSource":  "example.partner.com",
		"eventName":    eventName,
		"requestID":    id,
	})
	if err != nil {
		panic(err)
	}
	digest := sha256.Sum256(payload)
	return string(payload), base64.StdEncoding.EncodeToString(digest[:])
}

// newChannelFixture creates an event data store and a channel targeting it,
// returning the channel input value (the ARN).
func newChannelFixture(t *testing.T, store *cloudtrailstore.CloudTrailStore, svc *CloudTrailService, name string) (*cloudtrailstore.EventDataStore, *cloudtrailstore.Channel) {
	t.Helper()
	eds, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore(name+"-eds", store.GetAccountID(), store.GetRegion()))
	require.NoError(t, err)
	created, err := svc.createChannelCore(store, CreateChannelInput{
		Name:   name,
		Source: "Custom",
		DestinationsRaw: []interface{}{map[string]interface{}{
			"Type":     cloudtrailstore.DestinationTypeEventDataStore,
			"Location": eds.EventDataStoreARN,
		}},
		DestinationsSet: true,
	})
	require.NoError(t, err)
	ch, err := store.GetChannel(created["ChannelArn"].(string))
	require.NoError(t, err)
	return eds, ch
}

// TestPutAuditEventsCoreIngestion is the channel end-to-end unit: entries
// that pass validation land in the channel's destination event data store
// with the ActivityAuditLog category and the payload verbatim, and the
// response separates survivors from failures with the model's member
// spellings.
func TestPutAuditEventsCoreIngestion(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)
	eds, ch := newChannelFixture(t, store, svc, "ingest")

	goodPayload, goodChecksum := auditPayload("r-1", "SignIn")
	badChecksumPayload, _ := auditPayload("r-2", "SignIn")
	resp, err := svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: ch.ChannelARN,
		AuditEvents: []AuditEventInput{
			{ID: "evt-1", EventData: goodPayload, EventDataChecksum: goodChecksum, ChecksumSet: true},
			{ID: "evt-2", EventData: badChecksumPayload, EventDataChecksum: "bm90LXRoZS1jaGVja3N1bQ==", ChecksumSet: true},
			{ID: "evt-3", EventData: "not json"},
			{ID: "", EventData: goodPayload},
		},
	})
	require.NoError(t, err)

	successful := resp["successful"].([]map[string]interface{})
	failed := resp["failed"].([]map[string]interface{})
	require.Len(t, successful, 1)
	require.Len(t, failed, 3)
	assert.Equal(t, "evt-1", successful[0]["id"])
	assert.NotEmpty(t, successful[0]["eventID"], "CloudTrail assigns the eventID")

	byID := map[string]map[string]interface{}{}
	for _, f := range failed {
		byID[f["id"].(string)] = f
	}
	assert.Equal(t, "InvalidChecksum", byID["evt-2"]["errorCode"])
	assert.Equal(t, "InvalidData", byID["evt-3"]["errorCode"])
	assert.Equal(t, "FieldNotFound", byID[""]["errorCode"])

	events, _, err := store.LookupEDSEvents(eds.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, successful[0]["eventID"], events[0].EventID)
	assert.Equal(t, "ActivityAuditLog", events[0].EventCategory)
	assert.Equal(t, goodPayload, events[0].CloudTrailEvent)
}

// TestPutAuditEventsWireNullMembers pins the wire-level null reading: an
// externalId of JSON null is the absent member (the request proceeds),
// the same reading auditEvents null already takes, while the explicit
// empty string stays a present member that fails the shape bound.
func TestPutAuditEventsWireNullMembers(t *testing.T) {
	tmpDir := "./tmp/cloudtrail-put-audit-wire-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: tmpDir})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { mgr.Close() })
	st, err := mgr.GetStorage("us-east-1")
	if err != nil {
		t.Fatalf("region storage: %v", err)
	}
	store := cloudtrailstore.NewCloudTrailStore(st, "acc123", "us-east-1")
	svc := newTrailTestService(store)
	_, ch := newChannelFixture(t, store, svc, "null-members")
	reqCtx := request.NewRequestContext(context.Background(), mgr, "acc123", "us-east-1")

	payload, _ := auditPayload("r-null", "SignIn")
	req := &request.ParsedRequest{Parameters: map[string]interface{}{
		"channelArn":  ch.ChannelARN,
		"externalId":  nil,
		"auditEvents": []interface{}{map[string]interface{}{"id": "evt-null-ext", "eventData": payload}},
	}}
	resp, err := svc.PutAuditEvents(context.Background(), reqCtx, req)
	require.NoError(t, err)
	successful := resp.(map[string]interface{})["successful"].([]map[string]interface{})
	require.Len(t, successful, 1, "a null externalId is the absent member and must not fail the shape bound")

	req = &request.ParsedRequest{Parameters: map[string]interface{}{
		"channelArn":  ch.ChannelARN,
		"externalId":  "",
		"auditEvents": []interface{}{map[string]interface{}{"id": "evt-empty-ext", "eventData": payload}},
	}}
	_, err = svc.PutAuditEvents(context.Background(), reqCtx, req)
	assert.Equal(t, "ChannelUnsupportedSchema", awsErrorCode(err), "an explicit empty externalId stays a present member")
}

// TestPutAuditEventsCoreRequestErrors pins the request-level refusals: an
// absent channelArn, an ARN that is not a CloudTrail channel ARN, an
// unknown channel (by ARN and by bare ID), duplicate event IDs, an empty
// auditEvents list, and an externalId outside the model's shape.
func TestPutAuditEventsCoreRequestErrors(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)
	_, ch := newChannelFixture(t, store, svc, "errors")

	_, err := svc.putAuditEventsCore(store, PutAuditEventsInput{})
	assert.Equal(t, "InvalidChannelARN", awsErrorCode(err))

	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: "arn:aws:s3:::not-a-channel",
	})
	assert.Equal(t, "InvalidChannelARN", awsErrorCode(err))

	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: "arn:aws:cloudtrail:us-east-1:acc123:channel/00000000-0000-0000-0000-000000000000",
	})
	assert.Equal(t, "ChannelNotFound", awsErrorCode(err))

	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: "unknown-suffix",
	})
	assert.Equal(t, "ChannelNotFound", awsErrorCode(err))

	payload, _ := auditPayload("r-dup", "SignIn")
	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: ch.ChannelARN,
		AuditEvents: []AuditEventInput{
			{ID: "dup", EventData: payload},
			{ID: "dup", EventData: payload},
		},
	})
	assert.Equal(t, "DuplicatedAuditEventId", awsErrorCode(err))

	// The AuditEvents list is @length 1-100; an empty list never reaches
	// the entry loop.
	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{ChannelARN: ch.ChannelARN})
	assert.Equal(t, "ChannelUnsupportedSchema", awsErrorCode(err))

	// A present externalId must satisfy the model's ExternalId shape:
	// length 2-1224 and the printable pattern.
	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: ch.ChannelARN, ExternalID: "x", ExternalIDSet: true,
		AuditEvents: []AuditEventInput{{ID: "evt-x", EventData: payload}},
	})
	assert.Equal(t, "ChannelUnsupportedSchema", awsErrorCode(err))

	_, err = svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: ch.ChannelARN, ExternalID: "has spaces!", ExternalIDSet: true,
		AuditEvents: []AuditEventInput{{ID: "evt-x", EventData: payload}},
	})
	assert.Equal(t, "ChannelUnsupportedSchema", awsErrorCode(err))

	// Two empty ids are per-entry FieldNotFound failures, not a duplicate
	// identity: the request answers 200 with both entries failed.
	resp, err := svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN:  ch.ChannelARN,
		AuditEvents: []AuditEventInput{{EventData: payload}, {EventData: payload}},
	})
	require.NoError(t, err)
	assert.Empty(t, resp["successful"].([]map[string]interface{}))
	for _, f := range resp["failed"].([]map[string]interface{}) {
		assert.Equal(t, "FieldNotFound", f["errorCode"])
	}
}

// TestPutAuditEventsCoreBareIDSuffix pins the documented ARN-or-ID member:
// the bare ARN suffix resolves the same channel the full ARN does.
func TestPutAuditEventsCoreBareIDSuffix(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)
	_, ch := newChannelFixture(t, store, svc, "suffix")

	suffix := ch.ChannelARN[strings.LastIndex(ch.ChannelARN, "/")+1:]

	payload, checksum := auditPayload("r-suffix", "SignIn")
	resp, err := svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: suffix,
		AuditEvents: []AuditEventInput{
			{ID: "evt-suffix", EventData: payload, EventDataChecksum: checksum, ChecksumSet: true},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp["successful"].([]map[string]interface{}), 1)
}

// TestPutAuditEventsCoreStoppedDestination pins the ingestion gate: a
// destination whose ingestion is stopped accepts nothing — its entries fail
// with InvalidRecipient while the events themselves stay valid.
func TestPutAuditEventsCoreStoppedDestination(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)
	eds, ch := newChannelFixture(t, store, svc, "stopped")

	_, err := store.MutateEventDataStore(eds.EventDataStoreID, func(e *cloudtrailstore.EventDataStore) error {
		e.IngestionEnabled = false
		return nil
	})
	require.NoError(t, err)

	payload, _ := auditPayload("r-stopped", "SignIn")
	resp, err := svc.putAuditEventsCore(store, PutAuditEventsInput{
		ChannelARN: ch.ChannelARN,
		AuditEvents: []AuditEventInput{
			{ID: "evt-stopped", EventData: payload},
		},
	})
	require.NoError(t, err)
	failed := resp["failed"].([]map[string]interface{})
	require.Len(t, failed, 1)
	assert.Equal(t, "InvalidRecipient", failed[0]["errorCode"])
}

// TestChannelDestinationValidation pins the CreateChannel/UpdateChannel
// destination resolution: a destination that is not an event data store ARN
// answers EventDataStoreARNInvalidException, one that names no store this
// account owns answers EventDataStoreNotFoundException.
func TestChannelDestinationValidation(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)

	_, err := svc.createChannelCore(store, CreateChannelInput{
		Name:   "bad-location",
		Source: "Custom",
		DestinationsRaw: []interface{}{map[string]interface{}{
			"Type":     cloudtrailstore.DestinationTypeEventDataStore,
			"Location": "not-an-arn",
		}},
		DestinationsSet: true,
	})
	assert.Equal(t, "EventDataStoreARNInvalidException", awsErrorCode(err))

	_, err = svc.createChannelCore(store, CreateChannelInput{
		Name:   "missing-eds",
		Source: "Custom",
		DestinationsRaw: []interface{}{map[string]interface{}{
			"Type":     cloudtrailstore.DestinationTypeEventDataStore,
			"Location": "arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/00000000-0000-0000-0000-000000000000",
		}},
		DestinationsSet: true,
	})
	assert.Equal(t, "EventDataStoreNotFoundException", awsErrorCode(err))

	// Update applies the same resolution.
	eds, ch := newChannelFixture(t, store, svc, "update-dest")
	_, err = svc.updateChannelCore(store, UpdateChannelInput{
		Channel: ch.ChannelARN,
		DestinationsRaw: []interface{}{map[string]interface{}{
			"Type":     cloudtrailstore.DestinationTypeEventDataStore,
			"Location": eds.EventDataStoreARN,
		}},
		DestinationsSet: true,
	})
	require.NoError(t, err)
}

// TestEnforceEDSRetentionForStore pins the hourly sweep: live stores keep
// only their retention window of events, a deleted store inside the
// restore window keeps everything, and one past the window is removed with
// its events.
func TestEnforceEDSRetentionForStore(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)

	// Shortest retention the platform accepts: seven days.
	eds := cloudtrailstore.NewEventDataStore("retention", store.GetAccountID(), store.GetRegion())
	eds.RetentionPeriod = cloudtrailstore.MinEventDataStoreRetentionDays
	_, err := store.CreateEventDataStore(eds)
	require.NoError(t, err)

	recent := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Millisecond)
	ancient := time.Now().UTC().Add(-8 * 24 * time.Hour).Truncate(time.Millisecond)
	for _, e := range []*cloudtrailstore.Event{
		{EventID: "keep", EventName: "CreateTrail", EventTime: recent, EventCategory: "Management"},
		{EventID: "expire", EventName: "DeleteTrail", EventTime: ancient, EventCategory: "Management"},
	} {
		require.NoError(t, store.PutEventIntoEDS(eds.EventDataStoreID, e))
	}

	inside := cloudtrailstore.NewEventDataStore("inside-window", store.GetAccountID(), store.GetRegion())
	_, err = store.CreateEventDataStore(inside)
	require.NoError(t, err)
	require.NoError(t, store.PutEventIntoEDS(inside.EventDataStoreID,
		&cloudtrailstore.Event{EventID: "ancient-too", EventName: "CreateTrail", EventTime: ancient, EventCategory: "Management"}))
	_, err = store.MutateEventDataStore(inside.EventDataStoreID, func(e *cloudtrailstore.EventDataStore) error {
		e.Status = "PENDING_DELETION"
		// The fixtures derive from the effective restore window: TEST_MODE
		// shortens it to half a second, and a fixed one-hour "inside"
		// fixture would land outside that window.
		deleted := time.Now().UTC().Add(-edsRestoreWindow() / 2)
		e.DeletedTimestamp = &deleted
		return nil
	})
	require.NoError(t, err)

	outside := cloudtrailstore.NewEventDataStore("outside-window", store.GetAccountID(), store.GetRegion())
	_, err = store.CreateEventDataStore(outside)
	require.NoError(t, err)
	require.NoError(t, store.PutEventIntoEDS(outside.EventDataStoreID,
		&cloudtrailstore.Event{EventID: "dropped", EventName: "CreateTrail", EventTime: ancient, EventCategory: "Management"}))
	_, err = store.MutateEventDataStore(outside.EventDataStoreID, func(e *cloudtrailstore.EventDataStore) error {
		e.Status = "PENDING_DELETION"
		deleted := time.Now().UTC().Add(-(edsRestoreWindow() + time.Hour))
		e.DeletedTimestamp = &deleted
		return nil
	})
	require.NoError(t, err)

	svc.enforceEDSRetentionForStore(store)

	events, _, err := store.LookupEDSEvents(eds.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "keep", events[0].EventID)

	insideEvents, _, err := store.LookupEDSEvents(inside.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	assert.Len(t, insideEvents, 1, "restore window keeps the deleted store's events")

	_, err = store.GetEventDataStore(outside.EventDataStoreID)
	assert.ErrorIs(t, err, cloudtrailstore.ErrEventDataStoreNotFound,
		"the past-window store's record is gone")
	outsideEvents, _, err := store.LookupEDSEvents(outside.EventDataStoreID, cloudtrailstore.EDSQuery{})
	require.NoError(t, err)
	assert.Empty(t, outsideEvents, "the past-window store's events are dropped")
}

// TestExecuteQueryIsEDSScoped pins the data boundary at the query engine:
// a query reads only its own store's ingested copies; an event recorded in
// the history and another store's bucket never leaks into the result.
func TestExecuteQueryIsEDSScoped(t *testing.T) {
	store := newQueryTestStore(t)
	svc := newTrailTestService(store)

	// The history event lands before any store exists, so no per-store copy
	// of it can exist.
	seedLakeEvent(t, store, "alice", "HistoryOnly", false, nil)

	mine := cloudtrailstore.NewEventDataStore("mine", store.GetAccountID(), store.GetRegion())
	theirs := cloudtrailstore.NewEventDataStore("theirs", store.GetAccountID(), store.GetRegion())
	for _, eds := range []*cloudtrailstore.EventDataStore{mine, theirs} {
		_, err := store.CreateEventDataStore(eds)
		require.NoError(t, err)
	}

	// The other store's direct write must not surface in mine's query.
	theirsEvent, err := cloudtrailstore.RecordEventFromPayload(
		`{"eventName":"TheirsOnly"}`, "theirs-1", "ActivityAuditLog")
	require.NoError(t, err)
	require.NoError(t, store.PutEventIntoEDS(theirs.EventDataStoreID, theirsEvent))
	mineEvent, err := cloudtrailstore.RecordEventFromPayload(
		`{"eventName":"MineOnly"}`, "mine-1", "ActivityAuditLog")
	require.NoError(t, err)
	require.NoError(t, store.PutEventIntoEDS(mine.EventDataStoreID, mineEvent))

	pq, err := parseQueryStatement("SELECT eventName FROM " + mine.EventDataStoreARN)
	require.NoError(t, err)
	rows, stats, timedOut, err := svc.executeQuery(store, pq, time.Time{})
	require.NoError(t, err)
	assert.False(t, timedOut)
	require.Len(t, rows, 1)
	assert.Equal(t, "MineOnly", rows[0][0]["eventName"])
	assert.Equal(t, int64(1), stats.eventsScanned, "one store, one event scanned")
}
