package errors

import (
	"strings"

	"connectrpc.com/connect"
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
