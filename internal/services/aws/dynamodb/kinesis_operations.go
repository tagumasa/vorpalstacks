// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeKinesisStreamingDestination returns the Kinesis streaming destination for a table.
func (s *DynamoDBService) DescribeKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.describeKinesisStreamingDestinationCore(ctx, reqCtx, describeKinesisStreamingDestinationInput{
		Parameters: req.Parameters,
	})
}

// EnableKinesisStreamingDestination enables Kinesis streaming for a DynamoDB table.
func (s *DynamoDBService) EnableKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.enableKinesisStreamingDestinationCore(ctx, reqCtx, enableKinesisStreamingDestinationInput{
		Parameters: req.Parameters,
	})
}

// DisableKinesisStreamingDestination disables Kinesis streaming for a DynamoDB table.
func (s *DynamoDBService) DisableKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.disableKinesisStreamingDestinationCore(ctx, reqCtx, disableKinesisStreamingDestinationInput{
		Parameters: req.Parameters,
	})
}

// UpdateKinesisStreamingDestination updates the Kinesis streaming destination for a DynamoDB table.
func (s *DynamoDBService) UpdateKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.updateKinesisStreamingDestinationCore(ctx, reqCtx, updateKinesisStreamingDestinationInput{
		Parameters: req.Parameters,
	})
}
