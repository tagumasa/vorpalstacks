package sqs

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// AddPermission adds a permission to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_AddPermission.html
func (s *SQSService) AddPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	label := request.GetParamCaseInsensitive(req.Parameters, "Label")

	var awsAccountIDs []string
	var actions []string

	// Parse AWS account IDs: support JSON array, flattened query, and member.N formats.
	if ids, ok := req.Parameters["AWSAccountIds"].([]interface{}); ok && len(ids) > 0 {
		for _, id := range ids {
			if s, ok := id.(string); ok && s != "" {
				awsAccountIDs = append(awsAccountIDs, s)
			}
		}
	} else {
		for i := 1; ; i++ {
			accountID := request.GetParamCaseInsensitive(req.Parameters, "AWSAccountId."+strconv.Itoa(i))
			if accountID == "" {
				aidKey := "AWSAccountId." + strconv.Itoa(i)
				if val, ok := req.Parameters[aidKey].(string); ok {
					accountID = val
				}
			}
			if accountID == "" {
				break
			}
			awsAccountIDs = append(awsAccountIDs, accountID)
		}
	}

	// Parse actions: support JSON array, flattened query (ActionName.N), and legacy (Action.N) formats.
	if acts, ok := req.Parameters["Actions"].([]interface{}); ok && len(acts) > 0 {
		for _, a := range acts {
			if s, ok := a.(string); ok && s != "" {
				actions = append(actions, s)
			}
		}
	} else {
		for i := 1; ; i++ {
			action := request.GetParamCaseInsensitive(req.Parameters, "ActionName."+strconv.Itoa(i))
			if action == "" {
				actKey := "ActionName." + strconv.Itoa(i)
				if val, ok := req.Parameters[actKey].(string); ok {
					action = val
				}
			}
			if action == "" {
				action = request.GetParamCaseInsensitive(req.Parameters, "Action."+strconv.Itoa(i))
			}
			if action == "" {
				break
			}
			actions = append(actions, action)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.addPermissionCore(store, AddPermissionInput{
		QueueURL:      queueURL,
		Label:         label,
		AWSAccountIDs: awsAccountIDs,
		Actions:       actions,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// RemovePermission removes a permission from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_RemovePermission.html
func (s *SQSService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	label := request.GetParamCaseInsensitive(req.Parameters, "Label")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.removePermissionCore(store, RemovePermissionInput{
		QueueURL: queueURL,
		Label:    label,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
