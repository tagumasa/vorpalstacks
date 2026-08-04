package rdsdata

import (
	"encoding/json"
	"errors"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
)

// AWS Data API parameter length bounds (see
// https://docs.aws.amazon.com/rdsdataservice/latest/APIReference/).
const (
	resourceArnMinLen   = 11
	resourceArnMaxLen   = 100
	secretArnMinLen     = 11
	secretArnMaxLen     = 100
	sqlMaxLen           = 65536
	transactionIDMaxLen = 192
	databaseMaxLen      = 64
	schemaMaxLen        = 64
)

// validateLength enforces AWS-spec string length constraints and returns a
// typed InvalidParameterException on violation. The min check is skipped
// when allowEmpty=true (AWS marks database / schema / transactionId as
// optional with min length 0).
func validateLength(field, value string, min, max int, allowEmpty bool) error {
	n := len(value)
	if allowEmpty && n == 0 {
		return nil
	}
	if n < min {
		return invalidParam(fmt.Sprintf("%s length %d is shorter than minimum %d", field, n, min))
	}
	if n > max {
		return invalidParam(fmt.Sprintf("%s length %d exceeds maximum %d", field, n, max))
	}
	return nil
}

// validateResourceArn enforces the AWS-spec resourceArn constraints.
func validateResourceArn(arn string) error {
	return validateLength("resourceArn", arn, resourceArnMinLen, resourceArnMaxLen, false)
}

// validateSecretArn enforces the AWS-spec secretArn constraints.
func validateSecretArn(arn string) error {
	return validateLength("secretArn", arn, secretArnMinLen, secretArnMaxLen, false)
}

// validateSQL enforces the AWS-spec sql constraints.
// AWS docs: "Minimum length of 0. Maximum length of 65536."
func validateSQL(sqlStr string) error {
	return validateLength("sql", sqlStr, 0, sqlMaxLen, true)
}

// validateTransactionID enforces the AWS-spec transactionId constraints.
// transactionId is optional for ExecuteStatement / BatchExecuteStatement but
// required for CommitTransaction / RollbackTransaction — callers that need
// 'required' semantics pass allowEmpty=false.
func validateTransactionID(id string, allowEmpty bool) error {
	return validateLength("transactionId", id, 0, transactionIDMaxLen, allowEmpty)
}

// validateDatabase enforces the AWS-spec database constraints.
func validateDatabase(db string) error {
	return validateLength("database", db, 0, databaseMaxLen, true)
}

// validateSchema enforces the AWS-spec schema constraints.
func validateSchema(schema string) error {
	return validateLength("schema", schema, 0, schemaMaxLen, true)
}

// validateCommon enforces the resourceArn + secretArn + database + schema
// constraints shared by every Data API operation.
func validateCommon(resourceArn, secretArn, database, schema string) error {
	if err := validateResourceArn(resourceArn); err != nil {
		return err
	}
	if err := validateSecretArn(secretArn); err != nil {
		return err
	}
	if err := validateDatabase(database); err != nil {
		return err
	}
	if err := validateSchema(schema); err != nil {
		return err
	}
	return nil
}

// mapSQLError preserves AWS-spec exception types emitted by executeSQL
// (StatementTimeoutException, DatabaseErrorException) and wraps any other
// engine-emitted error in DatabaseErrorException so callers see the
// correct exception name rather than a generic BadRequestException.
func mapSQLError(err error) error {
	if err == nil {
		return nil
	}
	var awsErr *awserrors.AWSError
	if errors.As(err, &awsErr) {
		return awsErr
	}
	return databaseError(err.Error())
}

// parseRequest deserialises the JSON request body into the provided struct.
func parseRequest(req *request.ParsedRequest, v interface{}) error {
	if req == nil || len(req.Body) == 0 {
		return invalidParam("request body is empty")
	}
	if err := json.Unmarshal(req.Body, v); err != nil {
		return badRequest(fmt.Sprintf("failed to parse request: %v", err))
	}
	return nil
}
