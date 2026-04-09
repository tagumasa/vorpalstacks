package sfn

import (
	"encoding/json"
	"time"
)

// StateMachine represents an AWS Step Functions state machine definition and metadata.
type StateMachine struct {
	StateMachineArn      string                `json:"stateMachineArn"`
	Name                 string                `json:"name"`
	Definition           string                `json:"definition"`
	RoleArn              string                `json:"roleArn"`
	Type                 string                `json:"type"`
	CreationDate         time.Time             `json:"creationDate"`
	UpdateDate           time.Time             `json:"updateDate"`
	Description          string                `json:"description"`
	Status               string                `json:"status"`
	Tags                 map[string]string     `json:"tags"`
	VariableReferences   map[string][]string   `json:"-"`
	LoggingConfiguration *LoggingConfiguration `json:"loggingConfiguration,omitempty"`
}

// LoggingConfiguration controls whether execution history is logged to
// CloudWatch Logs and where.
type LoggingConfiguration struct {
	Level                string           `json:"level,omitempty"`
	IncludeExecutionData bool             `json:"includeExecutionData,omitempty"`
	Destinations         []LogDestination `json:"destinations,omitempty"`
}

// LogDestination specifies a CloudWatch Logs log group for execution history.
type LogDestination struct {
	CloudWatchLogsLogGroup *CloudWatchLogsLogGroup `json:"cloudWatchLogsLogGroup,omitempty"`
}

// CloudWatchLogsLogGroup identifies the log group for SFN execution history.
type CloudWatchLogsLogGroup struct {
	LogGroupArn string `json:"logGroupArn,omitempty"`
}

// Execution represents a single execution of a Step Functions state machine.
type Execution struct {
	ExecutionArn    string                  `json:"executionArn"`
	StateMachineArn string                  `json:"stateMachineArn"`
	Name            string                  `json:"name"`
	Status          string                  `json:"status"`
	Input           string                  `json:"input"`
	Output          string                  `json:"output"`
	TraceHeader     string                  `json:"traceHeader"`
	InputDetails    *ExecutionInputDetails  `json:"inputDetails,omitempty"`
	OutputDetails   *ExecutionOutputDetails `json:"outputDetails,omitempty"`
	StartDate       time.Time               `json:"startDate"`
	StopDate        time.Time               `json:"stopDate,omitempty"`
	Error           string                  `json:"error,omitempty"`
	Cause           string                  `json:"cause,omitempty"`
}

// ExecutionInputDetails describes the input included in an execution.
type ExecutionInputDetails struct {
	Included bool   `json:"included"`
	Type     string `json:"type"`
}

// ExecutionOutputDetails describes the output produced by an execution.
type ExecutionOutputDetails struct {
	Included bool   `json:"included"`
	Type     string `json:"type"`
}

