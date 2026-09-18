package cloudtrail

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vorpalstacks/internal/core/storage"
)

func newEDSEventsTestStore(t *testing.T) (*CloudTrailStore, *storage.PebbleStorage) {
	t.Helper()
	tmpDir := "./tmp/cloudtrail-eds-events-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() { s.Close() })

	return NewCloudTrailStore(s, "acc123", "us-east-1"), s
}

// edsTestEvent builds an event with an explicit category: ingestion
// selectors evaluate the category, and store-level fixtures bypass NewEvent
// (which fixes it to Management).
func edsTestEvent(id, name, category string, at time.Time) *Event {
	return &Event{
		EventID:       id,
		EventName:     name,
		EventSource:   "cloudtrail.amazonaws.com",
		EventTime:     at,
		EventCategory: category,
		UserIdentity:  &UserIdentity{Type: "IAMUser", UserName: "alice"},
	}
}

// TestNewEventDataStoreMaterialisesDefaultSelectors pins the register-3
// re-adjudication: a store created without selectors carries the documented
// default — management events — as a materialised selector set.
func TestNewEventDataStoreMaterialisesDefaultSelectors(t *testing.T) {
	eds := NewEventDataStore("defaulted", "acc123", "us-east-1")
	require.Len(t, eds.AdvancedEventSelectors, 1)
	assert.Equal(t, "Default management events", eds.AdvancedEventSelectors[0].Name)
	require.Len(t, eds.AdvancedEventSelectors[0].FieldSelectors, 1)
	fs := eds.AdvancedEventSelectors[0].FieldSelectors[0]
	assert.Equal(t, "eventCategory", fs.Field)
	assert.Equal(t, []string{"Management"}, fs.Equals)
}

// TestPutEventFansOutToMatchingEDSs verifies the ingestion matrix: the
// default selector copies management events only, a custom selector copies
// its own matches regardless of category, and ingestion-stopped or
// pending-deletion stores copy nothing. The event history itself keeps
// every event.
func TestPutEventFansOutToMatchingEDSs(t *testing.T) {
	store, _ := newEDSEventsTestStore(t)

	defaultEDS := NewEventDataStore("defaulted", "acc123", "us-east-1")
	customEDS := NewEventDataStore("custom", "acc123", "us-east-1")
	customEDS.AdvancedEventSelectors = []AdvancedEventSelector{{
		Name: "CreateTrail only",
		FieldSelectors: []AdvancedFieldSelector{{
			Field:  "eventName",
			Equals: []string{"CreateTrail"},
		}},
	}}
	stoppedEDS := NewEventDataStore("stopped", "acc123", "us-east-1")
	stoppedEDS.IngestionEnabled = false
	pendingEDS := NewEventDataStore("pending", "acc123", "us-east-1")
	pendingEDS.Status = "PENDING_DELETION"
	for _, eds := range []*EventDataStore{defaultEDS, customEDS, stoppedEDS, pendingEDS} {
		_, err := store.CreateEventDataStore(eds)
		require.NoError(t, err)
	}

	at := time.Now().UTC().Truncate(time.Millisecond)
	mgmt := edsTestEvent("fan-mgmt", "CreateTrail", "Management", at)
	data := edsTestEvent("fan-data", "DeleteObject", "Data", at.Add(time.Millisecond))
	for _, e := range []*Event{mgmt, data} {
		require.NoError(t, store.PutEvent(e))
	}

	// The default store holds the management event alone.
	defEvents, _, err := store.LookupEDSEvents(defaultEDS.EventDataStoreID, EDSQuery{})
	require.NoError(t, err)
	require.Len(t, defEvents, 1)
	assert.Equal(t, "fan-mgmt", defEvents[0].EventID)

	// The custom selector ignores category: the matching event name copies,
	// the other does not.
	customEvents, _, err := store.LookupEDSEvents(customEDS.EventDataStoreID, EDSQuery{})
	require.NoError(t, err)
	require.Len(t, customEvents, 1)
	assert.Equal(t, "fan-mgmt", customEvents[0].EventID)

	for _, eds := range []*EventDataStore{stoppedEDS, pendingEDS} {
		events, _, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{})
		require.NoError(t, err)
		assert.Empty(t, events, "store %s must not ingest", eds.Name)
	}

	// The event history keeps both events regardless of the stores.
	assert.Equal(t, 2, store.eventsStore.Count())
}

