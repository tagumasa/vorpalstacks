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

	if attrName != "" {
		if !validateTimeToLiveAttributeName(attrName) {
			return nil, ErrInvalidParameter
		}
	}

	if enabled && attrName == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Read the existing TTL state to determine the correct response status.
	existingTTL, _ := s.describeTimeToLiveCore(store, table.Name)

	ttl, err := s.updateTimeToLiveCore(ctx, store, UpdateTimeToLiveInput{
		TableName:     table.Name,
		Enabled:       enabled,
		AttributeName: attrName,
	})
	if err != nil {
		return nil, err
	}

	// If the requested state matches the existing state, the table is
	// already in the desired state — return ENABLED/DISABLED instead of
	// the transition state ENABLING/DISABLING.
	status := "ENABLING"
	if !ttl.Enabled {
		status = "DISABLING"
	}
	if existingTTL != nil && existingTTL.Enabled == ttl.Enabled {
		if ttl.Enabled {
			status = "ENABLED"
		} else {
			status = "DISABLED"
		}
	}

	return map[string]interface{}{
		"TimeToLiveSpecification": map[string]interface{}{
			"Enabled":          ttl.Enabled,
			"AttributeName":    ttl.AttributeName,
			"TimeToLiveStatus": status,
		},
	}, nil
}
