package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/eventbus"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// Destination delivery fakes: each embeds its invoker interface so only
// the methods the delivery path calls need implementations.

type fakeDestinationBus struct {
	eventbus.Bus
	sqs *fakeSQSDestination
	sns *fakeSNSDestination
	s3  *fakeS3Destination
}

func (b *fakeDestinationBus) SQSInvoker() eventbus.SQSInvoker { return b.sqs }
func (b *fakeDestinationBus) SNSInvoker() eventbus.SNSInvoker { return b.sns }
func (b *fakeDestinationBus) S3Invoker() eventbus.S3Invoker   { return b.s3 }

type fakeSQSDestination struct {
	eventbus.SQSInvoker
	queueURL string
	bodies   []string
}

func (f *fakeSQSDestination) GetQueueByName(ctx context.Context, region, name string) (string, error) {
	if f.queueURL != "" {
		return f.queueURL, nil
	}
	return "", fmt.Errorf("queue %s not found", name)
}

func (f *fakeSQSDestination) SendMessage(ctx context.Context, region, queueURL, body string, opts eventbus.SQSSendOptions) (string, string, error) {
	f.bodies = append(f.bodies, body)
	return "id", "md5", nil
}

type fakeSNSDestination struct {
	eventbus.SNSInvoker
	messages []string
}

func (f *fakeSNSDestination) PublishToTopic(ctx context.Context, topicARN, message, subject string, attributes map[string]string) (string, error) {
	f.messages = append(f.messages, message)
	return "id", nil
}

type fakeS3Destination struct {
	eventbus.S3Invoker
	keys    []string
	bodies  []string
	regions []string
}

func (f *fakeS3Destination) PutObject(ctx context.Context, region, bucket, key string, data []byte, contentType string) error {
	f.keys = append(f.keys, bucket+"/"+key)
	f.bodies = append(f.bodies, string(data))
	f.regions = append(f.regions, region)
	return nil
}

func destinationTestPoller(bus eventbus.Bus) *esmPoller {
	return &esmPoller{bus: bus}
}

func kinesisDestinationItems() []streamBatchItem {
	render := func(seq string, arrival time.Time) map[string]interface{} {
		return map[string]interface{}{
			"kinesis": map[string]interface{}{
				"sequenceNumber":              seq,
				"approximateArrivalTimestamp": arrival.UTC().Format("2006-01-02T15:04:05.000Z"),
			},
		}
	}
	first := time.Unix(1700000000, 0)
	return []streamBatchItem{
		{record: render("100", first), seq: "100"},
		{record: render("101", first.Add(time.Second)), seq: "101"},
		{record: render("102", first.Add(2*time.Second)), seq: "102"},
	}
}

func destinationTestMapping(dest string) *lambdastore.EventSourceMapping {
	return &lambdastore.EventSourceMapping{
		UUID:                 "11111111-2222-3333-4444-555555555555",
		FunctionArn:          "arn:aws:lambda:us-east-1:123456789012:function:esm-fail",
		MaximumRetryAttempts: 1,
		DestinationConfig: &lambdastore.DestinationConfig{
			OnFailure: &lambdastore.OnFailure{Destination: dest},
		},
	}
}

