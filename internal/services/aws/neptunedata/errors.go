package neptunedata

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
)

// NewNeptuneDataError creates a NeptuneDataError with the given error code,
// message, and HTTP status code.
func NewNeptuneDataError(code, message string, httpStatus int) *awserrors.AWSError {
	return awserrors.NewAWSError(code, message, httpStatus)
}

// badRequest returns a BadRequestException (400).
func badRequest(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("BadRequestException", msg, http.StatusBadRequest)
}

// malformedQuery returns a MalformedQueryException (400) for syntactically
// invalid query strings.
func malformedQuery(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("MalformedQueryException", msg, http.StatusBadRequest)
}

// unsupported returns an UnsupportedOperationException (400) for operations
// that are not supported by vorpalstacks (e.g. SPARQL, ML operations).
func unsupported(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("UnsupportedOperationException", msg, http.StatusBadRequest)
}

// invalidParameter returns an InvalidParameterException (400).
func invalidParameter(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("InvalidParameterException", msg, http.StatusBadRequest)
}

// missingParameter returns a MissingParameterException (400).
func missingParameter(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("MissingParameterException", msg, http.StatusBadRequest)
}

// bulkLoadNotFound returns a BulkLoadIdNotFoundException (404).
func bulkLoadNotFound(id string) *awserrors.AWSError {
	return NewNeptuneDataError("BulkLoadIdNotFoundException", fmt.Sprintf("Load ID not found: %s", id), http.StatusNotFound)
}

// accessDenied returns an AccessDeniedException (403).
func accessDenied(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("AccessDeniedException", msg, http.StatusForbidden)
}

// throttling returns a ThrottlingException (429).
func throttling(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ThrottlingException", msg, http.StatusTooManyRequests)
}

// preconditionFailed returns a PreconditionsFailedException (412), used when
// a FastReset token is invalid or expired.
func preconditionFailed(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("PreconditionsFailedException", msg, http.StatusPreconditionFailed)
}

// methodNotAllowed returns a MethodNotAllowedException (405).
func methodNotAllowed(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("MethodNotAllowedException", msg, http.StatusMethodNotAllowed)
}

// queryLimitExceeded returns a QueryLimitExceededException (429).
func queryLimitExceeded(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("QueryLimitExceededException", msg, http.StatusTooManyRequests)
}

// internalFailure returns an InternalFailureException (500).
func internalFailure(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("InternalFailureException", msg, http.StatusInternalServerError)
}

// timeLimitExceeded returns a TimeLimitExceededException (503).
func timeLimitExceeded(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("TimeLimitExceededException", msg, http.StatusServiceUnavailable)
}

// constraintViolation returns a ConstraintViolationException (409).
func constraintViolation(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ConstraintViolationException", msg, http.StatusConflict)
}

// invalidNumericData returns an InvalidNumericDataException (400).
func invalidNumericData(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("InvalidNumericDataException", msg, http.StatusBadRequest)
}

// parsingError returns a ParsingException (400).
func parsingError(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ParsingException", msg, http.StatusBadRequest)
}

// readOnlyViolation returns a ReadOnlyViolationException (400).
func readOnlyViolation(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ReadOnlyViolationException", msg, http.StatusBadRequest)
}

// cancelledByUser returns a CancelledByUserException (400).
func cancelledByUser(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("CancelledByUserException", msg, http.StatusBadRequest)
}

// concurrentModification returns a ConcurrentModificationException (409).
func concurrentModification(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ConcurrentModificationException", msg, http.StatusConflict)
}

// clientTimeout returns a ClientTimeoutException (408).
func clientTimeout(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ClientTimeoutException", msg, http.StatusRequestTimeout)
}

// statisticsNotAvailable returns a StatisticsNotAvailableException (404).
func statisticsNotAvailable(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("StatisticsNotAvailableException", msg, http.StatusNotFound)
}

// s3Exception returns an S3Exception (400).
func s3Exception(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("S3Exception", msg, http.StatusBadRequest)
}

// illegalArgument returns an IllegalArgumentException (400).
func illegalArgument(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("IllegalArgumentException", msg, http.StatusBadRequest)
}

// invalidArgument returns an InvalidArgumentException (400).
func invalidArgument(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("InvalidArgumentException", msg, http.StatusBadRequest)
}

// memoryLimitExceeded returns a MemoryLimitExceededException (500).
func memoryLimitExceeded(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("MemoryLimitExceededException", msg, http.StatusInternalServerError)
}

// queryTooLarge returns a QueryTooLargeException (400).
func queryTooLarge(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("QueryTooLargeException", msg, http.StatusBadRequest)
}

// tooManyRequests returns a TooManyRequestsException (429).
func tooManyRequests(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("TooManyRequestsException", msg, http.StatusTooManyRequests)
}

// serverShutdown returns a ServerShutdownException (500).
func serverShutdown(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ServerShutdownException", msg, http.StatusInternalServerError)
}

// loadUrlAccessDenied returns a LoadUrlAccessDeniedException (400).
func loadUrlAccessDenied(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("LoadUrlAccessDeniedException", msg, http.StatusBadRequest)
}

// expiredStream returns an ExpiredStreamException (400).
func expiredStream(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("ExpiredStreamException", msg, http.StatusBadRequest)
}

// streamRecordsNotFound returns a StreamRecordsNotFoundException (404).
func streamRecordsNotFound(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("StreamRecordsNotFoundException", msg, http.StatusNotFound)
}

// mlResourceNotFound returns an MLResourceNotFoundException (404).
func mlResourceNotFound(msg string) *awserrors.AWSError {
	return NewNeptuneDataError("MLResourceNotFoundException", msg, http.StatusNotFound)
}
