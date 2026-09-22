package cloudwatchlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The E1 pins: the vended-logs delivery family's validation rows and the
// delivery engine's execution semantics — the documented five-minute
// cadence pass, the AWSLogs/<account-id>/ S3 prefix, the cursor's
// at-least-once advance, and the CWL destination's flow through the
// shared ingestion seam.

const (
	deliveryTestSourceGroup = "vended-source-group"
	deliveryTestDestGroup   = "vended-dest-group"
	deliveryTestSourceName  = "vended-source"
	deliveryTestDestName    = "vended-dest"
)

// newVendedDeliveryTestEnv builds a service with both fixture groups and
// the recording S3 invoker wired through the event bus.
func newVendedDeliveryTestEnv(t *testing.T) (*LogsService, *logsstore.Store, *recordingS3Invoker) {
	t.Helper()
	svc, store := newReadTestService(t, deliveryTestSourceGroup)
	if err := store.CreateLogGroup(logsstore.NewLogGroup(deliveryTestDestGroup, "us-east-1", "000000000000")); err != nil {
		t.Fatalf("create destination group: %v", err)
	}

	s3 := newRecordingS3Invoker()
	bus := eventbus.NewEventBus()
	bus.SetS3Invoker(s3)
	busCtx, busCancel := context.WithCancel(context.Background())
	if err := bus.Start(busCtx); err != nil {
		t.Fatalf("bus start: %v", err)
	}
	t.Cleanup(busCancel)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatalf("SetEventBus: %v", err)
	}
	return svc, store, s3
}

// deliveryTestGroupARN mints a fixture group's ARN.
func deliveryTestGroupARN(t *testing.T, svc *LogsService, group string) string {
	t.Helper()
	store, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	return store.ARNBuilder().CloudWatch().LogGroup(group)
}

// putDeliverySourceOn builds a delivery source over the given group.
func putDeliverySourceOn(t *testing.T, svc *LogsService, group string) *logsstore.DeliverySource {
	t.Helper()
	view, err := svc.putDeliverySourceCore(&PutDeliverySourceInput{
		Name:        deliveryTestSourceName,
		ResourceArn: deliveryTestGroupARN(t, svc, group),
		LogType:     "APPLICATION_LOGS",
		Region:      "us-east-1",
	})
	if err != nil {
		t.Fatalf("put delivery source: %v", err)
	}
	return view.Source
}

// putCWLTestDestination builds a CWL-typed delivery destination over the
// destination fixture group.
func putCWLTestDestination(t *testing.T, svc *LogsService) *logsstore.DeliveryDestination {
	t.Helper()
	dest, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name:                   deliveryTestDestName,
		DestinationResourceArn: deliveryTestGroupARN(t, svc, deliveryTestDestGroup),
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("put delivery destination: %v", err)
	}
	return dest
}

// ensureTestBucket registers the bucket with the recording invoker so
// Put-time bucket existence checks answer true.
func ensureTestBucket(t *testing.T, s3 *recordingS3Invoker, bucket string) {
	t.Helper()
	if err := s3.EnsureBucket(context.Background(), "us-east-1", bucket); err != nil {
		t.Fatalf("ensure bucket: %v", err)
	}
}

// --- Validation rows (the member-trait census) ---

func TestPutDeliverySourceValidationRows(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	groupArn := store.ARNBuilder().CloudWatch().LogGroup(deliveryTestSourceGroup)

	cases := []struct {
		name    string
		input   *PutDeliverySourceInput
		wantErr string
	}{
		{"missing name", &PutDeliverySourceInput{ResourceArn: groupArn, LogType: "APPLICATION_LOGS", Region: "us-east-1"}, "ValidationException"},
		{"name pattern violation", &PutDeliverySourceInput{Name: "bad name!", ResourceArn: groupArn, LogType: "APPLICATION_LOGS", Region: "us-east-1"}, "ValidationException"},
		{"missing resourceArn", &PutDeliverySourceInput{Name: "src", LogType: "APPLICATION_LOGS", Region: "us-east-1"}, "ValidationException"},
		{"unsupported source service", &PutDeliverySourceInput{
			Name:        "src",
			ResourceArn: "arn:aws:vpc:us-east-1:000000000000:vpc/vpc-0123456789abcdef0",
			LogType:     "FLOW_LOGS", Region: "us-east-1"}, "ValidationException"},
		{"non-ARN resource", &PutDeliverySourceInput{Name: "src", ResourceArn: "not-an-arn", LogType: "APPLICATION_LOGS", Region: "us-east-1"}, "ValidationException"},
		{"missing logType", &PutDeliverySourceInput{Name: "src", ResourceArn: groupArn, Region: "us-east-1"}, "ValidationException"},
		{"logType pattern violation", &PutDeliverySourceInput{Name: "src", ResourceArn: groupArn, LogType: "APPLICATION LOGS", Region: "us-east-1"}, "ValidationException"},
		{"missing source group", &PutDeliverySourceInput{
			Name:        "src",
			ResourceArn: store.ARNBuilder().CloudWatch().LogGroup("no-such-group"),
			LogType:     "APPLICATION_LOGS", Region: "us-east-1"}, "ResourceNotFoundException"},
		{"cross-region resourceArn", &PutDeliverySourceInput{
			Name:        "src",
			ResourceArn: strings.Replace(groupArn, ":us-east-1:", ":us-west-2:", 1),
			LogType:     "APPLICATION_LOGS", Region: "us-east-1"}, "ValidationException"},
	}
	for _, tc := range cases {
		_, err := svc.putDeliverySourceCore(tc.input)
		if code := logsErrorCode(err); code != tc.wantErr {
			t.Fatalf("%s: code=%q want %q (err=%v)", tc.name, code, tc.wantErr, err)
		}
	}
}

