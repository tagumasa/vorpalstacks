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
		// The stored status is the TTL state machine's value — ENABLING or
		// DISABLING while the transition goroutine has not yet committed,
		// ENABLED or DISABLED after it has. A record without a status falls
		// back to the derived value.
		status := "DISABLED"
		if ttl.Status != "" {
			status = string(ttl.Status)
		} else if ttl.Enabled {
			status = "ENABLED"
		}
		desc := map[string]interface{}{
			"TimeToLiveStatus": status,
		}
		if ttl.AttributeName != "" {
			desc["AttributeName"] = ttl.AttributeName
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

	in := parseTimeToLiveSpec(req.Parameters["TimeToLiveSpecification"])
	in.TableName = table.Name

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ttl, err := s.updateTimeToLiveCore(ctx, store, in)
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
