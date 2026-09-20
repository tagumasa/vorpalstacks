package sns

import (
	"encoding/json"
	"fmt"
	"strings"

	"vorpalstacks/internal/core/logs"
	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateTopicInput carries every field that CreateTopic needs, in a format
// independent of the wire protocol (HTTP Query vs gRPC-Web). Both the HTTP
// API handler (topic_operations.go) and the admin gRPC handler
// (admin_handler.go) build this struct and delegate to createTopicCore,
// ensuring that name validation, FIFO consistency checks, attribute
// validation, and persistence follow a single code path.
type CreateTopicInput struct {
	Name       string
	Attributes map[string]string
	Tags       map[string]string

	// DataProtectionPolicy is the optional inline data protection policy
	// (CreateTopicInput.DataProtectionPolicy member). It is stored under the
	// internal "DataProtectionPolicy" attribute key, retrievable only via
	// GetDataProtectionPolicy, and excluded from GetTopicAttributes output.
	DataProtectionPolicy string
}

// TopicResult is the transport-agnostic result of creating or looking up a
// topic. The create response carries the topic's ARN only (the model's
// CreateTopicResult declares TopicArn as its sole member).
type TopicResult struct {
	Arn string
}

// ListTopicsInput carries the pagination token for ListTopics.
type ListTopicsInput struct {
	NextToken string
}

// TopicSummary is the transport-agnostic summary of a topic in list results.
type TopicSummary struct {
	TopicArn string
}

// ListTopicsResult is the transport-agnostic result of listing topics.
type ListTopicsResult struct {
	Topics    []TopicSummary
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createTopicCore is the single entry point for topic creation logic shared
// by the HTTP API and the admin gRPC handler. It performs all name
// validation, FIFO consistency checks, attribute validation (including
// value length caps), and persistence.
func (s *SNSService) createTopicCore(store snsstore.SNSStoreInterface, in CreateTopicInput) (*TopicResult, error) {
	if err := validateTopicName(in.Name); err != nil {
		return nil, err
	}

	isFifo := strings.HasSuffix(in.Name, ".fifo")

	if in.Attributes == nil {
		in.Attributes = make(map[string]string)
	}

	// FIFO consistency: respect the user-provided FifoTopic attribute
	// rather than silently overriding it. If the name has a .fifo suffix,
	// FifoTopic must be "true" (either explicitly or by default injection).
	// If the name lacks the suffix, FifoTopic must not be "true".
	if isFifo {
		userVal, hasAttr := in.Attributes[snsstore.AttrFifoTopic]
		if !hasAttr {
			in.Attributes[snsstore.AttrFifoTopic] = "true"
		} else if userVal != "true" {
			return nil, NewInvalidParameter(
				"FifoTopic attribute must be \"true\" when topic name ends with \".fifo\"")
		}

		if _, ok := in.Attributes[snsstore.AttrContentBasedDedup]; !ok {
			in.Attributes[snsstore.AttrContentBasedDedup] = "false"
		}
	} else {
		if in.Attributes[snsstore.AttrFifoTopic] == "true" {
			return nil, NewInvalidParameter(
				"FIFO Topic names must end with \".fifo\"")
		}
	}

	// Validate all attribute values, including the DoS length cap.
	for attrName, attrValue := range in.Attributes {
		if err := validateTopicAttribute(attrName, attrValue, isFifo); err != nil {
			return nil, err
		}
	}

	// The inline DataProtectionPolicy parameter is validated with the same
	// rules as PutDataProtectionPolicy and persisted under its reserved
	// attribute key so GetDataProtectionPolicy returns it from creation.
	if in.DataProtectionPolicy != "" {
		if err := validateDataProtectionPolicy(in.DataProtectionPolicy); err != nil {
			return nil, err
		}
		in.Attributes[snsstore.AttrDataProtectionPolicy] = in.DataProtectionPolicy
	}

	// CreateTopic declares the model's TagLimitExceededException, so the
	// inline tag set carries the same constraints the tag operations
	// enforce (key/value lengths, reserved prefix, request-level count).
	if err := validateTopicTagMap(in.Tags); err != nil {
		return nil, err
	}

	topic := &snsstore.Topic{
		Name:       in.Name,
		Attributes: in.Attributes,
	}

	created, err := store.CreateTopic(topic, in.Tags)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &TopicResult{
		Arn: created.Arn,
	}, nil
}

// deleteTopicCore is the single entry point for topic deletion shared by the
// HTTP API and the admin gRPC handler.
func (s *SNSService) deleteTopicCore(store snsstore.SNSStoreInterface, topicArn string) error {
	if topicArn == "" {
		return NewInvalidParameter("TopicArn is required")
	}

	if err := store.DeleteTopic(topicArn); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// listTopicsCore is the single entry point for topic listing shared by the
// HTTP API and the admin gRPC handler.
func (s *SNSService) listTopicsCore(store snsstore.SNSStoreInterface, in ListTopicsInput) (*ListTopicsResult, error) {
	result, err := store.ListTopics(storecommon.ListOptions{Marker: in.NextToken})
	if err != nil {
		return nil, mapStoreError(err)
	}

	topics := make([]TopicSummary, 0, len(result.Items))
	for _, topic := range result.Items {
		topics = append(topics, TopicSummary{TopicArn: topic.Arn})
	}

	nextToken := ""
	if result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}

	return &ListTopicsResult{
		Topics:    topics,
		NextToken: nextToken,
	}, nil
}

// GetTopicAttributesInput carries the topic ARN whose attributes are read.
type GetTopicAttributesInput struct {
	TopicArn string
}

// SetTopicAttributesInput carries a single attribute update for a topic.
type SetTopicAttributesInput struct {
	TopicArn       string
	AttributeName  string
	AttributeValue string
}

// getTopicAttributesCore is the single validation and persistence path for
// GetTopicAttributes, including the default-policy synthesis and
// AddPermission statement injection that shape the returned Policy attribute.
//
// The emitted set follows the documented GetTopicAttributes enumeration:
// DisplayName, SignatureVersion and MaximumMessageSize appear only when
// explicitly set ("Amazon SNS returns this attribute only if you explicitly
// set it"; "If the API response does not include the SignatureVersion
// attribute, it means that the SignatureVersion for the topic has value 1"),
// EffectiveDeliveryPolicy is always synthesised from the stored
// DeliveryPolicy and the system defaults, and bookkeeping timestamps are
// not attributes at all.
func (s *SNSService) getTopicAttributesCore(store snsstore.SNSStoreInterface, in GetTopicAttributesInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	attrs := make(map[string]string)
	attrs[snsstore.AttrTopicArn] = topic.Arn
	attrs[snsstore.AttrOwner] = topic.Owner
	attrs[snsstore.AttrSubscriptionsConfirmed] = fmt.Sprintf("%d", topic.SubscriptionsConfirmed)
	attrs[snsstore.AttrSubscriptionsDeleted] = fmt.Sprintf("%d", topic.SubscriptionsDeleted)
	attrs[snsstore.AttrSubscriptionsPending] = fmt.Sprintf("%d", topic.SubscriptionsPending)

	for k, v := range topic.Attributes {
		if k == snsstore.AttrDataProtectionPolicy {
			continue
		}
		if k == snsstore.AttrPolicy && v == "" {
			continue
		}
		attrs[k] = v
	}

	if _, hasPolicy := attrs[snsstore.AttrPolicy]; !hasPolicy {
		// Default policy uses version 2012-10-17 (the current AWS standard),
		// replacing the legacy 2008-10-17 default.
		attrs[snsstore.AttrPolicy] = formatDefaultPolicy(topic.Arn, topic.Owner)
	}

	if len(topic.Permissions) > 0 {
		attrs[snsstore.AttrPolicy] = injectPermissionsIntoPolicy(attrs[snsstore.AttrPolicy], topic.Arn, topic.Permissions)
	}

	// EffectiveDeliveryPolicy: the resolved policy "taking system defaults
	// into account" — the stored topic policy with every omitted field at
	// its default, or the pure system default when none is stored. The
	// stored value passed validation at write time; a value that no longer
	// parses degrades to the system default rather than failing the read.
	effective, err := parseTopicDeliveryPolicy(topic.Attributes[snsstore.AttrDeliveryPolicy])
	if err != nil {
		effective = nil
	}
	if effective == nil {
		effective = systemDefaultDeliveryPolicy()
	}
	if rendered, renderErr := effective.effectiveTopicJSON(); renderErr == nil {
		attrs[snsstore.AttrEffectiveDelivery] = rendered
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// setTopicAttributesCore is the single validation and persistence path for
// SetTopicAttributes. The topic is read first: the FIFO-only attributes
// validate against the topic's type, and the create-only and read-only
// keys are refused on this plane.
func (s *SNSService) setTopicAttributesCore(store snsstore.SNSStoreInterface, in SetTopicAttributesInput) error {
	if in.TopicArn == "" {
		return NewInvalidParameter("TopicArn is required")
	}
	if in.AttributeName == "" {
		return NewInvalidParameter("AttributeName is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		return mapStoreError(err)
	}
	if err := rejectNonSettableTopicAttribute(in.AttributeName); err != nil {
		return err
	}
	if err := validateTopicAttribute(in.AttributeName, in.AttributeValue, topic.IsFifoTopic()); err != nil {
		return err
	}

	attrs := map[string]string{in.AttributeName: in.AttributeValue}

	if err := store.SetTopicAttributes(in.TopicArn, attrs); err != nil {
		return mapStoreError(err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Shared helpers used by HTTP handlers
// ---------------------------------------------------------------------------

// formatDefaultPolicy returns the default SNS topic policy JSON with the
// given topic ARN and owner. The policy version is 2012-10-17 (the
// current AWS standard, replacing the legacy 2008-10-17 default).
func formatDefaultPolicy(topicArn, owner string) string {
	return fmt.Sprintf(
		`{"Version":"2012-10-17","Id":"__default_policy_ID","Statement":[{"Sid":"__default_statement_ID","Effect":"Allow","Principal":{"AWS":"*"},"Action":["SNS:GetTopicAttributes","SNS:SetTopicAttributes","SNS:AddPermission","SNS:RemovePermission","SNS:DeleteTopic","SNS:Subscribe","SNS:ListSubscriptionsByTopic","SNS:Publish","SNS:Receive"],"Resource":%q,"Condition":{"StringEquals":{"AWS:SourceOwner":%q}}}]}`,
		topicArn, owner)
}

// injectPermissionsIntoPolicy merges AddPermission entries into the topic's
// resource policy JSON, returning the updated policy string.
func injectPermissionsIntoPolicy(policyJSON, topicArn string, permissions []snsstore.Permission) string {
	var policyMap struct {
		Version   string                   `json:"Version"`
		Id        string                   `json:"Id"`
		Statement []map[string]interface{} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policyJSON), &policyMap); err != nil {
		// The stored policy leaves the synthesis — AddPermission already
		// reported success, so the omission is a degradation the operator
		// must be able to see.
		logs.Warn("SNS: the topic policy does not parse into the synthesis shape; AddPermission statements are omitted from the read view",
			logs.String("topicArn", topicArn),
			logs.Err(err))
		return policyJSON
	}

	for _, perm := range permissions {
		principals := make([]string, len(perm.Principals))
		for i, p := range perm.Principals {
			principals[i] = arnutil.NewARNBuilder(p, "").IAM().Root()
		}
		actions := make([]string, len(perm.Actions))
		copy(actions, perm.Actions)
		policyMap.Statement = append(policyMap.Statement, map[string]interface{}{
			"Sid":       perm.Label,
			"Effect":    "Allow",
			"Principal": map[string]interface{}{"AWS": principals},
			"Action":    actions,
			"Resource":  topicArn,
		})
	}

	updated, err := json.Marshal(policyMap)
	if err != nil {
		logs.Warn("SNS: synthesising the topic policy failed; the stored policy is returned without the AddPermission statements",
			logs.String("topicArn", topicArn),
			logs.Err(err))
		return policyJSON
	}
	return string(updated)
}
