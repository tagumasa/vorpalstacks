package appsync

import (
	"context"
	"encoding/json"
	"fmt"

	"vorpalstacks/internal/common/request"
)

// ExecuteGraphQLMutation runs one GraphQL mutation against the named API
// for an internal cross-service caller (an EventBridge rule target). The
// operation document and the variables object arrive as delivered by the
// caller; AWS_IAM authorisation is satisfied by the internal call itself —
// the pre-execution authoriser treats an internal, key-less request on an
// AWS_IAM-configured API as authorised, matching the documented contract
// that AppSync EventBridge targets use AWS_IAM.
func (s *AppSyncService) ExecuteGraphQLMutation(ctx context.Context, region, apiID, operation string, variablesJSON []byte) error {
	if apiID == "" {
		return fmt.Errorf("appsync mutation requires an API ID")
	}
	if operation == "" {
		return fmt.Errorf("appsync mutation requires a GraphQL operation document")
	}

	var variables map[string]interface{}
	if len(variablesJSON) > 0 {
		if err := json.Unmarshal(variablesJSON, &variables); err != nil {
			return fmt.Errorf("appsync mutation variables must be a JSON object: %w", err)
		}
	}

	body, err := json.Marshal(graphqlRequest{Query: operation, Variables: variables})
	if err != nil {
		return fmt.Errorf("failed to marshal appsync mutation request: %w", err)
	}

	store, err := s.GetStoreForRegion(region)
	if err != nil {
		return fmt.Errorf("failed to get appsync store for region %s: %w", region, err)
	}
	reqCtx := request.NewRequestContext(ctx, s.storageManager, s.accountID, region)

	result, wireErr := s.executeGraphQLCore(ctx, reqCtx, store, &graphqlExecutionInput{
		ApiId: apiID,
		Body:  body,
	})
	if wireErr != nil {
		return fmt.Errorf("appsync mutation rejected by API %s: %s: %s", apiID, wireErr.ErrorType, wireErr.Message)
	}
	if result != nil && len(result.Errors) > 0 {
		first := result.Errors[0]
		return fmt.Errorf("appsync mutation failed on API %s: %s: %s", apiID, first.ErrorType, first.Message)
	}
	return nil
}
