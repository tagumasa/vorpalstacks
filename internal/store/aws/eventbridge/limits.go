package eventbridge

// Retry-policy bounds. The AWS EventBridge API reference (RetryPolicy)
// documents the valid ranges: MaximumRetryAttempts "Valid Range: Minimum
// value of 0. Maximum value of 185." and MaximumEventAgeInSeconds "Valid
// Range: Minimum value of 60. Maximum value of 86400." The vendored model
// carries both as modelled @range traits (MaximumRetryAttempts 0-185,
// MaximumEventAgeInSeconds 60-86400).
//
// The same values double as the production delivery defaults: "By default,
// EventBridge retries sending the event for 24 hours and up to 185 times
// with an exponential back off and jitter" (EventBridge user guide, "How
// EventBridge retries delivering events").
const (
	// RetryPolicyMaxRetryAttempts is the ceiling for
	// RetryPolicy.MaximumRetryAttempts and the production default retry
	// count.
	RetryPolicyMaxRetryAttempts = 185
	// RetryPolicyMaxEventAgeSeconds is the ceiling for
	// RetryPolicy.MaximumEventAgeInSeconds and the production default retry
	// window (24 hours).
	RetryPolicyMaxEventAgeSeconds = 86400
	// RetryPolicyMinEventAgeSeconds is the floor for
	// RetryPolicy.MaximumEventAgeInSeconds.
	RetryPolicyMinEventAgeSeconds = 60
)

// KmsKeyIdentifierMaxLength is the @length(0,2048) ceiling the model puts
// on the KmsKeyIdentifier shape (CreateConnection, CreateArchive,
// CreateEventBus), counted in Unicode scalar values like every modelled
// string length.
const KmsKeyIdentifierMaxLength = 2048

// PutEvents entry bounds. The entry count is the Smithy @length(1,10)
// trait on PutEventsRequestEntryList. The request size ceiling comes from
// the user guide's "Calculating PutEvents event entry size" section: the
// sum of all entries in one request must stay under 1 MB (1,048,576
// bytes), the limit applying to the request as a whole rather than to
// individual entries; the per-entry charge is 14 bytes for a specified
// Time plus the UTF-8 byte lengths of Source, DetailType, Detail (when
// specified) and every Resources member, with TraceHeader excluded ("The
// trace header doesn't count towards the PutEventsRequestEntry event
// size", X-Ray integration page). The vendored 2015-10-07 model's
// operation documentation still phrases the same request-total rule with
// the pre-raise 256 KB number — recorded drift; the live quota governs.
const (
	// MaxPutEventsEntries bounds the Entries list of one PutEvents
	// request.
	MaxPutEventsEntries = 10
	// MaxPutEventsRequestSizeBytes bounds the summed entry size of one
	// PutEvents request.
	MaxPutEventsRequestSizeBytes = 1048576
	// PutEventsTimeFieldSizeBytes is the documented per-entry charge of a
	// specified Time member.
	PutEventsTimeFieldSizeBytes = 14
)

// Target input-configuration bounds. The Smithy model carries these as
// mechanical traits — @length on TargetInput (0-8192), TargetInputPath
// (0-256, also the TransformerPaths value shape), TransformerInput
// (1-8192) and TransformerPaths (0-100 entries), @length 1-256 plus
// @pattern ^[A-Za-z0-9_\-]+$ on InputTransformerPathKey — and the API
// reference renders the same values as its Length Constraints lines.
const (
	// TargetInputMaxLength bounds Target.Input ("Valid JSON text passed to
	// the target").
	TargetInputMaxLength = 8192
	// TargetInputPathMaxLength bounds Target.InputPath and every
	// InputPathsMap path value (both target the TargetInputPath shape).
	TargetInputPathMaxLength = 256
	// InputTemplateMinLength is the floor for
	// InputTransformer.InputTemplate (a required member).
	InputTemplateMinLength = 1
	// InputTemplateMaxLength bounds InputTransformer.InputTemplate.
	InputTemplateMaxLength = 8192
	// InputPathsMapMaxEntries bounds the InputTransformer.InputPathsMap
	// entry count.
	InputPathsMapMaxEntries = 100
	// InputPathsMapKeyMaxLength bounds every InputPathsMap key.
	InputPathsMapKeyMaxLength = 256
)

// ListLimitMaximum is the upper bound of every list operation's Limit
// member: all eight of them (ListArchives, ListConnections, ListReplays,
// ListApiDestinations, ListEventBuses, ListRules, ListTargetsByRule,
// ListRuleNamesByTarget) target the LimitMax100 shape, @range(min=1,
// max=100) in the vendored model. An omitted Limit arrives as the int32
// zero value on the wire and defaults to this maximum.
const ListLimitMaximum = 100

// Replay lifecycle quotas. The archives user guide states both: "You can
// have a maximum of ten active concurrent replays per account per AWS
// Region." and "EventBridge deletes replays after 90 days." The Smithy
// model carries LimitExceededException on StartReplay for the concurrency
// breach; the 90-day deletion is wired as the retention sweep anchored on
// the replay record's creation time.
const (
	// MaxConcurrentReplays bounds the simultaneously active replays
	// (STARTING, RUNNING or CANCELLING) per account per region.
	MaxConcurrentReplays = 10
	// ReplayRetentionDays is the age at which a replay record is deleted
	// by the retention sweep, counted from its creation.
	ReplayRetentionDays = 90
)

// Rule-count quota. The EventBridge quotas page documents "Number of
// rules | ... Each of the other supported Regions: 300 | Adjustable |
// Maximum number of rules an account can have per event bus" (the
// af-south-1 and eu-south-1 defaults are 100; the platform enforces the
// general 300). The Smithy model carries LimitExceededException on
// PutRule for the breach.
const (
	// MaxRulesPerEventBus bounds how many rules one event bus may
	// carry; creating a further rule on a full bus answers
	// LimitExceededException, while updating an existing rule at the
	// quota stays allowed.
	MaxRulesPerEventBus = 300
)

// Connection HTTP-parameter bounds. The Smithy model carries these as
// mechanical traits: @length(0,100) on ConnectionHeaderParametersList,
// ConnectionQueryStringParametersList and ConnectionBodyParametersList,
// and @length(0,512) plus @pattern on HeaderKey, HeaderValueSensitive,
// QueryStringKey and QueryStringValueSensitive. ConnectionBodyParameter's
// Key and Value target unadorned String/SensitiveString shapes, so the
// body family carries the list cap alone.
const (
	// ConnectionHttpParametersMaxEntries bounds each of the header,
	// query-string and body parameter lists a connection may declare.
	ConnectionHttpParametersMaxEntries = 100
	// ConnectionParameterMaxLength bounds every modelled HTTP-parameter
	// key and value (header and query-string families).
	ConnectionParameterMaxLength = 512
)
