// Package events provides EventBridge storage functionality for vorpalstacks.
package eventbridge

import (
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// EventBusState represents the state of an EventBridge event bus.
type EventBusState string

// EventBusState constants define the possible states of an EventBridge event bus.
const (
	EventBusStateActive   EventBusState = "ACTIVE"
	EventBusStateCreating EventBusState = "CREATING"
	EventBusStateDeleting EventBusState = "DELETING"
)

// RuleState represents the state of an EventBridge rule.
type RuleState string

// RuleState constants define the possible states of an EventBridge rule.
const (
	RuleStateEnabled  RuleState = "ENABLED"
	RuleStateDisabled RuleState = "DISABLED"
)

// ArchiveState represents the state of an EventBridge archive.
type ArchiveState string

// ArchiveState constants define the possible states of an EventBridge archive.
const (
	ArchiveStateEnabled  ArchiveState = "ENABLED"
	ArchiveStateDisabled ArchiveState = "DISABLED"
	ArchiveStateCreating ArchiveState = "CREATING"
	ArchiveStateUpdating ArchiveState = "UPDATING"
	ArchiveStateDeleting ArchiveState = "DELETING"
)

// ConnectionState represents the state of an EventBridge connection.
type ConnectionState string

// ConnectionState constants define the possible states of an EventBridge connection.
const (
	ConnectionStateAuthorized   ConnectionState = "AUTHORIZED"
	ConnectionStateDeauthorized ConnectionState = "DEAUTHORIZED"
	ConnectionStateCreating     ConnectionState = "CREATING"
	ConnectionStateUpdating     ConnectionState = "UPDATING"
	ConnectionStateDeleting     ConnectionState = "DELETING"
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

// ReplayState constants define the possible states of an EventBridge replay.
const (
	ReplayStateStarting  ReplayState = "STARTING"
	ReplayStateRunning   ReplayState = "RUNNING"
	ReplayStateCancelled ReplayState = "CANCELLED"
	ReplayStateCompleted ReplayState = "COMPLETED"
	ReplayStateFailed    ReplayState = "FAILED"
)

// EventBus represents an EventBridge event bus.
type EventBus struct {
	Name           string      `json:"name"`
	ARN            string      `json:"arn"`
	Region         string      `json:"region"`
	AccountID      string      `json:"accountId"`
	Description    string      `json:"description,omitempty"`
	Policy         string      `json:"policy,omitempty"`
	Tags           []types.Tag `json:"tags,omitempty"`
	CreatedAt      time.Time   `json:"createdAt"`
	LastModifiedAt time.Time   `json:"lastModifiedAt,omitempty"`
}

// Rule represents an EventBridge rule.
type Rule struct {
	Name               string      `json:"name"`
	ARN                string      `json:"arn"`
	Region             string      `json:"region"`
	AccountID          string      `json:"accountId"`
	EventBusName       string      `json:"eventBusName"`
	Description        string      `json:"description,omitempty"`
	EventPattern       string      `json:"eventPattern,omitempty"`
	ScheduleExpression string      `json:"scheduleExpression,omitempty"`
	State              RuleState   `json:"state"`
	ManagedBy          string      `json:"managedBy,omitempty"`
	RoleARN            string      `json:"roleArn,omitempty"`
	Tags               []types.Tag `json:"tags,omitempty"`
	CreatedAt          time.Time   `json:"createdAt"`
	LastModifiedAt     time.Time   `json:"lastModifiedAt"`
}

// Target represents an EventBridge target.
type Target struct {
	ID               string            `json:"id"`
	RuleName         string            `json:"ruleName"`
	EventBusName     string            `json:"eventBusName"`
	ARN              string            `json:"arn"`
	Input            string            `json:"input,omitempty"`
	InputPath        string            `json:"inputPath,omitempty"`
	InputTransformer *InputTransformer `json:"inputTransformer,omitempty"`
	RoleARN          string            `json:"roleArn,omitempty"`
	DeadLetterConfig *DeadLetterConfig `json:"deadLetterConfig,omitempty"`
	RetryPolicy      *RetryPolicy      `json:"retryPolicy,omitempty"`
	SqsParameters    *SqsParameters    `json:"sqsParameters,omitempty"`
	HttpParameters   *HttpParameters   `json:"httpParameters,omitempty"`
	CreatedAt        time.Time         `json:"createdAt"`
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
	Name           string       `json:"name"`
	ARN            string       `json:"arn"`
	Region         string       `json:"region"`
	AccountID      string       `json:"accountId"`
	EventBusName   string       `json:"eventBusName"`
	EventSourceARN string       `json:"eventSourceArn"`
	Description    string       `json:"description,omitempty"`
	EventPattern   string       `json:"eventPattern,omitempty"`
	RetentionDays  int32        `json:"retentionDays,omitempty"`
	State          ArchiveState `json:"state"`
	EventCount     int64        `json:"eventCount"`
	SizeBytes      int64        `json:"sizeBytes"`
	CreatedAt      time.Time    `json:"createdAt"`
}

// Connection represents an EventBridge connection.
type Connection struct {
	Name              string          `json:"name"`
	ARN               string          `json:"arn"`
	Region            string          `json:"region"`
	AccountID         string          `json:"accountId"`
	Description       string          `json:"description,omitempty"`
	AuthorizationType string          `json:"authorizationType"`
	State             ConnectionState `json:"state"`
	StateReason       string          `json:"stateReason,omitempty"`
	Tags              []types.Tag     `json:"tags,omitempty"`
	CreatedAt         time.Time       `json:"createdAt"`
	LastModifiedAt    time.Time       `json:"lastModifiedAt,omitempty"`
	LastAuthorizedAt  time.Time       `json:"lastAuthorizedAt,omitempty"`
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
	Tags                         []types.Tag         `json:"tags,omitempty"`
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
}

// ReplayDestination represents the destination configuration for an EventBridge replay.
type ReplayDestination struct {
	Arn        string   `json:"arn"`
	FilterArns []string `json:"filterArns,omitempty"`
}

// ArchivedEvent represents an archived EventBridge event.
type ArchivedEvent struct {
	ID          string                 `json:"id"`
	EventBusARN string                 `json:"eventBusArn"`
	Event       map[string]interface{} `json:"event"`
	Timestamp   time.Time              `json:"timestamp"`
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
