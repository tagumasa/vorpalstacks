package sfn

import "time"

// This file holds the execution history event record and its detail
// shapes — one struct per modelled EventDetails member.

// ExecutionHistoryEvent represents a single event in the execution history of a
// state machine execution. The detail members and their payload fields follow
// the HistoryEvent shape of the Smithy model exactly: every event type carries
// at most one detail member, and state-transition events share the generic
// stateEnteredEventDetails / stateExitedEventDetails pair.
type ExecutionHistoryEvent struct {
	ExecutionArn                             string                                    `json:"executionArn"`
	EventId                                  int64                                     `json:"eventId"`
	PreviousEventId                          int64                                     `json:"previousEventId"`
	Timestamp                                time.Time                                 `json:"timestamp"`
	Type                                     string                                    `json:"type"`
	ExecutionStartedEventDetails             *ExecutionStartedEventDetails             `json:"executionStartedEventDetails,omitempty"`
	ExecutionSucceededEventDetails           *ExecutionSucceededEventDetails           `json:"executionSucceededEventDetails,omitempty"`
	ExecutionFailedEventDetails              *ExecutionFailedEventDetails              `json:"executionFailedEventDetails,omitempty"`
	ExecutionAbortedEventDetails             *ExecutionAbortedEventDetails             `json:"executionAbortedEventDetails,omitempty"`
	ExecutionTimedOutEventDetails            *ExecutionTimedOutEventDetails            `json:"executionTimedOutEventDetails,omitempty"`
	ExecutionRedrivenEventDetails            *ExecutionRedrivenEventDetails            `json:"executionRedrivenEventDetails,omitempty"`
	TaskScheduledEventDetails                *TaskScheduledEventDetails                `json:"taskScheduledEventDetails,omitempty"`
	TaskStartedEventDetails                  *TaskStartedEventDetails                  `json:"taskStartedEventDetails,omitempty"`
	TaskStartFailedEventDetails              *TaskStartFailedEventDetails              `json:"taskStartFailedEventDetails,omitempty"`
	TaskSubmittedEventDetails                *TaskSubmittedEventDetails                `json:"taskSubmittedEventDetails,omitempty"`
	TaskSubmitFailedEventDetails             *TaskSubmitFailedEventDetails             `json:"taskSubmitFailedEventDetails,omitempty"`
	TaskSucceededEventDetails                *TaskSucceededEventDetails                `json:"taskSucceededEventDetails,omitempty"`
	TaskFailedEventDetails                   *TaskFailedEventDetails                   `json:"taskFailedEventDetails,omitempty"`
	TaskTimedOutEventDetails                 *TaskTimedOutEventDetails                 `json:"taskTimedOutEventDetails,omitempty"`
	LambdaFunctionScheduledEventDetails      *LambdaFunctionScheduledEventDetails      `json:"lambdaFunctionScheduledEventDetails,omitempty"`
	LambdaFunctionScheduleFailedEventDetails *LambdaFunctionScheduleFailedEventDetails `json:"lambdaFunctionScheduleFailedEventDetails,omitempty"`
	LambdaFunctionStartFailedEventDetails    *LambdaFunctionStartFailedEventDetails    `json:"lambdaFunctionStartFailedEventDetails,omitempty"`
	LambdaFunctionFailedEventDetails         *LambdaFunctionFailedEventDetails         `json:"lambdaFunctionFailedEventDetails,omitempty"`
	LambdaFunctionTimedOutEventDetails       *LambdaFunctionTimedOutEventDetails       `json:"lambdaFunctionTimedOutEventDetails,omitempty"`
	LambdaFunctionSucceededEventDetails      *LambdaFunctionSucceededEventDetails      `json:"lambdaFunctionSucceededEventDetails,omitempty"`
	ActivityTaskScheduledEventDetails        *ActivityTaskScheduledEventDetails        `json:"activityScheduledEventDetails,omitempty"`
	ActivityScheduleFailedEventDetails       *ActivityScheduleFailedEventDetails       `json:"activityScheduleFailedEventDetails,omitempty"`
	ActivityTaskStartedEventDetails          *ActivityTaskStartedEventDetails          `json:"activityStartedEventDetails,omitempty"`
	ActivityTaskSucceededEventDetails        *ActivityTaskSucceededEventDetails        `json:"activitySucceededEventDetails,omitempty"`
	ActivityTaskFailedEventDetails           *ActivityTaskFailedEventDetails           `json:"activityFailedEventDetails,omitempty"`
	ActivityTaskTimedOutEventDetails         *ActivityTaskTimedOutEventDetails         `json:"activityTimedOutEventDetails,omitempty"`
	StateEnteredEventDetails                 *StateEnteredEventDetails                 `json:"stateEnteredEventDetails,omitempty"`
	StateExitedEventDetails                  *StateExitedEventDetails                  `json:"stateExitedEventDetails,omitempty"`
	MapStateStartedEventDetails              *MapStateStartedEventDetails              `json:"mapStateStartedEventDetails,omitempty"`
	MapRunStartedEventDetails                *MapRunStartedEventDetails                `json:"mapRunStartedEventDetails,omitempty"`
	MapRunFailedEventDetails                 *MapRunFailedEventDetails                 `json:"mapRunFailedEventDetails,omitempty"`
	MapRunRedrivenEventDetails               *MapRunRedrivenEventDetails               `json:"mapRunRedrivenEventDetails,omitempty"`
	MapIterationEventDetails                 *MapIterationEventDetails                 `json:"mapIterationEventDetails,omitempty"`
	EvaluationFailedEventDetails             *EvaluationFailedEventDetails             `json:"evaluationFailedEventDetails,omitempty"`
}

