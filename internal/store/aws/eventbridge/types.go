package eventbridge

import (
	"sort"
	"time"
)

// RuleState represents the state of an EventBridge rule.
type RuleState string

// RuleState constants define the possible states of an EventBridge rule.
const (
	RuleStateEnabled                                  RuleState = "ENABLED"
	RuleStateDisabled                                 RuleState = "DISABLED"
	RuleStateEnabledWithAllCloudtrailManagementEvents RuleState = "ENABLED_WITH_ALL_CLOUDTRAIL_MANAGEMENT_EVENTS"
)

// ArchiveState represents the state of an EventBridge archive.
type ArchiveState string

// ArchiveState constants are the complete ArchiveState enum from the
// Smithy model (2015-10-07): CREATE_FAILED, CREATING, DISABLED, ENABLED,
// UPDATE_FAILED, UPDATING.
const (
	ArchiveStateEnabled      ArchiveState = "ENABLED"
	ArchiveStateDisabled     ArchiveState = "DISABLED"
	ArchiveStateCreating     ArchiveState = "CREATING"
	ArchiveStateUpdating     ArchiveState = "UPDATING"
	ArchiveStateCreateFailed ArchiveState = "CREATE_FAILED"
	ArchiveStateUpdateFailed ArchiveState = "UPDATE_FAILED"
)

// ConnectionState represents the state of an EventBridge connection.
type ConnectionState string

// ConnectionState constants are the complete ConnectionState enum from
// the Smithy model (2015-10-07): ACTIVE, AUTHORIZED, AUTHORIZING,
// CREATING, DEAUTHORIZED, DEAUTHORIZING, DELETING, FAILED_CONNECTIVITY,
// UPDATING.
const (
	ConnectionStateAuthorized         ConnectionState = "AUTHORIZED"
	ConnectionStateDeauthorized       ConnectionState = "DEAUTHORIZED"
	ConnectionStateCreating           ConnectionState = "CREATING"
	ConnectionStateUpdating           ConnectionState = "UPDATING"
	ConnectionStateDeleting           ConnectionState = "DELETING"
	ConnectionStateActive             ConnectionState = "ACTIVE"
	ConnectionStateAuthorizing        ConnectionState = "AUTHORIZING"
	ConnectionStateDeauthorizing      ConnectionState = "DEAUTHORIZING"
	ConnectionStateFailedConnectivity ConnectionState = "FAILED_CONNECTIVITY"
)

// ApiDestinationState represents the state of an EventBridge API destination.
type ApiDestinationState string

// ApiDestinationState constants define the possible states of an EventBridge API destination.
const (
	ApiDestinationStateActive   ApiDestinationState = "ACTIVE"
	ApiDestinationStateInactive ApiDestinationState = "INACTIVE"
)

// ReplayState represents the state of an EventBridge replay.
type ReplayState string

// ReplayState constants are the complete ReplayState enum from the Smithy
// model (2015-10-07): CANCELLED, CANCELLING, COMPLETED, FAILED, RUNNING,
// STARTING.
const (
	ReplayStateStarting   ReplayState = "STARTING"
	ReplayStateRunning    ReplayState = "RUNNING"
	ReplayStateCancelling ReplayState = "CANCELLING"
	ReplayStateCancelled  ReplayState = "CANCELLED"
	ReplayStateCompleted  ReplayState = "COMPLETED"
	ReplayStateFailed     ReplayState = "FAILED"
)

// The state-membership maps below hold each enum's complete model
// vocabulary; list filters validate against them so an out-of-enum state
// is a ValidationException rather than a silent empty match.
var (
	archiveStateMembers = map[ArchiveState]bool{
		ArchiveStateEnabled:      true,
		ArchiveStateDisabled:     true,
		ArchiveStateCreating:     true,
		ArchiveStateUpdating:     true,
		ArchiveStateCreateFailed: true,
		ArchiveStateUpdateFailed: true,
	}
	connectionStateMembers = map[ConnectionState]bool{
		ConnectionStateAuthorized:         true,
		ConnectionStateDeauthorized:       true,
		ConnectionStateCreating:           true,
		ConnectionStateUpdating:           true,
		ConnectionStateDeleting:           true,
		ConnectionStateActive:             true,
		ConnectionStateAuthorizing:        true,
		ConnectionStateDeauthorizing:      true,
		ConnectionStateFailedConnectivity: true,
	}
	replayStateMembers = map[ReplayState]bool{
		ReplayStateStarting:   true,
		ReplayStateRunning:    true,
		ReplayStateCancelling: true,
		ReplayStateCancelled:  true,
		ReplayStateCompleted:  true,
		ReplayStateFailed:     true,
	}
)

