package kinesis

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	kmsutil "vorpalstacks/internal/common/kmsutil"
	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// The pins below cover the member-trait validation matrix: every constraint
// the Core functions enforce from the model's member traits and documented
// cross-member rules is pinned with its error identity, and the
// boundary-adjacent valid values are pinned as accepted so the negative
// edge is exact.

// requireAWSCode fails the test unless err carries the given AWS error code.
func requireAWSCode(t *testing.T, name string, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected %s, got success", name, code)
	}
	apiErr, ok := err.(*awserrors.AWSError)
	if !ok {
		t.Fatalf("%s: error %v is not an AWSError", name, err)
	}
	if apiErr.GetCode() != code {
		t.Fatalf("%s: code = %q, want %q", name, apiErr.GetCode(), code)
	}
}

// TestValidationMatrixRejections pins every rejection row of the matrix
// against a live store: each row names the member trait or documented rule
// it enforces and the error identity the model declares for it.
func TestValidationMatrixRejections(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	ctx := context.Background()

	if _, err := store.CreateStream("matrix_prov", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create provisioned fixture: %v", err)
	}
	onDemand, err := store.CreateStream("matrix_od", 1, kinesisstore.StreamModeOnDemand, 0, 0, nil)
	if err != nil {
		t.Fatalf("create on-demand fixture: %v", err)
	}

	cases := []struct {
		name   string
		invoke func() error
		code   string
	}{
		{"StopStreamEncryption without EncryptionType", func() error {
			_, err := svc.stopStreamEncryptionCore(reqCtx, StopStreamEncryptionInput{
				StreamName: "matrix_prov", KeyId: "alias/key",
			})
			return err
		}, "InvalidArgumentException"},
		{"StopStreamEncryption with non-KMS EncryptionType", func() error {
			_, err := svc.stopStreamEncryptionCore(reqCtx, StopStreamEncryptionInput{
				StreamName: "matrix_prov", EncryptionType: "NONE", KeyId: "alias/key",
			})
			return err
		}, "InvalidArgumentException"},
		{"StopStreamEncryption without KeyId", func() error {
			_, err := svc.stopStreamEncryptionCore(reqCtx, StopStreamEncryptionInput{
				StreamName: "matrix_prov", EncryptionType: "KMS",
			})
			return err
		}, "InvalidArgumentException"},

		{"UpdateStreamMode with absent StreamModeDetails", func() error {
			err := svc.updateStreamModeCore(reqCtx, UpdateStreamModeInput{StreamARN: onDemand.StreamARN})
			return err
		}, "InvalidArgumentException"},
		{"UpdateStreamMode with details that lack the inner mode", func() error {
			_, err := svc.UpdateStreamMode(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamARN":         onDemand.StreamARN,
				"StreamModeDetails": map[string]interface{}{"ignored": "member"},
			}))
			return err
		}, "InvalidArgumentException"},

		{"EnableEnhancedMonitoring with an empty metrics list", func() error {
			_, err := svc.enableEnhancedMonitoringCore(reqCtx, EnhancedMonitoringInput{
				StreamName: "matrix_prov",
			})
			return err
		}, "InvalidArgumentException"},
		{"EnableEnhancedMonitoring with eight metrics", func() error {
			_, err := svc.enableEnhancedMonitoringCore(reqCtx, EnhancedMonitoringInput{
				StreamName:        "matrix_prov",
				ShardLevelMetrics: strings.Fields("IncomingBytes IncomingRecords OutgoingBytes OutgoingRecords WriteProvisionedThroughputExceeded ReadProvisionedThroughputExceeded IteratorAgeMilliseconds ALL"),
			})
			return err
		}, "InvalidArgumentException"},

		{"GetShardIterator AT_SEQUENCE_NUMBER without StartingSequenceNumber", func() error {
			_, err := svc.getShardIteratorCore(reqCtx, GetShardIteratorInput{
				StreamName: "matrix_prov", ShardId: "shardId-000000000000",
				ShardIteratorType: "AT_SEQUENCE_NUMBER",
			})
			return err
		}, "InvalidArgumentException"},
		{"GetShardIterator AFTER_SEQUENCE_NUMBER without StartingSequenceNumber", func() error {
			_, err := svc.getShardIteratorCore(reqCtx, GetShardIteratorInput{
				StreamName: "matrix_prov", ShardId: "shardId-000000000000",
				ShardIteratorType: "AFTER_SEQUENCE_NUMBER",
			})
			return err
		}, "InvalidArgumentException"},
		{"GetShardIterator with a non-pattern StartingSequenceNumber", func() error {
			_, err := svc.getShardIteratorCore(reqCtx, GetShardIteratorInput{
				StreamName: "matrix_prov", ShardId: "shardId-000000000000",
				ShardIteratorType: "AT_SEQUENCE_NUMBER", StartingSequenceNumber: "0123",
			})
			return err
		}, "InvalidArgumentException"},
		{"GetShardIterator AT_TIMESTAMP without Timestamp", func() error {
			_, err := svc.getShardIteratorCore(reqCtx, GetShardIteratorInput{
				StreamName: "matrix_prov", ShardId: "shardId-000000000000",
				ShardIteratorType: "AT_TIMESTAMP",
			})
			return err
		}, "InvalidArgumentException"},

		{"SplitShard without NewStartingHashKey", func() error {
			_, err := svc.splitShardCore(reqCtx, SplitShardInput{
				StreamName: "matrix_prov", ShardToSplit: "shardId-000000000000",
			})
			return err
		}, "InvalidArgumentException"},
		{"SplitShard at the parent range's starting boundary", func() error {
			_, err := svc.splitShardCore(reqCtx, SplitShardInput{
				StreamName: "matrix_prov", ShardToSplit: "shardId-000000000000",
				NewStartingHashKey: "0",
			})
			return err
		}, "InvalidArgumentException"},
		{"SplitShard on an on-demand stream", func() error {
			_, err := svc.splitShardCore(reqCtx, SplitShardInput{
				StreamName: "matrix_od", ShardToSplit: "shardId-000000000000",
				NewStartingHashKey: "9223372036854775808",
			})
			return err
		}, "InvalidArgumentException"},
		{"MergeShards on an on-demand stream", func() error {
			_, err := svc.mergeShardsCore(reqCtx, MergeShardsInput{
				StreamName:           "matrix_od",
				ShardToMerge:         "shardId-000000000000",
				AdjacentShardToMerge: "shardId-000000000001",
			})
			return err
		}, "InvalidArgumentException"},

		{"UpdateShardCount without ScalingType", func() error {
			_, err := svc.updateShardCountCore(reqCtx, UpdateShardCountInput{
				StreamName: "matrix_prov", TargetShardCount: 2,
			})
			return err
		}, "InvalidArgumentException"},
		{"UpdateShardCount on an on-demand stream", func() error {
			_, err := svc.updateShardCountCore(reqCtx, UpdateShardCountInput{
				StreamName: "matrix_od", TargetShardCount: 2, ScalingType: "UNIFORM_SCALING",
			})
			return err
		}, "InvalidArgumentException"},

		{"UpdateStreamWarmThroughput with a negative value", func() error {
			_, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
				StreamName: "matrix_od", WarmThroughputMiBps: -1,
			})
			return err
		}, "InvalidArgumentException"},
		{"UpdateStreamWarmThroughput on a provisioned stream", func() error {
			_, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
				StreamName: "matrix_prov", WarmThroughputMiBps: 128,
			})
			return err
		}, "InvalidArgumentException"},

		{"UpdateAccountSettings without the commitment member", func() error {
			_, err := svc.updateAccountSettingsCore(reqCtx, UpdateAccountSettingsInput{})
			return err
		}, "InvalidArgumentException"},
		{"UpdateAccountSettings with a non-enum status", func() error {
			_, err := svc.updateAccountSettingsCore(reqCtx, UpdateAccountSettingsInput{Status: "PAUSED"})
			return err
		}, "InvalidArgumentException"},

		{"IncreaseStreamRetentionPeriod to the current value", func() error {
			_, _, err := svc.updateRetentionPeriodCore(reqCtx, UpdateRetentionPeriodInput{
				StreamName: "matrix_prov", RetentionPeriodHours: 24,
			}, true)
			return err
		}, "InvalidArgumentException"},
		{"DecreaseStreamRetentionPeriod to the current value", func() error {
			_, _, err := svc.updateRetentionPeriodCore(reqCtx, UpdateRetentionPeriodInput{
				StreamName: "matrix_prov", RetentionPeriodHours: 24,
			}, false)
			return err
		}, "InvalidArgumentException"},

		{"ListShards without stream identification", func() error {
			_, err := svc.listShardsCore(reqCtx, ListShardsInput{})
			return err
		}, "InvalidArgumentException"},

		{"PutRecord with a non-pattern SequenceNumberForOrdering", func() error {
			_, err := svc.putRecordCore(reqCtx, PutRecordInput{
				StreamName: "matrix_prov", Data: "data", HasData: true, PartitionKey: "pk",
				SequenceNumberForOrdering: "00",
			})
			return err
		}, "InvalidArgumentException"},

		{"CreateStream with a non-enum stream mode", func() error {
			_, err := svc.createStreamCore(reqCtx, CreateStreamInput{
				StreamName: "matrix_badmode", StreamMode: "TURBO", HasStreamModeDetails: true,
			})
			return err
		}, "InvalidArgumentException"},

		{"CreateStream with an explicit zero shard count", func() error {
			_, err := svc.createStreamCore(reqCtx, CreateStreamInput{
				StreamName: "matrix_zerocount", ShardCount: 0, HasShardCount: true,
			})
			return err
		}, "InvalidArgumentException"},

		{"PutRecords entry without its required Data", func() error {
			_, err := svc.putRecordsCore(reqCtx, PutRecordsInput{
				StreamName: "matrix_prov",
				Records: []interface{}{
					map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
					map[string]interface{}{"PartitionKey": "pk"},
				},
			})
			return err
		}, "InvalidArgumentException"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			requireAWSCode(t, tc.name, tc.invoke(), tc.code)
		})
	}
}

