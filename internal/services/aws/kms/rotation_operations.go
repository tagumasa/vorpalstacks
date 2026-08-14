package kms

// Package kms provides KMS (Key Management Service) operations for vorpalstacks.

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// defaultRotationPeriodInDays is the AWS default rotation period when
// EnableKeyRotation is called without an explicit RotationPeriodInDays.
const defaultRotationPeriodInDays = int32(365)

// EnableKeyRotation enables automatic key rotation for a symmetric key.
// Key rotation automatically generates new cryptographic material every year.
// Per AWS spec, automatic rotation is supported only on SYMMETRIC_DEFAULT keys
// with AWS_KMS origin. The key must be enabled and not pending deletion or import.
func (s *KMSService) EnableKeyRotation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "EnableKeyRotation", nil)
	if err != nil {
		return nil, err
	}

	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStatePendingImport {
		// AWS: PendingImport is a distinct KMSInvalidStateException. The
		// previous code conflated the two states under ErrKeyPendingDeletion,
		// which is misleading for PendingImport keys.
		return nil, ErrKeyPendingImport
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return nil, ErrKeyDisabled
	}
	if err := validateRotationKeyEligibility(key.KeySpec, key.Origin); err != nil {
		return nil, err
	}

	// AWS: RotationPeriodInDays is optional, range 90-2560, default 365.
	rotationPeriod := defaultRotationPeriodInDays
	if _, ok := req.Parameters["RotationPeriodInDays"]; ok {
		rotationPeriod = int32(request.GetIntParam(req.Parameters, "RotationPeriodInDays"))
		if rotationPeriod < 90 || rotationPeriod > 2560 {
			return nil, ErrValidation
		}
	}

	if err := stores.keys.SetKeyRotationWithPeriod(key.KeyID, true, rotationPeriod); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableKeyRotation disables automatic key rotation for a symmetric key.
// Per AWS spec, rotation status cannot be changed while the key is pending
// deletion or disabled.
func (s *KMSService) DisableKeyRotation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "DisableKeyRotation", nil)
	if err != nil {
		return nil, err
	}

	if key.KeyState == kmsstore.KeyStatePendingDeletion {
		return nil, ErrKeyPendingDeletion
	}
	if key.KeyState == kmsstore.KeyStatePendingImport {
		return nil, ErrKeyPendingImport
	}
	if key.KeyState == kmsstore.KeyStateDisabled || !key.Enabled {
		return nil, ErrKeyDisabled
	}
	// DisableKeyRotation must apply the same KeySpec/Origin
	// eligibility check as EnableKeyRotation. Without this, callers
	// can invoke DisableKeyRotation on asymmetric/HMAC keys without
	// error, violating the API contract.
	if err := validateRotationKeyEligibility(key.KeySpec, key.Origin); err != nil {
		return nil, err
	}

	if err := stores.keys.SetKeyRotation(key.KeyID, false); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetKeyRotationStatus returns the status of key rotation for a specified key.
// It indicates whether automatic key rotation is enabled and the next rotation date.
func (s *KMSService) GetKeyRotationStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GetKeyRotationStatus", nil)
	if err != nil {
		return nil, err
	}

	rotationPeriod := key.RotationPeriodInDays
	if rotationPeriod == 0 {
		rotationPeriod = defaultRotationPeriodInDays
	}

	response := map[string]interface{}{
		"KeyRotationEnabled":   key.KeyRotationEnabled,
		"KeyId":                key.Arn,
		"RotationPeriodInDays": rotationPeriod,
	}
	if key.OnDemandRotationStartDate != nil {
		response["OnDemandRotationStartDate"] = key.OnDemandRotationStartDate.Unix()
	}

	// Compute NextRotationDate per the AWS documented behaviour:
	// "approximately RotationPeriodInDays after the key was enabled
	// for automatic rotation, or after the key was last rotated,
	// whichever is most recent." Using CreationDate as the base
	// produced past-due dates for keys created long before rotation
	// was enabled.
	if key.KeyRotationEnabled && !key.KeyRotationEnabledAt.IsZero() {
		base := key.KeyRotationEnabledAt
		if len(key.RotationHistory) > 0 {
			lastRotation := key.RotationHistory[len(key.RotationHistory)-1].RotationDate
			if lastRotation.After(base) {
				base = lastRotation
			}
		}
		nextRotation := base.AddDate(0, 0, int(rotationPeriod))
		response["NextRotationDate"] = nextRotation.Unix()
	}

	return response, nil
}

// maxOnDemandRotations is the AWS limit for on-demand rotations per key.
const maxOnDemandRotations = 25

// RotateKeyOnDemand triggers an immediate key rotation for a symmetric
// KMS key. Per AWS spec, only SYMMETRIC_DEFAULT keys with AWS_KMS origin
// are eligible. The key must be enabled and not pending deletion or
// import. AWS enforces a limit of 25 on-demand rotations per key.
func (s *KMSService) RotateKeyOnDemand(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "RotateKeyOnDemand", nil)
	if err != nil {
		return nil, err
	}

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

	return map[string]interface{}{
		"KeyId": key.Arn,
	}, nil
}

// GetKeyLastUsage returns the timestamp and operation of the most recent
// cryptographic operation that used the specified key. Read-only operations
// (DescribeKey, ListResourceTags) do not count as usage.
//
// Per the Smithy GetKeyLastUsageResponse shape, the response contains:
//   - KeyId (string)
//   - KeyLastUsage (KeyLastUsageData: Timestamp, Operation)
//   - KeyCreationDate (date)
//
// CloudTrailEventId and KmsRequestId from KeyLastUsageData are omitted
// because CloudTrail event correlation is out of scope.
func (s *KMSService) GetKeyLastUsage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	key, err := s.resolveAndAuthorizeKey(reqCtx, req, stores, "GetKeyLastUsage", nil)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"KeyId": key.Arn,
	}
	if !key.CreationDate.IsZero() {
		result["KeyCreationDate"] = key.CreationDate.Unix()
		// Per the Smithy GetKeyLastUsageResponse shape,
		// TrackingStartDate is the later of the key creation date or the
		// date KMS began recording activity. For this platform, usage
		// tracking starts at key creation, so TrackingStartDate equals
		// KeyCreationDate.
		result["TrackingStartDate"] = key.CreationDate.Unix()
	}
	if key.LastUsedAt != nil {
		result["KeyLastUsage"] = map[string]interface{}{
			"Timestamp": key.LastUsedAt.Unix(),
			"Operation": key.LastUsageOperation,
		}
	}
	return result, nil
}
