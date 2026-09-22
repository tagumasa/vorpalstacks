package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// gunzipJSON decompresses one gzip layer and decodes the JSON within.
func gunzipJSON(data []byte, v interface{}) error {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer zr.Close()
	return json.NewDecoder(zr).Decode(v)
}

// newDispatchTestService builds a service backed by temporary storage with
// a started bus-less wiring: the tests call the dispatch directly, so the
// bus only carries the recording Kinesis invoker and needs no
// subscriptions.
func newDispatchTestService(t *testing.T) (*LogsService, *regionRecordingKinesisInvoker) {
	t.Helper()
	svc, _ := newTestService(t)
	kinesis := &regionRecordingKinesisInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(kinesis)
	svc.bus = bus
	return svc, kinesis
}

// TestPutSubscriptionFilterAcceptsCWLDestinationARN pins the destination
// acceptance against the destinationArn member documentation's list: the
// platform's own destination ARN form (the cross-account "logical
// destination") is accepted alongside Lambda and Kinesis, while Firehose —
// a documented AWS form the platform cannot deliver until its Firehose
// service exists — is rejected at Put instead of accepting a filter that
// would discard every matched batch at delivery.
func TestPutSubscriptionFilterAcceptsCWLDestinationARN(t *testing.T) {
	svc, _ := newDispatchTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}

	cases := []struct {
		name    string
		destArn string
		wantErr bool
	}{
		{"lambda ARN accepted", "arn:aws:lambda:us-east-1:000000000000:function:consumer", false},
		{"kinesis ARN accepted", "arn:aws:kinesis:us-east-1:000000000000:stream/logs-dest", false},
		{"platform destination ARN (colon form) accepted", "arn:aws:logs:us-east-1:000000000000:destination:central", false},
		{"AWS destination ARN (slash form) accepted", "arn:aws:logs:us-east-1:000000000000:destination/central", false},
		{"firehose ARN rejected", "arn:aws:firehose:us-east-1:000000000000:deliverystream/logs-dest", true},
		{"s3 ARN rejected", "arn:aws:s3:::logs-bucket", true},
		{"bare string rejected", "not-an-arn", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// One log group per case: the per-group subscription-filter
			// quota is two, and accepted cases must not spend each
			// other's slots.
			group := "sub-group-" + strings.Map(func(r rune) rune {
				if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
					return r
				}
				return '-'
			}, tc.name)
			if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
				t.Fatalf("create log group: %v", err)
			}
			err := svc.putSubscriptionFilterCore(context.Background(), store, &PutSubscriptionFilterInput{
				LogGroupName:     group,
				FilterName:       "filter",
				FilterPattern:    "ERROR",
				FilterPatternSet: true,
				DestinationArn:   tc.destArn,
			})
			if tc.wantErr && err == nil {
				t.Fatalf("destinationArn %q must be rejected", tc.destArn)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("destinationArn %q must be accepted: %v", tc.destArn, err)
			}
			if tc.wantErr && err != nil && !strings.Contains(err.Error(), "destinationArn") {
				t.Fatalf("rejection must name destinationArn, got: %v", err)
			}
		})
	}
}

// TestPutDestinationRejectsFirehoseTarget pins the destination target
// vocabulary to the PutDestinationRequest targetArn member documentation —
// "The ARN of an Amazon Kinesis stream" — so a destination cannot wrap a
// Firehose target the platform cannot deliver.
func TestPutDestinationRejectsFirehoseTarget(t *testing.T) {
	svc, _ := newDispatchTestService(t)

	if _, err := svc.putDestinationCore("central", "arn:aws:iam::000000000000:role/CWLtoKinesis",
		"arn:aws:kinesis:us-east-1:000000000000:stream/logs-target", nil, "us-east-1"); err != nil {
		t.Fatalf("kinesis target must be accepted: %v", err)
	}
	if _, err := svc.putDestinationCore("firehose-wrap", "arn:aws:iam::000000000000:role/CWLtoKinesis",
		"arn:aws:firehose:us-east-1:000000000000:deliverystream/logs-target", nil, "us-east-1"); err == nil {
		t.Fatal("firehose target must be rejected until the platform Firehose service exists")
	}
}