// TestValidationMatrixBoundaryAcceptances pins the boundary-adjacent valid
// values so the rejection edges above are exact, together with the
// sequence-number pattern window, the warm-throughput range floor and the
// split range window on a narrowed shard.
func TestValidationMatrixBoundaryAcceptances(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	ctx := context.Background()

	if _, err := store.CreateStream("matrix_edge", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create provisioned fixture: %v", err)
	}
	if _, err := store.CreateStream("matrix_edge_od", 1, kinesisstore.StreamModeOnDemand, 0, 0, nil); err != nil {
		t.Fatalf("create on-demand fixture: %v", err)
	}
	// The on-demand creation default: four shards regardless of the
	// ShardCount member — the throughput-equivalent of the documented 4
	// MB/s write default.
	odShards, err := store.ListShards("matrix_edge_od", nil, "", 0)
	if err != nil {
		t.Fatalf("list on-demand fixture shards: %v", err)
	}
	if len(odShards) != int(kinesisstore.DefaultOnDemandShardCount) {
		t.Fatalf("on-demand creation shards: got %d, want the four-shard default", len(odShards))
	}

	// The SequenceNumber pattern: "0", a 129-digit non-zero-leading decimal
	// and unset are the accepted forms; a 130-digit value and a leading zero
	// reject.
	if !validateSequenceNumber("") || !validateSequenceNumber("0") {
		t.Fatal("validateSequenceNumber rejected an accepted form (unset or zero)")
	}
	if !validateSequenceNumber("1" + strings.Repeat("0", 128)) {
		t.Fatal("validateSequenceNumber rejected the 129-digit boundary value")
	}
	if validateSequenceNumber("1" + strings.Repeat("0", 129)) {
		t.Fatal("validateSequenceNumber accepted a 130-digit value")
	}
	if validateSequenceNumber("00") {
		t.Fatal("validateSequenceNumber accepted a leading zero")
	}

	// The warm-throughput range floor: the NaturalIntegerObject range trait
	// carries min 0 alone — zero releases the warm capacity and is valid.
	if !validateWarmThroughputMiBps(0) || validateWarmThroughputMiBps(-1) {
		t.Fatal("validateWarmThroughputMiBps window is wrong at the zero boundary")
	}

	// The shard-count window at its edges.
	if !validateShardCount(1) || !validateShardCount(kinesisstore.MaxShardCount) {
		t.Fatal("validateShardCount window is wrong at its inclusive edges")
	}
	if validateShardCount(kinesisstore.MaxShardCount+1) || validateShardCount(0) {
		t.Fatal("validateShardCount accepted an out-of-window value")
	}

	// Seven metrics is the MetricsNameList length ceiling.
	seven := strings.Fields("IncomingBytes IncomingRecords OutgoingBytes OutgoingRecords WriteProvisionedThroughputExceeded ReadProvisionedThroughputExceeded IteratorAgeMilliseconds")
	if !validateShardLevelMetricsList(seven) {
		t.Fatal("validateShardLevelMetricsList rejected the seven-metric ceiling")
	}

	// GetShardIterator with a pattern-valid sequence number creates the
	// iterator; the sequence-zero position reads an empty page.
	if _, err := svc.getShardIteratorCore(reqCtx, GetShardIteratorInput{
		StreamName: "matrix_edge", ShardId: "shardId-000000000000",
		ShardIteratorType: "AT_SEQUENCE_NUMBER", StartingSequenceNumber: "0",
	}); err != nil {
		t.Fatalf("get shard iterator at sequence zero: %v, want accepted", err)
	}

	// SplitShard accepts a key strictly inside the parent range, and the
	// two children tile the parent exactly. (The store's unfiltered list
	// carries closed shards as well; the open set is the split's result.)
	if _, err := svc.splitShardCore(reqCtx, SplitShardInput{
		StreamName: "matrix_edge", ShardToSplit: "shardId-000000000000",
		NewStartingHashKey: "9223372036854775808",
	}); err != nil {
		t.Fatalf("split inside the parent range: %v, want accepted", err)
	}
	allShards, err := store.ListShards("matrix_edge", nil, "", 0)
	if err != nil {
		t.Fatalf("list shards after split: %v", err)
	}
	var shards []*kinesisstore.Shard
	for _, shard := range allShards {
		if shard.SequenceNumberRange == nil || shard.SequenceNumberRange.EndingSequenceNumber == "" {
			shards = append(shards, shard)
		}
	}
	if len(shards) != 2 {
		t.Fatalf("shards after split: got %d open shards, want 2", len(shards))
	}
	low, high := shards[0], shards[1]
	if low.HashKeyRange.StartingHashKey != "0" || low.HashKeyRange.EndingHashKey != "9223372036854775807" {
		t.Fatalf("lower child range: got [%s, %s]", low.HashKeyRange.StartingHashKey, low.HashKeyRange.EndingHashKey)
	}
	if high.HashKeyRange.StartingHashKey != "9223372036854775808" || high.HashKeyRange.EndingHashKey != kinesisstore.MaxShardHashKey {
		t.Fatalf("upper child range: got [%s, %s]", high.HashKeyRange.StartingHashKey, high.HashKeyRange.EndingHashKey)
	}

	// A narrowed shard pins the range window itself: a key above the shard's
	// end but inside the platform space rejects, as does a key below its
	// start, and the ending boundary itself accepts (the upper child becomes
	// a single-point range).
	requireAWSCode(t, "split above the narrowed range", func() error {
		_, err := svc.splitShardCore(reqCtx, SplitShardInput{
			StreamName: "matrix_edge", ShardToSplit: low.ShardID,
			NewStartingHashKey: "9223372036854775900",
		})
		return err
	}(), "InvalidArgumentException")
	requireAWSCode(t, "split below the narrowed range", func() error {
		_, err := svc.splitShardCore(reqCtx, SplitShardInput{
			StreamName: "matrix_edge", ShardToSplit: high.ShardID,
			NewStartingHashKey: "9223372036854775807",
		})
		return err
	}(), "InvalidArgumentException")
	if _, err := svc.splitShardCore(reqCtx, SplitShardInput{
		StreamName: "matrix_edge", ShardToSplit: low.ShardID,
		NewStartingHashKey: "9223372036854775807",
	}); err != nil {
		t.Fatalf("split at the narrowed range's ending boundary: %v, want accepted", err)
	}

	// UpdateStreamWarmThroughput resolves by StreamName alone, accepts the
	// zero floor on an on-demand stream and leaves the mode configuration
	// untouched.
	result, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
		StreamName: "matrix_edge_od", WarmThroughputMiBps: 0, HasWarmThroughputMiBps: true,
	})
	if err != nil {
		t.Fatalf("warm throughput by name at the zero floor: %v, want accepted", err)
	}
	if result.WarmThroughputMiBps != 0 {
		t.Fatalf("warm throughput result: got %d, want 0", result.WarmThroughputMiBps)
	}
	updated, err := store.GetStream("matrix_edge_od")
	if err != nil {
		t.Fatalf("get stream after warm update: %v", err)
	}
	if updated.StreamModeDetails == nil || updated.StreamModeDetails.StreamMode != kinesisstore.StreamModeOnDemand {
		t.Fatalf("stream mode after warm update: got %+v, want ON_DEMAND unchanged", updated.StreamModeDetails)
	}

	// The warm-throughput report channel: a configured figure travels in
	// the summary, and the zero release floor reports as the absence of a
	// configured figure — the update response echoes the accepted target
	// (pinned above at the result value), so the release stays observable
	// without the summary reporting a zero capacity.
	summaryOf := func() map[string]interface{} {
		resp, err := svc.DescribeStreamSummary(ctx, reqCtx, parsedRequest(map[string]interface{}{
			"StreamName": "matrix_edge_od",
		}))
		if err != nil {
			t.Fatalf("describe stream summary: %v", err)
		}
		body, ok := resp.(map[string]interface{})
		if !ok {
			t.Fatalf("summary response shape: %T", resp)
		}
		summary, ok := body["StreamDescriptionSummary"].(map[string]interface{})
		if !ok {
			t.Fatalf("summary description shape: %T", body["StreamDescriptionSummary"])
		}
		return summary
	}
	if _, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
		StreamName: "matrix_edge_od", WarmThroughputMiBps: 128, HasWarmThroughputMiBps: true,
	}); err != nil {
		t.Fatalf("configure warm throughput: %v", err)
	}
	warm, ok := summaryOf()["WarmThroughput"].(map[string]interface{})
	if !ok || warm["TargetMiBps"] != int32(128) {
		t.Fatalf("summary with a configured figure: got %+v, want WarmThroughput TargetMiBps 128", summaryOf()["WarmThroughput"])
	}
	if _, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
		StreamName: "matrix_edge_od", WarmThroughputMiBps: 0, HasWarmThroughputMiBps: true,
	}); err != nil {
		t.Fatalf("release warm throughput: %v", err)
	}
	if _, still := summaryOf()["WarmThroughput"]; still {
		t.Fatal("summary after the zero release still carries WarmThroughput")
	}

	// The commitment status accepts both enum values.
	for _, status := range []string{"ENABLED", "DISABLED"} {
		if _, err := svc.updateAccountSettingsCore(reqCtx, UpdateAccountSettingsInput{Status: status}); err != nil {
			t.Fatalf("account settings status %q: %v, want accepted", status, err)
		}
	}

	// The ResourceARN shape accepts channel ARNs alongside stream ARNs —
	// the post-2026-09-07 model widened the pattern.
	if !validateResourceARN("arn:aws:kinesis:us-east-1:000000000000:channel/some-channel") {
		t.Fatal("validateResourceARN rejected a channel ARN the model's pattern accepts")
	}
	if !validateResourceARN("arn:aws:kinesis:us-east-1:000000000000:stream/some-stream") {
		t.Fatal("validateResourceARN rejected a stream ARN")
	}

	// A directional retention change still applies.
	if _, _, err := svc.updateRetentionPeriodCore(reqCtx, UpdateRetentionPeriodInput{
		StreamName: "matrix_edge", RetentionPeriodHours: 25,
	}, true); err != nil {
		t.Fatalf("increase retention to 25: %v, want accepted", err)
	}

	// A pattern-valid SequenceNumberForOrdering writes.
	if _, err := svc.putRecordCore(reqCtx, PutRecordInput{
		StreamName: "matrix_edge", Data: "data", HasData: true, PartitionKey: "pk",
		SequenceNumberForOrdering: "4961406234703423420635712463",
	}); err != nil {
		t.Fatalf("put record with a sequence number for ordering: %v, want accepted", err)
	}

	// The required Data member: an absent one is a shape violation, not an
	// empty record to ingest.
	if _, err := svc.putRecordCore(reqCtx, PutRecordInput{
		StreamName: "matrix_edge", PartitionKey: "pk",
	}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("put record without Data: %v, want InvalidArgumentException", err)
	}

	// A non-string Data member is a wire-type violation — the blob's wire
	// form is a base64 string — rejected at the strict read before the
	// Core, never coerced to an empty payload the lenient readers produce.
	if _, err := svc.PutRecord(ctx, reqCtx, &request.ParsedRequest{
		Operation: "PutRecord",
		Parameters: map[string]interface{}{
			"StreamName":   "matrix_edge",
			"PartitionKey": "pk",
			"Data":         123,
		},
	}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("put record with a non-string Data: %v, want InvalidArgumentException", err)
	}

	// A null Data member reads as absent (the protocol drops null structure
	// values), so the required-member check rejects the request.
	if _, err := svc.PutRecord(ctx, reqCtx, &request.ParsedRequest{
		Operation: "PutRecord",
		Parameters: map[string]interface{}{
			"StreamName":   "matrix_edge",
			"PartitionKey": "pk",
			"Data":         nil,
		},
	}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("put record with a null Data: %v, want InvalidArgumentException", err)
	}

	// An absent StreamModeDetails still defaults the creation to
	// provisioned, and a present one without the inner member rejects.
	if _, err := svc.createStreamCore(reqCtx, CreateStreamInput{StreamName: "matrix_edge_default", ShardCount: 1}); err != nil {
		t.Fatalf("create without mode details: %v, want the provisioned default", err)
	}
	if _, err := svc.createStreamCore(reqCtx, CreateStreamInput{StreamName: "matrix_edge_innerless", HasStreamModeDetails: true}); err == nil {
		t.Fatal("create with mode details lacking the inner member: accepted, want rejection")
	}
}

