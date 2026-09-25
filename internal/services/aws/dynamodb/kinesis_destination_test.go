// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestDescribeKinesisStreamingDestinationReportsDestinations(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "DescTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	resp, err := svc.DescribeKinesisStreamingDestination(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "DescTable",
	}})
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	describeResp, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("describe response shape: %#v", resp)
	}
	if describeResp["TableName"] != "DescTable" {
		t.Errorf("TableName: %v", describeResp["TableName"])
	}
	destinations, ok := describeResp["KinesisDataStreamDestinations"].([]map[string]interface{})
	if !ok || len(destinations) != 0 {
		t.Errorf("KinesisDataStreamDestinations: %#v", describeResp["KinesisDataStreamDestinations"])
	}
}
