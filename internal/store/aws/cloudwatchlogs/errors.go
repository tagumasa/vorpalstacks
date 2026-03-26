package cloudwatchlogs

import "errors"

var (
	// ErrLogGroupNotFound is returned when the specified CloudWatch Logs
	// log group does not exist.
	ErrLogGroupNotFound = errors.New("log group not found")

	// ErrLogGroupAlreadyExists is returned when attempting to create a log group
	// that already exists.
	ErrLogGroupAlreadyExists = errors.New("log group already exists")

	// ErrLogStreamNotFound is returned when the specified log stream
	// does not exist.
	ErrLogStreamNotFound = errors.New("log stream not found")

	// ErrLogStreamAlreadyExists is returned when attempting to create a log stream
	// that already exists.
	ErrLogStreamAlreadyExists = errors.New("log stream already exists")

	// ErrInvalidParameter is returned when a parameter is not valid.
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrResourceAlreadyExists is returned when attempting to create a resource
	// that already exists.
	ErrResourceAlreadyExists = errors.New("resource already exists")

	// ErrResourceNotFound is returned when the specified resource does not exist.
	ErrResourceNotFound = errors.New("resource not found")

	// ErrLimitExceeded is returned when a service limit has been exceeded.
	ErrLimitExceeded = errors.New("limit exceeded")

	// ErrDataAlreadyAccepted is returned when the log data has already been
	// accepted (for put events operations).
	ErrDataAlreadyAccepted = errors.New("data already accepted")

	// ErrInvalidSequenceToken is returned when the sequence token is not valid
	// for the log stream.
	ErrInvalidSequenceToken = errors.New("invalid sequence token")

	// ErrMetricFilterNotFound is returned when the specified metric filter
	// does not exist.
	ErrMetricFilterNotFound = errors.New("metric filter not found")

	// ErrMetricFilterAlreadyExists is returned when attempting to create a
	// metric filter that already exists.
	ErrMetricFilterAlreadyExists = errors.New("metric filter already exists")

	// ErrSubscriptionFilterNotFound is returned when the specified subscription filter
	// does not exist.
	ErrSubscriptionFilterNotFound = errors.New("subscription filter not found")

	// ErrDestinationNotFound is returned when the specified destination does not exist.
	ErrDestinationNotFound = errors.New("destination not found")

	// ErrDestinationAlreadyExists is returned when attempting to create a destination
	// that already exists.
	ErrDestinationAlreadyExists = errors.New("destination already exists")
)
