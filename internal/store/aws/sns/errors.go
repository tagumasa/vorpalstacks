package sns

// Package sns provides SNS (Simple Notification Service) data store implementations
// for vorpalstacks.

import "errors"

var (
	// ErrTopicNotFound is returned when the specified SNS topic does not exist.
	ErrTopicNotFound = errors.New("topic not found")

	// ErrTopicAlreadyExists is returned when attempting to create a topic
	// that already exists.
	ErrTopicAlreadyExists = errors.New("topic already exists")

	// ErrSubscriptionNotFound is returned when the specified subscription
	// does not exist.
	ErrSubscriptionNotFound = errors.New("subscription not found")

	// ErrSubscriptionAlreadyExists is returned when attempting to create a
	// subscription that already exists.
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")

	// ErrInvalidParameter is returned when a required parameter is missing
	// or has an invalid value.
	ErrInvalidParameter = errors.New("invalid parameter")

	// ErrInvalidTopicName is returned when the topic name does not meet
	// SNS naming requirements.
	ErrInvalidTopicName = errors.New("invalid topic name")

	// ErrInvalidProtocol is returned when the specified protocol is not
	// valid (e.g., email, SMS, HTTP, etc.).
	ErrInvalidProtocol = errors.New("invalid protocol")

	// ErrInvalidEndpoint is returned when the endpoint is not valid for
	// the specified protocol.
	ErrInvalidEndpoint = errors.New("invalid endpoint")

	// ErrConfirmationFailed is returned when the subscription confirmation
	// fails.
	ErrConfirmationFailed = errors.New("subscription confirmation failed")

	// ErrInvalidToken is returned when the confirmation token is not valid
	// or has expired.
	ErrInvalidToken = errors.New("invalid confirmation token")

	// ErrMessageTooLarge is returned when the message exceeds the maximum
	// allowed size for SNS.
	ErrMessageTooLarge = errors.New("message too large")

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
