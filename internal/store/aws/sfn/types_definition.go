package sfn

import "encoding/json"

// This file holds the state machine definition model: the definition
// envelope, the State interface and every typed state with its retry,
// catch, payload-template and Map/Parallel configuration shapes.

// StateMachineDefinition represents the Amazon States Language definition of a state
// machine.
type StateMachineDefinition struct {
	StartAt        string                 `json:"StartAt"`
	States         map[string]interface{} `json:"States"`
	TimeoutSeconds int32                  `json:"TimeoutSeconds,omitempty"`
	Version        string                 `json:"Version,omitempty"`
	QueryLanguage  string                 `json:"QueryLanguage,omitempty"`
	// ProcessorConfig carries the ItemProcessor processing-mode settings
	// (Mode INLINE or DISTRIBUTED plus ExecutionType). It only carries
	// meaning on definitions used as a Map ItemProcessor; the top-level
	// machine definition never sets it.
	ProcessorConfig *ProcessorConfig `json:"ProcessorConfig,omitempty"`
}

// State defines the common interface for all state machine state types.
type State interface {
	GetType() string
	GetNext() string
	GetEnd() bool
	GetQueryLanguage() string
}

func getOutputPathFromInputOutput(io *InputOutput) string {
	if io == nil {
		return ""
	}
	return io.Path
}

// PassState represents a Pass state that simply passes its input to its output,
// optionally transforming it.
// statePathFilters carries the path-filter configuration shared by every
// state type: the JSONPath-era InputPath/OutputPath fields and the
// JSONata Output union object that may carry them. The embedded fields
// marshal flat, so each state's wire shape is unchanged. A state-level
// Input member carries no documented AWS field semantics and is not
// parsed — the definition validator rejects it as an unknown field.
type statePathFilters struct {
	Output     *InputOutput `json:"Output,omitempty"`
	InputPath  string       `json:"InputPath,omitempty"`
	OutputPath string       `json:"OutputPath,omitempty"`
}

// GetInputPath returns the input path filter applied to the state input.
func (f statePathFilters) GetInputPath() string {
	return f.InputPath
}

// GetOutputPath returns the output path filter applied to the state output.
func (f statePathFilters) GetOutputPath() string {
	if f.OutputPath != "" {
		return f.OutputPath
	}
	return getOutputPathFromInputOutput(f.Output)
}

