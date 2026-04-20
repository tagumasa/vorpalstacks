package neptunedata

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// badRequest returns a BadRequestException (400).
func badRequest(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("BadRequestException", msg, http.StatusBadRequest)
}

// malformedQuery returns a MalformedQueryException (400) for syntactically
// invalid query strings.
func malformedQuery(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("MalformedQueryException", msg, http.StatusBadRequest)
}

// unsupported returns an UnsupportedOperationException (400) for operations
// that are not supported by vorpalstacks (e.g. SPARQL, ML operations).
func unsupported(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("UnsupportedOperationException", msg, http.StatusBadRequest)
}

// invalidParameter returns an InvalidParameterException (400).
func invalidParameter(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InvalidParameterException", msg, http.StatusBadRequest)
}

// missingParameter returns a MissingParameterException (400).
func missingParameter(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("MissingParameterException", msg, http.StatusBadRequest)
}

// bulkLoadNotFound returns a BulkLoadIdNotFoundException (404).
func bulkLoadNotFound(id string) *awserrors.AWSError {
	return awserrors.NewAWSError("BulkLoadIdNotFoundException", fmt.Sprintf("Load ID not found: %s", id), http.StatusNotFound)
}

// preconditionFailed returns a PreconditionsFailedException (412), used when
// a FastReset token is invalid or expired.
func preconditionFailed(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("PreconditionsFailedException", msg, http.StatusPreconditionFailed)
}

// internalFailure returns an InternalFailureException (500).
func internalFailure(msg string) *awserrors.AWSError {
	return awserrors.NewAWSError("InternalFailureException", msg, http.StatusInternalServerError)
}
