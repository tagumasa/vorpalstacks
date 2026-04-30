package testutil

import (
	"context"
	"errors"
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
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "PutItem_NonExistentTable", func() error {
		_, err := client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("NoSuchTable_xyz"),
			Item: map[string]types.AttributeValue{
				"id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DescribeTable_NonExistentTable", func() error {
		_, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "DeleteTable_NonExistentTable", func() error {
		_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Query_NonExistentTable", func() error {
		_, err := client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String("NoSuchTable_xyz"),
			KeyConditionExpression: aws.String("id = :id"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":id": &types.AttributeValueMemberS{Value: "k"},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "Scan_NonExistentTable", func() error {
		_, err := client.Scan(ctx, &dynamodb.ScanInput{
			TableName: aws.String("NoSuchTable_xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent table")
		}
		var rnf *types.ResourceNotFoundException
		if !errors.As(err, &rnf) {
			return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
		}
		return nil
	}))

	return results
}

func (r *TestRunner) dynamoDBConditionalCheckTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("dynamodb", "UpdateItem_ConditionalCheckFail", func() error {
		errTable := fmt.Sprintf("CondTable-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(errTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(errTable)})

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
		if err == nil {
			return fmt.Errorf("expected ConditionalCheckFailedException")
		}
		var ccf *types.ConditionalCheckFailedException
		if !errors.As(err, &ccf) {
			return fmt.Errorf("expected ConditionalCheckFailedException, got: %T: %v", err, err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "GetItem_NonExistentKey", func() error {
		errTable := fmt.Sprintf("GetItemErr-%d", time.Now().UnixNano())
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(errTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(errTable)})

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
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(errTable),
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
			},
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
			},
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			return fmt.Errorf("create table: %v", err)
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(errTable)})

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
		if err == nil {
			return fmt.Errorf("expected ConditionalCheckFailedException")
		}
		var ccf *types.ConditionalCheckFailedException
		if !errors.As(err, &ccf) {
			return fmt.Errorf("expected ConditionalCheckFailedException, got: %T: %v", err, err)
		}
		return nil
	}))

	return results
}
