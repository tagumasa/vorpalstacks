package sfn

import (
	"time"

	tagutil "vorpalstacks/internal/common/tags"
)

// AWS specification limits for Step Functions. These are the single source
// of truth; service handlers, validators and tests must reference them
// instead of inlining numeric literals.
const (
	// MaxResourceNameLength is the general quota for state machine,
	// execution, activity and alias names (Smithy Name shape @length(1, 80)).
	MaxResourceNameLength = 80

	// DefaultPageSize is the page size applied when a List operation
	// request omits maxResults.
	DefaultPageSize = 100

	// MaxVariableNameLength is the documented maximum length of a JSONata
	// state-variable name.
	MaxVariableNameLength = 80

	// MaxPerVariableBytes is the documented maximum size of a single
	// JSONata state-variable value.
	MaxPerVariableBytes = 256 * 1024

	// MaxAssignTotalBytes is the documented maximum combined size of one
	// Assign block's variable values.
	MaxAssignTotalBytes = 256 * 1024

	// MaxTotalVariableBytes is the documented maximum combined size of all
	// JSONata state variables in a state machine.
	MaxTotalVariableBytes = 10 * 1024 * 1024

	// MaxPageSize is the maximum for maxResults on List operations. The
	// 0-1000 bound is stated in the PageSize-typed members' documentation
	// trait; the PageShape itself carries no @range.
	MaxPageSize = 1000

	// MaxValidateDefinitionResults is the Smithy @range(0, 100) maximum for
	// maxResults on ValidateStateMachineDefinition. The diagnostics paging
	// shape is distinct from the List-operation PageSize shape and does
	// not share its bound.
	MaxValidateDefinitionResults = 100

	// MaxWaitSeconds is the Wait-state ceiling for Seconds (and for the
	// value a JSONata Seconds expression evaluates to): "You must specify
	// time as an integer value from 0 to 99,999,999" (AWS Wait state
	// documentation).
	MaxWaitSeconds = 99999999

	// MaxTestStateRunSeconds is the TestState run ceiling: "The TestState
	// API can run for up to five minutes. If the execution of a state
	// exceeds this duration, it fails with the States.Timeout error."
	MaxTestStateRunSeconds = 300

	// MaxTagsPerResource is the hard tagging quota: a maximum of fifty tags
	// per resource, not modifiable.
	MaxTagsPerResource = tagutil.MaxTagsPerResource

	// MaxArnLength is the Smithy Arn shape ceiling @length(1, 256) shared by
	// stateMachineArn, roleArn, executionArn and their qualified variants.
	MaxArnLength = 256

	// MaxExecutionDataBytes is the Smithy SensitiveData shape ceiling
	// @length(0, 262144): the UTF-8 byte bound for StartExecution input,
	// StartSyncExecution input and SendTaskSuccess output. It is also the
	// payload size limit the intrinsic-functions documentation places on
	// array arguments to States.ArrayPartition and States.ArrayContains.
	MaxExecutionDataBytes = 262144

	// MaxIntrinsicStringChars is the intrinsic-functions documentation bound
	// on the data string argument of States.Base64Encode,
	// States.Base64Decode and States.Hash: 10,000 characters.
	MaxIntrinsicStringChars = 10000

	// MaxIntrinsicNesting is the intrinsic-functions documentation bound on
	// nesting: "You can nest up to 10 intrinsic functions within a field".
	MaxIntrinsicNesting = 10

	// MaxArrayRangeElements is the intrinsic-functions documentation bound on
	// States.ArrayRange: "The new array can contain up to 1000 elements".
	MaxArrayRangeElements = 1000

	// MaxErrorLength is the Smithy SensitiveError shape ceiling
	// @length(0, 256): the bound for StopExecution and SendTaskFailure
	// error strings.
	MaxErrorLength = 256

	// MaxCauseLength is the Smithy SensitiveCause shape ceiling
	// @length(0, 32768): the bound for StopExecution and SendTaskFailure
	// cause strings.
	MaxCauseLength = 32768

	// MaxDefinitionLength is the Smithy Definition shape ceiling
	// @length(1, 1048576): the bound for CreateStateMachine and
	// UpdateStateMachine definitions.
	MaxDefinitionLength = 1048576

	// MaxTaskTokenLength is the Smithy TaskToken shape ceiling
	// @length(1, 2048).
	MaxTaskTokenLength = 2048

	// MaxTraceHeaderLength is the Smithy TraceHeader shape ceiling
	// @length(0, 256) with the @pattern(^\\p{ASCII}*$) ASCII-only profile.
	MaxTraceHeaderLength = 256

	// MaxVersionDescriptionLength is the versionDescription ceiling
	// @length(0, 256) on PublishStateMachineVersion and UpdateStateMachine.
	MaxVersionDescriptionLength = 256

	// MaxAliasesPerStateMachine is the alias quota per state machine
	// (AWS documentation on state machine aliases).
	MaxAliasesPerStateMachine = 100

	// MaxVersionsPerStateMachine is the version quota per state machine
	// (PublishStateMachineVersion documentation).
	MaxVersionsPerStateMachine = 1000

	// MaxRoutingConfigEntries is the number of versions an alias routing
	// configuration may point at (AWS documentation on state machine
	// aliases).
	MaxRoutingConfigEntries = 2

	// RedriveWindowDays bounds the redrivable period: executions may be
	// redriven within fourteen days of completing (RedriveExecution
	// documentation).
	RedriveWindowDays = 14

	// MaxRedriveEventHistory is the event-history ceiling for redrive
	// eligibility: the execution history must hold fewer than 24,999 events
	// to accommodate the ExecutionRedriven event and at least one more
	// event (RedriveExecution documentation).
	MaxRedriveEventHistory = 24999

	// MaxItemReaderItems is the ItemReader MaxItems ceiling: "You can
	// specify a limit of up to 100,000,000 after which the Distributed Map
	// stops reading items" (ItemReader documentation).
	MaxItemReaderItems = 100000000

	// MaxCSVHeaderBytes is the CSV header size ceiling for text delimited
	// ItemReader datasets: "Step Functions supports headers of up to 10 KiB
	// for text delimited files" (ItemReader documentation).
	MaxCSVHeaderBytes = 10240

	// MaxMapLabelLength is the Distributed Map Label field ceiling: labels
	// "can't exceed 40 characters in length" (Distributed Map state
	// documentation).
	MaxMapLabelLength = 40

	// MaxItemReaderFileBytes is the ItemReader single-file ceiling:
	// "Step Functions supports 10 GB as the maximum size of an individual
	// file in S3" (ItemReader documentation).
	MaxItemReaderFileBytes = 10 * 1024 * 1024 * 1024

	// MaxStateNameLength is the state-name ceiling from the Amazon States
	// Language specification: "its length MUST BE less than or equal to 80
	// Unicode characters".
	MaxStateNameLength = 80

	// MaxRetryIntervalSeconds is the Retrier IntervalSeconds ceiling:
	// "IntervalSeconds has a maximum value of 99,999,999" (error handling
	// documentation).
	MaxRetryIntervalSeconds = 99999999

	// MaxRetryAttempts is the Retrier MaxAttempts ceiling: "MaxAttempts
	// has a maximum value of 99,999,999" (error handling documentation).
	MaxRetryAttempts = 99999999

	// MaxRetryDelaySeconds is the exclusive Retrier MaxDelaySeconds upper
	// bound: "You must specify a value greater than 0 and less than
	// 31622401 for MaxDelaySeconds" (error handling documentation).
	MaxRetryDelaySeconds = 31622401
)

