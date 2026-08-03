package errors

import (
	"net/http"
	"strings"

	"connectrpc.com/connect"
	stderrors "errors"
)

// StoreErrorToGRPC maps a store-layer error to an appropriate gRPC status code.
// "not found" → NotFound, "already exists"/"duplicate" → AlreadyExists,
// "invalid"/"empty"/"required" → InvalidArgument, everything else → Internal.
func StoreErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	code := connect.CodeInternal
	switch {
	case strings.Contains(msg, "not found"):
		code = connect.CodeNotFound
	case strings.Contains(msg, "already exists"), strings.Contains(msg, "duplicate"):
		code = connect.CodeAlreadyExists
	case strings.Contains(msg, "invalid"), strings.Contains(msg, "empty"), strings.Contains(msg, "required"):
		code = connect.CodeInvalidArgument
	}
	return connect.NewError(code, err)
}

// AWSErrorToGRPC maps a service-layer *AWSError to an appropriate gRPC status
// code using the canonical HTTPStatus field. Non-AWSError values are returned
// as CodeInternal — these represent unexpected infrastructure failures, not
// client-correctable conditions.
func AWSErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}
	var awsErr *AWSError
	if stderrors.As(err, &awsErr) {
		return connect.NewError(httpStatusToConnectCode(awsErr.HTTPStatus), err)
	}
	return connect.NewError(connect.CodeInternal, err)
}

// httpStatusToConnectCode converts an HTTP status code to the closest
// connect.Code equivalent, following the canonical mapping from
// https://grpc.io/docs/guides/status-codes/.
func httpStatusToConnectCode(status int) connect.Code {
	switch status {
	case http.StatusBadRequest:
		return connect.CodeInvalidArgument
	case http.StatusUnauthorized:
		return connect.CodeUnauthenticated
	case http.StatusForbidden:
		return connect.CodePermissionDenied
	case http.StatusNotFound:
		return connect.CodeNotFound
	case http.StatusConflict:
		return connect.CodeAborted
	case http.StatusTooManyRequests:
		return connect.CodeResourceExhausted
	case http.StatusInternalServerError:
		return connect.CodeInternal
	case http.StatusServiceUnavailable:
		return connect.CodeUnavailable
	default:
		return connect.CodeUnknown
	}
}
