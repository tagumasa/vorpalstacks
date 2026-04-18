package sns

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// CreateTopic creates a new SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_CreateTopic.html
func (s *SNSService) CreateTopic(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewInvalidParameterException("Topic name is required")
	}
	if len(name) > 256 {
		return nil, NewInvalidParameterException("Topic name must not exceed 256 characters")
	}
	isFifo := strings.HasSuffix(name, ".fifo")
	baseName := name
	if isFifo {
		baseName = strings.TrimSuffix(name, ".fifo")
	}
	for _, c := range baseName {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return nil, NewInvalidParameterException("Topic name can only contain alphanumeric characters, hyphens, and underscores")
		}
	}

	topic := &snsstore.Topic{
		Name: name,
	}

	topic.Tags = tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	topic.Attributes = parseAttributes(req.Parameters)

	if isFifo {
		if topic.Attributes == nil {
			topic.Attributes = make(map[string]string)
		}
		topic.FifoTopic = true
		topic.Attributes["FifoTopic"] = "true"
		if _, ok := topic.Attributes["ContentBasedDeduplication"]; !ok {
			topic.Attributes["ContentBasedDeduplication"] = "false"
		}
	}

	for k, v := range topic.Attributes {
		if k == "FifoTopic" && v == "true" {
			topic.FifoTopic = true
		}
		if k == "ContentBasedDeduplication" && v == "true" {
			topic.ContentBasedDeduplication = true
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.CreateTopic(topic)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TopicArn": created.Arn,
	}, nil
}

// DeleteTopic deletes an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_DeleteTopic.html
func (s *SNSService) DeleteTopic(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, NewInvalidParameterException("TopicArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteTopic(topicArn); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetTopicAttributes returns the attributes of an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_GetTopicAttributes.html
func (s *SNSService) GetTopicAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, NewInvalidParameterException("TopicArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	topic, err := store.GetTopic(topicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	attrs := make(map[string]string)
	attrs["TopicArn"] = topic.Arn
	attrs["DisplayName"] = topic.DisplayName
	attrs["Owner"] = topic.Owner
	attrs["SubscriptionsConfirmed"] = fmt.Sprintf("%d", topic.SubscriptionsConfirmed)
	attrs["SubscriptionsDeleted"] = fmt.Sprintf("%d", topic.SubscriptionsDeleted)
	attrs["SubscriptionsPending"] = fmt.Sprintf("%d", topic.SubscriptionsPending)

	for k, v := range topic.Attributes {
		if k == "Policy" && v == "" {
			continue
		}
		attrs[k] = v
	}

	if _, hasPolicy := attrs["Policy"]; !hasPolicy {
		defaultPolicy := fmt.Sprintf(`{"Version":"2008-10-17","Id":"__default_policy_ID","Statement":[{"Sid":"__default_statement_ID","Effect":"Allow","Principal":{"AWS":"*"},"Action":["SNS:GetTopicAttributes","SNS:SetTopicAttributes","SNS:AddPermission","SNS:RemovePermission","SNS:DeleteTopic","SNS:Subscribe","SNS:ListSubscriptionsByTopic","SNS:Publish","SNS:Receive"],"Resource":"%s","Condition":{"StringEquals":{"AWS:SourceOwner":"%s"}}}]}`, topic.Arn, topic.Owner)
		attrs["Policy"] = defaultPolicy
	}

	if topic.FifoTopic {
		attrs["FifoTopic"] = "true"
	}
	if topic.ContentBasedDeduplication {
		attrs["ContentBasedDeduplication"] = "true"
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// SetTopicAttributes sets the attributes of an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_SetTopicAttributes.html
func (s *SNSService) SetTopicAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	attributeName := request.GetParamLowerFirst(req.Parameters, "AttributeName")
	attributeValue := request.GetParamLowerFirst(req.Parameters, "AttributeValue")

	if topicArn == "" {
		return nil, NewInvalidParameterException("TopicArn is required")
	}
	if attributeName == "" {
		return nil, NewInvalidParameterException("AttributeName is required")
	}

	attrs := map[string]string{attributeName: attributeValue}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.SetTopicAttributes(topicArn, attrs); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTopics lists the SNS topics.
// https://docs.aws.amazon.com/sns/latest/api/API_ListTopics.html
func (s *SNSService) ListTopics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	nextToken := request.GetStringParam(req.Parameters, "NextToken")
	result, err := store.ListTopics(common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, err
	}

	topics := make([]map[string]interface{}, 0, len(result.Items))
	for _, topic := range result.Items {
		topics = append(topics, map[string]interface{}{
			"TopicArn": topic.Arn,
		})
	}

	nextToken = ""
	if result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}
	return pagination.BuildListResponse("Topics", topics, nextToken), nil
}

func parseAttributes(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	if attrs, ok := params["Attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}
	if attrs, ok := params["attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}

	for i := 1; ; i++ {
		keyKey := "Attributes.entry." + strconv.Itoa(i) + ".key"
		valueKey := "Attributes.entry." + strconv.Itoa(i) + ".value"

		key := request.GetStringParam(params, keyKey)
		if key == "" {
			break
		}

		value := request.GetStringParam(params, valueKey)
		result[key] = value
	}

	return result
}