func TestPutDeliveryDestinationValidationRows(t *testing.T) {
	svc, store, s3 := newVendedDeliveryTestEnv(t)
	groupArn := store.ARNBuilder().CloudWatch().LogGroup(deliveryTestDestGroup)
	ensureTestBucket(t, s3, "vended-bucket")

	if _, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name: "cwl-ok", DestinationResourceArn: groupArn, Region: "us-east-1",
	}); err != nil {
		t.Fatalf("CWL destination over log group: %v", err)
	}
	if _, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name: "s3-ok", DestinationResourceArn: "arn:aws:s3:::vended-bucket", OutputFormat: "json", Region: "us-east-1",
	}); err != nil {
		t.Fatalf("S3 destination: %v", err)
	}

	// The CWL admission check addresses the ARN's own region — the region
	// the delivery engine writes to — not the call's region: a group that
	// exists only in the ARN's region admits, and one that exists only in
	// the call's region rejects.
	west, err := svc.getLogsStoreByRegion("us-west-2")
	if err != nil {
		t.Fatal(err)
	}
	if err := west.CreateLogGroup(logsstore.NewLogGroup("vended-west-dest", "us-west-2", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name: "cwl-west", DestinationResourceArn: west.ARNBuilder().CloudWatch().LogGroup("vended-west-dest"), Region: "us-east-1",
	}); err != nil {
		t.Fatalf("CWL destination over the ARN region's group: %v", err)
	}
	_, err = svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name: "cwl-west-missing", DestinationResourceArn: west.ARNBuilder().CloudWatch().LogGroup("no-such-west-group"), Region: "us-east-1",
	})
	if code := logsErrorCode(err); code != "ResourceNotFoundException" {
		t.Fatalf("missing ARN-region group: code=%q want ResourceNotFoundException (err=%v)", code, err)
	}

	cases := []struct {
		name    string
		input   *PutDeliveryDestinationInput
		wantErr string
	}{
		{"missing name", &PutDeliveryDestinationInput{DestinationResourceArn: groupArn, Region: "us-east-1"}, "ValidationException"},
		{"name pattern violation", &PutDeliveryDestinationInput{Name: "bad name!", DestinationResourceArn: groupArn, Region: "us-east-1"}, "ValidationException"},
		{"missing configuration", &PutDeliveryDestinationInput{Name: "d", Region: "us-east-1"}, "ValidationException"},
		{"firehose destination", &PutDeliveryDestinationInput{
			Name: "fh", DestinationResourceArn: "arn:aws:firehose:us-east-1:000000000000:deliverystream/vended", Region: "us-east-1"}, "ValidationException"},
		{"xray explicit type", &PutDeliveryDestinationInput{
			Name: "xray", ExplicitDestinationType: "XRAY", DestinationResourceArn: groupArn, Region: "us-east-1"}, "ValidationException"},
		{"parquet output format", &PutDeliveryDestinationInput{
			Name: "pq", DestinationResourceArn: groupArn, OutputFormat: "parquet", Region: "us-east-1"}, "ValidationException"},
		{"unknown output format", &PutDeliveryDestinationInput{
			Name: "fmt", DestinationResourceArn: groupArn, OutputFormat: "avro", Region: "us-east-1"}, "ValidationException"},
		{"type disagrees with resource", &PutDeliveryDestinationInput{
			Name: "mismatch", DestinationResourceArn: groupArn, ExplicitDestinationType: "S3", Region: "us-east-1"}, "ValidationException"},
		{"destination group missing", &PutDeliveryDestinationInput{
			Name: "gone", DestinationResourceArn: store.ARNBuilder().CloudWatch().LogGroup("no-such-group"), Region: "us-east-1"}, "ResourceNotFoundException"},
		{"destination bucket missing", &PutDeliveryDestinationInput{
			Name: "gone-bucket", DestinationResourceArn: "arn:aws:s3:::no-such-bucket", Region: "us-east-1"}, "ResourceNotFoundException"},
	}
	for _, tc := range cases {
		_, err := svc.putDeliveryDestinationCore(tc.input)
		if code := logsErrorCode(err); code != tc.wantErr {
			t.Fatalf("%s: code=%q want %q (err=%v)", tc.name, code, tc.wantErr, err)
		}
	}
}

// The record-shaping members the create and update operations share.
func TestDeliveryShapingValidationRows(t *testing.T) {
	cases := []struct {
		name          string
		recordFields  []string
		delimiter     string
		suffixPath    string
		destType      string
		wantRejection bool
	}{
		{"vocabulary fields", []string{"time", "stream", "message"}, "", "", "S3", false},
		{"listed delimiter", []string{"time", "message"}, ",", "", "S3", false},
		{"suffix on S3", []string{"time"}, "", "inbox/{yyyy}", "S3", false},
		{"unknown record field", []string{"time", "bucket"}, "", "", "S3", true},
		{"empty record field", []string{"time", ""}, "", "", "S3", true},
		{"oversize delimiter", []string{"time"}, "abcdef", "", "S3", true},
		{"unlisted delimiter", []string{"time"}, "##", "", "S3", true},
		{"suffix on CWL", []string{"time"}, "", "inbox/{yyyy}", "CWL", true},
	}
	for _, tc := range cases {
		err := validateDeliveryShaping(tc.recordFields, tc.delimiter, tc.suffixPath, tc.destType, tc.suffixPath != "")
		if tc.wantRejection && err == nil {
			t.Fatalf("%s: accepted", tc.name)
		}
		if !tc.wantRejection && err != nil {
			t.Fatalf("%s: rejected: %v", tc.name, err)
		}
		if tc.wantRejection && logsErrorCode(err) != "ValidationException" {
			t.Fatalf("%s: code=%q (err=%v)", tc.name, logsErrorCode(err), err)
		}
	}

	var many []string
	for i := 0; i < logsstore.DeliveryRecordFieldsMax+1; i++ {
		many = append(many, "time")
	}
	if err := validateDeliveryShaping(many, "", "", "S3", false); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("recordFields over the ceiling: %v", err)
	}
}

