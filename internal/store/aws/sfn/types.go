// Package stepfunction provides StepFunctions data store functionality for vorpalstacks.
package sfn

import (
	"encoding/json"
	"time"
)

// StateMachine represents a StepFunctions state machine.
type StateMachine struct {
	StateMachineArn string            `json:"stateMachineArn"`
	Name            string            `json:"name"`
	Definition      string            `json:"definition"`
	RoleArn         string            `json:"roleArn"`
	Type            string            `json:"type"`
	CreationDate    time.Time         `json:"creationDate"`
	UpdateDate      time.Time         `json:"updateDate"`
	Description     string            `json:"description"`
	Status          string            `json:"status"`
	Tags            map[string]string `json:"tags"`
}

// Execution represents a StepFunctions execution.
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

// ExecutionInputDetails represents details about the input to an execution.
type ExecutionInputDetails struct {
	Included bool   `json:"included"`
	Type     string `json:"type"`
}

// ExecutionOutputDetails represents details about the output of an execution.
type ExecutionOutputDetails struct {
	Included bool   `json:"included"`
	Type     string `json:"type"`
}

// ExecutionHistoryEvent represents an event in the execution history.
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
}

// ExecutionFailedEventDetails represents details of an execution failure event.
type ExecutionFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// TaskStartedEventDetails represents details of a task started event.
type TaskStartedEventDetails struct {
	Resource string `json:"resource"`
	Input    string `json:"input"`
}

// TaskSucceededEventDetails represents details of a task succeeded event.
type TaskSucceededEventDetails struct {
	Resource string `json:"resource"`
	Output   string `json:"output"`
}

// TaskFailedEventDetails represents details of a task failed event.
type TaskFailedEventDetails struct {
	Resource string `json:"resource"`
	Error    string `json:"error"`
	Cause    string `json:"cause"`
}

// ExecutionStartedEventDetails represents details of an execution started event.
type ExecutionStartedEventDetails struct {
	Input           string `json:"input"`
	RoleArn         string `json:"roleArn"`
	StateMachineArn string `json:"stateMachineArn"`
	Name            string `json:"name"`
}

// ExecutionSucceededEventDetails represents details of an execution succeeded event.
type ExecutionSucceededEventDetails struct {
	Output string `json:"output"`
}

// ExecutionAbortedEventDetails represents details of an execution aborted event.
type ExecutionAbortedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ExecutionTimedOutEventDetails represents details of an execution timed out event.
type ExecutionTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ChoiceStateEnteredEventDetails represents details when entering a choice state.
type ChoiceStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// ChoiceStateExitedEventDetails represents details when exiting a choice state.
type ChoiceStateExitedEventDetails struct {
	Output    string `json:"output"`
	Name      string `json:"name"`
	NextState string `json:"nextState"`
}

// WaitStateEnteredEventDetails represents details when entering a wait state.
type WaitStateEnteredEventDetails struct {
	Input     string `json:"input"`
	Name      string `json:"name"`
	Seconds   int64  `json:"seconds"`
	Timestamp string `json:"timestamp"`
}

// WaitStateExitedEventDetails represents details when exiting a wait state.
type WaitStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

