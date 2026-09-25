package dynamodb

import (
	"reflect"
	"testing"
	"time"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The renderer goldens pin one full response per shape per state: the
// member list is the vendored Smithy model's member list for the shape,
// and the values come from the persisted record alone. A member absent
// from the golden is absent from the wire.

// goldenBackup is the record every backup-plane golden renders: a USER
// backup of a provisioned table, captured with its source-table fields.
func goldenBackup() *dbstore.Backup {
	return &dbstore.Backup{
		BackupName:              "nightly",
		BackupArn:               "arn:aws:dynamodb:us-east-1:123456789012:table/GoldenTable/backup/01695097218000-d6299cbd",
		SourceTableName:         "GoldenTable",
		SourceTableArn:          "arn:aws:dynamodb:us-east-1:123456789012:table/GoldenTable",
		SourceTableId:           "aaaa1111-bbbb-2222-cccc-3333dddd4444",
		SourceTableCreationTime: time.Unix(1695000000, 0).UTC(),
		SourceTableSizeBytes:    4096,
		SourceTableItemCount:    42,
		BackupStatus:            dbstore.BackupStatusAvailable,
		BackupType:              dbstore.BackupTypeUser,
		BackupCreationDateTime:  time.Unix(1695097218, 0).UTC(),
		BackupSizeBytes:         2048,
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
		BillingMode: dbstore.BillingModeProvisioned,
		ProvisionedThroughput: &dbstore.ProvisionedThroughput{
			ReadCapacityUnits:  25,
			WriteCapacityUnits: 10,
		},
	}
}

func TestBackupDetailsRendererGolden(t *testing.T) {
	backup := goldenBackup()

	// The available render: the model's member set from the record — the
	// source-table name and ARN are BackupSummary members, not
	// BackupDetails members, so they stay out of this shape.
	available := buildBackupDetailsResponse(backup, backup.BackupStatus)
	wantAvailable := map[string]interface{}{
		"BackupArn":              backup.BackupArn,
		"BackupName":             "nightly",
		"BackupSizeBytes":        int64(2048),
		"BackupStatus":           "AVAILABLE",
		"BackupType":             "USER",
		"BackupCreationDateTime": int64(1695097218),
	}
	if !reflect.DeepEqual(available, wantAvailable) {
		t.Fatalf("available backup details: got %#v, want %#v", available, wantAvailable)
	}

	// The delete render reports the terminal status over the same record.
	deleted := buildBackupDetailsResponse(backup, dbstore.BackupStatusDeleted)
	wantDeleted := map[string]interface{}{
		"BackupArn":              backup.BackupArn,
		"BackupName":             "nightly",
		"BackupSizeBytes":        int64(2048),
		"BackupStatus":           "DELETED",
		"BackupType":             "USER",
		"BackupCreationDateTime": int64(1695097218),
	}
	if !reflect.DeepEqual(deleted, wantDeleted) {
		t.Fatalf("deleted backup details: got %#v, want %#v", deleted, wantDeleted)
	}

	// A backup carrying an expiry renders the member; one without keeps it
	// absent (the record's zero time is "no expiry", not the epoch).
	backup.BackupExpiryDateTime = time.Unix(1696000000, 0).UTC()
	withExpiry := buildBackupDetailsResponse(backup, backup.BackupStatus)
	if got, ok := withExpiry["BackupExpiryDateTime"]; !ok || got != int64(1696000000) {
		t.Fatalf("expiry member: got %#v", withExpiry["BackupExpiryDateTime"])
	}
	backup.BackupExpiryDateTime = time.Time{}
	if _, ok := buildBackupDetailsResponse(backup, backup.BackupStatus)["BackupExpiryDateTime"]; ok {
		t.Fatalf("zero expiry: member must stay absent")
	}
}

func TestSourceTableDetailsRendererGolden(t *testing.T) {
	backup := goldenBackup()

	// The model's member set from the captured source fields; the shape's
	// ProvisionedThroughput is the capacity pair (no history members —
	// dates and the decreases counter are TableDescription-plane members).
	got := buildSourceTableDetailsResponse(backup)
	want := map[string]interface{}{
		"TableName":             "GoldenTable",
		"TableArn":              backup.SourceTableArn,
		"TableId":               backup.SourceTableId,
		"TableSizeBytes":        int64(4096),
		"TableCreationDateTime": int64(1695000000),
		"ItemCount":             int64(42),
		"KeySchema": []map[string]interface{}{
			{"AttributeName": "pk", "KeyType": "HASH"},
			{"AttributeName": "sk", "KeyType": "RANGE"},
		},
		"BillingMode":           "PROVISIONED",
		"ProvisionedThroughput": map[string]interface{}{"ReadCapacityUnits": int64(25), "WriteCapacityUnits": int64(10)},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("source table details: got %#v, want %#v", got, want)
	}

	// A backup without captured throughput still renders the pair at zero:
	// the shape's members are required.
	backup.ProvisionedThroughput = nil
	pt := buildSourceTableDetailsResponse(backup)["ProvisionedThroughput"].(map[string]interface{})
	if pt["ReadCapacityUnits"] != int64(0) || pt["WriteCapacityUnits"] != int64(0) {
		t.Fatalf("zero pair: got %#v", pt)
	}
	if _, ok := pt["NumberOfDecreasesToday"]; ok {
		t.Fatalf("pair shape carries the decreases counter: %#v", pt)
	}
}

func TestProvisionedThroughputRendererGolden(t *testing.T) {
	// The table plane: capacity pair plus the change history, dates only
	// once a change happened.
	history := &dbstore.ProvisionedThroughput{
		ReadCapacityUnits:      5,
		WriteCapacityUnits:     5,
		LastIncreaseDateTime:   time.Unix(1695100000, 0).UTC(),
		NumberOfDecreasesToday: 2,
	}
	got := buildProvisionedThroughputDescriptionResponse(history, true)
	want := map[string]interface{}{
		"ReadCapacityUnits":      int64(5),
		"WriteCapacityUnits":     int64(5),
		"LastIncreaseDateTime":   int64(1695100000),
		"NumberOfDecreasesToday": int64(2),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("table plane history: got %#v, want %#v", got, want)
	}

	// The zero fallback is typed once: int64 zeros, the JSON number the
	// NonNegativeLongObject members define — never a bare Go int.
	zero := buildProvisionedThroughputDescriptionResponse(nil, true)
	if !reflect.DeepEqual(zero, map[string]interface{}{
		"ReadCapacityUnits":      int64(0),
		"WriteCapacityUnits":     int64(0),
		"NumberOfDecreasesToday": int64(0),
	}) {
		t.Fatalf("table plane zero fallback: got %#v", zero)
	}

	// The GSI plane: the capacity pair alone, with the same typed zero
	// fallback for a PAY_PER_REQUEST index.
	gsi := buildProvisionedThroughputDescriptionResponse(history, false)
	if !reflect.DeepEqual(gsi, map[string]interface{}{
		"ReadCapacityUnits":  int64(5),
		"WriteCapacityUnits": int64(5),
	}) {
		t.Fatalf("GSI plane: got %#v", gsi)
	}
	if got := buildProvisionedThroughputDescriptionResponse(nil, false); !reflect.DeepEqual(got, map[string]interface{}{
		"ReadCapacityUnits":  int64(0),
		"WriteCapacityUnits": int64(0),
	}) {
		t.Fatalf("GSI plane zero fallback: got %#v", got)
	}
}

func TestPointInTimeRecoveryRendererGolden(t *testing.T) {
	now := time.Unix(1695200000, 0).UTC()

	// No record and a disabled record both render the status alone.
	for _, pitr := range []*dbstore.PointInTimeRecoveryDescription{
		nil,
		{Status: dbstore.PITRStatusDisabled},
	} {
		got := buildPointInTimeRecoveryDescriptionResponse(pitr, now)
		want := map[string]interface{}{"PointInTimeRecoveryStatus": "DISABLED"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("disabled state: got %#v, want %#v", got, want)
		}
	}

	// An enabled record with a configured period: the window is the
	// trailing period ending at now — the record's own earliest moment is
	// older than the trailing edge, so the edge wins.
	enabled := &dbstore.PointInTimeRecoveryDescription{
		Status:                     dbstore.PITRStatusEnabled,
		EarliestRestorableDateTime: now.AddDate(0, 0, -100),
		RecoveryPeriodInDays:       7,
	}
	got := buildPointInTimeRecoveryDescriptionResponse(enabled, now)
	want := map[string]interface{}{
		"PointInTimeRecoveryStatus":  "ENABLED",
		"EarliestRestorableDateTime": now.AddDate(0, 0, -7).Unix(),
		"LatestRestorableDateTime":   now.Unix(),
		"RecoveryPeriodInDays":       7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("enabled with period: got %#v, want %#v", got, want)
	}

	// A record that never stored its period: the window falls back to the
	// default 35 days and the period member stays absent.
	enabled.RecoveryPeriodInDays = 0
	got = buildPointInTimeRecoveryDescriptionResponse(enabled, now)
	if got["EarliestRestorableDateTime"] != now.AddDate(0, 0, -35).Unix() {
		t.Fatalf("default period window: got %#v", got["EarliestRestorableDateTime"])
	}
	if _, ok := got["RecoveryPeriodInDays"]; ok {
		t.Fatalf("unstored period: member must stay absent: %#v", got)
	}
}

func TestStreamDescriptionRendererGolden(t *testing.T) {
	// The complete state: open and closed shards, and a page boundary.
	full := &DescribeStreamResult{
		StreamArn:      "arn:aws:dynamodb:us-east-1:123456789012:table/GoldenTable/stream/2026-09-24T00:00:00.000",
		StreamLabel:    "2026-09-24T00:00:00.000",
		StreamStatus:   "ENABLED",
		StreamViewType: "NEW_AND_OLD_IMAGES",
		TableName:      "GoldenTable",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
		},
		Shards: []ShardInfo{
			{ShardID: "shard-1", StartingSequenceNumber: "100", EndingSequenceNumber: "200"},
			{ShardID: "shard-2", StartingSequenceNumber: "201"},
		},
		CreationRequestDateTime: 1695300000,
		LastEvaluatedShardId:    "shard-2",
	}
	want := map[string]interface{}{
		"StreamArn":      full.StreamArn,
		"StreamLabel":    full.StreamLabel,
		"StreamStatus":   "ENABLED",
		"StreamViewType": "NEW_AND_OLD_IMAGES",
		"TableName":      "GoldenTable",
		"KeySchema":      []map[string]interface{}{{"AttributeName": "pk", "KeyType": "HASH"}},
		"Shards": []interface{}{
			map[string]interface{}{
				"ShardId": "shard-1",
				"SequenceNumberRange": map[string]interface{}{
					"StartingSequenceNumber": "100",
					"EndingSequenceNumber":   "200",
				},
			},
			map[string]interface{}{
				"ShardId": "shard-2",
				"SequenceNumberRange": map[string]interface{}{
					"StartingSequenceNumber": "201",
				},
			},
		},
		"CreationRequestDateTime": int64(1695300000),
		"LastEvaluatedShardId":    "shard-2",
	}
	if got := buildStreamDescriptionResponse(full); !reflect.DeepEqual(got, want) {
		t.Fatalf("full stream description: got %#v, want %#v", got, want)
	}

	// The terminal page: no boundary member, an empty shard list, and a
	// stream without a view type keeps the optional member absent.
	terminal := &DescribeStreamResult{
		StreamArn: full.StreamArn, StreamLabel: full.StreamLabel,
		StreamStatus: "ENABLED", StreamViewType: "",
		TableName: full.TableName, KeySchema: full.KeySchema,
		CreationRequestDateTime: full.CreationRequestDateTime,
	}
	got := buildStreamDescriptionResponse(terminal)
	if _, ok := got["LastEvaluatedShardId"]; ok {
		t.Fatalf("terminal page: boundary member must stay absent")
	}
	if _, ok := got["StreamViewType"]; ok {
		t.Fatalf("stream without a view type: member must stay absent")
	}
	if shards, ok := got["Shards"].([]interface{}); !ok || len(shards) != 0 {
		t.Fatalf("terminal page: empty shard list, got %#v", got["Shards"])
	}
}
