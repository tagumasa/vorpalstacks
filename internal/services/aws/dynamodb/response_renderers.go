package dynamodb

import (
	"time"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Response renderers — one function per model response shape, shared by
// every operation that emits the shape. Each renderer builds the member
// list the vendored Smithy model defines for the shape from the persisted
// record alone: optional members the record does not carry stay absent,
// and members the shape requires are rendered even at their zero value.
// The table description family lives in table_response.go.
// ---------------------------------------------------------------------------

// buildBackupDetailsResponse renders the model's BackupDetails shape. The
// status is the caller's reported status, not the record's stored one: the
// delete response reports the terminal DELETED status of a backup whose
// returned record still carries its pre-delete state.
func buildBackupDetailsResponse(b *dbstore.Backup, status dbstore.BackupStatus) map[string]interface{} {
	details := map[string]interface{}{
		"BackupArn":              b.BackupArn,
		"BackupName":             b.BackupName,
		"BackupSizeBytes":        b.BackupSizeBytes,
		"BackupStatus":           string(status),
		"BackupType":             string(b.BackupType),
		"BackupCreationDateTime": b.BackupCreationDateTime.Unix(),
	}
	if !b.BackupExpiryDateTime.IsZero() {
		details["BackupExpiryDateTime"] = b.BackupExpiryDateTime.Unix()
	}
	return details
}

// buildSourceTableDetailsResponse renders the model's SourceTableDetails
// shape from the backup record's captured source-table fields. The shape's
// ProvisionedThroughput member targets the capacity pair (the same shape
// CreateTableInput carries), not the table description's history-bearing
// variant — the change dates and the decreases counter are not members
// here.
func buildSourceTableDetailsResponse(b *dbstore.Backup) map[string]interface{} {
	return map[string]interface{}{
		"TableName":             b.SourceTableName,
		"TableArn":              b.SourceTableArn,
		"TableId":               b.SourceTableId,
		"TableSizeBytes":        b.SourceTableSizeBytes,
		"TableCreationDateTime": b.SourceTableCreationTime.Unix(),
		"ItemCount":             b.SourceTableItemCount,
		"KeySchema":             buildKeySchemaResponse(b.KeySchema),
		"BillingMode":           string(b.BillingMode),
		"ProvisionedThroughput": buildProvisionedThroughputPairResponse(b.ProvisionedThroughput),
	}
}

// buildProvisionedThroughputPairResponse renders the capacity-pair
// ProvisionedThroughput shape — the member SourceTableDetails carries.
// Both members are required, so a record without provisioned throughput
// still renders the pair at zero.
func buildProvisionedThroughputPairResponse(pt *dbstore.ProvisionedThroughput) map[string]interface{} {
	if pt == nil {
		return map[string]interface{}{
			"ReadCapacityUnits":  int64(0),
			"WriteCapacityUnits": int64(0),
		}
	}
	return map[string]interface{}{
		"ReadCapacityUnits":  pt.ReadCapacityUnits,
		"WriteCapacityUnits": pt.WriteCapacityUnits,
	}
}

// buildProvisionedThroughputDescriptionResponse renders the model's
// ProvisionedThroughputDescription shape — the response member
// TableDescription and GlobalSecondaryIndexDescription carry. history
// includes the change-history members: the table plane emits the change
// dates and the decreases counter, the GSI plane renders the capacity pair
// alone. A record without provisioned throughput (a PAY_PER_REQUEST table)
// renders the documented zero pair — the members are NonNegativeLongObjects
// set to 0 in that mode — typed once here.
func buildProvisionedThroughputDescriptionResponse(pt *dbstore.ProvisionedThroughput, history bool) map[string]interface{} {
	desc := map[string]interface{}{}
	if pt != nil {
		desc["ReadCapacityUnits"] = pt.ReadCapacityUnits
		desc["WriteCapacityUnits"] = pt.WriteCapacityUnits
		if history {
			if !pt.LastDecreaseDateTime.IsZero() {
				desc["LastDecreaseDateTime"] = pt.LastDecreaseDateTime.Unix()
			}
			if !pt.LastIncreaseDateTime.IsZero() {
				desc["LastIncreaseDateTime"] = pt.LastIncreaseDateTime.Unix()
			}
			desc["NumberOfDecreasesToday"] = pt.NumberOfDecreasesToday
		}
		return desc
	}
	desc["ReadCapacityUnits"] = int64(0)
	desc["WriteCapacityUnits"] = int64(0)
	if history {
		desc["NumberOfDecreasesToday"] = int64(0)
	}
	return desc
}

// buildPointInTimeRecoveryDescriptionResponse renders the model's
// PointInTimeRecoveryDescription shape: the status alone when recovery is
// disabled, and the restorable window with the configured period when
// enabled. The window is the trailing recovery period ending at the
// present; mutations commit synchronously on this platform, so now is
// restorable.
func buildPointInTimeRecoveryDescriptionResponse(pitr *dbstore.PointInTimeRecoveryDescription, now time.Time) map[string]interface{} {
	description := map[string]interface{}{
		"PointInTimeRecoveryStatus": "DISABLED",
	}
	if pitr == nil || pitr.Status != dbstore.PITRStatusEnabled {
		return description
	}
	description["PointInTimeRecoveryStatus"] = "ENABLED"
	description["EarliestRestorableDateTime"] = pitrEarliestRestorable(pitr, now).Unix()
	description["LatestRestorableDateTime"] = now.Unix()
	if pitr.RecoveryPeriodInDays > 0 {
		description["RecoveryPeriodInDays"] = pitr.RecoveryPeriodInDays
	}
	return description
}

// buildStreamDescriptionResponse renders the model's StreamDescription
// shape from the describe result: the shards as sequence-number ranges
// whose ending number appears only once the shard closed, and the
// last-evaluated shard id only while more shards remain.
func buildStreamDescriptionResponse(result *DescribeStreamResult) map[string]interface{} {
	shards := make([]interface{}, 0, len(result.Shards))
	for _, sh := range result.Shards {
		seqRange := map[string]interface{}{
			"StartingSequenceNumber": sh.StartingSequenceNumber,
		}
		if sh.EndingSequenceNumber != "" {
			seqRange["EndingSequenceNumber"] = sh.EndingSequenceNumber
		}
		shards = append(shards, map[string]interface{}{
			"ShardId":             sh.ShardID,
			"SequenceNumberRange": seqRange,
		})
	}
	description := map[string]interface{}{
		"StreamArn":               result.StreamArn,
		"StreamLabel":             result.StreamLabel,
		"StreamStatus":            result.StreamStatus,
		"TableName":               result.TableName,
		"KeySchema":               buildKeySchemaResponse(result.KeySchema),
		"Shards":                  shards,
		"CreationRequestDateTime": result.CreationRequestDateTime,
	}
	// The view type is an optional member: a stream that carries none
	// keeps it absent — the protocol omits unset members.
	if result.StreamViewType != "" {
		description["StreamViewType"] = result.StreamViewType
	}
	if result.LastEvaluatedShardId != "" {
		description["LastEvaluatedShardId"] = result.LastEvaluatedShardId
	}
	return description
}
