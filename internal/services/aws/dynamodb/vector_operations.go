package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// SearchVectors performs a vector similarity search on a vector index.
// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_SearchVectors.html
func (s *DynamoDBService) SearchVectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.searchVectorsCore(ctx, store, req.Parameters)
	if err != nil {
		return nil, err
	}

	searchResults := make([]map[string]interface{}, len(result.Items))
	for i, item := range result.Items {
		searchResults[i] = map[string]interface{}{
			"Item":  item,
			"Score": result.Scores[i],
		}
	}
	resp := map[string]interface{}{
		"SearchResults": searchResults,
	}
	if getReturnConsumedCapacity(req.Parameters) != "NONE" {
		resp["ConsumedCapacity"] = map[string]interface{}{
			"VectorSearchRequestBytes": float64(result.VectorSearchRequestBytes),
		}
	}
	return resp, nil
}
