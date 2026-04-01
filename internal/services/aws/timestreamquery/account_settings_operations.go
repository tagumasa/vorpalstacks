package timestreamquery

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
)

// DescribeAccountSettings returns the account settings for Timestream Query.
func (s *Service) DescribeAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	settings, err := store.accountSettingsStore.GetAccountSettings()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"MaxQueryTCU":       settings.MaxQueryTCU,
		"QueryPricingModel": settings.QueryPricingMode,
		"QueryCompute": map[string]interface{}{
			"ComputeMode": settings.QueryComputeType,
		},
	}, nil
}

// UpdateAccountSettings updates the account settings for Timestream Query.
func (s *Service) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var maxQueryTCU *int64
	if val := request.GetIntParam(req.Parameters, "MaxQueryTCU"); val > 0 {
		tcu := int64(val)
		maxQueryTCU = &tcu
	}

	queryPricingModel := request.GetParamCaseInsensitive(req.Parameters, "QueryPricingModel")
	queryComputeType := ""

	if qcMap := request.GetMapParamCaseInsensitive(req.Parameters, "QueryCompute"); qcMap != nil {
		if mode, ok := qcMap["ComputeMode"].(string); ok {
			queryComputeType = mode
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	settings, err := store.accountSettingsStore.UpdateAccountSettings(maxQueryTCU, queryPricingModel, queryComputeType, "")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"MaxQueryTCU":       settings.MaxQueryTCU,
		"QueryPricingModel": settings.QueryPricingMode,
		"QueryCompute": map[string]interface{}{
			"ComputeMode": settings.QueryComputeType,
		},
	}, nil
}
