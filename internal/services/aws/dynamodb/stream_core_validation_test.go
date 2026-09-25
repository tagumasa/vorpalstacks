// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"testing"

	"context"
	"errors"
	"strings"
	"time"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestIntParamWithPresenceRejectsNonIntegralAndOutOfRangeNumbers pins the
// wire-number discipline of the integer extraction every limit and
// time-range member flows through: a JSON number converts silently only
// when it is integral and inside the platform's integer range — a
// non-integral number, one outside the range, or a present value of a
// non-numeric type is the invalid-parameter error, never a silently
// truncated or overflowed value.
func TestIntParamWithPresenceRejectsNonIntegralAndOutOfRangeNumbers(t *testing.T) {
	value, present, err := intParamWithPresence(map[string]interface{}{"Limit": float64(7)}, "Limit")
	if err != nil || !present || value != 7 {
		t.Fatalf("integral float64: got (%d, %v, %v)", value, present, err)
	}
	value, present, err = intParamWithPresence(map[string]interface{}{"Limit": 9}, "Limit")
	if err != nil || !present || value != 9 {
		t.Fatalf("int-typed value: got (%d, %v, %v)", value, present, err)
	}
	if value, present, err = intParamWithPresence(map[string]interface{}{}, "Limit"); err != nil || present || value != 0 {
		t.Fatalf("absent member: got (%d, %v, %v)", value, present, err)
	}
	for _, bad := range []interface{}{float64(1.5), float64(-0.25), 1e20, -1e20, "7", true, nil} {
		if _, _, err := intParamWithPresence(map[string]interface{}{"Limit": bad}, "Limit"); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("Limit %v: expected ErrInvalidParameter, got %v", bad, err)
		}
	}
}

func TestGetRecordsCoreValidatesLimitAndIterator(t *testing.T) {
	svc := &DynamoDBService{}

	// An explicit Limit below 1 is invalid (model PositiveLongObject
	// minimum 1) even though an omitted Limit defaults. These checks run
	// before any store access, which the nil store proves.
	if _, err := svc.getRecordsCore(nil, "t|0|0", 0, true); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("explicit Limit 0: expected ErrInvalidParameter, got %v", err)
	}

	// Values above the documented bound of 1000 are the model's
	// LimitExceededException, not a ValidationException.
	if _, err := svc.getRecordsCore(nil, "t|0|0", 1001, true); !errors.Is(err, ErrStreamsLimitExceeded) {
		t.Fatalf("Limit 1001: expected ErrStreamsLimitExceeded, got %v", err)
	}

	// Iterator verification needs the store's signing key, so the
	// iterator-shape cases run against a real store.
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	defer st.Close()
	store := dbstore.NewDynamoDBStore(st, st, "123456789012", "us-east-1")
	signingKey, err := store.Streams().IteratorSigningKey()
	if err != nil {
		t.Fatalf("signing key: %v", err)
	}

	if _, err := svc.getRecordsCore(store, "not-an-iterator", 10, true); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("malformed iterator: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.getRecordsCore(store, "t|0|0", 10, true); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("hand-crafted plaintext iterator: expected ErrInvalidParameter, got %v", err)
	}

	// A shard iterator older than the fifteen-minute lifetime is expired.
	restoreNow := streamTimeNow
	streamTimeNow = func() time.Time { return time.Now().Add(-20 * time.Minute) }
	stale := encodeShardIterator(signingKey, "arn:aws:dynamodb:us-east-1:123456789012:table/t/stream/old", "t", 0, "LATEST")
	streamTimeNow = restoreNow
	if _, err := svc.getRecordsCore(store, stale, 10, true); !errors.Is(err, ErrExpiredIterator) {
		t.Fatalf("expired iterator: expected ErrExpiredIterator, got %v", err)
	}
}

func TestGetShardIteratorCoreRequiresParameters(t *testing.T) {
	svc := &DynamoDBService{}

	if _, err := svc.getShardIteratorCore(nil, "", "shard-1", "LATEST", ""); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty StreamArn: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.getShardIteratorCore(nil, "arn", "", "LATEST", ""); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty ShardId: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.getShardIteratorCore(nil, "arn", "shard-1", "", ""); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("empty ShardIteratorType: expected ErrInvalidParameter, got %v", err)
	}
}

