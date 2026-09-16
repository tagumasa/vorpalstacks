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

	// ErrReplayCapReached is returned by CreateReplayCapped when the store
	// already holds the maximum of active (non-terminal) replays — the
	// documented concurrent-replay cap. The service maps it to
	// LimitExceededException.
	ErrReplayCapReached = errors.New("active replay cap reached")

	// ErrRuleCapReached is returned by CreateRuleCapped when creating a
	// NEW rule would push its event bus past the documented per-bus
	// rule-count quota. The service maps it to
	// LimitExceededException.
	ErrRuleCapReached = errors.New("rule count cap reached")

	// ErrEmptyResourceName is returned by the Create* methods when the
	// resource name is empty — a write under an empty key would create a
	// phantom record, so the store enforces the invariant itself. The
	// service cores reject empty names before reaching the store; this
	// sentinel is the store-level backstop.
	ErrEmptyResourceName = errors.New("resource name is required")
)
