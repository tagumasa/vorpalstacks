package cloudfront

import (
	awserrors "vorpalstacks/internal/common/errors"
)

// requireID returns the standard InvalidArgument error when an
// operation's Id parameter is missing.
func requireID(id string) error {
	if id == "" {
		return invalidArgument("Id is required")
	}
	return nil
}

// verifyIfMatch enforces optimistic concurrency for update and delete
// operations: the If-Match header must be present, and unless it is the
// wildcard "*" it must match the resource's current ETag.
func verifyIfMatch(ifMatch, etag string) error {
	if ifMatch == "" {
		return awserrors.NewAWSError("InvalidIfMatchVersion",
			"The If-Match version is missing or not valid", 400)
	}
	if ifMatch != "*" && ifMatch != etag {
		return awserrors.NewAWSError("PreconditionFailed", preconditionFailedETagMsg, 412)
	}
	return nil
}

// ensureNameAvailable rejects a rename when another resource of the same
// type already uses the new name. Keeping a resource's current name is
// always allowed.
func ensureNameAvailable(currentName, newName string, nameTaken func(name string) bool, dupErr error) error {
	if newName == currentName {
		return nil
	}
	if nameTaken(newName) {
		return dupErr
	}
	return nil
}
