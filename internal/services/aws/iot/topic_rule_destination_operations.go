package iot

import (
	"context"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---- Topic Rule Destinations --------------------------------------
// AWS identifies destinations by ARN (auto-generated at create time, derived
// from a UUID) and resolves Confirm by confirmationToken. The earlier
// "destinationName" keying was a misreading of the Smithy model: neither
// CreateTopicRuleDestinationRequest nor any other request shape carries a
// destinationName member. The handlers below persist by ARN and produce the
// canonical TopicRuleDestination response shape.

// topicRuleDestinationResponse shapes a stored record into the AWS
// TopicRuleDestination output structure. The configuration sub-map is echoed
// verbatim so callers see the httpUrlProperties / vpcProperties they supplied.
func topicRuleDestinationResponse(rec map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{
		"arn":          rec["arn"],
		"status":       rec["status"],
		"statusReason": rec["statusReason"],
	}
	if cfg, ok := rec["destinationConfiguration"].(map[string]interface{}); ok {
		if http, ok := cfg["httpUrlConfiguration"].(map[string]interface{}); ok && len(http) > 0 {
			out["httpUrlProperties"] = http
		}
		if vpc, ok := cfg["vpcConfiguration"].(map[string]interface{}); ok && len(vpc) > 0 {
			out["vpcProperties"] = vpc
		}
	}
	if v, ok := rec["createdAt"]; ok {
		out["createdAt"] = v
	}
	if v, ok := rec["lastUpdatedAt"]; ok {
		out["lastUpdatedAt"] = v
	}
	return out
}

func (s *IoTService) CreateTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	cfg := request.GetMapParamCaseInsensitive(req.Parameters, "destinationConfiguration")
	if len(cfg) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// AWS assigns a UUID-derived identifier and builds the ARN from it.
	// destinationName is not part of the AWS API.
	uid := uuid.New().String()
	arn := iotstore.BuildTopicRuleDestinationARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), uid)
	confirmationToken := uuid.New().String()
	now := time.Now().UTC().Unix()
	rec := map[string]interface{}{
		"arn":                      arn,
		"status":                   "IN_PROGRESS",
		"confirmationToken":        confirmationToken,
		"destinationConfiguration": cfg,
		"createdAt":                now,
		"lastUpdatedAt":            now,
	}
	if err := store.PutGeneric("topicRuleDestination/"+arn, rec); err != nil {
		return nil, err
	}
	// Also index by confirmationToken so ConfirmTopicRuleDestination can look
	// the record up by token alone.
	if err := store.PutGeneric("topicRuleDestinationToken/"+confirmationToken, map[string]interface{}{"arn": arn}); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"topicRuleDestination": topicRuleDestinationResponse(rec),
	}, nil
}
func (s *IoTService) DeleteTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "arn")
	if arn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return nil, err
	}
	if token, ok := rec["confirmationToken"].(string); ok && token != "" {
		if err := store.DeleteGeneric("topicRuleDestinationToken/" + token); err != nil {
			return nil, err
		}
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) GetTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "arn")
	if arn == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("topicRuleDestination/"+arn, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	return map[string]interface{}{
		"topicRuleDestination": topicRuleDestinationResponse(rec),
	}, nil
}
func (s *IoTService) ListTopicRuleDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("topicRuleDestination/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		summaries = append(summaries, topicRuleDestinationResponse(rec))
	}
	return paginatedMaps("destinationSummaries", summaries, req.Parameters), nil
}
func (s *IoTService) UpdateTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetParamCaseInsensitive(req.Parameters, "arn")
	status := request.GetParamCaseInsensitive(req.Parameters, "status")
	if arn == "" || status == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	rec["status"] = status
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
func (s *IoTService) ConfirmTopicRuleDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	token := request.GetParamCaseInsensitive(req.Parameters, "confirmationToken")
	if token == "" {
		return nil, iotstore.ErrMissingParam
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// Resolve token -> arn, then flip the destination status to ENABLED.
	tokenRec := map[string]interface{}{}
	tokenExists, err := store.GetGenericExists("topicRuleDestinationToken/"+token, &tokenRec)
	if err != nil {
		return nil, err
	}
	if !tokenExists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	arn, _ := tokenRec["arn"].(string)
	if arn == "" {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	destKey := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(destKey, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrTopicRuleDestinationNotFound
	}
	rec["status"] = "ENABLED"
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(destKey, rec); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}
