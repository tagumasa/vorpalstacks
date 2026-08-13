package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putDataProtectionPolicyCore(logGroupIdentifier, policyDocument, region string) (*logsstore.DataProtectionPolicy, error) {
	if logGroupIdentifier == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}

	if _, err = store.GetLogGroup(logGroupIdentifier); err != nil {
		return nil, mapStoreError(err)
	}

	dpp := &logsstore.DataProtectionPolicy{
		LogGroupIdentifier: logGroupIdentifier,
		PolicyDocument:     policyDocument,
	}

	if err := store.PutDataProtectionPolicy(dpp); err != nil {
		return nil, mapStoreError(err)
	}
	return dpp, nil
}

func (s *LogsService) getDataProtectionPolicyCore(logGroupIdentifier, region string) (*logsstore.DataProtectionPolicy, error) {
	if logGroupIdentifier == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}

	return store.GetDataProtectionPolicy(logGroupIdentifier)
}

func (s *LogsService) deleteDataProtectionPolicyCore(logGroupIdentifier, region string) error {
	if logGroupIdentifier == "" {
		return ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteDataProtectionPolicy(logGroupIdentifier); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- HTTP handlers ---

func (s *LogsService) PutDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")

	dpp, err := s.putDataProtectionPolicyCore(logGroupIdentifier, policyDocument, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"logGroupIdentifier": dpp.LogGroupIdentifier,
		"policyDocument":     dpp.PolicyDocument,
		"lastUpdatedTime":    dpp.LastUpdatedTime,
	}, nil
}

func (s *LogsService) GetDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")

	dpp, err := s.getDataProtectionPolicyCore(logGroupIdentifier, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"logGroupIdentifier": dpp.LogGroupIdentifier,
		"policyDocument":     dpp.PolicyDocument,
		"lastUpdatedTime":    dpp.LastUpdatedTime,
	}, nil
}

func (s *LogsService) DeleteDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")

	if err := s.deleteDataProtectionPolicyCore(logGroupIdentifier, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
