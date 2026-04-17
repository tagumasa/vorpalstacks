package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Event defines the interface that all event bus events must implement.
type Event interface {
	EventType() string
	EventID() string
	EventTimestamp() time.Time
	EventSource() string
	EventRegion() string
	EventAccountID() string
	EventDepth() int
	SetEventDepth(depth int)
	EventCaller() CallerContext
}

// CallerContext represents the IAM identity that originated the event,
// including the invoking service principal and any assumed role.
type CallerContext struct {
	ServicePrincipal string `json:"service_principal"`
	SourceARN        string `json:"source_arn"`
	AssumedRoleARN   string `json:"assumed_role_arn"`
	AccountID        string `json:"account_id"`
}

// EventBase provides a default implementation of the Event interface,
// embedding common fields (ID, timestamp, source, region, account, depth, caller).
type EventBase struct {
	ID        string        `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Source    string        `json:"source"`
	Region    string        `json:"region"`
	AccountID string        `json:"account_id"`
	depth     atomic.Int32  `json:"-"`
	Caller    CallerContext `json:"caller"`
}

// EventType returns the event type identifier. Subtypes override this
// to supply their specific type string (e.g. "service:invoke").
func (e *EventBase) EventType() string { return "" }

// EventID returns the unique identifier for this event.
func (e *EventBase) EventID() string { return e.ID }

// EventTimestamp returns the time at which the event was created.
func (e *EventBase) EventTimestamp() time.Time { return e.Timestamp }

// EventSource returns the service that published the event.
func (e *EventBase) EventSource() string { return e.Source }

// EventRegion returns the AWS region associated with the event.
func (e *EventBase) EventRegion() string { return e.Region }

// EventAccountID returns the AWS account ID associated with the event.
func (e *EventBase) EventAccountID() string { return e.AccountID }

// EventDepth returns the current dispatch depth, used to prevent
// infinite event cycles.
func (e *EventBase) EventDepth() int { return int(e.depth.Load()) }

func (e *EventBase) SetEventDepth(d int) { e.depth.Store(int32(d)) }

func (e *EventBase) getEventBase() *EventBase { return e }

// EventCaller returns the IAM identity that originated the event.
func (e *EventBase) EventCaller() CallerContext { return e.Caller }

// ServiceInvokeRequest represents a synchronous service-to-service invocation
// dispatched through the event bus, carrying a target ARN, payload, and
// HTTP-style routing metadata.
type ServiceInvokeRequest struct {
	EventBase
	TargetARN         string              `json:"target_arn"`
	Payload           []byte              `json:"payload,omitempty"`
	Headers           map[string]string   `json:"headers,omitempty"`
	MultiValueHeaders map[string][]string `json:"multi_value_headers,omitempty"`
	QueryParams       map[string]string   `json:"query_params,omitempty"`
	PathParams        map[string]string   `json:"path_params,omitempty"`
	StageVariables    map[string]string   `json:"stage_variables,omitempty"`
	RequestContext    map[string]string   `json:"request_context,omitempty"`
}

// EventType returns "service:invoke" for this event type.
func (e *ServiceInvokeRequest) EventType() string { return "service:invoke" }

// SNSDeliveryEvent is published when an SNS message needs to be delivered
// to topic subscribers. MessageAttributes are serialised as raw JSON to
// avoid coupling the eventbus package to store-layer types.
type SNSDeliveryEvent struct {
	EventBase
	TopicARN          string                     `json:"topic_arn"`
	MessageID         string                     `json:"message_id"`
	Message           string                     `json:"message"`
	Subject           string                     `json:"subject,omitempty"`
	MessageGroupId    string                     `json:"message_group_id,omitempty"`
	MessageAttributes map[string]json.RawMessage `json:"message_attributes,omitempty"`
}

// EventType returns "sns:deliver" for this event type.
func (e *SNSDeliveryEvent) EventType() string { return "sns:deliver" }

// EventBridgeDeliveryEvent is published when EventBridge needs to deliver
// an event to a target. The input is already transformed by InputTransformer.
type EventBridgeDeliveryEvent struct {
	EventBase
	RuleARN   string `json:"rule_arn"`
	TargetID  string `json:"target_id"`
	TargetARN string `json:"target_arn"`
	Input     []byte `json:"input"`
}

// EventType returns "events:deliver" for this event type.
func (e *EventBridgeDeliveryEvent) EventType() string { return "events:deliver" }

// ScheduleFiredEvent is published when the Scheduler engine fires a schedule.
type ScheduleFiredEvent struct {
	EventBase
	ScheduleName string `json:"schedule_name"`
	ScheduleArn  string `json:"schedule_arn"`
	TargetArn    string `json:"target_arn"`
	Input        string `json:"input"`
}

// EventType returns "scheduler:fired" for this event type.
func (e *ScheduleFiredEvent) EventType() string { return "scheduler:fired" }

// CloudWatchLogDeliveryEvent is published when CloudWatch Logs subscription
// filters match log events. The payload is gzip-compressed JSON.
type CloudWatchLogDeliveryEvent struct {
	EventBase
	LogGroup       string `json:"log_group"`
	LogStream      string `json:"log_stream"`
	DestinationArn string `json:"destination_arn"`
	Payload        []byte `json:"payload"`
}

// EventType returns "logs:deliver" for this event type.
func (e *CloudWatchLogDeliveryEvent) EventType() string { return "logs:deliver" }

// S3ObjectOp represents the type of S3 object operation that triggered an event.
type S3ObjectOp string

// S3 object event type constants matching AWS event names.
const (
	S3ObjectCreatedPut                     S3ObjectOp = "ObjectCreated:Put"
	S3ObjectCreatedCopy                    S3ObjectOp = "ObjectCreated:Copy"
	S3ObjectCreatedCompleteMultipartUpload S3ObjectOp = "ObjectCreated:CompleteMultipartUpload"
	S3ObjectCreatedPost                    S3ObjectOp = "ObjectCreated:Post"
	S3ObjectRemovedDelete                  S3ObjectOp = "ObjectRemoved:Delete"
	S3ObjectRemovedDeleteMarkerCreated     S3ObjectOp = "ObjectRemoved:DeleteMarkerCreated"
)

// S3ObjectEvent is published after a successful S3 object operation (PutObject,
// CopyObject, DeleteObject, CompleteMultipartUpload). The handler reads the
// bucket's NotificationConfiguration and dispatches to configured destinations.
type S3ObjectEvent struct {
	EventBase
	Bucket    string     `json:"bucket"`
	Key       string     `json:"key"`
	Size      int64      `json:"size"`
	VersionID string     `json:"version_id"`
	ETag      string     `json:"etag"`
	Op        S3ObjectOp `json:"op"`
}

// EventType returns "s3:ObjectEvent" for this event type.
func (e *S3ObjectEvent) EventType() string { return "s3:ObjectEvent" }

// CloudWatchAlarmStateEvent is published when a CloudWatch alarm transitions
// between states (OK, ALARM, INSUFFICIENT_DATA). The handler resolves action
// ARNs and dispatches notifications to the configured targets.
type CloudWatchAlarmStateEvent struct {
	EventBase
	AlarmName     string `json:"alarm_name"`
	AlarmARN      string `json:"alarm_arn"`
	PreviousState string `json:"previous_state"`
	NewState      string `json:"new_state"`
	Reason        string `json:"reason"`
}

// EventType returns "cloudwatch:AlarmStateChange" for this event type.
func (e *CloudWatchAlarmStateEvent) EventType() string { return "cloudwatch:AlarmStateChange" }

// CognitoTriggerEvent is published synchronously during Cognito user pool
// lifecycle operations (SignUp, ConfirmSignUp, authentication, token
// generation, etc.). The handler invokes the configured Lambda trigger
// function identified by LambdaARN and returns the trigger response
// payload to the caller.
type CognitoTriggerEvent struct {
	EventBase
	TriggerSource string `json:"trigger_source"` // AWS trigger source, e.g. "PreSignUp_SignUp"
	UserPoolID    string `json:"user_pool_id"`
	Username      string `json:"username"`
	ClientID      string `json:"client_id"`
	LambdaARN     string `json:"lambda_arn"` // resolved Lambda ARN from pool LambdaConfig
	Version       int    `json:"version"`
	Payload       []byte `json:"payload"` // full AWS Cognito trigger payload JSON
}

// EventType returns "cognito:Trigger" for this event type.
func (e *CognitoTriggerEvent) EventType() string { return "cognito:Trigger" }

// SecretRotationEvent is published synchronously for each step of a Secrets
// Manager rotation (createSecret, setSecret, testSecret). The handler invokes
// the rotation Lambda with the step-specific payload.
type SecretRotationEvent struct {
	EventBase
	SecretId  string `json:"secret_id"`
	Step      string `json:"step"`
	VersionId string `json:"version_id"`
}

// EventType returns "secretsmanager:Rotation" for this event type.
func (e *SecretRotationEvent) EventType() string { return "secretsmanager:Rotation" }

// DynamoDBRecordEvent is published when DynamoDB Streams records are available
// for a Lambda event source mapping. Each event carries one or more change
// records from the source table's stream.
type DynamoDBRecordEvent struct {
	EventBase
	TableName string                 `json:"table_name"`
	StreamArn string                 `json:"stream_arn"`
	Records   []DynamoDBChangeRecord `json:"records"`
}

// DynamoDBChangeRecord represents a single data modification captured in a
// DynamoDB stream (INSERT, MODIFY, REMOVE).
type DynamoDBChangeRecord struct {
	SequenceNumber string                 `json:"sequence_number"`
	Keys           map[string]interface{} `json:"keys"`
	NewImage       map[string]interface{} `json:"new_image,omitempty"`
	OldImage       map[string]interface{} `json:"old_image,omitempty"`
	EventName      string                 `json:"event_name"`
}

// EventType returns "dynamodb:Record" for this event type.
func (e *DynamoDBRecordEvent) EventType() string { return "dynamodb:Record" }

// LogEntry carries a single CloudWatch Logs log event (timestamp + message).
// Used by LambdaLogWriteEvent to transport Lambda execution logs through the
// service bus without importing the store-layer logsstore package.
type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
}

// APIGatewayAccessLogEvent is published by the API Gateway runtime after each
// request when the stage has AccessLogSettings configured. The CloudWatch Logs
// handler writes the formatted log entry to the destination log group/stream
// specified in the access log settings.
type APIGatewayAccessLogEvent struct {
	EventBase
	RestApiId      string `json:"rest_api_id"`
	StageName      string `json:"stage_name"`
	DestinationArn string `json:"destination_arn"`
	LogGroup       string `json:"log_group"`
	LogStream      string `json:"log_stream"`
	FormattedLog   string `json:"formatted_log"`
}

// EventType returns "apigateway:AccessLog" for this event type.
func (e *APIGatewayAccessLogEvent) EventType() string { return "apigateway:AccessLog" }

// CloudWatchLogsPutEvent is published when EventBridge, Scheduler, or Step
// Functions needs to write log events directly to a CloudWatch Logs log group.
// This avoids coupling those services to the CloudWatch Logs store layer.
type CloudWatchLogsPutEvent struct {
	EventBase
	LogGroup  string     `json:"log_group"`
	LogStream string     `json:"log_stream"`
	LogEvents []LogEntry `json:"log_events"`
}

// EventType returns "logs:PutLogEvents" for this event type.
func (e *CloudWatchLogsPutEvent) EventType() string { return "logs:PutLogEvents" }

// LambdaLogWriteEvent is published by the Lambda service after a function
// invocation completes. The CloudWatch Logs handler ingests the log entries
// into the appropriate log group/stream, then applies metric filters and
// subscription filters. This replaces the previous direct logsStore call from
// Lambda, enabling subscription and metric filter evaluation on Lambda-produced
// logs (previously bypassed).
//
// Cycle prevention: the bus increments Depth on each publish. If a subscription
// filter triggers further events (e.g. logs:deliver → Lambda → lambda:LogWrite),
// the depth increases. Events at depth >= MaxEventDepth (default 3) are dropped.
type LambdaLogWriteEvent struct {
	EventBase
	FunctionName string     `json:"function_name"`
	Version      string     `json:"version"`
	LogGroup     string     `json:"log_group"`
	LogStream    string     `json:"log_stream"`
	LogEvents    []LogEntry `json:"log_events"`
}

// EventType returns "lambda:LogWrite" for this event type.
func (e *LambdaLogWriteEvent) EventType() string { return "lambda:LogWrite" }

// StepFunctionsStartExecutionEvent is published when EventBridge, Scheduler,
// or CloudWatch Alarms needs to start a Step Functions execution. The SFN
// service subscribes to this event and creates + executes the state machine
// run asynchronously.
type StepFunctionsStartExecutionEvent struct {
	EventBase
	ExecutionArn    string `json:"execution_arn"`
	StateMachineArn string `json:"state_machine_arn"`
	Input           string `json:"input"`
}

// EventType returns "states:startExecution" for this event type.
func (e *StepFunctionsStartExecutionEvent) EventType() string { return "states:startExecution" }

// EventBridgePutEventsEvent is published when Scheduler or another service
// needs to put events into an EventBridge event bus. The EventBridge service
// subscribes to this event and delivers to matching rules/targets.
type EventBridgePutEventsEvent struct {
	EventBase
	EventBusName string `json:"event_bus_name"`
	Input        string `json:"input"`
}

// EventType returns "events:putEvents" for this event type.
func (e *EventBridgePutEventsEvent) EventType() string { return "events:putEvents" }

// HandlerResult is returned by event handlers. StatusCode and Payload carry the
// response for synchronous dispatch; Error indicates handler failure.
type HandlerResult struct {
	StatusCode int64
	Payload    []byte
	Error      error
}

// EventRegistry maintains a mapping of event type strings to zero-value
// constructors, enabling deserialisation of events from persistent storage.
type EventRegistry struct {
	constructors map[string]func() Event
	mu           sync.RWMutex
}

// NewEventRegistry creates an empty EventRegistry ready for use.
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		constructors: make(map[string]func() Event),
	}
}

// Register associates an event type string with its zero-value constructor.
func (r *EventRegistry) Register(eventType string, constructor func() Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.constructors[eventType] = constructor
}

// New creates a new zero-value instance of the named event type.
func (r *EventRegistry) New(eventType string) (Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.constructors[eventType]
	if !ok {
		return nil, fmt.Errorf("eventbus: unknown event type %q", eventType)
	}
	return c(), nil
}

// Deserialize instantiates an event of the given type and unmarshals the
// provided JSON bytes into it.
func (r *EventRegistry) Deserialize(eventType string, data []byte) (Event, error) {
	event, err := r.New(eventType)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, event); err != nil {
		return nil, fmt.Errorf("eventbus: failed to deserialize event %q: %w", eventType, err)
	}
	return event, nil
}

// SerializeEvent marshals an event to JSON bytes for persistent storage
// or outbox persistence.
func SerializeEvent(event Event) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("eventbus: failed to serialize event %q: %w", event.EventType(), err)
	}
	return data, nil
}

// EventMarshaller provides a generic marshal/unmarshal pair for any Event
// type, satisfying the outbox serialisation contract.
type EventMarshaller[T Event] struct{}

// MarshalEvent serialises the given event to JSON bytes.
func (m *EventMarshaller[T]) MarshalEvent(ctx context.Context, event T) ([]byte, error) {
	return SerializeEvent(event)
}

// UnmarshalEvent deserialises JSON bytes into a concrete event of type T.
func (m *EventMarshaller[T]) UnmarshalEvent(ctx context.Context, data []byte) (T, error) {
	var zero T
	err := json.Unmarshal(data, &zero)
	if err != nil {
		return zero, err
	}
	return zero, nil
}
