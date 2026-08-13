package dynamodb

import (
	"fmt"
	"strconv"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
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

// GetResourcePolicyResult is the service-layer result of GetResourcePolicy.
type GetResourcePolicyResult struct {
	Policy     string
	RevisionId string
}

// getResourcePolicyCore returns the resource-based policy for the named
// table. Returns ErrPolicyNotFound when no policy is attached.
func (s *DynamoDBService) getResourcePolicyCore(store dbstore.DynamoDBStoreInterface, tableName string) (*GetResourcePolicyResult, error) {
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
	TableName          string
	Policy             string
	ExpectedRevisionId string // optional; empty skips optimistic-lock check
}

// PutResourcePolicyResult is the service-layer result of PutResourcePolicy.
type PutResourcePolicyResult struct {
	RevisionId string
}

// putResourcePolicyCore creates or replaces the resource-based policy for
// the named table. When ExpectedRevisionId is non-empty, it must match the
// current revision or ErrPolicyNotFound is returned.
func (s *DynamoDBService) putResourcePolicyCore(store dbstore.DynamoDBStoreInterface, in PutResourcePolicyInput) (*PutResourcePolicyResult, error) {
	if in.ExpectedRevisionId != "" {
		currentRev, revErr := store.Tables().GetResourcePolicyRevisionId(in.TableName)
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

	if err := store.Tables().SetResourcePolicy(in.TableName, in.Policy); err != nil {
		return nil, err
	}

	newRev, revErr := store.Tables().GetResourcePolicyRevisionId(in.TableName)
	if revErr != nil {
		return nil, revErr
	}
	return &PutResourcePolicyResult{
		RevisionId: fmt.Sprintf("v%d", newRev),
	}, nil
}

// deleteResourcePolicyCore removes the resource-based policy from the
// named table. When ExpectedRevisionId is non-empty, it must match the
// current revision or ErrPolicyNotFound is returned.
func (s *DynamoDBService) deleteResourcePolicyCore(store dbstore.DynamoDBStoreInterface, tableName, expectedRevisionId string) error {
	if expectedRevisionId != "" {
		currentRev, revErr := store.Tables().GetResourcePolicyRevisionId(tableName)
		if revErr != nil {
			return revErr
		}
		matched, matchErr := revisionMatches(expectedRevisionId, currentRev)
		if matchErr != nil {
			return matchErr
		}
		if !matched {
			return ErrPolicyNotFound
		}
	}
	return store.Tables().DeleteResourcePolicy(tableName)
}