// IsValidArchiveState reports whether s is a member of the model's
// ArchiveState vocabulary.
func IsValidArchiveState(s ArchiveState) bool { return archiveStateMembers[s] }

// IsValidConnectionState reports whether s is a member of the model's
// ConnectionState vocabulary.
func IsValidConnectionState(s ConnectionState) bool { return connectionStateMembers[s] }

// IsValidReplayState reports whether s is a member of the model's
// ReplayState vocabulary.
func IsValidReplayState(s ReplayState) bool { return replayStateMembers[s] }

// sortedVocabulary renders an enum-membership map as a sorted value list
// for validation messages.
func sortedVocabulary[T ~string](m map[T]bool) []string {
	values := make([]string, 0, len(m))
	for v := range m {
		values = append(values, string(v))
	}
	sort.Strings(values)
	return values
}

// ArchiveStateVocabulary returns the sorted ArchiveState vocabulary.
func ArchiveStateVocabulary() []string { return sortedVocabulary(archiveStateMembers) }

// ConnectionStateVocabulary returns the sorted ConnectionState vocabulary.
func ConnectionStateVocabulary() []string { return sortedVocabulary(connectionStateMembers) }

// ReplayStateVocabulary returns the sorted ReplayState vocabulary.
func ReplayStateVocabulary() []string { return sortedVocabulary(replayStateMembers) }

// EventBus represents an EventBridge event bus. Resource tags live in the
// TagStore keyed by ARN; the record carries no tag mirror.
type EventBus struct {
	Name             string            `json:"name"`
	ARN              string            `json:"arn"`
	Region           string            `json:"region"`
	AccountID        string            `json:"accountId"`
	Description      string            `json:"description,omitempty"`
	Policy           string            `json:"policy,omitempty"`
	KmsKeyIdentifier string            `json:"kmsKeyIdentifier,omitempty"`
	DeadLetterConfig *DeadLetterConfig `json:"deadLetterConfig,omitempty"`
	LogConfig        *BusLogConfig     `json:"logConfig,omitempty"`
	CreatedAt        time.Time         `json:"createdAt"`
	LastModifiedAt   time.Time         `json:"lastModifiedAt,omitempty"`
}

// BusLogConfig configures the EventBridge bus-level logging destination.
// Per Smithy LogConfig shape: IncludeDetail (NONE|FULL) and Level
// (OFF|ERROR|INFO|TRACE).
type BusLogConfig struct {
	IncludeDetail string `json:"includeDetail,omitempty"`
	Level         string `json:"level,omitempty"`
}

