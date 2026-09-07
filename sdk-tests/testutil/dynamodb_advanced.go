package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (r *TestRunner) dynamoDBMainAdvancedTests(ctx context.Context, client *dynamodb.Client, tableName, tableARN string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "TagResource", func() error {
		_, err := client.TagResource(ctx, &dynamodb.TagResourceInput{
			ResourceArn: aws.String(tableARN),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("Test")},
			},
		})
		if err != nil {
			return err
		}
		listResp, listErr := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
			ResourceArn: aws.String(tableARN),
		})
		if listErr != nil {
			return fmt.Errorf("TagResource verification failed: %w", listErr)
		}
		found := false
		for _, t := range listResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				if t.Value == nil || *t.Value != "Test" {
					return fmt.Errorf("tag value mismatch: expected Test, got %v", t.Value)
				}
				found = true
			}
		}
		if !found {
			return fmt.Errorf("tag Environment not found after TagResource")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ListTagsOfResource", func() error {
		resp, err := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
			ResourceArn: aws.String(tableARN),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) == 0 {
			return fmt.Errorf("no tags found")
		}
		found := false
		for _, t := range resp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				if t.Value == nil || *t.Value != "Test" {
					return fmt.Errorf("tag value mismatch: expected Test, got %v", t.Value)
				}
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected tag with key 'Environment'")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &dynamodb.UntagResourceInput{
			ResourceArn: aws.String(tableARN),
			TagKeys:     []string{"Environment"},
		})
		if err != nil {
			return err
		}
		verifyResp, verifyErr := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
			ResourceArn: aws.String(tableARN),
		})
		if verifyErr != nil {
			return fmt.Errorf("UntagResource verification failed: %w", verifyErr)
		}
		for _, t := range verifyResp.Tags {
			if t.Key != nil && *t.Key == "Environment" {
				return fmt.Errorf("tag Environment should be removed after UntagResource")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateTimeToLive", func() error {
		resp, err := client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(tableName),
			TimeToLiveSpecification: &types.TimeToLiveSpecification{
				AttributeName: aws.String("ttl"),
				Enabled:       aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		if resp.TimeToLiveSpecification == nil {
			return fmt.Errorf("TimeToLiveSpecification is nil")
		}
		if !*resp.TimeToLiveSpecification.Enabled {
			return fmt.Errorf("expected Enabled=true, got false")
		}
		if resp.TimeToLiveSpecification.AttributeName == nil || *resp.TimeToLiveSpecification.AttributeName != "ttl" {
			return fmt.Errorf("expected AttributeName=ttl, got %v", resp.TimeToLiveSpecification.AttributeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeTimeToLive", func() error {
		resp, err := client.DescribeTimeToLive(ctx, &dynamodb.DescribeTimeToLiveInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.TimeToLiveDescription == nil {
			return fmt.Errorf("TTL description not found")
		}
		if resp.TimeToLiveDescription.AttributeName == nil || *resp.TimeToLiveDescription.AttributeName != "ttl" {
			return fmt.Errorf("expected AttributeName=ttl, got %v", resp.TimeToLiveDescription.AttributeName)
		}
		if resp.TimeToLiveDescription.TimeToLiveStatus != types.TimeToLiveStatusEnabled {
			return fmt.Errorf("expected TimeToLiveStatus=ENABLED, got %v", resp.TimeToLiveDescription.TimeToLiveStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateBackup", func() error {
		backupName := fmt.Sprintf("TestBackup-%d", time.Now().UnixNano())
		resp, err := client.CreateBackup(ctx, &dynamodb.CreateBackupInput{
			TableName:  aws.String(tableName),
			BackupName: aws.String(backupName),
		})
		if err != nil {
			return err
		}
		if resp.BackupDetails == nil {
			return fmt.Errorf("BackupDetails is nil")
		}
		if resp.BackupDetails.BackupName == nil || *resp.BackupDetails.BackupName != backupName {
			return fmt.Errorf("expected BackupName=%s, got %v", backupName, resp.BackupDetails.BackupName)
		}
		if resp.BackupDetails.BackupArn == nil || *resp.BackupDetails.BackupArn == "" {
			return fmt.Errorf("expected non-empty BackupArn")
		}
		client.DeleteBackup(ctx, &dynamodb.DeleteBackupInput{
			BackupArn: resp.BackupDetails.BackupArn,
		})
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ListBackups", func() error {
		resp, err := client.ListBackups(ctx, &dynamodb.ListBackupsInput{})
		if err != nil {
			return err
		}
		if resp.BackupSummaries == nil {
			return fmt.Errorf("backup summaries is nil")
		}
		for _, bs := range resp.BackupSummaries {
			if bs.BackupArn == nil || *bs.BackupArn == "" {
				return fmt.Errorf("expected BackupArn in BackupSummary")
			}
			if bs.TableName == nil || *bs.TableName == "" {
				return fmt.Errorf("expected TableName in BackupSummary")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "ListBackups_BackupTypeFilter_Pagination", func() error {
		bpTable := fmt.Sprintf("BkPag-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, bpTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		backupArnSet := map[string]bool{}
		defer func() {
			for arn := range backupArnSet {
				client.DeleteBackup(ctx, &dynamodb.DeleteBackupInput{BackupArn: aws.String(arn)})
			}
		}()
		for i := 0; i < 3; i++ {
			resp, err := client.CreateBackup(ctx, &dynamodb.CreateBackupInput{
				TableName:  aws.String(bpTable),
				BackupName: aws.String(fmt.Sprintf("%s-%d", bpTable, i)),
			})
			if err != nil {
				return fmt.Errorf("create backup %d: %v", i, err)
			}
			backupArnSet[*resp.BackupDetails.BackupArn] = true
		}

		// BackupTypeFilter ALL lists every on-demand backup type, so the
		// walk must reach all three backups even though each page is limited
		// to one entry.
		var listed []string
		var startArn *string
		pageCount := 0
		for {
			resp, err := client.ListBackups(ctx, &dynamodb.ListBackupsInput{
				TableName:               aws.String(bpTable),
				BackupType:              types.BackupTypeFilterAll,
				Limit:                   aws.Int32(1),
				ExclusiveStartBackupArn: startArn,
			})
			if err != nil {
				return fmt.Errorf("list backups page: %v", err)
			}
			pageCount++
			for _, bs := range resp.BackupSummaries {
				if bs.BackupArn == nil {
					return fmt.Errorf("backup summary without BackupArn")
				}
				listed = append(listed, *bs.BackupArn)
			}
			if resp.LastEvaluatedBackupArn == nil || *resp.LastEvaluatedBackupArn == "" {
				break
			}
			startArn = resp.LastEvaluatedBackupArn
			if pageCount > 10 {
				return fmt.Errorf("pagination did not terminate after %d pages", pageCount)
			}
		}

		if len(listed) != 3 {
			return fmt.Errorf("expected 3 backups through the ALL filter, got %d (%v)", len(listed), listed)
		}
		if pageCount < 3 {
			return fmt.Errorf("expected at least 3 pages with Limit=1, got %d", pageCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeContinuousBackups", func() error {
		resp, err := client.DescribeContinuousBackups(ctx, &dynamodb.DescribeContinuousBackupsInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			return err
		}
		if resp.ContinuousBackupsDescription == nil {
			return fmt.Errorf("continuous backups description not found")
		}
		if resp.ContinuousBackupsDescription.PointInTimeRecoveryDescription == nil {
			return fmt.Errorf("expected PointInTimeRecoveryDescription")
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "UpdateContinuousBackups", func() error {
		resp, err := client.UpdateContinuousBackups(ctx, &dynamodb.UpdateContinuousBackupsInput{
			TableName: aws.String(tableName),
			PointInTimeRecoverySpecification: &types.PointInTimeRecoverySpecification{
				PointInTimeRecoveryEnabled: aws.Bool(true),
			},
		})
		if err != nil {
			return err
		}
		if resp.ContinuousBackupsDescription == nil {
			return fmt.Errorf("ContinuousBackupsDescription is nil")
		}
		if resp.ContinuousBackupsDescription.PointInTimeRecoveryDescription == nil {
			return fmt.Errorf("expected PointInTimeRecoveryDescription")
		}
		if resp.ContinuousBackupsDescription.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus != types.PointInTimeRecoveryStatusEnabled {
			return fmt.Errorf("expected PointInTimeRecoveryStatus=ENABLED, got %v", resp.ContinuousBackupsDescription.PointInTimeRecoveryDescription.PointInTimeRecoveryStatus)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBNonExistentTableTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "GetItem_NonExistentTable", func() error {
		_, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String("NoSuchTable_xyz"),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		return expectResourceNotFound(err)
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_NonExistentTable", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("NoSuchTable_xyz"),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		return expectResourceNotFound(err)
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeTable_NonExistentTable", func() error {
		_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		return expectResourceNotFound(err)
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteTable_NonExistentTable", func() error {
		_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		return expectResourceNotFound(err)
	}))

	results = append(results, r.RunTest("dynamodb", "Query_NonExistentTable", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String("NoSuchTable_xyz"),
			KeyConditionExpression: aws.String("id = :id"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		return expectResourceNotFound(err)
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_NonExistentTable", func() error {
		_, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		return expectResourceNotFound(err)
	}))

	return results
}

func (r *TestRunner) dynamoDBConditionalCheckTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "UpdateItem_ConditionalCheckFail", func() error {
		errTable := fmt.Sprintf("CondTable-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, errTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(errTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "cond1"},
				"status": &types.AttributeValueMemberS{Value: "active"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		_, err = client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: aws.String(errTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "cond1"},
			},
			UpdateExpression:    aws.String("SET #s = :val"),
			ConditionExpression: aws.String("#s = :expected"),
			ExpressionAttributeNames: map[string]string{
				"#s": "status",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":val":      &types.AttributeValueMemberS{Value: "inactive"},
				":expected": &types.AttributeValueMemberS{Value: "deleted"},
			},
		})
		return expectConditionalCheckFailed(err)
	}))

	results = append(results, r.RunTest("dynamodb", "GetItem_NonExistentKey", func() error {
		errTable := fmt.Sprintf("GetItemErr-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, errTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		resp, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(errTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "nonexistent"},
			},
		})
		if err != nil {
			return fmt.Errorf("GetItem on non-existent key should not error, got: %v", err)
		}
		if len(resp.Item) != 0 {
			return fmt.Errorf("expected empty item for non-existent key, got %d attributes", len(resp.Item))
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteItem_ConditionalCheckFail", func() error {
		errTable := fmt.Sprintf("DelCondTable-%d", time.Now().UnixNano())
		cleanupTable, err := createDynamoTestTable(ctx, client, errTable)
		if err != nil {
			return err
		}
		defer cleanupTable()

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(errTable),
			Item: map[string]types.AttributeValue{
				"id":     &types.AttributeValueMemberS{Value: "del1"},
				"status": &types.AttributeValueMemberS{Value: "active"},
			},
		})
		if err != nil {
			return fmt.Errorf("put item: %v", err)
		}

		_, err = client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(errTable),
			Key: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "del1"},
			},
			ConditionExpression: aws.String("attribute_not_exists(id)"),
		})
		return expectConditionalCheckFailed(err)
	}))

	return results
}