func TestCreateDeliveryPairingRows(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	if _, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	}); err != nil {
		t.Fatalf("create delivery: %v", err)
	}

	dup, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	})
	if err == nil || logsErrorCode(err) != "ConflictException" {
		t.Fatalf("duplicate pairing: %v", err)
	}
	if dup != nil {
		t.Fatal("duplicate pairing returned a delivery")
	}

	_, err = svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     "no-such-source",
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	})
	if logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("missing source: %v", err)
	}
	_, err = svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName: source.Name,
		DeliveryDestinationArn: store.ARNBuilder().CloudWatch().
			DeliveryDestination("no-such-dest"),
		Region: "us-east-1",
	})
	if logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("missing destination: %v", err)
	}
}

func TestDeliveryLifecycleAndGuards(t *testing.T) {
	svc, _, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		RecordFields:           []string{"time", "message"},
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("create delivery: %v", err)
	}
	if svcarn.ExtractDeliveryIdFromARN(delivery.Arn) != delivery.Id {
		t.Fatalf("delivery ARN does not address the id: %q vs %q", delivery.Arn, delivery.Id)
	}

	// The delete guards: a referenced source or destination refuses.
	if err := svc.deleteDeliverySourceCore("us-east-1", source.Name); logsErrorCode(err) != "ConflictException" {
		t.Fatalf("delete referenced source: %v", err)
	}
	if err := svc.deleteDeliveryDestinationCore("us-east-1", dest.Name); logsErrorCode(err) != "ConflictException" {
		t.Fatalf("delete referenced destination: %v", err)
	}
	if err := svc.deleteDeliverySourceCore("us-east-1", "never-was"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("delete missing source: %v", err)
	}
	if err := svc.deleteDeliveryCore("us-east-1", "no-such-delivery"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("delete missing delivery: %v", err)
	}

	// The update re-shapes the delivery without changing its endpoints.
	if err := svc.updateDeliveryConfigurationCore(&UpdateDeliveryConfigurationInput{
		Id: delivery.Id, RecordFields: []string{"message"}, Region: "us-east-1",
	}); err != nil {
		t.Fatalf("update delivery configuration: %v", err)
	}
	saved, err := svc.getLogsStoreByRegion("us-east-1")
	if err != nil {
		t.Fatal(err)
	}
	updated, err := saved.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.RecordFields) != 1 || updated.RecordFields[0] != "message" {
		t.Fatalf("updated record fields: %v", updated.RecordFields)
	}

	if err := svc.deleteDeliveryCore("us-east-1", delivery.Id); err != nil {
		t.Fatalf("delete delivery: %v", err)
	}
	if err := svc.deleteDeliverySourceCore("us-east-1", source.Name); err != nil {
		t.Fatalf("delete unreferenced source: %v", err)
	}
	if err := svc.deleteDeliveryDestinationCore("us-east-1", dest.Name); err != nil {
		t.Fatalf("delete unreferenced destination: %v", err)
	}
}

