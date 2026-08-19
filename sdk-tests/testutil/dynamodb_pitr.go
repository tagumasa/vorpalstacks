package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// dynamoDBPITRTests pins the point-in-time recovery contract: the
// continuous-backups status shape, restores that reconstruct the table
// state at the requested time, and the restore-window error behaviour.
func (r *TestRunner) dynamoDBPITRTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()

	// DescribeContinuousBackups reports the outer continuous-backups
	// status as ENABLED on every table; only the point-in-time recovery
	// status depends on the table settings, and no restorable window is
	// reported while recovery is disabled.
	descTable := fmt.Sprintf("PitrDesc-%d", suffix)
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(descTable),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	results = append(results, r.RunTest("dynamodb", "DescribeContinuousBackups_DisabledShape", func() error {
		if err != nil {
			return fmt.Errorf("create table: %w", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(descTable)})

		resp, err := client.DescribeContinuousBackups(ctx, &dynamodb.DescribeContinuousBackupsInput{
			TableName: aws.String(descTable),
		})
		if err != nil {
			return fmt.Errorf("describe continuous backups: %w", err)
		}
		desc := resp.ContinuousBackupsDescription
		if desc.ContinuousBackupsStatus != types.ContinuousBackupsStatusEnabled {
			return fmt.Errorf("outer ContinuousBackupsStatus = %q, want ENABLED (continuous backups exist on every table)", desc.ContinuousBackupsStatus)
		}
		pitr := desc.PointInTimeRecoveryDescription
		if pitr == nil || pitr.PointInTimeRecoveryStatus != types.PointInTimeRecoveryStatusDisabled {
			return fmt.Errorf("point-in-time recovery status = %+v, want DISABLED", pitr)
		}
		if pitr.EarliestRestorableDateTime != nil || pitr.LatestRestorableDateTime != nil {
			return fmt.Errorf("restorable window must be absent while recovery is disabled, got earliest=%v latest=%v", pitr.EarliestRestorableDateTime, pitr.LatestRestorableDateTime)
		}
		return nil
	}))

	pitrTable := fmt.Sprintf("PitrRestore-%d", suffix)
	restoreAt, restoreErr := r.pitrSeedTableWithHistory(ctx, client, pitrTable)
	defer func() {
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(pitrTable)})
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(fmt.Sprintf("PitrRestored-%d", suffix))})
	}()
	results = append(results, r.RunTest("dynamodb", "DescribeContinuousBackups_EnabledWindow", func() error {
		if restoreErr != nil {
			return fmt.Errorf("seed table: %w", restoreErr)
		}
		resp, err := client.DescribeContinuousBackups(ctx, &dynamodb.DescribeContinuousBackupsInput{
			TableName: aws.String(pitrTable),
		})
		if err != nil {
			return fmt.Errorf("describe continuous backups: %w", err)
		}
		pitr := resp.ContinuousBackupsDescription.PointInTimeRecoveryDescription
		if pitr.PointInTimeRecoveryStatus != types.PointInTimeRecoveryStatusEnabled {
			return fmt.Errorf("point-in-time recovery status = %q, want ENABLED", pitr.PointInTimeRecoveryStatus)
		}
		if pitr.EarliestRestorableDateTime == nil || pitr.LatestRestorableDateTime == nil {
			return fmt.Errorf("enabled recovery must report the restorable window, got %+v", pitr)
		}
		if pitr.EarliestRestorableDateTime.After(restoreAt) {
			return fmt.Errorf("earliest restorable %v must precede the seed writes at %v", pitr.EarliestRestorableDateTime, restoreAt)
		}
		return nil
	}))

	restoredTable := fmt.Sprintf("PitrRestored-%d", suffix)
	results = append(results, r.RunTest("dynamodb", "RestoreTableToPointInTime_RestoresStateAtTime", func() error {
		if restoreErr != nil {
			return fmt.Errorf("seed table: %w", restoreErr)
		}
		out, err := client.RestoreTableToPointInTime(ctx, &dynamodb.RestoreTableToPointInTimeInput{
			SourceTableName: aws.String(pitrTable),
			TargetTableName: aws.String(restoredTable),
			RestoreDateTime: aws.Time(restoreAt),
		})
		if err != nil {
			return fmt.Errorf("restore to point in time: %w", err)
		}
		summary := out.TableDescription.RestoreSummary
		if summary == nil {
			return fmt.Errorf("restore response must carry RestoreSummary")
		}
		if summary.RestoreInProgress == nil || *summary.RestoreInProgress {
			return fmt.Errorf("completed restore must report RestoreInProgress=false, got %v", summary.RestoreInProgress)
		}
		if summary.SourceTableArn == nil || !strings.Contains(*summary.SourceTableArn, pitrTable) {
			return fmt.Errorf("RestoreSummary source table ARN = %v, want the source table", summary.SourceTableArn)
		}
		if summary.RestoreDateTime == nil || diffSeconds(*summary.RestoreDateTime, restoreAt) > 2 {
			return fmt.Errorf("RestoreSummary restore time = %v, want ~%v", summary.RestoreDateTime, restoreAt)
		}

		// The item overwritten after the restore point must come back with
		// its value at the point; the item deleted after the point must be
		// resurrected.
		item, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName:      aws.String(restoredTable),
			Key:            map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "k"}},
			ConsistentRead: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("get restored item: %w", err)
		}
		attr, ok := item.Item["v"].(*types.AttributeValueMemberS)
		if !ok || attr.Value != "v1" {
			return fmt.Errorf("restored value = %+v, want v1 (the value at the restore point)", item.Item["v"])
		}
		deleted, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName:      aws.String(restoredTable),
			Key:            map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "k2"}},
			ConsistentRead: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("get resurrected item: %w", err)
		}
		if len(deleted.Item) == 0 {
			return fmt.Errorf("item deleted after the restore point must exist again in the restored table")
		}
		scan, err := client.Scan(ctx, &dynamodb.ScanInput{TableName: aws.String(restoredTable)})
		if err != nil {
			return fmt.Errorf("scan restored table: %w", err)
		}
		if scan.Count != 2 {
			return fmt.Errorf("restored table count = %d, want 2", scan.Count)
		}
		return nil
	}))

	noPitrTable := fmt.Sprintf("PitrNone-%d", suffix)
	results = append(results, r.RunTest("dynamodb", "RestoreTableToPointInTime_WithoutRecovery_Rejected", func() error {
		if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(noPitrTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		}); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(noPitrTable)})

		_, err := client.RestoreTableToPointInTime(ctx, &dynamodb.RestoreTableToPointInTimeInput{
			SourceTableName:         aws.String(noPitrTable),
			TargetTableName:         aws.String(noPitrTable + "-restored"),
			UseLatestRestorableTime: aws.Bool(true),
		})
		var unavailable *types.PointInTimeRecoveryUnavailableException
		if err == nil {
			return fmt.Errorf("restore without recovery enabled must fail")
		}
		if !errors.As(err, &unavailable) {
			return fmt.Errorf("error = %v, want PointInTimeRecoveryUnavailableException", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "RestoreTableToPointInTime_OutOfWindow_Rejected", func() error {
		if restoreErr != nil {
			return fmt.Errorf("seed table: %w", restoreErr)
		}
		for _, when := range []time.Time{time.Now().Add(-time.Hour), time.Now().Add(time.Hour)} {
			_, err := client.RestoreTableToPointInTime(ctx, &dynamodb.RestoreTableToPointInTimeInput{
				SourceTableName: aws.String(pitrTable),
				TargetTableName: aws.String(fmt.Sprintf("PitrBad-%d", suffix)),
				RestoreDateTime: aws.Time(when),
			})
			var invalidTime *types.InvalidRestoreTimeException
			if err == nil {
				return fmt.Errorf("restore at %v (outside the restorable window) must fail", when)
			}
			if !errors.As(err, &invalidTime) {
				return fmt.Errorf("error for %v = %v, want InvalidRestoreTimeException", when, err)
			}
		}
		return nil
	}))

	latestTable := fmt.Sprintf("PitrLatest-%d", suffix)
	results = append(results, r.RunTest("dynamodb", "RestoreTableToPointInTime_UseLatestRestorableTime", func() error {
		if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(latestTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		}); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(latestTable)})

		if _, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
			TableName: aws.String(latestTable),
			PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
				PointInTimeRecoveryEnabled: aws.Bool(true),
			},
		}); err != nil {
			return fmt.Errorf("enable recovery: %w", err)
		}
		if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(latestTable),
			Item:      map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "k"}, "v": &types.AttributeValueMemberS{Value: "now"}},
		}); err != nil {
			return fmt.Errorf("put item: %w", err)
		}

		restored := latestTable + "-restored"
		out, err := client.RestoreTableToPointInTime(ctx, &dynamodb.RestoreTableToPointInTimeInput{
			SourceTableName:         aws.String(latestTable),
			TargetTableName:         aws.String(restored),
			UseLatestRestorableTime: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("restore to latest restorable time: %w", err)
		}
		if out.TableDescription == nil || out.TableDescription.TableStatus != types.TableStatusActive {
			return fmt.Errorf("restored table must be ACTIVE, got %+v", out.TableDescription)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(restored)})
		item, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName:      aws.String(restored),
			Key:            map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "k"}},
			ConsistentRead: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("get restored item: %w", err)
		}
		if attr, ok := item.Item["v"].(*types.AttributeValueMemberS); !ok || attr.Value != "now" {
			return fmt.Errorf("restored value = %+v, want now (latest restorable state)", item.Item["v"])
		}
		return nil
	}))

	overrideTable := fmt.Sprintf("PitrOverride-%d", suffix)
	results = append(results, r.RunTest("dynamodb", "RestoreTableToPointInTime_IndexOverrideSelection", func() error {
		if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(overrideTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("lk"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("a"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("sk"), KeyType: types.KeyTypeRange},
			},
			LocalSecondaryIndexes: []types.LocalSecondaryIndex{{
				IndexName: aws.String("lsi-1"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("lk"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			}},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("gsi-a"),
					KeySchema: []types.KeySchemaElement{{AttributeName: aws.String("a"), KeyType: types.KeyTypeHash}},
					Projection: &types.Projection{
						ProjectionType:   types.ProjectionTypeInclude,
						NonKeyAttributes: []string{"v"},
					},
				},
				{
					IndexName:  aws.String("gsi-b"),
					KeySchema:  []types.KeySchemaElement{{AttributeName: aws.String("a"), KeyType: types.KeyTypeHash}},
					Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
				},
			},
			BillingMode: types.BillingModePayPerRequest,
		}); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(overrideTable)})

		if _, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
			TableName: aws.String(overrideTable),
			PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
				PointInTimeRecoveryEnabled: aws.Bool(true),
			},
		}); err != nil {
			return fmt.Errorf("enable recovery: %w", err)
		}

		restored := overrideTable + "-restored"
		_, err := client.RestoreTableToPointInTime(ctx, &dynamodb.RestoreTableToPointInTimeInput{
			SourceTableName:         aws.String(overrideTable),
			TargetTableName:         aws.String(restored),
			UseLatestRestorableTime: aws.Bool(true),
			GlobalSecondaryIndexOverride: []types.GlobalSecondaryIndex{{
				IndexName: aws.String("gsi-b"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("a"), KeyType: types.KeyTypeHash},
				},
				Projection: &types.Projection{
					ProjectionType:   types.ProjectionTypeInclude,
					NonKeyAttributes: []string{"v"},
				},
			}},
			LocalSecondaryIndexOverride: []types.LocalSecondaryIndex{{
				IndexName: aws.String("lsi-1"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("lk"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{
					ProjectionType:   types.ProjectionTypeInclude,
					NonKeyAttributes: []string{"v"},
				},
			}},
		})
		if err != nil {
			return fmt.Errorf("restore with index overrides: %w", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(restored)})

		desc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(restored)})
		if err != nil {
			return fmt.Errorf("describe restored table: %w", err)
		}
		if len(desc.Table.GlobalSecondaryIndexes) != 1 {
			return fmt.Errorf("override must select only gsi-b, got %d GSIs", len(desc.Table.GlobalSecondaryIndexes))
		}
		gsi := desc.Table.GlobalSecondaryIndexes[0]
		if aws.ToString(gsi.IndexName) != "gsi-b" {
			return fmt.Errorf("selected GSI = %q, want gsi-b", aws.ToString(gsi.IndexName))
		}
		if gsi.Projection == nil || len(gsi.Projection.NonKeyAttributes) != 1 || gsi.Projection.NonKeyAttributes[0] != "v" {
			return fmt.Errorf("selected GSI projection = %+v, want NonKeyAttributes [v]", gsi.Projection)
		}
		if len(desc.Table.LocalSecondaryIndexes) != 1 {
			return fmt.Errorf("override must keep exactly the selected LSI, got %d", len(desc.Table.LocalSecondaryIndexes))
		}
		lsi := desc.Table.LocalSecondaryIndexes[0]
		if lsi.Projection == nil || len(lsi.Projection.NonKeyAttributes) != 1 || lsi.Projection.NonKeyAttributes[0] != "v" {
			return fmt.Errorf("LSI projection = %+v, want NonKeyAttributes [v]", lsi.Projection)
		}
		return nil
	}))

	return results
}

