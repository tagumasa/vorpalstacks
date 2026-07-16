// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DescribeKinesisStreamingDestination returns the Kinesis streaming destination for a table.
func (s *DynamoDBService) DescribeKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	var destinations []map[string]interface{}
	for _, d := range table.KinesisDataStreamDestinations {
		destinations = append(destinations, map[string]interface{}{
			"StreamArn":                    d.StreamArn,
			"DestinationStatus":            d.DestinationStatus,
			"DestinationStatusDescription": d.DestinationStatusDescription,
		})
	}

	return map[string]interface{}{
		"KinesisDataStreamDestinations": destinations,
		"TableName":                     table.Name,
	}, nil
}

// EnableKinesisStreamingDestination enables Kinesis streaming for a DynamoDB table.
func (s *DynamoDBService) EnableKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	streamArn := request.GetStringParam(req.Parameters, "StreamArn")
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	for _, d := range table.KinesisDataStreamDestinations {
		if d.StreamArn == streamArn {
			return nil, ErrResourceAlreadyExists
		}
	}

	newDestination := &dbstore.KinesisDataStreamDestination{
		StreamArn:         streamArn,
		DestinationStatus: "ACTIVE",
	}
	destinations := append(table.KinesisDataStreamDestinations, newDestination)

	if err := store.Tables().SetKinesisStreamingDestination(tableName, destinations); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": "ACTIVE",
	}, nil
}

// DisableKinesisStreamingDestination disables Kinesis streaming for a DynamoDB table.
func (s *DynamoDBService) DisableKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	streamArn := request.GetStringParam(req.Parameters, "StreamArn")
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var remainingDestinations []*dbstore.KinesisDataStreamDestination
	found := false
	for _, d := range table.KinesisDataStreamDestinations {
		if d.StreamArn == streamArn {
			found = true
			continue
		}
		remainingDestinations = append(remainingDestinations, d)
	}

	if !found {
		return nil, ErrResourceNotFound
	}

	if err := store.Tables().SetKinesisStreamingDestination(tableName, remainingDestinations); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": "DISABLED",
	}, nil
}

// UpdateKinesisStreamingDestination updates the Kinesis streaming destination for a DynamoDB table.
func (s *DynamoDBService) UpdateKinesisStreamingDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	streamArn := request.GetStringParam(req.Parameters, "StreamArn")
	if streamArn == "" {
		return nil, ErrInvalidParameter
	}

	found := false
	for _, d := range table.KinesisDataStreamDestinations {
		if d.StreamArn == streamArn {
			found = true
			break
		}
	}

	if !found {
		return nil, ErrResourceNotFound
	}

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": "ACTIVE",
	}, nil
}
