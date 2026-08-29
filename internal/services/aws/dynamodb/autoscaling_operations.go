// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeTableReplicaAutoScaling returns the auto scaling settings for table replicas.
func (s *DynamoDBService) DescribeTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeTableReplicaAutoScalingCore(ctx, reqCtx, describeTableReplicaAutoScalingInput{
		Parameters: req.Parameters,
	})
}

// UpdateTableReplicaAutoScaling updates the auto scaling settings for table replicas.
func (s *DynamoDBService) UpdateTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateTableReplicaAutoScalingCore(ctx, reqCtx, updateTableReplicaAutoScalingInput{
		Parameters: req.Parameters,
	})
}