// StateMachine represents an AWS Step Functions state machine definition and metadata.
type StateMachine struct {
	StateMachineArn         string                   `json:"stateMachineArn"`
	Name                    string                   `json:"name"`
	Definition              string                   `json:"definition"`
	RoleArn                 string                   `json:"roleArn"`
	Type                    string                   `json:"type"`
	CreationDate            time.Time                `json:"creationDate"`
	UpdateDate              time.Time                `json:"updateDate"`
	Status                  string                   `json:"status"`
	VariableReferences      map[string][]string      `json:"-"`
	LoggingConfiguration    *LoggingConfiguration    `json:"loggingConfiguration,omitempty"`
	EncryptionConfiguration *EncryptionConfiguration `json:"encryptionConfiguration,omitempty"`
	TracingConfiguration    *TracingConfiguration    `json:"tracingConfiguration,omitempty"`
	RevisionId              string                   `json:"revisionId,omitempty"`
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

// EncryptionConfiguration specifies KMS key for encrypting state machine data.
type EncryptionConfiguration struct {
	Type                  string `json:"type,omitempty"`
	KmsKeyId              string `json:"kmsKeyId,omitempty"`
	KmsDataKeyReusePeriod int32  `json:"kmsDataKeyReusePeriodSeconds,omitempty"`
}

// TracingConfiguration controls X-Ray tracing on state machine executions.
type TracingConfiguration struct {
	Enabled bool `json:"enabled"`
}

// Execution represents a single execution of a Step Functions state machine.
type Execution struct {
	ExecutionArn           string                         `json:"executionArn"`
	StateMachineArn        string                         `json:"stateMachineArn"`
	StateMachineVersionArn string                         `json:"stateMachineVersionArn,omitempty"`
	StateMachineAliasArn   string                         `json:"stateMachineAliasArn,omitempty"`
	Name                   string                         `json:"name"`
	Status                 string                         `json:"status"`
	Input                  string                         `json:"input"`
	Output                 string                         `json:"output"`
	TraceHeader            string                         `json:"traceHeader"`
	StartDate              time.Time                      `json:"startDate"`
	StopDate               time.Time                      `json:"stopDate,omitempty"`
	Error                  string                         `json:"error,omitempty"`
	Cause                  string                         `json:"cause,omitempty"`
	MapRunArn              string                         `json:"mapRunArn,omitempty"`
	ItemCount              int64                          `json:"itemCount,omitempty"`
	RedriveCount           int64                          `json:"redriveCount,omitempty"`
	RedriveDate            time.Time                      `json:"redriveDate,omitempty"`
	ParallelCheckpoints    map[string]*ParallelCheckpoint `json:"parallelCheckpoints,omitempty"`
}

// ParallelCheckpoint stores the per-branch results of a Parallel state so that
// a redriven execution can skip already-completed branches and re-run only the
// failed ones.
type ParallelCheckpoint struct {
	BranchResults map[int]string `json:"branchResults"`
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