func TestDescribeStreamValidatesAndHonoursRequestMembers(t *testing.T) {
	svc, store, reqCtx, table := streamPlaneFixture(t)
	streamArn := table.StreamArn

	if _, err := svc.describeStreamCore(store, DescribeStreamInput{StreamArn: streamArn, Limit: 0, LimitSet: true}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("explicit Limit 0: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.describeStreamCore(store, DescribeStreamInput{StreamArn: streamArn, Limit: -3, LimitSet: true}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("negative Limit: expected ErrInvalidParameter, got %v", err)
	}
	if _, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "Limit": float64(0),
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("wire Limit 0: expected ErrInvalidParameter, got %v", err)
	}

	resp, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"StreamArn": streamArn}})
	if err != nil {
		t.Fatalf("describe stream: %v", err)
	}
	desc := resp.(map[string]interface{})["StreamDescription"].(map[string]interface{})
	if _, present := desc["LastEvaluatedShardId"]; present {
		t.Fatalf("single-shard page must not carry LastEvaluatedShardId: %v", desc)
	}
	shard := desc["Shards"].([]interface{})[0].(map[string]interface{})
	seqRange := shard["SequenceNumberRange"].(map[string]interface{})
	for _, member := range []string{"StartingSequenceNumber", "EndingSequenceNumber"} {
		v, _ := seqRange[member].(string)
		if len(v) < 21 || len(v) > 40 {
			t.Fatalf("%s = %q is outside the model's 21-40 sequence-number length", member, v)
		}
	}

	shardID := shard["ShardId"].(string)

	// The cursor positions the page after the matching shard: the last
	// shard leaves nothing to serve, and a cursor naming no shard of the
	// stream leaves nothing either — the same policy ListStreams applies.
	afterLast, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ExclusiveStartShardId": shardID,
	}})
	if err != nil {
		t.Fatalf("describe after last shard: %v", err)
	}
	if shards := afterLast.(map[string]interface{})["StreamDescription"].(map[string]interface{})["Shards"].([]interface{}); len(shards) != 0 {
		t.Fatalf("cursor at the last shard: expected an empty page, got %d shards", len(shards))
	}
	if _, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ExclusiveStartShardId": "shardId-00000001741631711871-1f6a72cf",
	}}); err != nil {
		t.Fatalf("unknown cursor: expected an empty page, got %v", err)
	}
	if _, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ExclusiveStartShardId": "short",
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("short ExclusiveStartShardId: expected ErrInvalidParameter, got %v", err)
	}

	// CHILD_SHARDS asks for the children of a shard; the platform's shards
	// never split, so the answer is empty. Any other type is off-enum.
	childShards, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn,
		"ShardFilter": map[string]interface{}{
			"Type":    "CHILD_SHARDS",
			"ShardId": shardID,
		},
	}})
	if err != nil {
		t.Fatalf("CHILD_SHARDS filter: %v", err)
	}
	if shards := childShards.(map[string]interface{})["StreamDescription"].(map[string]interface{})["Shards"].([]interface{}); len(shards) != 0 {
		t.Fatalf("CHILD_SHARDS over a never-splitting shard: expected no children, got %d", len(shards))
	}
	if _, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn":   streamArn,
		"ShardFilter": map[string]interface{}{"Type": "LATEST"},
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("off-enum ShardFilter type: expected ErrInvalidParameter, got %v", err)
	}

	// Limit 1 serves the single shard as one complete page.
	onePage, err := svc.DescribeStream(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "Limit": float64(1),
	}})
	if err != nil {
		t.Fatalf("Limit 1: %v", err)
	}
	if shards := onePage.(map[string]interface{})["StreamDescription"].(map[string]interface{})["Shards"].([]interface{}); len(shards) != 1 {
		t.Fatalf("Limit 1: expected the one shard, got %d", len(shards))
	}
}

