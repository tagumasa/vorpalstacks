package dynamodb

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestQueryRecordsContributorEvent drives the Query handler end to end and
// verifies the single partition-series read event lands in the contributor
// counters.
func TestQueryRecordsContributorEvent(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	ctx := context.Background()
	reqCtx := request.NewRequestContext(ctx, mgr, "123456789012", "us-east-1")
	s := NewDynamoDBService("123456789012")

	store, err := s.store(reqCtx)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	if _, err := store.Tables().Create(
		"QTbl",
		[]*dbstore.KeySchemaElement{{AttributeName: "pk", KeyType: dbstore.KeyTypeHash}, {AttributeName: "sk", KeyType: dbstore.KeyTypeRange}},
		[]*dbstore.AttributeDefinition{{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS}, {AttributeName: "sk", AttributeType: dbstore.ScalarAttributeTypeS}},
		dbstore.BillingModePayPerRequest, nil, nil, nil, nil, nil, false,
	); err != nil {
		t.Fatalf("create table: %v", err)
	}
	tbl, err := store.Tables().Get("QTbl")
	if err != nil {
		t.Fatalf("get table: %v", err)
	}
	tbl.ContributorInsightsEnabled = true
	if err := store.Tables().Put(tbl); err != nil {
		t.Fatalf("enable insights: %v", err)
	}

	resp, err := s.Query(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "QTbl",
		"KeyConditionExpression":    "pk = :pk",
		"ExpressionAttributeValues": map[string]interface{}{":pk": map[string]interface{}{"S": "A"}},
	}})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if count := resp.(map[string]interface{})["ScannedCount"]; count != 0 {
		t.Fatalf("expected an empty partition, got %v", resp)
	}

	at := time.Now()
	stats, err := store.Contributors().TopKeys("QTbl", dbstore.ContributorLayoutPartitionKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys: %v", err)
	}
	if len(stats) != 1 || stats[0].Key != `["s:A"]` || stats[0].Count != 1 || stats[0].Units != 1 {
		t.Fatalf("expected one query event on the partition series, got %+v", stats)
	}
	full, err := store.Contributors().TopKeys("QTbl", dbstore.ContributorLayoutFullKey, at.Add(-time.Minute), at.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("top keys full: %v", err)
	}
	if len(full) != 0 {
		t.Fatalf("a query must not touch the full-key series, got %+v", full)
	}
}
