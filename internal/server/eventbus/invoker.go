package eventbus

import "context"

type ServiceInvoker interface {
	Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult
	ServiceType() string
}

type LambdaInvokerAdapter struct {
	Invoker LambdaInvoker
}

type LambdaInvoker interface {
	InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error)
}

func (a *LambdaInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	statusCode, respPayload, err := a.Invoker.InvokeForGateway(ctx, action.Identifier, payload)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: statusCode, Payload: respPayload}
}

func (a *LambdaInvokerAdapter) ServiceType() string { return "lambda" }

type SQSInvokerAdapter struct {
	Store SQSSender
}

type SQSSender interface {
	SendMessage(ctx context.Context, queueURL string, messageBody string, delaySeconds int64, messageAttributes map[string]string) (string, string, error)
}

func (a *SQSInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	messageID, _, err := a.Store.SendMessage(ctx, action.Identifier, string(payload), 0, nil)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: 200, Payload: []byte(messageID)}
}

func (a *SQSInvokerAdapter) ServiceType() string { return "sqs" }

type SNSInvokerAdapter struct {
	Publisher SNSPublisher
}

type SNSPublisher interface {
	PublishToTopic(ctx context.Context, topicArn string, message string, subject string, messageAttributes map[string]string) (string, error)
}

func (a *SNSInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	messageID, err := a.Publisher.PublishToTopic(ctx, action.Identifier, string(payload), "", nil)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: 200, Payload: []byte(messageID)}
}

func (a *SNSInvokerAdapter) ServiceType() string { return "sns" }

type KinesisInvokerAdapter struct {
	Store KinesisPutter
}

type KinesisPutter interface {
	PutRecord(ctx context.Context, streamName string, partitionKey string, data []byte) (string, error)
}

func (a *KinesisInvokerAdapter) Invoke(ctx context.Context, action *TargetAction, payload []byte) HandlerResult {
	sequenceNumber, err := a.Store.PutRecord(ctx, action.Identifier, "eventbus", payload)
	if err != nil {
		return HandlerResult{Error: err}
	}
	return HandlerResult{StatusCode: 200, Payload: []byte(sequenceNumber)}
}

func (a *KinesisInvokerAdapter) ServiceType() string { return "kinesis" }
