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
		"MaxQueryTCU":             settings.MaxQueryTCU,
		"QueryPricingMode":        settings.QueryPricingMode,
		"QueryComputeType":        settings.QueryComputeType,
		"EncryptionConfiguration": settings.EncryptionConfiguration,
	}, nil
}

// UpdateAccountSettings updates the account settings for Timestream Query.
func (s *Service) UpdateAccountSettings(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var maxQueryTCU *int64
	if val := request.GetIntParam(req.Parameters, "MaxQueryTCU"); val > 0 {
		tcu := int64(val)
		maxQueryTCU = &tcu
	}

	queryPricingMode := request.GetParamCaseInsensitive(req.Parameters, "QueryPricingMode")
	queryComputeType := request.GetParamCaseInsensitive(req.Parameters, "QueryComputeType")
	encryptionConfiguration := ""

	if encMap := request.GetMapParamCaseInsensitive(req.Parameters, "EncryptionConfiguration"); encMap != nil {
		if opt, ok := encMap["KmsKeyId"].(string); ok {
			encryptionConfiguration = opt
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	settings, err := store.accountSettingsStore.UpdateAccountSettings(maxQueryTCU, queryPricingMode, queryComputeType, encryptionConfiguration)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"MaxQueryTCU":             settings.MaxQueryTCU,
		"QueryPricingMode":        settings.QueryPricingMode,
		"QueryComputeType":        settings.QueryComputeType,
		"EncryptionConfiguration": settings.EncryptionConfiguration,
	}, nil
}