// ParallelStateEnteredEventDetails represents details when entering a parallel state.
type ParallelStateEnteredEventDetails struct {
	Input    string   `json:"input"`
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

// ParallelStateExitedEventDetails represents details when exiting a parallel state.
type ParallelStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

// MapStateEnteredEventDetails represents details when entering a map state.
type MapStateEnteredEventDetails struct {
	Input          string `json:"input"`
	Name           string `json:"name"`
	ItemsProcessed int64  `json:"itemsProcessed"`
	ItemsFailed    int64  `json:"itemsFailed"`
}

// MapStateExitedEventDetails represents details when exiting a map state.
type MapStateExitedEventDetails struct {
	Output         string `json:"output"`
	Name           string `json:"name"`
	ItemsProcessed int64  `json:"itemsProcessed,omitempty"`
	ItemsFailed    int64  `json:"itemsFailed,omitempty"`
}

// PassStateEnteredEventDetails represents details when entering a pass state.
type PassStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// PassStateExitedEventDetails represents details when exiting a pass state.
type PassStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

// FailStateEnteredEventDetails represents details when entering a fail state.
type FailStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// SucceedStateEnteredEventDetails represents details when entering a succeed state.
type SucceedStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

// Activity represents a StepFunctions activity.
type Activity struct {
	ActivityArn  string    `json:"activityArn"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

// StateMachineDefinition represents the definition of a state machine.
type StateMachineDefinition struct {
	StartAt        string                 `json:"StartAt"`
	States         map[string]interface{} `json:"States"`
	TimeoutSeconds int32                  `json:"TimeoutSeconds,omitempty"`
	Version        string                 `json:"Version,omitempty"`
}

// State represents the interface for states in a state machine.
type State interface {
	GetType() string
	GetNext() string
	GetEnd() bool
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

// PassState represents a pass state in a StepFunctions state machine.
type PassState struct {
	Name           string          `json:"-"`
	Type           string          `json:"Type"`
	Next           string          `json:"Next,omitempty"`
	End            bool            `json:"End,omitempty"`
	Comment        string          `json:"Comment,omitempty"`
	Input          *InputOutput    `json:"Input,omitempty"`
	Output         *InputOutput    `json:"Output,omitempty"`
	Result         interface{}     `json:"Result,omitempty"`
	ResultPath     string          `json:"ResultPath,omitempty"`
	ResultSelector *ResultSelector `json:"ResultSelector,omitempty"`
	Parameters     *Parameters     `json:"Parameters,omitempty"`
}

// GetType returns the type of the pass state.
func (s *PassState) GetType() string { return s.Type }

// GetNext returns the next state after the pass state.
func (s *PassState) GetNext() string { return s.Next }

// GetEnd returns whether this is an end state.
func (s *PassState) GetEnd() bool { return s.End }

// GetInputPath returns the input path for the pass state.
func (s *PassState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the pass state.
func (s *PassState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetResultSelector returns the result selector for the pass state.
func (s *PassState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// TaskState represents a task state in a StepFunctions state machine.
type TaskState struct {
	Name             string          `json:"-"`
	Type             string          `json:"Type"`
	Next             string          `json:"Next,omitempty"`
	End              bool            `json:"End,omitempty"`
	Comment          string          `json:"Comment,omitempty"`
	Input            *InputOutput    `json:"Input,omitempty"`
	Output           *InputOutput    `json:"Output,omitempty"`
	Resource         string          `json:"Resource"`
	Parameters       *Parameters     `json:"Parameters,omitempty"`
	ResultPath       string          `json:"ResultPath,omitempty"`
	ResultSelector   *ResultSelector `json:"ResultSelector,omitempty"`
	Retry            []*RetryPolicy  `json:"Retry,omitempty"`
	Catch            []*CatchPolicy  `json:"Catch,omitempty"`
	TimeoutSeconds   int32           `json:"TimeoutSeconds,omitempty"`
	HeartbeatSeconds int32           `json:"HeartbeatSeconds,omitempty"`
}

// GetType returns the type of the task state.
func (s *TaskState) GetType() string { return s.Type }

// GetNext returns the next state after the task state.
func (s *TaskState) GetNext() string { return s.Next }

// GetEnd returns whether this is an end state.
func (s *TaskState) GetEnd() bool { return s.End }

// GetInputPath returns the input path for the task state.
func (s *TaskState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the task state.
func (s *TaskState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// GetResultSelector returns the result selector for the task state.
func (s *TaskState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// ChoiceState represents a choice state in a StepFunctions state machine.
type ChoiceState struct {
	Name    string        `json:"-"`
	Type    string        `json:"Type"`
	Comment string        `json:"Comment,omitempty"`
	Input   *InputOutput  `json:"Input,omitempty"`
	Output  *InputOutput  `json:"Output,omitempty"`
	Choices []*ChoiceRule `json:"Choices"`
	Default string        `json:"Default,omitempty"`
}

// GetType returns the type of the choice state.
func (s *ChoiceState) GetType() string { return s.Type }

// GetNext returns the next state (always empty for choice).
func (s *ChoiceState) GetNext() string { return "" }

// GetEnd returns whether this is an end state (always false for choice).
func (s *ChoiceState) GetEnd() bool { return false }

// GetInputPath returns the input path for the choice state.
func (s *ChoiceState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the choice state.
func (s *ChoiceState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// WaitState represents a wait state in a StepFunctions state machine.
type WaitState struct {
	Name          string       `json:"-"`
	Type          string       `json:"Type"`
	Next          string       `json:"Next,omitempty"`
	End           bool         `json:"End,omitempty"`
	Comment       string       `json:"Comment,omitempty"`
	Input         *InputOutput `json:"Input,omitempty"`
	Output        *InputOutput `json:"Output,omitempty"`
	Seconds       int32        `json:"Seconds,omitempty"`
	TimestampPath string       `json:"TimestampPath,omitempty"`
	SecondsPath   string       `json:"SecondsPath,omitempty"`
	Timestamp     string       `json:"Timestamp,omitempty"`
}

// GetType returns the type of the wait state.
func (s *WaitState) GetType() string { return s.Type }

// GetNext returns the next state after the wait state.
func (s *WaitState) GetNext() string { return s.Next }

// GetEnd returns whether this is an end state.
func (s *WaitState) GetEnd() bool { return s.End }

// GetInputPath returns the input path for the wait state.
func (s *WaitState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the wait state.
func (s *WaitState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// ParallelState represents a parallel state in a StepFunctions state machine.
type ParallelState struct {
	Name     string                    `json:"-"`
	Type     string                    `json:"Type"`
	Next     string                    `json:"Next,omitempty"`
	End      bool                      `json:"End,omitempty"`
	Comment  string                    `json:"Comment,omitempty"`
	Input    *InputOutput              `json:"Input,omitempty"`
	Output   *InputOutput              `json:"Output,omitempty"`
	Branches []*StateMachineDefinition `json:"Branches"`
	Retry    []*RetryPolicy            `json:"Retry,omitempty"`
	Catch    []*CatchPolicy            `json:"Catch,omitempty"`
}

// GetType returns the type of the parallel state.
func (s *ParallelState) GetType() string { return s.Type }

// GetNext returns the next state after the parallel state.
func (s *ParallelState) GetNext() string { return s.Next }

// GetEnd returns whether this is an end state.
func (s *ParallelState) GetEnd() bool { return s.End }

// GetInputPath returns the input path for the parallel state.
func (s *ParallelState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the parallel state.
func (s *ParallelState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// MapState represents a map state in a StepFunctions state machine.
type MapState struct {
	Name           string                  `json:"-"`
	Type           string                  `json:"Type"`
	Next           string                  `json:"Next,omitempty"`
	End            bool                    `json:"End,omitempty"`
	Comment        string                  `json:"Comment,omitempty"`
	Input          *InputOutput            `json:"Input,omitempty"`
	Output         *InputOutput            `json:"Output,omitempty"`
	Iterator       *StateMachineDefinition `json:"Iterator"`
	ItemsPath      string                  `json:"ItemsPath,omitempty"`
	MaxConcurrency int32                   `json:"MaxConcurrency,omitempty"`
	Retry          []*RetryPolicy          `json:"Retry,omitempty"`
	Catch          []*CatchPolicy          `json:"Catch,omitempty"`
}

// GetType returns the type of the map state.
func (s *MapState) GetType() string { return s.Type }

// GetNext returns the next state after the map state.
func (s *MapState) GetNext() string { return s.Next }

// GetEnd returns whether this is an end state.
func (s *MapState) GetEnd() bool { return s.End }

// GetInputPath returns the input path for the map state.
func (s *MapState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the map state.
func (s *MapState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// FailState represents a fail state in a StepFunctions state machine.
type FailState struct {
	Name    string `json:"-"`
	Type    string `json:"Type"`
	Comment string `json:"Comment,omitempty"`
	Cause   string `json:"Cause,omitempty"`
	Error   string `json:"Error,omitempty"`
}

// GetType returns the type of the fail state.
func (s *FailState) GetType() string { return s.Type }

// GetNext returns the next state (always empty for fail).
func (s *FailState) GetNext() string { return "" }

// GetEnd returns whether this is an end state (always true for fail).
func (s *FailState) GetEnd() bool { return true }

// SucceedState represents a succeed state in a StepFunctions state machine.
type SucceedState struct {
	Name    string       `json:"-"`
	Type    string       `json:"Type"`
	Comment string       `json:"Comment,omitempty"`
	Input   *InputOutput `json:"Input,omitempty"`
	Output  *InputOutput `json:"Output,omitempty"`
}

// GetType returns the type of the succeed state.
func (s *SucceedState) GetType() string { return s.Type }

// GetNext returns the next state (always empty for succeed).
func (s *SucceedState) GetNext() string { return "" }

// GetEnd returns whether this is an end state (always true for succeed).
func (s *SucceedState) GetEnd() bool { return true }

// GetInputPath returns the input path for the succeed state.
func (s *SucceedState) GetInputPath() string { return getInputPathFromInputOutput(s.Input) }

// GetOutputPath returns the output path for the succeed state.
func (s *SucceedState) GetOutputPath() string { return getOutputPathFromInputOutput(s.Output) }

// RetryPolicy represents a retry policy for a task state.
type RetryPolicy struct {
	ErrorEquals     []string `json:"ErrorEquals"`
	IntervalSeconds int32    `json:"IntervalSeconds,omitempty"`
	MaxAttempts     int32    `json:"MaxAttempts,omitempty"`
	BackoffRate     float64  `json:"BackoffRate,omitempty"`
}

// CatchPolicy represents a catch policy for a task state.
type CatchPolicy struct {
	ErrorEquals []string `json:"ErrorEquals"`
	ResultPath  string   `json:"ResultPath,omitempty"`
	Next        string   `json:"Next"`
}

// ChoiceRule represents a choice rule for a choice state.
type ChoiceRule struct {
	Variable             string             `json:"Variable,omitempty"`
	Next                 string             `json:"Next,omitempty"`
	And                  []*ChoiceRule      `json:"And,omitempty"`
	Or                   []*ChoiceRule      `json:"Or,omitempty"`
	Not                  *ChoiceRule        `json:"Not,omitempty"`
	StringEquals         map[string]string  `json:"StringEquals,omitempty"`
	StringLessThan       map[string]string  `json:"StringLessThan,omitempty"`
	StringGreaterThan    map[string]string  `json:"StringGreaterThan,omitempty"`
	NumericEquals        map[string]float64 `json:"NumericEquals,omitempty"`
	NumericLessThan      map[string]float64 `json:"NumericLessThan,omitempty"`
	NumericGreaterThan   map[string]float64 `json:"NumericGreaterThan,omitempty"`
	BooleanEquals        map[string]bool    `json:"BooleanEquals,omitempty"`
	TimestampEquals      map[string]string  `json:"TimestampEquals,omitempty"`
	TimestampLessThan    map[string]string  `json:"TimestampLessThan,omitempty"`
	TimestampGreaterThan map[string]string  `json:"TimestampGreaterThan,omitempty"`
	IsPresent            *bool              `json:"IsPresent,omitempty"`
}

// UnmarshalJSON handles both scalar and map formats for comparison operators
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

// InputOutput represents input/output configuration for a state.
type InputOutput struct {
	Value                interface{} `json:"-"`
	Path                 string      `json:"Path,omitempty"`
	Parameters           *Parameters `json:"Parameters,omitempty"`
	PayloadTemplate      string      `json:"-"`
	PayloadTemplateValue interface{} `json:"-"`
	InputPath            string      `json:"InputPath,omitempty"`
	ItemsPath            string      `json:"ItemsPath,omitempty"`
}

// ResultSelector represents a result selector for a state.
type ResultSelector struct {
	Fields map[string]interface{} `json:"-"`
}

// UnmarshalJSON unmarshals JSON data into the ResultSelector fields.
func (rs *ResultSelector) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &rs.Fields)
}

// MarshalJSON marshals the ResultSelector fields to JSON.
func (rs *ResultSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.Fields)
}

// Parameters represents parameters for a state.
type Parameters struct {
	Values map[string]interface{} `json:"-"`
}

// UnmarshalJSON unmarshals JSON data into the Parameters values.
func (p *Parameters) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Values)
}

// MarshalJSON marshals the Parameters values to JSON.
func (p *Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Values)
}

// StateMachineListResult represents the result of listing state machines.
type StateMachineListResult struct {
	StateMachines []*StateMachine
	NextToken     string
}

// ExecutionListResult represents the result of listing executions.
type ExecutionListResult struct {
	Executions []*Execution
	NextToken  string
}

// ActivityListResult represents the result of listing activities.
type ActivityListResult struct {
	Activities []*Activity
	NextToken  string
}

// ActivityTask represents a task in a StepFunctions activity.
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

// ActivityTaskResult represents the result of an activity task.
type ActivityTaskResult struct {
	TaskToken string
	Output    string
	Error     error
}

// ActivityTaskScheduledEventDetails represents details of an activity task scheduled event.
type ActivityTaskScheduledEventDetails struct {
	Resource  string `json:"resource"`
	Input     string `json:"input"`
	TaskToken string `json:"taskToken"`
}

// ActivityTaskStartedEventDetails represents details of an activity task started event.
type ActivityTaskStartedEventDetails struct {
	WorkerName string `json:"workerName"`
}

// ActivityTaskSucceededEventDetails represents details of an activity task succeeded event.
type ActivityTaskSucceededEventDetails struct {
	Output string `json:"output"`
}

// ActivityTaskFailedEventDetails represents details of an activity task failed event.
type ActivityTaskFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// ActivityTaskTimedOutEventDetails represents details of an activity task timed out event.
type ActivityTaskTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}
