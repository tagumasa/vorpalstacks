package testutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// createDynamoTestTable creates a throwaway on-demand table and returns a
// cleanup closure deleting it. The default shape is a single string "id"
// hash key under PAY_PER_REQUEST billing, which is the fixture most
// item-level tests need; opts override the input with declarative members
// (key schema, secondary indexes, stream specification, deletion
// protection) so every fixture table shares one creation, error-wrapping
// and cleanup path. Creation is wait-free: the service returns tables
// ACTIVE, and flows that need an explicit active-wait poll separately.
//
// Inline CreateTable calls remain only where the create input is itself
// the assertion target: the CreateTable operation tests, negative
// validation paths, and a fixture shared by a whole test family (one
// definition consumed by many tests).
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

// withDynamoKeySchema replaces the default single-string-key shape with the
// given attribute definitions and key schema, for composite keys and
// non-string (N/B) key types. The schema stays declarative at the call
// site because the key shape is part of what the test exercises.
func withDynamoKeySchema(defs []types.AttributeDefinition, schema []types.KeySchemaElement) func(*dynamodb.CreateTableInput) {
	return func(input *dynamodb.CreateTableInput) {
		input.AttributeDefinitions = defs
		input.KeySchema = schema
	}
}

// withDynamoGSI appends global secondary indexes to the fixture table.
func withDynamoGSI(indexes ...types.GlobalSecondaryIndex) func(*dynamodb.CreateTableInput) {
	return func(input *dynamodb.CreateTableInput) {
		input.GlobalSecondaryIndexes = append(input.GlobalSecondaryIndexes, indexes...)
	}
}

// withDynamoLSI appends local secondary indexes to the fixture table.
func withDynamoLSI(indexes ...types.LocalSecondaryIndex) func(*dynamodb.CreateTableInput) {
	return func(input *dynamodb.CreateTableInput) {
		input.LocalSecondaryIndexes = append(input.LocalSecondaryIndexes, indexes...)
	}
}

// withDynamoStream enables a stream with the given view type on the
// fixture table.
func withDynamoStream(view types.StreamViewType) func(*dynamodb.CreateTableInput) {
	return func(input *dynamodb.CreateTableInput) {
		input.StreamSpecification = &types.StreamSpecification{
			StreamEnabled:  aws.Bool(true),
			StreamViewType: view,
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