// TestDeliverDiscardedBatch_SQSRecordShape pins the documented SQS/SNS
// invocation record: requestContext (with the only documented condition
// value), responseContext, version, timestamp and the batch placement
// members, without a payload member.
func TestDeliverDiscardedBatch_SQSRecordShape(t *testing.T) {
	sqs := &fakeSQSDestination{queueURL: "http://localhost/queue/dest"}
	bus := &fakeDestinationBus{sqs: sqs, sns: &fakeSNSDestination{}, s3: &fakeS3Destination{}}
	p := destinationTestPoller(bus)

	items := kinesisDestinationItems()
	src := testStreamSource("kinesis", "shardId-000000000001")
	info := streamFailureBatchInfoOf(src, items)
	p.deliverDiscardedBatch(context.Background(), destinationTestMapping("arn:aws:sqs:us-east-1:123456789012:dest-queue"),
		src, info, marshalStreamBatch(items), 2, discardedBatchResponse(&esmFunctionError{classification: "Unhandled"}))

	if len(sqs.bodies) != 1 {
		t.Fatalf("expected 1 SQS delivery, got %d", len(sqs.bodies))
	}
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(sqs.bodies[0]), &record); err != nil {
		t.Fatalf("delivery body is not JSON: %v", err)
	}
	reqCtx, ok := record["requestContext"].(map[string]interface{})
	if !ok {
		t.Fatalf("requestContext missing: %s", sqs.bodies[0])
	}
	if reqCtx["condition"] != "RetryAttemptsExhausted" {
		t.Fatalf("condition = %v, want RetryAttemptsExhausted", reqCtx["condition"])
	}
	if reqCtx["functionArn"] != "arn:aws:lambda:us-east-1:123456789012:function:esm-fail" {
		t.Fatalf("functionArn = %v", reqCtx["functionArn"])
	}
	if reqCtx["approximateInvokeCount"] != float64(2) {
		t.Fatalf("approximateInvokeCount = %v, want 2", reqCtx["approximateInvokeCount"])
	}
	if reqCtx["requestId"] == "" {
		t.Fatal("requestId is empty")
	}
	if record["version"] != "1.0" {
		t.Fatalf("version = %v, want 1.0", record["version"])
	}
	if record["timestamp"] == "" {
		t.Fatal("timestamp is empty")
	}
	respCtx, ok := record["responseContext"].(map[string]interface{})
	if !ok {
		t.Fatalf("responseContext missing: %s", sqs.bodies[0])
	}
	if respCtx["functionError"] != "Unhandled" {
		t.Fatalf("functionError = %v, want Unhandled", respCtx["functionError"])
	}
	batch, ok := record["KinesisBatchInfo"].(map[string]interface{})
	if !ok {
		t.Fatalf("KinesisBatchInfo missing: %s", sqs.bodies[0])
	}
	if batch["shardId"] != "shardId-000000000001" {
		t.Fatalf("shardId = %v", batch["shardId"])
	}
	if batch["startSequenceNumber"] != "100" || batch["endSequenceNumber"] != "102" {
		t.Fatalf("sequence span = %v..%v, want 100..102", batch["startSequenceNumber"], batch["endSequenceNumber"])
	}
	if batch["batchSize"] != float64(3) {
		t.Fatalf("batchSize = %v, want 3", batch["batchSize"])
	}
	if batch["streamArn"] != src.streamArn {
		t.Fatalf("streamArn = %v", batch["streamArn"])
	}
	if batch["approximateArrivalOfFirstRecord"] == "" || batch["approximateArrivalOfLastRecord"] == "" {
		t.Fatal("arrival fields are empty")
	}
	if _, has := record["payload"]; has {
		t.Fatal("SQS destination record must not carry a payload member")
	}
	if _, has := record["DDBStreamBatchInfo"]; has {
		t.Fatal("kinesis record must not carry DDBStreamBatchInfo")
	}
}

