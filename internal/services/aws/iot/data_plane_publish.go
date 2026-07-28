package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// Publish implements the iotdataplane Publish operation (POST /topics/{topic}).
// It publishes the message to the regional MQTT broker, which delivers to MQTT
// subscribers and triggers TopicRule SQL evaluation via the OnPublished hook.
// When the broker is unavailable (e.g. failed to start in TEST_MODE), it falls
// back to feeding the rule executor directly.
func (s *IoTService) Publish(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if s.throttle != nil && !s.throttle.allow() {
		return nil, iotstore.ErrThrottling
	}
	topic := request.GetParamCaseInsensitive(req.Parameters, "topic")
	if topic == "" {
		return nil, nil
	}
	// The payload arrives as the raw request body; the framework stores it
	// under the "body" key for operations without a structured input shape.
	payload := []byte{}
	if body, ok := req.Parameters["body"].(string); ok && body != "" {
		payload = []byte(body)
	} else if raw, ok := req.Parameters["payload"].(string); ok && raw != "" {
		payload = []byte(raw)
	}

	// Publish to the broker for MQTT subscriber delivery, and trigger the
	// rule executor directly. The broker's OnPublished hook only fires for
	// client-originated publishes, not server.Publish calls, so we must
	// invoke the executor explicitly for rule evaluation.
	if brk := s.brokerForReq(reqCtx); brk != nil {
		_ = brk.Publish(topic, payload)
	}
	executor := s.executorForReq(reqCtx)
	if executor != nil {
		executor.OnMessage(topic, payload)
	}

	return nil, nil
}
