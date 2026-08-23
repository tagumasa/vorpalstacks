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
	// The restjson1 Publish operation carries the payload blob as the raw
	// request body (the httpPayload trait): the generic body parser merges
	// a JSON object body into req.Parameters field by field, so the payload
	// must be recovered from req.Body. Wrapped "body"/"payload" string
	// parameters keep working for callers that send an envelope.
	payload := []byte{}
	if body, ok := req.Parameters["body"].(string); ok && body != "" {
		payload = []byte(body)
	} else if raw, ok := req.Parameters["payload"].(string); ok && raw != "" {
		payload = []byte(raw)
	} else if len(req.Body) > 0 {
		payload = req.Body
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
