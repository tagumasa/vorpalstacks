package dynamodb

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestKinesisRecordForDestinationAppliesPerDestinationPrecision(t *testing.T) {
	table := guardTestTable()
	now := time.Unix(1700000000, 123456789)
	keys := map[string]*dbstore.AttributeValue{
		"pk": matchSAttr("hashval"),
		"sk": matchSAttr("sortval"),
	}

	milli := &dbstore.KinesisDataStreamDestination{DestinationStatus: kinesisDestinationActive}
	micro := &dbstore.KinesisDataStreamDestination{
		DestinationStatus:                    kinesisDestinationActive,
		ApproximateCreationDateTimePrecision: dbstore.ACDTPrecisionMicrosecond,
	}

	milliPayload, partitionKey := kinesisRecordForDestination(milli, table, dbstore.StreamEventInsert, keys, nil, nil, now)
	if partitionKey != "hashval" {
		t.Fatalf("partition key = %q, want the HASH attribute value", partitionKey)
	}
	microPayload, _ := kinesisRecordForDestination(micro, table, dbstore.StreamEventInsert, keys, nil, nil, now)

	var milliRecord kinesisDestinationEnvelope
	if err := json.Unmarshal(milliPayload, &milliRecord); err != nil {
		t.Fatalf("decode millisecond payload: %v", err)
	}
	if milliRecord.EventName != "INSERT" {
		t.Fatalf("eventName = %s, want INSERT", milliRecord.EventName)
	}
	if milliRecord.DynamoDB.ApproximateCreationDateTimePrecision != string(dbstore.ACDTPrecisionMillisecond) {
		t.Fatalf("default precision = %s, want MILLISECOND", milliRecord.DynamoDB.ApproximateCreationDateTimePrecision)
	}
	if milliRecord.DynamoDB.ApproximateCreationDateTime != now.UnixMilli() {
		t.Fatalf("millisecond timestamp = %d, want %d", milliRecord.DynamoDB.ApproximateCreationDateTime, now.UnixMilli())
	}
	if s, ok := milliRecord.DynamoDB.Keys["pk"].(map[string]interface{}); !ok || s["S"] != "hashval" {
		t.Fatalf("keys not carried in the dynamodb envelope: %v", milliRecord.DynamoDB.Keys)
	}

	var microRecord kinesisDestinationEnvelope
	if err := json.Unmarshal(microPayload, &microRecord); err != nil {
		t.Fatalf("decode microsecond payload: %v", err)
	}
	if microRecord.DynamoDB.ApproximateCreationDateTimePrecision != string(dbstore.ACDTPrecisionMicrosecond) {
		t.Fatalf("microsecond destination precision = %s, want MICROSECOND", microRecord.DynamoDB.ApproximateCreationDateTimePrecision)
	}
	if microRecord.DynamoDB.ApproximateCreationDateTime != now.UnixMicro() {
		t.Fatalf("microsecond timestamp = %d, want %d", microRecord.DynamoDB.ApproximateCreationDateTime, now.UnixMicro())
	}
}

func TestExtractPartitionKeyForKinesisUsesHashAttribute(t *testing.T) {
	keys := map[string]*dbstore.AttributeValue{
		"pk": matchSAttr("hashval"),
		"sk": matchSAttr("aaa"),
	}
	// The key must be selected by name, never by map iteration order.
	for i := 0; i < 100; i++ {
		if got := extractPartitionKeyForKinesis(keys, "pk"); got != "hashval" {
			t.Fatalf("iteration %d: partition key = %q, want hashval", i, got)
		}
	}

	num := "123"
	if got := extractPartitionKeyForKinesis(map[string]*dbstore.AttributeValue{"pk": {N: &num}}, "pk"); got != "123" {
		t.Fatalf("numeric partition key = %q, want 123", got)
	}
	if got := extractPartitionKeyForKinesis(map[string]*dbstore.AttributeValue{}, "pk"); got != "default" {
		t.Fatalf("empty keys partition key = %q, want default", got)
	}
}

