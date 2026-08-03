package kms

import (
	stderrors "errors"
	"fmt"

	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/kms/hsm"
	kmsstore "vorpalstacks/internal/store/aws/kms"
	"vorpalstacks/internal/utils/aws/types"
)

// ---------------------------------------------------------------------------
// Transport-agnostic DTOs for CreateKey
// ---------------------------------------------------------------------------

// CreateKeyInput carries every field that CreateKey needs, in a format
// independent of the wire protocol. Both the HTTP API handler and the
// admin gRPC handler build this struct and delegate to createKeyCore,
// ensuring that validation, key creation, and policy handling follow a
// single code path.
type CreateKeyInput struct {
	Description        string
	KeyUsage           kmsstore.KeyUsage
	KeySpec            kmsstore.KeySpec
	Origin             kmsstore.OriginType
	MultiRegion        bool
	CustomKeyStoreID   string
	XksKeyID           string
	Policy             string
	BypassLockoutCheck bool
	Tags               []types.Tag
	AccountID          string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createKeyCore is the single entry point for key creation shared by the
// HTTP API and the admin gRPC handler. It applies defaults, performs all
// validation, generates the key material (HSM or import-pending), persists
// the key policy and tags, and returns the fully-created key. On any
// failure after key creation it cleans up partial state via
// CascadeDeleteKey.
func (s *KMSService) createKeyCore(stores *kmsStores, in CreateKeyInput) (*kmsstore.Key, error) {
	// 1. Apply defaults (Smithy: KeyUsage/KeySpec/Origin are OPTIONAL).
	if in.KeyUsage == "" {
		in.KeyUsage = kmsstore.KeyUsageEncryptDecrypt
	}
	if in.KeySpec == "" {
		in.KeySpec = kmsstore.KeySpecSymmetricDefault
	}
	if in.Origin == "" {
		in.Origin = kmsstore.OriginTypeAWSKMS
	}

	// 2. AWS_CLOUDHSM origin requires a CustomKeyStoreId pointing at a real
	// CloudHSM custom key store. CloudHSM integration is not implemented,
	// so reject AWS_CLOUDHSM explicitly rather than silently storing a key
	// that no HSM backend can serve.
	if in.Origin == kmsstore.OriginTypeAWSCloudHSM {
		return nil, ErrUnsupportedOperation
	}

	// 3. CustomKeyStoreId and XksKeyId must not be silently dropped.
	if err := validateOriginParams(in.CustomKeyStoreID, in.XksKeyID); err != nil {
		return nil, err
	}

	// 4. KeyUsage/KeySpec combination validation (Smithy matrix).
	if err := validateKeyUsageSpecCombo(in.KeyUsage, in.KeySpec); err != nil {
		return nil, err
	}

	// 5. Description length (Smithy DescriptionType: 0-8192).
	if err := validateDescriptionLength(in.Description); err != nil {
		return nil, err
	}

	// 6. Validate tags BEFORE creating any state so that an invalid tag
	// list cannot leave a half-created key behind.
	if len(in.Tags) > 0 {
		if err := validateKMSTags(in.Tags); err != nil {
			return nil, err
		}
	}

	// 7. Policy defaults and size validation (Smithy PolicyType: 1-131072).
	policyStr := in.Policy
	if policyStr == "" {
		policyStr = kmsstore.DefaultKeyPolicy
	}
	if err := validatePolicySize(policyStr); err != nil {
		return nil, err
	}

	// 8. Validate the supplied policy does not lock out the root principal
	// unless BypassPolicyLockoutSafetyCheck is explicitly true.
	if !in.BypassLockoutCheck {
		if err := validatePolicyDoesNotLockOutRoot(policyStr, in.AccountID); err != nil {
			return nil, err
		}
	}

	// 9. Generate key ID.
	keyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, err
	}

	// 10. Create key in store.
	key, err := stores.keys.Create(
		keyID,
		in.KeyUsage,
		in.KeySpec,
		in.Description,
		in.Origin,
		in.MultiRegion,
	)
	if err != nil {
		return nil, err
	}

	// 11. Apply tags.
	if len(in.Tags) > 0 {
		if err := stores.keys.TagStore.Tag(keyID, tagutil.ToMap(in.Tags)); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after Tag failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, err
		}
	}

	// 12. Set up key material.
	if in.Origin == kmsstore.OriginTypeExternal {
		if err := stores.keys.SetPendingImport(keyID); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after SetPendingImport failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, err
		}
		key.KeyState = kmsstore.KeyStatePendingImport
		key.Enabled = false
	} else {
		if err := s.hsmBackend.GenerateKey(keyID, hsm.KeySpec(in.KeySpec)); err != nil {
			if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after HSM GenerateKey failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, err
		}
	}

	// 13. Set bypass flag and persist.
	key.BypassPolicyLockoutSafetyCheck = in.BypassLockoutCheck
	if err := stores.keys.Update(key); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
			logs.Error("Failed to cascade-delete key after Update", logs.Err(delErr), logs.String("keyId", keyID))
		}
		return nil, err
	}

	// 14. Persist key policy.
	if err := stores.keyPolicies.PutDefault(keyID, policyStr); err != nil {
		if delErr := stores.CascadeDeleteKey(s.hsmBackend, keyID); delErr != nil {
			logs.Error("Failed to cascade-delete key after PutDefault policy failure", logs.Err(delErr), logs.String("keyId", keyID))
		}
		return nil, err
	}

	return key, nil
}

// scheduleKeyDeletionCore is the single entry point for scheduling key
// deletion shared by the HTTP API and the admin gRPC handler. It
// validates the pending window (AWS range 7-30), schedules deletion in
// the store, and returns the updated key.
func (s *KMSService) scheduleKeyDeletionCore(stores *kmsStores, keyID string, pendingWindowInDays int) (*kmsstore.Key, int, error) {
	// Validate pending window (AWS: 7-30).
	if pendingWindowInDays < 7 || pendingWindowInDays > 30 {
		return nil, 0, NewValidationError(fmt.Sprintf("PendingWindowInDays %d does not fit range 7-30", pendingWindowInDays))
	}

	// Schedule deletion. Translate the store-layer ErrInvalidKeyState
	// (key is already PendingDeletion) into the service-layer
	// KMSInvalidStateException so that both the HTTP path and the admin
	// gRPC path surface the correct HTTP 400 / CodeInvalidArgument.
	if err := stores.keys.ScheduleDeletion(keyID, pendingWindowInDays); err != nil {
		if stderrors.Is(err, kmsstore.ErrInvalidKeyState) {
			return nil, 0, ErrKeyPendingDeletion
		}
		return nil, 0, err
	}

	// Retrieve updated key.
	updatedKey, err := stores.keys.Get(keyID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve key after scheduling deletion: %w", err)
	}
	if updatedKey == nil {
		return nil, 0, fmt.Errorf("failed to retrieve key after scheduling deletion: store returned nil")
	}

	return updatedKey, pendingWindowInDays, nil
}
