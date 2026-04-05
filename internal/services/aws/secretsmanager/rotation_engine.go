package secretsmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

const (
	rotationCheckInterval = 60 * time.Second
	rotationStepCreate    = "createSecret"
	rotationStepSet       = "setSecret"
	rotationStepTest      = "testSecret"
)

// rotationChecker periodically scans secrets for those whose NextRotationDate
// has passed and triggers automatic rotation via the rotation engine.
type rotationChecker struct {
	svc    *SecretsManagerService
	logger logs.Logger
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// newRotationChecker creates a new rotation checker bound to the given service.
func newRotationChecker(svc *SecretsManagerService, logger logs.Logger) *rotationChecker {
	return &rotationChecker{
		svc:    svc,
		logger: logger,
		stopCh: make(chan struct{}),
	}
}

// start launches the background rotation check loop.
func (rc *rotationChecker) start(ctx context.Context) {
	rc.wg.Add(1)
	go rc.loop(ctx)
}

// stop signals the rotation checker to terminate and waits for it to finish.
func (rc *rotationChecker) stop() {
	close(rc.stopCh)
	rc.wg.Wait()
}

// loop ticks at rotationCheckInterval and checks for secrets due for rotation.
func (rc *rotationChecker) loop(ctx context.Context) {
	defer rc.wg.Done()
	ticker := time.NewTicker(rotationCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rc.checkDueRotations(ctx)
		case <-rc.stopCh:
			return
		}
	}
}

// checkDueRotations lists all secrets across all regions and triggers
// rotation for any whose NextRotationDate has passed.
func (rc *rotationChecker) checkDueRotations(ctx context.Context) {
	if rc.svc.storageManager == nil {
		return
	}

	regions := rc.svc.storageManager.GetActiveRegions()
	for _, region := range regions {
		rc.checkDueRotationsForRegion(ctx, region)
	}
}

func (rc *rotationChecker) checkDueRotationsForRegion(ctx context.Context, region string) {
	storage, err := rc.svc.storageManager.GetStorage(region)
	if err != nil {
		rc.logError("failed to get storage for rotation check", logs.String("region", region), logs.Err(err))
		return
	}
	store := secretsmanagerstore.NewSecretStore(storage, rc.svc.accountID, region)

	opts := common.ListOptions{MaxItems: 100}
	result, err := common.List[secretsmanagerstore.Secret](store.GetBaseStore(), opts, nil)
	if err != nil {
		rc.logError("failed to list secrets for rotation check", logs.Err(err))
		return
	}

	now := time.Now().UTC()
	for _, sec := range result.Items {
		if !sec.RotationEnabled || sec.RotationLambdaARN == "" {
			continue
		}
		if sec.DeletedDate != nil {
			continue
		}
		if !sec.NextRotationDate.IsZero() && sec.NextRotationDate.Before(now) {
			if err := rc.svc.executeRotation(ctx, store, sec); err != nil {
				rc.logError("automatic rotation failed",
					logs.String("secret", sec.Name),
					logs.Err(err))
			}
		}
	}
}

// logError logs a structured error message if a logger is available.
func (rc *rotationChecker) logError(msg string, fields ...logs.Field) {
	if rc.logger != nil {
		rc.logger.Error(msg, fields...)
	}
}

// rotationPayload represents the JSON payload sent to a rotation Lambda
// function. This matches the AWS Secrets Manager rotation request format.
type rotationPayload struct {
	SecretId           string `json:"SecretId"`
	ClientRequestToken string `json:"ClientRequestToken"`
	Step               string `json:"Step"`
}

// rotationLambdaResponse is the expected JSON response from a rotation Lambda.
// The function may return an empty object or additional metadata.
type rotationLambdaResponse struct {
	SecretId string `json:"SecretId"`
}