// ExecutionFailedEventDetails contains error and cause information for a failed execution.
type ExecutionFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// EvaluationFailedEventDetails contains details about a failed Choice state
// condition evaluation.
type EvaluationFailedEventDetails struct {
	State    string `json:"state"`
	Cause    string `json:"cause"`
	Error    string `json:"error"`
	Location string `json:"location"`
}

// TaskStartedEventDetails contains the resource members emitted when a task
// integration starts (model members: resource, resourceType).
type TaskStartedEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
}

// TaskScheduledEventDetails contains the members emitted when a task
// integration is scheduled (model members: resource, resourceType, region,
// parameters, timeoutInSeconds, heartbeatInSeconds).
type TaskScheduledEventDetails struct {
	Resource           string                    `json:"resource"`
	ResourceType       string                    `json:"resourceType"`
	Region             string                    `json:"region,omitempty"`
	Parameters         interface{}               `json:"parameters,omitempty"`
	TimeoutInSeconds   int32                     `json:"timeoutInSeconds,omitempty"`
	HeartbeatInSeconds int32                     `json:"heartbeatInSeconds,omitempty"`
	TaskCredentials    *TaskScheduledCredentials `json:"taskCredentials,omitempty"`
}

// TaskScheduledCredentials is the taskCredentials member of the scheduled
// event details: the role ARN the Credentials payload template resolved
// to (TaskCredentials shape, sole member roleArn).
type TaskScheduledCredentials struct {
	RoleArn string `json:"roleArn"`
}

// TaskStartFailedEventDetails contains the members emitted when a task
// integration failed to start (model members: resource, resourceType, error,
// cause).
type TaskStartFailedEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
	Error        string `json:"error"`
	Cause        string `json:"cause"`
}

// TaskSubmittedEventDetails contains the members emitted when a callback
// task's payload is accepted by the integrated service (model members:
// resource, resourceType, output, outputDetails).
type TaskSubmittedEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
	Output       string `json:"output,omitempty"`
}

// TaskSubmitFailedEventDetails contains the members emitted when a callback
// task's submission to the integrated service failed (model members:
// resource, resourceType, error, cause).
type TaskSubmitFailedEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
	Error        string `json:"error"`
	Cause        string `json:"cause"`
}

// TaskSucceededEventDetails contains the output produced by a successful task
// integration (model members: output, outputDetails, resource, resourceType).
type TaskSucceededEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
	Output       string `json:"output"`
}

// TaskFailedEventDetails contains error and cause information for a failed
// task integration (model members: cause, error, resource, resourceType).
type TaskFailedEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
	Error        string `json:"error"`
	Cause        string `json:"cause"`
}

// TaskTimedOutEventDetails contains the members emitted when a task
// integration times out (model members: cause, error, resource, resourceType).
type TaskTimedOutEventDetails struct {
	Resource     string `json:"resource"`
	ResourceType string `json:"resourceType"`
	Error        string `json:"error"`
	Cause        string `json:"cause"`
}

// LambdaFunctionScheduledEventDetails contains the members emitted when a
// Lambda task is scheduled (model members: input, inputDetails, resource,
// taskCredentials, timeoutInSeconds).
type LambdaFunctionScheduledEventDetails struct {
	Resource         string                    `json:"resource"`
	Input            string                    `json:"input"`
	TimeoutInSeconds int32                     `json:"timeoutInSeconds,omitempty"`
	TaskCredentials  *TaskScheduledCredentials `json:"taskCredentials,omitempty"`
}

