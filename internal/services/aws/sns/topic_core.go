package sns

import (
	"fmt"
	"strings"
	"time"

	storecommon "vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
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
// topic.
type TopicResult struct {
	Arn  string
	Name string
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
// validation (including reserved-prefix rejection), FIFO consistency checks,
// attribute validation (including value length caps), and persistence.
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
		userVal, hasAttr := in.Attributes["FifoTopic"]
		if !hasAttr {
			in.Attributes["FifoTopic"] = "true"
		} else if userVal != "true" {
			return nil, NewInvalidParameter(
				"FifoTopic attribute must be \"true\" when topic name ends with \".fifo\"")
		}

		if _, ok := in.Attributes["ContentBasedDeduplication"]; !ok {
			in.Attributes["ContentBasedDeduplication"] = "false"
		}
	} else {
		if in.Attributes["FifoTopic"] == "true" {
			return nil, NewInvalidParameter(
				"FIFO Topic names must end with \".fifo\"")
		}
	}

	// Validate all attribute values, including the DoS length cap.
	for attrName, attrValue := range in.Attributes {
		if err := validateTopicAttribute(attrName, attrValue); err != nil {
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
		if in.Attributes == nil {
			in.Attributes = make(map[string]string)
		}
		in.Attributes["DataProtectionPolicy"] = in.DataProtectionPolicy
	}

	topic := &snsstore.Topic{
		Name:       in.Name,
		Attributes: in.Attributes,
		Tags:       in.Tags,
	}

	created, err := store.CreateTopic(topic)
	if err != nil {
		return nil, err
	}

	return &TopicResult{
		Arn:  created.Arn,
		Name: created.Name,
	}, nil
}

// deleteTopicCore is the single entry point for topic deletion shared by the
// HTTP API and the admin gRPC handler.
func (s *SNSService) deleteTopicCore(store snsstore.SNSStoreInterface, topicArn string) error {
	if topicArn == "" {
		return NewInvalidParameter("TopicArn is required")
	}

	if err := store.DeleteTopic(topicArn); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return ErrTopicNotFound
		}
		return err
	}
	return nil
}

// listTopicsCore is the single entry point for topic listing shared by the
// HTTP API and the admin gRPC handler.
func (s *SNSService) listTopicsCore(store snsstore.SNSStoreInterface, in ListTopicsInput) (*ListTopicsResult, error) {
	result, err := store.ListTopics(storecommon.ListOptions{Marker: in.NextToken})
	if err != nil {
		return nil, err
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
func (s *SNSService) getTopicAttributesCore(store snsstore.SNSStoreInterface, in GetTopicAttributesInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	attrs := make(map[string]string)
	attrs["TopicArn"] = topic.Arn
	attrs["DisplayName"] = topic.GetDisplayName()
	attrs["Owner"] = topic.Owner
	attrs["SubscriptionsConfirmed"] = fmt.Sprintf("%d", topic.SubscriptionsConfirmed)
	attrs["SubscriptionsDeleted"] = fmt.Sprintf("%d", topic.SubscriptionsDeleted)
	attrs["SubscriptionsPending"] = fmt.Sprintf("%d", topic.SubscriptionsPending)

	for k, v := range topic.Attributes {
		if k == "DataProtectionPolicy" {
			continue
		}
		if k == "Policy" && v == "" {
			continue
		}
		attrs[k] = v
	}

	if _, hasPolicy := attrs["Policy"]; !hasPolicy {
		// Default policy uses version 2012-10-17 (the current AWS standard),
		// replacing the legacy 2008-10-17 default.
		attrs["Policy"] = formatDefaultPolicy(topic.Arn, topic.Owner)
	}

	if len(topic.Permissions) > 0 {
		attrs["Policy"] = injectPermissionsIntoPolicy(attrs["Policy"], topic.Arn, topic.Permissions)
	}

	if !topic.CreatedDate.IsZero() {
		attrs["CreatedDate"] = topic.CreatedDate.UTC().Format(time.RFC3339)
	}
	if !topic.LastModifiedTime.IsZero() {
		attrs["LastModifiedTime"] = topic.LastModifiedTime.UTC().Format(time.RFC3339)
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// setTopicAttributesCore is the single validation and persistence path for
// SetTopicAttributes.
func (s *SNSService) setTopicAttributesCore(store snsstore.SNSStoreInterface, in SetTopicAttributesInput) error {
	if in.TopicArn == "" {
		return NewInvalidParameter("TopicArn is required")
	}
	if in.AttributeName == "" {
		return NewInvalidParameter("AttributeName is required")
	}

	if err := validateTopicAttribute(in.AttributeName, in.AttributeValue); err != nil {
		return err
	}

	attrs := map[string]string{in.AttributeName: in.AttributeValue}

	if err := store.SetTopicAttributes(in.TopicArn, attrs); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return ErrTopicNotFound
		}
		return err
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
