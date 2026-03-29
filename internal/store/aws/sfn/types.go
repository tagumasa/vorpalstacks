package sfn

import (
	"encoding/json"
	"time"
)

type StateMachine struct {
	StateMachineArn    string              `json:"stateMachineArn"`
	Name               string              `json:"name"`
	Definition         string              `json:"definition"`
	RoleArn            string              `json:"roleArn"`
	Type               string              `json:"type"`
	CreationDate       time.Time           `json:"creationDate"`
	UpdateDate         time.Time           `json:"updateDate"`
	Description        string              `json:"description"`
	Status             string              `json:"status"`
	Tags               map[string]string   `json:"tags"`
	VariableReferences map[string][]string `json:"-"`
}

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

type ExecutionInputDetails struct {
	Included bool   `json:"included"`
	Type     string `json:"type"`
}

type ExecutionOutputDetails struct {
	Included bool   `json:"included"`
	Type     string `json:"type"`
}

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

type ExecutionFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

type EvaluationFailedEventDetails struct {
	State    string `json:"state"`
	Cause    string `json:"cause"`
	Error    string `json:"error"`
	Location string `json:"location"`
}

type TaskStartedEventDetails struct {
	Resource string `json:"resource"`
	Input    string `json:"input"`
}

type TaskSucceededEventDetails struct {
	Resource string `json:"resource"`
	Output   string `json:"output"`
}

type TaskFailedEventDetails struct {
	Resource string `json:"resource"`
	Error    string `json:"error"`
	Cause    string `json:"cause"`
}

type ExecutionStartedEventDetails struct {
	Input           string `json:"input"`
	RoleArn         string `json:"roleArn"`
	StateMachineArn string `json:"stateMachineArn"`
	Name            string `json:"name"`
}

type ExecutionSucceededEventDetails struct {
	Output string `json:"output"`
}

type ExecutionAbortedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

type ExecutionTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

type ChoiceStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

type ChoiceStateExitedEventDetails struct {
	Output    string `json:"output"`
	Name      string `json:"name"`
	NextState string `json:"nextState"`
}

type WaitStateEnteredEventDetails struct {
	Input     string `json:"input"`
	Name      string `json:"name"`
	Seconds   int64  `json:"seconds"`
	Timestamp string `json:"timestamp"`
}

type WaitStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

type ParallelStateEnteredEventDetails struct {
	Input    string   `json:"input"`
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

type ParallelStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

type MapStateEnteredEventDetails struct {
	Input          string `json:"input"`
	Name           string `json:"name"`
	ItemsProcessed int64  `json:"itemsProcessed"`
	ItemsFailed    int64  `json:"itemsFailed"`
}

type MapStateExitedEventDetails struct {
	Output         string `json:"output"`
	Name           string `json:"name"`
	ItemsProcessed int64  `json:"itemsProcessed,omitempty"`
	ItemsFailed    int64  `json:"itemsFailed,omitempty"`
}

type PassStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

type PassStateExitedEventDetails struct {
	Output string `json:"output"`
	Name   string `json:"name"`
}

type FailStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
	Error string `json:"error"`
	Cause string `json:"cause"`
}

type SucceedStateEnteredEventDetails struct {
	Input string `json:"input"`
	Name  string `json:"name"`
}

type Activity struct {
	ActivityArn  string    `json:"activityArn"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creationDate"`
}

type StateMachineDefinition struct {
	StartAt        string                 `json:"StartAt"`
	States         map[string]interface{} `json:"States"`
	TimeoutSeconds int32                  `json:"TimeoutSeconds,omitempty"`
	Version        string                 `json:"Version,omitempty"`
	QueryLanguage  string                 `json:"QueryLanguage,omitempty"`
}

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

func (s *PassState) GetType() string                    { return s.Type }
func (s *PassState) GetNext() string                    { return s.Next }
func (s *PassState) GetEnd() bool                       { return s.End }
func (s *PassState) GetInputPath() string               { return getInputPathFromInputOutput(s.Input) }
func (s *PassState) GetOutputPath() string              { return getOutputPathFromInputOutput(s.Output) }
func (s *PassState) GetResultSelector() *ResultSelector { return s.ResultSelector }
func (s *PassState) GetQueryLanguage() string           { return s.QueryLanguage }

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

func (s *TaskState) GetType() string                    { return s.Type }
func (s *TaskState) GetNext() string                    { return s.Next }
func (s *TaskState) GetEnd() bool                       { return s.End }
func (s *TaskState) GetInputPath() string               { return getInputPathFromInputOutput(s.Input) }
func (s *TaskState) GetOutputPath() string              { return getOutputPathFromInputOutput(s.Output) }
func (s *TaskState) GetResultSelector() *ResultSelector { return s.ResultSelector }
func (s *TaskState) GetQueryLanguage() string           { return s.QueryLanguage }

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

func (s *ChoiceState) GetType() string          { return s.Type }
func (s *ChoiceState) GetNext() string          { return "" }
func (s *ChoiceState) GetEnd() bool             { return false }
func (s *ChoiceState) GetInputPath() string     { return getInputPathFromInputOutput(s.Input) }
func (s *ChoiceState) GetOutputPath() string    { return getOutputPathFromInputOutput(s.Output) }
func (s *ChoiceState) GetQueryLanguage() string { return s.QueryLanguage }

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

