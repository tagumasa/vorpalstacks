package iot

import (
	"time"

	"github.com/google/uuid"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Topic-rule-destination Core.
// AWS identifies destinations by ARN (auto-generated at create time, derived
// from a UUID) and resolves Confirm by confirmationToken. The earlier
// "destinationName" keying was a misreading of the Smithy model: neither
// CreateTopicRuleDestinationRequest nor any other request shape carries a
// destinationName member. Records are persisted by ARN under
// "topicRuleDestination/<arn>" with a "topicRuleDestinationToken/<token>"
// confirmationToken index.
// ---------------------------------------------------------------------------

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

// createTopicRuleDestinationCore validates and persists a destination with
// its confirmationToken index.
func (s *IoTService) createTopicRuleDestinationCore(store iotstore.IotStoreInterface, cfg map[string]interface{}) (map[string]interface{}, error) {
	if len(cfg) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	// AWS assigns a UUID-derived identifier and builds the ARN from it.
	// destinationName is not part of the AWS API.
	uid := uuid.New().String()
	arn := iotstore.BuildTopicRuleDestinationARN(store.GetAccountID(), store.GetRegion(), uid)
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

// deleteTopicRuleDestinationCore removes a destination record and its
// confirmationToken index.
func (s *IoTService) deleteTopicRuleDestinationCore(store iotstore.IotStoreInterface, arn string) error {
	if arn == "" {
		return iotstore.ErrMissingParam
	}
	key := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrTopicRuleDestinationNotFound
	}
	if err := store.DeleteGeneric(key); err != nil {
		return err
	}
	if token, ok := rec["confirmationToken"].(string); ok && token != "" {
		if err := store.DeleteGeneric("topicRuleDestinationToken/" + token); err != nil {
			return err
		}
	}
	return nil
}

// getTopicRuleDestinationCore retrieves a destination wrapped in the
// topicRuleDestination response member.
func (s *IoTService) getTopicRuleDestinationCore(store iotstore.IotStoreInterface, arn string) (map[string]interface{}, error) {
	if arn == "" {
		return nil, iotstore.ErrMissingParam
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

// listTopicRuleDestinationsCore lists destination summaries.
func (s *IoTService) listTopicRuleDestinationsCore(store iotstore.IotStoreInterface) ([]map[string]interface{}, error) {
	items, err := store.ListGeneric("topicRuleDestination/")
	if err != nil {
		return nil, err
	}
	summaries := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		summaries = append(summaries, topicRuleDestinationResponse(rec))
	}
	return summaries, nil
}

// topicRuleDestinationStatuses is the documented TopicRuleDestinationStatus
// enum shared by the destination output shape and UpdateTopicRuleDestination.
var topicRuleDestinationStatuses = map[string]bool{
	"ENABLED":     true,
	"IN_PROGRESS": true,
	"DISABLED":    true,
	"ERROR":       true,
	"DELETING":    true,
}

// updateTopicRuleDestinationCore applies a status change to a destination.
// The status must be a member of the documented enum; anything else is an
// InvalidRequestException.
func (s *IoTService) updateTopicRuleDestinationCore(store iotstore.IotStoreInterface, arn, status string) error {
	if arn == "" || status == "" {
		return iotstore.ErrMissingParam
	}
	if !topicRuleDestinationStatuses[status] {
		return iotstore.ErrInvalidRequest
	}
	key := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrTopicRuleDestinationNotFound
	}
	rec["status"] = status
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	return store.PutGeneric(key, rec)
}

// confirmTopicRuleDestinationCore resolves a confirmationToken and flips
// the destination status to ENABLED.
func (s *IoTService) confirmTopicRuleDestinationCore(store iotstore.IotStoreInterface, token string) error {
	if token == "" {
		return iotstore.ErrMissingParam
	}
	// Resolve token -> arn, then flip the destination status to ENABLED.
	tokenRec := map[string]interface{}{}
	tokenExists, err := store.GetGenericExists("topicRuleDestinationToken/"+token, &tokenRec)
	if err != nil {
		return err
	}
	if !tokenExists {
		return iotstore.ErrTopicRuleDestinationNotFound
	}
	arn, _ := tokenRec["arn"].(string)
	if arn == "" {
		return iotstore.ErrTopicRuleDestinationNotFound
	}
	destKey := "topicRuleDestination/" + arn
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(destKey, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrTopicRuleDestinationNotFound
	}
	rec["status"] = "ENABLED"
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	return store.PutGeneric(destKey, rec)
}
