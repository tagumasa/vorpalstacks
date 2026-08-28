// Transport-agnostic Core functions for IAM SSH public keys: validation
// and store operations shared by the AWS-compatible HTTP API handlers and
// any admin plane paths (the xxxCore pattern).
package iam

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"vorpalstacks/internal/common/pagination"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// UploadSSHPublicKeyInput holds the parameters for uploading an SSH public
// key.
type UploadSSHPublicKeyInput struct {
	UserName         string
	SSHPublicKeyBody string
}

// GetSSHPublicKeyInput holds the parameters for retrieving an SSH public
// key.
type GetSSHPublicKeyInput struct {
	UserName       string
	SSHPublicKeyId string
	Encoding       string
}

// UpdateSSHPublicKeyInput holds the parameters for updating an SSH public
// key status.
type UpdateSSHPublicKeyInput struct {
	UserName       string
	SSHPublicKeyId string
	Status         string
}

// SSHPublicKeyListResult holds a paginated SSH public key listing.
type SSHPublicKeyListResult struct {
	Keys        []*iamstore.SSHPublicKey
	IsTruncated bool
	NextMarker  string
}

// uploadSSHPublicKeyCore validates input and uploads an SSH public key for
// the specified IAM user.  The key body is canonicalised before storage.
func (s *IAMService) uploadSSHPublicKeyCore(store *iamstore.IAMStore, input *UploadSSHPublicKeyInput) (*iamstore.SSHPublicKey, error) {
	if input.UserName == "" {
		return nil, NewValidationError("UserName")
	}
	sshPublicKeyBody := input.SSHPublicKeyBody
	if sshPublicKeyBody == "" {
		return nil, NewValidationError("SSHPublicKeyBody")
	}
	// publicKeyMaterialType carries a Latin-1 pattern, so lengths count
	// Unicode characters.
	if utf8.RuneCountInString(sshPublicKeyBody) > maxSSHPublicKeyLength {
		return nil, NewInvalidInputError("SSHPublicKeyBody", fmt.Sprintf("must be 1 to %d characters", maxSSHPublicKeyLength))
	}

	parsedKey, err := parseSSHPublicKey(sshPublicKeyBody)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}
	canonicalBody := canonicalSSHPublicKeyBody(parsedKey)

	if !store.Users().Exists(input.UserName) {
		return nil, NewNoSuchUserError(input.UserName)
	}

	key, err := store.SSHPublicKeys().UploadWithGuards(input.UserName, canonicalBody)
	if err != nil {
		if errors.Is(err, iamstore.ErrDuplicateSSHPublicKey) {
			return nil, ErrDuplicateSSHPublicKey
		}
		if errors.Is(err, iamstore.ErrSSHPublicKeyLimitExceeded) {
			return nil, ErrLimitExceededSSHPublicKeys
		}
		return nil, err
	}
	return key, nil
}

// getSSHPublicKeyCore validates input and retrieves the specified SSH
// public key.
func (s *IAMService) getSSHPublicKeyCore(store *iamstore.IAMStore, input *GetSSHPublicKeyInput) (*iamstore.SSHPublicKey, error) {
	if input.SSHPublicKeyId == "" {
		return nil, NewValidationError("SSHPublicKeyId")
	}
	if input.UserName == "" {
		return nil, NewValidationError("UserName")
	}
	encoding := input.Encoding
	if encoding == "" {
		return nil, NewValidationError("Encoding")
	}
	if encoding != "SSH" && encoding != "PEM" {
		return nil, NewInvalidInputError("Encoding", "must be SSH or PEM")
	}

	if !store.Users().Exists(input.UserName) {
		return nil, NewNoSuchUserError(input.UserName)
	}

	key, err := store.SSHPublicKeys().Get(input.SSHPublicKeyId)
	if err != nil {
		return nil, NewNoSuchEntityError("SSH public key", input.SSHPublicKeyId)
	}
	// The key is retrieved scoped to the named user; a key owned by
	// another user is reported as not existing.
	if key.UserName != input.UserName {
		return nil, NewNoSuchEntityError("SSH public key", input.SSHPublicKeyId)
	}
	return key, nil
}

// updateSSHPublicKeyCore validates input and changes the status of the
// specified SSH public key to Active or Inactive.  The named user must own
// the key; otherwise the operation reports the key as not existing for that
// user.
func (s *IAMService) updateSSHPublicKeyCore(store *iamstore.IAMStore, input *UpdateSSHPublicKeyInput) error {
	if input.SSHPublicKeyId == "" {
		return NewValidationError("SSHPublicKeyId")
	}
	if input.UserName == "" {
		return NewValidationError("UserName")
	}
	status := input.Status
	if status == "" {
		return NewValidationError("Status")
	}
	if status != "Active" && status != "Inactive" {
		return NewInvalidInputError("Status", "must be Active or Inactive")
	}

	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}
	key, err := store.SSHPublicKeys().Get(input.SSHPublicKeyId)
	if err != nil {
		return NewNoSuchEntityError("SSH public key", input.SSHPublicKeyId)
	}
	// The named user must own the key; otherwise the operation reports the
	// key as not existing for that user.
	if key.UserName != input.UserName {
		return NewNoSuchEntityError("SSH public key", input.SSHPublicKeyId)
	}

	return store.SSHPublicKeys().UpdateStatus(input.SSHPublicKeyId, status)
}

// listSSHPublicKeysCore returns a paginated list of the SSH public keys
// associated with the specified IAM user.
func (s *IAMService) listSSHPublicKeysCore(store *iamstore.IAMStore, userName, marker string, maxItems int) (*SSHPublicKeyListResult, error) {
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	keys, err := store.SSHPublicKeys().ListByUserName(userName)
	if err != nil {
		return nil, err
	}

	paged := pagination.PaginateSlice(keys, marker, maxItems, func(k *iamstore.SSHPublicKey) string {
		return k.SSHPublicKeyId
	})

	return &SSHPublicKeyListResult{
		Keys:        paged.Items,
		IsTruncated: paged.IsTruncated,
		NextMarker:  paged.NextMarker,
	}, nil
}

// deleteSSHPublicKeyCore validates input and deletes the specified SSH
// public key for the specified IAM user.  The named user must own the key;
// otherwise the operation reports the key as not existing for that user.
func (s *IAMService) deleteSSHPublicKeyCore(store *iamstore.IAMStore, input *UpdateSSHPublicKeyInput) error {
	if input.SSHPublicKeyId == "" {
		return NewValidationError("SSHPublicKeyId")
	}
	if input.UserName == "" {
		return NewValidationError("UserName")
	}

	if !store.Users().Exists(input.UserName) {
		return NewNoSuchUserError(input.UserName)
	}
	key, err := store.SSHPublicKeys().Get(input.SSHPublicKeyId)
	if err != nil {
		return NewNoSuchEntityError("SSH public key", input.SSHPublicKeyId)
	}
	// The named user must own the key; otherwise the operation reports the
	// key as not existing for that user.
	if key.UserName != input.UserName {
		return NewNoSuchEntityError("SSH public key", input.SSHPublicKeyId)
	}

	return store.SSHPublicKeys().Delete(input.SSHPublicKeyId)
}
