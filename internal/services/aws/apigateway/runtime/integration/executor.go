// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"context"
	"net/http"

	"vorpalstacks/internal/eventbus"
)

// IntegrationRequest represents an API Gateway integration request.
type IntegrationRequest struct {
	Method                      string
	URI                         string
	Headers                     map[string]string
	MultiValueHeaders           map[string][]string
	Body                        []byte
	QueryParams                 map[string]string
	MultiValueQueryStringParams map[string][]string
	PathParams                  map[string]string
	Path                        string
	StageVariables              map[string]string
	IntegrationType             string
	RequestParameters           map[string]string
	RequestTemplates            map[string]string
	RestApiId                   string
	StageName                   string
	SourceIP                    string
	IntegrationResponses        map[string]*IntegrationResponseConfig
}

// IntegrationResponseConfig represents the configuration for an integration response.
type IntegrationResponseConfig struct {
	StatusCode        string
	SelectionPattern  string
	ResponseHeaders   map[string]string
	ResponseTemplates map[string]string
	ContentHandling   string
}

// IntegrationResponse represents an API Gateway integration response.
type IntegrationResponse struct {
	StatusCode      int
	Headers         map[string]string
	Body            []byte
	IsBase64Encoded bool
}

// Executor executes API Gateway integrations.
type Executor interface {
	Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error)
}

// ExecutorFactory creates integration executors.
type ExecutorFactory struct {
	bus       eventbus.Bus
	accountID string
	region    string
}

// NewExecutorFactory creates a new executor factory.
func NewExecutorFactory(bus eventbus.Bus) *ExecutorFactory {
	return &ExecutorFactory{
		bus: bus,
	}
}

// SetEventBus sets the event bus for cross-service delivery.
func (f *ExecutorFactory) SetEventBus(bus eventbus.Bus) {
	f.bus = bus
}

// SetAccountAndRegion sets the account ID and region for ARN construction.
func (f *ExecutorFactory) SetAccountAndRegion(accountID, region string) {
	f.accountID = accountID
	f.region = region
}

// CreateExecutor creates an executor for the given integration type.
func (f *ExecutorFactory) CreateExecutor(integrationType string) (Executor, error) {
	switch integrationType {
	case "MOCK":
		return NewMockExecutor(), nil
	case "HTTP", "HTTP_PROXY":
		return NewHTTPExecutor(), nil
	case "AWS", "AWS_PROXY":
		return NewAWSExecutor(f.bus, f.accountID, f.region), nil
	default:
		return nil, &IntegrationError{
			Message:  "Unsupported integration type: " + integrationType,
			Type:     "BadRequestException",
			HTTPCode: http.StatusBadRequest,
		}
	}
}

// IntegrationError represents an error during integration execution.
type IntegrationError struct {
	Message  string
	Type     string
	HTTPCode int
}

// Error returns the error message for the integration error.
func (e *IntegrationError) Error() string {
	return e.Message
}