// ExecutionHistoryEvent represents a single event in the execution history of a
// state machine execution.
type ExecutionHistoryEvent struct {
	ExecutionArn                      string                             `json:"executionArn"`
	EventId                           int64                              `json:"eventId"`
	PreviousEventId                   int64                              `json:"previousEventId"`
	Timestamp                         time.Time                          `json:"timestamp"`
	Type                              string                             `json:"type"`
	ExecutionFailedEventDetails       *ExecutionFailedEventDetails       `json:"executionFailedEventDetails,omitempty"`
	TaskStartedEventDetails           *TaskStartedEventDetails           `json:"taskStartedEventDetails,omitempty"`
	TaskSucceededEventDetails         *TaskSucceededEventDetails         `json:"taskSucceededEventDetails,omitempty"`
	TaskFailedEventDetails            *TaskFailedEventDetails            `json:"taskFailedEventDetails,omitempty"`
	ExecutionStartedEventDetails      *ExecutionStartedEventDetails      `json:"executionStartedEventDetails,omitempty"`
	ExecutionSucceededEventDetails    *ExecutionSucceededEventDetails    `json:"executionSucceededEventDetails,omitempty"`
	ExecutionAbortedEventDetails      *ExecutionAbortedEventDetails      `json:"executionAbortedEventDetails,omitempty"`
	ExecutionTimedOutEventDetails     *ExecutionTimedOutEventDetails     `json:"executionTimedOutEventDetails,omitempty"`
	ChoiceStateEnteredEventDetails    *ChoiceStateEnteredEventDetails    `json:"choiceStateEnteredEventDetails,omitempty"`
	ChoiceStateExitedEventDetails     *ChoiceStateExitedEventDetails     `json:"choiceStateExitedEventDetails,omitempty"`
	WaitStateEnteredEventDetails      *WaitStateEnteredEventDetails      `json:"waitStateEnteredEventDetails,omitempty"`
	WaitStateExitedEventDetails       *WaitStateExitedEventDetails       `json:"waitStateExitedEventDetails,omitempty"`
	ParallelStateEnteredEventDetails  *ParallelStateEnteredEventDetails  `json:"parallelStateEnteredEventDetails,omitempty"`
	ParallelStateExitedEventDetails   *ParallelStateExitedEventDetails   `json:"parallelStateExitedEventDetails,omitempty"`
	MapStateEnteredEventDetails       *MapStateEnteredEventDetails       `json:"mapStateEnteredEventDetails,omitempty"`
	MapStateExitedEventDetails        *MapStateExitedEventDetails        `json:"mapStateExitedEventDetails,omitempty"`
	PassStateEnteredEventDetails      *PassStateEnteredEventDetails      `json:"passStateEnteredEventDetails,omitempty"`
	PassStateExitedEventDetails       *PassStateExitedEventDetails       `json:"passStateExitedEventDetails,omitempty"`
	FailStateEnteredEventDetails      *FailStateEnteredEventDetails      `json:"failStateEnteredEventDetails,omitempty"`
	SucceedStateEnteredEventDetails   *SucceedStateEnteredEventDetails   `json:"succeedStateEnteredEventDetails,omitempty"`
	ActivityTaskScheduledEventDetails *ActivityTaskScheduledEventDetails `json:"activityTaskScheduledEventDetails,omitempty"`
	ActivityTaskStartedEventDetails   *ActivityTaskStartedEventDetails   `json:"activityTaskStartedEventDetails,omitempty"`
	ActivityTaskSucceededEventDetails *ActivityTaskSucceededEventDetails `json:"activityTaskSucceededEventDetails,omitempty"`
	ActivityTaskFailedEventDetails    *ActivityTaskFailedEventDetails    `json:"activityTaskFailedEventDetails,omitempty"`
	ActivityTaskTimedOutEventDetails  *ActivityTaskTimedOutEventDetails  `json:"activityTaskTimedOutEventDetails,omitempty"`
	EvaluationFailedEventDetails      *EvaluationFailedEventDetails      `json:"evaluationFailedEventDetails,omitempty"`
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

// TaskStartedEventDetails contains details emitted when a Task state starts execution.
type TaskStartedEventDetails struct {
	Resource string `json:"resource"`
	Input    string `json:"input"`
}

// TaskSucceededEventDetails contains the output produced by a successful Task state.
type TaskSucceededEventDetails struct {
	Resource string `json:"resource"`
	Output   string `json:"output"`
}

// TaskFailedEventDetails contains error and cause information for a failed Task state.
type TaskFailedEventDetails struct {
	Resource string `json:"resource"`
	Error    string `json:"error"`
	Cause    string `json:"cause"`
}

// ExecutionStartedEventDetails contains details emitted when an execution starts.
type ExecutionStartedEventDetails struct {
	Input           string `json:"input"`
	RoleArn         string `json:"roleArn"`
	StateMachineArn string `json:"stateMachineArn"`
	Name            string `json:"name"`
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

// ChoiceStateEnteredEventDetails contains details emitted when a Choice state is
// entered.
type ChoiceStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// ChoiceStateExitedEventDetails contains details emitted when a Choice state is exited.
type ChoiceStateExitedEventDetails struct {
	Output    string `json:"output"`
	Name      string `json:"name"`
	NextState string `json:"nextState"`
}

// WaitStateEnteredEventDetails contains details emitted when a Wait state is entered.
type WaitStateEnteredEventDetails struct {
	Input     string `json:"input"`
	Name      string `json:"name"`
	Seconds   int64  `json:"seconds"`
	Timestamp string `json:"timestamp"`
}

// WaitStateExitedEventDetails contains details emitted when a Wait state is exited.
type WaitStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

// ParallelStateEnteredEventDetails contains details emitted when a Parallel state
// is entered.
type ParallelStateEnteredEventDetails struct {
	Input    string   `json:"input"`
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

// ParallelStateExitedEventDetails contains details emitted when a Parallel state is
// exited.
type ParallelStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

// MapStateEnteredEventDetails contains details emitted when a Map state is entered.
type MapStateEnteredEventDetails struct {
	Input          string `json:"input"`
	Name           string `json:"name"`
	ItemsProcessed int64  `json:"itemsProcessed"`
	ItemsFailed    int64  `json:"itemsFailed"`
}

// MapStateExitedEventDetails contains details emitted when a Map state is exited.
type MapStateExitedEventDetails struct {
	Output         string `json:"output"`
	Name           string `json:"name"`
	ItemsProcessed int64  `json:"itemsProcessed,omitempty"`
	ItemsFailed    int64  `json:"itemsFailed,omitempty"`
}

// PassStateEnteredEventDetails contains details emitted when a Pass state is entered.
type PassStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// PassStateExitedEventDetails contains details emitted when a Pass state is exited.
type PassStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

// FailStateEnteredEventDetails contains details emitted when a Fail state is entered.
type FailStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// SucceedStateEnteredEventDetails contains details emitted when a Succeed state is
// entered.
type SucceedStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// Activity represents an AWS Step Functions activity used by Task states.
type Activity struct {
	ActivityArn  string    `json:"activityArn"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

// StateMachineDefinition represents the Amazon States Language definition of a state
// machine.
type StateMachineDefinition struct {
	StartAt        string                 `json:"StartAt"`
	States         map[string]interface{} `json:"States"`
	TimeoutSeconds int32                  `json:"TimeoutSeconds,omitempty"`
	Version        string                 `json:"Version,omitempty"`
	QueryLanguage  string                 `json:"QueryLanguage,omitempty"`
}

// State defines the common interface for all state machine state types.
type State interface {
	GetType() string
	GetNext() string
	GetEnd() bool
	GetQueryLanguage() string
}

func getInputPathFromInputOutput(io *InputOutput) string {
	if io == nil {
		return ""
	}
	return io.InputPath
}

func getOutputPathFromInputOutput(io *InputOutput) string {
	if io == nil {
		return ""
	}
	return io.Path
}

// PassState represents a Pass state that simply passes its input to its output,
// optionally transforming it.
type PassState struct {
	Name           string                 `json:"-"`
	Type           string                 `json:"Type"`
	Next           string                 `json:"Next,omitempty"`
	End            bool                   `json:"End,omitempty"`
	Comment        string                 `json:"Comment,omitempty"`
	Input          *InputOutput           `json:"Input,omitempty"`
	Output         *InputOutput           `json:"Output,omitempty"`
	OutputRaw      json.RawMessage        `json:"-"`
	Result         interface{}            `json:"Result,omitempty"`
	ResultPath     string                 `json:"ResultPath,omitempty"`
	ResultSelector *ResultSelector        `json:"ResultSelector,omitempty"`
	Parameters     *Parameters            `json:"Parameters,omitempty"`
	QueryLanguage  string                 `json:"QueryLanguage,omitempty"`
	Assign         map[string]interface{} `json:"Assign,omitempty"`
	JSONataOutput  interface{}            `json:"-"`
}

// GetType returns the state type identifier (e.g. "Pass").
func (s *PassState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to.
func (s *PassState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *PassState) GetEnd() bool { return s.End }

// GetInputPath returns the input path filter applied to the state input.
func (s *PassState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *PassState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetResultSelector returns the result selector applied to the state output.
func (s *PassState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// GetQueryLanguage returns the query language used by this state.
func (s *PassState) GetQueryLanguage() string { return s.QueryLanguage }

// TaskState represents a Task state that executes a unit of work, such as calling an
// activity or a Lambda function.
type TaskState struct {
	Name             string                 `json:"-"`
	Type             string                 `json:"Type"`
	Next             string                 `json:"Next,omitempty"`
	End              bool                   `json:"End,omitempty"`
	Comment          string                 `json:"Comment,omitempty"`
	Input            *InputOutput           `json:"Input,omitempty"`
	Output           *InputOutput           `json:"Output,omitempty"`
	OutputRaw        json.RawMessage        `json:"-"`
	Resource         string                 `json:"Resource"`
	Parameters       *Parameters            `json:"Parameters,omitempty"`
	ResultPath       string                 `json:"ResultPath,omitempty"`
	ResultSelector   *ResultSelector        `json:"ResultSelector,omitempty"`
	Retry            []*RetryPolicy         `json:"Retry,omitempty"`
	Catch            []*CatchPolicy         `json:"Catch,omitempty"`
	TimeoutSeconds   interface{}            `json:"TimeoutSeconds,omitempty"`
	HeartbeatSeconds interface{}            `json:"HeartbeatSeconds,omitempty"`
	QueryLanguage    string                 `json:"QueryLanguage,omitempty"`
	Arguments        interface{}            `json:"Arguments,omitempty"`
	Assign           map[string]interface{} `json:"Assign,omitempty"`
	JSONataOutput    interface{}            `json:"-"`
}

// GetType returns the state type identifier (e.g. "Task").
func (s *TaskState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to.
func (s *TaskState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *TaskState) GetEnd() bool { return s.End }

// GetInputPath returns the input path filter applied to the state input.
func (s *TaskState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *TaskState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetResultSelector returns the result selector applied to the state output.
func (s *TaskState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// GetQueryLanguage returns the query language used by this state.
func (s *TaskState) GetQueryLanguage() string { return s.QueryLanguage }

// GetTimeoutSeconds returns the timeout in seconds for the task execution.
func (s *TaskState) GetTimeoutSeconds() int32 {
	if s.TimeoutSeconds == nil {
		return 0
	}
	switch v := s.TimeoutSeconds.(type) {
	case float64:
		return int32(v)
	case int32:
		return v
	case json.Number:
		n, _ := v.Int64()
		return int32(n)
	}
	return 0
}

// GetHeartbeatSeconds returns the heartbeat interval in seconds for the task.
func (s *TaskState) GetHeartbeatSeconds() int32 {
	if s.HeartbeatSeconds == nil {
		return 0
	}
	switch v := s.HeartbeatSeconds.(type) {
	case float64:
		return int32(v)
	case int32:
		return v
	case json.Number:
		n, _ := v.Int64()
		return int32(n)
	}
	return 0
}

// ChoiceState represents a Choice state that adds branching logic to a state machine.
type ChoiceState struct {
	Name          string        `json:"-"`
	Type          string        `json:"Type"`
	Comment       string        `json:"Comment,omitempty"`
	Input         *InputOutput  `json:"Input,omitempty"`
	Output        *InputOutput  `json:"Output,omitempty"`
	Choices       []*ChoiceRule `json:"Choices"`
	Default       string        `json:"Default,omitempty"`
	QueryLanguage string        `json:"QueryLanguage,omitempty"`
}

// GetType returns the state type identifier (e.g. "Choice").
func (s *ChoiceState) GetType() string { return s.Type }

// GetNext always returns empty for Choice states; transitions are determined by rules.
func (s *ChoiceState) GetNext() string { return "" }

// GetEnd always returns false for Choice states.
func (s *ChoiceState) GetEnd() bool { return false }

// GetInputPath returns the input path filter applied to the state input.
func (s *ChoiceState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *ChoiceState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetQueryLanguage returns the query language used by this state.
func (s *ChoiceState) GetQueryLanguage() string { return s.QueryLanguage }

// WaitState represents a Wait state that delays the state machine from continuing for
// a specified time.
type WaitState struct {
	Name          string                 `json:"-"`
	Type          string                 `json:"Type"`
	Next          string                 `json:"Next,omitempty"`
	End           bool                   `json:"End,omitempty"`
	Comment       string                 `json:"Comment,omitempty"`
	Input         *InputOutput           `json:"Input,omitempty"`
	Output        *InputOutput           `json:"Output,omitempty"`
	Seconds       interface{}            `json:"Seconds,omitempty"`
	TimestampPath string                 `json:"TimestampPath,omitempty"`
	SecondsPath   string                 `json:"SecondsPath,omitempty"`
	Timestamp     string                 `json:"Timestamp,omitempty"`
	QueryLanguage string                 `json:"QueryLanguage,omitempty"`
	Assign        map[string]interface{} `json:"Assign,omitempty"`
}

// GetType returns the state type identifier (e.g. "Wait").
func (s *WaitState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to after waiting.
func (s *WaitState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *WaitState) GetEnd() bool { return s.End }

// GetInputPath returns the input path filter applied to the state input.
func (s *WaitState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *WaitState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetQueryLanguage returns the query language used by this state.
func (s *WaitState) GetQueryLanguage() string { return s.QueryLanguage }

// GetSeconds returns the number of seconds to wait before transitioning.
func (s *WaitState) GetSeconds() int32 {
	if s.Seconds == nil {
		return 0
	}
	switch v := s.Seconds.(type) {
	case float64:
		return int32(v)
	case int32:
		return v
	case json.Number:
		n, _ := v.Int64()
		return int32(n)
	}
	return 0
}

// ParallelState represents a Parallel state that executes multiple branches
// concurrently.
type ParallelState struct {
	Name          string                    `json:"-"`
	Type          string                    `json:"Type"`
	Next          string                    `json:"Next,omitempty"`
	End           bool                      `json:"End,omitempty"`
	Comment       string                    `json:"Comment,omitempty"`
	Input         *InputOutput              `json:"Input,omitempty"`
	Output        *InputOutput              `json:"Output,omitempty"`
	OutputRaw     json.RawMessage           `json:"-"`
	Branches      []*StateMachineDefinition `json:"Branches"`
	Retry         []*RetryPolicy            `json:"Retry,omitempty"`
	Catch         []*CatchPolicy            `json:"Catch,omitempty"`
	QueryLanguage string                    `json:"QueryLanguage,omitempty"`
	Arguments     interface{}               `json:"Arguments,omitempty"`
	Assign        map[string]interface{}    `json:"Assign,omitempty"`
	JSONataOutput interface{}               `json:"-"`
}

// GetType returns the state type identifier (e.g. "Parallel").
func (s *ParallelState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to after all branches complete.
func (s *ParallelState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *ParallelState) GetEnd() bool { return s.End }

// GetInputPath returns the input path filter applied to the state input.
func (s *ParallelState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *ParallelState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetQueryLanguage returns the query language used by this state.
func (s *ParallelState) GetQueryLanguage() string { return s.QueryLanguage }

// MapState represents a Map state that processes a collection of items by applying a
// sub-state machine to each.
type MapState struct {
	Name           string                  `json:"-"`
	Type           string                  `json:"Type"`
	Next           string                  `json:"Next,omitempty"`
	End            bool                    `json:"End,omitempty"`
	Comment        string                  `json:"Comment,omitempty"`
	Input          *InputOutput            `json:"Input,omitempty"`
	Output         *InputOutput            `json:"Output,omitempty"`
	OutputRaw      json.RawMessage         `json:"-"`
	Iterator       *StateMachineDefinition `json:"Iterator,omitempty"`
	ItemProcessor  *StateMachineDefinition `json:"ItemProcessor,omitempty"`
	ItemsPath      string                  `json:"ItemsPath,omitempty"`
	MaxConcurrency int32                   `json:"MaxConcurrency,omitempty"`
	Retry          []*RetryPolicy          `json:"Retry,omitempty"`
	Catch          []*CatchPolicy          `json:"Catch,omitempty"`
	QueryLanguage  string                  `json:"QueryLanguage,omitempty"`
	Items          interface{}             `json:"Items,omitempty"`
	ItemSelector   interface{}             `json:"ItemSelector,omitempty"`
	Assign         map[string]interface{}  `json:"Assign,omitempty"`
	JSONataOutput  interface{}             `json:"-"`
}

// GetType returns the state type identifier (e.g. "Map").
func (s *MapState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to after map processing completes.
func (s *MapState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *MapState) GetEnd() bool { return s.End }

// GetInputPath returns the input path filter applied to the state input.
func (s *MapState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *MapState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetQueryLanguage returns the query language used by this state.
func (s *MapState) GetQueryLanguage() string { return s.QueryLanguage }

// GetIterator returns the sub-state machine definition used to process each item,
// preferring the Iterator field and falling back to ItemProcessor.
func (s *MapState) GetIterator() *StateMachineDefinition {
	if s.Iterator != nil {
		return s.Iterator
	}
	return s.ItemProcessor
}

// FailState represents a Fail state that stops the execution and marks it as failed.
type FailState struct {
	Name          string      `json:"-"`
	Type          string      `json:"Type"`
	Comment       string      `json:"Comment,omitempty"`
	Cause         interface{} `json:"Cause,omitempty"`
	Error         string      `json:"Error,omitempty"`
	QueryLanguage string      `json:"QueryLanguage,omitempty"`
}

// GetType returns the state type identifier (e.g. "Fail").
func (s *FailState) GetType() string { return s.Type }

// GetNext always returns empty for Fail states as they are terminal.
func (s *FailState) GetNext() string { return "" }

// GetEnd always returns true for Fail states.
func (s *FailState) GetEnd() bool { return true }

// GetQueryLanguage returns the query language used by this state.
func (s *FailState) GetQueryLanguage() string { return s.QueryLanguage }

// GetCause returns the failure cause as a string.
func (s *FailState) GetCause() string {
	if s.Cause == nil {
		return ""
	}
	switch v := s.Cause.(type) {
	case string:
		return v
	default:
		return ""
	}
}

// SucceedState represents a Succeed state that stops the execution successfully.
type SucceedState struct {
	Name          string          `json:"-"`
	Type          string          `json:"Type"`
	Comment       string          `json:"Comment,omitempty"`
	Input         *InputOutput    `json:"Input,omitempty"`
	Output        *InputOutput    `json:"Output,omitempty"`
	OutputRaw     json.RawMessage `json:"-"`
	QueryLanguage string          `json:"QueryLanguage,omitempty"`
	JSONataOutput interface{}     `json:"-"`
}

// GetType returns the state type identifier (e.g. "Succeed").
func (s *SucceedState) GetType() string { return s.Type }

// GetNext always returns empty for Succeed states as they are terminal.
func (s *SucceedState) GetNext() string { return "" }

// GetEnd always returns true for Succeed states.
func (s *SucceedState) GetEnd() bool { return true }

// GetInputPath returns the input path filter applied to the state input.
func (s *SucceedState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path filter applied to the state output.
func (s *SucceedState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetQueryLanguage returns the query language used by this state.
func (s *SucceedState) GetQueryLanguage() string { return s.QueryLanguage }

// RetryPolicy defines the retry behaviour for a Task state when errors occur.
type RetryPolicy struct {
	ErrorEquals     []string `json:"ErrorEquals"`
	IntervalSeconds int32    `json:"IntervalSeconds,omitempty"`
	MaxAttempts     int32    `json:"MaxAttempts,omitempty"`
	BackoffRate     float64  `json:"BackoffRate,omitempty"`
}

// CatchPolicy defines fallback behaviour for a Task state when an error occurs.
type CatchPolicy struct {
	ErrorEquals []string               `json:"ErrorEquals"`
	ResultPath  string                 `json:"ResultPath,omitempty"`
	Next        string                 `json:"Next"`
	Assign      map[string]interface{} `json:"Assign,omitempty"`
	Output      interface{}            `json:"Output,omitempty"`
}

// ChoiceRule defines a condition used in a Choice state to determine the next state
// transition.
type ChoiceRule struct {
	Variable             string                 `json:"Variable,omitempty"`
	Next                 string                 `json:"Next,omitempty"`
	And                  []*ChoiceRule          `json:"And,omitempty"`
	Or                   []*ChoiceRule          `json:"Or,omitempty"`
	Not                  *ChoiceRule            `json:"Not,omitempty"`
	StringEquals         map[string]string      `json:"StringEquals,omitempty"`
	StringLessThan       map[string]string      `json:"StringLessThan,omitempty"`
	StringGreaterThan    map[string]string      `json:"StringGreaterThan,omitempty"`
	NumericEquals        map[string]float64     `json:"NumericEquals,omitempty"`
	NumericLessThan      map[string]float64     `json:"NumericLessThan,omitempty"`
	NumericGreaterThan   map[string]float64     `json:"NumericGreaterThan,omitempty"`
	BooleanEquals        map[string]bool        `json:"BooleanEquals,omitempty"`
	TimestampEquals      map[string]string      `json:"TimestampEquals,omitempty"`
	TimestampLessThan    map[string]string      `json:"TimestampLessThan,omitempty"`
	TimestampGreaterThan map[string]string      `json:"TimestampGreaterThan,omitempty"`
	IsPresent            *bool                  `json:"IsPresent,omitempty"`
	Condition            string                 `json:"Condition,omitempty"`
	Assign               map[string]interface{} `json:"Assign,omitempty"`
}

// UnmarshalJSON deserialises a ChoiceRule from JSON, supporting both shorthand
// (single-value) and longhand (variable-keyed map) operator formats.
func (r *ChoiceRule) UnmarshalJSON(data []byte) error {
	type Alias ChoiceRule
	aux := &struct {
		*Alias
		StringEqualsRaw         interface{} `json:"StringEquals,omitempty"`
		StringLessThanRaw       interface{} `json:"StringLessThan,omitempty"`
		StringGreaterThanRaw    interface{} `json:"StringGreaterThan,omitempty"`
		NumericEqualsRaw        interface{} `json:"NumericEquals,omitempty"`
		NumericLessThanRaw      interface{} `json:"NumericLessThan,omitempty"`
		NumericGreaterThanRaw   interface{} `json:"NumericGreaterThan,omitempty"`
		BooleanEqualsRaw        interface{} `json:"BooleanEquals,omitempty"`
		TimestampEqualsRaw      interface{} `json:"TimestampEquals,omitempty"`
		TimestampLessThanRaw    interface{} `json:"TimestampLessThan,omitempty"`
		TimestampGreaterThanRaw interface{} `json:"TimestampGreaterThan,omitempty"`
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.StringEquals = convertToMapString(aux.StringEqualsRaw, r.Variable)
	r.StringLessThan = convertToMapString(aux.StringLessThanRaw, r.Variable)
	r.StringGreaterThan = convertToMapString(aux.StringGreaterThanRaw, r.Variable)
	r.NumericEquals = convertToMapFloat(aux.NumericEqualsRaw, r.Variable)
	r.NumericLessThan = convertToMapFloat(aux.NumericLessThanRaw, r.Variable)
	r.NumericGreaterThan = convertToMapFloat(aux.NumericGreaterThanRaw, r.Variable)
	r.BooleanEquals = convertToMapBool(aux.BooleanEqualsRaw, r.Variable)
	r.TimestampEquals = convertToMapString(aux.TimestampEqualsRaw, r.Variable)
	r.TimestampLessThan = convertToMapString(aux.TimestampLessThanRaw, r.Variable)
	r.TimestampGreaterThan = convertToMapString(aux.TimestampGreaterThanRaw, r.Variable)
	return nil
}

func convertToMapString(v interface{}, variable string) map[string]string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case string:
		if variable == "" {
			return nil
		}
		return map[string]string{variable: val}
	case map[string]interface{}:
		result := make(map[string]string)
		for k, iv := range val {
			if s, ok := iv.(string); ok {
				result[k] = s
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		return nil
	}
}

func convertToMapFloat(v interface{}, variable string) map[string]float64 {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case float64:
		if variable == "" {
			return nil
		}
		return map[string]float64{variable: val}
	case map[string]interface{}:
		result := make(map[string]float64)
		for k, iv := range val {
			if f, ok := iv.(float64); ok {
				result[k] = f
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		return nil
	}
}

func convertToMapBool(v interface{}, variable string) map[string]bool {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case bool:
		if variable == "" {
			return nil
		}
		return map[string]bool{variable: val}
	case map[string]interface{}:
		result := make(map[string]bool)
		for k, iv := range val {
			if b, ok := iv.(bool); ok {
				result[k] = b
			}
		}
		if len(result) == 0 {
			return nil
		}
		return result
	default:
		return nil
	}
}

// InputOutput represents input and output configuration for a state, supporting
// JSONPath, Parameters, and payload template values.
type InputOutput struct {
	Value                interface{} `json:"-"`
	Path                 string      `json:"Path,omitempty"`
	Parameters           *Parameters `json:"Parameters,omitempty"`
	PayloadTemplate      string      `json:"-"`
	PayloadTemplateValue interface{} `json:"-"`
	InputPath            string      `json:"InputPath,omitempty"`
	ItemsPath            string      `json:"ItemsPath,omitempty"`
}

// UnmarshalJSON deserialises an InputOutput from JSON, accepting either a raw string
// (payload template) or a structured object.
func (io *InputOutput) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		io.PayloadTemplate = s
		return nil
	}
	type Alias InputOutput
	return json.Unmarshal(data, (*Alias)(io))
}

// ResultSelector specifies a JSONPath filter applied to the output of a state before
// it is passed to the next state.
type ResultSelector struct {
	Fields map[string]interface{} `json:"-"`
}

// UnmarshalJSON deserialises a ResultSelector from JSON into its Fields map.
func (rs *ResultSelector) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &rs.Fields)
}

// MarshalJSON serialises a ResultSelector's Fields map to JSON.
func (rs *ResultSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.Fields)
}

// Parameters represents a map of parameter values passed to a state's resource.
type Parameters struct {
	Values map[string]interface{} `json:"-"`
}

// UnmarshalJSON deserialises Parameters from JSON into its Values map.
func (p *Parameters) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Values)
}

// MarshalJSON serialises Parameters' Values map to JSON.
func (p *Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Values)
}

// StateMachineListResult holds a paginated list of state machines.
type StateMachineListResult struct {
	StateMachines []*StateMachine
	NextToken     string
}

// ExecutionListResult holds a paginated list of state machine executions.
type ExecutionListResult struct {
	Executions []*Execution
	NextToken  string
}

// ActivityListResult holds a paginated list of activities.
type ActivityListResult struct {
	Activities []*Activity
	NextToken  string
}

// ActivityTask represents a single activity task dispatched to a worker.
type ActivityTask struct {
	TaskToken    string    `json:"taskToken"`
	ActivityArn  string    `json:"activityArn"`
	ExecutionArn string    `json:"executionArn"`
	Input        string    `json:"input"`
	Status       string    `json:"status"`
	Output       string    `json:"output,omitempty"`
	Error        string    `json:"error,omitempty"`
	Cause        string    `json:"cause,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	CompletedAt  time.Time `json:"completedAt,omitempty"`
	WorkerName   string    `json:"workerName,omitempty"`
}

// ActivityTaskResult contains the outcome of an activity task, including output or
// error.
type ActivityTaskResult struct {
	TaskToken string
	Output    string
	Error     error
}

// ActivityTaskScheduledEventDetails contains details emitted when an activity task is
// scheduled.
type ActivityTaskScheduledEventDetails struct {
	Resource         string `json:"resource"`
	Input            string `json:"input"`
	TaskToken        string `json:"taskToken"`
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

// StateMachineVersion represents a published version of a state machine definition.
type StateMachineVersion struct {
	StateMachineVersionArn string    `json:"stateMachineVersionArn"`
	StateMachineArn        string    `json:"stateMachineArn"`
	Version                int64     `json:"version"`
	Description            string    `json:"description,omitempty"`
	CreationDate           time.Time `json:"creationDate"`
	Definition             string    `json:"definition,omitempty"`
	RevisionId             string    `json:"revisionId,omitempty"`
}

// RoutingConfiguration maps a state machine version ARN to a traffic weight for use
// with aliases.
type RoutingConfiguration struct {
	StateMachineVersionArn string `json:"stateMachineVersionArn"`
	Weight                 int32  `json:"weight"`
}

// StateMachineAlias represents an alias that routes invocations to one or more
// versions of a state machine.
type StateMachineAlias struct {
	StateMachineAliasArn string                 `json:"stateMachineAliasArn"`
	StateMachineArn      string                 `json:"stateMachineArn,omitempty"`
	Name                 string                 `json:"name"`
	Description          string                 `json:"description,omitempty"`
	RoutingConfiguration []RoutingConfiguration `json:"routingConfiguration"`
	CreationDate         time.Time              `json:"creationDate"`
	UpdateDate           time.Time              `json:"updateDate,omitempty"`
}

// StateMachineVersionListResult holds a paginated list of state machine versions.
type StateMachineVersionListResult struct {
	Versions  []*StateMachineVersion
	NextToken string
}

// StateMachineAliasListResult holds a paginated list of state machine aliases.
type StateMachineAliasListResult struct {
	Aliases   []*StateMachineAlias
	NextToken string
}

// MapRun represents a distributed map state execution within a Step Functions
// state machine. Map runs are persisted in Pebble so they survive restarts.
type MapRun struct {
	MapRunArn                  string  `json:"mapRunArn"`
	ExecutionArn               string  `json:"executionArn"`
	StateMachineArn            string  `json:"stateMachineArn"`
	Name                       string  `json:"name"`
	Status                     string  `json:"status"`
	StartDate                  int64   `json:"startDate"`
	StopDate                   int64   `json:"stopDate,omitempty"`
	ItemCount                  int64   `json:"itemCount"`
	MaxConcurrency             int64   `json:"maxConcurrency"`
	ToleratedFailureCount      int64   `json:"toleratedFailureCount,omitempty"`
	ToleratedFailurePercentage float32 `json:"toleratedFailurePercentage,omitempty"`
	ItemsProcessedCount        int64   `json:"itemsProcessedCount"`
	ItemsFailedCount           int64   `json:"itemsFailedCount"`
	ItemsCancelledCount        int64   `json:"itemsCancelledCount"`
}

// MapRunListResult holds a paginated list of map runs.
type MapRunListResult struct {
	MapRuns   []*MapRun
	NextToken string
}