// TestDeliverDiscardedBatch_SNSAndS3 pins the SNS delivery and the S3
// difference: the record gains the escaped payload member and lands under
// the documented object naming convention.
func TestDeliverDiscardedBatch_SNSAndS3(t *testing.T) {
	sns := &fakeSNSDestination{}
	s3 := &fakeS3Destination{}
	bus := &fakeDestinationBus{sqs: &fakeSQSDestination{}, sns: sns, s3: s3}
	p := destinationTestPoller(bus)

	items := kinesisDestinationItems()
	src := testStreamSource("kinesis", "shardId-000000000001")
	info := streamFailureBatchInfoOf(src, items)
	payload := marshalStreamBatch(items)

	p.deliverDiscardedBatch(context.Background(),
		destinationTestMapping("arn:aws:sns:us-east-1:123456789012:dest-topic"),
		src, info, payload, 1, uninvokedBatchResponse())
	if len(sns.messages) != 1 {
		t.Fatalf("expected 1 SNS delivery, got %d", len(sns.messages))
	}
	if !strings.Contains(sns.messages[0], "KinesisBatchInfo") {
		t.Fatalf("SNS message lacks KinesisBatchInfo: %s", sns.messages[0])
	}

	p.deliverDiscardedBatch(context.Background(),
		destinationTestMapping("arn:aws:s3:::dest-bucket"),
		src, info, payload, 1, uninvokedBatchResponse())
	if len(s3.bodies) != 1 {
		t.Fatalf("expected 1 S3 delivery, got %d", len(s3.bodies))
	}
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(s3.bodies[0]), &record); err != nil {
		t.Fatalf("S3 body is not JSON: %v", err)
	}
	if record["payload"] == "" {
		t.Fatal("S3 destination record lacks the payload member")
	}
	if !strings.HasPrefix(s3.keys[0], "dest-bucket/aws/lambda/11111111-2222-3333-4444-555555555555/shardId-000000000001/") {
		t.Fatalf("S3 key %q does not follow aws/lambda/<uuid>/<shardID>/ convention", s3.keys[0])
	}
	rest := strings.TrimPrefix(s3.keys[0], "dest-bucket/aws/lambda/11111111-2222-3333-4444-555555555555/shardId-000000000001/")
	parts := strings.SplitN(rest, "/", 4)
	if len(parts) != 4 {
		t.Fatalf("S3 key tail %q lacks the date path", rest)
	}
	datePath := parts[0] + "/" + parts[1] + "/" + parts[2]
	if _, err := time.Parse("2006/01/02", datePath); err != nil {
		t.Fatalf("S3 key date path %q is not YYYY/MM/DD: %v", datePath, err)
	}
	tail := parts[3]
	if len(tail) < 37 {
		t.Fatalf("S3 key tail %q lacks the timestamp and UUID", tail)
	}
	stamp := tail[:len(tail)-37] // drop the "-<36-character UUID>" suffix
	if _, err := time.Parse("2006-01-02T15.04.05", stamp); err != nil {
		t.Fatalf("S3 key timestamp %q is not YYYY-MM-DDTHH.MM.SS: %v", stamp, err)
	}
}

// TestDeliverDiscardedBatch_NoDestination ensures mappings without an
// on-failure destination deliver nothing.
func TestDeliverDiscardedBatch_NoDestination(t *testing.T) {
	sqs := &fakeSQSDestination{queueURL: "u"}
	bus := &fakeDestinationBus{sqs: sqs, sns: &fakeSNSDestination{}, s3: &fakeS3Destination{}}
	p := destinationTestPoller(bus)

	mapping := destinationTestMapping("")
	mapping.DestinationConfig = nil
	items := kinesisDestinationItems()
	src := testStreamSource("kinesis", "shardId-000000000001")
	p.deliverDiscardedBatch(context.Background(), mapping, src, streamFailureBatchInfoOf(src, items),
		marshalStreamBatch(items), 1, uninvokedBatchResponse())
	if len(sqs.bodies) != 0 {
		t.Fatalf("expected no delivery, got %d", len(sqs.bodies))
	}
}

