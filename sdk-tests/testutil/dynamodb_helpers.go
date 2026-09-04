package testutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// createDynamoTestTable creates a throwaway hash-only on-demand table and
// returns a cleanup closure deleting it. The default shape is a single
// string "id" hash key under PAY_PER_REQUEST billing, which is the fixture
// most item-level tests need; opts may override the input (for example a
// different hash-key attribute or extra members such as deletion
// protection). Creation is wait-free: the service returns tables ACTIVE,
// and flows that need an explicit active-wait poll separately.
//
// Tables whose create input is itself the scenario — GSI/LSI/stream/SSE
// specifications, composite key schemas, negative validation paths — keep
// their inline CreateTable calls so the exercised input stays visible.
func createDynamoTestTable(ctx context.Context, client *dynamodb.Client, name string, opts ...func(*dynamodb.CreateTableInput)) (func(), error) {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := client.CreateTable(ctx, input); err != nil {
		return func() {}, fmt.Errorf("create table %s: %w", name, err)
	}
	return func() {
		_, _ = client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: aws.String(name)})
	}, nil
}

// withDynamoHashKey replaces the default key with a single string hash key
// of the given attribute name, for fixtures whose items address the table
// through an attribute other than "id".
func withDynamoHashKey(attr string) func(*dynamodb.CreateTableInput) {
	return func(input *dynamodb.CreateTableInput) {
		input.AttributeDefinitions = []types.AttributeDefinition{
			{AttributeName: aws.String(attr), AttributeType: types.ScalarAttributeTypeS},
		}
		input.KeySchema = []types.KeySchemaElement{
			{AttributeName: aws.String(attr), KeyType: types.KeyTypeHash},
		}
	}
}

// expectResourceNotFound asserts that err is a DynamoDB
// ResourceNotFoundException.
func expectResourceNotFound(err error) error {
	if err == nil {
		return fmt.Errorf("expected ResourceNotFoundException")
	}
	var rnf *types.ResourceNotFoundException
	if !errors.As(err, &rnf) {
		return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
	}
	return nil
}

// expectConditionalCheckFailed asserts that err is a DynamoDB
// ConditionalCheckFailedException.
func expectConditionalCheckFailed(err error) error {
	if err == nil {
		return fmt.Errorf("expected ConditionalCheckFailedException")
	}
	var ccf *types.ConditionalCheckFailedException
	if !errors.As(err, &ccf) {
		return fmt.Errorf("expected ConditionalCheckFailedException, got: %T: %v", err, err)
	}
	return nil
}
