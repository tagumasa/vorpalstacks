package dynamodb

import (
	"context"
	"errors"
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

// expectedPolicyRevision converts an already format-validated
// ExpectedRevisionId ("v<N>") into the revision number the store compares
// under the table record's lock. An absent member reports
// PolicyRevisionUnchecked — the wire contract's skip-optimistic-locking
// value.
func expectedPolicyRevision(expected string) (int, error) {
	if expected == "" {
		return dbstore.PolicyRevisionUnchecked, nil
	}
	num, err := strconv.Atoi(strings.TrimPrefix(expected, "v"))
	if err != nil {
		return 0, ErrInvalidParameter
	}
	return num, nil
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

	expectedRev, revErr := expectedPolicyRevision(in.ExpectedRevisionId)
	if revErr != nil {
		return nil, revErr
	}

	// The revision check and the write run inside one locked
	// read-modify-write of the table record: two concurrent writers that
	// both observed the same revision cannot both apply — the loser is
	// told the revision moved, never silently overwritten.
	newRev, setErr := store.Tables().SetResourcePolicyExpected(tableName, in.Policy, expectedRev)
	if setErr != nil {
		if errors.Is(setErr, dbstore.ErrPolicyRevisionMismatch) {
			return nil, ErrPolicyNotFound
		}
		return nil, setErr
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

// DeleteResourcePolicyResult is the service-layer result of
// DeleteResourcePolicy.
type DeleteResourcePolicyResult struct {
	RevisionId string
}

// deleteResourcePolicyCore removes the resource-based policy from the table
// named by the request ARN. When ExpectedRevisionId is non-empty, it must
// match the current revision or ErrPolicyNotFound is returned.
func (s *DynamoDBService) deleteResourcePolicyCore(ctx context.Context, reqCtx *request.RequestContext, in DeleteResourcePolicyInput) (*DeleteResourcePolicyResult, error) {
	if !validatePolicyRevisionId(in.ExpectedRevisionId) {
		return nil, ErrInvalidParameter
	}

	store, tableName, err := s.resolveResourcePolicyTable(reqCtx, in.ResourceArn)
	if err != nil {
		return nil, err
	}

	expectedRev, revErr := expectedPolicyRevision(in.ExpectedRevisionId)
	if revErr != nil {
		return nil, revErr
	}

	// The same one-lock write the put path applies: the check and the
	// delete commit together, and the delete advances the revision so a
	// concurrent put holding the same expected revision cannot apply past
	// it.
	newRev, delErr := store.Tables().DeleteResourcePolicyExpected(tableName, expectedRev)
	if delErr != nil {
		if errors.Is(delErr, dbstore.ErrPolicyRevisionMismatch) {
			return nil, ErrPolicyNotFound
		}
		return nil, delErr
	}
	return &DeleteResourcePolicyResult{
		RevisionId: fmt.Sprintf("v%d", newRev),
	}, nil
}
