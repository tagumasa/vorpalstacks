package eventbridge

import "errors"

var (
	// ErrEventBusNotFound is returned when the specified EventBridge event bus
	// does not exist.
	ErrEventBusNotFound = errors.New("event bus not found")

	// ErrEventBusAlreadyExists is returned when attempting to create an event bus
	// that already exists.
	ErrEventBusAlreadyExists = errors.New("event bus already exists")

	// ErrRuleNotFound is returned when the specified EventBridge rule
	// does not exist.
	ErrRuleNotFound = errors.New("rule not found")

	// ErrRuleAlreadyExists is returned when attempting to create a rule
	// that already exists.
	ErrRuleAlreadyExists = errors.New("rule already exists")

	// ErrTargetNotFound is returned when the specified rule target
	// does not exist.
	ErrTargetNotFound = errors.New("target not found")

	// ErrArchiveNotFound is returned when the specified event archive
	// does not exist.
	ErrArchiveNotFound = errors.New("archive not found")

	// ErrArchiveAlreadyExists is returned when attempting to create an archive
	// that already exists.
	ErrArchiveAlreadyExists = errors.New("archive already exists")

	// ErrConnectionNotFound is returned when the specified EventBridge connection
	// does not exist.
	ErrConnectionNotFound = errors.New("connection not found")

	// ErrConnectionAlreadyExists is returned when attempting to create a connection
	// that already exists.
	ErrConnectionAlreadyExists = errors.New("connection already exists")

	// ErrApiDestinationNotFound is returned when the specified API destination
	// does not exist.
	ErrApiDestinationNotFound = errors.New("api destination not found")

	// ErrApiDestinationAlreadyExists is returned when attempting to create an API
	// destination that already exists.
	ErrApiDestinationAlreadyExists = errors.New("api destination already exists")

	// ErrReplayNotFound is returned when the specified replay does not exist.
	ErrReplayNotFound = errors.New("replay not found")

	// ErrReplayAlreadyExists is returned when attempting to create a replay
	// that already exists.
	ErrReplayAlreadyExists = errors.New("replay already exists")

	// ErrInvalidEventPattern is returned when the event pattern is not valid.
	ErrInvalidEventPattern = errors.New("invalid event pattern")

	// ErrInvalidARN is returned when the Amazon Resource Name (ARN) is not valid.
	ErrInvalidARN = errors.New("invalid ARN")
)