// TestPutRecordsWholeRequestCap pins the documented whole-request bound:
// the sum of every entry's decoded payload and partition key must fit the
// ten-MiB request cap — each entry fits its own per-record bound, so only
// the aggregate can stop a batch of valid records.
func TestPutRecordsWholeRequestCap(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("reqcap", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}

	// One MiB minus the partition key's share keeps every entry inside the
	// per-record bound, so ten entries sit under the request cap and eleven
	// cross it.
	payload := base64.StdEncoding.EncodeToString(make([]byte, 1024*1024-len("pk")))
	build := func(n int) []interface{} {
		entries := make([]interface{}, n)
		for i := range entries {
			entries[i] = map[string]interface{}{"Data": payload, "PartitionKey": "pk"}
		}
		return entries
	}

	if _, err := svc.putRecordsCore(reqCtx, PutRecordsInput{StreamName: "reqcap", Records: build(10)}); err != nil {
		t.Fatalf("batch under the request cap: %v, want accepted", err)
	}
	_, err := svc.putRecordsCore(reqCtx, PutRecordsInput{StreamName: "reqcap", Records: build(11)})
	requireAWSCode(t, "batch past the request cap", err, "InvalidArgumentException")
}

// TestValidationMatrixErrorIdentity pins that the trio's status guard
// answers the model-declared ResourceInUseException — the guard is
// unreachable through the public surface while every stream this platform
// creates is ACTIVE, so the identity is pinned at the helper itself.
func TestValidationMatrixErrorIdentity(t *testing.T) {
	err := ensureReshapableStream(&kinesisstore.Stream{StreamStatus: kinesisstore.StreamStatusUpdating})
	requireAWSCode(t, "non-ACTIVE guard", err, "ResourceInUseException")
}