// TestLookupEDSEventsOrderBoundsAndToken pins the walk contract: ascending
// event-time order, both time bounds applied to the ordered keys, a
// pagination token that resumes strictly after the last served key, and an
// undecodable token refused as invalid.
func TestLookupEDSEventsOrderBoundsAndToken(t *testing.T) {
	store, _ := newEDSEventsTestStore(t)
	eds := NewEventDataStore("paged", "acc123", "us-east-1")
	_, err := store.CreateEventDataStore(eds)
	require.NoError(t, err)

	base := time.Now().UTC().Add(-time.Hour).Truncate(time.Millisecond)
	for i, id := range []string{"e1", "e2", "e3", "e4", "e5"} {
		require.NoError(t, store.PutEventIntoEDS(eds.EventDataStoreID,
			edsTestEvent(id, "CreateTrail", "Management", base.Add(time.Duration(i)*time.Minute))))
	}

	// Full walk, ascending.
	events, next, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{MaxResults: 2})
	require.NoError(t, err)
	require.Len(t, events, 2)
	assert.Equal(t, "e1", events[0].EventID)
	assert.NotEmpty(t, next)

	// Follow the token to exhaustion.
	var order []string
	for {
		events, next, err = store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{MaxResults: 2, NextToken: next})
		require.NoError(t, err)
		for _, e := range events {
			order = append(order, e.EventID)
		}
		if next == "" {
			break
		}
	}
	assert.Equal(t, []string{"e3", "e4", "e5"}, order)

	// Time bounds: the middle two — e2 (60s) and e3 (120s) — between a
	// 30-second start and a 150-second end.
	start := base.Add(30 * time.Second)
	end := base.Add(150 * time.Second)
	bounded, _, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{StartTime: &start, EndTime: &end})
	require.NoError(t, err)
	require.Len(t, bounded, 2)
	assert.Equal(t, "e2", bounded[0].EventID)
	assert.Equal(t, "e3", bounded[1].EventID)

	// The end bound is inclusive at EndTime itself and admits nothing
	// after it: a channel-style event (full nanosecond precision) half a
	// millisecond past the bound stays out, one at exactly the bound is
	// served.
	exact := end
	subMsLate := end.Add(900 * time.Microsecond)
	require.NoError(t, store.PutEventIntoEDS(eds.EventDataStoreID, edsTestEvent("e6", "CreateTrail", "Management", exact)))
	require.NoError(t, store.PutEventIntoEDS(eds.EventDataStoreID, edsTestEvent("e7", "CreateTrail", "Management", subMsLate)))
	bounded, _, err = store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{StartTime: &start, EndTime: &end})
	require.NoError(t, err)
	require.Len(t, bounded, 3)
	assert.Equal(t, "e6", bounded[2].EventID)

	_, _, err = store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{NextToken: "not-a-key"})
	assert.ErrorIs(t, err, ErrInvalidNextToken)

	// A key-shaped token minted by another event data store must be
	// refused: the shared "<nanos>#<id>" key shape would otherwise resume
	// this walk at an arbitrary point.
	other := NewEventDataStore("other", "acc123", "us-east-1")
	_, err = store.CreateEventDataStore(other)
	require.NoError(t, err)
	require.NoError(t, store.PutEventIntoEDS(other.EventDataStoreID,
		edsTestEvent("o1", "CreateTrail", "Management", base)))
	_, foreignNext, err := store.LookupEDSEvents(other.EventDataStoreID, EDSQuery{MaxResults: 1})
	require.NoError(t, err)
	require.NotEmpty(t, foreignNext)
	_, _, err = store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{NextToken: foreignNext})
	assert.ErrorIs(t, err, ErrInvalidNextToken)
}

