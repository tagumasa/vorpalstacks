package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

import (
	"errors"

	"vorpalstacks/internal/store/aws/common"
)

// notFoundOr maps a storage-layer miss to the caller's not-found sentinel;
// every other error — an infrastructure fault — surfaces as itself, so the
// service layer answers InternalError instead of a phantom 404 for a record
// that exists.
func notFoundOr(err, sentinel error) error {
	if common.IsNotFound(err) {
		return sentinel
	}
	return err
}

// Every sentinel here has a production return site and a service-layer
// translation in the sns service's mapStoreError; parameter-shape validation
// lives in the service's Cores (the single validation path), so the store
// carries no invalid-parameter vocabulary of its own.
var (
	// ErrTopicNotFound is returned when the specified SNS topic does not exist.
	ErrTopicNotFound = errors.New("topic not found")

	// ErrSubscriptionNotFound is returned when the specified subscription
	// does not exist.
	ErrSubscriptionNotFound = errors.New("subscription not found")

	// ErrTopicLimitExceeded is returned when creating a topic past the
	// documented per-account topic quota for its type.
	ErrTopicLimitExceeded = errors.New("topic limit exceeded")

	// ErrSubscriptionLimitExceeded is returned when subscribing past the
	// documented per-topic subscription quota.
	ErrSubscriptionLimitExceeded = errors.New("subscription limit exceeded")

	// ErrFilterPolicyLimitExceeded is returned when attaching a filter
	// policy past the documented per-topic filter-policy quota.
	ErrFilterPolicyLimitExceeded = errors.New("filter policy limit exceeded")

	// ErrPermissionLabelExists is returned when AddPermission names a
	// label an existing statement on the topic already carries — the
	// model declares Label "A unique identifier for the new policy
	// statement", so a duplicate is a constraint violation, not an
	// overwrite.
	ErrPermissionLabelExists = errors.New("permission label already exists")

	// ErrInvalidToken is returned when the confirmation token is not valid
	// or has expired.
	ErrInvalidToken = errors.New("invalid confirmation token")

	// ErrPlatformApplicationNotFound is returned when the specified platform
	// application does not exist.
	ErrPlatformApplicationNotFound = errors.New("platform application not found")

	// ErrEndpointNotFound is returned when the specified platform endpoint
	// does not exist.
	ErrEndpointNotFound = errors.New("endpoint not found")

	// ErrPlatformApplicationAlreadyExists is returned when attempting to create
	// a platform application that already exists.
	ErrPlatformApplicationAlreadyExists = errors.New("platform application already exists")
)
