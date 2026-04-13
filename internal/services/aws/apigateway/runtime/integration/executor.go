// Package integration provides API Gateway integration functionality for vorpalstacks.
package integration

import (
	"context"
	"net/http"

	"vorpalstacks/internal/eventbus"
	sns "vorpalstacks/internal/store/aws/sns"
	sqs "vorpalstacks/internal/store/aws/sqs"
)

// LambdaInvoker invokes Lambda functions for API Gateway integrations.
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}

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
	lambdaInvoker LambdaInvoker
	sqsStore      sqs.SQSStoreInterface
	snsStore      sns.SNSStoreInterface
	accountID     string
	region        string
	bus           *eventbus.EventBus
}

// NewExecutorFactory creates a new executor factory.
func NewExecutorFactory(lambdaInvoker LambdaInvoker) *ExecutorFactory {
	return &ExecutorFactory{
		lambdaInvoker: lambdaInvoker,
	}
}

// NewExecutorFactoryWithStores creates a new executor factory with SQS and SNS stores.
func NewExecutorFactoryWithStores(lambdaInvoker LambdaInvoker, sqsStore sqs.SQSStoreInterface, snsStore sns.SNSStoreInterface, accountID, region string) *ExecutorFactory {
	return &ExecutorFactory{
		lambdaInvoker: lambdaInvoker,
		sqsStore:      sqsStore,
		snsStore:      snsStore,
		accountID:     accountID,
		region:        region,
	}
}

// SetEventBus sets the event bus for SNS fan-out bug fix.
func (f *ExecutorFactory) SetEventBus(bus *eventbus.EventBus) {
	f.bus = bus
}

// SetSQSStore sets the SQS store for SQS integration targets.
func (f *ExecutorFactory) SetSQSStore(store sqs.SQSStoreInterface) {
	f.sqsStore = store
}

// SetSNSStore sets the SNS store for SNS integration targets.
func (f *ExecutorFactory) SetSNSStore(store sns.SNSStoreInterface) {
	f.snsStore = store
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
		return NewAWSExecutorWithStores(f.lambdaInvoker, f.sqsStore, f.snsStore, f.accountID, f.region, f.bus), nil
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
