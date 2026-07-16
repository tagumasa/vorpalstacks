package iotevents

import (
	"context"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// BatchPutMessage evaluates a batch of input messages against all loaded
// detector models in the state machine. This is a thin handler that parses
// the request and delegates evaluation to the store layer.
func (s *IoTEventsService) BatchPutMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	concrete, ok := store.(*iotstore.IotStore)
	if !ok || concrete.StateMachine() == nil {
		return map[string]interface{}{
			"BatchPutMessageErrorEntries": []map[string]interface{}{},
		}, nil
	}

	messagesRaw, ok := req.Parameters["messages"]
	if !ok {
		return map[string]interface{}{
			"BatchPutMessageErrorEntries": []map[string]interface{}{},
		}, nil
	}

	messages, err := parseBatchMessages(messagesRaw)
	if err != nil {
		return nil, awserrors.NewValidationException(fmt.Sprintf("invalid messages parameter: %v", err))
	}

	errs := concrete.BatchEvaluate(ctx, messages)

	return map[string]interface{}{
		"BatchPutMessageErrorEntries": errs,
	}, nil
}

// parseBatchMessages extracts a slice of InputMessage from the raw messages
// parameter, which may be []interface{} or []map[string]interface{}.
func parseBatchMessages(raw interface{}) ([]iotstore.InputMessage, error) {
	switch m := raw.(type) {
	case []interface{}:
		messages := make([]iotstore.InputMessage, 0, len(m))
		for _, item := range m {
			if msg, ok := item.(map[string]interface{}); ok {
				inputName, _ := msg["inputName"].(string)
				payload := iotstore.ExtractPayload(msg["messagePayload"])
				messages = append(messages, iotstore.InputMessage{InputName: inputName, Payload: payload})
			}
		}
		return messages, nil
	case []map[string]interface{}:
		messages := make([]iotstore.InputMessage, 0, len(m))
		for _, msg := range m {
			inputName, _ := msg["inputName"].(string)
			payload := iotstore.ExtractPayload(msg["messagePayload"])
			messages = append(messages, iotstore.InputMessage{InputName: inputName, Payload: payload})
		}
		return messages, nil
	default:
		return nil, awserrors.NewValidationException("messages must be an array")
	}
}
