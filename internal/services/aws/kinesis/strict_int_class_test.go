package kinesis

import (
	"context"
	"errors"
	"testing"
	"time"

	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

// TestStrictIntParamStates pins the three states of the read every typed
// integer member goes through: an absent member keeps the operation's
// documented default (present=false), integral values — JSON float64 or
// query-wire string — parse, and fractional, int32-overflowing, or
// non-numeric values are wire-type violations reported as
// InvalidArgumentException.
func TestStrictIntParamStates(t *testing.T) {
	params := map[string]interface{}{
		"Integral":   float64(7),
		"Negative":   float64(-3),
		"AsString":   "42",
		"Fractional": 7.5,
		"Overflow":   float64(5e9),
		"Exponent":   "5e9",
		"NotNumeric": "seven",
		"lowerFirst": float64(9),
		"AsNull":     nil,
	}

	if v, present, err := strictIntParam(params, "Absent"); present || err != nil || v != 0 {
		t.Fatalf("absent member: want (0, false, nil), got (%d, %t, %v)", v, present, err)
	}
	// A present null reads as absent — the awsJson1_1 conformance model
	// drops null structure values.
	if v, present, err := strictIntParam(params, "AsNull"); present || err != nil || v != 0 {
		t.Fatalf("null member: want (0, false, nil), got (%d, %t, %v)", v, present, err)
	}
	if v, present, err := strictIntParam(params, "Integral"); !present || err != nil || v != 7 {
		t.Fatalf("integral member: want (7, true, nil), got (%d, %t, %v)", v, present, err)
	}
	if v, present, err := strictIntParam(params, "Negative"); !present || err != nil || v != -3 {
		t.Fatalf("negative integral member: want (-3, true, nil), got (%d, %t, %v)", v, present, err)
	}
	if v, present, err := strictIntParam(params, "AsString"); !present || err != nil || v != 42 {
		t.Fatalf("string-form member: want (42, true, nil), got (%d, %t, %v)", v, present, err)
	}
	if v, present, err := strictIntParam(params, "LowerFirst"); !present || err != nil || v != 9 {
		t.Fatalf("lower-first key: want (9, true, nil), got (%d, %t, %v)", v, present, err)
	}
	for _, key := range []string{"Fractional", "Overflow", "Exponent", "NotNumeric"} {
		if _, present, err := strictIntParam(params, key); !present || err == nil || !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("%s member: want (0, true, InvalidArgumentException), got present=%t err=%v", key, present, err)
		}
	}
}

// TestStrictStringParamStates pins the three states of the typed String
// read: an absent member reports present=false, a string value reads in
// every key casing, and a present-but-non-string value — the wire form of a
// String or blob member is a string — is a wire-type violation reported as
// InvalidArgumentException, never the empty string the lenient readers
// coerce it to.
func TestStrictStringParamStates(t *testing.T) {
	params := map[string]interface{}{
		"AsString":   "YWJj",
		"lowerFirst": "eHl6",
		"AsNumber":   123,
		"AsBoolean":  true,
		"AsNull":     nil,
	}

	if v, present, err := strictStringParam(params, "Absent"); present || err != nil || v != "" {
		t.Fatalf("absent member: want (\"\", false, nil), got (%q, %t, %v)", v, present, err)
	}
	// A present null reads as absent — the protocol's dropped-null rule.
	if v, present, err := strictStringParam(params, "AsNull"); present || err != nil || v != "" {
		t.Fatalf("null member: want (\"\", false, nil), got (%q, %t, %v)", v, present, err)
	}
	if v, present, err := strictStringParam(params, "AsString"); !present || err != nil || v != "YWJj" {
		t.Fatalf("string member: want (\"YWJj\", true, nil), got (%q, %t, %v)", v, present, err)
	}
	if v, present, err := strictStringParam(params, "LowerFirst"); !present || err != nil || v != "eHl6" {
		t.Fatalf("lower-first-cased member: want (\"eHl6\", true, nil), got (%q, %t, %v)", v, present, err)
	}
	for _, key := range []string{"AsNumber", "AsBoolean"} {
		if _, present, err := strictStringParam(params, key); !present || !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("%s: want (present, InvalidArgumentException), got (%t, %v)", key, present, err)
		}
	}
}