// pitrSeedTableWithHistory creates a table with recovery enabled and writes
// a mutation history around the returned restore point: both items hold
// their "at the point" values before it, and are overwritten or deleted
// after it. The caller deletes the table.
func (r *TestRunner) pitrSeedTableWithHistory(ctx context.Context, client *dynamodb.Client, tableName string) (time.Time, error) {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("create table: %w", err)
	}
	if _, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
		TableName: aws.String(tableName),
		PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: aws.Bool(true),
		},
	}); err != nil {
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		return time.Time{}, fmt.Errorf("enable recovery: %w", err)
	}

	put := func(id, v string) error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: id},
				"v":  &types.AttributeValueMemberS{Value: v},
			},
		})
		return err
	}
	if err := put("k", "v1"); err != nil {
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		return time.Time{}, fmt.Errorf("put v1: %w", err)
	}
	if err := put("k2", "keep"); err != nil {
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		return time.Time{}, fmt.Errorf("put keep: %w", err)
	}

	// The restore point must separate the "at the point" writes from the
	// later mutations; DynamoDB timestamps carry second-level granularity,
	// so the margins exceed one second on both sides.
	time.Sleep(1100 * time.Millisecond)
	restoreAt := time.Now()
	time.Sleep(1100 * time.Millisecond)

	if err := put("k", "v2"); err != nil {
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		return time.Time{}, fmt.Errorf("put v2: %w", err)
	}
	if _, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: "k2"}},
	}); err != nil {
		client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})
		return time.Time{}, fmt.Errorf("delete k2: %w", err)
	}
	return restoreAt, nil
}

func diffSeconds(a, b time.Time) int64 {
	d := a.Sub(b)
	if d < 0 {
		d = -d
	}
	return int64(d.Seconds())
}
