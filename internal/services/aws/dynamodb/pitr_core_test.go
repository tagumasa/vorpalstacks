package dynamodb

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestPitrEarliestRestorableTrailingWindow(t *testing.T) {
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

	// Recovery enabled an hour ago: the journal only reaches back to the
	// enable moment, so that remains the earliest restorable time.
	fresh := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.Add(-time.Hour),
	}
	if got := pitrEarliestRestorable(fresh, now); !got.Equal(fresh.EarliestRestorableDateTime) {
		t.Fatalf("fresh enable: earliest = %v, want the enable time %v", got, fresh.EarliestRestorableDateTime)
	}

	// Recovery enabled 40 days ago with the default period: the trailing
	// 35-day edge is later than the enable moment and wins.
	aged := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.AddDate(0, 0, -40),
	}
	want := now.AddDate(0, 0, -pitrDefaultRecoveryPeriodDays)
	if got := pitrEarliestRestorable(aged, now); !got.Equal(want) {
		t.Fatalf("aged enable: earliest = %v, want the trailing edge %v", got, want)
	}

	// A configured 1-day period on a 3-day-old enable moves the edge to
	// now-1d.
	shortPeriod := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.AddDate(0, 0, -3),
		RecoveryPeriodInDays:       1,
	}
	want = now.AddDate(0, 0, -1)
	if got := pitrEarliestRestorable(shortPeriod, now); !got.Equal(want) {
		t.Fatalf("1-day period: earliest = %v, want %v", got, want)
	}

	// The trailing edge never precedes the enable moment: the journal holds
	// nothing older than it.
	shortEnable := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.Add(-12 * time.Hour),
		RecoveryPeriodInDays:       1,
	}
	if got := pitrEarliestRestorable(shortEnable, now); !got.Equal(shortEnable.EarliestRestorableDateTime) {
		t.Fatalf("recent enable with 1-day period: earliest = %v, want the enable time %v", got, shortEnable.EarliestRestorableDateTime)
	}
}

// TestUpdateContinuousBackupsResponseCarriesWindow pins the update
// response against the documented response syntax: the operation returns
// the same ContinuousBackupsDescription shape as the describe plane, so an
// enabling request answers with the recovery status AND the restorable
// window it just persisted (EarliestRestorableDateTime,
// LatestRestorableDateTime, the configured RecoveryPeriodInDays), and a
// disabling request answers with the status alone — the description of a
// table with nothing restorable.
func TestUpdateContinuousBackupsResponseCarriesWindow(t *testing.T) {
	sm, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("storage manager: %v", err)
	}
	defer sm.Close()
	svc := &DynamoDBService{}
	svc.SetStorageManager(sm)

	store, err := svc.GetCachedStoreForRegion("us-east-1")
	if err != nil {
		t.Fatalf("region store: %v", err)
	}
	if _, err := store.Tables().Create(dbstore.CreateTableParams{
		Name:                 "PitrWindowTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "pk", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "pk", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	reqCtx := request.NewRequestContext(context.Background(), sm, "123456789012", "us-east-1")
	enable := func() map[string]interface{} {
		t.Helper()
		resp, err := svc.UpdateContinuousBackups(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "PitrWindowTable",
			"PointInTimeRecoverySpecification": map[string]interface{}{
				"PointInTimeRecoveryEnabled": true,
				"RecoveryPeriodInDays":       7,
			},
		}})
		if err != nil {
			t.Fatalf("enable: %v", err)
		}
		return resp.(map[string]interface{})["ContinuousBackupsDescription"].(map[string]interface{})
	}

	before := time.Now().Add(-time.Second)
	desc := enable()
	after := time.Now().Add(time.Second)
	pitr := desc["PointInTimeRecoveryDescription"].(map[string]interface{})
	if pitr["PointInTimeRecoveryStatus"] != "ENABLED" {
		t.Fatalf("enable: status = %#v", pitr["PointInTimeRecoveryStatus"])
	}
	if pitr["RecoveryPeriodInDays"] != 7 {
		t.Fatalf("enable: configured period missing: %#v", pitr)
	}
	earliest, ok := pitr["EarliestRestorableDateTime"].(int64)
	if !ok {
		t.Fatalf("enable: earliest restorable missing: %#v", pitr)
	}
	latest, ok := pitr["LatestRestorableDateTime"].(int64)
	if !ok {
		t.Fatalf("enable: latest restorable missing: %#v", pitr)
	}
	// The window the request persisted: recovery starts at the enable
	// moment (the journal holds nothing older) and the present is
	// restorable, so both bounds land inside the request's own instants.
	if earliest < before.Unix() || earliest > after.Unix() {
		t.Fatalf("enable: earliest = %d, want the enable moment", earliest)
	}
	if latest < before.Unix() || latest > after.Unix() {
		t.Fatalf("enable: latest = %d, want the request's now", latest)
	}

	// Disabling answers with the status alone: there is nothing restorable
	// to describe.
	resp, err := svc.UpdateContinuousBackups(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "PitrWindowTable",
		"PointInTimeRecoverySpecification": map[string]interface{}{
			"PointInTimeRecoveryEnabled": false,
		},
	}})
	if err != nil {
		t.Fatalf("disable: %v", err)
	}
	desc = resp.(map[string]interface{})["ContinuousBackupsDescription"].(map[string]interface{})
	pitr = desc["PointInTimeRecoveryDescription"].(map[string]interface{})
	if len(pitr) != 1 || pitr["PointInTimeRecoveryStatus"] != "DISABLED" {
		t.Fatalf("disable: want the status alone, got %#v", pitr)
	}
}
