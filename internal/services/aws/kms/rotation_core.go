package kms

// rotation_core.go carries the Core functions of the KMS key-rotation
// family. The state/eligibility gates run before any store or HSM
// mutation, matching the original failure precedence; the on-demand
// rotation enforces the AWS per-key limit before rotating.

import (
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// EnableKeyRotationInput carries the EnableKeyRotation members. The
// rotation-period presence flag distinguishes an explicitly provided
// value (range-checked in the Core) from the documented default of 365
// days.
type EnableKeyRotationInput struct {
	RotationPeriodInDays int32
	HasRotationPeriod    bool
}

// enableKeyRotationCore is the single entry point for enabling automatic
// key rotation. Per AWS spec, automatic rotation is supported only on
// SYMMETRIC_DEFAULT keys with AWS_KMS origin; the key must be enabled and
// not pending deletion or import. RotationPeriodInDays is optional,
// range 90-2560, default 365.
func (s *KMSService) enableKeyRotationCore(stores *kmsStores, key *kmsstore.Key, in EnableKeyRotationInput) error {
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStatePendingImport {
		// AWS: PendingImport is a distinct KMSInvalidStateException. The
		// previous code conflated the two states under ErrKeyPendingDeletion,
		// which is misleading for PendingImport keys.
		return ErrKeyPendingImport
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return ErrKeyDisabled
	}
	if err := validateRotationKeyEligibility(key.KeySpec, key.Origin); err != nil {
		return err
	}

	rotationPeriod := defaultRotationPeriodInDays
	if in.HasRotationPeriod {
		rotationPeriod = in.RotationPeriodInDays
		if rotationPeriod < 90 || rotationPeriod > 2560 {
			return ErrValidation
		}
	}

	return stores.keys.SetKeyRotationWithPeriod(key.KeyID, true, rotationPeriod)
}

// disableKeyRotationCore is the single entry point for disabling automatic
// key rotation. Per AWS spec, rotation status cannot be changed while the
// key is pending deletion or disabled.
func (s *KMSService) disableKeyRotationCore(stores *kmsStores, key *kmsstore.Key) error {
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStatePendingImport {
		return ErrKeyPendingImport
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return ErrKeyDisabled
	}
	// DisableKeyRotation must apply the same KeySpec/Origin
	// eligibility check as EnableKeyRotation. Without this, callers
	// can invoke DisableKeyRotation on asymmetric/HMAC keys without
	// error, violating the API contract.
	if err := validateRotationKeyEligibility(key.KeySpec, key.Origin); err != nil {
		return err
	}

	return stores.keys.SetKeyRotation(key.KeyID, false)
}

// rotateKeyOnDemandCore is the single entry point for triggering an
// immediate rotation. Only SYMMETRIC_DEFAULT keys with AWS_KMS origin are
// eligible; AWS enforces a limit of 25 on-demand rotations per key. The
// key record is returned with the rotation history appended so the caller
// can serialise the key ARN.
func (s *KMSService) rotateKeyOnDemandCore(stores *kmsStores, key *kmsstore.Key) (*kmsstore.Key, error) {
	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStatePendingImport {
		return nil, ErrKeyPendingImport
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return nil, ErrKeyDisabled
	}
	if key.KeySpec != kmsstore.KeySpecSymmetricDefault || key.Origin != kmsstore.OriginTypeAWSKMS {
		return nil, ErrUnsupportedOperation
	}

	if key.OnDemandRotationCount >= maxOnDemandRotations {
		return nil, awserrors.NewAWSError("LimitExceededException",
			fmt.Sprintf("Maximum on-demand rotations (%d) exceeded for this key", maxOnDemandRotations),
			http.StatusTooManyRequests)
	}

	// Rotate the key material in the HSM. Previous versions are preserved
	// so that ciphertexts encrypted before rotation remain decryptable.
	if err := s.hsmBackend.RotateKey(key.KeyID); err != nil {
		return nil, err
	}

	// Record the rotation event.
	now := time.Now().UTC()
	key.OnDemandRotationCount++
	key.OnDemandRotationStartDate = &now
	key.RotationHistory = append(key.RotationHistory, kmsstore.RotationEntry{
		RotationDate:  now,
		RotationType:  "ON_DEMAND",
		KeyMaterialId: fmt.Sprintf("v%d", key.OnDemandRotationCount),
	})
	if err := stores.keys.Update(key); err != nil {
		return nil, err
	}

	return key, nil
}