func TestGetShardIteratorValidatesShardIdAndSequenceNumber(t *testing.T) {
	svc, store, reqCtx, table := streamPlaneFixture(t)
	streamArn := table.StreamArn
	shardID := streamShardID(t, svc, reqCtx, streamArn)

	getIterator := func(shard, iterType, seq string) error {
		_, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"StreamArn":         streamArn,
			"ShardId":           shard,
			"ShardIteratorType": iterType,
			"SequenceNumber":    seq,
		}})
		return err
	}
	iterator := func(iterType, seq string) error { return getIterator(shardID, iterType, seq) }

	// The model bounds ShardId at 28-65 characters; anything shorter is a
	// validation error before the shard is resolved.
	if err := getIterator("x", "LATEST", ""); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("1-char ShardId: expected ErrInvalidParameter, got %v", err)
	}
	// A well-formed id that is not the stream's single shard names a shard
	// that does not exist.
	if err := getIterator("shardId-00000001741631711871-1f6a72cf", "LATEST", ""); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("unknown well-formed ShardId: expected ErrResourceNotFound, got %v", err)
	}

	// Drain the shard the way a client does and keep the wire sequence
	// numbers it hands back.
	itResp, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ShardId": shardID, "ShardIteratorType": "TRIM_HORIZON",
	}})
	if err != nil {
		t.Fatalf("TRIM_HORIZON iterator: %v", err)
	}
	recResp, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": itResp.(map[string]interface{})["ShardIterator"],
	}})
	if err != nil {
		t.Fatalf("get records: %v", err)
	}
	records := recResp.(map[string]interface{})["Records"].([]interface{})
	if len(records) != 2 {
		t.Fatalf("expected the two captured records, got %d", len(records))
	}
	firstSeq := records[0].(*dbstore.StreamRecord).Dynamodb.SequenceNumber
	secondSeq := records[1].(*dbstore.StreamRecord).Dynamodb.SequenceNumber
	if len(firstSeq) < 21 {
		t.Fatalf("record sequence number %q is shorter than the model's 21-character minimum", firstSeq)
	}

	// The emitted form is the accepted form: AT the first record's number
	// serves that record first; AFTER the last serves nothing.
	atResp, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ShardId": shardID,
		"ShardIteratorType": "AT_SEQUENCE_NUMBER", "SequenceNumber": firstSeq,
	}})
	if err != nil {
		t.Fatalf("AT_SEQUENCE_NUMBER round-trip: %v", err)
	}
	atRecords, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": atResp.(map[string]interface{})["ShardIterator"],
	}})
	if err != nil {
		t.Fatalf("get records at position: %v", err)
	}
	atList := atRecords.(map[string]interface{})["Records"].([]interface{})
	if len(atList) == 0 || atList[0].(*dbstore.StreamRecord).Dynamodb.SequenceNumber != firstSeq {
		t.Fatalf("AT_SEQUENCE_NUMBER must serve the named record first, got %v", atList)
	}
	afterResp, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ShardId": shardID,
		"ShardIteratorType": "AFTER_SEQUENCE_NUMBER", "SequenceNumber": secondSeq,
	}})
	if err != nil {
		t.Fatalf("AFTER_SEQUENCE_NUMBER round-trip: %v", err)
	}
	afterRecords, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": afterResp.(map[string]interface{})["ShardIterator"],
	}})
	if err != nil {
		t.Fatalf("get records after position: %v", err)
	}
	if afterList := afterRecords.(map[string]interface{})["Records"].([]interface{}); len(afterList) != 0 {
		t.Fatalf("AFTER the last record must serve nothing, got %d records", len(afterList))
	}

	if err := iterator("AT_SEQUENCE_NUMBER", "5"); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("short SequenceNumber: expected ErrInvalidParameter, got %v", err)
	}
	if err := iterator("AT_SEQUENCE_NUMBER", strings.Repeat("9", 41)); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("41-digit SequenceNumber: expected ErrInvalidParameter, got %v", err)
	}

	// Trim the whole stream and pin the documented iterator-time trimmed
	// error: AT or AFTER a number whose first servable record is gone.
	if err := store.Streams().TrimOlderThan("StreamsPinTable", time.Now().Add(time.Hour)); err != nil {
		t.Fatalf("trim: %v", err)
	}
	if err := iterator("AT_SEQUENCE_NUMBER", firstSeq); !errors.Is(err, ErrTrimmedDataAccess) {
		t.Fatalf("AT a trimmed record: expected ErrTrimmedDataAccess, got %v", err)
	}
	if err := iterator("AFTER_SEQUENCE_NUMBER", firstSeq); !errors.Is(err, ErrTrimmedDataAccess) {
		t.Fatalf("AFTER into trimmed records: expected ErrTrimmedDataAccess, got %v", err)
	}
	// AFTER the trim floor itself stays readable: the position's first
	// record is the oldest survivor.
	if err := iterator("AFTER_SEQUENCE_NUMBER", secondSeq); err != nil {
		t.Fatalf("AFTER the trim floor: expected a readable position, got %v", err)
	}

	// A LATEST iterator frozen before the trim follows the floor the same
	// way: the fully trimmed stream (floor past the frozen position, no
	// survivor yet) serves an empty page rather than the trimmed-data
	// error a caller-named position earns.
	latestResp, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ShardId": shardID, "ShardIteratorType": "LATEST",
	}})
	if err != nil {
		t.Fatalf("fully trimmed LATEST iterator: %v", err)
	}
	latestRecords, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": latestResp.(map[string]interface{})["ShardIterator"],
	}})
	if err != nil {
		t.Fatalf("fully trimmed LATEST read: %v", err)
	}
	if latestList := latestRecords.(map[string]interface{})["Records"].([]interface{}); len(latestList) != 0 {
		t.Fatalf("fully trimmed LATEST read must serve an empty page, got %d records", len(latestList))
	}

	// A TRIM_HORIZON position follows the floor: the trim above removed
	// both records, a fresh write lands beyond the floor, and the horizon
	// read serves that oldest (and only) survivor — the documented horizon
	// semantics, not the trimmed-data error a caller-named position earns.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "StreamsPinTable",
		"Item":      map[string]interface{}{"id": map[string]interface{}{"S": "c"}},
	}}); err != nil {
		t.Fatalf("post-trim put: %v", err)
	}
	horizonResp, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": streamArn, "ShardId": shardID, "ShardIteratorType": "TRIM_HORIZON",
	}})
	if err != nil {
		t.Fatalf("post-trim TRIM_HORIZON iterator: %v", err)
	}
	horizonRecords, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": horizonResp.(map[string]interface{})["ShardIterator"],
	}})
	if err != nil {
		t.Fatalf("post-trim TRIM_HORIZON read: %v", err)
	}
	horizonList := horizonRecords.(map[string]interface{})["Records"].([]interface{})
	if len(horizonList) != 1 ||
		horizonList[0].(*dbstore.StreamRecord).Dynamodb.Keys["id"].(map[string]interface{})["S"] != "c" {
		t.Fatalf("post-trim TRIM_HORIZON must serve the oldest retained record alone, got %v", horizonList)
	}
}

