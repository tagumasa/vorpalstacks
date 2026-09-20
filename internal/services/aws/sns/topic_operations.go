package sns

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
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

		DataProtectionPolicy: request.GetParamLowerFirst(req.Parameters, "DataProtectionPolicy"),
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
//
// Attribute assembly, default-policy synthesis, and AddPermission statement
// injection live inside the Core.
func (s *SNSService) GetTopicAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getTopicAttributesCore(store, GetTopicAttributesInput{
		TopicArn: request.GetParamLowerFirst(req.Parameters, "TopicArn"),
	})
}

// SetTopicAttributes sets the attributes of an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_SetTopicAttributes.html
func (s *SNSService) SetTopicAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setTopicAttributesCore(store, SetTopicAttributesInput{
		TopicArn:       request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		AttributeName:  request.GetParamLowerFirst(req.Parameters, "AttributeName"),
		AttributeValue: request.GetParamLowerFirst(req.Parameters, "AttributeValue"),
	}); err != nil {
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
