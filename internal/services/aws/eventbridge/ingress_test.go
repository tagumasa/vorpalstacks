package eventbridge

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// newIngressStore returns a store with one event bus, one ENABLED archive on
// it, and — when disabledArchive is set — one DISABLED archive, so ingress
// rows can observe which archives captured the event.
func newIngressStore(t *testing.T, busName string, disabledArchive bool) *eventsstore.EventsStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	store := eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
	ctx := context.Background()
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: busName}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateArchive(ctx, &eventsstore.Archive{
		Name:         "ingress-archive",
		EventBusName: busName,
		State:        eventsstore.ArchiveStateEnabled,
	}); err != nil {
		t.Fatal(err)
	}
	if disabledArchive {
		if err := store.CreateArchive(ctx, &eventsstore.Archive{
			Name:         "ingress-archive-off",
			EventBusName: busName,
			State:        eventsstore.ArchiveStateDisabled,
		}); err != nil {
			t.Fatal(err)
		}
	}
	return store
}

func archivedCount(t *testing.T, store *eventsstore.EventsStore, archiveName string) int {
	t.Helper()
	count := 0
	nextToken := ""
	for {
		result, err := store.ListArchiveEvents(context.Background(), archiveName, time.Time{}, time.Now().Add(time.Hour), eventsstore.ListLimitMaximum, nextToken)
		if err != nil {
			t.Fatalf("failed to read archive %s: %v", archiveName, err)
		}
		count += len(result.Events)
		if result.NextToken == "" {
			return count
		}
		nextToken = result.NextToken
	}
}

// TestPutEventsCoreIngressMatrix pins the HTTP ingress plane: a valid entry
// succeeds and reaches the bus's enabled archives (never the disabled ones);
// the ARN form of EventBusName resolves to the canonical bus; and a
// non-existent bus is not an entry failure — the documented AWS behaviour is
// a 200 response whose matching finds no rules and drops the event.
func TestPutEventsCoreIngressMatrix(t *testing.T) {
	ctx := context.Background()
	svc := NewEventsService(nil, "000000000000")

	busName := "ingress-http"
	store := newIngressStore(t, busName, true)

	validEntry := func(bus string) map[string]interface{} {
		return map[string]interface{}{
			"Source":       "com.example.ingress",
			"DetailType":   "IngressTest",
			"Detail":       `{"k":"v"}`,
			"EventBusName": bus,
		}
	}

	// valid × archive-on: the entry succeeds and only the ENABLED archive
	// captures the event.
	result, err := svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{validEntry(busName)}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.FailedEntryCount != 0 || len(result.Entries) != 1 {
		t.Fatalf("valid entry must succeed, got %+v", result)
	}
	if result.Entries[0]["EventId"] == nil || result.Entries[0]["EventId"].(string) == "" {
		t.Fatalf("valid entry must return an EventId, got %+v", result.Entries[0])
	}
	if got := archivedCount(t, store, "ingress-archive"); got != 1 {
		t.Fatalf("the enabled archive must capture the ingressed event, got %d events", got)
	}
	if got := archivedCount(t, store, "ingress-archive-off"); got != 0 {
		t.Fatalf("the disabled archive must stay empty, got %d events", got)
	}

	// ARN-form EventBusName resolves to the canonical bus: the archive is
	// keyed by the bus name, so capture proves the resolution.
	arn := "arn:aws:events:us-east-1:000000000000:event-bus/" + busName
	result, err = svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{validEntry(arn)}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.FailedEntryCount != 0 {
		t.Fatalf("an ARN-form EventBusName must be accepted, got %+v", result)
	}
	if got := archivedCount(t, store, "ingress-archive"); got != 2 {
		t.Fatalf("the ARN form must resolve to the canonical bus and be archived, got %d events", got)
	}

	// bus-missing: the documented AWS drop — a 200-equivalent success entry,
	// not a failed one.
	result, err = svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{validEntry("no-such-bus")}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.FailedEntryCount != 0 {
		t.Fatalf("a non-existent bus must not fail the entry (documented drop), got %+v", result)
	}
	if result.Entries[0]["EventId"] == nil {
		t.Fatalf("a dropped-bus entry still reports an EventId, got %+v", result.Entries[0])
	}
}

