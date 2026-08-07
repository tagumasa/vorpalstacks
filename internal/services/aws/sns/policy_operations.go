package sns

import (
	"context"
	"encoding/json"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// validActionNames lists the SNS action names accepted by AddPermission.
var validActionNames = map[string]bool{
	"Publish":                  true,
	"GetTopicAttributes":       true,
	"SetTopicAttributes":       true,
	"Subscribe":                true,
	"ListSubscriptionsByTopic": true,
	"DeleteTopic":              true,
	"Receive":                  true,
	"AddPermission":            true,
	"RemovePermission":         true,
}

// GetDataProtectionPolicy retrieves the data protection policy for the specified SNS topic.
func (s *SNSService) GetDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if topicArn == "" {
		topicArn = request.GetParamLowerFirst(req.Parameters, "TopicArn")
	}
	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("ResourceArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.GetDataProtectionPolicy(topicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
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
		return nil, awserrors.NewInvalidParameterException("ResourceArn is required")
	}

	policy := request.GetParamLowerFirst(req.Parameters, "DataProtectionPolicy")
	if policy == "" {
		return nil, awserrors.NewInvalidParameterException("DataProtectionPolicy is required")
	}

	var policyCheck interface{}
	if err := json.Unmarshal([]byte(policy), &policyCheck); err != nil {
		return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid DataProtectionPolicy: not valid JSON: %s", err.Error()))
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.PutDataProtectionPolicy(topicArn, policy); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// AddPermission adds permissions to the access policy of the specified SNS topic.
func (s *SNSService) AddPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}

	label := request.GetParamLowerFirst(req.Parameters, "Label")
	if label == "" {
		return nil, awserrors.NewInvalidParameterException("Label is required")
	}

	awsAccountIdsRaw := request.GetListParamLowerFirst(req.Parameters, "AWSAccountId")
	actionNamesRaw := request.GetListParamLowerFirst(req.Parameters, "ActionName")

	var awsAccountIds []string
	for _, item := range awsAccountIdsRaw {
		if val, ok := item["value"].(string); ok {
			awsAccountIds = append(awsAccountIds, val)
		} else if val, ok := item["AwsAccountId"].(string); ok {
			awsAccountIds = append(awsAccountIds, val)
		}
	}

	var actionNames []string
	for _, item := range actionNamesRaw {
		if val, ok := item["value"].(string); ok {
			actionNames = append(actionNames, val)
		} else if val, ok := item["ActionName"].(string); ok {
			actionNames = append(actionNames, val)
		}
	}

	if len(awsAccountIds) == 0 {
		return nil, awserrors.NewInvalidParameterException("AwsAccountId is required")
	}
	if len(actionNames) == 0 {
		return nil, awserrors.NewInvalidParameterException("ActionName is required")
	}

	for _, id := range awsAccountIds {
		if len(id) != 12 {
			return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid AWS account ID %q: must be 12 digits", id))
		}
		for _, c := range id {
			if c < '0' || c > '9' {
				return nil, awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid AWS account ID %q: must be numeric", id))
			}
		}
	}

	for _, action := range actionNames {
		if !validActionNames[action] {
			return nil, awserrors.NewInvalidParameterException(fmt.Sprintf(
				"Invalid action name %q. Valid values: GetTopicAttributes, SetTopicAttributes, AddPermission, RemovePermission, DeleteTopic, Subscribe, ListSubscriptionsByTopic, Publish, Receive",
				action))
		}
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
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemovePermission removes permissions from the access policy of the specified SNS topic.
func (s *SNSService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}

	label := request.GetParamLowerFirst(req.Parameters, "Label")
	if label == "" {
		return nil, awserrors.NewInvalidParameterException("Label is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.RemovePermission(topicArn, label); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
