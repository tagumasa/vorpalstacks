package kinesis

import (
	"context"
	"sort"
	"strings"
	"testing"

	"vorpalstacks/internal/common/request"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// The emitted-key vocabulary pin: every operation's response carries only
// the model's output members (emitted is a subset of the model's set) and
// every member the model marks required is present. Optional members are
// omitted when not applicable — an absent continuation token, a
// successful batch entry's failure members. The expected sets are
// transcribed from the vendored Smithy model (models/kinesis
// /2013-12-02): an operation whose output targets Unit has no members.
// StreamDescriptionSummary's StreamId is marked "Not Implemented.
// Reserved for future use." by the model and is absent.
func TestEmittedKeyVocabulary(t *testing.T) {
	svc, store, reqCtx := newVocabularyEnv(t)

	records := []interface{}{
		map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
	}
	iterator := func() string {
		it, err := store.CreateShardIterator("vocab", "shardId-000000000000", "TRIM_HORIZON", "", nil)
		if err != nil {
			t.Fatalf("create iterator: %v", err)
		}
		return it.IteratorID
	}

	cases := []struct {
		op       string
		invoke   func() (interface{}, error)
		model    []string
		required []string
	}{
		{"CreateStream", func() (interface{}, error) {
			return nil, nil // covered by CreateStream_Vocabulary below
		}, nil, nil},
		{"DeleteStream", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"DescribeStream", func() (interface{}, error) {
			return svc.DescribeStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab"}))
		}, []string{"StreamDescription"}, []string{"StreamDescription"}},
		{"DescribeStreamSummary", func() (interface{}, error) {
			return svc.DescribeStreamSummary(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab"}))
		}, []string{"StreamDescriptionSummary"}, []string{"StreamDescriptionSummary"}},
		{"ListStreams", func() (interface{}, error) {
			return svc.ListStreams(context.Background(), reqCtx, parsedRequest(nil))
		}, []string{"HasMoreStreams", "NextToken", "StreamNames", "StreamSummaries"}, []string{"HasMoreStreams", "StreamNames"}},
		{"UpdateStreamMode", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"DescribeStreamConsumer", func() (interface{}, error) {
			return svc.DescribeStreamConsumer(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"ConsumerARN": vocabConsumerARN(t, store)}))
		}, []string{"ConsumerDescription"}, []string{"ConsumerDescription"}},
		{"ListStreamConsumers", func() (interface{}, error) {
			return svc.ListStreamConsumers(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"StreamARN": vocabStreamARN(t, store)}))
		}, []string{"Consumers", "NextToken"}, []string{"Consumers"}},
		{"RegisterStreamConsumer", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"DeregisterStreamConsumer", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"GetRecords", func() (interface{}, error) {
			return svc.GetRecords(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"ShardIterator": iterator()}))
		}, []string{"ChildShards", "MillisBehindLatest", "NextShardIterator", "Records"}, []string{"Records"}},
		{"GetShardIterator", func() (interface{}, error) {
			return svc.GetShardIterator(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "ShardId": "shardId-000000000000", "ShardIteratorType": "LATEST",
			}))
		}, []string{"ShardIterator"}, []string{"ShardIterator"}},
		{"PutRecord", func() (interface{}, error) {
			return svc.PutRecord(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "Data": "data", "PartitionKey": "pk",
			}))
		}, []string{"EncryptionType", "SequenceNumber", "ShardId"}, []string{"SequenceNumber", "ShardId"}},
		{"PutRecords", func() (interface{}, error) {
			return svc.PutRecords(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "Records": records,
			}))
		}, []string{"EncryptionType", "FailedRecordCount", "Records"}, []string{"Records"}},
		{"ListShards", func() (interface{}, error) {
			return svc.ListShards(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab"}))
		}, []string{"NextToken", "Shards"}, []string{"Shards"}},
		{"SplitShard", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"MergeShards", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"UpdateShardCount", func() (interface{}, error) {
			return svc.UpdateShardCount(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "TargetShardCount": 1, "ScalingType": "UNIFORM_SCALING",
			}))
		}, []string{"CurrentShardCount", "StreamARN", "StreamName", "TargetShardCount"}, []string{"CurrentShardCount", "StreamName", "StreamARN", "TargetShardCount"}},
		{"IncreaseStreamRetentionPeriod", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"DecreaseStreamRetentionPeriod", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"DescribeLimits", func() (interface{}, error) {
			return svc.DescribeLimits(context.Background(), reqCtx, parsedRequest(nil))
		}, []string{"ChannelCount", "ChannelCountLimit", "OnDemandStreamCount", "OnDemandStreamCountLimit", "OpenShardCount", "ShardLimit"}, []string{"ShardLimit", "OpenShardCount", "OnDemandStreamCount", "OnDemandStreamCountLimit"}},
		{"EnableEnhancedMonitoring", func() (interface{}, error) {
			return svc.EnableEnhancedMonitoring(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "ShardLevelMetrics": []interface{}{"IncomingBytes"},
			}))
		}, []string{"CurrentShardLevelMetrics", "DesiredShardLevelMetrics", "StreamARN", "StreamName"}, []string{"StreamName", "CurrentShardLevelMetrics", "DesiredShardLevelMetrics", "StreamARN"}},
		{"DisableEnhancedMonitoring", func() (interface{}, error) {
			return svc.DisableEnhancedMonitoring(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "ShardLevelMetrics": []interface{}{"IncomingBytes"},
			}))
		}, []string{"CurrentShardLevelMetrics", "DesiredShardLevelMetrics", "StreamARN", "StreamName"}, []string{"StreamName", "CurrentShardLevelMetrics", "DesiredShardLevelMetrics", "StreamARN"}},
		{"StartStreamEncryption", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"StopStreamEncryption", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"GetResourcePolicy", func() (interface{}, error) {
			return svc.GetResourcePolicy(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"ResourceARN": vocabStreamARN(t, store),
			}))
		}, []string{"Policy"}, []string{"Policy"}},
		{"UpdateMaxRecordSize", func() (interface{}, error) {
			return nil, nil
		}, nil, nil},
		{"UpdateStreamWarmThroughput", func() (interface{}, error) {
			// The operation is scoped to on-demand streams, so the warm
			// fixture is one.
			if _, err := store.CreateStream("vocab_warm", 1, kinesisstore.StreamModeOnDemand, 0, 0, nil); err != nil {
				t.Fatalf("create on-demand warm fixture: %v", err)
			}
			return svc.UpdateStreamWarmThroughput(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab_warm", "WarmThroughputMiBps": 128,
			}))
		}, []string{"StreamARN", "StreamName", "WarmThroughput"}, []string{"StreamARN", "StreamName", "WarmThroughput"}},
		{"ListTagsForStream", func() (interface{}, error) {
			return svc.ListTagsForStream(context.Background(), reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab"}))
		}, []string{"HasMoreTags", "Tags"}, []string{"Tags", "HasMoreTags"}},
		{"ListTagsForResource", func() (interface{}, error) {
			return svc.ListTagsForResource(context.Background(), reqCtx, parsedRequest(map[string]interface{}{
				"ResourceARN": vocabStreamARN(t, store),
			}))
		}, []string{"Tags"}, []string{"Tags"}},
	}

	for _, tc := range cases {
		if tc.invoke == nil {
			continue
		}
		result, err := tc.invoke()
		if err != nil {
			t.Fatalf("%s: %v", tc.op, err)
		}
		if result == nil {
			continue
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: response type %T, want a map", tc.op, result)
		}
		assertVocabulary(t, tc.op, m, tc.model, tc.required)
	}
}