// executeRotation performs the full multi-step rotation sequence for a secret.
//
// The AWS rotation protocol requires four sequential steps:
//  1. createSecret — Lambda generates new secret material and stores it
//     via PutSecretValue with AWSPENDING stage.
//  2. setSecret    — Lambda configures the target resource (database, etc.)
//     with the new secret value from AWSPENDING.
//  3. testSecret   — Lambda verifies the new secret works by testing
//     connectivity with the AWSPENDING value.
//  4. finishSecret — (implicit) the service promotes AWSPENDING to AWSCURRENT
//     and demotes AWSCURRENT to AWSPREVIOUS.
//
// Steps 1-3 invoke the rotation Lambda; step 4 is handled locally by the
// service after all Lambda invocations succeed.
func (s *SecretsManagerService) executeRotation(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, secret *secretsmanagerstore.Secret) error {
	if s.lambdaInvoker == nil {
		return fmt.Errorf("lambda invoker not configured for secret rotation")
	}

	clientToken := uuid.New().String()[:32]

	steps := []string{rotationStepCreate, rotationStepSet, rotationStepTest}

	for _, step := range steps {
		payload := rotationPayload{
			SecretId:           secret.ARN,
			ClientRequestToken: clientToken,
			Step:               step,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal rotation payload for step %s: %w", step, err)
		}

		statusCode, respBytes, invokeErr := s.lambdaInvoker.InvokeForGateway(ctx, secret.RotationLambdaARN, payloadBytes)
		if invokeErr != nil {
			return fmt.Errorf("rotation Lambda invocation failed at step %s: %w", step, invokeErr)
		}

		if statusCode != 200 {
			return fmt.Errorf("rotation Lambda returned status %d at step %s: %s", statusCode, step, string(respBytes))
		}

		var lambdaResp rotationLambdaResponse
		if len(respBytes) > 0 {
			if jsonErr := json.Unmarshal(respBytes, &lambdaResp); jsonErr != nil {
				return fmt.Errorf("failed to parse rotation Lambda response at step %s: %w", step, jsonErr)
			}
		}
	}

	if err := s.finishRotation(store, secret, clientToken); err != nil {
		return fmt.Errorf("failed to finish rotation: %w", err)
	}

	return nil
}

// finishRotation promotes the AWSPENDING version to AWSCURRENT and demotes
// the current AWSCURRENT to AWSPREVIOUS. This is the local counterpart of
// the AWS finishSecret step.
func (s *SecretsManagerService) finishRotation(store secretsmanagerstore.SecretStoreInterface, secret *secretsmanagerstore.Secret, pendingVersionID string) error {
	_, err := store.GetSecretVersion(secret.Name, pendingVersionID)
	if err != nil {
		return fmt.Errorf("AWSPENDING version %s not found: %w", pendingVersionID, err)
	}

	oldCurrentID := secret.CurrentVersion

	oldPrevious, prevErr := store.GetSecretVersionByStage(secret.Name, "AWSPREVIOUS")
	if prevErr == nil && oldPrevious.VersionId != pendingVersionID {
		cleanedStages := []string{}
		for _, st := range oldPrevious.VersionStages {
			if st != "AWSPREVIOUS" {
				cleanedStages = append(cleanedStages, st)
			}
		}
		_ = store.UpdateSecretVersionStage(secret.Name, oldPrevious.VersionId, cleanedStages)
	}

	if oldCurrentID != "" && oldCurrentID != pendingVersionID {
		_ = store.UpdateSecretVersionStage(secret.Name, oldCurrentID, []string{"AWSPREVIOUS"})
	}

	_ = store.UpdateSecretVersionStage(secret.Name, pendingVersionID, []string{"AWSCURRENT"})

	secret.CurrentVersion = pendingVersionID
	secret.LastRotatedDate = storeClock()

	if secret.RotationRules != nil && secret.RotationRules.AutomaticallyAfterDays > 0 {
		secret.NextRotationDate = secret.LastRotatedDate.AddDate(0, 0, secret.RotationRules.AutomaticallyAfterDays)
	}

	if err := store.UpdateSecretMetadata(secret); err != nil {
		return fmt.Errorf("failed to update secret metadata after rotation: %w", err)
	}

	s.logRotation(secret.Name, pendingVersionID)
	return nil
}

// logRotation emits a structured info log for a successful rotation.
func (s *SecretsManagerService) logRotation(secretName, versionID string) {
	if s.logger != nil {
		s.logger.Info("secret rotation completed",
			logs.String("secret", secretName),
			logs.String("version", versionID))
	}
}
