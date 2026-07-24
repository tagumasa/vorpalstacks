// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
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
		DestinationStatus: "ENABLING",
	}
	destinations := append(table.KinesisDataStreamDestinations, newDestination)

	if err := store.Tables().SetKinesisStreamingDestination(tableName, destinations); err != nil {
		return nil, err
	}

	// Background transition to ACTIVE.
	go func() {
		time.Sleep(1 * time.Second)
		for _, d := range destinations {
			if d.StreamArn == streamArn {
				d.DestinationStatus = "ACTIVE"
				break
			}
		}
		if err := store.Tables().SetKinesisStreamingDestination(tableName, destinations); err != nil {
			logs.Error("Failed to transition Kinesis destination to ACTIVE",
				logs.Err(err),
				logs.String("tableName", tableName),
			)
		}
	}()

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": "ENABLING",
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

	// Set DISABLING status first, then remove in background.

	// Keep the destination with DISABLING status for the response.
	disableDestinations := append(table.KinesisDataStreamDestinations[:0:0], table.KinesisDataStreamDestinations...)
	for _, d := range disableDestinations {
		if d.StreamArn == streamArn {
			d.DestinationStatus = "DISABLING"
			break
		}
	}
	if err := store.Tables().SetKinesisStreamingDestination(tableName, disableDestinations); err != nil {
		return nil, err
	}

	// Background transition: remove destination after brief delay.
	go func() {
		time.Sleep(1 * time.Second)
		if err := store.Tables().SetKinesisStreamingDestination(tableName, remainingDestinations); err != nil {
			logs.Error("Failed to remove Kinesis destination after DISABLING",
				logs.Err(err),
				logs.String("tableName", tableName),
			)
		}
	}()

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": "DISABLING",
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

	var precision int
	if configMap, ok := req.Parameters["UpdateKinesisStreamingConfiguration"].(map[string]interface{}); ok {
		if p, ok := configMap["ApproximateCreationDateTimePrecision"].(float64); ok {
			precision = int(p)
		}
	}

	found := false
	for _, d := range table.KinesisDataStreamDestinations {
		if d.StreamArn == streamArn {
			found = true
			if precision > 0 {
				d.ApproximateCreationDateTimePrecision = precision
			}
			break
		}
	}

	if !found {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().SetKinesisStreamingDestination(tableName, table.KinesisDataStreamDestinations); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TableName":         tableName,
		"StreamArn":         streamArn,
		"DestinationStatus": "ACTIVE",
	}, nil
}
