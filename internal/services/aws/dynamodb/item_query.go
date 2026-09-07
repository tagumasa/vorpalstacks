package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// Query retrieves items based on their key condition expression.
func (s *DynamoDBService) Query(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.queryCore(ctx, reqCtx, queryInput{Parameters: req.Parameters})
}

// Scan retrieves all items in the specified table or index.
func (s *DynamoDBService) Scan(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.scanCore(ctx, store, table, req.Parameters)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"Count":        result.Count,
		"ScannedCount": result.ScannedCount,
	}
	if result.IncludeItems {
		resp["Items"] = buildItemsResponse(result.Items)
	}
	if result.LastEvaluatedKey != nil {
		resp["LastEvaluatedKey"] = buildItemResponse(result.LastEvaluatedKey)
	}

	returnConsumedCapacity := getReturnConsumedCapacity(req.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		isLSI := result.IndexName != "" && !isGSI(table, result.IndexName)
		resp["ConsumedCapacity"] = buildConsumedCapacityResponseWithIndex(table.Name, result.IndexName, result.CapacityUnits, isLSI)
	}

	return resp, nil
}