func TestStreamRecordWriteFailureFailsTheItemWrite(t *testing.T) {
	svc, store, _, table := streamPlaneFixture(t)

	// Positive control: the same write through the healthy store commits.
	if _, err := svc.putItemCore(context.Background(), store, "us-east-1", PutItemCoreInput{
		Table: table,
		Item:  map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("healthy")},
	}); err != nil {
		t.Fatalf("healthy control write: %v", err)
	}

	faulted := streamsFaultStore{DynamoDBStoreInterface: store}
	if _, err := svc.putItemCore(context.Background(), faulted, "us-east-1", PutItemCoreInput{
		Table: table,
		Item:  map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("doomed")},
	}); err == nil {
		t.Fatal("item write with a failing stream record write: expected the transaction to fail")
	}
	item, getErr := store.Items().Get("StreamsPinTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("doomed")})
	if getErr == nil && item != nil {
		t.Fatal("the item committed without its stream record")
	}
}

func TestShardIteratorIsBoundToItsStreamGeneration(t *testing.T) {
	svc, store, reqCtx, table := streamPlaneFixture(t)
	firstArn := table.StreamArn

	oldIter, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": firstArn, "ShardId": streamShardID(t, svc, reqCtx, firstArn), "ShardIteratorType": "TRIM_HORIZON",
	}})
	if err != nil {
		t.Fatalf("iterator of the first generation: %v", err)
	}
	oldIterator := oldIter.(map[string]interface{})["ShardIterator"].(string)

	// Disable, then re-enable (the enabling update carries the view type
	// the model requires): the re-enabled stream is a new generation — a
	// fresh ARN and an empty record space, so its records begin with the
	// new generation's own writes.
	streamSpec := func(enabled bool) map[string]interface{} {
		spec := map[string]interface{}{"StreamEnabled": enabled}
		if enabled {
			spec["StreamViewType"] = "NEW_AND_OLD_IMAGES"
		}
		return spec
	}
	for _, enabled := range []bool{false, true} {
		if _, err := svc.UpdateTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":           "StreamsPinTable",
			"StreamSpecification": streamSpec(enabled),
		}}); err != nil {
			t.Fatalf("update table stream enabled=%v: %v", enabled, err)
		}
	}
	updated, err := store.Tables().Get("StreamsPinTable")
	if err != nil {
		t.Fatalf("reload table: %v", err)
	}
	if updated.StreamArn == "" || updated.StreamArn == firstArn {
		t.Fatalf("re-enable must generate a fresh stream ARN, got %q after %q", updated.StreamArn, firstArn)
	}

	if _, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": oldIterator,
	}}); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("superseded-generation iterator: expected ErrResourceNotFound, got %v", err)
	}

	// The successor generation starts empty: a fresh TRIM_HORIZON iterator
	// serves none of the superseded generation's records, and the record
	// space itself holds nothing until the new generation's writes arrive.
	newIter, err := svc.GetShardIterator(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"StreamArn": updated.StreamArn, "ShardId": streamShardID(t, svc, reqCtx, updated.StreamArn), "ShardIteratorType": "TRIM_HORIZON",
	}})
	if err != nil {
		t.Fatalf("iterator of the successor generation: %v", err)
	}
	resp, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": newIter.(map[string]interface{})["ShardIterator"],
	}})
	if err != nil {
		t.Fatalf("successor-generation read: %v", err)
	}
	if records := resp.(map[string]interface{})["Records"].([]interface{}); len(records) != 0 {
		t.Fatalf("successor generation must start empty, got %d records of the superseded generation", len(records))
	}
	if latest, latestErr := store.Streams().GetLatestSequenceForStream("StreamsPinTable", ""); latestErr != nil || latest != 0 {
		t.Fatalf("the generation change must leave an empty record space, latest seq = %d (err %v)", latest, latestErr)
	}

	// The successor generation's own write reads back under its own ARN.
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "StreamsPinTable",
		"Item":      map[string]interface{}{"id": map[string]interface{}{"S": "c"}},
	}}); err != nil {
		t.Fatalf("post-re-enable put: %v", err)
	}
	succResp, err := svc.GetRecords(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ShardIterator": resp.(map[string]interface{})["NextShardIterator"],
	}})
	if err != nil {
		t.Fatalf("successor-generation follow-up read: %v", err)
	}
	succRecords := succResp.(map[string]interface{})["Records"].([]interface{})
	if len(succRecords) != 1 {
		t.Fatalf("successor generation must serve its own write, got %d records", len(succRecords))
	}
	if rec := succRecords[0].(*dbstore.StreamRecord); rec.EventSourceARN != updated.StreamArn {
		t.Fatalf("successor record eventSourceARN = %q, want the successor generation's %q", rec.EventSourceARN, updated.StreamArn)
	}
}

