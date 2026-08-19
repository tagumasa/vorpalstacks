package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeTimeToLive returns the Time to Live settings for a table.
func (s *DynamoDBService) DescribeTimeToLive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ttl, err := s.describeTimeToLiveCore(store, table.Name)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"TimeToLiveDescription": map[string]interface{}{
			"TimeToLiveStatus": "DISABLED",
		},
	}

	if ttl != nil {
		desc := map[string]interface{}{
			"TimeToLiveStatus": "DISABLED",
		}
		if ttl.AttributeName != "" {
			desc["AttributeName"] = ttl.AttributeName
		}
		if ttl.Enabled {
			desc["TimeToLiveStatus"] = "ENABLED"
		}
		resp["TimeToLiveDescription"] = desc
	}

	return resp, nil
}

// UpdateTimeToLive enables or disables Time to Live for a table.
func (s *DynamoDBService) UpdateTimeToLive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	ttlSpec, ok := req.Parameters["TimeToLiveSpecification"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	enabled, err := validateBoolParam(ttlSpec, "Enabled", false)
	if err != nil {
		return nil, err
	}
	attrName, ok := ttlSpec["AttributeName"].(string)
	if !ok {
		return nil, ErrInvalidParameter
	}

	if attrName == "" {
		return nil, ErrInvalidParameter
	}
	if !validateTimeToLiveAttributeName(attrName) {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Enabling TTL on a table that already has TTL enabled is rejected;
	// renaming the TTL attribute requires disabling TTL first. Disabling
	// is always allowed and is the documented path for changing the
	// attribute.
	existingTTL, _ := s.describeTimeToLiveCore(store, table.Name)
	if enabled && existingTTL != nil && existingTTL.Enabled {
		return nil, ErrInvalidParameter
	}

	ttl, err := s.updateTimeToLiveCore(ctx, store, UpdateTimeToLiveInput{
		TableName:     table.Name,
		Enabled:       enabled,
		AttributeName: attrName,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TimeToLiveSpecification": map[string]interface{}{
			"Enabled":       ttl.Enabled,
			"AttributeName": ttl.AttributeName,
		},
	}, nil
}
