package sns

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// GetDataProtectionPolicy retrieves the data protection policy for the specified SNS topic.
func (s *SNSService) GetDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if topicArn == "" {
		topicArn = request.GetParamLowerFirst(req.Parameters, "TopicArn")
	}
	if topicArn == "" {
		return nil, NewInvalidParameterException("ResourceArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.GetDataProtectionPolicy(topicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	return map[string]interface{}{
		"DataProtectionPolicy": policy,
	}, nil
}

// PutDataProtectionPolicy sets the data protection policy for the specified SNS topic.
func (s *SNSService) PutDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if topicArn == "" {
		topicArn = request.GetParamLowerFirst(req.Parameters, "TopicArn")
	}
	if topicArn == "" {
		return nil, NewInvalidParameterException("ResourceArn is required")
	}

	policy := request.GetParamLowerFirst(req.Parameters, "DataProtectionPolicy")
	if policy == "" {
		return nil, NewInvalidParameterException("DataProtectionPolicy is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.PutDataProtectionPolicy(topicArn, policy); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AddPermission adds permissions to the access policy of the specified SNS topic.
func (s *SNSService) AddPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, NewInvalidParameterException("TopicArn is required")
	}

	label := request.GetParamLowerFirst(req.Parameters, "Label")
	if label == "" {
		return nil, NewInvalidParameterException("Label is required")
	}

	awsAccountIdsRaw := request.GetListParamLowerFirst(req.Parameters, "AWSAccountId")
	actionNamesRaw := request.GetListParamLowerFirst(req.Parameters, "ActionName")

	var awsAccountIds []string
	for _, item := range awsAccountIdsRaw {
		if s, ok := item["value"].(string); ok {
			awsAccountIds = append(awsAccountIds, s)
		} else if s, ok := item["AwsAccountId"].(string); ok {
			awsAccountIds = append(awsAccountIds, s)
		}
	}

	var actionNames []string
	for _, item := range actionNamesRaw {
		if s, ok := item["value"].(string); ok {
			actionNames = append(actionNames, s)
		} else if s, ok := item["ActionName"].(string); ok {
			actionNames = append(actionNames, s)
		}
	}

	if len(awsAccountIds) == 0 {
		return nil, NewInvalidParameterException("AwsAccountId is required")
	}
	if len(actionNames) == 0 {
		return nil, NewInvalidParameterException("ActionName is required")
	}

	permission := &snsstore.Permission{
		Label:      label,
		Principals: awsAccountIds,
		Actions:    actionNames,
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.AddPermission(topicArn, permission); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemovePermission removes permissions from the access policy of the specified SNS topic.
func (s *SNSService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, NewInvalidParameterException("TopicArn is required")
	}

	label := request.GetParamLowerFirst(req.Parameters, "Label")
	if label == "" {
		return nil, NewInvalidParameterException("Label is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.RemovePermission(topicArn, label); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
