package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Resource Policy Core — single validation + persistence path
//
// These methods encapsulate resource-policy lifecycle logic. Both the HTTP
// API handlers (resource_policy_operations.go) and any future admin handler
// delegate to these methods to ensure identical behaviour.
// ---------------------------------------------------------------------------

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

// resolveResourcePolicyTable validates the request ARN and resolves it to an
// existing table, acquiring the store on the way.
func (s *DynamoDBService) resolveResourcePolicyTable(reqCtx *request.RequestContext, resourceArn string) (dbstore.DynamoDBStoreInterface, string, error) {
	if !validateResourceArnString(resourceArn) {
		return nil, "", ErrInvalidParameter
	}

	tableName := svcarn.ParseTableARN(resourceArn)
	if tableName == "" {
		return nil, "", ErrResourceNotFound
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}
	if _, err := store.Tables().Get(tableName); err != nil {
		return nil, "", ErrResourceNotFound
	}
	return store, tableName, nil
}

// GetResourcePolicyInput is the service-layer DTO for GetResourcePolicy.
type GetResourcePolicyInput struct {
	ResourceArn string
}

// GetResourcePolicyResult is the service-layer result of GetResourcePolicy.
type GetResourcePolicyResult struct {
	Policy     string
	RevisionId string
}

// getResourcePolicyCore returns the resource-based policy for the table
// named by the request ARN. Returns ErrPolicyNotFound when no policy is
// attached.
func (s *DynamoDBService) getResourcePolicyCore(ctx context.Context, reqCtx *request.RequestContext, in GetResourcePolicyInput) (*GetResourcePolicyResult, error) {
	store, tableName, err := s.resolveResourcePolicyTable(reqCtx, in.ResourceArn)
	if err != nil {
		return nil, err
	}

	policy, err := store.Tables().GetResourcePolicy(tableName)
	if err != nil {
		return nil, err
	}
	if policy == "" {
		return nil, ErrPolicyNotFound
	}
	rev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
	if revErr != nil {
		return nil, revErr
	}
	return &GetResourcePolicyResult{
		Policy:     policy,
		RevisionId: fmt.Sprintf("v%d", rev),
	}, nil
}

// PutResourcePolicyInput is the service-layer DTO for PutResourcePolicy.
type PutResourcePolicyInput struct {
	ResourceArn        string
	Policy             string
	ExpectedRevisionId string // optional; empty skips optimistic-lock check
}

// PutResourcePolicyResult is the service-layer result of PutResourcePolicy.
type PutResourcePolicyResult struct {
	RevisionId string
}

// putResourcePolicyCore creates or replaces the resource-based policy for
// the table named by the request ARN. When ExpectedRevisionId is non-empty,
// it must match the current revision or ErrPolicyNotFound is returned.
func (s *DynamoDBService) putResourcePolicyCore(ctx context.Context, reqCtx *request.RequestContext, in PutResourcePolicyInput) (*PutResourcePolicyResult, error) {
	if in.Policy == "" {
		return nil, ErrInvalidParameter
	}

	if !validatePolicyRevisionId(in.ExpectedRevisionId) {
		return nil, ErrInvalidParameter
	}

	store, tableName, err := s.resolveResourcePolicyTable(reqCtx, in.ResourceArn)
	if err != nil {
		return nil, err
	}

	if in.ExpectedRevisionId != "" {
		currentRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
		if revErr != nil {
			return nil, revErr
		}
		matched, matchErr := revisionMatches(in.ExpectedRevisionId, currentRev)
		if matchErr != nil {
			return nil, matchErr
		}
		if !matched {
			return nil, ErrPolicyNotFound
		}
	}

	if err := store.Tables().SetResourcePolicy(tableName, in.Policy); err != nil {
		return nil, err
	}

	newRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
	if revErr != nil {
		return nil, revErr
	}
	return &PutResourcePolicyResult{
		RevisionId: fmt.Sprintf("v%d", newRev),
	}, nil
}

// DeleteResourcePolicyInput is the service-layer DTO for
// DeleteResourcePolicy.
type DeleteResourcePolicyInput struct {
	ResourceArn        string
	ExpectedRevisionId string // optional; empty skips optimistic-lock check
}

// deleteResourcePolicyCore removes the resource-based policy from the table
// named by the request ARN. When ExpectedRevisionId is non-empty, it must
// match the current revision or ErrPolicyNotFound is returned.
func (s *DynamoDBService) deleteResourcePolicyCore(ctx context.Context, reqCtx *request.RequestContext, in DeleteResourcePolicyInput) error {
	if !validatePolicyRevisionId(in.ExpectedRevisionId) {
		return ErrInvalidParameter
	}

	store, tableName, err := s.resolveResourcePolicyTable(reqCtx, in.ResourceArn)
	if err != nil {
		return err
	}

	if in.ExpectedRevisionId != "" {
		currentRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
		if revErr != nil {
			return revErr
		}
		matched, matchErr := revisionMatches(in.ExpectedRevisionId, currentRev)
		if matchErr != nil {
			return matchErr
		}
		if !matched {
			return ErrPolicyNotFound
		}
	}
	return store.Tables().DeleteResourcePolicy(tableName)
}