// TestRecordSizeCountsPartitionKey pins the Data member's combined-size
// rule: the decoded payload added to the partition key size must fit the
// stream's maximum record size — a payload that fits alone fails through
// its key, identically on the single-record path and inside a batch (an
// entry-shape violation rejects the whole request; no earlier entry is
// written).
func TestRecordSizeCountsPartitionKey(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)

	// The minimum record-size window keeps the fixture payload small.
	if _, err := store.CreateStream("matrix_size", 1, kinesisstore.StreamModeProvisioned, kinesisstore.MinMaxRecordSizeInKiB, 0, nil); err != nil {
		t.Fatalf("create size fixture: %v", err)
	}

	fitsAlone := base64.StdEncoding.EncodeToString(make([]byte, kinesisstore.MinMaxRecordSizeInKiB*1024-100))
	failsThroughKey := strings.Repeat("k", 200)

	if _, err := svc.putRecordCore(reqCtx, PutRecordInput{
		StreamName: "matrix_size", Data: fitsAlone, HasData: true, PartitionKey: "pk",
	}); err != nil {
		t.Fatalf("payload that fits with its key: %v, want accepted", err)
	}

	_, err := svc.putRecordCore(reqCtx, PutRecordInput{
		StreamName: "matrix_size", Data: fitsAlone, PartitionKey: failsThroughKey,
	})
	requireAWSCode(t, "single record whose key pushes it over", err, "InvalidArgumentException")

	_, err = svc.putRecordsCore(reqCtx, PutRecordsInput{
		StreamName: "matrix_size",
		Records: []interface{}{
			map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
			map[string]interface{}{"Data": fitsAlone, "PartitionKey": failsThroughKey},
		},
	})
	requireAWSCode(t, "batch carrying one oversized entry", err, "InvalidArgumentException")

	// The request-level rejection wrote nothing: the first, valid entry
	// must not have landed when its successor failed validation.
	records, _, err := store.GetRecords("matrix_size", "shardId-000000000000", "", 10, false, time.Time{})
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records after rejected batch: got %d, want the single pre-batch record alone", len(records))
	}
}