// TestReplicationRetriesTransientFailures pins the bounded-retry delivery
// contract: a replica write that fails transiently is retried inside the
// same delivery and converges, and a persistently failing write exhausts
// the in-flight delivery's bounded attempts and escalates to the
// convergence worker — the attempts continue past the bound, the
// eventual-consistency posture AWS's global tables run.
func TestReplicationRetriesTransientFailures(t *testing.T) {
	svc, eastCtx, _, _ := globalTablePlaneFixture(t)
	createPinGlobalTable(t, svc, eastCtx)
	eastStore, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("east store: %v", err)
	}

	waitForCalls := func(target int32, calls *atomic.Int32) bool {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			if calls.Load() >= target {
				return true
			}
			time.Sleep(5 * time.Millisecond)
		}
		return false
	}

	// A transient fault on the first two attempts heals on the third.
	var healed atomic.Int32
	svc.replicateToGlobalTableReplicas(eastStore, "us-east-1", "GtPinTable", func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		if healed.Add(1) <= 2 {
			return errors.New("transient store fault")
		}
		return nil
	})
	if !waitForCalls(3, &healed) {
		t.Fatalf("transient fault never healed: %d attempts", healed.Load())
	}

	// A persistent fault exhausts the in-flight delivery's bounded
	// attempts and escalates: the convergence worker keeps the delivery
	// trying past the bound instead of abandoning it.
	var persistent atomic.Int32
	svc.replicateToGlobalTableReplicas(eastStore, "us-east-1", "GtPinTable", func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		persistent.Add(1)
		return errors.New("persistent store fault")
	})
	if !waitForCalls(replicateRetryAttempts, &persistent) {
		t.Fatalf("persistent fault not attempted %d times: %d", replicateRetryAttempts, persistent.Load())
	}
	// The escalation lands after the bounded attempt returns — the attempt
	// counter observes the fault function's entry, not its aftermath — so
	// the check waits for the enqueue instead of racing it.
	escalationDeadline := time.Now().Add(5 * time.Second)
	for svc.pendingEscalationCount() == 0 && time.Now().Before(escalationDeadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if svc.pendingEscalationCount() == 0 {
		t.Fatal("exhausted delivery never escalated")
	}
	if !waitForCalls(replicateRetryAttempts*3, &persistent) {
		t.Fatalf("escalated delivery stopped retrying: %d attempts", persistent.Load())
	}
}

// TestCloseStopsALateEscalationWorker pins the escalation worker's lifetime
// binding: the worker derives from the service's background context, so a
// worker spawned after the shutdown cancellation has already run still
// observes it and exits — the background WaitGroup settles instead of
// waiting forever on a worker nothing can stop.
func TestCloseStopsALateEscalationWorker(t *testing.T) {
	svc := &DynamoDBService{}
	bgCtx, bgCancel := context.WithCancel(context.Background())
	svc.bgCtx, svc.bgCancel = bgCtx, bgCancel
	// The cancellation runs before the delivery exists: the worker the
	// delivery spawns is born after cancellation, the ordering a Close
	// racing a late delivery produces.
	bgCancel()

	svc.bgWg.Add(1)
	go func() {
		defer svc.bgWg.Done()
		svc.escalateReplicationDelivery("LateTable", "us-west-2", func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
			return nil
		}, errors.New("delivery exhausted"))
	}()

	settled := make(chan struct{})
	go func() {
		svc.bgWg.Wait()
		close(settled)
	}()
	select {
	case <-settled:
	case <-time.After(5 * time.Second):
		t.Fatal("background wait-group never settled: a worker spawned after the shutdown cancellation never stopped")
	}

	// A service built without the constructor closes cleanly too: the
	// lifecycle the worker binds to is ensured lazily, so Close cancels a
	// real context rather than reading a nil.
	bare := &DynamoDBService{}
	bareClosed := make(chan struct{})
	go func() {
		bare.Close()
		close(bareClosed)
	}()
	select {
	case <-bareClosed:
	case <-time.After(5 * time.Second):
		t.Fatal("Close never returned on a service built without the constructor")
	}
}

// closeWaitRegressionWindow bounds the Close-wait pins' negative halves:
// while a delivery is parked, Close must not return, and a Close that
// skips the wait-group returns inside this window — the bounded window is
// the only witness available for that negative without production hooks,
// so it stays minimal rather than standing in for synchronisation.
const closeWaitRegressionWindow = 150 * time.Millisecond

// blockingKinesisInvoker stalls the PutRecord delivery until the test
// releases it, so the WaitGroup pin can observe whether Close waits for an
// in-flight emit.
type blockingKinesisInvoker struct {
	started chan struct{}
	release chan struct{}
}

func (b blockingKinesisInvoker) ListShards(ctx context.Context, region, streamName string) ([]invokers.ShardInfo, error) {
	return nil, errors.New("unused in this pin")
}