// TestParseTimestampMember pins the strict Timestamp-member read: epoch
// seconds parse (fractional honoured), the RFC 3339 notation the member's
// documentation shows parses, and an unparseable value is a wire-type
// violation, never a silently dropped position.
func TestParseTimestampMember(t *testing.T) {
	ts, err := parseTimestampMember("1758000000")
	if err != nil || ts.Unix() != 1758000000 {
		t.Fatalf("integral epoch seconds: got (%v, %v)", ts, err)
	}
	frac, err := parseTimestampMember("1758000000.5")
	if err != nil || frac.Unix() != 1758000000 || frac.Nanosecond() != 500000000 {
		t.Fatalf("fractional epoch seconds: got (%v, %v)", frac, err)
	}
	rfc, err := parseTimestampMember("2016-04-04T19:58:46.480Z")
	if err != nil || rfc.Unix() != 1459799926 || rfc.Nanosecond() != 480000000 {
		t.Fatalf("RFC 3339 notation: got (%v, %v)", rfc, err)
	}
	if _, err := parseTimestampMember("not-a-timestamp"); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("unparseable timestamp: want InvalidArgumentException, got: %v", err)
	}
	// Non-finite spellings ParseFloat accepts and epoch values beyond the
	// int64-seconds range narrow implementation-dependently — both are
	// wire-type violations, never garbage instants.
	for _, invalid := range []string{"NaN", "nan", "Inf", "+Inf", "-Inf", "Infinity", "1e300", "-1e300", "9223372036854775808"} {
		if _, err := parseTimestampMember(invalid); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("%q: want InvalidArgumentException, got: %v", invalid, err)
		}
	}
	// The upper boundary is the representation's usable maximum: time.Time's
	// internal seconds field stores sec + unixToInternal (62135596800), so
	// values past math.MaxInt64−unixToInternal overflow it. The float64
	// grid steps 1024 there, so the largest value that still fits is
	// 1024-aligned below the ceiling and the next step is already past it;
	// the int64-boundary float further up is not a working instant either.
	// The low end (MinInt64) stays usable: adding the offset moves the sum
	// toward zero.
	if _, err := parseTimestampMember("9223371974719178752"); err != nil {
		t.Fatalf("largest usable epoch seconds: unexpected error %v", err)
	}
	for _, beyondUsable := range []string{"9223371974719179776", "9223372036854774784"} {
		if _, err := parseTimestampMember(beyondUsable); !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("%q: past the usable ceiling, want InvalidArgumentException, got: %v", beyondUsable, err)
		}
	}
	if _, err := parseTimestampMember("-9223372036854775808"); err != nil {
		t.Fatalf("min representable epoch seconds: unexpected error %v", err)
	}
}

// TestStrictTimestampParamStates pins the three states of the typed
// Timestamp read: an absent member reports present=false, the epoch-second
// JSON number the AWS SDKs serialise reads (the string-only readers drop it
// and silently disable the member), both documented string notations read,
// and a present-but-unreadable value is a wire-type violation reported as
// InvalidArgumentException.
func TestStrictTimestampParamStates(t *testing.T) {
	params := map[string]interface{}{
		"AsNumber":   1758000000.5,
		"AsEpochStr": "1758000000.5",
		"AsRFC3339":  "2016-04-04T19:58:46.480Z",
		"Garbage":    "not-a-timestamp",
		"AsBoolean":  true,
		"AsNull":     nil,
	}

	if v, present, err := strictTimestampParam(params, "Absent"); present || err != nil || v != "" {
		t.Fatalf("absent member: want (\"\", false, nil), got (%q, %t, %v)", v, present, err)
	}
	// A present null reads as absent — the protocol's dropped-null rule.
	if v, present, err := strictTimestampParam(params, "AsNull"); present || err != nil || v != "" {
		t.Fatalf("null member: want (\"\", false, nil), got (%q, %t, %v)", v, present, err)
	}
	if v, present, err := strictTimestampParam(params, "AsNumber"); !present || err != nil || v != "1758000000.5" {
		t.Fatalf("numeric member: want (\"1758000000.5\", true, nil), got (%q, %t, %v)", v, present, err)
	}
	if v, present, err := strictTimestampParam(params, "AsEpochStr"); !present || err != nil || v != "1758000000.5" {
		t.Fatalf("epoch string member: want (\"1758000000.5\", true, nil), got (%q, %t, %v)", v, present, err)
	}
	rfc, present, err := strictTimestampParam(params, "AsRFC3339")
	if !present || err != nil {
		t.Fatalf("RFC 3339 member: got (%q, %t, %v)", rfc, present, err)
	}
	if parsed, perr := parseTimestampMember(rfc); perr != nil || parsed.Unix() != 1459799926 {
		t.Fatalf("RFC 3339 canonical form: got (%q → %v, %v)", rfc, parsed, perr)
	}
	for _, key := range []string{"Garbage", "AsBoolean"} {
		if _, present, err := strictTimestampParam(params, key); !present || err == nil || !errors.Is(err, ErrInvalidArgument) {
			t.Fatalf("%s member: want (\"\", true, InvalidArgumentException), got present=%t err=%v", key, present, err)
		}
	}
}

