package iot

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// ---- Topic Rule Destinations --------------------------------------
// AWS identifies destinations by ARN (auto-generated at create time, derived
// from a UUID) and resolves Confirm by confirmationToken. The earlier
// "destinationName" keying was a misreading of the Smithy model: neither
// CreateTopicRuleDestinationRequest nor any other request shape carries a
// destinationName member. The handlers below persist by ARN and produce the
// canonical TopicRuleDestination response shape.

func (s *IoTService) CreateTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.createTopicRuleDestinationCore(store, request.GetMapParamCaseInsensitive(req.Parameters, "destinationConfiguration"))
}

func (s *IoTService) DeleteTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteTopicRuleDestinationCore(store, request.GetParamCaseInsensitive(req.Parameters, "arn")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) GetTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.getTopicRuleDestinationCore(store, request.GetParamCaseInsensitive(req.Parameters, "arn"))
}

func (s *IoTService) ListTopicRuleDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	summaries, err := s.listTopicRuleDestinationsCore(store)
	if err != nil {
		return nil, err
	}
	return paginatedMaps("destinationSummaries", summaries, req.Parameters)
}

func (s *IoTService) UpdateTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.updateTopicRuleDestinationCore(store,
		request.GetParamCaseInsensitive(req.Parameters, "arn"),
		request.GetParamCaseInsensitive(req.Parameters, "status")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

func (s *IoTService) ConfirmTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.confirmTopicRuleDestinationCore(store, request.GetParamCaseInsensitive(req.Parameters, "confirmationToken")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