// TestWarmThroughputRequiredAndShardIdTraits pins two contract rows the
// matrix had missed: WarmThroughputMiBps is a required member of
// UpdateStreamWarmThroughput (an absent member is rejected, not read as the
// range-legal zero that would overwrite the stored target with nothing),
// and every shard-id member carries the ShardId shape's pattern — a
// malformed ID answers invalid input on all four operations that take one
// instead of falling into a not-found lookup.
func TestWarmThroughputRequiredAndShardIdTraits(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("warmreq", 1, kinesisstore.StreamModeOnDemand, 0, 0, nil)
	if err != nil {
		t.Fatalf("create on-demand fixture: %v", err)
	}

	if _, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
		StreamName: "warmreq", HasWarmThroughputMiBps: false,
	}); err == nil {
		t.Fatal("absent WarmThroughputMiBps accepted, want the required member rejected")
	} else {
		requireAWSCode(t, "absent WarmThroughputMiBps", err, "InvalidArgumentException")
	}

	if _, err := svc.updateStreamWarmThroughputCore(reqCtx, UpdateStreamWarmThroughputInput{
		StreamName: "warmreq", WarmThroughputMiBps: 0, HasWarmThroughputMiBps: true,
	}); err != nil {
		t.Fatalf("explicit zero warm throughput rejected: %v, want the range-legal value accepted", err)
	}

	if _, err := svc.getShardIteratorCore(reqCtx, GetShardIteratorInput{
		StreamName: "warmreq", ShardId: "shard id/with spaces", ShardIteratorType: "LATEST",
	}); err == nil {
		t.Fatal("malformed GetShardIterator ShardId accepted")
	} else {
		requireAWSCode(t, "GetShardIterator malformed ShardId", err, "InvalidArgumentException")
	}

	if _, err := svc.splitShardCore(reqCtx, SplitShardInput{
		StreamName: "warmreq", ShardToSplit: "no/such;id", NewStartingHashKey: "1",
	}); err == nil {
		t.Fatal("malformed SplitShard ShardToSplit accepted")
	} else {
		requireAWSCode(t, "SplitShard malformed ShardToSplit", err, "InvalidArgumentException")
	}

	if _, err := svc.mergeShardsCore(reqCtx, MergeShardsInput{
		StreamName: "warmreq", ShardToMerge: "shardId-000000000000", AdjacentShardToMerge: "bad id",
	}); err == nil {
		t.Fatal("malformed MergeShards AdjacentShardToMerge accepted")
	} else {
		requireAWSCode(t, "MergeShards malformed AdjacentShardToMerge", err, "InvalidArgumentException")
	}

	consumer, err := store.RegisterStreamConsumer(stream.StreamARN, "shardid-reader", nil)
	if err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	if _, err := svc.subscribeToShardCore(context.Background(), reqCtx, SubscribeToShardInput{
		ConsumerARN: consumer.ConsumerARN, ShardId: "bad id", StartingPositionType: "LATEST",
	}); err == nil {
		t.Fatal("malformed SubscribeToShard ShardId accepted")
	} else {
		requireAWSCode(t, "SubscribeToShard malformed ShardId", err, "InvalidArgumentException")
	}
}

