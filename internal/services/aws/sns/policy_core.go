package sns

import (
	"fmt"

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

// GetDataProtectionPolicyInput carries the resolved topic ARN for reading the
// data protection policy.
type GetDataProtectionPolicyInput struct {
	TopicArn string
}

// PutDataProtectionPolicyInput carries the topic ARN and policy document for
// storing the data protection policy.
type PutDataProtectionPolicyInput struct {
	TopicArn string
	Policy   string
}

// AddPermissionInput carries the parsed permission statement for AddPermission.
type AddPermissionInput struct {
	TopicArn      string
	Label         string
	AWSAccountIds []string
	ActionNames   []string
}

// RemovePermissionInput carries the topic ARN and statement label for
// RemovePermission.
type RemovePermissionInput struct {
	TopicArn string
	Label    string
}

// getDataProtectionPolicyCore is the single validation and persistence path
// for GetDataProtectionPolicy.
func (s *SNSService) getDataProtectionPolicyCore(store snsstore.SNSStoreInterface, in GetDataProtectionPolicyInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("ResourceArn is required")
	}

	policy, err := store.GetDataProtectionPolicy(in.TopicArn)
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

// putDataProtectionPolicyCore is the single validation and persistence path
// for PutDataProtectionPolicy.
func (s *SNSService) putDataProtectionPolicyCore(store snsstore.SNSStoreInterface, in PutDataProtectionPolicyInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("ResourceArn is required")
	}

	if in.Policy == "" {
		return nil, NewInvalidParameter("DataProtectionPolicy is required")
	}

	if err := validateDataProtectionPolicy(in.Policy); err != nil {
		return nil, err
	}

	if err := store.PutDataProtectionPolicy(in.TopicArn, in.Policy); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// addPermissionCore is the single validation and persistence path for
// AddPermission.
func (s *SNSService) addPermissionCore(store snsstore.SNSStoreInterface, in AddPermissionInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}

	if in.Label == "" {
		return nil, NewInvalidParameter("Label is required")
	}

	if len(in.AWSAccountIds) == 0 {
		return nil, NewInvalidParameter("AwsAccountId is required")
	}
	if len(in.ActionNames) == 0 {
		return nil, NewInvalidParameter("ActionName is required")
	}

	for _, id := range in.AWSAccountIds {
		if len(id) != 12 {
			return nil, NewInvalidParameter(fmt.Sprintf("Invalid AWS account ID %q: must be 12 digits", id))
		}
		for _, c := range id {
			if c < '0' || c > '9' {
				return nil, NewInvalidParameter(fmt.Sprintf("Invalid AWS account ID %q: must be numeric", id))
			}
		}
	}

	for _, action := range in.ActionNames {
		if !validActionNames[action] {
			return nil, NewInvalidParameter(fmt.Sprintf(
				"Invalid action name %q. Valid values: GetTopicAttributes, SetTopicAttributes, AddPermission, RemovePermission, DeleteTopic, Subscribe, ListSubscriptionsByTopic, Publish, Receive",
				action))
		}
	}

	permission := &snsstore.Permission{
		Label:      in.Label,
		Principals: in.AWSAccountIds,
		Actions:    in.ActionNames,
	}

	if err := store.AddPermission(in.TopicArn, permission); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// removePermissionCore is the single validation and persistence path for
// RemovePermission.
func (s *SNSService) removePermissionCore(store snsstore.SNSStoreInterface, in RemovePermissionInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}

	if in.Label == "" {
		return nil, NewInvalidParameter("Label is required")
	}

	if err := store.RemovePermission(in.TopicArn, in.Label); err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}