// TestPurgeAndDropEDSEvents pins the per-store retention operations: the
// purge removes only events older than the cutoff and is idempotent; the
// drop removes everything and leaves the event history untouched.
func TestPurgeAndDropEDSEvents(t *testing.T) {
	store, _ := newEDSEventsTestStore(t)
	eds := NewEventDataStore("retained", "acc123", "us-east-1")
	_, err := store.CreateEventDataStore(eds)
	require.NoError(t, err)

	old := time.Now().UTC().Add(-48 * time.Hour).Truncate(time.Millisecond)
	fresh := time.Now().UTC().Add(-time.Hour).Truncate(time.Millisecond)
	for _, e := range []*Event{
		edsTestEvent("old-1", "CreateTrail", "Management", old),
		edsTestEvent("old-2", "DeleteTrail", "Management", old.Add(time.Minute)),
		edsTestEvent("fresh-1", "StartLogging", "Management", fresh),
	} {
		require.NoError(t, store.PutEventIntoEDS(eds.EventDataStoreID, e))
	}

	cutoff := time.Now().UTC().Add(-24 * time.Hour)
	removed, err := store.PurgeEDSEventsBefore(eds.EventDataStoreID, cutoff, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, removed)

	again, err := store.PurgeEDSEventsBefore(eds.EventDataStoreID, cutoff, 0)
	require.NoError(t, err)
	assert.Zero(t, again)

	remaining, _, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{})
	require.NoError(t, err)
	require.Len(t, remaining, 1)
	assert.Equal(t, "fresh-1", remaining[0].EventID)

	require.NoError(t, store.DropEDSEvents(eds.EventDataStoreID))
	dropped, _, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{})
	require.NoError(t, err)
	assert.Empty(t, dropped)
}

// TestRecordEventFromPayload pins the record-payload parsing: record
// spellings map onto the Event, the payload rides verbatim in
// CloudTrailEvent, a boolean readOnly becomes the wire string, an absent
// eventTime falls back to the ingestion time rather than rejecting the
// event, and the identity/category overrides carry the two callers'
// contracts — the channel path forces its assigned ID and the
// ActivityAuditLog category, the import path keeps the record's own values
// and defaults an absent category to Management.
func TestRecordEventFromPayload(t *testing.T) {
	payload := `{"eventVersion":"1.08","eventTime":"2026-09-17T10:00:00Z","eventSource":"example.partner.com",` +
		`"eventName":"SignIn","awsRegion":"us-east-1","readOnly":false,` +
		`"userIdentity":{"type":"IAMUser","userName":"partner"},"resources":[{"ARN":"arn:aws:s3:::b","type":"AWS::S3::Bucket"}]}`
	event, err := RecordEventFromPayload(payload, "assigned-id", "ActivityAuditLog")
	require.NoError(t, err)
	assert.Equal(t, "assigned-id", event.EventID)
	assert.Equal(t, "SignIn", event.EventName)
	assert.Equal(t, "example.partner.com", event.EventSource)
	assert.Equal(t, "us-east-1", event.AwsRegion)
	assert.Equal(t, "false", event.ReadOnly)
	assert.Equal(t, "ActivityAuditLog", event.EventCategory)
	assert.Equal(t, payload, event.CloudTrailEvent)
	require.NotNil(t, event.UserIdentity)
	assert.Equal(t, "partner", event.UserIdentity.UserName)
	require.Len(t, event.Resources, 1)
	assert.Equal(t, "arn:aws:s3:::b", event.Resources[0].ResourceName)
	wantTime, err := time.Parse(time.RFC3339, "2026-09-17T10:00:00Z")
	require.NoError(t, err)
	assert.Equal(t, wantTime, event.EventTime.UTC())

	noTime, err := RecordEventFromPayload(`{"eventName":"SignIn"}`, "id-2", "ActivityAuditLog")
	require.NoError(t, err)
	assert.False(t, noTime.EventTime.IsZero())

	// The import form: no overrides — the record's own eventID and category
	// survive, and a record without a category is a management event.
	imported, err := RecordEventFromPayload(
		`{"eventID":"record-own-id","eventCategory":"Data","eventName":"GetObject"}`, "", "")
	require.NoError(t, err)
	assert.Equal(t, "record-own-id", imported.EventID)
	assert.Equal(t, "Data", imported.EventCategory)
	bare, err := RecordEventFromPayload(`{"eventName":"ListBuckets"}`, "", "")
	require.NoError(t, err)
	assert.Equal(t, "Management", bare.EventCategory)
	assert.Empty(t, bare.EventID) // the write path assigns one

	_, err = RecordEventFromPayload("not json", "id-3", "ActivityAuditLog")
	assert.Error(t, err)
}

