package cloudtrail

import "errors"

var (
	// ErrTrailNotFound is returned when the specified CloudTrail trail
	// does not exist.
	ErrTrailNotFound = errors.New("trail not found")

	// ErrTrailAlreadyExists is returned when attempting to create a trail
	// that already exists.
	ErrTrailAlreadyExists = errors.New("trail already exists")

	// ErrInvalidTrailName is returned when the trail name is not valid.
	ErrInvalidTrailName = errors.New("invalid trail name")

	// ErrEventNotFound is returned when the specified CloudTrail event
	// does not exist.
	ErrEventNotFound = errors.New("event not found")

	// ErrEventDataStoreNotFound is returned when the specified event data
	// store does not exist.
	ErrEventDataStoreNotFound = errors.New("event data store not found")

	// ErrEventDataStoreAlreadyExists is returned when attempting to create
	// an event data store with a name that is already in use.
	ErrEventDataStoreAlreadyExists = errors.New("event data store already exists")

	// ErrEventDataStoreNotPendingDeletion is returned when attempting to
	// restore an event data store that is not in PENDING_DELETION state.
	ErrEventDataStoreNotPendingDeletion = errors.New("event data store is not pending deletion")

	// ErrQueryNotFound is returned when the specified query does not exist.
	ErrQueryNotFound = errors.New("query not found")

	// ErrChannelNotFound is returned when the specified channel does not exist.
	ErrChannelNotFound = errors.New("channel not found")

	// ErrChannelAlreadyExists is returned when a channel name already exists.
	ErrChannelAlreadyExists = errors.New("channel already exists")

	// ErrChannelSourceInUse is returned when a channel source is already
	// carried by another channel: a maximum of one channel exists per
	// source.
	ErrChannelSourceInUse = errors.New("channel source already in use")

	// ErrEventConfigurationNotFound is returned when no event configuration
	// is stored for the requested trail or event data store.
	ErrEventConfigurationNotFound = errors.New("event configuration not found")

	// ErrInvalidNextToken is returned when a LookupEvents NextToken is not
	// a token this store issued — undecodable, or from a different
	// parameterisation.
	ErrInvalidNextToken = errors.New("invalid next token")

	// ErrMaxConcurrentQueries is returned when a query admission would
	// exceed the store's bound on concurrently running queries.
	ErrMaxConcurrentQueries = errors.New("maximum number of concurrent queries exceeded")
)
