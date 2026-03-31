package neptunedata

import (
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
)

// NeptuneDataError wraps an AWSError with Neptune Data API-specific unwrapping.
// This ensures the dispatcher can extract the error code and HTTP status for
// the X-Amzn-ErrorType response header.
type NeptuneDataError struct {
	*awserrors.AWSError
}

// Unwrap returns the underlying AWSError for error chain walking by the
// dispatcher's extractAWSError function.
func (e *NeptuneDataError) Unwrap() error {
	return e.AWSError
}

// NewNeptuneDataError creates a NeptuneDataError with the given error code,
// message, and HTTP status code.
func NewNeptuneDataError(code, message string, httpStatus int) *NeptuneDataError {
	return &NeptuneDataError{
		AWSError: awserrors.NewAWSError(code, message, httpStatus),
	}
}

// badRequest returns a BadRequestException (400).
func badRequest(msg string) *NeptuneDataError {
	return NewNeptuneDataError("BadRequestException", msg, http.StatusBadRequest)
}

// malformedQuery returns a MalformedQueryException (400) for syntactically
// invalid query strings.
func malformedQuery(msg string) *NeptuneDataError {
	return NewNeptuneDataError("MalformedQueryException", msg, http.StatusBadRequest)
}

// unsupported returns an UnsupportedOperationException (400) for operations
// that are not supported by vorpalstacks (e.g. SPARQL, ML operations).
func unsupported(msg string) *NeptuneDataError {
	return NewNeptuneDataError("UnsupportedOperationException", msg, http.StatusBadRequest)
}

// invalidParameter returns an InvalidParameterException (400).
func invalidParameter(msg string) *NeptuneDataError {
	return NewNeptuneDataError("InvalidParameterException", msg, http.StatusBadRequest)
}

// missingParameter returns a MissingParameterException (400).
func missingParameter(msg string) *NeptuneDataError {
	return NewNeptuneDataError("MissingParameterException", msg, http.StatusBadRequest)
}

// bulkLoadNotFound returns a BulkLoadIdNotFoundException (404).
func bulkLoadNotFound(id string) *NeptuneDataError {
	return NewNeptuneDataError("BulkLoadIdNotFoundException", fmt.Sprintf("Load ID not found: %s", id), http.StatusNotFound)
}

// accessDenied returns an AccessDeniedException (403).
func accessDenied(msg string) *NeptuneDataError {
	return NewNeptuneDataError("AccessDeniedException", msg, http.StatusForbidden)
}

// throttling returns a ThrottlingException (429).
func throttling(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ThrottlingException", msg, http.StatusTooManyRequests)
}

// preconditionFailed returns a PreconditionsFailedException (412), used when
// a FastReset token is invalid or expired.
func preconditionFailed(msg string) *NeptuneDataError {
	return NewNeptuneDataError("PreconditionsFailedException", msg, http.StatusPreconditionFailed)
}

// methodNotAllowed returns a MethodNotAllowedException (405).
func methodNotAllowed(msg string) *NeptuneDataError {
	return NewNeptuneDataError("MethodNotAllowedException", msg, http.StatusMethodNotAllowed)
}

// queryLimitExceeded returns a QueryLimitExceededException (429).
func queryLimitExceeded(msg string) *NeptuneDataError {
	return NewNeptuneDataError("QueryLimitExceededException", msg, http.StatusTooManyRequests)
}

// internalFailure returns an InternalFailureException (500).
func internalFailure(msg string) *NeptuneDataError {
	return NewNeptuneDataError("InternalFailureException", msg, http.StatusInternalServerError)
}

// timeLimitExceeded returns a TimeLimitExceededException (503).
func timeLimitExceeded(msg string) *NeptuneDataError {
	return NewNeptuneDataError("TimeLimitExceededException", msg, http.StatusServiceUnavailable)
}

// constraintViolation returns a ConstraintViolationException (409).
func constraintViolation(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ConstraintViolationException", msg, http.StatusConflict)
}

// invalidNumericData returns an InvalidNumericDataException (400).
func invalidNumericData(msg string) *NeptuneDataError {
	return NewNeptuneDataError("InvalidNumericDataException", msg, http.StatusBadRequest)
}

// parsingError returns a ParsingException (400).
func parsingError(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ParsingException", msg, http.StatusBadRequest)
}

// readOnlyViolation returns a ReadOnlyViolationException (400).
func readOnlyViolation(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ReadOnlyViolationException", msg, http.StatusBadRequest)
}

// cancelledByUser returns a CancelledByUserException (400).
func cancelledByUser(msg string) *NeptuneDataError {
	return NewNeptuneDataError("CancelledByUserException", msg, http.StatusBadRequest)
}

// concurrentModification returns a ConcurrentModificationException (409).
func concurrentModification(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ConcurrentModificationException", msg, http.StatusConflict)
}

// clientTimeout returns a ClientTimeoutException (408).
func clientTimeout(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ClientTimeoutException", msg, http.StatusRequestTimeout)
}

// statisticsNotAvailable returns a StatisticsNotAvailableException (404).
func statisticsNotAvailable(msg string) *NeptuneDataError {
	return NewNeptuneDataError("StatisticsNotAvailableException", msg, http.StatusNotFound)
}

// s3Exception returns an S3Exception (400).
func s3Exception(msg string) *NeptuneDataError {
	return NewNeptuneDataError("S3Exception", msg, http.StatusBadRequest)
}

// illegalArgument returns an IllegalArgumentException (400).
func illegalArgument(msg string) *NeptuneDataError {
	return NewNeptuneDataError("IllegalArgumentException", msg, http.StatusBadRequest)
}

// invalidArgument returns an InvalidArgumentException (400).
func invalidArgument(msg string) *NeptuneDataError {
	return NewNeptuneDataError("InvalidArgumentException", msg, http.StatusBadRequest)
}

// memoryLimitExceeded returns a MemoryLimitExceededException (500).
func memoryLimitExceeded(msg string) *NeptuneDataError {
	return NewNeptuneDataError("MemoryLimitExceededException", msg, http.StatusInternalServerError)
}

// queryTooLarge returns a QueryTooLargeException (400).
func queryTooLarge(msg string) *NeptuneDataError {
	return NewNeptuneDataError("QueryTooLargeException", msg, http.StatusBadRequest)
}

// tooManyRequests returns a TooManyRequestsException (429).
func tooManyRequests(msg string) *NeptuneDataError {
	return NewNeptuneDataError("TooManyRequestsException", msg, http.StatusTooManyRequests)
}

// serverShutdown returns a ServerShutdownException (500).
func serverShutdown(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ServerShutdownException", msg, http.StatusInternalServerError)
}

// loadUrlAccessDenied returns a LoadUrlAccessDeniedException (400).
func loadUrlAccessDenied(msg string) *NeptuneDataError {
	return NewNeptuneDataError("LoadUrlAccessDeniedException", msg, http.StatusBadRequest)
}

// expiredStream returns an ExpiredStreamException (400).
func expiredStream(msg string) *NeptuneDataError {
	return NewNeptuneDataError("ExpiredStreamException", msg, http.StatusBadRequest)
}

// streamRecordsNotFound returns a StreamRecordsNotFoundException (404).
func streamRecordsNotFound(msg string) *NeptuneDataError {
	return NewNeptuneDataError("StreamRecordsNotFoundException", msg, http.StatusNotFound)
}

// mlResourceNotFound returns an MLResourceNotFoundException (404).
func mlResourceNotFound(msg string) *NeptuneDataError {
	return NewNeptuneDataError("MLResourceNotFoundException", msg, http.StatusNotFound)
}
