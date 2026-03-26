// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"context"
)

// MockExecutor handles mock integrations for API Gateway.
type MockExecutor struct{}

// NewMockExecutor creates a new mock executor instance.
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{}
}

// Execute returns a mock response for API Gateway integration testing.
func (e *MockExecutor) Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	return &IntegrationResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            []byte(`{"message": "Mock integration response"}`),
		IsBase64Encoded: false,
	}, nil
}
