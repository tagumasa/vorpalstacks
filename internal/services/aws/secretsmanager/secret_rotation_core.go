package secretsmanager

import (
	"context"
	stderrors "errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/scheduleexpr"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs shared by HTTP API and admin handler
// ---------------------------------------------------------------------------

// RestoreSecretInput carries all fields needed for RestoreSecret.
type RestoreSecretInput struct {
	SecretId string
}

// RotateSecretInput carries all fields needed for RotateSecret.
type RotateSecretInput struct {
	SecretId                       string
	RotationLambdaARN              string
	AutomaticallyAfterDays         int
	ScheduleExpression             string
	Duration                       string
	ClientRequestToken             string
	ExternalSecretRotationMetadata []secretsmanagerstore.ExternalSecretRotationMetadataItem
	ExternalSecretRotationRoleArn  string
	// RotateImmediately defaults to true; HasRotateImmediately=false means
	// the caller left the parameter unset.
	RotateImmediately    bool
	HasRotateImmediately bool
}

// CancelRotateSecretInput carries all fields needed for CancelRotateSecret.
type CancelRotateSecretInput struct {
	SecretId string
}

// ---------------------------------------------------------------------------
// Result structs — transport-agnostic results
// ---------------------------------------------------------------------------

// RestoreSecretResult holds the transport-agnostic result of RestoreSecret.
type RestoreSecretResult struct {
	ARN  string
	Name string
}

// RotateSecretResult holds the transport-agnostic result of RotateSecret.
type RotateSecretResult struct {
	ARN       string
	Name      string
	VersionID string
}

// CancelRotateSecretResult holds the transport-agnostic result of
// CancelRotateSecret.
type CancelRotateSecretResult struct {
	ARN       string
	Name      string
	VersionID string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// restoreSecretCore is the single entry point for restoring a previously
// deleted secret. Secrets Manager keeps deleted secrets recoverable for a
// minimum of 30 days.
func (s *SecretsManagerService) restoreSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in RestoreSecretInput) (*RestoreSecretResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecretForMetadata(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	if err := store.RestoreSecret(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return &RestoreSecretResult{
		ARN:  secret.ARN,
		Name: secret.Name,
	}, nil
}

// rotateSecretCore is the single entry point for configuring rotation and,
// if a RotationLambdaARN is available and RotateImmediately is true,
// performing the full AWS rotation protocol:
//
//  1. Update rotation metadata (RotationEnabled, RotationLambdaARN, RotationRules).
//  2. If RotateImmediately is true (the default) and a rotation Lambda is
//     configured, execute the multi-step rotation:
//     a. createSecret — Lambda generates new secret material with AWSPENDING stage.
//     b. setSecret    — Lambda configures the target resource with the new value.
//     c. testSecret   — Lambda verifies the new secret works.
//     d. finishSecret — Service promotes AWSPENDING to AWSCURRENT.
//  3. If RotateImmediately is false, only the rotation configuration is updated;
//     the actual rotation is deferred to the next scheduled window.
//  4. If no rotation Lambda is set, fall back to copying the current version
//     (metadata-only rotation, preserving backward compatibility).
func (s *SecretsManagerService) rotateSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in RotateSecretInput) (*RotateSecretResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecret(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	if err := validateRotationLambdaARN(in.RotationLambdaARN); err != nil {
		return nil, err
	}
	if err := validateRotationRules(in.AutomaticallyAfterDays, in.ScheduleExpression, in.Duration); err != nil {
		return nil, err
	}

	// ClientRequestToken — idempotency token for rotation. When
	// provided, passed to executeRotation as the rotation cycle token.
	if err := validateClientRequestToken(in.ClientRequestToken); err != nil {
		return nil, err
	}
	if err := validateExternalSecretRotation(in.ExternalSecretRotationMetadata, in.ExternalSecretRotationRoleArn); err != nil {
		return nil, err
	}

	// RotateImmediately defaults to true per AWS spec.
	rotateImmediately := true
	if in.HasRotateImmediately {
		rotateImmediately = in.RotateImmediately
	}

	// RotateImmediately=false tests the rotation configuration by running
	// the testSecret step of the Lambda rotation function ("This test
	// creates an AWSPENDING version of the secret and then removes it").
	// Without a rotation Lambda there is nothing to probe and only the
	// configuration is updated.
	if !rotateImmediately && secret.RotationLambdaARN != "" {
		if s.bus == nil {
			return nil, errors.NewAWSError("InvalidRequestException",
				"Rotation Lambda is configured but the event bus is not available.", http.StatusBadRequest)
		}
		token := in.ClientRequestToken
		if token == "" {
			token = uuid.New().String()[:32]
		}
		if probeErr := s.invokeRotationLambda(ctx, secret.RotationLambdaARN, secret.ARN, token, rotationStepTest); probeErr != nil {
			if stderrors.Is(probeErr, errRotationLambdaNotFound) {
				return nil, errors.NewAWSError("ResourceNotFoundException",
					probeErr.Error(), http.StatusNotFound)
			}
			return nil, errors.NewAWSError("InvalidRequestException",
				probeErr.Error(), http.StatusBadRequest)
		}
	}

	secret.RotationEnabled = true
	if in.RotationLambdaARN != "" {
		secret.RotationLambdaARN = in.RotationLambdaARN
	}
	if in.AutomaticallyAfterDays > 0 || in.ScheduleExpression != "" {
		if secret.RotationRules == nil {
			secret.RotationRules = &secretsmanagerstore.RotationRules{}
		}
		if in.AutomaticallyAfterDays > 0 {
			secret.RotationRules.AutomaticallyAfterDays = in.AutomaticallyAfterDays
		}
		if in.ScheduleExpression != "" {
			secret.RotationRules.ScheduleExpression = in.ScheduleExpression
		}
	}
	if in.Duration != "" {
		if secret.RotationRules == nil {
			secret.RotationRules = &secretsmanagerstore.RotationRules{}
		}
		secret.RotationRules.Duration = in.Duration
	}
	if len(in.ExternalSecretRotationMetadata) > 0 {
		secret.ExternalSecretRotationMetadata = in.ExternalSecretRotationMetadata
	}
	if in.ExternalSecretRotationRoleArn != "" {
		secret.ExternalSecretRotationRoleArn = in.ExternalSecretRotationRoleArn
	}

	var versionId string
	rotatedViaLambda := false

	if rotateImmediately {
		if secret.RotationLambdaARN != "" && s.bus == nil {
			return nil, errors.NewAWSError("InvalidRequestException",
				"Rotation Lambda is configured but the event bus is not available.", http.StatusBadRequest)
		}
		if secret.RotationLambdaARN != "" && s.bus != nil {
			// executeRotation mutates LastRotatedDate/NextRotationDate on the
			// secret pointer before calling finishRotation. If finishRotation
			// fails, those dates must be restored so the error handler does
			// not persist a failed rotation as if it succeeded.
			origLastRotated := secret.LastRotatedDate
			origNextRotation := secret.NextRotationDate
			if rotErr := s.executeRotation(ctx, store, secret, in.ClientRequestToken); rotErr != nil {
				secret.LastRotatedDate = origLastRotated
				secret.NextRotationDate = origNextRotation
				secret.LastChangedDate = time.Now().UTC()
				_ = store.UpdateSecretMetadata(secret)
				// Distinguish Lambda-not-found (ResourceNotFoundException)
				// from contract violations (InvalidRequestException).
				if stderrors.Is(rotErr, errRotationLambdaNotFound) {
					return nil, errors.NewAWSError("ResourceNotFoundException",
						rotErr.Error(), http.StatusNotFound)
				}
				return nil, errors.NewAWSError("InvalidRequestException",
					rotErr.Error(), http.StatusBadRequest)
			}
			versionId = secret.CurrentVersion
			rotatedViaLambda = true
		} else {
			versionId, err = s.executeMetadataOnlyRotation(store, secret)
			if err != nil {
				return nil, err
			}
			secret.LastRotatedDate = storeClock()
		}
	}

	// Calculate NextRotationDate. executeRotation already set both
	// LastRotatedDate and NextRotationDate for the Lambda path, so only
	// recalculate for metadata-only or deferred rotations.
	if !rotatedViaLambda {
		if next := computeNextRotationDate(secret.RotationRules, secret.LastRotatedDate, storeClock()); !next.IsZero() {
			secret.NextRotationDate = next
		}
	}

	if err := store.UpdateSecretMetadata(secret); err != nil {
		return nil, mapStoreError(err)
	}

	return &RotateSecretResult{
		ARN:       secret.ARN,
		Name:      secret.Name,
		VersionID: versionId,
	}, nil
}

// executeMetadataOnlyRotation performs a rotation without invoking a Lambda
// function. It copies the current version to a new version, mirroring the
// pre-existing behaviour when no rotation Lambda is configured.
func (s *SecretsManagerService) executeMetadataOnlyRotation(store secretsmanagerstore.SecretStoreInterface, secret *secretsmanagerstore.Secret) (string, error) {
	if secret.CurrentVersion == "" {
		return "", nil
	}

	currentVer, verErr := store.GetSecretVersion(secret.Name, secret.CurrentVersion)
	if verErr != nil {
		return "", verErr
	}

	newVersion := secretsmanagerstore.NewSecretVersion("")
	newVersion.SecretName = secret.Name
	newVersion.SecretString = currentVer.SecretString
	newVersion.SecretBinary = make([]byte, len(currentVer.SecretBinary))
	copy(newVersion.SecretBinary, currentVer.SecretBinary)
	newVersion.VersionStages = []string{"AWSCURRENT"}

	if createErr := store.CreateVersionDirect(secret.Name, newVersion); createErr != nil {
		return "", createErr
	}

	oldPrevious, prevErr := store.GetSecretVersionByStage(secret.Name, "AWSPREVIOUS")
	if prevErr == nil && oldPrevious.VersionId != newVersion.VersionId {
		prevStages := []string{}
		for _, st := range oldPrevious.VersionStages {
			if st != "AWSPREVIOUS" {
				prevStages = append(prevStages, st)
			}
		}
		if err := store.UpdateSecretVersionStage(secret.Name, oldPrevious.VersionId, prevStages); err != nil {
			return "", mapStoreError(err)
		}
	}

	if err := store.UpdateSecretVersionStage(secret.Name, secret.CurrentVersion, []string{"AWSPREVIOUS"}); err != nil {
		return "", mapStoreError(err)
	}

	secret.VersionIDs = append(secret.VersionIDs, newVersion.VersionId)
	secret.CurrentVersion = newVersion.VersionId

	return newVersion.VersionId, nil
}

// cancelRotateSecretCore is the single entry point for disabling automatic
// rotation for a secret. If a version with the AWSPENDING stage exists, that
// stage label is removed before the rotation configuration is cleared.
func (s *SecretsManagerService) cancelRotateSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in CancelRotateSecretInput) (*CancelRotateSecretResult, error) {
	if err := validateSecretId(in.SecretId); err != nil {
		return nil, err
	}

	secret, err := resolveSecret(store, in.SecretId)
	if err != nil {
		return nil, err
	}

	pendingVer, err := store.GetSecretVersionByStage(secret.Name, "AWSPENDING")
	versionId := ""
	if err == nil {
		versionId = pendingVer.VersionId
		newStages := []string{}
		for _, st := range pendingVer.VersionStages {
			if st != "AWSPENDING" {
				newStages = append(newStages, st)
			}
		}
		if err := store.UpdateSecretVersionStage(secret.Name, pendingVer.VersionId, newStages); err != nil {
			return nil, mapStoreError(err)
		}
	}

	if err := store.CancelRotation(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return &CancelRotateSecretResult{
		ARN:       secret.ARN,
		Name:      secret.Name,
		VersionID: versionId,
	}, nil
}

// computeNextRotationDate derives the next scheduled rotation time from the
// rotation rules. A ScheduleExpression is evaluated with the shared AWS
// schedule-expression engine, anchored at the last rotation — or at the
// configuration moment when the secret has never rotated, so a fresh rate()
// schedule first fires one full period after configuration. cron() resolves
// to its first matching minute at or after now; rate() resolves to the
// current period boundary, and when a boundary has already elapsed past the
// anchor the owed boundary itself is returned as a past due time so the
// checker fires immediately instead of skipping the period. RotationRules
// Duration bounds the rotation window start; the checker fires from the
// scheduled time onwards, which lies inside the window. A zero time means
// "no schedule".
func computeNextRotationDate(rules *secretsmanagerstore.RotationRules, lastRotated, now time.Time) time.Time {
	if rules == nil {
		return time.Time{}
	}
	if rules.ScheduleExpression != "" {
		anchor := lastRotated
		if anchor.IsZero() {
			// A schedule first configured on an never-rotated secret
			// counts from the configuration moment, not the creation
			// date.
			anchor = now
		}
		if strings.HasPrefix(rules.ScheduleExpression, "rate(") {
			boundary, err := scheduleexpr.NextExecutionTime(rules.ScheduleExpression, now, anchor, nil)
			if err != nil {
				return time.Time{}
			}
			duration, ok := scheduleexpr.ParseRateDuration(rules.ScheduleExpression)
			if !ok {
				return time.Time{}
			}
			if boundary.After(anchor) {
				// A period boundary has elapsed since the anchor
				// rotation: that rotation is owed now.
				return boundary
			}
			return boundary.Add(duration)
		}
		next, err := scheduleexpr.NextExecutionTime(rules.ScheduleExpression, now, anchor, nil)
		if err != nil {
			return time.Time{}
		}
		return next
	}
	if rules.AutomaticallyAfterDays > 0 {
		base := lastRotated
		if base.IsZero() {
			base = now
		}
		return base.AddDate(0, 0, rules.AutomaticallyAfterDays)
	}
	return time.Time{}
}
