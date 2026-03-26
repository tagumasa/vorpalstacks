// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
)

// DescribeTableReplicaAutoScaling returns the auto scaling settings for table replicas.
func (s *DynamoDBService) DescribeTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    []map[string]interface{}{},
		},
	}, nil
}

// UpdateTableReplicaAutoScaling updates the auto scaling settings for table replicas.
func (s *DynamoDBService) UpdateTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    []map[string]interface{}{},
		},
	}, nil
}
