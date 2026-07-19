package rdsdata

import (
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// AWS RDS Data API exception names and HTTP status codes are taken from
// https://docs.aws.amazon.com/rdsdataservice/latest/APIReference/

func badRequest(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("BadRequestException", msg, http.StatusBadRequest)
}

func invalidParam(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidParameterException", msg, http.StatusBadRequest)
}

func accessDenied(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("AccessDeniedException", msg, http.StatusForbidden)
}

func internalError(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InternalServerErrorException", msg, http.StatusInternalServerError)
}

func transactionNotFound(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("TransactionNotFoundException", msg, http.StatusNotFound)
}

// databaseError corresponds to DatabaseErrorException (HTTP 400): "There was
// an error in processing the SQL statement." Use this for any failure emitted
// by the underlying SQL engine (parse error, runtime error, constraint
// violation, etc.) as opposed to BadRequestException which is reserved for
// request-shape problems.
func databaseError(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("DatabaseErrorException", msg, http.StatusBadRequest)
}

// notFoundError corresponds to NotFoundException (HTTP 404): one of
// resourceArn / secretArn / transactionId supplied does not exist.
func notFoundError(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("NotFoundException", msg, http.StatusNotFound)
}

// databaseNotFound corresponds to DatabaseNotFoundException (HTTP 404): "The
// DB cluster doesn't have a DB instance." Returned when a cluster ARN
// resolves to zero writer instances.
func databaseNotFound(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("DatabaseNotFoundException", msg, http.StatusNotFound)
}

// databaseUnavailable corresponds to DatabaseUnavailableException (HTTP 504):
// "The writer instance in the DB cluster isn't available."
func databaseUnavailable(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("DatabaseUnavailableException", msg, http.StatusGatewayTimeout)
}

// invalidSecret corresponds to InvalidSecretException (HTTP 400): the secret
// referenced by secretArn exists but is not a valid RDS database secret
// (missing username/password, wrong format, etc.). Exported as InvalidSecret
// so external SecretResolver implementations can return the correctly-typed
// AWS error without re-implementing the wrapper.
func invalidSecret(msg string) *awserrors.AWSError {
	return InvalidSecret(msg)
}

// InvalidSecret constructs an InvalidSecretException (HTTP 400).
func InvalidSecret(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidSecretException", msg, http.StatusBadRequest)
}

// secretsError corresponds to SecretsErrorException (HTTP 400): a problem
// retrieving the secret itself (not found, decrypt failure, timeout).
func secretsError(msg string) *awserrors.AWSError {
	return SecretsError(msg)
}

// SecretsError constructs a SecretsErrorException (HTTP 400).
func SecretsError(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("SecretsErrorException", msg, http.StatusBadRequest)
}

// statementTimeout corresponds to StatementTimeoutException (HTTP 400). The
// optional dbConnectionId field is not surfaced in the AWSError helper; if
// needed it can be threaded through later.
func statementTimeout(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("StatementTimeoutException", msg, http.StatusBadRequest)
}

// serviceUnavailable corresponds to ServiceUnavailableError (HTTP 503).
func serviceUnavailable(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("ServiceUnavailableError", msg, http.StatusServiceUnavailable)
}

// forbiddenError corresponds to ForbiddenException (HTTP 403).
func forbiddenError(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("ForbiddenException", msg, http.StatusForbidden)
}
