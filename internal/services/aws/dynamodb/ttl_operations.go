// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DescribeTimeToLive returns the Time to Live settings for a table.
func (s *DynamoDBService) DescribeTimeToLive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	ttl, err := store.Tables().GetTimeToLive(tableName)
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
	tableName := table.Name

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
		if err := validateTimeToLiveAttributeName(attrName); err != nil {
			return nil, err
		}
	}

	if enabled && attrName == "" {
		return nil, ErrInvalidParameter
	}

	ttl := &dbstore.TimeToLiveSpecification{
		Enabled:       enabled,
		AttributeName: attrName,
	}
	if enabled {
		ttl.Status = dbstore.TTLStatusEnabling
	} else {
		ttl.Status = dbstore.TTLStatusDisabling
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().SetTimeToLive(tableName, ttl); err != nil {
		return nil, err
	}

	// Background transition to final state.
	go func() {
		time.Sleep(1 * time.Second)
		if enabled {
			ttl.Status = dbstore.TTLStatusEnabled
		} else {
			ttl.Status = dbstore.TTLStatusDisabled
		}
		if err := store.Tables().SetTimeToLive(tableName, ttl); err != nil {
			logs.Error("Failed to transition TTL to final state",
				logs.Err(err),
				logs.String("tableName", tableName),
			)
		}
	}()

	// Response shows the transition state.
	status := "ENABLING"
	if !enabled {
		status = "DISABLING"
	}

	return map[string]interface{}{
		"TimeToLiveSpecification": map[string]interface{}{
			"Enabled":          enabled,
			"AttributeName":    attrName,
			"TimeToLiveStatus": status,
		},
	}, nil
}
