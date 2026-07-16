package iotevents

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iotutil"
	"vorpalstacks/internal/eventbus"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// iotEventsRecursionTimeout bounds how long a recursive iotEvents action
// chain is allowed to run before it is cancelled. This prevents infinite loops
// when two detector models trigger each other.
const iotEventsRecursionTimeout = 5 * time.Second

type DetectorActionAdapter struct {
	bus           eventbus.Bus
	logger        *slog.Logger
	batchEvaluate func(ctx context.Context, messages []iotstore.InputMessage) []map[string]interface{}
	republishFn   func(topic string, payload []byte) error
}

func NewDetectorActionAdapter(bus eventbus.Bus, logger *slog.Logger) *DetectorActionAdapter {
	if logger == nil {
		logger = slog.Default()
	}
	return &DetectorActionAdapter{bus: bus, logger: logger}
}

func (a *DetectorActionAdapter) SetBatchEvaluate(fn func(ctx context.Context, messages []iotstore.InputMessage) []map[string]interface{}) {
	a.batchEvaluate = fn
}

func (a *DetectorActionAdapter) SetRepublishFn(fn func(topic string, payload []byte) error) {
	a.republishFn = fn
}

func (a *DetectorActionAdapter) OnAction() func(string, string, string, map[string]interface{}) {
	return func(modelName, key, actionType string, payload map[string]interface{}) {
		if a.bus == nil {
			a.logger.Warn("detector action: eventbus not available", "model", modelName, "actionType", actionType)
			return
		}
		if err := a.dispatch(context.Background(), modelName, key, actionType, payload); err != nil {
			a.logger.Warn("detector action dispatch failed", "model", modelName, "key", key, "actionType", actionType, "error", err)
		}
	}
}

func (a *DetectorActionAdapter) dispatch(ctx context.Context, modelName, key, actionType string, params map[string]interface{}) error {
	payloadJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal detector action payload: %w", err)
	}

	switch actionType {
	case "iotTopicPublish":
		return a.dispatchRepublish(modelName, params, payloadJSON)
	case "sqs":
		return a.dispatchSQS(ctx, params, payloadJSON)
	case "sns":
		return a.dispatchSNS(ctx, params, payloadJSON)
	case "lambda":
		return a.dispatchLambda(ctx, params, payloadJSON)
	case "firehose":
		return fmt.Errorf("firehose action not available in vorpalstacks")
	case "iotEvents":
		return a.dispatchIoTEvents(ctx, params)
	default:
		return awserrors.NewValidationException(fmt.Sprintf("unsupported detector action type: %s", actionType))
	}
}

func (a *DetectorActionAdapter) dispatchSQS(ctx context.Context, params map[string]interface{}, payloadJSON []byte) error {
	invoker := a.bus.SQSInvoker()
	if invoker == nil {
		return fmt.Errorf("sqs invoker not available")
	}
	queueURL := iotutil.StrFromMap(params, "queueUrl")
	if queueURL == "" {
		queueName := iotutil.StrFromMap(params, "queueName")
		if queueName == "" {
			return fmt.Errorf("sqs: no queueUrl or queueName specified")
		}
		resolved, err := invoker.GetQueueByName(ctx, queueName)
		if err != nil {
			return fmt.Errorf("sqs: failed to resolve queue %q: %w", queueName, err)
		}
		queueURL = resolved
	}
	_, _, err := invoker.SendMessage(ctx, queueURL, string(payloadJSON), eventbus.SQSSendOptions{})
	return err
}

func (a *DetectorActionAdapter) dispatchSNS(ctx context.Context, params map[string]interface{}, payloadJSON []byte) error {
	invoker := a.bus.SNSInvoker()
	if invoker == nil {
		return fmt.Errorf("sns invoker not available")
	}
	topicARN := iotutil.StrFromMap(params, "topicArn")
	if topicARN == "" {
		topicARN = iotutil.StrFromMap(params, "targetArn")
	}
	if topicARN == "" {
		return fmt.Errorf("sns: no topicArn or targetArn specified")
	}
	_, err := invoker.PublishToTopic(ctx, topicARN, string(payloadJSON), "", nil)
	return err
}

func (a *DetectorActionAdapter) dispatchLambda(ctx context.Context, params map[string]interface{}, payloadJSON []byte) error {
	invoker := a.bus.LambdaInvoker()
	if invoker == nil {
		return fmt.Errorf("lambda invoker not available")
	}
	functionName := iotutil.StrFromMap(params, "functionArn")
	if functionName == "" {
		functionName = iotutil.StrFromMap(params, "functionName")
	}
	if functionName == "" {
		return fmt.Errorf("lambda: no functionArn or functionName specified")
	}
	_, _, err := invoker.InvokeForGateway(ctx, functionName, payloadJSON)
	return err
}

func (a *DetectorActionAdapter) dispatchRepublish(modelName string, params map[string]interface{}, payloadJSON []byte) error {
	if a.republishFn == nil {
		return fmt.Errorf("iotTopicPublish: republish not configured (model %s)", modelName)
	}
	topic := iotutil.StrFromMap(params, "mqttTopic", "topic")
	if topic == "" {
		return fmt.Errorf("iotTopicPublish: no topic specified")
	}
	return a.republishFn(topic, payloadJSON)
}

func (a *DetectorActionAdapter) dispatchIoTEvents(_ context.Context, params map[string]interface{}) error {
	if a.batchEvaluate == nil {
		return fmt.Errorf("iotEvents: BatchEvaluate not configured")
	}
	inputName := iotutil.StrFromMap(params, "inputName")
	if inputName == "" {
		return fmt.Errorf("iotEvents: no inputName specified")
	}

	// Run in a goroutine to prevent unbounded recursive call stack growth
	// when detector models trigger each other. By the time this callback
	// fires, EvaluateEvent has already released sm.mu (actions are
	// dispatched outside the lock), so there is no deadlock risk.
	// Fire-and-forget with a 5s timeout bounds the recursion depth.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.logger.Error("detector action iotEvents: panic in recursive BatchEvaluate",
					"inputName", inputName, "panic", r)
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), iotEventsRecursionTimeout)
		defer cancel()

		payload := make(map[string]interface{})
		for k, v := range params {
			if k != "inputName" {
				payload[k] = v
			}
		}

		msg := iotstore.InputMessage{InputName: inputName, Payload: payload}
		errs := a.batchEvaluate(ctx, []iotstore.InputMessage{msg})
		if len(errs) > 0 {
			a.logger.Warn("detector action iotEvents: recursive evaluation errors", "inputName", inputName, "errors", len(errs))
			return
		}
		a.logger.Info("detector action iotEvents: recursive BatchPutMessage dispatched", "inputName", inputName)
	}()

	return nil
}