func (b blockingKinesisInvoker) PutRecord(ctx context.Context, region, streamName, partitionKey string, data []byte) (string, string, error) {
	close(b.started)
	<-b.release
	return "1", "shardId-00000001741631711871-1f6a72cf", nil
}

func (b blockingKinesisInvoker) CreateShardIterator(ctx context.Context, region, streamName string, shardID string, iteratorType string, startingSequenceNumber string, timestamp *time.Time) (string, error) {
	return "", errors.New("unused in this pin")
}

func (b blockingKinesisInvoker) GetRecords(ctx context.Context, region, streamName string, shardID string, startingSequenceNumber string, limit int32, includeStart bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", errors.New("unused in this pin")
}

func (b blockingKinesisInvoker) StreamExists(ctx context.Context, region, streamARN string) (bool, error) {
	return false, errors.New("unused in this pin")
}

// TestBackgroundDeliveriesHoldTheServiceOpen pins the shutdown contract of
// the two asynchronous delivery workers: the Kinesis-destination emit and
// the global-table replication register on the service's background
// WaitGroup, so Close waits for an in-flight delivery to finish instead of
// killing it mid-flight — each worker's own bounded timeout keeps that wait
// finite. A parked delivery holds Close open; releasing it lets Close
// return.
func TestBackgroundDeliveriesHoldTheServiceOpen(t *testing.T) {
	// The Kinesis emit half: one active destination, an invoker whose
	// PutRecord parks until the test releases it.
	bus := eventbus.NewEventBus()
	emitStarted := make(chan struct{})
	emitRelease := make(chan struct{})
	bus.SetKinesisInvoker(blockingKinesisInvoker{started: emitStarted, release: emitRelease})
	emitSvc := NewDynamoDBService("123456789012")
	emitSvc.SetEventBus(bus)
	destTable := &dbstore.Table{
		Name: "WgEmitTable",
		KinesisDataStreamDestinations: []*dbstore.KinesisDataStreamDestination{{
			StreamArn:         "arn:aws:kinesis:us-east-1:123456789012:stream/WgDest",
			DestinationStatus: kinesisDestinationActive,
		}},
	}
	emitSvc.sendToKinesisDestinations(destTable, dbstore.StreamEventInsert, map[string]*dbstore.AttributeValue{"pk": matchSAttr("hashval")}, nil, nil)
	<-emitStarted
	emitClosed := make(chan struct{})
	go func() {
		emitSvc.Close()
		close(emitClosed)
	}()
	// The negative half of the Close wait: while the emit is parked, Close
	// must not return. The channel is the observation; the bounded window
	// is the only available witness for the negative (a Close that skips
	// the wait-group returns inside it, an unregistered-slot regression),
	// so it stays minimal rather than standing in for synchronisation.
	select {
	case <-emitClosed:
		t.Fatal("Close returned while a Kinesis destination emit was still in flight")
	case <-time.After(closeWaitRegressionWindow):
	}
	close(emitRelease)
	select {
	case <-emitClosed:
	case <-time.After(5 * time.Second):
		t.Fatal("Close never returned after the Kinesis emit completed")
	}

	// The replication half: a two-region global record, a delivery callback
	// that parks until released.
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	t.Cleanup(func() { sm.Close() })
	repSvc := NewDynamoDBService("123456789012")
	repSvc.SetStorageManager(sm)
	store, err := repSvc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.GlobalTables().Create("WgReplTable", []*dbstore.Replica{
		{RegionName: "us-east-1"},
		{RegionName: "eu-west-1"},
	}); err != nil {
		t.Fatalf("create global record: %v", err)
	}
	repStarted := make(chan struct{})
	repRelease := make(chan struct{})
	repSvc.replicateToGlobalTableReplicas(store, "us-east-1", "WgReplTable", func(ctx context.Context, destStore dbstore.DynamoDBStoreInterface) error {
		close(repStarted)
		<-repRelease
		return nil
	})
	<-repStarted
	repClosed := make(chan struct{})
	go func() {
		repSvc.Close()
		close(repClosed)
	}()
	// The same negative half as the Kinesis emit above, for the
	// replication delivery: bounded window, channel observation.
	select {
	case <-repClosed:
		t.Fatal("Close returned while a global-table replication delivery was still in flight")
	case <-time.After(closeWaitRegressionWindow):
	}
	close(repRelease)
	select {
	case <-repClosed:
	case <-time.After(5 * time.Second):
		t.Fatal("Close never returned after the replication delivery completed")
	}
}