// TestTimestampMembersAcceptNumericWireForm pins the fix at the operation
// boundary: the AWS SDKs serialise every Timestamp member as an epoch-second
// JSON number, so the handlers must read the numeric form end to end —
// GetShardIterator's AT_TIMESTAMP position, ListShards' ShardFilter and
// StreamCreationTimestamp, and ListStreamConsumers' StreamCreationTimestamp
// all receive float64 values here, not strings.
func TestTimestampMembersAcceptNumericWireForm(t *testing.T) {
	svc, store, reqCtx := newListPageTestEnv(t)
	if _, err := store.CreateStream("tsnum", 1, kinesisstore.StreamModeProvisioned, 0, 0, nil); err != nil {
		t.Fatalf("create stream: %v", err)
	}
	ctx := context.Background()

	iter, err := svc.GetShardIterator(ctx, reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "tsnum", "ShardId": "shardId-000000000000",
		"ShardIteratorType": "AT_TIMESTAMP", "Timestamp": 1758000000.5,
	}))
	if err != nil {
		t.Fatalf("AT_TIMESTAMP with a numeric Timestamp: %v", err)
	}
	iterMap, ok := iter.(map[string]interface{})
	if !ok || iterMap["ShardIterator"] == nil || iterMap["ShardIterator"] == "" {
		t.Fatalf("iterator response: got %#v", iter)
	}

	shards, err := svc.ListShards(ctx, reqCtx, parsedRequest(map[string]interface{}{
		"StreamName": "tsnum",
		// A future moment keeps the open set: the shard started before it
		// and remains open.
		"ShardFilter": map[string]interface{}{"Type": "AT_TIMESTAMP", "Timestamp": float64(time.Now().Add(time.Hour).Unix())},
	}))
	if err != nil {
		t.Fatalf("ListShards with a numeric filter timestamp: %v", err)
	}
	if m, ok := shards.(map[string]interface{}); !ok || len(m["Shards"].([]map[string]interface{})) != 1 {
		t.Fatalf("ListShards shards: got %#v", shards)
	}
}

// TestTypedMembersRejectNonIntegerWireValues pins the float-truncation
// class dead at the operation boundary: a fractional or exponent-notation
// count is rejected as InvalidArgumentException before any store is
// acquired. The lenient reader this replaces truncated float64(5e9) to
// 705032704, turning one CreateStream call into that many shard writes.
func TestTypedMembersRejectNonIntegerWireValues(t *testing.T) {
	svc := NewKinesisService("000000000000")
	reqCtx := failingStorageReqCtx(t)
	ctx := context.Background()

	cases := []struct {
		name   string
		invoke func() error
	}{
		{"CreateStream exponent shard count", func() error {
			_, err := svc.CreateStream(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "exp", "ShardCount": float64(5e9),
			}))
			return err
		}},
		{"CreateStream query-wire exponent shard count", func() error {
			_, err := svc.CreateStream(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "exp", "ShardCount": "5e9",
			}))
			return err
		}},
		{"CreateStream fractional max record size", func() error {
			_, err := svc.CreateStream(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "frac", "MaxRecordSizeInKiB": float64(2048.5),
			}))
			return err
		}},
		// ShardCount is a PositiveIntegerObject (range min 1): an explicit
		// zero on the wire is an out-of-range value, not the absent member
		// whose default is a single shard.
		{"CreateStream explicit zero shard count", func() error {
			_, err := svc.CreateStream(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"StreamName": "zero", "ShardCount": float64(0),
			}))
			return err
		}},
		{"ListStreams fractional limit", func() error {
			_, err := svc.ListStreams(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"Limit": 1.5,
			}))
			return err
		}},
		{"ListShards fractional max results", func() error {
			_, err := svc.ListShards(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"MaxResults": 1.5,
			}))
			return err
		}},
		{"ListStreamConsumers fractional max results", func() error {
			_, err := svc.ListStreamConsumers(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"MaxResults": 2.5,
			}))
			return err
		}},
		{"GetRecords fractional limit", func() error {
			_, err := svc.GetRecords(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"ShardIterator": "iterator", "Limit": 0.5,
			}))
			return err
		}},
		{"UpdateShardCount fractional target", func() error {
			_, err := svc.UpdateShardCount(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"TargetShardCount": 2.5,
			}))
			return err
		}},
		{"IncreaseStreamRetentionPeriod fractional hours", func() error {
			_, err := svc.IncreaseStreamRetentionPeriod(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"RetentionPeriodHours": float64(24.5),
			}))
			return err
		}},
		{"DecreaseStreamRetentionPeriod fractional hours", func() error {
			_, err := svc.DecreaseStreamRetentionPeriod(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"RetentionPeriodHours": "36.0",
			}))
			return err
		}},
		{"UpdateMaxRecordSize fractional size", func() error {
			_, err := svc.UpdateMaxRecordSize(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"MaxRecordSizeInKiB": float64(2048.5),
			}))
			return err
		}},
		{"UpdateStreamWarmThroughput fractional throughput", func() error {
			_, err := svc.UpdateStreamWarmThroughput(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"WarmThroughputMiBps": "3.5",
			}))
			return err
		}},
		{"UpdateStreamMode fractional throughput", func() error {
			_, err := svc.UpdateStreamMode(ctx, reqCtx, parsedRequest(map[string]interface{}{
				"WarmThroughputMiBps": 1.5,
			}))
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.invoke(); !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("expected InvalidArgumentException, got: %v", err)
			}
		})
	}
}
