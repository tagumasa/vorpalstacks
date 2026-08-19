package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
)

// dynamoDBInputValidationTests pins input-validation contracts that the
// service enforces server-side: TransactGetItems resolves every referenced
// table (a missing table is a ResourceNotFoundException and key types are
// checked against the table's attribute definitions), PAY_PER_REQUEST
// tables cannot carry ProvisionedThroughput, and a configured
// WarmThroughput is echoed by table descriptions while an unconfigured one
// is omitted.
func (r *TestRunner) dynamoDBInputValidationTests(ctx context.Context, client *dynamodb.Client) []TestResult {
	var results []TestResult
	suffix := time.Now().UnixNano()

	tableName := fmt.Sprintf("input-validate-%d", suffix)
	if err := createTTLTestTable(ctx, client, tableName); err != nil {
		return append(results, TestResult{
			Service:  "dynamodb",
			TestName: "TransactGetItems_MissingTable_ResourceNotFound",
			Status:   "FAIL",
			Error:    fmt.Sprintf("create table %s: %v", tableName, err),
		})
	}
	defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(tableName)})

	results = append(results, r.RunTest("dynamodb", "TransactGetItems_MissingTable_ResourceNotFound", func() error {
		_, err := client.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
			TransactItems: []dynamodbtypes.TransactGetItem{{
				Get: &dynamodbtypes.Get{
					TableName: aws.String(fmt.Sprintf("no-such-table-%d", suffix)),
					Key: map[string]dynamodbtypes.AttributeValue{
						"pk": &dynamodbtypes.AttributeValueMemberS{Value: "k"},
					},
				},
			}},
		})
		if err == nil {
			return fmt.Errorf("expected ResourceNotFoundException for a missing table")
		}
		var notFound *dynamodbtypes.ResourceNotFoundException
		if !errors.As(err, &notFound) {
			return fmt.Errorf("expected ResourceNotFoundException, got %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TransactGetItems_KeyTypeMismatch_Rejected", func() error {
		_, err := client.TransactGetItems(ctx, &dynamodb.TransactGetItemsInput{
			TransactItems: []dynamodbtypes.TransactGetItem{{
				Get: &dynamodbtypes.Get{
					TableName: aws.String(tableName),
					Key: map[string]dynamodbtypes.AttributeValue{
						"pk": &dynamodbtypes.AttributeValueMemberN{Value: "1"},
					},
				},
			}},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for a numeric key on a string table")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected APIError, got %v", err)
		}
		if apiErr.ErrorCode() != "ValidationException" {
			return fmt.Errorf("expected ValidationException, got %s", apiErr.ErrorCode())
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "CreateTable_PayPerRequestWithProvisionedThroughput_Rejected", func() error {
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(fmt.Sprintf("ppr-throughput-%d", suffix)),
			AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			},
			KeySchema: []dynamodbtypes.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
			},
			BillingMode: dynamodbtypes.BillingModePayPerRequest,
			ProvisionedThroughput: &dynamodbtypes.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
		if err == nil {
			return fmt.Errorf("expected ValidationException for ProvisionedThroughput on a PAY_PER_REQUEST table")
		}
		var apiErr smithy.APIError
		if !errors.As(err, &apiErr) {
			return fmt.Errorf("expected APIError, got %v", err)
		}
		if apiErr.ErrorCode() != "ValidationException" {
			return fmt.Errorf("expected ValidationException, got %s", apiErr.ErrorCode())
		}
		return nil
	}))

	results = append(results, r.RunTest("dynamodb", "TableDescription_WarmThroughputEchoedAndOmitted", func() error {
		warmName := fmt.Sprintf("warm-throughput-%d", suffix)
		_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			TableName: aws.String(warmName),
			AttributeDefinitions: []dynamodbtypes.AttributeDefinition{
				{AttributeName: aws.String("pk"), AttributeType: dynamodbtypes.ScalarAttributeTypeS},
			},
			KeySchema: []dynamodbtypes.KeySchemaElement{
				{AttributeName: aws.String("pk"), KeyType: dynamodbtypes.KeyTypeHash},
			},
			BillingMode: dynamodbtypes.BillingModeProvisioned,
			WarmThroughput: &dynamodbtypes.WarmThroughput{
				ReadUnitsPerSecond:  aws.Int64(3000),
				WriteUnitsPerSecond: aws.Int64(1000),
			},
			ProvisionedThroughput: &dynamodbtypes.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})
		if err != nil {
			return err
		}
		defer client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(warmName)})

		warmDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(warmName)})
		if err != nil {
			return err
		}
		if warmDesc.Table.WarmThroughput == nil {
			return fmt.Errorf("expected WarmThroughput on a table created with one")
		}
		if aws.ToInt64(warmDesc.Table.WarmThroughput.ReadUnitsPerSecond) != 3000 {
			return fmt.Errorf("expected ReadUnitsPerSecond=3000, got %v", warmDesc.Table.WarmThroughput.ReadUnitsPerSecond)
		}
		if aws.ToInt64(warmDesc.Table.WarmThroughput.WriteUnitsPerSecond) != 1000 {
			return fmt.Errorf("expected WriteUnitsPerSecond=1000, got %v", warmDesc.Table.WarmThroughput.WriteUnitsPerSecond)
		}

		plainDesc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(tableName)})
		if err != nil {
			return err
		}
		if plainDesc.Table.WarmThroughput != nil {
			return fmt.Errorf("expected no WarmThroughput on a table created without one, got %v", plainDesc.Table.WarmThroughput)
		}
		return nil
	}))

	return results
}
