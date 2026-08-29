// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeKinesisStreamingDestination returns the Kinesis streaming destination for a table.
func (s *DynamoDBService) DescribeKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	var destinations []map[string]interface{}
	for _, d := range table.KinesisDataStreamDestinations {
		dest := map[string]interface{}{
			"StreamArn":                    d.StreamArn,
			"DestinationStatus":            d.DestinationStatus,
			"DestinationStatusDescription": d.DestinationStatusDescription,
		}
		if d.ApproximateCreationDateTimePrecision != "" {
			dest["ApproximateCreationDateTimePrecision"] = d.ApproximateCreationDateTimePrecision
		}
		destinations = append(destinations, dest)
	}

	return map[string]interface{}{
		"KinesisDataStreamDestinations": destinations,
		"TableName":                     table.Name,
	}, nil
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