func TestListStreamsRejectsBelowOneLimit(t *testing.T) {
	svc, _, reqCtx, _ := streamPlaneFixture(t)

	for _, limit := range []float64{0, -1} {
		if _, err := svc.ListStreams(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"Limit": limit,
		}}); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("ListStreams Limit %v: expected ErrInvalidParameter, got %v", limit, err)
		}
	}

	resp, err := svc.ListStreams(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{}})
	if err != nil {
		t.Fatalf("ListStreams default limit: %v", err)
	}
	streams := resp.(map[string]interface{})["Streams"].([]map[string]interface{})
	if len(streams) != 1 {
		t.Fatalf("expected the fixture's stream listed, got %v", streams)
	}
}

// TestStreamSpecificationMemberTypeContract pins the stream face's
// typed-member presence contract: the specification member is typed as a
// structure and StreamEnabled as a Boolean, so a present value of another
// type on either is the invalid-parameter error, never a silently skipped
// specification.
func TestStreamSpecificationMemberTypeContract(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)

	create := func(table string, streamSpec interface{}) error {
		_, err := svc.CreateTable(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":            table,
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
			"StreamSpecification":  streamSpec,
		}})
		return err
	}

	if err := create("StreamTypeFaceTable", "not-a-structure"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-structure StreamSpecification: expected the invalid-parameter rejection, got %v", err)
	}
	if err := create("StreamTypeFaceBoolTable", map[string]interface{}{"StreamEnabled": "yes"}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-boolean StreamEnabled: expected the invalid-parameter rejection, got %v", err)
	}
	if err := create("StreamTypeFaceOkTable", map[string]interface{}{"StreamEnabled": true, "StreamViewType": "NEW_AND_OLD_IMAGES"}); err != nil {
		t.Fatalf("well-typed stream specification: %v", err)
	}
}
