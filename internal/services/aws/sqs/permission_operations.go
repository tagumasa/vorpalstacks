package sqs

import (
	"context"
	"strconv"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
)

// AddPermission adds permission to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_AddPermission.html
func (s *SQSService) AddPermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	label := request.GetParamCaseInsensitive(req.Parameters, "Label")
	if label == "" {
		return nil, ErrMissingParameter
	}

	var awsAccountIDs []string
	var actions []string

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

	for i := 1; ; i++ {
		action := request.GetParamCaseInsensitive(req.Parameters, "Action."+strconv.Itoa(i))
		if action == "" {
			actKey := "Action." + strconv.Itoa(i)
			if val, ok := req.Parameters[actKey].(string); ok {
				action = val
			}
		}
		if action == "" {
			break
		}
		actions = append(actions, action)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.AddPermission(queueURL, label, awsAccountIDs, actions); err != nil {
		return nil, convertStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// RemovePermission removes permission from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_RemovePermission.html
func (s *SQSService) RemovePermission(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	label := request.GetParamCaseInsensitive(req.Parameters, "Label")
	if label == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.RemovePermission(queueURL, label); err != nil {
		return nil, convertStoreError(err)
	}

	return response.EmptyResponse(), nil
}