// TestEmittedKeyVocabularyUnitOutputs drives the Unit-output operations
// whose side effects the shared fixture cannot repeat (create/delete on
// taken names): each returns an empty body — the model types their
// outputs as Unit.
func TestEmittedKeyVocabularyUnitOutputs(t *testing.T) {
	svc, store, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()

	unit := []struct {
		name   string
		invoke func() (interface{}, error)
	}{
		{"CreateStream", func() (interface{}, error) {
			return svc.CreateStream(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab_unit", "ShardCount": 1}))
		}},
		{"UpdateStreamMode", func() (interface{}, error) {
			return svc.UpdateStreamMode(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamARN": vocabStreamARN(t, store), "StreamModeDetails": map[string]interface{}{"StreamMode": "PROVISIONED"}}))
		}},
		{"IncreaseStreamRetentionPeriod", func() (interface{}, error) {
			// The change must be directional: an equal value rejects.
			return svc.IncreaseStreamRetentionPeriod(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab", "RetentionPeriodHours": kinesisstore.MinRetentionPeriodHours + 1}))
		}},
		{"DecreaseStreamRetentionPeriod", func() (interface{}, error) {
			return svc.DecreaseStreamRetentionPeriod(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab", "RetentionPeriodHours": kinesisstore.MinRetentionPeriodHours}))
		}},
		{"SplitShard", func() (interface{}, error) {
			return svc.SplitShard(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab_unit", "ShardToSplit": "shardId-000000000000", "NewStartingHashKey": "9223372036854775808"}))
		}},
		{"MergeShards", func() (interface{}, error) {
			return svc.MergeShards(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab_unit", "ShardToSplit": "", "ShardToMerge": "shardId-000000000001", "AdjacentShardToMerge": "shardId-000000000002",
			}))
		}},
		{"StartStreamEncryption", func() (interface{}, error) {
			return svc.StartStreamEncryption(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab", "EncryptionType": "KMS", "KeyId": "key-1"}))
		}},
		{"StopStreamEncryption", func() (interface{}, error) {
			return svc.StopStreamEncryption(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "vocab", "EncryptionType": "KMS", "KeyId": "key-1",
			}))
		}},
		{"UpdateMaxRecordSize", func() (interface{}, error) {
			return svc.UpdateMaxRecordSize(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamARN": vocabStreamARN(t, store), "MaxRecordSizeInKiB": 1024}))
		}},
	}
	for _, tc := range unit {
		result, err := tc.invoke()
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		m, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("%s: response type %T, want a map", tc.name, result)
		}
		if len(m) != 0 {
			t.Fatalf("%s: Unit-output response carries members %v, want none", tc.name, keysOf(m))
		}
	}
}

// TestEmittedKeyVocabularyNestedShapes pins the nested response shapes the
// model defines separately: Consumer vs ConsumerDescription, the
// PutRecords result entry, the summary's WarmThroughput object, and the
// Record's fractional arrival timestamp.
func TestEmittedKeyVocabularyNestedShapes(t *testing.T) {
	svc, store, reqCtx := newVocabularyEnv(t)
	ctx := context.Background()

	t.Run("Consumer", func(t *testing.T) {
		result, err := svc.ListStreamConsumers(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamARN": vocabStreamARN(t, store)}))
		if err != nil {
			t.Fatalf("list consumers: %v", err)
		}
		consumers := result.(map[string]interface{})["Consumers"].([]map[string]interface{})
		assertVocabulary(t, "Consumer", consumers[0], []string{"ConsumerARN", "ConsumerCreationTimestamp", "ConsumerName", "ConsumerStatus"}, []string{"ConsumerARN", "ConsumerCreationTimestamp", "ConsumerName", "ConsumerStatus"})
	})
	t.Run("ConsumerDescription", func(t *testing.T) {
		result, err := svc.DescribeStreamConsumer(ctx, reqCtx, parsedRequest(map[string]interface{}{"ConsumerARN": vocabConsumerARN(t, store)}))
		if err != nil {
			t.Fatalf("describe consumer: %v", err)
		}
		desc := result.(map[string]interface{})["ConsumerDescription"].(map[string]interface{})
		assertVocabulary(t, "ConsumerDescription", desc, []string{"ConsumerARN", "ConsumerCreationTimestamp", "ConsumerName", "ConsumerStatus", "StreamARN"}, []string{"ConsumerARN", "ConsumerCreationTimestamp", "ConsumerName", "ConsumerStatus", "StreamARN"})
	})
	t.Run("PutRecordsResultEntry", func(t *testing.T) {
		// The per-outcome form: a successful record includes SequenceNumber
		// and ShardId and nothing of the failure form. The failed form
		// (ErrorCode and ErrorMessage) is the write-time channel the
		// ErrorCode member bounds to its two documented outcomes — an
		// entry that fails validation has no per-entry identity, so the
		// request itself rejects.
		result, err := svc.PutRecords(ctx, reqCtx, parsedRequest(map[string]interface{}{
			"StreamName": "vocab",
			"Records": []interface{}{
				map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
				map[string]interface{}{"Data": "data", "PartitionKey": "pk2"},
			},
		}))
		if err != nil {
			t.Fatalf("put records: %v", err)
		}
		entries := result.(map[string]interface{})["Records"].([]map[string]interface{})
		if len(entries) != 2 {
			t.Fatalf("put records entries: got %d, want 2", len(entries))
		}
		for i := range entries {
			assertVocabulary(t, "PutRecordsResultEntry(success)", entries[i], []string{"SequenceNumber", "ShardId"}, []string{"SequenceNumber", "ShardId"})
		}

		_, err = svc.PutRecords(ctx, reqCtx, parsedRequest(map[string]interface{}{
			"StreamName": "vocab",
			"Records": []interface{}{
				map[string]interface{}{"Data": "data", "PartitionKey": "pk"},
				map[string]interface{}{"Data": "data", "PartitionKey": strings.Repeat("k", 300)},
			},
		}))
		requireAWSCode(t, "batch entry that fails its own shape", err, "InvalidArgumentException")
	})
	t.Run("StreamDescriptionSummary", func(t *testing.T) {
		result, err := svc.DescribeStreamSummary(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab"}))
		if err != nil {
			t.Fatalf("describe summary: %v", err)
		}
		summary := result.(map[string]interface{})["StreamDescriptionSummary"].(map[string]interface{})
		assertVocabulary(t, "StreamDescriptionSummary", summary, []string{
			"ChannelCount", "ConsumerCount", "EncryptionType", "EnhancedMonitoring",
			"MaxRecordSizeInKiB", "OpenShardCount", "RetentionPeriodHours", "StreamARN",
			"StreamCreationTimestamp", "StreamModeDetails", "StreamName", "StreamStatus",
		}, []string{
			"ChannelCount", "ConsumerCount", "EncryptionType", "EnhancedMonitoring",
			"MaxRecordSizeInKiB", "OpenShardCount", "RetentionPeriodHours", "StreamARN",
			"StreamCreationTimestamp", "StreamModeDetails", "StreamName", "StreamStatus",
		})
		// KeyId is optional: an unencrypted stream omits the member rather
		// than answering an empty identifier.
		if _, ok := summary["KeyId"]; ok {
			t.Fatal("StreamDescriptionSummary emitted KeyId on an unencrypted stream")
		}
		// WarmThroughput is optional the same way: a stream with no
		// configured warm throughput omits the member rather than
		// reporting a zero capacity.
		if _, ok := summary["WarmThroughput"]; ok {
			t.Fatal("StreamDescriptionSummary emitted WarmThroughput on an unconfigured stream")
		}
		if _, err := store.CreateStream("vocab_warm", 1, kinesisstore.StreamModeProvisioned, 0, 1024, nil); err != nil {
			t.Fatalf("create warm fixture: %v", err)
		}
		configured, err := svc.DescribeStreamSummary(ctx, reqCtx, parsedRequest(map[string]interface{}{"StreamName": "vocab_warm"}))
		if err != nil {
			t.Fatalf("describe warm summary: %v", err)
		}
		wt := configured.(map[string]interface{})["StreamDescriptionSummary"].(map[string]interface{})["WarmThroughput"].(map[string]interface{})
		assertVocabulary(t, "WarmThroughput", wt, []string{"CurrentMiBps", "TargetMiBps"}, []string{"CurrentMiBps", "TargetMiBps"})
	})
	t.Run("Record arrival timestamp is fractional", func(t *testing.T) {
		if _, _, err := store.PutRecordWithShardSelection("vocab", "pk", "data", ""); err != nil {
			t.Fatalf("put record: %v", err)
		}
		it, err := store.CreateShardIterator("vocab", "shardId-000000000000", "TRIM_HORIZON", "", nil)
		if err != nil {
			t.Fatalf("create iterator: %v", err)
		}
		result, err := svc.GetRecords(ctx, reqCtx, parsedRequest(map[string]interface{}{"ShardIterator": it.IteratorID}))
		if err != nil {
			t.Fatalf("get records: %v", err)
		}
		records := result.(map[string]interface{})["Records"].([]map[string]interface{})
		assertVocabulary(t, "Record", records[0], []string{"ApproximateArrivalTimestamp", "Data", "EncryptionType", "PartitionKey", "SequenceNumber"}, []string{"ApproximateArrivalTimestamp", "Data", "PartitionKey", "SequenceNumber"})
		ts, ok := records[0]["ApproximateArrivalTimestamp"].(float64)
		if !ok {
			t.Fatalf("arrival timestamp type %T, want float64", records[0]["ApproximateArrivalTimestamp"])
		}
		if ts < 1e9 {
			t.Fatalf("arrival timestamp %v does not look like epoch seconds", ts)
		}
	})
	t.Run("DryRun answers the model's rejection and performs no write", func(t *testing.T) {
		before, err := store.GetShard("vocab", "shardId-000000000000")
		if err != nil {
			t.Fatalf("get shard: %v", err)
		}
		_, err = svc.PutRecord(ctx, reqCtx, parsedRequest(map[string]interface{}{
			"StreamName": "vocab", "Data": "data", "PartitionKey": "pk", "DryRun": true,
		}))
		requireAWSCode(t, "dry-run put", err, "DryRunOperationException")
		after, err := store.GetShard("vocab", "shardId-000000000000")
		if err != nil {
			t.Fatalf("get shard after dry-run: %v", err)
		}
		if after.LatestSequenceNumber != before.LatestSequenceNumber {
			t.Fatal("dry-run put advanced the shard's sequence cursor")
		}
	})
	t.Run("GetShardIterator DryRun answers the model's rejection", func(t *testing.T) {
		// The dry run validates everything — a bad shard still answers the
		// shard's own error — then rejects with the model's dedicated
		// shape; no cursor is created.
		_, err := svc.GetShardIterator(ctx, reqCtx, parsedRequest(map[string]interface{}{
			"StreamName": "vocab", "ShardId": "shardId-000000000000", "ShardIteratorType": "LATEST", "DryRun": true,
		}))
		requireAWSCode(t, "dry-run iterator", err, "DryRunOperationException")
		_, err = svc.GetShardIterator(ctx, reqCtx, parsedRequest(map[string]interface{}{
			"StreamName": "vocab", "ShardId": "shardId-000000000999", "ShardIteratorType": "LATEST", "DryRun": true,
		}))
		requireAWSCode(t, "dry-run iterator on a bad shard", err, "ResourceNotFoundException")
	})
}

// newVocabularyEnv opens one environment whose fixtures the vocabulary
// table reads: a stream with tags and a policy, and a registered consumer.
func newVocabularyEnv(t *testing.T) (*KinesisService, *kinesisstore.KinesisStore, *request.RequestContext) {
	t.Helper()
	svc, store, reqCtx := newListPageTestEnv(t)
	stream, err := store.CreateStream("vocab", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil)
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	if err := store.Tag("vocab", map[string]string{"env": "vocab"}); err != nil {
		t.Fatalf("tag stream: %v", err)
	}
	if err := store.PutResourcePolicy(stream.StreamARN, `{"Version":"2012-10-17"}`); err != nil {
		t.Fatalf("put policy: %v", err)
	}
	if _, err := store.RegisterStreamConsumer(stream.StreamARN, "vocab-reader", nil); err != nil {
		t.Fatalf("register consumer: %v", err)
	}
	return svc, store, reqCtx
}

func vocabStreamARN(t *testing.T, store *kinesisstore.KinesisStore) string {
	t.Helper()
	stream, err := store.GetStream("vocab")
	if err != nil {
		t.Fatalf("get stream: %v", err)
	}
	return stream.StreamARN
}

func vocabConsumerARN(t *testing.T, store *kinesisstore.KinesisStore) string {
	t.Helper()
	consumers, err := store.ListStreamConsumers("vocab")
	if err != nil || len(consumers) == 0 {
		t.Fatalf("list consumers: %v (%d)", err, len(consumers))
	}
	return consumers[0].ConsumerARN
}

// assertVocabulary checks the model contract on one response map: every
// emitted key is a modelled member, and every required member is present.
func assertVocabulary(t *testing.T, name string, m map[string]interface{}, model, required []string) {
	t.Helper()
	modelSet := make(map[string]bool, len(model))
	for _, k := range model {
		modelSet[k] = true
	}
	for _, k := range keysOf(m) {
		if !modelSet[k] {
			t.Fatalf("%s emitted the unmodelled member %q", name, k)
		}
	}
	for _, k := range required {
		if _, ok := m[k]; !ok {
			t.Fatalf("%s omitted the required member %q", name, k)
		}
	}
}

func keysOf(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
