package sns

import (
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
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
// validation (including reserved-prefix rejection), FIFO consistency checks
// (M1), attribute validation (including value length caps), and persistence.
func (s *SNSService) createTopicCore(store snsstore.SNSStoreInterface, in CreateTopicInput) (*TopicResult, error) {
	if err := validateTopicName(in.Name); err != nil {
		return nil, err
	}

	isFifo := strings.HasSuffix(in.Name, ".fifo")

	if in.Attributes == nil {
		in.Attributes = make(map[string]string)
	}

	// M1: FIFO consistency — respect the user-provided FifoTopic attribute
	// rather than silently overriding it. If the name has a .fifo suffix,
	// FifoTopic must be "true" (either explicitly or by default injection).
	// If the name lacks the suffix, FifoTopic must not be "true".
	if isFifo {
		userVal, hasAttr := in.Attributes["FifoTopic"]
		if !hasAttr {
			in.Attributes["FifoTopic"] = "true"
		} else if userVal != "true" {
			return nil, awserrors.NewInvalidParameterException(
				"FifoTopic attribute must be \"true\" when topic name ends with \".fifo\"")
		}

		if _, ok := in.Attributes["ContentBasedDeduplication"]; !ok {
			in.Attributes["ContentBasedDeduplication"] = "false"
		}
	} else {
		if in.Attributes["FifoTopic"] == "true" {
			return nil, awserrors.NewInvalidParameterException(
				"FIFO Topic names must end with \".fifo\"")
		}
	}

	// Validate all attribute values (M6: includes DoS length cap).
	for attrName, attrValue := range in.Attributes {
		if err := validateTopicAttribute(attrName, attrValue); err != nil {
			return nil, err
		}
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
		return awserrors.NewInvalidParameterException("TopicArn is required")
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

// ---------------------------------------------------------------------------
// Shared helpers used by HTTP handlers
// ---------------------------------------------------------------------------

// formatDefaultPolicy returns the default SNS topic policy JSON with the
// given topic ARN and owner. The policy version is 2012-10-17 (L7: updated
// from the legacy 2008-10-17).
func formatDefaultPolicy(topicArn, owner string) string {
	return fmt.Sprintf(
		`{"Version":"2012-10-17","Id":"__default_policy_ID","Statement":[{"Sid":"__default_statement_ID","Effect":"Allow","Principal":{"AWS":"*"},"Action":["SNS:GetTopicAttributes","SNS:SetTopicAttributes","SNS:AddPermission","SNS:RemovePermission","SNS:DeleteTopic","SNS:Subscribe","SNS:ListSubscriptionsByTopic","SNS:Publish","SNS:Receive"],"Resource":%q,"Condition":{"StringEquals":{"AWS:SourceOwner":%q}}}]}`,
		topicArn, owner)
}
