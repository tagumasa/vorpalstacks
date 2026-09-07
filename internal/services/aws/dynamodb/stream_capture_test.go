package dynamodb

import (
	"encoding/json"
	"testing"
	"time"

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
		ApproximateCreationDateTimePrecision: kinesisPrecisionMicrosecond,
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
	if milliRecord.DynamoDB.ApproximateCreationDateTimePrecision != kinesisPrecisionMillisecond {
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
	if microRecord.DynamoDB.ApproximateCreationDateTimePrecision != kinesisPrecisionMicrosecond {
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
