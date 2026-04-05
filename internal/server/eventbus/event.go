package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

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

type CallerContext struct {
	ServicePrincipal string `json:"service_principal"`
	SourceARN        string `json:"source_arn"`
	AssumedRoleARN   string `json:"assumed_role_arn"`
	AccountID        string `json:"account_id"`
}

type EventBase struct {
	ID        string        `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Source    string        `json:"source"`
	Region    string        `json:"region"`
	AccountID string        `json:"account_id"`
	Depth     int           `json:"depth"`
	Caller    CallerContext `json:"caller"`
}

func (e *EventBase) EventType() string          { return "" }
func (e *EventBase) EventID() string            { return e.ID }
func (e *EventBase) EventTimestamp() time.Time  { return e.Timestamp }
func (e *EventBase) EventSource() string        { return e.Source }
func (e *EventBase) EventRegion() string        { return e.Region }
func (e *EventBase) EventAccountID() string     { return e.AccountID }
func (e *EventBase) EventDepth() int            { return e.Depth }
func (e *EventBase) SetEventDepth(d int)        { e.Depth = d }
func (e *EventBase) EventCaller() CallerContext { return e.Caller }

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

func (e *ServiceInvokeRequest) EventType() string { return "service:invoke" }

// SNSDeliveryEvent is published when an SNS message needs to be delivered
// to topic subscribers. The handler looks up subscriptions, applies
// protocol-specific envelope formatting, and dispatches via invoker.
type SNSDeliveryEvent struct {
	EventBase
	TopicARN       string `json:"topic_arn"`
	MessageID      string `json:"message_id"`
	Message        string `json:"message"`
	Subject        string `json:"subject,omitempty"`
	MessageGroupId string `json:"message_group_id,omitempty"`
}

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

func (e *EventBridgeDeliveryEvent) EventType() string { return "events:deliver" }

// ScheduleFiredEvent is published when the Scheduler engine fires a schedule.
type ScheduleFiredEvent struct {
	EventBase
	ScheduleName string `json:"schedule_name"`
	ScheduleArn  string `json:"schedule_arn"`
	TargetArn    string `json:"target_arn"`
	Input        string `json:"input"`
}

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

func (e *StepFunctionsStartExecutionEvent) EventType() string { return "states:startExecution" }

// EventBridgePutEventsEvent is published when Scheduler or another service
// needs to put events into an EventBridge event bus. The EventBridge service
// subscribes to this event and delivers to matching rules/targets.
type EventBridgePutEventsEvent struct {
	EventBase
	EventBusName string `json:"event_bus_name"`
	Input        string `json:"input"`
}

func (e *EventBridgePutEventsEvent) EventType() string { return "events:putEvents" }

// HandlerResult is returned by event handlers. StatusCode and Payload carry the
// response for synchronous dispatch; Error indicates handler failure.
type HandlerResult struct {
	StatusCode int64
	Payload    []byte
	Error      error
}

type EventRegistry struct {
	constructors map[string]func() Event
	mu           sync.RWMutex
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		constructors: make(map[string]func() Event),
	}
}

func (r *EventRegistry) Register(eventType string, constructor func() Event) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.constructors[eventType] = constructor
}

func (r *EventRegistry) New(eventType string) (Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.constructors[eventType]
	if !ok {
		return nil, fmt.Errorf("eventbus: unknown event type %q", eventType)
	}
	return c(), nil
}

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

func SerializeEvent(event Event) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("eventbus: failed to serialize event %q: %w", event.EventType(), err)
	}
	return data, nil
}

type EventMarshaller[T Event] struct{}

func (m *EventMarshaller[T]) MarshalEvent(ctx context.Context, event T) ([]byte, error) {
	return SerializeEvent(event)
}

func (m *EventMarshaller[T]) UnmarshalEvent(ctx context.Context, data []byte) (T, error) {
	var zero T
	err := json.Unmarshal(data, &zero)
	if err != nil {
		return zero, err
	}
	return zero, nil
}
