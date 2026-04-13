package eventbus

import "context"

// ServiceInvoker defines the contract for invoking a target service as
// part of event dispatch.
type ServiceInvoker interface {
	Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult
	ServiceType() string
}

// LambdaInvokerAdapter adapts a LambdaInvoker to the ServiceInvoker interface.
type LambdaInvokerAdapter struct {
	Invoker LambdaInvoker
}

// LambdaInvoker invokes a Lambda function by name and returns the status
// code and response payload.
type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}

// Invoke dispatches a Lambda invocation via the underlying LambdaInvoker.
func (a *LambdaInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	statusCode, respPayload, err := a.Invoker.InvokeForGateway(ctx, action.Identifier, payload)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: statusCode, Payload: respPayload}
}

// ServiceType returns "lambda" for this adapter.
func (a *LambdaInvokerAdapter) ServiceType() string { return "lambda" }

// SQSInvokerAdapter adapts an SQSSender to the ServiceInvoker interface.
type SQSInvokerAdapter struct {
	Store SQSSender
}

// SQSSender sends a message to an SQS queue, returning the message ID and MD5.
type SQSSender interface {
	SendMessage(ctx context.Context, queueURL string, messageBody string, delaySeconds int64, messageAttributes map[string]string) (string, string, error)
}

// Invoke dispatches an SQS SendMessage via the underlying SQSSender.
func (a *SQSInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	messageID, _, err := a.Store.SendMessage(ctx, action.Identifier, string(payload), 0, nil)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: 200, Payload: []byte(messageID)}
}

// ServiceType returns "sqs" for this adapter.
func (a *SQSInvokerAdapter) ServiceType() string { return "sqs" }

// SNSInvokerAdapter adapts an SNSPublisher to the ServiceInvoker interface.
type SNSInvokerAdapter struct {
	Publisher SNSPublisher
}

// SNSPublisher publishes a message to an SNS topic, returning the message ID.
type SNSPublisher interface {
	PublishToTopic(ctx context.Context, topicArn string, message string, subject string, messageAttributes map[string]string) (string, error)
}

// Invoke dispatches an SNS PublishToTopic via the underlying SNSPublisher.
func (a *SNSInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	messageID, err := a.Publisher.PublishToTopic(ctx, action.Identifier, string(payload), "", nil)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: 200, Payload: []byte(messageID)}
}

// ServiceType returns "sns" for this adapter.
func (a *SNSInvokerAdapter) ServiceType() string { return "sns" }

// KinesisInvokerAdapter adapts a KinesisPutter to the ServiceInvoker interface.
type KinesisInvokerAdapter struct {
	Store KinesisPutter
}

// KinesisPutter puts a record into a Kinesis stream, returning the sequence number.
type KinesisPutter interface {
	PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) (string, error)
}

// Invoke dispatches a Kinesis PutRecord via the underlying KinesisPutter.
func (a *KinesisInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	sequenceNumber, err := a.Store.PutRecord(ctx, action.Identifier, "eventbus", payload)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: 200, Payload: []byte(sequenceNumber)}
}

// ServiceType returns "kinesis" for this adapter.
func (a *KinesisInvokerAdapter) ServiceType() string { return "kinesis" }
