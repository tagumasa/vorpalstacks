package sns

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// parsePermissionList extracts string values from a list parameter that
// arrives in two wire formats: {"value": "..."} entries or member-named
// entries ({"AwsAccountId": "..."} / {"ActionName": "..."}).
func parsePermissionList(items []map[string]interface{}, member string) []string {
	var result []string
	for _, item := range items {
		if val, ok := item["value"].(string); ok {
			result = append(result, val)
		} else if val, ok := item[member].(string); ok {
			result = append(result, val)
		}
	}
	return result
}

// GetDataProtectionPolicy retrieves the data protection policy for the specified SNS topic.
func (s *SNSService) GetDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	topicArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if topicArn == "" {
		topicArn = request.GetParamLowerFirst(req.Parameters, "TopicArn")
	}

	return s.getDataProtectionPolicyCore(store, GetDataProtectionPolicyInput{TopicArn: topicArn})
}

// PutDataProtectionPolicy sets the data protection policy for the specified SNS topic.
func (s *SNSService) PutDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	topicArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	if topicArn == "" {
		topicArn = request.GetParamLowerFirst(req.Parameters, "TopicArn")
	}

	return s.putDataProtectionPolicyCore(store, PutDataProtectionPolicyInput{
		TopicArn: topicArn,
		Policy:   request.GetParamLowerFirst(req.Parameters, "DataProtectionPolicy"),
	})
}

// AddPermission adds permissions to the access policy of the specified SNS topic.
func (s *SNSService) AddPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Wire parsing only; the required-member, account-ID, and action-name
	// validations live in the Core.
	awsAccountIds := parsePermissionList(request.GetListParamLowerFirst(req.Parameters, "AWSAccountId"), "AwsAccountId")
	actionNames := parsePermissionList(request.GetListParamLowerFirst(req.Parameters, "ActionName"), "ActionName")

	return s.addPermissionCore(store, AddPermissionInput{
		TopicArn:      request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		Label:         request.GetParamLowerFirst(req.Parameters, "Label"),
		AWSAccountIds: awsAccountIds,
		ActionNames:   actionNames,
	})
}

// RemovePermission removes permissions from the access policy of the specified SNS topic.
func (s *SNSService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.removePermissionCore(store, RemovePermissionInput{
		TopicArn: request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		Label:    request.GetParamLowerFirst(req.Parameters, "Label"),
	})
}
