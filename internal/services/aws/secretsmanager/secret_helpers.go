package secretsmanager

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

var (
	// ErrSecretNotFound indicates that the specified secret does not exist.
	ErrSecretNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Secrets Manager can't find the specified secret.", http.StatusNotFound)
	// ErrSecretAlreadyExists indicates that a resource with the specified ID already exists.
	ErrSecretAlreadyExists = awserrors.NewAWSError("ResourceExistsException", "A resource with the ID you requested already exists.", http.StatusBadRequest)
	// ErrInvalidSecretName indicates that the parameter name is not valid.
	ErrInvalidSecretName = awserrors.NewAWSError("InvalidParameterException", "The parameter name is not valid.", http.StatusBadRequest)
	// ErrInvalidVersionId indicates that the version ID is not valid.
	ErrInvalidVersionId = awserrors.NewAWSError("InvalidParameterException", "The version ID is not valid.", http.StatusBadRequest)
	// ErrInvalidRequest indicates that the request is not valid.
	ErrInvalidRequest = awserrors.NewAWSError("InvalidParameterException", "The request is not valid.", http.StatusBadRequest)
)

func mapStoreError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, secretsmanagerstore.ErrSecretNotFound) {
		return ErrSecretNotFound
	}
	if errors.Is(err, secretsmanagerstore.ErrSecretAlreadyExists) {
		return ErrSecretAlreadyExists
	}
	if errors.Is(err, secretsmanagerstore.ErrInvalidSecretName) {
		return ErrInvalidSecretName
	}
	if errors.Is(err, secretsmanagerstore.ErrInvalidVersionId) {
		return ErrInvalidVersionId
	}
	return err
}

func (s *SecretsManagerService) resolveSecret(reqCtx *request.RequestContext, secretId string) (*secretsmanagerstore.Secret, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if len(secretId) >= 3 && secretId[:3] == "arn" {
		secret, err := store.GetSecretByARN(secretId)
		return secret, mapStoreError(err)
	}
	secret, err := store.GetSecret(secretId)
	return secret, mapStoreError(err)
}

func (s *SecretsManagerService) resolveSecretForMetadata(reqCtx *request.RequestContext, secretId string) (*secretsmanagerstore.Secret, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if len(secretId) >= 3 && secretId[:3] == "arn" {
		opts := common.ListOptions{MaxItems: 1000}
		result, err := common.List[secretsmanagerstore.Secret](store.GetBaseStore(), opts, func(sec *secretsmanagerstore.Secret) bool {
			return sec.ARN == secretId
		})
		if err != nil {
			return nil, err
		}
		if len(result.Items) == 0 {
			return nil, ErrSecretNotFound
		}
		return result.Items[0], nil
	}
	secret, err := store.GetSecretForMetadata(secretId)
	return secret, mapStoreError(err)
}

func (s *SecretsManagerService) buildSecretVersionsToStages(reqCtx *request.RequestContext, secret *secretsmanagerstore.Secret) map[string][]string {
	result := make(map[string][]string)
	store, err := s.store(reqCtx)
	if err != nil {
		for _, versionId := range secret.VersionIDs {
			if versionId == secret.CurrentVersion {
				result[versionId] = []string{"AWSCURRENT"}
			} else {
				result[versionId] = []string{}
			}
		}
		return result
	}

	for _, versionId := range secret.VersionIDs {
		version, err := store.GetSecretVersion(secret.Name, versionId)
		if err == nil {
			result[versionId] = version.VersionStages
		} else {
			if versionId == secret.CurrentVersion {
				result[versionId] = []string{"AWSCURRENT"}
			} else {
				result[versionId] = []string{}
			}
		}
	}
	return result
}

func (s *SecretsManagerService) buildTagsList(secret *secretsmanagerstore.Secret) []interface{} {
	tags := tagutil.MapToResponse(secret.Tags)
	result := make([]interface{}, len(tags))
	for i, t := range tags {
		result[i] = t
	}
	return result
}

func (s *SecretsManagerService) addRotationFields(m map[string]interface{}, secret *secretsmanagerstore.Secret) {
	if secret.RotationEnabled {
		m["RotationEnabled"] = true
	}
	if secret.RotationLambdaARN != "" {
		m["RotationLambdaARN"] = secret.RotationLambdaARN
	}
	if secret.RotationRules != nil {
		m["RotationRules"] = map[string]interface{}{
			"AutomaticallyAfterDays": secret.RotationRules.AutomaticallyAfterDays,
		}
	}
	if !secret.LastRotatedDate.IsZero() {
		m["LastRotatedDate"] = secret.LastRotatedDate.Unix()
	}
	if secret.DeletedDate != nil {
		m["DeletedDate"] = secret.DeletedDate.Unix()
	}
}

func decodeSecretBinary(secretBinaryStr string) ([]byte, error) {
	if secretBinaryStr == "" {
		return nil, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(secretBinaryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SecretBinary encoding: %w", err)
	}
	return decoded, nil
}

func storeClock() time.Time {
	return time.Now().UTC()
}