type PassState struct {
	statePathFilters
	Name           string                 `json:"-"`
	Type           string                 `json:"Type"`
	Next           string                 `json:"Next,omitempty"`
	End            bool                   `json:"End,omitempty"`
	Comment        string                 `json:"Comment,omitempty"`
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

// GetResultSelector returns the result selector applied to the state output.
func (s *PassState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// GetQueryLanguage returns the query language used by this state.
func (s *PassState) GetQueryLanguage() string { return s.QueryLanguage }

// TaskState represents a Task state that executes a unit of work, such as calling an
// activity or a Lambda function.
type TaskState struct {
	statePathFilters
	Name                 string                 `json:"-"`
	Type                 string                 `json:"Type"`
	Next                 string                 `json:"Next,omitempty"`
	End                  bool                   `json:"End,omitempty"`
	Comment              string                 `json:"Comment,omitempty"`
	OutputRaw            json.RawMessage        `json:"-"`
	Resource             string                 `json:"Resource"`
	Parameters           *Parameters            `json:"Parameters,omitempty"`
	ResultPath           string                 `json:"ResultPath,omitempty"`
	ResultSelector       *ResultSelector        `json:"ResultSelector,omitempty"`
	Retry                []*RetryPolicy         `json:"Retry,omitempty"`
	Catch                []*CatchPolicy         `json:"Catch,omitempty"`
	TimeoutSeconds       interface{}            `json:"TimeoutSeconds,omitempty"`
	HeartbeatSeconds     interface{}            `json:"HeartbeatSeconds,omitempty"`
	TimeoutSecondsPath   string                 `json:"TimeoutSecondsPath,omitempty"`
	HeartbeatSecondsPath string                 `json:"HeartbeatSecondsPath,omitempty"`
	QueryLanguage        string                 `json:"QueryLanguage,omitempty"`
	Arguments            interface{}            `json:"Arguments,omitempty"`
	Credentials          interface{}            `json:"Credentials,omitempty"`
	Assign               map[string]interface{} `json:"Assign,omitempty"`
	JSONataOutput        interface{}            `json:"-"`
}

// GetType returns the state type identifier (e.g. "Task").
func (s *TaskState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to.
func (s *TaskState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *TaskState) GetEnd() bool { return s.End }

// GetResultSelector returns the result selector applied to the state output.
func (s *TaskState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// GetQueryLanguage returns the query language used by this state.
func (s *TaskState) GetQueryLanguage() string { return s.QueryLanguage }

// numericAsInt32 coerces a decoded JSON number held as interface{} —
// float64 through the default decoder, json.Number through the
// UseNumber decoder, int32 from typed construction — to an int32
// seconds value, with zero for absent or non-numeric values.
func numericAsInt32(v interface{}) int32 {
	switch n := v.(type) {
	case float64:
		return int32(n)
	case int32:
		return n
	case json.Number:
		i, _ := n.Int64()
		return int32(i)
	}
	return 0
}

// GetTimeoutSeconds returns the timeout in seconds for the task execution.
func (s *TaskState) GetTimeoutSeconds() int32 {
	return numericAsInt32(s.TimeoutSeconds)
}

// GetHeartbeatSeconds returns the heartbeat interval in seconds for the task.
func (s *TaskState) GetHeartbeatSeconds() int32 {
	return numericAsInt32(s.HeartbeatSeconds)
}

// ChoiceState represents a Choice state that adds branching logic to a state machine.
type ChoiceState struct {
	statePathFilters
	Name          string                 `json:"-"`
	Type          string                 `json:"Type"`
	Comment       string                 `json:"Comment,omitempty"`
	Choices       []*ChoiceRule          `json:"Choices"`
	Default       string                 `json:"Default,omitempty"`
	QueryLanguage string                 `json:"QueryLanguage,omitempty"`
	Assign        map[string]interface{} `json:"Assign,omitempty"`
}

// GetType returns the state type identifier (e.g. "Choice").
func (s *ChoiceState) GetType() string { return s.Type }

// GetNext always returns empty for Choice states; transitions are determined by rules.
func (s *ChoiceState) GetNext() string { return "" }

// GetEnd always returns false for Choice states.
func (s *ChoiceState) GetEnd() bool { return false }

// GetQueryLanguage returns the query language used by this state.
func (s *ChoiceState) GetQueryLanguage() string { return s.QueryLanguage }

// WaitState represents a Wait state that delays the state machine from continuing for
// a specified time.
type WaitState struct {
	statePathFilters
	Name          string                 `json:"-"`
	Type          string                 `json:"Type"`
	Next          string                 `json:"Next,omitempty"`
	End           bool                   `json:"End,omitempty"`
	Comment       string                 `json:"Comment,omitempty"`
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

// GetQueryLanguage returns the query language used by this state.
func (s *WaitState) GetQueryLanguage() string { return s.QueryLanguage }

// GetSeconds returns the number of seconds to wait before transitioning.
func (s *WaitState) GetSeconds() int32 {
	return numericAsInt32(s.Seconds)
}

// ParallelState represents a Parallel state that executes multiple branches
// concurrently.
type ParallelState struct {
	statePathFilters
	Name           string                    `json:"-"`
	Type           string                    `json:"Type"`
	Next           string                    `json:"Next,omitempty"`
	End            bool                      `json:"End,omitempty"`
	Comment        string                    `json:"Comment,omitempty"`
	OutputRaw      json.RawMessage           `json:"-"`
	Branches       []*StateMachineDefinition `json:"Branches"`
	Parameters     *Parameters               `json:"Parameters,omitempty"`
	ResultPath     string                    `json:"ResultPath,omitempty"`
	ResultSelector *ResultSelector           `json:"ResultSelector,omitempty"`
	Retry          []*RetryPolicy            `json:"Retry,omitempty"`
	Catch          []*CatchPolicy            `json:"Catch,omitempty"`
	QueryLanguage  string                    `json:"QueryLanguage,omitempty"`
	Arguments      interface{}               `json:"Arguments,omitempty"`
	Assign         map[string]interface{}    `json:"Assign,omitempty"`
	JSONataOutput  interface{}               `json:"-"`
}

// GetType returns the state type identifier (e.g. "Parallel").
func (s *ParallelState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to after all branches complete.
func (s *ParallelState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *ParallelState) GetEnd() bool { return s.End }

// GetQueryLanguage returns the query language used by this state.
func (s *ParallelState) GetQueryLanguage() string { return s.QueryLanguage }

// GetResultSelector returns the result selector applied to the combined
// branch output array before ResultPath.
func (s *ParallelState) GetResultSelector() *ResultSelector { return s.ResultSelector }

// MapState represents a Map state that processes a collection of items by applying a
// sub-state machine to each.
type MapState struct {
	statePathFilters
	Name                           string                  `json:"-"`
	Type                           string                  `json:"Type"`
	Next                           string                  `json:"Next,omitempty"`
	End                            bool                    `json:"End,omitempty"`
	Comment                        string                  `json:"Comment,omitempty"`
	OutputRaw                      json.RawMessage         `json:"-"`
	Iterator                       *StateMachineDefinition `json:"Iterator,omitempty"`
	ItemProcessor                  *StateMachineDefinition `json:"ItemProcessor,omitempty"`
	ItemsPath                      string                  `json:"ItemsPath,omitempty"`
	MaxConcurrency                 interface{}             `json:"MaxConcurrency,omitempty"`
	MaxConcurrencyPath             string                  `json:"MaxConcurrencyPath,omitempty"`
	Parameters                     *Parameters             `json:"Parameters,omitempty"`
	ResultPath                     string                  `json:"ResultPath,omitempty"`
	ResultSelector                 *ResultSelector         `json:"ResultSelector,omitempty"`
	Retry                          []*RetryPolicy          `json:"Retry,omitempty"`
	Catch                          []*CatchPolicy          `json:"Catch,omitempty"`
	QueryLanguage                  string                  `json:"QueryLanguage,omitempty"`
	Items                          interface{}             `json:"Items,omitempty"`
	ItemSelector                   interface{}             `json:"ItemSelector,omitempty"`
	Assign                         map[string]interface{}  `json:"Assign,omitempty"`
	Label                          string                  `json:"Label,omitempty"`
	ItemReader                     *ItemReaderConfig       `json:"ItemReader,omitempty"`
	ItemBatcher                    *ItemBatcherConfig      `json:"ItemBatcher,omitempty"`
	ResultWriter                   *ResultWriterConfig     `json:"ResultWriter,omitempty"`
	ToleratedFailureCount          interface{}             `json:"ToleratedFailureCount,omitempty"`
	ToleratedFailureCountPath      string                  `json:"ToleratedFailureCountPath,omitempty"`
	ToleratedFailurePercentage     interface{}             `json:"ToleratedFailurePercentage,omitempty"`
	ToleratedFailurePercentagePath string                  `json:"ToleratedFailurePercentagePath,omitempty"`
	JSONataOutput                  interface{}             `json:"-"`
}

// ProcessorConfig is the ItemProcessor processing-mode configuration of a
// Map state: Mode INLINE or DISTRIBUTED, plus the child execution type
// required when the mode is DISTRIBUTED (Distributed Map state
// documentation).
type ProcessorConfig struct {
	Mode          string `json:"Mode,omitempty"`
	ExecutionType string `json:"ExecutionType,omitempty"`
}

// ItemReaderConfig is the Distributed Map ItemReader: the dataset source
// and its parsing configuration. Parameters carries the JSONPath argument
// object (Bucket, Key, Prefix, VersionId, ExpectedBucketOwner); Arguments
// is the JSONata equivalent (ItemReader documentation).
type ItemReaderConfig struct {
	Resource     string                  `json:"Resource"`
	Parameters   json.RawMessage         `json:"Parameters,omitempty"`
	Arguments    json.RawMessage         `json:"Arguments,omitempty"`
	ReaderConfig *ItemReaderReaderConfig `json:"ReaderConfig,omitempty"`
}

// ItemReaderReaderConfig is the ReaderConfig member of an ItemReader: the
// dataset format and limits applied while reading items.
type ItemReaderReaderConfig struct {
	InputType         string   `json:"InputType,omitempty"`
	CSVHeaderLocation string   `json:"CSVHeaderLocation,omitempty"`
	CSVHeaders        []string `json:"CSVHeaders,omitempty"`
	CSVDelimiter      string   `json:"CSVDelimiter,omitempty"`
	MaxItems          *int64   `json:"MaxItems,omitempty"`
	MaxItemsPath      string   `json:"MaxItemsPath,omitempty"`
	ItemsPointer      string   `json:"ItemsPointer,omitempty"`
	Transformation    string   `json:"Transformation,omitempty"`
	ManifestType      string   `json:"ManifestType,omitempty"`
}

// ItemBatcherConfig is the Map ItemBatcher: the batch sizing and fixed
// input applied when a Map state groups items into the Items array each
// child workflow execution receives. The per-batch byte cap shares the
// 256 KiB child-execution input bound, and an unspecified cap defaults to
// it (ItemBatcher documentation). MaxItemsPerBatch is an interface so a
// JSONata state can carry an expression string instead of a literal count.
type ItemBatcherConfig struct {
	MaxItemsPerBatch          interface{} `json:"MaxItemsPerBatch,omitempty"`
	MaxItemsPerBatchPath      string      `json:"MaxItemsPerBatchPath,omitempty"`
	MaxInputBytesPerBatch     interface{} `json:"MaxInputBytesPerBatch,omitempty"`
	MaxInputBytesPerBatchPath string      `json:"MaxInputBytesPerBatchPath,omitempty"`
	BatchInput                interface{} `json:"BatchInput,omitempty"`
	BatchInputPath            string      `json:"BatchInputPath,omitempty"`
}

// ResultWriterConfig is the Distributed Map ResultWriter: the S3 export
// destination and the output formatting options. At least WriterConfig
// alone, or Resource with Parameters, must be present (ResultWriter
// documentation).
type ResultWriterConfig struct {
	Resource     string                    `json:"Resource,omitempty"`
	Parameters   json.RawMessage           `json:"Parameters,omitempty"`
	Arguments    json.RawMessage           `json:"Arguments,omitempty"`
	WriterConfig *ResultWriterWriterConfig `json:"WriterConfig,omitempty"`
}

// ResultWriterWriterConfig is the WriterConfig member of a ResultWriter:
// the result transformation and the exported file format.
type ResultWriterWriterConfig struct {
	Transformation string `json:"Transformation,omitempty"`
	OutputType     string `json:"OutputType,omitempty"`
}

// GetType returns the state type identifier (e.g. "Map").
func (s *MapState) GetType() string { return s.Type }

// GetNext returns the name of the next state to transition to after map processing completes.
func (s *MapState) GetNext() string { return s.Next }

// GetEnd reports whether this state is a terminal state.
func (s *MapState) GetEnd() bool { return s.End }

// GetQueryLanguage returns the query language used by this state.
func (s *MapState) GetQueryLanguage() string { return s.QueryLanguage }

// GetResultSelector returns the result selector applied to the Map result
// array before ResultPath.
func (s *MapState) GetResultSelector() *ResultSelector { return s.ResultSelector }

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
	Error         interface{} `json:"Error,omitempty"`
	ErrorPath     string      `json:"ErrorPath,omitempty"`
	CausePath     string      `json:"CausePath,omitempty"`
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

// GetError returns the failure error name as a string.
func (s *FailState) GetError() string {
	if s.Error == nil {
		return ""
	}
	switch v := s.Error.(type) {
	case string:
		return v
	default:
		return ""
	}
}

// SucceedState represents a Succeed state that stops the execution successfully.
type SucceedState struct {
	statePathFilters
	Name          string          `json:"-"`
	Type          string          `json:"Type"`
	Comment       string          `json:"Comment,omitempty"`
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

// GetQueryLanguage returns the query language used by this state.
func (s *SucceedState) GetQueryLanguage() string { return s.QueryLanguage }

// RetryPolicy defines the retry behaviour for a Task state when errors occur.
type RetryPolicy struct {
	ErrorEquals     []string `json:"ErrorEquals"`
	IntervalSeconds int32    `json:"IntervalSeconds,omitempty"`
	MaxAttempts     int32    `json:"MaxAttempts,omitempty"`
	BackoffRate     float64  `json:"BackoffRate,omitempty"`
	MaxDelaySeconds int32    `json:"MaxDelaySeconds,omitempty"`
	JitterStrategy  string   `json:"JitterStrategy,omitempty"`
}

// UnmarshalJSON applies the documented default for MaxAttempts when the
// member is absent ("A positive integer that represents the maximum number
// of retry attempts (3 by default). … A value of 0 specifies that the
// error is never retried"): an explicit zero keeps its never-retry meaning,
// so the default cannot live at the use site where zero is
// indistinguishable from absent.
func (r *RetryPolicy) UnmarshalJSON(data []byte) error {
	var raw struct {
		ErrorEquals     []string `json:"ErrorEquals"`
		IntervalSeconds int32    `json:"IntervalSeconds"`
		MaxAttempts     *int32   `json:"MaxAttempts"`
		BackoffRate     float64  `json:"BackoffRate"`
		MaxDelaySeconds int32    `json:"MaxDelaySeconds"`
		JitterStrategy  string   `json:"JitterStrategy"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.ErrorEquals = raw.ErrorEquals
	r.IntervalSeconds = raw.IntervalSeconds
	r.MaxAttempts = 3
	if raw.MaxAttempts != nil {
		r.MaxAttempts = *raw.MaxAttempts
	}
	r.BackoffRate = raw.BackoffRate
	r.MaxDelaySeconds = raw.MaxDelaySeconds
	r.JitterStrategy = raw.JitterStrategy
	return nil
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
	Variable string        `json:"Variable,omitempty"`
	Next     string        `json:"Next,omitempty"`
	And      []*ChoiceRule `json:"And,omitempty"`
	Or       []*ChoiceRule `json:"Or,omitempty"`
	Not      *ChoiceRule   `json:"Not,omitempty"`
	// The string- and timestamp-valued operators are pointers so a literal
	// empty string stays distinguishable from an absent operator.
	StringEquals                   *string                `json:"StringEquals,omitempty"`
	StringLessThan                 *string                `json:"StringLessThan,omitempty"`
	StringGreaterThan              *string                `json:"StringGreaterThan,omitempty"`
	StringLessThanEquals           *string                `json:"StringLessThanEquals,omitempty"`
	StringGreaterThanEquals        *string                `json:"StringGreaterThanEquals,omitempty"`
	StringMatches                  *string                `json:"StringMatches,omitempty"`
	NumericEquals                  *float64               `json:"NumericEquals,omitempty"`
	NumericLessThan                *float64               `json:"NumericLessThan,omitempty"`
	NumericGreaterThan             *float64               `json:"NumericGreaterThan,omitempty"`
	NumericLessThanEquals          *float64               `json:"NumericLessThanEquals,omitempty"`
	NumericGreaterThanEquals       *float64               `json:"NumericGreaterThanEquals,omitempty"`
	BooleanEquals                  *bool                  `json:"BooleanEquals,omitempty"`
	TimestampEquals                *string                `json:"TimestampEquals,omitempty"`
	TimestampLessThan              *string                `json:"TimestampLessThan,omitempty"`
	TimestampGreaterThan           *string                `json:"TimestampGreaterThan,omitempty"`
	TimestampLessThanEquals        *string                `json:"TimestampLessThanEquals,omitempty"`
	TimestampGreaterThanEquals     *string                `json:"TimestampGreaterThanEquals,omitempty"`
	IsPresent                      *bool                  `json:"IsPresent,omitempty"`
	IsNull                         *bool                  `json:"IsNull,omitempty"`
	IsBoolean                      *bool                  `json:"IsBoolean,omitempty"`
	IsString                       *bool                  `json:"IsString,omitempty"`
	IsNumeric                      *bool                  `json:"IsNumeric,omitempty"`
	IsTimestamp                    *bool                  `json:"IsTimestamp,omitempty"`
	StringEqualsPath               string                 `json:"StringEqualsPath,omitempty"`
	StringLessThanPath             string                 `json:"StringLessThanPath,omitempty"`
	StringGreaterThanPath          string                 `json:"StringGreaterThanPath,omitempty"`
	StringLessThanEqualsPath       string                 `json:"StringLessThanEqualsPath,omitempty"`
	StringGreaterThanEqualsPath    string                 `json:"StringGreaterThanEqualsPath,omitempty"`
	NumericEqualsPath              string                 `json:"NumericEqualsPath,omitempty"`
	NumericLessThanPath            string                 `json:"NumericLessThanPath,omitempty"`
	NumericGreaterThanPath         string                 `json:"NumericGreaterThanPath,omitempty"`
	NumericLessThanEqualsPath      string                 `json:"NumericLessThanEqualsPath,omitempty"`
	NumericGreaterThanEqualsPath   string                 `json:"NumericGreaterThanEqualsPath,omitempty"`
	BooleanEqualsPath              string                 `json:"BooleanEqualsPath,omitempty"`
	TimestampEqualsPath            string                 `json:"TimestampEqualsPath,omitempty"`
	TimestampLessThanPath          string                 `json:"TimestampLessThanPath,omitempty"`
	TimestampGreaterThanPath       string                 `json:"TimestampGreaterThanPath,omitempty"`
	TimestampLessThanEqualsPath    string                 `json:"TimestampLessThanEqualsPath,omitempty"`
	TimestampGreaterThanEqualsPath string                 `json:"TimestampGreaterThanEqualsPath,omitempty"`
	Condition                      string                 `json:"Condition,omitempty"`
	Assign                         map[string]interface{} `json:"Assign,omitempty"`
}

// UnmarshalJSON deserialises a ChoiceRule from JSON. Each comparison
// operator carries its comparison value directly ("the value may be a
// string, number, boolean, or timestamp") — no variable-keyed map form
// exists in the States Language. A literal whose JSON type does not match
// its operator decodes as the absent operator: "For each operator which
// compares values, if the values are not both of the appropriate type
// (String, number, boolean, or Timestamp) the comparison will return
// false", so the rule never matches rather than rejecting the definition.
// The *Path operators ("the value MUST be a Path"), StringMatches ("The
// value MUST be a String") and the Is* booleans keep strict typing — the
// spec attaches explicit MUSTs to those values.
func (r *ChoiceRule) UnmarshalJSON(data []byte) error {
	type Alias ChoiceRule
	aux := &struct {
		*Alias
		StringEqualsRaw               interface{} `json:"StringEquals,omitempty"`
		StringLessThanRaw             interface{} `json:"StringLessThan,omitempty"`
		StringGreaterThanRaw          interface{} `json:"StringGreaterThan,omitempty"`
		StringLessThanEqualsRaw       interface{} `json:"StringLessThanEquals,omitempty"`
		StringGreaterThanEqualsRaw    interface{} `json:"StringGreaterThanEquals,omitempty"`
		NumericEqualsRaw              interface{} `json:"NumericEquals,omitempty"`
		NumericLessThanRaw            interface{} `json:"NumericLessThan,omitempty"`
		NumericGreaterThanRaw         interface{} `json:"NumericGreaterThan,omitempty"`
		NumericLessThanEqualsRaw      interface{} `json:"NumericLessThanEquals,omitempty"`
		NumericGreaterThanEqualsRaw   interface{} `json:"NumericGreaterThanEquals,omitempty"`
		BooleanEqualsRaw              interface{} `json:"BooleanEquals,omitempty"`
		TimestampEqualsRaw            interface{} `json:"TimestampEquals,omitempty"`
		TimestampLessThanRaw          interface{} `json:"TimestampLessThan,omitempty"`
		TimestampGreaterThanRaw       interface{} `json:"TimestampGreaterThan,omitempty"`
		TimestampLessThanEqualsRaw    interface{} `json:"TimestampLessThanEquals,omitempty"`
		TimestampGreaterThanEqualsRaw interface{} `json:"TimestampGreaterThanEquals,omitempty"`
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.StringEquals = stringOperatorValue(aux.StringEqualsRaw)
	r.StringLessThan = stringOperatorValue(aux.StringLessThanRaw)
	r.StringGreaterThan = stringOperatorValue(aux.StringGreaterThanRaw)
	r.StringLessThanEquals = stringOperatorValue(aux.StringLessThanEqualsRaw)
	r.StringGreaterThanEquals = stringOperatorValue(aux.StringGreaterThanEqualsRaw)
	r.NumericEquals = floatOperatorValue(aux.NumericEqualsRaw)
	r.NumericLessThan = floatOperatorValue(aux.NumericLessThanRaw)
	r.NumericGreaterThan = floatOperatorValue(aux.NumericGreaterThanRaw)
	r.NumericLessThanEquals = floatOperatorValue(aux.NumericLessThanEqualsRaw)
	r.NumericGreaterThanEquals = floatOperatorValue(aux.NumericGreaterThanEqualsRaw)
	r.BooleanEquals = boolOperatorValue(aux.BooleanEqualsRaw)
	r.TimestampEquals = stringOperatorValue(aux.TimestampEqualsRaw)
	r.TimestampLessThan = stringOperatorValue(aux.TimestampLessThanRaw)
	r.TimestampGreaterThan = stringOperatorValue(aux.TimestampGreaterThanRaw)
	r.TimestampLessThanEquals = stringOperatorValue(aux.TimestampLessThanEqualsRaw)
	r.TimestampGreaterThanEquals = stringOperatorValue(aux.TimestampGreaterThanEqualsRaw)
	return nil
}

// stringOperatorValue decodes one string-valued choice operator's
// comparison value, returning nil for an absent operator so a literal
// empty string stays a usable comparison operand — and for a value whose
// JSON type is not a string, the mistyped-literal representation whose
// comparison "will return false".
func stringOperatorValue(v interface{}) *string {
	if s, ok := v.(string); ok {
		return &s
	}
	return nil
}

// floatOperatorValue decodes one numeric choice operator's comparison
// value; the pointer distinguishes an absent operator from a zero value.
func floatOperatorValue(v interface{}) *float64 {
	switch val := v.(type) {
	case float64:
		return &val
	case json.Number:
		if f, err := val.Float64(); err == nil {
			return &f
		}
	}
	return nil
}

// boolOperatorValue decodes one boolean choice operator's comparison
// value; the pointer distinguishes an absent operator from a false value.
func boolOperatorValue(v interface{}) *bool {
	if b, ok := v.(bool); ok {
		return &b
	}
	return nil
}

// InputOutput represents input and output configuration for a state, supporting
// JSONPath, Parameters, and payload template values.
type InputOutput struct {
	Path       string      `json:"Path,omitempty"`
	Parameters *Parameters `json:"Parameters,omitempty"`
	InputPath  string      `json:"InputPath,omitempty"`
	ItemsPath  string      `json:"ItemsPath,omitempty"`
}

// UnmarshalJSON deserialises an InputOutput from JSON, accepting either a raw string
// (payload template) or a structured object. The payload template is stored
// on InputPath because the SL spec uses InputPath as the canonical field for
// inline string payloads when no other structured fields are present.
func (io *InputOutput) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		io.InputPath = s
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