// TestPutEventsEntryContract pins the per-entry error vocabulary and the
// request-level rules: a malformed Detail fails its entry as MalformedDetail,
// an incomplete entry fails as InvalidArgument while complete siblings
// succeed, a request with no complete entry fails entirely, and the
// entry-count and summed-size ceilings reject the whole request.
func TestPutEventsEntryContract(t *testing.T) {
	ctx := context.Background()
	svc := NewEventsService(nil, "000000000000")
	store := newIngressStore(t, "contract-bus", false)

	completeEntry := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"DetailType":   "ContractTest",
		"Detail":       `{"k":"v"}`,
	}

	// Malformed Detail → MalformedDetail; the complete sibling succeeds.
	malformed := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"DetailType":   "ContractTest",
		"Detail":       `not-json{`,
	}
	result, err := svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{malformed, completeEntry}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.FailedEntryCount != 1 {
		t.Fatalf("expected one failed entry, got %d", result.FailedEntryCount)
	}
	if code, _ := result.Entries[0]["ErrorCode"].(string); code != "MalformedDetail" {
		t.Fatalf("malformed Detail must fail as MalformedDetail, got %+v", result.Entries[0])
	}
	if result.Entries[1]["EventId"] == nil {
		t.Fatalf("the complete sibling must succeed, got %+v", result.Entries[1])
	}

	// Non-object Detail (a JSON array) is also MalformedDetail.
	arrayDetail := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"DetailType":   "ContractTest",
		"Detail":       `[1,2]`,
	}
	result, err = svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{arrayDetail, completeEntry}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if code, _ := result.Entries[0]["ErrorCode"].(string); code != "MalformedDetail" {
		t.Fatalf("a non-object Detail must fail as MalformedDetail, got %+v", result.Entries[0])
	}

	// The null literal unmarshals into a nil map without an error, yet it
	// is not a JSON object — MalformedDetail like any other scalar.
	nullDetail := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"DetailType":   "ContractTest",
		"Detail":       `null`,
	}
	result, err = svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{nullDetail, completeEntry}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if code, _ := result.Entries[0]["ErrorCode"].(string); code != "MalformedDetail" {
		t.Fatalf("a null Detail must fail as MalformedDetail, got %+v", result.Entries[0])
	}

	// An entry missing one of the trio fails as InvalidArgument.
	noDetailType := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"Detail":       `{"k":"v"}`,
	}
	result, err = svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{noDetailType, completeEntry}, Region: "us-east-1"})
	if err != nil {
		t.Fatal(err)
	}
	if code, _ := result.Entries[0]["ErrorCode"].(string); code != "InvalidArgument" {
		t.Fatalf("an incomplete entry must fail as InvalidArgument, got %+v", result.Entries[0])
	}

	// No complete entry in the request → the entire request fails.
	if _, err := svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{noDetailType}, Region: "us-east-1"}); err == nil {
		t.Fatal("a request with no complete entry must fail entirely")
	}

	// Entry count over the bound → the entire request fails.
	many := make([]interface{}, 0, eventsstore.MaxPutEventsEntries+1)
	for i := 0; i <= eventsstore.MaxPutEventsEntries; i++ {
		many = append(many, completeEntry)
	}
	if _, err := svc.putEventsCore(ctx, store, PutEventsInput{Entries: many, Region: "us-east-1"}); err == nil {
		t.Fatal("an over-count request must fail entirely")
	}

	// Summed entry size at the ceiling → the entire request fails; one
	// byte under passes.
	bigDetail := `{"pad":"` + string(make([]byte, eventsstore.MaxPutEventsRequestSizeBytes)) + `"}`
	bigEntry := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"DetailType":   "ContractTest",
		"Detail":       bigDetail,
	}
	if _, err := svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{bigEntry}, Region: "us-east-1"}); err == nil {
		t.Fatal("a request at the size ceiling must fail entirely")
	}
	fits := map[string]interface{}{
		"EventBusName": "contract-bus",
		"Source":       "com.example.contract",
		"DetailType":   "ContractTest",
		"Detail":       `{"pad":"` + string(make([]byte, eventsstore.MaxPutEventsRequestSizeBytes-64)) + `"}`,
	}
	if _, err := svc.putEventsCore(ctx, store, PutEventsInput{Entries: []interface{}{fits}, Region: "us-east-1"}); err != nil {
		t.Fatalf("a request one allowance under the ceiling must pass: %v", err)
	}
}