// Rule represents an EventBridge rule. Resource tags live in the TagStore
// keyed by ARN; the record carries no tag mirror.
type Rule struct {
	Name               string    `json:"name"`
	ARN                string    `json:"arn"`
	Region             string    `json:"region"`
	AccountID          string    `json:"accountId"`
	EventBusName       string    `json:"eventBusName"`
	Description        string    `json:"description,omitempty"`
	EventPattern       string    `json:"eventPattern,omitempty"`
	ScheduleExpression string    `json:"scheduleExpression,omitempty"`
	State              RuleState `json:"state"`
	ManagedBy          string    `json:"managedBy,omitempty"`
	RoleARN            string    `json:"roleArn,omitempty"`
	CreatedBy          string    `json:"createdBy,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	LastModifiedAt     time.Time `json:"lastModifiedAt"`
	// LastFiredAt records the most recent schedule boundary this rule
	// fired under. It is an internal durability marker for the scheduler
	// worker (a restart re-seeds its dedup cache from it); it never
	// appears in API responses, and recording it does not count as a
	// rule modification.
	LastFiredAt time.Time `json:"lastFiredAt"`
}

// Target represents an EventBridge target.
type Target struct {
	ID                string             `json:"id"`
	RuleName          string             `json:"ruleName"`
	EventBusName      string             `json:"eventBusName"`
	ARN               string             `json:"arn"`
	Input             string             `json:"input,omitempty"`
	InputPath         string             `json:"inputPath,omitempty"`
	InputTransformer  *InputTransformer  `json:"inputTransformer,omitempty"`
	RoleARN           string             `json:"roleArn,omitempty"`
	DeadLetterConfig  *DeadLetterConfig  `json:"deadLetterConfig,omitempty"`
	RetryPolicy       *RetryPolicy       `json:"retryPolicy,omitempty"`
	SqsParameters     *SqsParameters     `json:"sqsParameters,omitempty"`
	HttpParameters    *HttpParameters    `json:"httpParameters,omitempty"`
	KinesisParameters *KinesisParameters `json:"kinesisParameters,omitempty"`
	AppSyncParameters *AppSyncParameters `json:"appSyncParameters,omitempty"`
	CreatedAt         time.Time          `json:"createdAt"`
}

// KinesisParameters represents the Kinesis parameters for an EventBridge target.
type KinesisParameters struct {
	PartitionKeyPath string `json:"partitionKeyPath,omitempty"`
}

// AppSyncParameters configures an AppSync GraphQL target.
type AppSyncParameters struct {
	GraphQLOperation string `json:"graphQLOperation,omitempty"`
}

// InputTransformer represents an input transformer for EventBridge targets.
type InputTransformer struct {
	InputPathsMap map[string]string `json:"inputPathsMap,omitempty"`
	InputTemplate string            `json:"inputTemplate,omitempty"`
}

// DeadLetterConfig represents the dead letter configuration for an EventBridge target.
type DeadLetterConfig struct {
	Arn string `json:"arn,omitempty"`
}

// RetryPolicy represents the retry policy for an EventBridge target.
type RetryPolicy struct {
	MaximumEventAgeInSeconds int32 `json:"maximumEventAgeInSeconds,omitempty"`
	MaximumRetryAttempts     int32 `json:"maximumRetryAttempts,omitempty"`
}

// SqsParameters represents the SQS parameters for an EventBridge target.
type SqsParameters struct {
	MessageGroupId string `json:"messageGroupId,omitempty"`
}

// HttpParameters represents the HTTP parameters for an EventBridge target.
type HttpParameters struct {
	HeaderParameters      map[string]string `json:"headerParameters,omitempty"`
	PathParameterValues   []string          `json:"pathParameterValues,omitempty"`
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
}

// Archive represents an EventBridge archive.
type Archive struct {
	Name             string       `json:"name"`
	ARN              string       `json:"arn"`
	Region           string       `json:"region"`
	AccountID        string       `json:"accountId"`
	EventBusName     string       `json:"eventBusName"`
	EventSourceARN   string       `json:"eventSourceArn"`
	Description      string       `json:"description,omitempty"`
	EventPattern     string       `json:"eventPattern,omitempty"`
	RetentionDays    int32        `json:"retentionDays,omitempty"`
	KmsKeyIdentifier string       `json:"kmsKeyIdentifier,omitempty"`
	State            ArchiveState `json:"state"`
	StateReason      string       `json:"stateReason,omitempty"`
	EventCount       int64        `json:"eventCount"`
	SizeBytes        int64        `json:"sizeBytes"`
	CreatedAt        time.Time    `json:"createdAt"`
}

// ConnectionHeaderParameter is one additional header parameter a connection
// applies to its HTTP requests (Smithy ConnectionHeaderParameter: Key, Value,
// IsValueSecret — "up to 100 additional header parameters per request").
type ConnectionHeaderParameter struct {
	Key           string `json:"key,omitempty"`
	Value         string `json:"value,omitempty"`
	IsValueSecret bool   `json:"isValueSecret,omitempty"`
}

// ConnectionQueryStringParameter is one additional query string parameter a
// connection applies to its HTTP requests (Smithy
// ConnectionQueryStringParameter: Key, Value, IsValueSecret).
type ConnectionQueryStringParameter struct {
	Key           string `json:"key,omitempty"`
	Value         string `json:"value,omitempty"`
	IsValueSecret bool   `json:"isValueSecret,omitempty"`
}

// ConnectionBodyParameter is one additional body parameter a connection
// applies to its HTTP requests (Smithy ConnectionBodyParameter: Key, Value,
// IsValueSecret).
type ConnectionBodyParameter struct {
	Key           string `json:"key,omitempty"`
	Value         string `json:"value,omitempty"`
	IsValueSecret bool   `json:"isValueSecret,omitempty"`
}

// ConnectionHttpParameters represents additional header/query/body parameters
// to send on the HTTP request when invoking an API destination target. Each
// family is a list of {Key, Value, IsValueSecret} entries per the Smithy
// ConnectionHttpParameters shape (each list bounded at 100 entries).
type ConnectionHttpParameters struct {
	HeaderParameters      []ConnectionHeaderParameter      `json:"headerParameters,omitempty"`
	QueryStringParameters []ConnectionQueryStringParameter `json:"queryStringParameters,omitempty"`
	BodyParameters        []ConnectionBodyParameter        `json:"bodyParameters,omitempty"`
}

// BasicAuthParameters holds Basic HTTP authentication credentials.
type BasicAuthParameters struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// OAuthClientParameters holds the OAuth client credentials used to obtain an
// access token from the authorization endpoint.
type OAuthClientParameters struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// OAuthParameters holds the OAuth client-credentials flow configuration.
type OAuthParameters struct {
	ClientParameters      *OAuthClientParameters    `json:"clientParameters,omitempty"`
	AuthorizationEndpoint string                    `json:"authorizationEndpoint,omitempty"`
	HttpMethod            string                    `json:"httpMethod,omitempty"`
	OAuthHttpParameters   *ConnectionHttpParameters `json:"oauthHttpParameters,omitempty"`
}

// ApiKeyAuthParameters holds the API key header used to authenticate the
// request to the API destination.
type ApiKeyAuthParameters struct {
	ApiKeyName  string `json:"apiKeyName"`
	ApiKeyValue string `json:"apiKeyValue"`
}

// AuthParameters holds the per-authorization-type credentials for a connection.
// Only one of BasicAuthParameters / OAuthParameters / ApiKeyAuthParameters is
// populated, matching the value of Connection.AuthorizationType.
type AuthParameters struct {
	BasicAuthParameters      *BasicAuthParameters      `json:"basicAuthParameters,omitempty"`
	OAuthParameters          *OAuthParameters          `json:"oauthParameters,omitempty"`
	ApiKeyAuthParameters     *ApiKeyAuthParameters     `json:"apiKeyAuthParameters,omitempty"`
	InvocationHttpParameters *ConnectionHttpParameters `json:"invocationHttpParameters,omitempty"`
}

// ConnectivityResourceParameters wraps the ResourceParameters sub-structure
// per the Smithy shape. The top-level request member
// InvocationConnectivityParameters maps to this shape; it contains a
// ResourceParameters member that itself carries the resource configuration
// ARN (and resource association ARN on responses).
type ConnectivityResourceParameters struct {
	ResourceParameters *ResourceConfiguration `json:"resourceParameters,omitempty"`
}

// ResourceConfiguration corresponds to the Smithy
// ConnectivityResourceConfigurationArn (request) and
// DescribeConnectionResourceParameters (response) shapes. The
// ResourceAssociationArn field is populated only on Describe responses.
type ResourceConfiguration struct {
	ResourceConfigurationArn string `json:"resourceConfigurationArn,omitempty"`
	ResourceAssociationArn   string `json:"resourceAssociationArn,omitempty"`
}

// Connection represents an EventBridge connection.
type Connection struct {
	Name                             string                          `json:"name"`
	ARN                              string                          `json:"arn"`
	Region                           string                          `json:"region"`
	AccountID                        string                          `json:"accountId"`
	Description                      string                          `json:"description,omitempty"`
	AuthorizationType                string                          `json:"authorizationType"`
	AuthParameters                   *AuthParameters                 `json:"authParameters,omitempty"`
	InvocationConnectivityParameters *ConnectivityResourceParameters `json:"invocationConnectivityParameters,omitempty"`
	KmsKeyIdentifier                 string                          `json:"kmsKeyIdentifier,omitempty"`
	SecretArn                        string                          `json:"secretArn,omitempty"`
	State                            ConnectionState                 `json:"state"`
	StateReason                      string                          `json:"stateReason,omitempty"`
	CreatedAt                        time.Time                       `json:"createdAt"`
	LastModifiedAt                   time.Time                       `json:"lastModifiedAt,omitempty"`
	LastAuthorizedAt                 time.Time                       `json:"lastAuthorizedAt,omitempty"`
}

// ApiDestination represents an EventBridge API destination.
type ApiDestination struct {
	Name                         string              `json:"name"`
	ARN                          string              `json:"arn"`
	Region                       string              `json:"region"`
	AccountID                    string              `json:"accountId"`
	ConnectionARN                string              `json:"connectionArn"`
	Description                  string              `json:"description,omitempty"`
	HttpMethod                   string              `json:"httpMethod"`
	InvocationEndpoint           string              `json:"invocationEndpoint"`
	InvocationRateLimitPerSecond int32               `json:"invocationRateLimitPerSecond,omitempty"`
	State                        ApiDestinationState `json:"state"`
	CreatedAt                    time.Time           `json:"createdAt"`
	LastModifiedAt               time.Time           `json:"lastModifiedAt,omitempty"`
}

// Event represents an EventBridge event.
type Event struct {
	ID           string                 `json:"id"`
	Version      string                 `json:"version"`
	DetailType   string                 `json:"detailType"`
	Source       string                 `json:"source"`
	Account      string                 `json:"account"`
	Time         time.Time              `json:"time"`
	Region       string                 `json:"region"`
	Resources    []string               `json:"resources,omitempty"`
	Detail       map[string]interface{} `json:"detail"`
	EventBusName string                 `json:"eventBusName"`
	TraceHeader  string                 `json:"traceHeader,omitempty"`
	// ReplayName carries the replay-name metadata field EventBridge adds to
	// an event while replaying it (the archives user guide documents the
	// field and the {"replay-name":[{"exists":false}]} managed-rule pattern
	// consumers use to ignore replays). It is a delivery-time property set
	// only by the replay path; PutEvents ingress cannot supply it.
	ReplayName string `json:"replayName,omitempty"`
	// IngestionTime records when the platform received the event: the
	// aws.events.event.ingestion-time reserved variable renders this
	// stamp — "The time at which the event was received by EventBridge.
	// This is an ISO 8601 timestamp. This variable is generated by
	// EventBridge and can't be overwritten" (input transformation page) —
	// never the publisher-suppliable Time member. Zero means the record
	// predates the stamp (a pre-index archive record, a synthetic event);
	// readers fall back to Time.
	IngestionTime time.Time `json:"ingestionTime,omitempty"`
}

// Replay represents an EventBridge replay.
type Replay struct {
	Name                  string             `json:"name"`
	ARN                   string             `json:"arn"`
	Region                string             `json:"region"`
	AccountID             string             `json:"accountId"`
	Description           string             `json:"description,omitempty"`
	State                 ReplayState        `json:"state"`
	StateReason           string             `json:"stateReason,omitempty"`
	EventSourceARN        string             `json:"eventSourceArn"`
	Destination           *ReplayDestination `json:"destination,omitempty"`
	EventStartTime        time.Time          `json:"eventStartTime"`
	EventEndTime          time.Time          `json:"eventEndTime"`
	EventLastReplayedTime time.Time          `json:"eventLastReplayedTime,omitempty"`
	ReplayStartTime       time.Time          `json:"replayStartTime,omitempty"`
	ReplayEndTime         time.Time          `json:"replayEndTime,omitempty"`
	// CreatedAt anchors the replay record's age: EventBridge deletes replay
	// records after 90 days (the archives user guide), counted from the
	// record's creation.
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// ReplayDestination represents the destination configuration for an EventBridge replay.
type ReplayDestination struct {
	Arn        string   `json:"arn"`
	FilterArns []string `json:"filterArns,omitempty"`
}

// ArchivedEvent represents an archived EventBridge event. The archived
// envelope is self-describing; the owning bus is recoverable from the
// archive record and is not mirrored per event.
type ArchivedEvent struct {
	ID        string                 `json:"id"`
	Event     map[string]interface{} `json:"event"`
	Timestamp time.Time              `json:"timestamp"`
}

// PutEventsRequestEntry represents an entry in a PutEvents request.
type PutEventsRequestEntry struct {
	Time         time.Time              `json:"time,omitempty"`
	Source       string                 `json:"source,omitempty"`
	Resources    []string               `json:"resources,omitempty"`
	DetailType   string                 `json:"detailType,omitempty"`
	Detail       map[string]interface{} `json:"detail,omitempty"`
	EventBusName string                 `json:"eventBusName,omitempty"`
	TraceHeader  string                 `json:"traceHeader,omitempty"`
}

// PutEventsResultEntry represents an entry in a PutEvents result.
type PutEventsResultEntry struct {
	EventId      string `json:"eventId,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// EventBusListResult represents the result of listing EventBridge event buses.
type EventBusListResult struct {
	EventBuses []*EventBus
	NextToken  string
}

// RuleListResult represents the result of listing EventBridge rules.
type RuleListResult struct {
	Rules     []*Rule
	NextToken string
}

// TargetListResult represents the result of listing EventBridge targets.
type TargetListResult struct {
	Targets   []*Target
	NextToken string
}
