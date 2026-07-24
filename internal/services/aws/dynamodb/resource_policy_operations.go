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
// current revision number stored on the table. Returns true if they match.
func revisionMatches(expected string, currentRev int) bool {
	expectedNum, err := strconv.Atoi(strings.TrimPrefix(expected, "v"))
	if err != nil {
		return false
	}
	return expectedNum == currentRev
}

// DeleteResourcePolicy deletes a resource policy from a DynamoDB table.
func (s *DynamoDBService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetStringParam(req.Parameters, "ResourceArn")
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	expectedRevisionId := request.GetStringParam(req.Parameters, "ExpectedRevisionId")

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
		currentRev, _ := store.Tables().GetResourcePolicyRevisionId(tableName)
		if !revisionMatches(expectedRevisionId, currentRev) {
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
	if resourceArn == "" {
		return nil, ErrInvalidParameter
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
	if resourceArn == "" {
		return nil, ErrInvalidParameter
	}

	policy, ok := req.Parameters["Policy"].(string)
	if !ok || policy == "" {
		return nil, ErrInvalidParameter
	}

	expectedRevisionId := request.GetStringParam(req.Parameters, "ExpectedRevisionId")

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
		currentRev, _ := store.Tables().GetResourcePolicyRevisionId(tableName)
		if !revisionMatches(expectedRevisionId, currentRev) {
			return nil, ErrPolicyNotFound
		}
	}

	if err := store.Tables().SetResourcePolicy(tableName, policy); err != nil {
		return nil, err
	}

	newRev, _ := store.Tables().GetResourcePolicyRevisionId(tableName)
	revisionId := fmt.Sprintf("v%d", newRev)

	return map[string]interface{}{
		"RevisionId": revisionId,
	}, nil
}