// LambdaFunctionScheduleFailedEventDetails contains the error pair emitted
// when scheduling a Lambda task fails.
type LambdaFunctionScheduleFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// LambdaFunctionStartFailedEventDetails contains the error pair emitted when
// starting a Lambda invocation fails at the service level.
type LambdaFunctionStartFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// LambdaFunctionFailedEventDetails contains the error pair reported by a
// Lambda function that failed.
type LambdaFunctionFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// LambdaFunctionTimedOutEventDetails contains the error pair emitted when a
// Lambda task times out.
type LambdaFunctionTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// LambdaFunctionSucceededEventDetails contains the output of a successful
// Lambda invocation (model members: output, outputDetails).
type LambdaFunctionSucceededEventDetails struct {
	Output string `json:"output"`
}

// ExecutionStartedEventDetails contains details emitted when an execution
// starts (model members: input, inputDetails, roleArn,
// stateMachineAliasArn, stateMachineVersionArn).
type ExecutionStartedEventDetails struct {
	Input                  string `json:"input"`
	RoleArn                string `json:"roleArn"`
	StateMachineAliasArn   string `json:"stateMachineAliasArn,omitempty"`
	StateMachineVersionArn string `json:"stateMachineVersionArn,omitempty"`
}

// ExecutionSucceededEventDetails contains the output produced by a successful execution.
type ExecutionSucceededEventDetails struct {
	Output string `json:"output"`
}

// ExecutionAbortedEventDetails contains error and cause information for an aborted
// execution.
type ExecutionAbortedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ExecutionTimedOutEventDetails contains error and cause information for a timed-out
// execution.
type ExecutionTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ExecutionRedrivenEventDetails contains details emitted when an execution is
// redriven (model member: redriveCount). The execution keeps its original ARN;
// this event marks the point in history where the redrive occurred.
type ExecutionRedrivenEventDetails struct {
	RedriveCount int64 `json:"redriveCount"`
}

// StateEnteredEventDetails is the generic detail shape every state-entry
// event carries, whatever the state type (model members: input, inputDetails,
// name).
type StateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// StateExitedEventDetails is the generic detail shape every state-exit event
// carries (model members: assignedVariables, assignedVariablesDetails, name,
// output, outputDetails).
type StateExitedEventDetails struct {
	Name              string                 `json:"name"`
	Output            string                 `json:"output"`
	AssignedVariables map[string]interface{} `json:"assignedVariables,omitempty"`
	// NextState records the transition a Choice state took. The HistoryEvent
	// model defines no such wire member — it exists so redrive and restart
	// recovery can resume at the state the choice actually selected; the
	// response serialiser never emits it.
	NextState string `json:"nextState,omitempty"`
}

// MapStateStartedEventDetails contains the item count of a Map state when it
// starts (model member: length).
type MapStateStartedEventDetails struct {
	Length int64 `json:"length"`
}

// MapRunStartedEventDetails contains the ARN of a Map Run when it starts
// (model member: mapRunArn).
type MapRunStartedEventDetails struct {
	MapRunArn string `json:"mapRunArn"`
}

// MapRunFailedEventDetails contains the error pair of a failed Map Run.
type MapRunFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// MapRunRedrivenEventDetails contains the members emitted when a Map Run is
// redriven (model members: mapRunArn, redriveCount).
type MapRunRedrivenEventDetails struct {
	MapRunArn    string `json:"mapRunArn"`
	RedriveCount int64  `json:"redriveCount"`
}

// MapIterationEventDetails contains the index and name of one Map state
// iteration, shared by the MapIterationStarted/Succeeded/Failed/Aborted
// events (model members: index, name).
type MapIterationEventDetails struct {
	Index int64  `json:"index"`
	Name  string `json:"name"`
}

// ActivityTaskScheduledEventDetails contains details emitted when an activity task is
// scheduled.
type ActivityTaskScheduledEventDetails struct {
	Resource         string `json:"resource"`
	Input            string `json:"input"`
	TimeoutInSeconds int32  `json:"timeoutInSeconds,omitempty"`
	HeartbeatSeconds int32  `json:"heartbeatInSeconds,omitempty"`
}

// ActivityTaskStartedEventDetails contains details emitted when an activity task is
// started by a worker.
type ActivityTaskStartedEventDetails struct {
	WorkerName string `json:"workerName"`
}

// ActivityTaskSucceededEventDetails contains the output of a successfully completed
// activity task.
type ActivityTaskSucceededEventDetails struct {
	Output string `json:"output"`
}

// ActivityTaskFailedEventDetails contains error and cause information for a failed
// activity task.
type ActivityTaskFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ActivityTaskTimedOutEventDetails contains error and cause information for a timed-out
// activity task.
type ActivityTaskTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ActivityScheduleFailedEventDetails contains the error pair emitted when
// scheduling an activity task fails.
type ActivityScheduleFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}