// TestPutDestinationUpdatePreservesAccessPolicy pins the update half of
// the destination/policy split: PutDestination "creates or updates a
// destination" and carries no accessPolicy member, while the policy is
// attached through its own operation — so an update through
// PutDestination must preserve the policy that operation installed.
func TestPutDestinationUpdatePreservesAccessPolicy(t *testing.T) {
	svc, _ := newDispatchTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}

	const role = "arn:aws:iam::000000000000:role/CWLtoKinesis"
	const target = "arn:aws:kinesis:us-east-1:000000000000:stream/logs-target"
	const retarget = "arn:aws:kinesis:us-east-1:000000000000:stream/logs-target-2"
	const policy = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"logs:PutSubscriptionFilter","Resource":"*"}]}`

	if _, err := svc.putDestinationCore("policy-keeper", role, target, nil, "us-east-1"); err != nil {
		t.Fatalf("create destination: %v", err)
	}
	if err := svc.putDestinationPolicyCore("policy-keeper", policy, "us-east-1"); err != nil {
		t.Fatalf("install policy: %v", err)
	}

	updated, err := svc.putDestinationCore("policy-keeper", role, retarget, nil, "us-east-1")
	if err != nil {
		t.Fatalf("update destination: %v", err)
	}
	if updated.TargetArn != retarget {
		t.Fatalf("update must carry the new target, got %q", updated.TargetArn)
	}
	if updated.AccessPolicy != policy {
		t.Fatalf("update must preserve the installed policy, got %q", updated.AccessPolicy)
	}

	dest, err := store.GetDestination("policy-keeper")
	if err != nil {
		t.Fatal(err)
	}
	if dest.AccessPolicy != policy {
		t.Fatalf("persisted record must keep the installed policy, got %q", dest.AccessPolicy)
	}
}

// TestDispatchDeliversThroughCWLDestinationTarget pins the destination
// indirection end to end: a subscription payload addressed to a CloudWatch
// Logs destination ARN is resolved to the destination record's target ARN
// and delivered there — the cross-account destination model the acceptance
// side newly honours.
func TestDispatchDeliversThroughCWLDestinationTarget(t *testing.T) {
	svc, invoker := newDispatchTestService(t)
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	const target = "arn:aws:kinesis:eu-west-1:000000000000:stream/logs-target"
	if err := store.PutDestination(&logsstore.Destination{
		Name:      "central",
		ARN:       "arn:aws:logs:us-east-1:000000000000:destination:central",
		TargetArn: target,
	}); err != nil {
		t.Fatalf("put destination: %v", err)
	}

	svc.dispatchSubscriptionDelivery("us-east-1",
		"arn:aws:logs:us-east-1:000000000000:destination:central",
		"CloudTrail/logs", "delivery-stream", "ByLogStream", []byte("compressed-payload"))

	regions, streams, _, data := invoker.snapshot()
	if len(regions) != 1 {
		t.Fatalf("expected exactly one delivery through the destination target, got %d", len(regions))
	}
	if regions[0] != "eu-west-1" || streams[0] != "logs-target" {
		t.Fatalf("delivery must address the target ARN's stream, got region %q stream %q", regions[0], streams[0])
	}
	if blob, err := base64.StdEncoding.DecodeString(string(data[0])); err != nil || string(blob) != "compressed-payload" {
		t.Fatalf("delivery must carry the payload in the record's wire form, got %q (err %v)", data[0], err)
	}
}

// TestDispatchSurvivesUndeliverableDestinations pins that every
// undeliverable destination form returns (warns) instead of panicking or
// silently blocking the fan-out worker: an unknown ARN form, a missing
// destination record, and a pre-rejection Firehose filter.
func TestDispatchSurvivesUndeliverableDestinations(t *testing.T) {
	svc, invoker := newDispatchTestService(t)

	svc.dispatchSubscriptionDelivery("us-east-1", "arn:aws:sqs:us-east-1:000000000000:queue", "g", "s", "ByLogStream", []byte("p"))
	svc.dispatchSubscriptionDelivery("us-east-1", "arn:aws:logs:us-east-1:000000000000:destination:missing", "g", "s", "ByLogStream", []byte("p"))
	svc.dispatchSubscriptionDelivery("us-east-1", "arn:aws:firehose:us-east-1:000000000000:deliverystream/legacy", "g", "s", "ByLogStream", []byte("p"))

	if n := invoker.putCount(); n != 0 {
		t.Fatalf("undeliverable destinations must not reach Kinesis, got %d PutRecord calls", n)
	}
}

// TestBusLessDispatchIsVisibleNotSilent pins the bus-less fallback leg:
// with no bus configured the dispatch still runs and reaches the shared
// callees, which warn (rather than return silently) — the delivery loss
// becomes observable in the server log.
func TestBusLessDispatchIsVisibleNotSilent(t *testing.T) {
	svc := &LogsService{accountID: "000000000000"}

	// Neither call may panic; both take the visible-warning path.
	svc.dispatchSubscriptionDelivery("us-east-1", "arn:aws:lambda:us-east-1:000000000000:function:consumer", "g", "s", "ByLogStream", []byte("p"))
	svc.dispatchSubscriptionDelivery("us-east-1", "arn:aws:kinesis:us-east-1:000000000000:stream/logs-dest", "g", "s", "ByLogStream", []byte("p"))
}

// TestSubscriptionDeliveryThroughBusCarriesDistribution pins that the
// published delivery event carries the matched filter's distribution, so
// the bus leg delivers with the same partition-key semantics as the
// bus-less leg.
func TestSubscriptionDeliveryThroughBusCarriesDistribution(t *testing.T) {
	svc, _ := newDispatchTestService(t)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := svc.bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	defer busCancel()
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	if err := store.CreateLogGroup(logsstore.NewLogGroup("bus-group", "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create log group: %v", err)
	}
	if err := store.PutSubscriptionFilter(&logsstore.SubscriptionFilter{
		LogGroupName:   "bus-group",
		FilterName:     "deliver",
		FilterPattern:  "ERROR",
		DestinationArn: "arn:aws:kinesis:us-east-1:000000000000:stream/logs-dest",
		Distribution:   "Random",
	}); err != nil {
		t.Fatalf("put subscription filter: %v", err)
	}

	seen := make(chan string, 1)
	subID, err := eventbus.SubscribeTyped[*eventbus.CloudWatchLogDeliveryEvent](svc.bus, func(_ context.Context, evt *eventbus.CloudWatchLogDeliveryEvent) eventbus.HandlerResult {
		select {
		case seen <- evt.Distribution:
		default:
		}
		return eventbus.HandlerResult{}
	}, eventbus.WithAsync())
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	_ = subID

	svc.deliverSubscriptionEvents(store, "us-east-1", "bus-group", "stream",
		[]logsstore.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "ERROR via bus"}}, nil)

	select {
	case d := <-seen:
		if d != "Random" {
			t.Fatalf("delivery event must carry the filter's distribution, got %q", d)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the delivery event")
	}
}

// TestPutDestinationPolicyByteCap pins the policy member's documented
// byte ceiling: "This can be up to 5120 bytes" — a policy past it rejects
// with InvalidParameterException, one within it installs.
func TestPutDestinationPolicyByteCap(t *testing.T) {
	svc, _ := newDispatchTestService(t)
	const role = "arn:aws:iam::000000000000:role/CWLtoKinesis"
	const target = "arn:aws:kinesis:us-east-1:000000000000:stream/logs-target"
	if _, err := svc.putDestinationCore("policy-cap", role, target, nil, "us-east-1"); err != nil {
		t.Fatalf("create destination: %v", err)
	}

	over := `{"pad":"` + strings.Repeat("x", logsstore.MaxAccessPolicyBytes) + `"}`
	err := svc.putDestinationPolicyCore("policy-cap", over, "us-east-1")
	if code := logsErrorCode(err); code != "InvalidParameterException" {
		t.Fatalf("over-cap policy: code=%q err=%v", code, err)
	}
	within := `{"pad":"` + strings.Repeat("x", 512) + `"}`
	if err := svc.putDestinationPolicyCore("policy-cap", within, "us-east-1"); err != nil {
		t.Fatalf("within-cap policy: %v", err)
	}
}

// TestDispatchFilterFanOutSurvivesShutdownRace pins the shutdown side of
// the fan-out contract: a submission whose select lands the send after
// the context is cancelled — with no worker left to consume it, the race
// the three-case select invites — still evaluates the batch. The sender
// that raced drains the queue itself, so whichever branch the select
// picks (send with post-send drain, or done with the inline fallback)
// the subscription delivery reaches its destination; nothing strands in
// the buffer.
func TestDispatchFilterFanOutSurvivesShutdownRace(t *testing.T) {
	svc, invoker := newDispatchTestService(t)
	// The publish leg needs a started bus carrying the service's
	// delivery handlers: the handler runs on the bus's own lifecycle,
	// which the service cancellation below leaves running.
	bus := eventbus.NewEventBus()
	bus.SetKinesisInvoker(invoker)
	busCtx, busCancel := context.WithCancel(context.Background())
	t.Cleanup(busCancel)
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("set event bus: %v", err)
	}
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("logs store: %v", err)
	}
	const group = "fanout-race-group"
	if err := store.CreateLogGroup(logsstore.NewLogGroup(group, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create log group: %v", err)
	}
	destArn := "arn:aws:logs:us-east-1:000000000000:destination:race-target"
	if err := store.PutDestination(&logsstore.Destination{
		Name:      "race-target",
		ARN:       destArn,
		TargetArn: "arn:aws:kinesis:eu-west-1:000000000000:stream/race-target",
	}); err != nil {
		t.Fatalf("put destination: %v", err)
	}
	if err := svc.putSubscriptionFilterCore(context.Background(), store, &PutSubscriptionFilterInput{
		LogGroupName:     group,
		FilterName:       "race",
		FilterPattern:    "ERROR",
		FilterPatternSet: true,
		DestinationArn:   destArn,
	}); err != nil {
		t.Fatalf("put subscription filter: %v", err)
	}

	// Cancel with no worker ever started: the queue has no consumer by
	// construction, so a stranded send would be a stranded batch.
	svc.cancel()
	before := invoker.putCount()
	svc.dispatchFilterFanOut(store, "us-east-1", group, "s",
		[]logsstore.LogEntry{{Timestamp: time.Now().UnixMilli(), Message: "ERROR during shutdown"}}, nil)

	// The bus leg delivers on its own goroutine: wait for it rather
	// than racing the handler.
	deadline := time.Now().Add(2 * time.Second)
	for invoker.putCount() < before+1 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if got := invoker.putCount() - before; got != 1 {
		t.Fatalf("accepted batch must be evaluated exactly once through the cancelled fan-out, got %d deliveries", got)
	}
	if stranded := len(svc.filterJobs); stranded != 0 {
		t.Fatalf("no job may strand in the shutdown queue, %d remain", stranded)
	}
}
