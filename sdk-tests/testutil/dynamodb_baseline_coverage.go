package testutil

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// dynamoDBBaselineCoverageTests adds baseline coverage for the operations
// that had no SDK test at all: backup describe/restore, global table
// settings, resource policies, and replica auto-scaling descriptions.
func (r *TestRunner) dynamoDBBaselineCoverageTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()
	setupErr := func(name string, err error) []TestResult {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: name,
			Status:   "FAIL",
			Error:    err.Error(),
		})
	}

	// --- Backup describe + restore --------------------------------------
	backupTable := fmt.Sprintf("baseline-backup-%d", suffix)
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(backupTable),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	}); err != nil {
		return setupErr("DescribeBackup_ReturnsBackupDetails", fmt.Errorf("create table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(backupTable)})
	if err := waitKinesisDestTableActive(ctx, client, backupTable); err != nil {
		return setupErr("DescribeBackup_ReturnsBackupDetails", fmt.Errorf("wait active: %v", err))
	}
	if _, err := client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(backupTable),
		Item: map[string]dynamodbtypes.AttributeValue{
			"pk": &dynamodbtypes.AttributeValueMemberS{Value: "kept"},
		},
	}); err != nil {
		return setupErr("DescribeBackup_ReturnsBackupDetails", fmt.Errorf("seed item: %v", err))
	}

	var backupArn string
	results = append(results, r.RunTest("dynamodb", "DescribeBackup_ReturnsBackupDetails", func() error {
		created, err := client.CreateBackup(ctx, &dynamodb.CreateBackupInput{
			TableName:  aws.String(backupTable),
			BackupName: aws.String(fmt.Sprintf("baseline-backup-%d", suffix)),
		})
		if err != nil {
			return err
		}
		if created.BackupDetails == nil || created.BackupDetails.BackupArn == nil {
			return fmt.Errorf("CreateBackup returned no BackupDetails")
		}
		backupArn = *created.BackupDetails.BackupArn

		desc, err := client.DescribeBackup(ctx, &dynamodb.DescribeBackupInput{BackupArn: aws.String(backupArn)})
		if err != nil {
			return err
		}
		if desc.BackupDescription == nil {
			return fmt.Errorf("DescribeBackup returned no BackupDescription")
		}
		bd := desc.BackupDescription.BackupDetails
		if bd == nil || aws.ToString(bd.BackupArn) != backupArn {
			return fmt.Errorf("DescribeBackup returned %+v", bd)
		}
		// The backup is materialised asynchronously; wait for AVAILABLE.
		deadline := time.Now().Add(10 * time.Second)
		for bd.BackupStatus != dynamodbtypes.BackupStatusAvailable {
			if time.Now().After(deadline) {
				return fmt.Errorf("backup did not become AVAILABLE, status %v", bd.BackupStatus)
			}
			time.Sleep(200 * time.Millisecond)
			desc, err = client.DescribeBackup(ctx, &dynamodb.DescribeBackupInput{BackupArn: aws.String(backupArn)})
			if err != nil {
				return err
			}
			bd = desc.BackupDescription.BackupDetails
		}
		if desc.BackupDescription.SourceTableDetails == nil ||
			aws.ToString(desc.BackupDescription.SourceTableDetails.TableName) != backupTable {
			return fmt.Errorf("SourceTableDetails missing or wrong table: %+v", desc.BackupDescription.SourceTableDetails)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "RestoreTableFromBackup_RestoresData", func() error {
		target := backupTable + "-restored"
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(target)})
		if _, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(backupTable),
			Key: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: "kept"},
			},
		}); err != nil {
			return err
		}

		resp, err := client.RestoreTableFromBackup(ctx, &dynamodb.RestoreTableFromBackupInput{
			BackupArn:       aws.String(backupArn),
			TargetTableName: aws.String(target),
		})
		if err != nil {
			return err
		}
		if resp.TableDescription == nil || aws.ToString(resp.TableDescription.TableName) != target {
			return fmt.Errorf("restore returned %+v", resp.TableDescription)
		}
		if err := waitKinesisDestTableActive(ctx, client, target); err != nil {
			return err
		}
		got, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(target),
			Key: map[string]dynamodbtypes.AttributeValue{
				"pk": &dynamodbtypes.AttributeValueMemberS{Value: "kept"},
			},
		})
		if err != nil {
			return err
		}
		if got.Item == nil {
			return fmt.Errorf("restored table lost the backed-up item")
		}

		// The restore summary persists on the table: DescribeTable keeps
		// reporting the source backup after the restore response is gone.
		targetDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(target)})
		if err != nil {
			return err
		}
		rs := targetDesc.Table.RestoreSummary
		if rs == nil {
			return fmt.Errorf("DescribeTable returned no RestoreSummary for a restored table")
		}
		if aws.ToString(rs.SourceBackupArn) != backupArn {
			return fmt.Errorf("expected SourceBackupArn=%s, got %v", backupArn, rs.SourceBackupArn)
		}
		if rs.SourceTableArn == nil || !strings.Contains(*rs.SourceTableArn, backupTable) {
			return fmt.Errorf("expected SourceTableArn of %s, got %v", backupTable, rs.SourceTableArn)
		}
		if rs.RestoreInProgress == nil || *rs.RestoreInProgress {
			return fmt.Errorf("expected RestoreInProgress=false, got %v", rs.RestoreInProgress)
		}
		return nil
	}))

	// --- Resource policy put/get/delete ---------------------------------
	policyTable := fmt.Sprintf("baseline-policy-%d", suffix)
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(policyTable),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	}); err != nil {
		return setupErr("PutResourcePolicy_RoundTrips", fmt.Errorf("create table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(policyTable)})
	if err := waitKinesisDestTableActive(ctx, client, policyTable); err != nil {
		return setupErr("PutResourcePolicy_RoundTrips", fmt.Errorf("wait active: %v", err))
	}
	tableDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(policyTable)})
	if err != nil || tableDesc.Table.TableArn == nil {
		return setupErr("PutResourcePolicy_RoundTrips", fmt.Errorf("describe for ARN: %v", err))
	}
	tableArn := *tableDesc.Table.TableArn

	results = append(results, r.RunTest("dynamodb", "PutResourcePolicy_RoundTrips", func() error {
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"dynamodb:GetItem","Resource":"` + tableArn + `"}]}`
		putResp, err := client.PutResourcePolicy(ctx, &dynamodb.PutResourcePolicyInput{
			ResourceArn: aws.String(tableArn),
			Policy:      aws.String(policy),
		})
		if err != nil {
			return err
		}
		if putResp.RevisionId == nil {
			return fmt.Errorf("PutResourcePolicy returned no RevisionId")
		}

		getResp, err := client.GetResourcePolicy(ctx, &dynamodb.GetResourcePolicyInput{
			ResourceArn: aws.String(tableArn),
		})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.Policy) != policy {
			return fmt.Errorf("policy round-trip mismatch: %v", getResp.Policy)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteResourcePolicy_RemovesPolicy", func() error {
		if _, err := client.DeleteResourcePolicy(ctx, &dynamodb.DeleteResourcePolicyInput{
			ResourceArn: aws.String(tableArn),
		}); err != nil {
			return err
		}
		_, err := client.GetResourcePolicy(ctx, &dynamodb.GetResourcePolicyInput{
			ResourceArn: aws.String(tableArn),
		})
		if err == nil {
			return fmt.Errorf("expected an error reading a deleted policy")
		}
		var policyMissing *dynamodbtypes.PolicyNotFoundException
		if !errors.As(err, &policyMissing) {
			return fmt.Errorf("expected PolicyNotFoundException, got %v", err)
		}
		return nil
	}))

	// --- Global table settings ------------------------------------------
	gtName := fmt.Sprintf("baseline-gt-settings-%d", suffix)
	// A global table links existing replica tables, so the backing table
	// shares the global table's name.
	if err := createGlobalTableTestTable(ctx, client, gtName, false); err != nil {
		return setupErr("GlobalTableSettings_DescribeAndUpdate", fmt.Errorf("create backing table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(gtName)})
	if _, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
		GlobalTableName: aws.String(gtName),
		ReplicationGroup: []dynamodbtypes.Replica{
			{RegionName: aws.String(r.region)},
		},
	}); err != nil {
		return setupErr("GlobalTableSettings_DescribeAndUpdate", fmt.Errorf("create global table: %v", err))
	}

	results = append(results, r.RunTest("dynamodb", "GlobalTableSettings_DescribeAndUpdate", func() error {
		desc, err := client.DescribeGlobalTableSettings(ctx, &dynamodb.DescribeGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
		})
		if err != nil {
			return err
		}
		if aws.ToString(desc.GlobalTableName) != gtName {
			return fmt.Errorf("expected GlobalTableName=%s, got %v", gtName, desc.GlobalTableName)
		}
		if len(desc.ReplicaSettings) != 1 || aws.ToString(desc.ReplicaSettings[0].RegionName) != r.region {
			return fmt.Errorf("expected one replica in %s, got %+v", r.region, desc.ReplicaSettings)
		}

		updated, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
			ReplicaSettingsUpdate: []dynamodbtypes.ReplicaSettingsUpdate{
				{
					RegionName:                          aws.String(r.region),
					ReplicaProvisionedReadCapacityUnits: aws.Int64(7),
				},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(updated.GlobalTableName) != gtName {
			return fmt.Errorf("update returned %v", updated.GlobalTableName)
		}

		after, err := client.DescribeGlobalTableSettings(ctx, &dynamodb.DescribeGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
		})
		if err != nil {
			return err
		}
		if aws.ToInt64(after.ReplicaSettings[0].ReplicaProvisionedReadCapacityUnits) != 7 {
			return fmt.Errorf("expected read units 7 after update, got %v", after.ReplicaSettings[0].ReplicaProvisionedReadCapacityUnits)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_GlobalMembersApplied", func() error {
		updated, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName:                          aws.String(gtName),
			GlobalTableBillingMode:                   dynamodbtypes.BillingModeProvisioned,
			GlobalTableProvisionedWriteCapacityUnits: aws.Int64(9),
			ReplicaSettingsUpdate: []dynamodbtypes.ReplicaSettingsUpdate{
				{
					RegionName:                          aws.String(r.region),
					ReplicaProvisionedReadCapacityUnits: aws.Int64(5),
				},
			},
		})
		if err != nil {
			return err
		}
		if aws.ToString(updated.GlobalTableName) != gtName {
			return fmt.Errorf("update returned %v", updated.GlobalTableName)
		}

		after, err := client.DescribeGlobalTableSettings(ctx, &dynamodb.DescribeGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
		})
		if err != nil {
			return err
		}
		if len(after.ReplicaSettings) != 1 {
			return fmt.Errorf("expected one replica, got %+v", after.ReplicaSettings)
		}
		rs := after.ReplicaSettings[0]
		if rs.ReplicaBillingModeSummary == nil || rs.ReplicaBillingModeSummary.BillingMode != dynamodbtypes.BillingModeProvisioned {
			return fmt.Errorf("expected PROVISIONED billing mode summary, got %+v", rs.ReplicaBillingModeSummary)
		}
		if aws.ToInt64(rs.ReplicaProvisionedWriteCapacityUnits) != 9 {
			return fmt.Errorf("expected write units 9, got %v", rs.ReplicaProvisionedWriteCapacityUnits)
		}
		if aws.ToInt64(rs.ReplicaProvisionedReadCapacityUnits) != 5 {
			return fmt.Errorf("expected read units 5, got %v", rs.ReplicaProvisionedReadCapacityUnits)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_InvalidBillingModeRejected", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName:        aws.String(gtName),
			GlobalTableBillingMode: dynamodbtypes.BillingMode("INVALID"),
		})
		if err == nil {
			return fmt.Errorf("expected error for off-enum billing mode")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_ZeroWriteCapacityRejected", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName:                          aws.String(gtName),
			GlobalTableProvisionedWriteCapacityUnits: aws.Int64(0),
		})
		if err == nil {
			return fmt.Errorf("expected error for zero write capacity")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	// The table-level write auto-scaling policy is echoed on every replica
	// together with its target tracking configuration.
	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_TargetTrackingRoundTrip", func() error {
		if _, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
			GlobalTableProvisionedWriteCapacityAutoScalingSettingsUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
				MinimumUnits: aws.Int64(5),
				MaximumUnits: aws.Int64(50),
				ScalingPolicyUpdate: &dynamodbtypes.AutoScalingPolicyUpdate{
					PolicyName: aws.String("gt-write-tracking"),
					TargetTrackingScalingPolicyConfiguration: &dynamodbtypes.AutoScalingTargetTrackingScalingPolicyConfigurationUpdate{
						TargetValue:      aws.Float64(50),
						DisableScaleIn:   aws.Bool(true),
						ScaleInCooldown:  aws.Int32(60),
						ScaleOutCooldown: aws.Int32(30),
					},
				},
			},
		}); err != nil {
			return err
		}
		after, err := client.DescribeGlobalTableSettings(ctx, &dynamodb.DescribeGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
		})
		if err != nil {
			return err
		}
		if len(after.ReplicaSettings) != 1 {
			return fmt.Errorf("expected one replica, got %+v", after.ReplicaSettings)
		}
		as := after.ReplicaSettings[0].ReplicaProvisionedWriteCapacityAutoScalingSettings
		if as == nil || len(as.ScalingPolicies) != 1 {
			return fmt.Errorf("expected one scaling policy, got %+v", as)
		}
		tt := as.ScalingPolicies[0].TargetTrackingScalingPolicyConfiguration
		if tt == nil {
			return fmt.Errorf("expected target tracking configuration on the policy")
		}
		if aws.ToFloat64(tt.TargetValue) != 50 || !aws.ToBool(tt.DisableScaleIn) ||
			aws.ToInt32(tt.ScaleInCooldown) != 60 || aws.ToInt32(tt.ScaleOutCooldown) != 30 {
			return fmt.Errorf("target tracking round-trip mismatch: %+v", tt)
		}
		return nil
	}))

	// A replica settings update naming a region that is not part of the
	// global table is rejected instead of silently ignored.
	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_UnknownReplicaRejected", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtName),
			ReplicaSettingsUpdate: []dynamodbtypes.ReplicaSettingsUpdate{
				{
					RegionName:                          aws.String("no-such-region-xyz"),
					ReplicaProvisionedReadCapacityUnits: aws.Int64(5),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for an unknown replica region")
		}
		return expectAWSErrorCode(err, "ReplicaNotFoundException")
	}))

	// The replica settings update list requires at least one entry.
	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_EmptyUpdateListRejected", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName:        aws.String(gtName),
			ReplicaSettingsUpdate: []dynamodbtypes.ReplicaSettingsUpdate{},
		})
		if err == nil {
			return fmt.Errorf("expected error for an empty replica settings update list")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	// --- Replica auto-scaling descriptions ------------------------------
	asTable := fmt.Sprintf("baseline-autoscaling-%d", suffix)
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(asTable),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	}); err != nil {
		return setupErr("TableReplicaAutoScaling_UpdateAndDescribe", fmt.Errorf("create table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(asTable)})
	if err := waitKinesisDestTableActive(ctx, client, asTable); err != nil {
		return setupErr("TableReplicaAutoScaling_UpdateAndDescribe", fmt.Errorf("wait active: %v", err))
	}

	results = append(results, r.RunTest("dynamodb", "TableReplicaAutoScaling_UpdateAndDescribe", func() error {
		before, err := client.DescribeTableReplicaAutoScaling(ctx, &dynamodb.DescribeTableReplicaAutoScalingInput{
			TableName: aws.String(asTable),
		})
		if err != nil {
			return err
		}
		if before.TableAutoScalingDescription == nil ||
			aws.ToString(before.TableAutoScalingDescription.TableName) != asTable {
			return fmt.Errorf("unexpected description: %+v", before.TableAutoScalingDescription)
		}

		updated, err := client.UpdateTableReplicaAutoScaling(ctx, &dynamodb.UpdateTableReplicaAutoScalingInput{
			TableName: aws.String(asTable),
			ReplicaUpdates: []dynamodbtypes.ReplicaAutoScalingUpdate{
				{
					RegionName: aws.String(r.region),
					ReplicaProvisionedReadCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
						MinimumUnits:        aws.Int64(1),
						MaximumUnits:        aws.Int64(10),
						AutoScalingDisabled: aws.Bool(false),
					},
				},
			},
		})
		if err != nil {
			return err
		}
		if updated.TableAutoScalingDescription == nil {
			return fmt.Errorf("update returned no description")
		}

		after, err := client.DescribeTableReplicaAutoScaling(ctx, &dynamodb.DescribeTableReplicaAutoScalingInput{
			TableName: aws.String(asTable),
		})
		if err != nil {
			return err
		}
		replicas := after.TableAutoScalingDescription.Replicas
		found := false
		for _, replica := range replicas {
			if aws.ToString(replica.RegionName) == r.region {
				found = true
				if replica.ReplicaProvisionedReadCapacityAutoScalingSettings == nil {
					return fmt.Errorf("read auto-scaling settings missing on the replica")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("replica %s not described: %+v", r.region, replicas)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTableReplicaAutoScaling_SettingsMemberValidationRejected", func() error {
		longArn := strings.Repeat("a", 1601)
		_, err := client.UpdateTableReplicaAutoScaling(ctx, &dynamodb.UpdateTableReplicaAutoScalingInput{
			TableName: aws.String(asTable),
			ReplicaUpdates: []dynamodbtypes.ReplicaAutoScalingUpdate{
				{
					RegionName: aws.String(r.region),
					ReplicaProvisionedReadCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
						AutoScalingRoleArn: aws.String(longArn),
					},
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for over-length nested AutoScalingRoleArn")
		}
		if err := expectAWSErrorCode(err, "ValidationException"); err != nil {
			return err
		}
		// TargetTrackingScalingPolicyConfiguration.TargetValue is bounded
		// by the documented metric range: zero, sub-minimum and
		// over-maximum values are all rejected.
		for _, bad := range []float64{0, 1e-200, 2e108} {
			_, err := client.UpdateTableReplicaAutoScaling(ctx, &dynamodb.UpdateTableReplicaAutoScalingInput{
				TableName: aws.String(asTable),
				ReplicaUpdates: []dynamodbtypes.ReplicaAutoScalingUpdate{
					{
						RegionName: aws.String(r.region),
						ReplicaProvisionedReadCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
							ScalingPolicyUpdate: &dynamodbtypes.AutoScalingPolicyUpdate{
								PolicyName: aws.String("range-policy"),
								TargetTrackingScalingPolicyConfiguration: &dynamodbtypes.AutoScalingTargetTrackingScalingPolicyConfigurationUpdate{
									TargetValue: aws.Float64(bad),
								},
							},
						},
					},
				},
			})
			if err == nil {
				return fmt.Errorf("expected error for out-of-range TargetValue %v", bad)
			}
			if err := expectAWSErrorCode(err, "ValidationException"); err != nil {
				return err
			}
		}
		return nil
	}))

	// --- Global-table settings optional members -------------------------
	gtGsiName := fmt.Sprintf("baseline-gt-gsi-%d", suffix)
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(gtGsiName),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("gpk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		GlobalSecondaryIndexes: []dynamodbtypes.GlobalSecondaryIndex{
			{
				IndexName: aws.String("gsi-index"),
				KeySchema: []dynamodbtypes.KeySchemaElement{
					{AttributeName: aws.String("gpk"), KeyType: dynamodbtypes.KeyTypeHash},
				},
				Projection: &dynamodbtypes.Projection{ProjectionType: dynamodbtypes.ProjectionTypeAll},
			},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
		StreamSpecification: &dynamodbtypes.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: dynamodbtypes.StreamViewTypeNewAndOldImages,
		},
	}); err != nil {
		return setupErr("UpdateGlobalTableSettings_ReplicaSettingsMembersRoundTrip", fmt.Errorf("create GSI table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(gtGsiName)})
	if err := waitKinesisDestTableActive(ctx, client, gtGsiName); err != nil {
		return setupErr("UpdateGlobalTableSettings_ReplicaSettingsMembersRoundTrip", fmt.Errorf("wait active: %v", err))
	}
	if _, err := client.CreateGlobalTable(ctx, &dynamodb.CreateGlobalTableInput{
		GlobalTableName: aws.String(gtGsiName),
		ReplicationGroup: []dynamodbtypes.Replica{
			{RegionName: aws.String(r.region)},
		},
	}); err != nil {
		return setupErr("UpdateGlobalTableSettings_ReplicaSettingsMembersRoundTrip", fmt.Errorf("create global table: %v", err))
	}

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_ReplicaSettingsMembersRoundTrip", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtGsiName),
			ReplicaSettingsUpdate: []dynamodbtypes.ReplicaSettingsUpdate{
				{
					RegionName: aws.String(r.region),
					ReplicaProvisionedReadCapacityAutoScalingSettingsUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
						MinimumUnits: aws.Int64(1),
						MaximumUnits: aws.Int64(10),
					},
					ReplicaGlobalSecondaryIndexSettingsUpdate: []dynamodbtypes.ReplicaGlobalSecondaryIndexSettingsUpdate{
						{
							IndexName:                    aws.String("gsi-index"),
							ProvisionedReadCapacityUnits: aws.Int64(5),
							ProvisionedReadCapacityAutoScalingSettingsUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
								MinimumUnits: aws.Int64(2),
								MaximumUnits: aws.Int64(8),
							},
						},
					},
					ReplicaTableClass: dynamodbtypes.TableClassStandardInfrequentAccess,
				},
			},
		})
		if err != nil {
			return err
		}

		after, err := client.DescribeGlobalTableSettings(ctx, &dynamodb.DescribeGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtGsiName),
		})
		if err != nil {
			return err
		}
		if len(after.ReplicaSettings) != 1 {
			return fmt.Errorf("expected one replica, got %+v", after.ReplicaSettings)
		}
		rs := after.ReplicaSettings[0]
		if rs.ReplicaProvisionedReadCapacityAutoScalingSettings == nil ||
			aws.ToInt64(rs.ReplicaProvisionedReadCapacityAutoScalingSettings.MinimumUnits) != 1 {
			return fmt.Errorf("expected replica read auto-scaling minimum 1, got %+v", rs.ReplicaProvisionedReadCapacityAutoScalingSettings)
		}
		if len(rs.ReplicaGlobalSecondaryIndexSettings) != 1 ||
			aws.ToString(rs.ReplicaGlobalSecondaryIndexSettings[0].IndexName) != "gsi-index" {
			return fmt.Errorf("expected one GSI setting for gsi-index, got %+v", rs.ReplicaGlobalSecondaryIndexSettings)
		}
		gsi := rs.ReplicaGlobalSecondaryIndexSettings[0]
		if aws.ToInt64(gsi.ProvisionedReadCapacityUnits) != 5 {
			return fmt.Errorf("expected GSI read units 5, got %v", gsi.ProvisionedReadCapacityUnits)
		}
		if gsi.ProvisionedReadCapacityAutoScalingSettings == nil ||
			aws.ToInt64(gsi.ProvisionedReadCapacityAutoScalingSettings.MinimumUnits) != 2 {
			return fmt.Errorf("expected GSI read auto-scaling minimum 2, got %+v", gsi.ProvisionedReadCapacityAutoScalingSettings)
		}
		if rs.ReplicaTableClassSummary == nil ||
			rs.ReplicaTableClassSummary.TableClass != dynamodbtypes.TableClassStandardInfrequentAccess {
			return fmt.Errorf("expected STANDARD_INFREQUENT_ACCESS table class, got %+v", rs.ReplicaTableClassSummary)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_GlobalGSISettingsRoundTrip", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtGsiName),
			GlobalTableProvisionedWriteCapacityAutoScalingSettingsUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
				MinimumUnits: aws.Int64(1),
				MaximumUnits: aws.Int64(20),
			},
			GlobalTableGlobalSecondaryIndexSettingsUpdate: []dynamodbtypes.GlobalTableGlobalSecondaryIndexSettingsUpdate{
				{
					IndexName:                     aws.String("gsi-index"),
					ProvisionedWriteCapacityUnits: aws.Int64(6),
					ProvisionedWriteCapacityAutoScalingSettingsUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
						MinimumUnits: aws.Int64(3),
						MaximumUnits: aws.Int64(9),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		after, err := client.DescribeGlobalTableSettings(ctx, &dynamodb.DescribeGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtGsiName),
		})
		if err != nil {
			return err
		}
		if len(after.ReplicaSettings) != 1 {
			return fmt.Errorf("expected one replica, got %+v", after.ReplicaSettings)
		}
		rs := after.ReplicaSettings[0]
		if rs.ReplicaProvisionedWriteCapacityAutoScalingSettings == nil ||
			aws.ToInt64(rs.ReplicaProvisionedWriteCapacityAutoScalingSettings.MinimumUnits) != 1 {
			return fmt.Errorf("expected replica write auto-scaling minimum 1, got %+v", rs.ReplicaProvisionedWriteCapacityAutoScalingSettings)
		}
		if len(rs.ReplicaGlobalSecondaryIndexSettings) != 1 {
			return fmt.Errorf("expected one merged GSI setting, got %+v", rs.ReplicaGlobalSecondaryIndexSettings)
		}
		gsi := rs.ReplicaGlobalSecondaryIndexSettings[0]
		if aws.ToInt64(gsi.ProvisionedWriteCapacityUnits) != 6 ||
			gsi.ProvisionedWriteCapacityAutoScalingSettings == nil ||
			aws.ToInt64(gsi.ProvisionedWriteCapacityAutoScalingSettings.MinimumUnits) != 3 {
			return fmt.Errorf("expected GSI write units 6 with auto-scaling minimum 3, got %+v", gsi)
		}
		// The replica-scoped read settings from the previous update must
		// survive the global-scope write update on the same index.
		if aws.ToInt64(gsi.ProvisionedReadCapacityUnits) != 5 ||
			gsi.ProvisionedReadCapacityAutoScalingSettings == nil ||
			aws.ToInt64(gsi.ProvisionedReadCapacityAutoScalingSettings.MinimumUnits) != 2 {
			return fmt.Errorf("expected GSI read settings preserved, got %+v", gsi)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateGlobalTableSettings_InvalidTableClassRejected", func() error {
		_, err := client.UpdateGlobalTableSettings(ctx, &dynamodb.UpdateGlobalTableSettingsInput{
			GlobalTableName: aws.String(gtGsiName),
			ReplicaSettingsUpdate: []dynamodbtypes.ReplicaSettingsUpdate{
				{
					RegionName:        aws.String(r.region),
					ReplicaTableClass: dynamodbtypes.TableClass("INVALID"),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for off-enum table class")
		}
		return expectAWSErrorCode(err, "ValidationException")
	}))

	// --- Table-level replica auto-scaling members ------------------------
	asGsiTable := fmt.Sprintf("baseline-autoscaling-gsi-%d", suffix)
	if _, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(asGsiTable),
		AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("gpk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamodbtypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
		},
		GlobalSecondaryIndexes: []dynamodbtypes.GlobalSecondaryIndex{
			{
				IndexName: aws.String("gsi-index"),
				KeySchema: []dynamodbtypes.KeySchemaElement{
					{AttributeName: aws.String("gpk"), KeyType: dynamodbtypes.KeyTypeHash},
				},
				Projection: &dynamodbtypes.Projection{ProjectionType: dynamodbtypes.ProjectionTypeAll},
			},
		},
		BillingMode: dynamodbtypes.BillingModePayPerRequest,
	}); err != nil {
		return setupErr("UpdateTableReplicaAutoScaling_TableLevelMembersRoundTrip", fmt.Errorf("create GSI table: %v", err))
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(asGsiTable)})
	if err := waitKinesisDestTableActive(ctx, client, asGsiTable); err != nil {
		return setupErr("UpdateTableReplicaAutoScaling_TableLevelMembersRoundTrip", fmt.Errorf("wait active: %v", err))
	}

	results = append(results, r.RunTest("dynamodb", "UpdateTableReplicaAutoScaling_TableLevelMembersRoundTrip", func() error {
		_, err := client.UpdateTableReplicaAutoScaling(ctx, &dynamodb.UpdateTableReplicaAutoScalingInput{
			TableName: aws.String(asGsiTable),
			ProvisionedWriteCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
				MinimumUnits: aws.Int64(2),
				MaximumUnits: aws.Int64(30),
			},
			GlobalSecondaryIndexUpdates: []dynamodbtypes.GlobalSecondaryIndexAutoScalingUpdate{
				{
					IndexName: aws.String("gsi-index"),
					ProvisionedWriteCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
						MinimumUnits: aws.Int64(4),
						MaximumUnits: aws.Int64(40),
					},
				},
			},
			ReplicaUpdates: []dynamodbtypes.ReplicaAutoScalingUpdate{
				{
					RegionName: aws.String(r.region),
					ReplicaProvisionedReadCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
						MinimumUnits: aws.Int64(1),
						MaximumUnits: aws.Int64(10),
					},
					ReplicaGlobalSecondaryIndexUpdates: []dynamodbtypes.ReplicaGlobalSecondaryIndexAutoScalingUpdate{
						{
							IndexName: aws.String("gsi-index"),
							ProvisionedReadCapacityAutoScalingUpdate: &dynamodbtypes.AutoScalingSettingsUpdate{
								MinimumUnits: aws.Int64(5),
								MaximumUnits: aws.Int64(50),
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		after, err := client.DescribeTableReplicaAutoScaling(ctx, &dynamodb.DescribeTableReplicaAutoScalingInput{
			TableName: aws.String(asGsiTable),
		})
		if err != nil {
			return err
		}
		replicas := after.TableAutoScalingDescription.Replicas
		found := false
		for _, replica := range replicas {
			if aws.ToString(replica.RegionName) != r.region {
				continue
			}
			found = true
			if replica.ReplicaProvisionedWriteCapacityAutoScalingSettings == nil ||
				aws.ToInt64(replica.ReplicaProvisionedWriteCapacityAutoScalingSettings.MinimumUnits) != 2 {
				return fmt.Errorf("expected table-level write auto-scaling minimum 2 on the replica, got %+v", replica.ReplicaProvisionedWriteCapacityAutoScalingSettings)
			}
			var gsi *dynamodbtypes.ReplicaGlobalSecondaryIndexAutoScalingDescription
			for i := range replica.GlobalSecondaryIndexes {
				if aws.ToString(replica.GlobalSecondaryIndexes[i].IndexName) == "gsi-index" {
					gsi = &replica.GlobalSecondaryIndexes[i]
					break
				}
			}
			if gsi == nil {
				return fmt.Errorf("gsi-index not described: %+v", replica.GlobalSecondaryIndexes)
			}
			if gsi.ProvisionedReadCapacityAutoScalingSettings == nil ||
				aws.ToInt64(gsi.ProvisionedReadCapacityAutoScalingSettings.MinimumUnits) != 5 {
				return fmt.Errorf("expected GSI read auto-scaling minimum 5, got %+v", gsi.ProvisionedReadCapacityAutoScalingSettings)
			}
			if gsi.ProvisionedWriteCapacityAutoScalingSettings == nil ||
				aws.ToInt64(gsi.ProvisionedWriteCapacityAutoScalingSettings.MinimumUnits) != 4 {
				return fmt.Errorf("expected GSI write auto-scaling minimum 4, got %+v", gsi.ProvisionedWriteCapacityAutoScalingSettings)
			}
			break
		}
		if !found {
			return fmt.Errorf("replica %s not described: %+v", r.region, replicas)
		}
		return nil
	}))

	return results
}
