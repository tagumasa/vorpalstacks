package secretsmanager

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	// maxTagsPerSecret is the maximum number of tags allowed on a single secret.
	maxTagsPerSecret = 50
	// maxTagKeyLength is the maximum length of a tag key in characters.
	maxTagKeyLength = 128
	// maxTagValueLength is the maximum length of a tag value in characters.
	maxTagValueLength = 256
)

// validateSecretTags validates tag count, key length, and value length against
// AWS Secrets Manager quotas.  Tag count overflow uses
// InvalidParameterException — it is documented for both CreateSecret and
// TagResource, whereas LimitExceededException is only documented for
// CreateSecret.
func validateSecretTags(tags []types.Tag) error {
	if len(tags) > maxTagsPerSecret {
		return awserrors.NewAWSError("InvalidParameterException",
			fmt.Sprintf("You can't have more than %d tags on a secret.", maxTagsPerSecret), http.StatusBadRequest)
	}
	for _, t := range tags {
		if t.Key == "" || len(t.Key) > maxTagKeyLength {
			return awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("Tag key length must be between 1 and %d characters.", maxTagKeyLength), http.StatusBadRequest)
		}
		if len(t.Value) > maxTagValueLength {
			return awserrors.NewAWSError("InvalidParameterException",
				fmt.Sprintf("Tag value length must be between 0 and %d characters.", maxTagValueLength), http.StatusBadRequest)
		}
	}
	return nil
}

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

// arnPrefix is the standard AWS ARN prefix used to detect ARN-based secret IDs.
const arnPrefix = "arn:"

// storeErrorMappings maps store-level sentinel errors to AWS API errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: secretsmanagerstore.ErrSecretNotFound, AWS: ErrSecretNotFound},
	{Store: secretsmanagerstore.ErrSecretAlreadyExists, AWS: ErrSecretAlreadyExists},
	{Store: secretsmanagerstore.ErrInvalidSecretName, AWS: ErrInvalidSecretName},
	{Store: secretsmanagerstore.ErrInvalidVersionId, AWS: ErrInvalidVersionId},
}

func mapStoreError(err error) error {
	return awserrors.MapStoreError(err, storeErrorMappings)
}

func (s *SecretsManagerService) resolveSecret(reqCtx *request.RequestContext, secretId string) (*secretsmanagerstore.Secret, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(secretId, arnPrefix) {
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
	if strings.HasPrefix(secretId, arnPrefix) {
		name, err := store.LookupNameByARN(secretId)
		if err != nil {
			return nil, mapStoreError(err)
		}
		return store.GetSecretForMetadata(name)
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
			}
		}
		return result
	}

	for _, versionId := range secret.VersionIDs {
		version, err := store.GetSecretVersion(secret.Name, versionId)
		if err == nil && len(version.VersionStages) > 0 {
			result[versionId] = version.VersionStages
		} else if versionId == secret.CurrentVersion {
			result[versionId] = []string{"AWSCURRENT"}
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
	if secret.RotationEnabled || secret.RotationLambdaARN != "" {
		m["RotationEnabled"] = secret.RotationEnabled
	}
	if secret.RotationLambdaARN != "" {
		m["RotationLambdaARN"] = secret.RotationLambdaARN
	}
	if secret.RotationRules != nil {
		rules := map[string]interface{}{}
		if secret.RotationRules.AutomaticallyAfterDays > 0 {
			rules["AutomaticallyAfterDays"] = secret.RotationRules.AutomaticallyAfterDays
		}
		if secret.RotationRules.ScheduleExpression != "" {
			rules["ScheduleExpression"] = secret.RotationRules.ScheduleExpression
		}
		if secret.RotationRules.Duration != "" {
			rules["Duration"] = secret.RotationRules.Duration
		}
		if len(rules) > 0 {
			m["RotationRules"] = rules
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
		return nil, awserrors.NewValidationException(fmt.Sprintf("invalid SecretBinary encoding: %v", err))
	}
	return decoded, nil
}

func storeClock() time.Time {
	return time.Now().UTC()
}
