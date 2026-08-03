// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// revisionMatches compares an ExpectedRevisionId ("v<N>") against the
// current revision number stored on the table. Returns true if they match,
// or an error if the format is invalid.
func revisionMatches(expected string, currentRev int) (bool, error) {
	trimmed := strings.TrimPrefix(expected, "v")
	if trimmed == "" || trimmed == expected {
		// Either "v" alone, or no "v" prefix at all.
		return false, ErrInvalidParameter
	}
	expectedNum, err := strconv.Atoi(trimmed)
	if err != nil {
		return false, ErrInvalidParameter
	}
	return expectedNum == currentRev, nil
}

// DeleteResourcePolicy deletes a resource policy from a DynamoDB table.
func (s *DynamoDBService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if err := validateResourceArnString(resourceArn); err != nil {
		return nil, err
	}

	expectedRevisionId := request.GetStringParam(req.Parameters, "ExpectedRevisionId")
	if err := validatePolicyRevisionId(expectedRevisionId); err != nil {
		return nil, err
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if expectedRevisionId != "" {
		currentRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
		if revErr != nil {
			return nil, revErr
		}
		matched, matchErr := revisionMatches(expectedRevisionId, currentRev)
		if matchErr != nil {
			return nil, matchErr
		}
		if !matched {
			return nil, ErrPolicyNotFound
		}
	}

	if err := store.Tables().DeleteResourcePolicy(tableName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetResourcePolicy returns the resource policy for a DynamoDB table.
func (s *DynamoDBService) GetResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if err := validateResourceArnString(resourceArn); err != nil {
		return nil, err
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	policy, err := store.Tables().GetResourcePolicy(tableName)
	if err != nil {
		return nil, err
	}

	if policy == "" {
		return nil, ErrPolicyNotFound
	}

	return map[string]interface{}{
		"Policy": policy,
	}, nil
}

// PutResourcePolicy creates or updates a resource policy for a DynamoDB table.
func (s *DynamoDBService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if err := validateResourceArnString(resourceArn); err != nil {
		return nil, err
	}

	policy, ok := req.Parameters["Policy"].(string)
	if !ok || policy == "" {
		return nil, ErrInvalidParameter
	}

	expectedRevisionId := request.GetStringParam(req.Parameters, "ExpectedRevisionId")
	if err := validatePolicyRevisionId(expectedRevisionId); err != nil {
		return nil, err
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	_, err = store.Tables().Get(tableName)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	if expectedRevisionId != "" {
		currentRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
		if revErr != nil {
			return nil, revErr
		}
		matched, matchErr := revisionMatches(expectedRevisionId, currentRev)
		if matchErr != nil {
			return nil, matchErr
		}
		if !matched {
			return nil, ErrPolicyNotFound
		}
	}

	if err := store.Tables().SetResourcePolicy(tableName, policy); err != nil {
		return nil, err
	}

	newRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
	if revErr != nil {
		return nil, revErr
	}
	revisionId := fmt.Sprintf("v%d", newRev)

	return map[string]interface{}{
		"RevisionId": revisionId,
	}, nil
}