// TestStreamFailureBatchInfoOf_DynamoDB pins the DDB record form: the
// DDBStreamBatchInfo member name and second-precision RFC 3339 arrival
// times.
func TestStreamFailureBatchInfoOf_DynamoDB(t *testing.T) {
	record := &eventbus.DynamoDBStreamRecord{
		Dynamodb: map[string]interface{}{
			"SequenceNumber":              "800000000003126276362",
			"ApproximateCreationDateTime": float64(1573695200),
		},
	}
	items := []streamBatchItem{{record: record, seq: "800000000003126276362"}}
	src := testStreamSource("dynamodb", "shardId-000000000001")
	info := streamFailureBatchInfoOf(src, items)
	if info.StartSequenceNumber != "800000000003126276362" || info.EndSequenceNumber != "800000000003126276362" {
		t.Fatalf("sequence span = %v..%v", info.StartSequenceNumber, info.EndSequenceNumber)
	}
	if info.BatchSize != 1 {
		t.Fatalf("batchSize = %d, want 1", info.BatchSize)
	}
	if !strings.HasSuffix(info.ApproximateArrivalOfFirstRecord, "Z") || len(info.ApproximateArrivalOfFirstRecord) != len("2019-11-14T00:13:19Z") {
		t.Fatalf("first arrival = %q, want second-precision RFC 3339", info.ApproximateArrivalOfFirstRecord)
	}
}

// TestDiscardExpiredDynamoDBRecords pins the age split: expired records
// reach the destination, fresh records survive, and -1 keeps everything.
func TestDiscardExpiredDynamoDBRecords(t *testing.T) {
	sqs := &fakeSQSDestination{queueURL: "u"}
	bus := &fakeDestinationBus{sqs: sqs, sns: &fakeSNSDestination{}, s3: &fakeS3Destination{}}
	p := destinationTestPoller(bus)
	src := testStreamSource("dynamodb", "shardId-000000000001")

	oldRecord := eventbus.DynamoDBStreamRecord{Dynamodb: map[string]interface{}{
		"SequenceNumber":              "1",
		"ApproximateCreationDateTime": float64(time.Now().Add(-2 * time.Hour).Unix()),
	}}
	freshRecord := eventbus.DynamoDBStreamRecord{Dynamodb: map[string]interface{}{
		"SequenceNumber":              "2",
		"ApproximateCreationDateTime": float64(time.Now().Unix()),
	}}

	mapping := destinationTestMapping("arn:aws:sqs:us-east-1:123456789012:dest-queue")
	mapping.MaximumRecordAgeInSeconds = 3600
	fresh := p.discardExpiredDynamoDBRecords(context.Background(), mapping, src,
		[]eventbus.DynamoDBStreamRecord{oldRecord, freshRecord})
	if len(fresh) != 1 || fresh[0].Dynamodb["SequenceNumber"] != "2" {
		t.Fatalf("fresh remainder = %+v, want only sequence 2", fresh)
	}
	if len(sqs.bodies) != 1 || !strings.Contains(sqs.bodies[0], "DDBStreamBatchInfo") {
		t.Fatalf("expected one DDBStreamBatchInfo delivery, got %d bodies", len(sqs.bodies))
	}

	mapping.MaximumRecordAgeInSeconds = -1
	kept := p.discardExpiredDynamoDBRecords(context.Background(), mapping, src,
		[]eventbus.DynamoDBStreamRecord{oldRecord, freshRecord})
	if len(kept) != 2 {
		t.Fatalf("-1 must keep every record, got %d", len(kept))
	}
}

// TestFailureDestinationObjectKey pins the documented naming convention
// layout.
func TestFailureDestinationObjectKey(t *testing.T) {
	now := time.Date(2019, 11, 14, 0, 38, 6, 0, time.UTC)
	key := failureDestinationObjectKey("uuid-1", "shardId-000000000001", now)
	wantPrefix := "aws/lambda/uuid-1/shardId-000000000001/2019/11/14/2019-11-14T00.38.06-"
	if !strings.HasPrefix(key, wantPrefix) {
		t.Fatalf("key %q lacks prefix %q", key, wantPrefix)
	}
	tail := strings.TrimPrefix(key, wantPrefix)
	if _, err := uuid.Parse(tail); err != nil {
		t.Fatalf("key tail %q is not a UUID: %v", tail, err)
	}
}
