package cloudtrail

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below drive the store's single locking design: every record
// family's read-modify-write runs under the store mutex through the
// Mutate*/DeleteChannelIf surface, the channel-to-EDS delete guards fire off
// the Destination.Location association, and the async executors coordinate
// with CancelQuery/StopImport through status-checked transitions. The
// concurrency pins run under -race in the A-phase gate.

func newLockTestStore(t *testing.T) *cloudtrailstore.CloudTrailStore {
	t.Helper()
	tmpDir := "./tmp/cloudtrail-lock-test"
	os.RemoveAll(tmpDir)
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	s, err := storage.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open test storage: %v", err)
	}
	t.Cleanup(func() { s.Close() })

	return cloudtrailstore.NewCloudTrailStore(s, "acc123", "us-east-1")
}

// Concurrent field updates through the mutate surface must all land: the
// closure runs under the store mutex against the loaded record, so no
// update can be lost to a concurrent writer.
func TestConcurrentEDSMutationsAllLand(t *testing.T) {
	store := newLockTestStore(t)
	created, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("eds-concurrent", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("create eds: %v", err)
	}

	const writers = 8
	var wg sync.WaitGroup
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("writer-%d", i)
			_, err := store.MutateEventDataStore(created.EventDataStoreID, func(eds *cloudtrailstore.EventDataStore) error {
				if eds.Tags == nil {
					eds.Tags = map[string]string{}
				}
				eds.Tags[key] = key
				return nil
			})
			if err != nil {
				t.Errorf("writer %d mutate: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	final, err := store.GetEventDataStore(created.EventDataStoreID)
	if err != nil {
		t.Fatalf("get eds: %v", err)
	}
	for i := 0; i < writers; i++ {
		key := fmt.Sprintf("writer-%d", i)
		if final.Tags[key] != key {
			t.Fatalf("update from writer %d lost; tags = %v", i, final.Tags)
		}
	}
}

// Concurrent Core operations touching disjoint members of the same event
// data store (ingestion toggle, member update, federation) must all be
// visible on the stored record afterwards.
func TestConcurrentCoreUpdatesDisjointMembersAllLand(t *testing.T) {
	svc := &CloudTrailService{}
	store := newLockTestStore(t)

	eds := cloudtrailstore.NewEventDataStore("eds-disjoint", "acc123", "us-east-1")
	created, err := store.CreateEventDataStore(eds)
	if err != nil {
		t.Fatalf("create eds: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if _, err := svc.setEventDataStoreIngestion(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN}, false); err != nil {
			t.Errorf("ingestion toggle: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := svc.updateEventDataStoreCore(store, UpdateEventDataStoreInput{
			EventDataStore: created.EventDataStoreARN,
			KmsKeyId:       "arn:aws:kms:us-east-1:acc123:key/1234abcd-12ab-34cd-56ef-1234567890ab",
		}); err != nil {
			t.Errorf("member update: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		if _, err := svc.enableFederationCore(context.Background(), store, EnableFederationInput{
			EventDataStore:    created.EventDataStoreARN,
			FederationRoleArn: "arn:aws:iam::acc123:role/federation",
		}); err != nil {
			t.Errorf("federation enable: %v", err)
		}
	}()
	wg.Wait()

	final, err := store.GetEventDataStore(created.EventDataStoreID)
	if err != nil {
		t.Fatalf("get eds: %v", err)
	}
	if final.IngestionEnabled {
		t.Fatal("ingestion toggle lost")
	}
	if final.KMSKeyID == "" {
		t.Fatal("KMS key update lost")
	}
	if final.FederationStatus != "ENABLED" {
		t.Fatalf("federation update lost; status = %q", final.FederationStatus)
	}
}

// A cancel that lands while the executor runs must hold: the executor's
// finalisation refuses to write over CANCELLED and discards its results.
// The inverse order refuses the cancel instead.
func TestCancelQueryCoordinationBothOrders(t *testing.T) {
	svc := &CloudTrailService{}
	store := newLockTestStore(t)

	seed := func(id string) {
		if err := store.SaveQuery(&cloudtrailstore.QueryRecord{
			QueryID:        id,
			EventDataStore: "eds",
			QueryStatement: "SELECT eventID FROM eds",
			QueryStatus:    "RUNNING",
			StartTime:      time.Now().UTC(),
		}); err != nil {
			t.Fatalf("seed query %s: %v", id, err)
		}
	}

	// Cancel first, then the executor's late completion.
	seed("q-cancel-first")
	resp, err := svc.cancelQueryCore(store, CancelQueryInput{QueryID: "q-cancel-first"})
	if err != nil {
		t.Fatalf("cancel: %v", err)
	}
	if resp["QueryStatus"] != "CANCELLED" {
		t.Fatalf("cancel status = %v, want CANCELLED", resp["QueryStatus"])
	}
	svc.finaliseQueryExecution(store, "q-cancel-first", [][]map[string]string{{{"eventID": "e1"}}}, nil, false, nil)
	qr, err := store.GetQuery("q-cancel-first")
	if err != nil {
		t.Fatalf("get query: %v", err)
	}
	if qr.QueryStatus != "CANCELLED" {
		t.Fatalf("executor overwrote CANCELLED with %q", qr.QueryStatus)
	}
	if len(qr.QueryResultRows) != 0 {
		t.Fatalf("cancelled query must not store results; got %d rows", len(qr.QueryResultRows))
	}

	// Completion first, then the late cancel refuses — with the model's
	// declared InactiveQueryException ("The specified query cannot be
	// canceled because it is in the FINISHED, FAILED, TIMED_OUT, or
	// CANCELLED state").
	seed("q-finish-first")
	svc.finaliseQueryExecution(store, "q-finish-first", nil, nil, false, nil)
	if _, err := svc.cancelQueryCore(store, CancelQueryInput{QueryID: "q-finish-first"}); err == nil {
		t.Fatal("cancel after FINISHED must be refused")
	} else {
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.GetCode() != "InactiveQueryException" {
			t.Fatalf("cancel after FINISHED: error = %v, want InactiveQueryException", err)
		}
	}
}

// A stop that lands before the import executor's transitions must hold:
// runImport refuses both the IN_PROGRESS and the COMPLETED write over
// STOPPED.
func TestStopImportCoordinationHoldsAgainstRunImport(t *testing.T) {
	svc := &CloudTrailService{}
	store := newLockTestStore(t)

	created, err := store.CreateImport(cloudtrailstore.NewImport(
		[]string{"arn:aws:cloudtrail:us-east-1:acc123:eventdatastore/uuid"},
		cloudtrailstore.ImportSource{S3LocationURI: "s3://bucket/prefix", S3BucketRegion: "us-east-1"},
	))
	if err != nil {
		t.Fatalf("create import: %v", err)
	}

	if _, err := svc.stopImportCore(store, ImportIDInput{ImportID: created.ImportID}); err != nil {
		t.Fatalf("stop import: %v", err)
	}

	svc.runImport(store, created.ImportID)

	imp, err := store.GetImport(created.ImportID)
	if err != nil {
		t.Fatalf("get import: %v", err)
	}
	if imp.ImportStatus != "STOPPED" {
		t.Fatalf("runImport overwrote STOPPED with %q", imp.ImportStatus)
	}
	if imp.ImportStatistics.PrefixesCompleted != 0 || imp.ImportStatistics.FilesCompleted != 0 {
		t.Fatalf("stopped import must not record completion statistics; got %+v", imp.ImportStatistics)
	}
}

// DeleteChannel documents no event-data-store precondition: a channel
// deletes even while its EVENT_DATA_STORE destination resolves to an
// ENABLED store. The association blocks in the other direction — the test
// below pins DeleteEventDataStore refusing while a channel points at the
// store — so teardown order is channel first, store second.
func TestDeleteChannelWithLiveEDSDestination(t *testing.T) {
	svc := &CloudTrailService{}
	store := newLockTestStore(t)

	eds := cloudtrailstore.NewEventDataStore("eds-guard", "acc123", "us-east-1")
	created, err := store.CreateEventDataStore(eds)
	if err != nil {
		t.Fatalf("create eds: %v", err)
	}

	ch := cloudtrailstore.NewChannel("ch-guard", "Custom", "acc123", "us-east-1")
	ch.Destinations = []cloudtrailstore.Destination{{
		Type:     cloudtrailstore.DestinationTypeEventDataStore,
		Location: created.EventDataStoreARN,
	}}
	createdCh, err := store.CreateChannel(ch)
	if err != nil {
		t.Fatalf("create channel: %v", err)
	}

	if _, err := svc.deleteChannelCore(store, ChannelInput{Channel: createdCh.ChannelARN}); err != nil {
		t.Fatalf("delete channel with live EDS destination: %v", err)
	}
	if _, err := store.GetChannel(createdCh.ChannelARN); err == nil {
		t.Fatal("channel still exists after successful delete")
	}
}

// The event-data-store delete guard fires off the same association from the
// other side: a channel destination pointing at the event data store
// refuses its deletion with ChannelExistsForEDSException.
func TestDeleteEDSGuardFiresForAssociatedChannel(t *testing.T) {
	svc := &CloudTrailService{}
	store := newLockTestStore(t)

	eds := cloudtrailstore.NewEventDataStore("eds-guard-2", "acc123", "us-east-1")
	eds.TerminationProtectionEnabled = false
	created, err := store.CreateEventDataStore(eds)
	if err != nil {
		t.Fatalf("create eds: %v", err)
	}

	ch := cloudtrailstore.NewChannel("ch-guard-2", "Custom", "acc123", "us-east-1")
	ch.Destinations = []cloudtrailstore.Destination{{
		Type:     cloudtrailstore.DestinationTypeEventDataStore,
		Location: created.EventDataStoreARN,
	}}
	if _, err := store.CreateChannel(ch); err != nil {
		t.Fatalf("create channel: %v", err)
	}

	_, err = svc.deleteEventDataStoreCore(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN})
	assertErrorCode(t, err, "ChannelExistsForEDSException", "delete EDS with associated channel")

	// Deleting the channel first clears the association and the delete
	// goes through (the destination retires with the EDS still ENABLED —
	// the channel guard allows deletion once the channel itself is gone).
	if err := store.DeleteChannelIf(ch.ChannelARN, func(*cloudtrailstore.Channel) error { return nil }); err != nil {
		t.Fatalf("delete channel: %v", err)
	}
	if _, err := svc.deleteEventDataStoreCore(store, EventDataStoreIDInput{EventDataStore: created.EventDataStoreARN}); err != nil {
		t.Fatalf("delete EDS after channel removed: %v", err)
	}
}

// Concurrent channel renames serialise on the store mutex: the final name
// is one of the contenders and the record stays consistent.
func TestConcurrentChannelRenamesSerialise(t *testing.T) {
	store := newLockTestStore(t)
	created, err := store.CreateChannel(cloudtrailstore.NewChannel("ch-race", "Custom", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("create channel: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := store.MutateChannel(created.ChannelARN, func(ch *cloudtrailstore.Channel) error {
				ch.Name = fmt.Sprintf("ch-race-%d", i)
				return nil
			})
			if err != nil {
				t.Errorf("rename %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()

	final, err := store.GetChannel(created.ChannelARN)
	if err != nil {
		t.Fatalf("get channel: %v", err)
	}
	if final.Name == "ch-race" {
		t.Fatal("no rename landed")
	}
}

// assertErrorCode pins the wire code of a Core refusal.
func assertErrorCode(t *testing.T, err error, code string, context string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected %s, got success", context, code)
	}
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) {
		t.Fatalf("%s: error %v is not an AWS error", context, err)
	}
	if awsErr.GetCode() != code {
		t.Fatalf("%s: code = %q, want %q", context, awsErr.GetCode(), code)
	}
}
