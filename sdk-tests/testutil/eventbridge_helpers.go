package testutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

// createEventBridgeTestBus creates a throwaway event bus and returns a
// cleanup closure deleting it. Creation is wait-free: the suite uses buses
// immediately after creation today and the server returns them usable.
//
// Buses whose create input is itself the scenario — log-config round
// trips, out-of-enum rejections, the duplicate-name negative — keep their
// inline CreateEventBus calls so the exercised input stays visible.
func createEventBridgeTestBus(ctx context.Context, client *eventbridge.Client, name string, opts ...func(*eventbridge.CreateEventBusInput)) (func(), error) {
	input := &eventbridge.CreateEventBusInput{
		Name: aws.String(name),
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := client.CreateEventBus(ctx, input); err != nil {
		return func() {}, fmt.Errorf("create event bus %s: %w", name, err)
	}
	return func() {
		_, _ = client.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{Name: aws.String(name)})
	}, nil
}

// createEventBridgeTestRule puts a throwaway rule on a bus and returns a
// cleanup closure deleting it. The default input is a rate(1 hour)
// schedule; opts adapt rules whose input differs (extra descriptions,
// event patterns).
func createEventBridgeTestRule(ctx context.Context, client *eventbridge.Client, busName, name string, opts ...func(*eventbridge.PutRuleInput)) (func(), error) {
	input := &eventbridge.PutRuleInput{
		Name:               aws.String(name),
		EventBusName:       aws.String(busName),
		ScheduleExpression: aws.String("rate(1 hour)"),
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := client.PutRule(ctx, input); err != nil {
		return func() {}, fmt.Errorf("put rule %s: %w", name, err)
	}
	return func() {
		_, _ = client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(name), EventBusName: aws.String(busName)})
	}, nil
}

// createEventBridgeTestConnection creates a throwaway connection with
// basic auth and returns a cleanup closure deleting it; opts adapt the
// authorisation input (API key connections, other credentials). Tests
// that delete the connection as the operation under test call the helper
// without keeping the cleanup.
func createEventBridgeTestConnection(ctx context.Context, client *eventbridge.Client, name string, opts ...func(*eventbridge.CreateConnectionInput)) (func(), error) {
	input := &eventbridge.CreateConnectionInput{
		Name:              aws.String(name),
		AuthorizationType: types.ConnectionAuthorizationTypeBasic,
		AuthParameters: &types.CreateConnectionAuthRequestParameters{
			BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
				Username: aws.String("u"),
				Password: aws.String("p"),
			},
		},
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := client.CreateConnection(ctx, input); err != nil {
		return func() {}, fmt.Errorf("create connection %s: %w", name, err)
	}
	return func() {
		_, _ = client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(name)})
	}, nil
}

// createEventBridgeTestApiDestination creates a throwaway API destination
// on a connection and returns a cleanup closure deleting it. The default
// input is a POST endpoint at https://example.com/webhook; opts adapt the
// endpoint or description.
func createEventBridgeTestApiDestination(ctx context.Context, client *eventbridge.Client, name, connectionARN string, opts ...func(*eventbridge.CreateApiDestinationInput)) (func(), error) {
	input := &eventbridge.CreateApiDestinationInput{
		Name:               aws.String(name),
		ConnectionArn:      aws.String(connectionARN),
		HttpMethod:         types.ApiDestinationHttpMethodPost,
		InvocationEndpoint: aws.String("https://example.com/webhook"),
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := client.CreateApiDestination(ctx, input); err != nil {
		return func() {}, fmt.Errorf("create api destination %s: %w", name, err)
	}
	return func() {
		_, _ = client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(name)})
	}, nil
}

// expectEventBridgeResourceNotFound asserts that err is an EventBridge
// ResourceNotFoundException. The DynamoDB helper of the same shape asserts
// the DynamoDB exception type, so the two cannot share an implementation.
func expectEventBridgeResourceNotFound(err error) error {
	if err == nil {
		return fmt.Errorf("expected ResourceNotFoundException")
	}
	var rnf *types.ResourceNotFoundException
	if !errors.As(err, &rnf) {
		return fmt.Errorf("expected ResourceNotFoundException, got: %T: %v", err, err)
	}
	return nil
}