// TestPutEventIntoEDSRoundTrip pins the channel write path: an event
// written directly into a store returns through the walk with its key
// carried time and identifier, without passing the ingestion selectors.
func TestPutEventIntoEDSRoundTrip(t *testing.T) {
	store, _ := newEDSEventsTestStore(t)
	eds := NewEventDataStore("channel-target", "acc123", "us-east-1")
	eds.IngestionEnabled = false // channel writes bypass the ingestion gate
	_, err := store.CreateEventDataStore(eds)
	require.NoError(t, err)

	at := time.Now().UTC().Add(-time.Minute).Truncate(time.Millisecond)
	event, err := RecordEventFromPayload(
		`{"eventTime":"`+at.Format(time.RFC3339Nano)+`","eventName":"SignIn"}`, "audit-1", "ActivityAuditLog")
	require.NoError(t, err)
	require.NoError(t, store.PutEventIntoEDS(eds.EventDataStoreID, event))

	back, _, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{})
	require.NoError(t, err)
	require.Len(t, back, 1)
	assert.Equal(t, "audit-1", back[0].EventID)
	assert.Equal(t, at.UnixNano(), back[0].EventTime.UnixNano())
	assert.Equal(t, "ActivityAuditLog", back[0].EventCategory)
}

// TestPutEventIntoEDSsAtomicFanOut pins the channel fan-out's atomicity:
// every destination's copy joins one transaction, so a failure discovered
// at a later destination (here: a destination record that no longer
// exists) leaves no copy behind in the earlier destination either — a
// surviving partial copy would be duplicated by the request's retry,
// which carries a fresh event record.
func TestPutEventIntoEDSsAtomicFanOut(t *testing.T) {
	store, _ := newEDSEventsTestStore(t)
	first := NewEventDataStore("fan-first", "acc123", "us-east-1")
	second := NewEventDataStore("fan-second", "acc123", "us-east-1")
	for _, eds := range []*EventDataStore{first, second} {
		_, err := store.CreateEventDataStore(eds)
		require.NoError(t, err)
	}

	event := edsTestEvent("audit-atomic", "SignIn", "ActivityAuditLog",
		time.Now().UTC().Add(-time.Minute).Truncate(time.Millisecond))

	// The healthy fan-out copies the event into every destination.
	require.NoError(t, store.PutEventIntoEDSs([]string{first.EventDataStoreID, second.EventDataStoreID}, event))
	for _, eds := range []*EventDataStore{first, second} {
		events, _, err := store.LookupEDSEvents(eds.EventDataStoreID, EDSQuery{})
		require.NoError(t, err)
		assert.Len(t, events, 1, "destination %s must hold the copy", eds.Name)
	}

	// The first destination is emptied so the failure path's "no copy
	// anywhere" is observable per destination.
	require.NoError(t, store.DropEDSEvents(first.EventDataStoreID))

	// The second destination's record is gone: the batch write aborts and
	// the first destination keeps nothing.
	_, err := store.DeleteEventDataStoreIf(second.EventDataStoreID, func(*EventDataStore) error { return nil })
	require.NoError(t, err)

	err = store.PutEventIntoEDSs([]string{first.EventDataStoreID, second.EventDataStoreID},
		edsTestEvent("audit-atomic-2", "SignIn", "ActivityAuditLog", time.Now().UTC().Truncate(time.Millisecond)))
	require.ErrorIs(t, err, ErrEventDataStoreNotFound)

	firstEvents, _, err := store.LookupEDSEvents(first.EventDataStoreID, EDSQuery{})
	require.NoError(t, err)
	assert.Empty(t, firstEvents, "no destination may hold a copy of the aborted write")
}
