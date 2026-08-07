package sns

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// CreateTopic creates a new SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_CreateTopic.html
func (s *SNSService) CreateTopic(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := CreateTopicInput{
		Name:       request.GetParamLowerFirst(req.Parameters, "Name"),
		Attributes: parseAttributes(req.Parameters),
		Tags:       tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")),
	}

	result, err := s.createTopicCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TopicArn": result.Arn,
	}, nil
}

// DeleteTopic deletes an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_DeleteTopic.html
func (s *SNSService) DeleteTopic(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if err := s.deleteTopicCore(store, topicArn); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetTopicAttributes returns the attributes of an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_GetTopicAttributes.html
func (s *SNSService) GetTopicAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	topic, err := store.GetTopic(topicArn)
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
		// L7: default policy version updated from legacy 2008-10-17 to 2012-10-17.
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

// SetTopicAttributes sets the attributes of an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_SetTopicAttributes.html
func (s *SNSService) SetTopicAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	attributeName := request.GetParamLowerFirst(req.Parameters, "AttributeName")
	attributeValue := request.GetParamLowerFirst(req.Parameters, "AttributeValue")

	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}
	if attributeName == "" {
		return nil, awserrors.NewInvalidParameterException("AttributeName is required")
	}

	if err := validateTopicAttribute(attributeName, attributeValue); err != nil {
		return nil, err
	}

	attrs := map[string]string{attributeName: attributeValue}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.SetTopicAttributes(topicArn, attrs); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
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

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	result, err := s.listTopicsCore(store, ListTopicsInput{NextToken: nextToken})
	if err != nil {
		return nil, err
	}

	topics := make([]map[string]interface{}, 0, len(result.Topics))
	for _, t := range result.Topics {
		topics = append(topics, map[string]interface{}{
			"TopicArn": t.TopicArn,
		})
	}

	return pagination.BuildListResponse("Topics", topics, result.NextToken), nil
}

func parseAttributes(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for _, key := range []string{"Attributes", "attributes"} {
		if attrs, ok := params[key].(map[string]interface{}); ok {
			for k, v := range attrs {
				if vs, ok := v.(string); ok {
					result[k] = vs
				}
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

// injectPermissionsIntoPolicy merges AddPermission entries into the topic's
// resource policy JSON, returning the updated policy string.
func injectPermissionsIntoPolicy(policyJSON, topicArn string, permissions []snsstore.Permission) string {
	var policyMap struct {
		Version   string                   `json:"Version"`
		Id        string                   `json:"Id"`
		Statement []map[string]interface{} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policyJSON), &policyMap); err != nil {
		return policyJSON
	}

	for _, perm := range permissions {
		principals := make([]string, len(perm.Principals))
		for i, p := range perm.Principals {
			principals[i] = "arn:aws:iam::" + p + ":root"
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
		return policyJSON
	}
	return string(updated)
}
