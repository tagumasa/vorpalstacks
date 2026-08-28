// Transport-agnostic Core functions for IAM access keys: validation and
// store operations shared by the AWS-compatible HTTP API handlers and any
// admin plane paths (the xxxCore pattern).
package iam

import (
	"errors"

	"vorpalstacks/internal/common/pagination"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// UpdateAccessKeyInput holds the parameters for UpdateAccessKey.  UserName
// is optional: when set it must own the key.
type UpdateAccessKeyInput struct {
	AccessKeyId string
	UserName    string
	Status      string
}

// AccessKeyListResult holds a paginated access key listing.
type AccessKeyListResult struct {
	Keys        []*iamstore.AccessKey
	IsTruncated bool
	NextMarker  string
}

// createAccessKeyCore validates that the user exists and creates a new
// access key.  The per-user quota is enforced atomically inside the store
// layer to prevent race conditions on concurrent requests.
func (s *IAMService) createAccessKeyCore(store *iamstore.IAMStore, userName string) (*iamstore.AccessKey, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	key, err := store.AccessKeys().CreateWithLimit(userName, iamstore.MaxAccessKeysPerUser)
	if err != nil {
		if errors.Is(err, iamstore.ErrAccessKeyLimitExceeded) {
			return nil, ErrAccessKeyLimitExceeded
		}
		return nil, err
	}
	return key, nil
}

// deleteAccessKeyCore validates input and deletes the specified access key.
// A non-empty userName that does not own the key yields NoSuchEntity.
func (s *IAMService) deleteAccessKeyCore(store *iamstore.IAMStore, accessKeyId, userName string) error {
	if accessKeyId == "" {
		return NewValidationError("AccessKeyId")
	}

	key, err := store.AccessKeys().Get(accessKeyId)
	if err != nil {
		return NewNoSuchAccessKeyError(accessKeyId)
	}

	if userName != "" && key.UserName != userName {
		return NewNoSuchAccessKeyError(accessKeyId)
	}

	return store.AccessKeys().Delete(accessKeyId)
}

// listAccessKeysCore returns a paginated list of the access keys owned by
// the specified user.
func (s *IAMService) listAccessKeysCore(store *iamstore.IAMStore, userName, marker string, maxItems int) (*AccessKeyListResult, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	keys, err := store.AccessKeys().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	result := pagination.PaginateSlice(keys, marker, maxItems, func(k *iamstore.AccessKey) string {
		return k.AccessKeyId
	})

	return &AccessKeyListResult{
		Keys:        result.Items,
		IsTruncated: result.IsTruncated,
		NextMarker:  result.NextMarker,
	}, nil
}

// getAccessKeyLastUsedCore retrieves the access key with its last-used
// tracking fields.
func (s *IAMService) getAccessKeyLastUsedCore(store *iamstore.IAMStore, accessKeyId string) (*iamstore.AccessKey, error) {
	if accessKeyId == "" {
		return nil, NewValidationError("AccessKeyId")
	}

	key, err := store.AccessKeys().Get(accessKeyId)
	if err != nil {
		return nil, NewNoSuchAccessKeyError(accessKeyId)
	}
	return key, nil
}

// updateAccessKeyCore validates input and updates the status of the
// specified access key.  A non-empty userName that does not own the key
// yields NoSuchEntity.
func (s *IAMService) updateAccessKeyCore(store *iamstore.IAMStore, input *UpdateAccessKeyInput) error {
	if input.AccessKeyId == "" {
		return NewValidationError("AccessKeyId")
	}

	key, err := store.AccessKeys().Get(input.AccessKeyId)
	if err != nil {
		return NewNoSuchAccessKeyError(input.AccessKeyId)
	}

	if input.UserName != "" && key.UserName != input.UserName {
		return NewNoSuchAccessKeyError(input.AccessKeyId)
	}

	newStatus := iamstore.AccessKeyStatus(input.Status)
	if newStatus != iamstore.AccessKeyStatusActive && newStatus != iamstore.AccessKeyStatusInactive {
		return NewInvalidInputError("Status", "must be Active or Inactive")
	}

	return store.AccessKeys().UpdateStatus(input.AccessKeyId, newStatus)
}
