package secretsmanager

import (
	"context"
	stderrors "errors"
	"net/http"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

func nestedInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return 0
}

func nestedString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// RestoreSecret restores a previously deleted secret. Secrets Manager keeps
// deleted secrets recoverable for a minimum of 30 days.
func (s *SecretsManagerService) RestoreSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecretForMetadata(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.RestoreSecret(secret.Name); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

// RotateSecret configures rotation for a secret and, if a RotationLambdaARN
// is available and RotateImmediately is true, performs the full AWS rotation
// protocol:
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
func (s *SecretsManagerService) RotateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	rotationLambdaARN := request.GetStringParam(req.Parameters, "RotationLambdaARN")
	rotationRulesRaw, _ := req.Parameters["RotationRules"].(map[string]interface{})
	if rotationRulesRaw == nil {
		rotationRulesRaw = make(map[string]interface{})
	}
	automaticallyAfterDays := nestedInt(rotationRulesRaw, "AutomaticallyAfterDays")
	scheduleExpression := nestedString(rotationRulesRaw, "ScheduleExpression")
	duration := nestedString(rotationRulesRaw, "Duration")
	if automaticallyAfterDays > 0 {
		if err := validateAutomaticallyAfterDays(automaticallyAfterDays); err != nil {
			return nil, err
		}
	}

	// ClientRequestToken — idempotency token for rotation. When
	// provided, passed to executeRotation as the rotation cycle token.
	// Smithy ClientRequestTokenType length 32-64.
	clientRequestToken := request.GetStringParam(req.Parameters, "ClientRequestToken")
	if clientRequestToken != "" {
		if len(clientRequestToken) < 32 || len(clientRequestToken) > 64 {
			return nil, errors.NewAWSError("InvalidParameterException",
				"ClientRequestToken must be 32 to 64 characters long.", http.StatusBadRequest)
		}
	}

	// RotateImmediately defaults to true per AWS spec.
	rotateImmediately := true
	if request.HasParam(req.Parameters, "RotateImmediately") {
		rotateImmediately = request.GetBoolParam(req.Parameters, "RotateImmediately")
	}

	secret.RotationEnabled = true
	if rotationLambdaARN != "" {
		secret.RotationLambdaARN = rotationLambdaARN
	}
	if automaticallyAfterDays > 0 || scheduleExpression != "" {
		if secret.RotationRules == nil {
			secret.RotationRules = &secretsmanagerstore.RotationRules{}
		}
		if automaticallyAfterDays > 0 {
			secret.RotationRules.AutomaticallyAfterDays = automaticallyAfterDays
		}
		if scheduleExpression != "" {
			secret.RotationRules.ScheduleExpression = scheduleExpression
		}
	}
	if duration != "" {
		if secret.RotationRules == nil {
			secret.RotationRules = &secretsmanagerstore.RotationRules{}
		}
		secret.RotationRules.Duration = duration
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
			if rotErr := s.executeRotation(ctx, store, secret, clientRequestToken); rotErr != nil {
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
		if automaticallyAfterDays > 0 {
			base := secret.LastRotatedDate
			if base.IsZero() {
				base = storeClock()
			}
			secret.NextRotationDate = base.AddDate(0, 0, automaticallyAfterDays)
		}
	}

	if err := store.UpdateSecretMetadata(secret); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":       secret.ARN,
		"Name":      secret.Name,
		"VersionId": versionId,
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

// CancelRotateSecret disables automatic rotation for a secret. If a version
// with the AWSPENDING stage exists, that stage label is removed before the
// rotation configuration is cleared.
func (s *SecretsManagerService) CancelRotateSecret(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, errors.ErrMissingParameter
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
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

	result := map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}
	if versionId != "" {
		result["VersionId"] = versionId
	}
	return result, nil
}
