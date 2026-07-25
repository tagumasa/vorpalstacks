package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// PutDataProtectionPolicy creates or updates a data protection policy for a log group.
func (s *LogsService) PutDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")

	if logGroupIdentifier == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetLogGroup(logGroupIdentifier)
	if err != nil {
		return nil, mapStoreError(err)
	}

	dpp := &logsstore.DataProtectionPolicy{
		LogGroupIdentifier: logGroupIdentifier,
		PolicyDocument:     policyDocument,
	}

	if err := store.PutDataProtectionPolicy(dpp); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"logGroupIdentifier": dpp.LogGroupIdentifier,
		"policyDocument":     dpp.PolicyDocument,
		"lastUpdatedTime":    dpp.LastUpdatedTime,
	}, nil
}

// GetDataProtectionPolicy retrieves the data protection policy for a log group.
func (s *LogsService) GetDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	if logGroupIdentifier == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	dpp, err := store.GetDataProtectionPolicy(logGroupIdentifier)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"logGroupIdentifier": dpp.LogGroupIdentifier,
		"policyDocument":     dpp.PolicyDocument,
		"lastUpdatedTime":    dpp.LastUpdatedTime,
	}, nil
}

// DeleteDataProtectionPolicy deletes the data protection policy for a log group.
func (s *LogsService) DeleteDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	if logGroupIdentifier == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteDataProtectionPolicy(logGroupIdentifier); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}
