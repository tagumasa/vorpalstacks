package secretsmanager

import (
	"net/http"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
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

// arnPrefix is the standard AWS ARN prefix used to detect ARN-based secret IDs.
const arnPrefix = "arn:"

// storeErrorMappings maps store-level sentinel errors to AWS API errors.
var storeErrorMappings = []awserrors.StoreErrorMapping{
	{Store: secretsmanagerstore.ErrSecretNotFound, AWS: ErrSecretNotFound},
	{Store: secretsmanagerstore.ErrSecretAlreadyExists, AWS: ErrSecretAlreadyExists},
	{Store: secretsmanagerstore.ErrInvalidSecretName, AWS: ErrInvalidSecretName},
	{Store: secretsmanagerstore.ErrInvalidVersionId, AWS: ErrInvalidVersionId},
}

// parseOffsetNextToken parses an integer-offset pagination token for the
// list operations (ListSecrets, ListSecretVersionIds, BatchGetSecretValue).
// Tokens are opaque to clients but internally a plain decimal offset, so
// anything non-numeric, negative, or beyond the result set is invalid —
// all three operations document InvalidNextTokenException (HTTP 400) for
// an invalid NextToken value.
func parseOffsetNextToken(token string, total int) (int, error) {
	if token == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(token)
	if err != nil || n < 0 || n >= total {
		return 0, awserrors.NewAWSError("InvalidNextTokenException",
			"Your request has an invalid next token.", http.StatusBadRequest)
	}
	return n, nil
}

func mapStoreError(err error) error {
	return awserrors.MapStoreError(err, storeErrorMappings)
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

func storeClock() time.Time {
	return time.Now().UTC()
}