func (s *WaitState) GetType() string          { return s.Type }
func (s *WaitState) GetNext() string          { return s.Next }
func (s *WaitState) GetEnd() bool             { return s.End }
func (s *WaitState) GetInputPath() string     { return getInputPathFromInputOutput(s.Input) }
func (s *WaitState) GetOutputPath() string    { return getOutputPathFromInputOutput(s.Output) }
func (s *WaitState) GetQueryLanguage() string { return s.QueryLanguage }

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

func (s *ParallelState) GetType() string          { return s.Type }
func (s *ParallelState) GetNext() string          { return s.Next }
func (s *ParallelState) GetEnd() bool             { return s.End }
func (s *ParallelState) GetInputPath() string     { return getInputPathFromInputOutput(s.Input) }
func (s *ParallelState) GetOutputPath() string    { return getOutputPathFromInputOutput(s.Output) }
func (s *ParallelState) GetQueryLanguage() string { return s.QueryLanguage }

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

func (s *MapState) GetType() string          { return s.Type }
func (s *MapState) GetNext() string          { return s.Next }
func (s *MapState) GetEnd() bool             { return s.End }
func (s *MapState) GetInputPath() string     { return getInputPathFromInputOutput(s.Input) }
func (s *MapState) GetOutputPath() string    { return getOutputPathFromInputOutput(s.Output) }
func (s *MapState) GetQueryLanguage() string { return s.QueryLanguage }

func (s *MapState) GetIterator() *StateMachineDefinition {
	if s.Iterator != nil {
		return s.Iterator
	}
	return s.ItemProcessor
}

type FailState struct {
	Name          string      `json:"-"`
	Type          string      `json:"Type"`
	Comment       string      `json:"Comment,omitempty"`
	Cause         interface{} `json:"Cause,omitempty"`
	Error         string      `json:"Error,omitempty"`
	QueryLanguage string      `json:"QueryLanguage,omitempty"`
}

func (s *FailState) GetType() string          { return s.Type }
func (s *FailState) GetNext() string          { return "" }
func (s *FailState) GetEnd() bool             { return true }
func (s *FailState) GetQueryLanguage() string { return s.QueryLanguage }

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

func (s *SucceedState) GetType() string          { return s.Type }
func (s *SucceedState) GetNext() string          { return "" }
func (s *SucceedState) GetEnd() bool             { return true }
func (s *SucceedState) GetInputPath() string     { return getInputPathFromInputOutput(s.Input) }
func (s *SucceedState) GetOutputPath() string    { return getOutputPathFromInputOutput(s.Output) }
func (s *SucceedState) GetQueryLanguage() string { return s.QueryLanguage }

type RetryPolicy struct {
	ErrorEquals     []string `json:"ErrorEquals"`
	IntervalSeconds int32    `json:"IntervalSeconds,omitempty"`
	MaxAttempts     int32    `json:"MaxAttempts,omitempty"`
	BackoffRate     float64  `json:"BackoffRate,omitempty"`
}

type CatchPolicy struct {
	ErrorEquals []string               `json:"ErrorEquals"`
	ResultPath  string                 `json:"ResultPath,omitempty"`
	Next        string                 `json:"Next"`
	Assign      map[string]interface{} `json:"Assign,omitempty"`
	Output      interface{}            `json:"Output,omitempty"`
}

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

type InputOutput struct {
	Value                interface{} `json:"-"`
	Path                 string      `json:"Path,omitempty"`
	Parameters           *Parameters `json:"Parameters,omitempty"`
	PayloadTemplate      string      `json:"-"`
	PayloadTemplateValue interface{} `json:"-"`
	InputPath            string      `json:"InputPath,omitempty"`
	ItemsPath            string      `json:"ItemsPath,omitempty"`
}

func (io *InputOutput) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		io.PayloadTemplate = s
		return nil
	}
	type Alias InputOutput
	return json.Unmarshal(data, (*Alias)(io))
}

type ResultSelector struct {
	Fields map[string]interface{} `json:"-"`
}

func (rs *ResultSelector) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &rs.Fields)
}

func (rs *ResultSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.Fields)
}

type Parameters struct {
	Values map[string]interface{} `json:"-"`
}

func (p *Parameters) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Values)
}

func (p *Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Values)
}

type StateMachineListResult struct {
	StateMachines []*StateMachine
	NextToken     string
}

type ExecutionListResult struct {
	Executions []*Execution
	NextToken  string
}

type ActivityListResult struct {
	Activities []*Activity
	NextToken  string
}

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

type ActivityTaskResult struct {
	TaskToken string
	Output    string
	Error     error
}

type ActivityTaskScheduledEventDetails struct {
	Resource         string `json:"resource"`
	Input            string `json:"input"`
	TaskToken        string `json:"taskToken"`
	HeartbeatSeconds int32  `json:"heartbeatInSeconds,omitempty"`
}

type ActivityTaskStartedEventDetails struct {
	WorkerName string `json:"workerName"`
}

type ActivityTaskSucceededEventDetails struct {
	Output string `json:"output"`
}

type ActivityTaskFailedEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

type ActivityTaskTimedOutEventDetails struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}