func TestDeliveryDestinationPolicyRoundtrip(t *testing.T) {
	svc, _, _ := newVendedDeliveryTestEnv(t)
	putCWLTestDestination(t, svc)

	policyDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"logs:CreateDelivery","Resource":"*"}]}`
	if _, err := svc.putDeliveryDestinationPolicyCore("us-east-1", deliveryTestDestName, policyDoc); err != nil {
		t.Fatalf("put policy: %v", err)
	}
	got, err := svc.getDeliveryDestinationPolicyCore("us-east-1", deliveryTestDestName)
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	if got.PolicyDocument != policyDoc {
		t.Fatal("policy round-trip mismatch")
	}
	// The family's declared error list carries ValidationException, not
	// the log-plane InvalidParameterException.
	if _, err := svc.putDeliveryDestinationPolicyCore("us-east-1", deliveryTestDestName, "not-json"); logsErrorCode(err) != "ValidationException" {
		t.Fatalf("non-JSON policy: %v", err)
	}
	if _, err := svc.getDeliveryDestinationPolicyCore("us-east-1", "never-was"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("policy of missing destination: %v", err)
	}
	if err := svc.deleteDeliveryDestinationPolicyCore("us-east-1", deliveryTestDestName); err != nil {
		t.Fatalf("delete policy: %v", err)
	}
	if _, err := svc.getDeliveryDestinationPolicyCore("us-east-1", deliveryTestDestName); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("get deleted policy: %v", err)
	}
	if err := svc.deleteDeliveryDestinationPolicyCore("us-east-1", deliveryTestDestName); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("delete deleted policy: %v", err)
	}
}

// The source's status members are computed from the underlying
// resource: ACTIVE while the group exists, INACTIVE with statusReason
// RESOURCE_DELETED after its deletion.
func TestDeliverySourceStatusTracksResource(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)

	view := svc.deliverySourceStatus("us-east-1", source)
	formatted := formatDeliverySource(view)
	if formatted["status"] != "ACTIVE" {
		t.Fatalf("status with live group: %v", formatted["status"])
	}
	if _, ok := formatted["statusReason"]; ok {
		t.Fatal("ACTIVE source carries statusReason")
	}

	if err := store.DeleteLogGroup(deliveryTestSourceGroup); err != nil {
		t.Fatal(err)
	}
	view = svc.deliverySourceStatus("us-east-1", source)
	formatted = formatDeliverySource(view)
	if formatted["status"] != "INACTIVE" || formatted["statusReason"] != "RESOURCE_DELETED" {
		t.Fatalf("status after group deletion: %v / %v", formatted["status"], formatted["statusReason"])
	}
}

// --- Engine execution ---

// The CWL-typed delivery flows through the shared ingestion seam: the
// source group's new events land in the destination group with their
// originating streams, once per cadence pass, and the persisted cursor
// stops the next pass re-delivering them.
func TestDeliveryEngineCWLFlow(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	const stream = "app/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now - 2000, Message: `{"level":"INFO","msg":"first"}`},
		{Timestamp: now - 1000, Message: `{"level":"ERROR","msg":"second"}`},
	}); err != nil {
		t.Fatal(err)
	}

	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		RecordFields:           []string{"time", "message"},
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	svc.tickDeliveries()

	events, err := fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("destination group events: got %d, want 2", len(events))
	}
	// The delivered form is shaped per the delivery's recordFields —
	// the default (unset outputFormat) renders the fields joined by the
	// default delimiter, so the raw message no longer re-ingests bare.
	firstShaped := fmt.Sprintf("%d\t%s", now-2000, `{"level":"INFO","msg":"first"}`)
	if events[0].Message != firstShaped {
		t.Fatalf("first delivered message: %q want %q", events[0].Message, firstShaped)
	}

	saved, err := store.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CursorTime != now-1000 || saved.CursorStream != stream {
		t.Fatalf("cursor after pass: %d/%q", saved.CursorTime, saved.CursorStream)
	}

	// A later event delivers on the next pass; the earlier two do not
	// re-deliver.
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now - 500, Message: "third"},
	}); err != nil {
		t.Fatal(err)
	}
	svc.tickDeliveries()
	events, err = fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("destination group events after second pass: got %d, want 3", len(events))
	}
	thirdShaped := fmt.Sprintf("%d\t%s", now-500, "third")
	if events[2].Message != thirdShaped {
		t.Fatalf("third delivered message: %q want %q", events[2].Message, thirdShaped)
	}
}

// Each stream's window is cut into batches the shared ingestion seam
// accepts — the count, byte and span bounds of the PutLogEvents batch
// contract — so a window past any one bound still delivers instead of
// rejecting outright on every pass.
func TestSplitDeliveryBatchesCutsOnContractBounds(t *testing.T) {
	batchSizes := func(batches [][]logsstore.LogEntry) []int {
		sizes := make([]int, len(batches))
		for i, b := range batches {
			sizes[i] = len(b)
		}
		return sizes
	}

	// The count cut: one event past the maximal batch splits the run.
	countRun := make([]logsstore.LogEntry, logsstore.MaxChunkSize+3)
	for i := range countRun {
		countRun[i] = logsstore.LogEntry{Timestamp: 5000, Message: "m"}
	}
	batches := splitDeliveryBatches(countRun)
	if len(batches) != 2 || len(batches[0]) != logsstore.MaxChunkSize || len(batches[1]) != 3 {
		t.Fatalf("count cut: batch sizes %v", batchSizes(batches))
	}

	// The byte cut: a run whose batch bytes exceed the ceiling splits so
	// every batch stays within it.
	const msgBytes = 1000
	n := logsstore.MaxPutLogEventsBatchBytes/(msgBytes+logsstore.PutLogEventOverheadBytes) + 5
	byteRun := make([]logsstore.LogEntry, n)
	for i := range byteRun {
		byteRun[i] = logsstore.LogEntry{Timestamp: 5000, Message: strings.Repeat("x", msgBytes)}
	}
	batches = splitDeliveryBatches(byteRun)
	if len(batches) < 2 {
		t.Fatalf("byte cut did not split: batch sizes %v", batchSizes(batches))
	}
	for bi, b := range batches {
		total := 0
		for _, e := range b {
			total += len(e.Message) + logsstore.PutLogEventOverheadBytes
		}
		if total > logsstore.MaxPutLogEventsBatchBytes {
			t.Fatalf("batch %d carries %d bytes over the ceiling", bi, total)
		}
	}

	// The span cut: a run spanning past 24 hours splits at the boundary.
	hour := int64(60 * 60 * 1000)
	spanRun := []logsstore.LogEntry{
		{Timestamp: 0, Message: "a"},
		{Timestamp: 10 * hour, Message: "b"},
		{Timestamp: 20 * hour, Message: "c"},
		{Timestamp: 30 * hour, Message: "d"},
		{Timestamp: 40 * hour, Message: "e"},
	}
	batches = splitDeliveryBatches(spanRun)
	if len(batches) != 2 || len(batches[0]) != 3 || len(batches[1]) != 2 {
		t.Fatalf("span cut: batch sizes %v", batchSizes(batches))
	}
}

// A CWL-destination window larger than one PutLogEvents batch delivers
// completely: without the cut, the single oversized ingestion rejects on
// every pass, the cursor never advances, and the delivery stalls until
// the events age out of every horizon.
func TestDeliveryEngineCWLWindowBeyondOneBatch(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	const stream = "bulk/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UnixMilli()
	// Two source batches — a maximal one and one more event — so the
	// pending window holds MaxChunkSize+1 events for the stream.
	big := make([]logsstore.LogEntry, logsstore.MaxChunkSize)
	for i := range big {
		big[i] = logsstore.LogEntry{Timestamp: ts, Message: "m"}
	}
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, big); err != nil {
		t.Fatal(err)
	}
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: ts, Message: "tail"},
	}); err != nil {
		t.Fatal(err)
	}

	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("create delivery: %v", err)
	}

	svc.tickDeliveries()

	events, err := fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != logsstore.MaxChunkSize+1 {
		t.Fatalf("destination group events: got %d, want %d", len(events), logsstore.MaxChunkSize+1)
	}
	saved, err := store.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CursorTime != ts || saved.CursorStream != stream {
		t.Fatalf("cursor after bulk pass: %d/%q", saved.CursorTime, saved.CursorStream)
	}
}

// Events the destination's retention horizon permanently rejects are
// dropped with a warning while the rest of the window delivers and the
// cursor advances: their age only grows, so holding them would park the
// delivery behind events no pass could ever deliver.
func TestDeliveryEngineDropsDestinationExpiredEvents(t *testing.T) {
	svc, store := newReadTestService(t, deliveryTestSourceGroup)
	destLg := logsstore.NewLogGroup(deliveryTestDestGroup, "us-east-1", "000000000000")
	destLg.RetentionInDays = 1
	if err := store.CreateLogGroup(destLg); err != nil {
		t.Fatalf("create retention-bound destination group: %v", err)
	}
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	const stream = "expiry/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	twoDays := int64(2 * 24 * 60 * 60 * 1000)
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now - twoDays, Message: "aged-out"},
		{Timestamp: now - 1000, Message: "fresh"},
	}); err != nil {
		t.Fatal(err)
	}

	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("create delivery: %v", err)
	}

	svc.tickDeliveries()

	events, err := fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	// The delivery carries the default recordFields, so the fresh event
	// arrives shaped (time, stream and message joined by the default
	// delimiter); the aged-out event never arrives.
	freshShaped := fmt.Sprintf("%d\t%s\tfresh", now-1000, stream)
	if len(events) != 1 || events[0].Message != freshShaped {
		t.Fatalf("destination events: got %v, want the fresh event alone shaped as %q", events, freshShaped)
	}
	saved, err := store.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CursorTime != now-1000 || saved.CursorStream != stream {
		t.Fatalf("cursor must advance past the dropped event: %d/%q", saved.CursorTime, saved.CursorStream)
	}
}

// The S3-typed delivery writes one gzipped object per pass under the
// documented AWSLogs/<account-id>/ prefix, shaped per the destination's
// output format and the delivery's record fields.
func TestDeliveryEngineS3Flow(t *testing.T) {
	svc, store, s3 := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)

	const bucket = "vended-delivery-bucket"
	ensureTestBucket(t, s3, bucket)
	dest, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name:                   "s3-dest",
		DestinationResourceArn: "arn:aws:s3:::" + bucket,
		OutputFormat:           "json",
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	const stream = "ingest/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now - 1000, Message: "line-one"},
		{Timestamp: now - 500, Message: "line-two"},
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		RecordFields:           []string{"time", "message"},
		Region:                 "us-east-1",
	}); err != nil {
		t.Fatal(err)
	}

	svc.tickDeliveries()

	objects := deliveryKeysFor(t, s3, bucket)
	if len(objects) != 1 {
		t.Fatalf("delivered objects: got %d, want 1 (%v)", len(objects), objects)
	}
	if !strings.HasPrefix(objects[0], "AWSLogs/000000000000/vended-source/") {
		t.Fatalf("delivered key prefix: %q", objects[0])
	}

	body := gunzipObject(t, s3, bucket, objects[0])
	lines := strings.Split(strings.TrimRight(body, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("delivered lines: %q", body)
	}
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &record); err != nil {
		t.Fatalf("delivered line is not JSON: %v", err)
	}
	if record["message"] != "line-one" {
		t.Fatalf("delivered record message: %v", record)
	}
	if _, ok := record["stream"]; ok {
		t.Fatalf("delivered record carries unselected stream field: %v", record)
	}
	if _, ok := record["time"]; !ok {
		t.Fatalf("delivered record missing time: %v", record)
	}

	// The second pass writes nothing new: the cursor advanced past both
	// events.
	svc.tickDeliveries()
	if got := len(deliveryKeysFor(t, s3, bucket)); got != 1 {
		t.Fatalf("delivered objects after second pass: got %d, want 1", got)
	}
}

// The suffixPath variables substitute into the delivered object's path,
// and the w3c output renders its #Fields header with the delimiter.
func TestDeliveryEngineS3PathForms(t *testing.T) {
	svc, store, s3 := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)

	const bucket = "vended-suffix-bucket"
	ensureTestBucket(t, s3, bucket)
	dest, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name:                   "s3-suffix-dest",
		DestinationResourceArn: "arn:aws:s3:::" + bucket,
		OutputFormat:           "w3c",
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	const stream = "sfx/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: time.Now().UnixMilli() - 1000, Message: "w3c-line"},
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		FieldDelimiter:         "\t",
		S3SuffixPath:           "inbox/{accountId}/{yyyy}",
		Region:                 "us-east-1",
	}); err != nil {
		t.Fatal(err)
	}
	svc.tickDeliveries()

	objects := deliveryKeysFor(t, s3, bucket)
	if len(objects) != 1 {
		t.Fatalf("delivered objects: %v", objects)
	}
	year := time.Now().UTC().Format("2006")
	wantPrefix := "AWSLogs/000000000000/vended-source/inbox/000000000000/" + year + "/"
	if !strings.HasPrefix(objects[0], wantPrefix) {
		t.Fatalf("suffix-path key: %q want prefix %q", objects[0], wantPrefix)
	}
	body := gunzipObject(t, s3, bucket, objects[0])
	if !strings.HasPrefix(body, "#Fields: time stream message\n") {
		t.Fatalf("w3c header: %q", body)
	}
	if !strings.HasSuffix(body, "w3c-line\n") {
		t.Fatalf("w3c body: %q", body)
	}
}

// A delivery whose destination is the source group itself never
// re-ingests: the self-delivery guard skips the pass.
func TestDeliveryEngineSelfDeliveryGuard(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)

	dest, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name:                   "self-dest",
		DestinationResourceArn: deliveryTestGroupARN(t, svc, deliveryTestSourceGroup),
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	const stream = "self/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: time.Now().UnixMilli(), Message: "loop"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	}); err != nil {
		t.Fatal(err)
	}

	svc.tickDeliveries()

	events, err := fetchAllLogEvents(store, deliveryTestSourceGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("self-delivery duplicated events: got %d, want 1", len(events))
	}
}

// The tags map traits hold on every Put operation of the family: 1-50
// entries ("Map Entries: Maximum number of 50 items", PutDeliveryDestination
// API reference) whose keys and values satisfy the TagKey and TagValue
// traits — raised as the family's ValidationException, the identity the
// three operations declare.
func TestDeliveryPutOperationsTagTraits(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	groupArn := store.ARNBuilder().CloudWatch().LogGroup(deliveryTestSourceGroup)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	oversize := make(map[string]string, tagutil.MaxTagsPerResource+1)
	for i := 0; i <= tagutil.MaxTagsPerResource; i++ {
		oversize[fmt.Sprintf("key%02d", i)] = "value"
	}

	cases := []struct {
		name    string
		tags    map[string]string
		wantErr string
	}{
		{"empty map", map[string]string{}, "ValidationException"},
		{"over the 50 ceiling", oversize, "ValidationException"},
		{"key charset violation", map[string]string{"bad key!": "v"}, "ValidationException"},
		{"key over 128 characters", map[string]string{strings.Repeat("k", tagutil.MaxTagKeyLength+1): "v"}, "ValidationException"},
		{"value over 256 characters", map[string]string{"k": strings.Repeat("v", tagutil.MaxTagValueLength+1)}, "ValidationException"},
		{"reserved aws prefix", map[string]string{"aws:team": "v"}, "ValidationException"},
	}
	for _, tc := range cases {
		if _, err := svc.putDeliverySourceCore(&PutDeliverySourceInput{
			Name: "tag-src", ResourceArn: groupArn, LogType: "APPLICATION_LOGS",
			Tags: tc.tags, Region: "us-east-1",
		}); logsErrorCode(err) != tc.wantErr {
			t.Fatalf("%s / source: code=%q want %q (err=%v)", tc.name, logsErrorCode(err), tc.wantErr, err)
		}
		if _, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
			Name: "tag-dest", DestinationResourceArn: groupArn,
			Tags: tc.tags, Region: "us-east-1",
		}); logsErrorCode(err) != tc.wantErr {
			t.Fatalf("%s / destination: code=%q want %q (err=%v)", tc.name, logsErrorCode(err), tc.wantErr, err)
		}
		if _, err := svc.createDeliveryCore(&CreateDeliveryInput{
			DeliverySourceName: source.Name, DeliveryDestinationArn: dest.Arn,
			Tags: tc.tags, Region: "us-east-1",
		}); logsErrorCode(err) != tc.wantErr {
			t.Fatalf("%s / delivery: code=%q want %q (err=%v)", tc.name, logsErrorCode(err), tc.wantErr, err)
		}
	}

	// A tag set inside every trait persists on the family's records.
	valid := map[string]string{"team": "logs", "cost-centre": "vended"}
	tagged, err := svc.putDeliverySourceCore(&PutDeliverySourceInput{
		Name: "tag-src", ResourceArn: groupArn, LogType: "APPLICATION_LOGS",
		Tags: valid, Region: "us-east-1",
	})
	if err != nil {
		t.Fatalf("valid tags rejected: %v", err)
	}
	if !reflect.DeepEqual(tagged.Source.Tags, valid) {
		t.Fatalf("source tags round-trip: %v", tagged.Source.Tags)
	}
}

// An S3 delivery destination admits only against the wired S3 substrate:
// with no invoker the bucket's existence cannot be verified, and admitting
// it would accept a delivery that could never flow — so the unwired state
// rejects with the operation's declared ServiceUnavailableException.
func TestPutDeliveryDestinationS3RequiresWiredSubstrate(t *testing.T) {
	svc, _ := newReadTestService(t, deliveryTestSourceGroup)
	_, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name: "s3-unwired", DestinationResourceArn: "arn:aws:s3:::vended-bucket", Region: "us-east-1",
	})
	if code := logsErrorCode(err); code != "ServiceUnavailableException" {
		t.Fatalf("unwired S3 substrate: code=%q want ServiceUnavailableException (err=%v)", code, err)
	}
}

// A deleted destination group fails the pass instead of being silently
// recreated with default settings: the destination is the customer's
// existing group ("Create a delivery destination for an existing log
// group", vended-logs setup) and the delivery flow holds stream creation
// and event writes on it alone. The cursor does not advance over the
// failed pass, and a customer recreation of the group resumes the
// delivery.
func TestDeliveryEngineDestinationGroupDeletionFailsPass(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	const stream = "recreate/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now, Message: "gone-dest"},
	}); err != nil {
		t.Fatal(err)
	}

	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("create delivery: %v", err)
	}

	if err := store.DeleteLogGroup(deliveryTestDestGroup); err != nil {
		t.Fatal(err)
	}

	svc.tickDeliveries()

	if _, err := store.GetLogGroup(deliveryTestDestGroup); err == nil {
		t.Fatal("deleted destination group was recreated by the delivery pass")
	}
	saved, err := store.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CursorTime != 0 {
		t.Fatalf("cursor advanced over a failed pass: %d", saved.CursorTime)
	}

	if err := store.CreateLogGroup(logsstore.NewLogGroup(deliveryTestDestGroup, "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	svc.tickDeliveries()
	events, err := fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events after destination group recreation: got %d, want 1", len(events))
	}
}

// DescribeConfigurationTemplates serves the platform catalog and honours
// its filters.
func TestDescribeConfigurationTemplatesCatalog(t *testing.T) {
	svc, _, _ := newVendedDeliveryTestEnv(t)
	resp, err := svc.DescribeConfigurationTemplates(context.Background(), nil, vocabRequest(nil))
	if err != nil {
		t.Fatalf("describe templates: %v", err)
	}
	list := resp.(map[string]interface{})["configurationTemplates"]
	if list == nil {
		list = []map[string]interface{}{}
	}
	rows, ok := list.([]map[string]interface{})
	if !ok || len(rows) != 2 {
		t.Fatalf("template catalog rows: %v", list)
	}
	types := map[string]bool{}
	for _, row := range rows {
		types[row["deliveryDestinationType"].(string)] = true
		if row["service"] != "logs" || row["resourceType"] != "logGroup" {
			t.Fatalf("template row identity: %v", row)
		}
	}
	if !types["CWL"] || !types["S3"] {
		t.Fatalf("template destination types: %v", types)
	}

	resp, err = svc.DescribeConfigurationTemplates(context.Background(), nil, vocabRequest(map[string]interface{}{
		"deliveryDestinationTypes": []interface{}{"S3"},
	}))
	if err != nil {
		t.Fatal(err)
	}
	rows = resp.(map[string]interface{})["configurationTemplates"].([]map[string]interface{})
	if len(rows) != 1 || rows[0]["deliveryDestinationType"] != "S3" {
		t.Fatalf("filtered catalog rows: %v", rows)
	}

	// An unknown service filter yields an empty catalog, not an error.
	resp, err = svc.DescribeConfigurationTemplates(context.Background(), nil, vocabRequest(map[string]interface{}{
		"service": "ec2",
	}))
	if err != nil {
		t.Fatalf("unknown service filter: %v", err)
	}
	rowsIface := resp.(map[string]interface{})["configurationTemplates"]
	if rowsIface != nil && len(rowsIface.([]map[string]interface{})) != 0 {
		t.Fatalf("unknown service filter rows: %v", rowsIface)
	}
}

// --- helpers ---

// deliveryKeysFor lists the object keys written to one bucket.
func deliveryKeysFor(t *testing.T, s3 *recordingS3Invoker, bucket string) []string {
	t.Helper()
	keys, err := s3.ListObjects(context.Background(), "us-east-1", bucket, "", 1000)
	if err != nil {
		t.Fatalf("list objects: %v", err)
	}
	return keys
}

// gunzipObject reads and decompresses one recorded object.
func gunzipObject(t *testing.T, s3 *recordingS3Invoker, bucket, key string) string {
	t.Helper()
	data, ok := s3.object(bucket, key)
	if !ok {
		t.Fatalf("object %s/%s not recorded", bucket, key)
	}
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer zr.Close()
	body, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("gzip read: %v", err)
	}
	return string(body)
}

// Byte-identical events at the same millisecond are distinct store
// records and each delivers once: the cursor triple alone would cut
// every later duplicate forever, so the cursor also counts the events
// sharing its triple.
func TestDeliveryCursorDuplicateEvents(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	source := putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	const stream = "dup/1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, deliveryTestSourceGroup)); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now, Message: "same"},
		{Timestamp: now, Message: "same"},
	}); err != nil {
		t.Fatal(err)
	}

	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     source.Name,
		DeliveryDestinationArn: dest.Arn,
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}

	svc.tickDeliveries()
	events, err := fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("identical events delivered: got %d, want 2", len(events))
	}
	saved, err := store.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CursorCount != 2 {
		t.Fatalf("cursor count = %d, want 2", saved.CursorCount)
	}

	// The second pass re-delivers nothing; a third identical event
	// appended later delivers on its pass.
	svc.tickDeliveries()
	events, _ = fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if len(events) != 2 {
		t.Fatalf("second pass re-delivered: got %d events", len(events))
	}
	if _, err := store.PutLogEvents(deliveryTestSourceGroup, stream, []logsstore.LogEntry{
		{Timestamp: now, Message: "same"},
	}); err != nil {
		t.Fatal(err)
	}
	svc.tickDeliveries()
	events, _ = fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if len(events) != 3 {
		t.Fatalf("later duplicate never delivered: got %d events", len(events))
	}

	// The third append extended the cursor triple: the count accumulates
	// the previously delivered events at it, and a fourth pass delivers
	// nothing. A replacement count would skip too few on every later
	// pass and re-deliver the triple's earlier events forever.
	saved, err = store.GetDelivery(delivery.Id)
	if err != nil {
		t.Fatal(err)
	}
	if saved.CursorCount != 3 {
		t.Fatalf("cursor count after triple extension = %d, want 3", saved.CursorCount)
	}
	svc.tickDeliveries()
	events, _ = fetchAllLogEvents(store, deliveryTestDestGroup, stream, 0, 0)
	if len(events) != 3 {
		t.Fatalf("fourth pass re-delivered after triple extension: got %d events", len(events))
	}
}

// The S3 object key's documented forms: the Hive-compatible account
// segment AWSLogs/aws-account-id=<id>/, Hive rendering of suffixPath
// variables (myFolder/{yyyy}/{MM}/{dd} yielding
// myFolder/year=.../month=.../day=.../), and the destination prefix
// replacing the default AWSLogs/<account>/<service>/ head.
func TestDeliveryS3HiveAndPrefixKeys(t *testing.T) {
	const account, source = "111122223333", "app-source"
	now := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	base := &logsstore.Delivery{}
	s3Dest := &logsstore.DeliveryDestination{}

	plain := deliveryS3ObjectKey(account, source, base, s3Dest, now, "")
	if !strings.HasPrefix(plain, "AWSLogs/"+account+"/"+source+"/2026/09/10/") {
		t.Fatalf("plain key: %q", plain)
	}

	hive := deliveryS3ObjectKey(account, source, &logsstore.Delivery{S3HiveCompatiblePath: true}, s3Dest, now, "")
	if !strings.HasPrefix(hive, "AWSLogs/aws-account-id="+account+"/"+source+"/year=2026/month=09/day=10/") {
		t.Fatalf("hive key: %q", hive)
	}

	hiveSuffix := deliveryS3ObjectKey(account, source, &logsstore.Delivery{
		S3HiveCompatiblePath: true,
		S3SuffixPath:         "myFolder/{yyyy}/{MM}/{dd}",
	}, s3Dest, now, "")
	if !strings.HasPrefix(hiveSuffix, "AWSLogs/aws-account-id="+account+"/"+source+"/myFolder/year=2026/month=09/day=10/") {
		t.Fatalf("hive suffix key: %q", hiveSuffix)
	}

	plainSuffix := deliveryS3ObjectKey(account, source, &logsstore.Delivery{
		S3SuffixPath: "myFolder/{yyyy}",
	}, s3Dest, now, "")
	if !strings.HasPrefix(plainSuffix, "AWSLogs/"+account+"/"+source+"/myFolder/2026/") {
		t.Fatalf("plain suffix key: %q", plainSuffix)
	}

	prefixed := deliveryS3ObjectKey(account, source, base, s3Dest, now, "MyLogPrefix")
	if !strings.HasPrefix(prefixed, "MyLogPrefix/2026/09/10/") {
		t.Fatalf("destination-prefix key: %q", prefixed)
	}
	if strings.Contains(prefixed, "AWSLogs/") {
		t.Fatalf("destination prefix did not replace the default head: %q", prefixed)
	}

	bucket, prefix, ok := s3BucketAndPrefixFromArn("arn:aws:s3:::bucket-name/MyLogPrefix")
	if !ok || bucket != "bucket-name" || prefix != "MyLogPrefix" {
		t.Fatalf("bucket/prefix split: %q %q %v", bucket, prefix, ok)
	}
	bucket, prefix, ok = s3BucketAndPrefixFromArn("arn:aws:s3:::bucket-name")
	if !ok || bucket != "bucket-name" || prefix != "" {
		t.Fatalf("bare bucket split: %q %q %v", bucket, prefix, ok)
	}
}

// "Both keys and values must be between 1 and 255 characters in
// length." — deliverySourceConfiguration members validate their length.
func TestPutDeliverySourceConfigTraits(t *testing.T) {
	svc, _, _ := newVendedDeliveryTestEnv(t)
	long := strings.Repeat("v", 256)
	for _, row := range []struct {
		name   string
		config map[string]string
	}{
		{"empty key", map[string]string{"": "value"}},
		{"overlong key", map[string]string{long: "value"}},
		{"empty value", map[string]string{"k": ""}},
		{"overlong value", map[string]string{"k": long}},
	} {
		if _, err := svc.putDeliverySourceCore(&PutDeliverySourceInput{
			Name:                        "cfg-" + strings.ReplaceAll(strings.ReplaceAll(row.name, " ", "-"), "overlong", "ol"),
			ResourceArn:                 "arn:aws:logs:us-east-1:000000000000:log-group:" + deliveryTestSourceGroup,
			LogType:                     "ACCESS",
			DeliverySourceConfiguration: row.config,
			Region:                      "us-east-1",
		}); logsErrorCode(err) != "ValidationException" {
			t.Errorf("%s: %v", row.name, err)
		}
	}
	// The boundary serves.
	if _, err := svc.putDeliverySourceCore(&PutDeliverySourceInput{
		Name:                        "cfg-boundary",
		ResourceArn:                 "arn:aws:logs:us-east-1:000000000000:log-group:" + deliveryTestSourceGroup,
		LogType:                     "ACCESS",
		DeliverySourceConfiguration: map[string]string{strings.Repeat("k", 255): strings.Repeat("v", 255)},
		Region:                      "us-east-1",
	}); err != nil {
		t.Fatalf("boundary config: %v", err)
	}
}

// The destination record is a resource of the call's region: a CWL
// destination whose resource ARN addresses another region's log group
// admits against that region's group but persists — record and ARN both
// — in the store the call resolved, where the family's reads find it.
func TestPutDeliveryDestinationPersistsInCallRegion(t *testing.T) {
	svc, store, _ := newVendedDeliveryTestEnv(t)
	_ = store
	westStore, err := svc.getLogsStoreByRegion("us-west-2")
	if err != nil {
		t.Fatalf("west store: %v", err)
	}
	if err := westStore.CreateLogGroup(logsstore.NewLogGroup("west-dest-group", "us-west-2", "000000000000")); err != nil {
		t.Fatalf("west group: %v", err)
	}
	westArn := westStore.ARNBuilder().CloudWatch().LogGroup("west-dest-group")

	dest, err := svc.putDeliveryDestinationCore(&PutDeliveryDestinationInput{
		Name:                   "cross-region-dest",
		DestinationResourceArn: westArn,
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("put cross-region destination: %v", err)
	}
	if got, err := store.GetDeliveryDestination("cross-region-dest"); err != nil || got == nil {
		t.Fatalf("call-region store must hold the destination: %v", err)
	}
	if _, err := westStore.GetDeliveryDestination("cross-region-dest"); err == nil {
		t.Fatal("destination resource's region must not hold the record")
	}
	if !strings.HasPrefix(dest.Arn, "arn:aws:logs:us-east-1:") {
		t.Fatalf("destination ARN must carry the call region: %s", dest.Arn)
	}
}

// The delivery record carries the destination's platform ARN, whatever
// spelling the creation request cast: the delete-time reference guard
// and the (source, destination) pair uniqueness compare against the
// platform-cast ARN, so a caller-cast ARN with foreign account digits
// must not slip a second delivery past the duplicate gate or a referenced
// destination past the delete guard.
func TestCreateDeliveryCanonicalisesDestinationArn(t *testing.T) {
	svc, _, _ := newVendedDeliveryTestEnv(t)
	putDeliverySourceOn(t, svc, deliveryTestSourceGroup)
	dest := putCWLTestDestination(t, svc)

	doctored := strings.Replace(dest.Arn, ":000000000000:", ":123456789012:", 1)
	delivery, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     deliveryTestSourceName,
		DeliveryDestinationArn: doctored,
		Region:                 "us-east-1",
	})
	if err != nil {
		t.Fatalf("create with caller-cast ARN: %v", err)
	}
	if delivery.DeliveryDestinationArn != dest.Arn {
		t.Fatalf("stored destination ARN must be the platform's own: got %s want %s",
			delivery.DeliveryDestinationArn, dest.Arn)
	}

	// The pair gate holds under a differently-cast repeat of the same
	// destination.
	if _, err := svc.createDeliveryCore(&CreateDeliveryInput{
		DeliverySourceName:     deliveryTestSourceName,
		DeliveryDestinationArn: strings.Replace(dest.Arn, ":000000000000:", ":999999999999:", 1),
		Region:                 "us-east-1",
	}); logsErrorCode(err) != "ConflictException" {
		t.Fatalf("differently-cast duplicate pair must conflict: %v", err)
	}

	// The reference guard holds: the destination is in use and refuses
	// to delete.
	if err := svc.deleteDeliveryDestinationCore("us-east-1", deliveryTestDestName); err == nil {
		t.Fatal("referenced destination must refuse to delete")
	}
}
