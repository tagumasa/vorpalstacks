// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"net/url"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/config"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DescribeEndpoints returns the service endpoints for DynamoDB.
func (s *DynamoDBService) DescribeEndpoints(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Report the endpoint the caller actually reached. The AWS hostname
	// format is unreachable on an edge or on-premises deployment, so the
	// client-facing host wins: the request Host header first, then the
	// configured endpoints.base_url.
	address := ""
	if req != nil && req.Headers != nil {
		address = req.Headers.Get("Host")
	}
	if address == "" {
		if parsed, err := url.Parse(config.BaseURL()); err == nil && parsed.Host != "" {
			address = parsed.Host
		}
	}
	if address == "" {
		address = "dynamodb." + reqCtx.GetRegion() + ".amazonaws.com"
	}
	return map[string]interface{}{
		"Endpoints": []map[string]interface{}{
			{
				"Address":              address,
				"CachePeriodInMinutes": 1440,
			},
		},
	}, nil
}

// DescribeLimits returns the limits for DynamoDB operations.
func (s *DynamoDBService) DescribeLimits(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"AccountMaxReadCapacityUnits":  dbstore.AccountMaxReadCapacityUnits,
		"AccountMaxWriteCapacityUnits": dbstore.AccountMaxWriteCapacityUnits,
		"TableMaxReadCapacityUnits":    dbstore.TableMaxReadCapacityUnits,
		"TableMaxWriteCapacityUnits":   dbstore.TableMaxWriteCapacityUnits,
	}, nil
}