// TestPutEventsBusHandlerIngressMatrix pins the internal bus ingress plane:
// a valid event is archived onto the named bus like the HTTP plane's ingress;
// the ARN form resolves; and the failure modes that the PutEvents API plane
// answers per entry — a missing bus, a malformed Detail, an out-of-bound
// DetailType or Source, non-string Resources — fail the handler so the
// publisher's retry and dead-letter policy observes them.
func TestPutEventsBusHandlerIngressMatrix(t *testing.T) {
	ctx := context.Background()

	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	busName := "ingress-bus"
	busARN := "arn:aws:events:us-east-1:000000000000:event-bus/" + busName

	// The store is wired through GetStoreForRegion, which seeds the default
	// bus; the named bus is created explicitly.
	seed := newIngressStore(t, busName, false)
	svc := NewEventsService(mgr, "000000000000")
	svc.SetEventsStore("us-east-1", seed)

	putEvent := func(bus, input string) eventbus.HandlerResult {
		evt := &eventbus.EventBridgePutEventsEvent{Input: input, EventBusName: bus}
		evt.Region = "us-east-1"
		return svc.handlePutEventsEvent(ctx, evt)
	}

	// valid × archive-on: the wrapped Detail form is archived.
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"test","Detail":{"k":"v"}}`); res.Error != nil {
		t.Fatalf("valid bus ingress must succeed: %v", res.Error)
	}
	if got := archivedCount(t, seed, "ingress-archive"); got != 1 {
		t.Fatalf("bus ingress must archive like the HTTP plane, got %d events", got)
	}

	// The Input-as-detail form (S3 notifications, scheduler target input)
	// also archives.
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"test","note":"payload-as-detail"}`); res.Error != nil {
		t.Fatalf("payload-as-detail ingress must succeed: %v", res.Error)
	}
	if got := archivedCount(t, seed, "ingress-archive"); got != 2 {
		t.Fatalf("the payload-as-detail form must archive, got %d events", got)
	}

	// ARN-form EventBusName resolves to the canonical bus.
	if res := putEvent(busARN, `{"Source":"aws.test","DetailType":"test","Detail":{}}`); res.Error != nil {
		t.Fatalf("ARN-form bus ingress must succeed: %v", res.Error)
	}
	if got := archivedCount(t, seed, "ingress-archive"); got != 3 {
		t.Fatalf("the ARN form must resolve to the canonical bus, got %d events", got)
	}

	// bus-missing: the handler failure that makes the drop observable to
	// the publisher.
	if res := putEvent("no-such-bus", `{"Source":"aws.test","DetailType":"test","Detail":{}}`); res.Error == nil {
		t.Fatal("a bus event targeting a non-existent bus must fail the handler")
	}

	// invalid rows.
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"test","Detail":"not-json{"}`); res.Error == nil {
		t.Fatal("a malformed string Detail must fail the handler")
	}
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"test","Detail":"42"}`); res.Error == nil {
		t.Fatal("a non-object Detail must fail the handler")
	}
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"test","Detail":"null"}`); res.Error == nil {
		t.Fatal("a null-literal Detail must fail the handler")
	}
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"`+string(make([]byte, 129))+`","Detail":{}}`); res.Error == nil {
		t.Fatal("an over-long DetailType must fail the handler")
	}
	if res := putEvent(busName, `{"Source":"aws.test","DetailType":"test","Detail":{},"Resources":[42]}`); res.Error == nil {
		t.Fatal("a non-string Resources member must fail the handler")
	}
	if got := archivedCount(t, seed, "ingress-archive"); got != 3 {
		t.Fatalf("failed ingress rows must not archive, got %d events", got)
	}
}
