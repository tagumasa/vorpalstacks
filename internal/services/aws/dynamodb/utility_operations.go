// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"vorpalstacks/internal/services/aws/common/request"
)

// DescribeEndpoints returns the service endpoints for DynamoDB.
func (s *DynamoDBService) DescribeEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"Endpoints": []map[string]interface{}{
			{
				"Address":              "dynamodb." + reqCtx.GetRegion() + ".amazonaws.com",
				"CachePeriodInMinutes": 1440,
			},
		},
	}, nil
}

// DescribeLimits returns the limits for DynamoDB operations.
func (s *DynamoDBService) DescribeLimits(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"AccountMaxReadCapacityUnits":        80000,
		"AccountMaxWriteCapacityUnits":       80000,
		"TableMaxReadCapacityUnits":          40000,
		"TableMaxWriteCapacityUnits":         40000,
		"AccountMaxReadCapacityUnitsChange":  100000,
		"AccountMaxWriteCapacityUnitsChange": 100000,
	}, nil
}