// fakeChecker pins the encryption core's key-check seam without a KMS
// service behind it.
type fakeChecker struct{ err error }

func (f fakeChecker) CheckKey(ctx context.Context, region, keyID string) error { return f.err }

// TestStartStreamEncryptionChecksKey pins the reachable KMS family: the
// key resolves through the injected checker before the stream adopts it,
// each sentinel failure answers its modelled identity, and a healthy key
// proceeds to encryption.
func TestStartStreamEncryptionChecksKey(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("kmschk", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	const keyARN = "arn:aws:kms:us-east-1:000000000000:key/12345678-1234-1234-1234-123456789012"

	for _, tc := range []struct {
		name string
		err  error
		code string
	}{
		{"missing key", kmsutil.ErrKeyNotFound, "KMSNotFoundException"},
		{"disabled key", kmsutil.ErrKeyDisabled, "KMSDisabledException"},
		{"invalid-state key", kmsutil.ErrKeyInvalidState, "KMSInvalidStateException"},
		{"invalid-usage key", kmsutil.ErrKeyInvalidUsage, "KMSAccessDeniedException"},
	} {
		svc.SetKMSChecker(fakeChecker{tc.err})
		_, err := svc.startStreamEncryptionCore(reqCtx, StartStreamEncryptionInput{
			StreamName: "kmschk", EncryptionType: "KMS", KeyId: keyARN,
		})
		requireAWSCode(t, tc.name, err, tc.code)
	}

	svc.SetKMSChecker(fakeChecker{nil})
	stream, err := svc.startStreamEncryptionCore(reqCtx, StartStreamEncryptionInput{
		StreamName: "kmschk", EncryptionType: "KMS", KeyId: keyARN,
	})
	if err != nil {
		t.Fatalf("healthy key: %v", err)
	}
	if stream.EncryptionType != "KMS" || stream.KeyID != keyARN {
		t.Fatalf("encrypted stream: got type %q key %q", stream.EncryptionType, stream.KeyID)
	}
}
