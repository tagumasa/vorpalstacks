package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

func (s *IoTService) BatchPutMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	messagesRaw, ok := req.Parameters["messages"]
	if !ok {
		return map[string]interface{}{
			"BatchPutMessageErrorEntries": []map[string]interface{}{},
		}, nil
	}

	var messages []map[string]interface{}
	switch m := messagesRaw.(type) {
	case []interface{}:
		for _, item := range m {
			if msg, ok := item.(map[string]interface{}); ok {
				messages = append(messages, msg)
			}
		}
	case []map[string]interface{}:
		messages = m
	}

	if s.stateMachine == nil {
		return map[string]interface{}{
			"BatchPutMessageErrorEntries": []map[string]interface{}{},
		}, nil
	}

	for _, msg := range messages {
		inputName, _ := msg["inputName"].(string)

		payloadRaw := msg["payload"]
		var payload map[string]interface{}
		switch p := payloadRaw.(type) {
		case string:
			continue
		case map[string]interface{}:
			payload = p
		}

		if inputName == "" || payload == nil {
			continue
		}

		s.stateMachine.mu.RLock()
		models := make([]string, 0, len(s.stateMachine.models))
		for modelName := range s.stateMachine.models {
			models = append(models, modelName)
		}
		s.stateMachine.mu.RUnlock()

		for _, modelName := range models {
			_ = s.stateMachine.EvaluateEvent(ctx, modelName, inputName, payload)
		}
	}

	return map[string]interface{}{
		"BatchPutMessageErrorEntries": []map[string]interface{}{},
	}, nil
}
