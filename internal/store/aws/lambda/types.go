// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"

	types "vorpalstacks/internal/common/tags"
)

// Runtime represents the Lambda function runtime environment.
type Runtime string

// Runtime constants define supported Lambda runtime environments.
const (
	RuntimeNodejs24X      Runtime = "nodejs24.x"
	RuntimeNodejs22X      Runtime = "nodejs22.x"
	RuntimePython314      Runtime = "python3.14"
	RuntimePython313      Runtime = "python3.13"
	RuntimePython312      Runtime = "python3.12"
	RuntimePython311      Runtime = "python3.11"
	RuntimePython310      Runtime = "python3.10"
	RuntimeJava25         Runtime = "java25"
	RuntimeJava21         Runtime = "java21"
	RuntimeJava17         Runtime = "java17"
	RuntimeJava11         Runtime = "java11"
	RuntimeJava17Al2023   Runtime = "java17.al2023"
	RuntimeJava11Al2023   Runtime = "java11.al2023"
	RuntimeJava8Al2023    Runtime = "java8.al2023"
	RuntimeJava8Al2       Runtime = "java8.al2"
	RuntimeDotnet10       Runtime = "dotnet10"
	RuntimeDotnet8        Runtime = "dotnet8"
	RuntimeRuby40         Runtime = "ruby4.0"
	RuntimeRuby34         Runtime = "ruby3.4"
	RuntimeRuby33         Runtime = "ruby3.3"
	RuntimeProvidedAl2023 Runtime = "provided.al2023"
	RuntimeProvidedAl2    Runtime = "provided.al2"

	// Deprecated runtimes — kept for existing functions only.
	RuntimeNodejs20X Runtime = "nodejs20.x"
	RuntimeNodejs18X Runtime = "nodejs18.x"
	RuntimeNodejs16X Runtime = "nodejs16.x"
	RuntimePython39  Runtime = "python3.9"
	RuntimePython38  Runtime = "python3.8"
	RuntimeDotnet6   Runtime = "dotnet6"
	RuntimeRuby32    Runtime = "ruby3.2"
	RuntimeGo1X      Runtime = "go1.x"
)

// CurrentRuntimes is the single source for the runtime identifiers
// accepted for new functions and layer compatible-runtime declarations,
// derived from the vendored Smithy model's Runtime enum: the
// feature:public values of the currently-supported generation. The
// deprecated constants above remain invocable for existing functions but
// are not in this set; preview values the model gates behind non-public
// feature tags (nodejs26.x, python3.15) are not part of the public enum
// surface and are excluded.
var CurrentRuntimes = []Runtime{
	RuntimeNodejs24X, RuntimeNodejs22X,
	RuntimePython314, RuntimePython313, RuntimePython312, RuntimePython311, RuntimePython310,
	RuntimeJava25, RuntimeJava21, RuntimeJava17, RuntimeJava11,
	RuntimeJava17Al2023, RuntimeJava11Al2023, RuntimeJava8Al2023, RuntimeJava8Al2,
	RuntimeDotnet10, RuntimeDotnet8,
	RuntimeRuby40, RuntimeRuby34, RuntimeRuby33,
	RuntimeProvidedAl2023, RuntimeProvidedAl2,
}

// AWS defaults for asynchronous invocation configuration, applied when
// PutFunctionEventInvokeConfig omits the corresponding member: "By default,
// Lambda retries an asynchronous invocation twice if the function returns
// an error. It retains events in a queue for up to six hours."
const (
	DefaultMaximumRetryAttempts     = int32(2)
	DefaultMaximumEventAgeInSeconds = int32(21600)
)

// Function configuration defaults applied when CreateFunction omits the
// member. The invocation plane also relies on the timeout default when it
// meets a stored configuration carrying a zero, so the exec deadline and
// the create default can never diverge.
const (
	DefaultFunctionTimeoutSeconds = int32(3)
	DefaultFunctionMemorySizeMB   = int32(128)
)

// Function configuration member ranges from the CreateFunction and
// UpdateFunctionConfiguration model @range traits: EphemeralStorageSize
// min 512 max 32768 MB, Timeout min 1 max 5400 seconds, MemorySize min
// 128 max 32768 MB.
const (
	MinEphemeralStorageSizeMB = int32(512)
	MaxEphemeralStorageSizeMB = int32(32768)
	MinTimeoutSeconds         = int32(1)
	MaxTimeoutSeconds         = int32(5400)
	MinMemorySizeMB           = int32(128)
	MaxMemorySizeMB           = int32(32768)
)

// DefaultLayerVersionListMaxItems is the ListLayerVersions documented
// default page size, applied when the request omits MaxItems.
const DefaultLayerVersionListMaxItems = 50

// MaxResourcePolicyLength is the maximum length of a resource-based
// policy document, from the ResourcePolicy shape @length trait
// (min 1, max 20480 characters) shared by PutResourcePolicy and
// GetResourcePolicy.
const MaxResourcePolicyLength = 20480

// MaxPolicyResourceArnLength is the maximum length of the ResourceArn
// member of the resource-policy operations, from the PolicyResourceArn
// shape @length trait (min 0, max 256 characters).
const MaxPolicyResourceArnLength = 256

// Sandbox pool bounds for image-package functions. The idle TTL follows the
// AWS model of terminating an inactive execution environment after a period
// (AWS publishes no fixed value; 300 seconds is this platform's policy,
// reaping every idle sandbox including the last one). The global limit caps
// concurrent sandboxes across functions without reserved concurrency —
// functions with reserved concurrency draw against their own reserve
// instead, mirroring the reserved/unreserved concurrency split.
const (
	DefaultSandboxIdleTTL     = 300 * time.Second
	DefaultGlobalSandboxLimit = 128
)

// Account-level limits reported by GetAccountSettings. UnreservedConcurrentExecutions
// is derived per region as the concurrency limit minus the sum of reserved
// concurrency across the region's functions.
const (
	AccountLimitTotalCodeSize        = int64(80530636800)
	AccountLimitCodeSizeUnzipped     = int64(262144000)
	AccountLimitCodeSizeZipped       = int64(52428800)
	AccountLimitConcurrentExecutions = int64(1000)
)

// Environment variable rules: "Environment variables beginning with
// AWS_LAMBDA_ are reserved" and "The total size of all environment
// variables doesn't exceed 4 KB" (AWS Lambda configuration documentation).
// Keys follow "Keys start with a letter and are at least two characters.
// Keys only contain letters, numbers, and the underscore character (_)".
const (
	MaxEnvironmentVariablesSizeBytes = 4096
	ReservedEnvironmentPrefix        = "AWS_LAMBDA_"
	// EnvironmentVariableKeyPattern is the documented key shape: an
	// initial letter, then letters, digits and underscores, minimum two
	// characters.
	EnvironmentVariableKeyPattern = `^[a-zA-Z][a-zA-Z0-9_]+$`
)

// Event source mapping parallelization factor range: "Valid Range:
// Minimum value of 1. Maximum value of 10." (CreateEventSourceMapping
// model).
const (
	MinParallelizationFactor = int32(1)
	MaxParallelizationFactor = int32(10)
)

// Event source mapping batch size range from the CreateEventSourceMapping
// model: "Valid Range: Minimum value of 1. Maximum value of 10000." The
// per-source defaults (100 for stream sources, 10 for queues) are
// service-plane policy applied by defaultESMBatchSize in the service
// package.
const (
	MinESMBatchSize = int32(1)
	MaxESMBatchSize = int32(10000)
)

// Event source mapping member ranges from the 2015-03-31 model shapes:
// MaximumBatchingWindowInSeconds @range(0, 300), MaximumRecordAgeInSeconds
// @range(-1, 604800), MaximumRetryAttempts @range(-1, 10000), and
// TumblingWindowInSeconds @range(0, 900).
const (
	MinESMBatchingWindowSeconds = int32(0)
	MaxESMBatchingWindowSeconds = int32(300)
	MinESMRecordAgeSeconds      = int32(-1)
	MaxESMRecordAgeSeconds      = int32(604800)
	MinESMRetryAttempts         = int32(-1)
	MaxESMRetryAttempts         = int32(10000)
	MinESMTumblingWindowSeconds = int32(0)
	MaxESMTumblingWindowSeconds = int32(900)
)

// Function event invoke config member ranges from the 2015-03-31 model
// shapes: MaximumEventAgeInSeconds @range(60, 21600) and
// MaximumRetryAttempts @range(0, 2).
const (
	MinEventAgeSeconds          = int32(60)
	MaxEventAgeSeconds          = int32(21600)
	MinEventInvokeRetryAttempts = int32(0)
	MaxEventInvokeRetryAttempts = int32(2)
)

// MaxStatementIdLength is the StatementId shape's @length(1, 100)
// maximum from the 2015-03-31 model.
const MaxStatementIdLength = 100

// LoggingConfig member constraints from the 2015-03-31 model: the LogGroup
// shape carries length 1-512 and the pattern below; the enum members
// (LogFormat, ApplicationLogLevel, SystemLogLevel) are enforced by the
// service-layer validators.
const (
	MaxLogGroupLength = 512
	LogGroupPattern   = `^[\.\-_/#A-Za-z0-9]+$`
)

// ImageConfig member constraints from the 2015-03-31 model: EntryPoint and
// Command target the StringList shape (length 0-1500) and WorkingDirectory
// carries length 0-1000.
const (
	MaxImageConfigListLength       = 1500
	MaxImageConfigWorkingDirectory = 1000
)

// Function URL CORS constraints from the 2015-03-31 model: the Cors.MaxAge
// member targets the MaxAge shape whose @range is 0 to 86400 seconds.
const (
	MaxCorsMaxAgeSeconds = int32(86400)
)

// State represents the current state of a Lambda function.
type State string

// State constants represent possible Lambda function states.
const (
	StatePending  State = "Pending"
	StateActive   State = "Active"
	StateInactive State = "Inactive"
	StateFailed   State = "Failed"
)

// LastUpdateStatus represents the status of the last update operation on a Lambda function.
type LastUpdateStatus string

// LastUpdateStatus constants represent possible update statuses.
const (
	LastUpdateStatusSuccessful LastUpdateStatus = "Successful"
	LastUpdateStatusFailed     LastUpdateStatus = "Failed"
	LastUpdateStatusInProgress LastUpdateStatus = "InProgress"
)

// Function represents a Lambda function configuration and metadata.
type Function struct {
	FunctionName               string                 `json:"function_name"`
	FunctionArn                string                 `json:"function_arn"`
	Runtime                    Runtime                `json:"runtime"`
	Role                       string                 `json:"role"`
	Handler                    string                 `json:"handler"`
	CodeSize                   int64                  `json:"code_size"`
	CodeSha256                 string                 `json:"code_sha256"`
	CodeLocation               string                 `json:"code_location"`
	ImageUri                   string                 `json:"image_uri,omitempty"`
	SourceCodeHash             string                 `json:"source_code_hash,omitempty"`
	Description                string                 `json:"description,omitempty"`
	Timeout                    int32                  `json:"timeout"`
	MemorySize                 int32                  `json:"memory_size"`
	EphemeralStorage           *EphemeralStorage      `json:"ephemeral_storage,omitempty"`
	Publish                    bool                   `json:"publish"`
	Architectures              []string               `json:"architectures,omitempty"`
	KMSKeyArn                  string                 `json:"kms_key_arn,omitempty"`
	RevisionId                 string                 `json:"revision_id"`
	State                      State                  `json:"state"`
	StateReason                string                 `json:"state_reason,omitempty"`
	StateReasonCode            string                 `json:"state_reason_code,omitempty"`
	LastUpdateStatus           LastUpdateStatus       `json:"last_update_status"`
	LastUpdateReason           string                 `json:"last_update_reason,omitempty"`
	LastUpdateStatusReason     string                 `json:"last_update_status_reason,omitempty"`
	LastUpdateStatusReasonCode string                 `json:"last_update_status_reason_code,omitempty"`
	LastModified               time.Time              `json:"last_modified"`
	LastModifiedUser           string                 `json:"last_modified_user,omitempty"`
	VpcConfig                  *VpcConfig             `json:"vpc_config,omitempty"`
	Environment                *Environment           `json:"environment,omitempty"`
	DeadLetterConfig           *DeadLetterConfig      `json:"dead_letter_config,omitempty"`
	TracingConfig              *TracingConfig         `json:"tracing_config,omitempty"`
	Layers                     []LayerReference       `json:"layers,omitempty"`
	Tags                       []types.Tag            `json:"tags,omitempty"`
	SnapStart                  *SnapStart             `json:"snap_start,omitempty"`
	PackageType                string                 `json:"package_type"`
	SigningProfileVersionArn   string                 `json:"signing_profile_version_arn,omitempty"`
	SigningJobArn              string                 `json:"signing_job_arn,omitempty"`
	UrlConfig                  *FunctionUrlConfig     `json:"url_config,omitempty"`
	CodeSigningConfigArn       string                 `json:"code_signing_config_arn,omitempty"`
	LoggingConfig              *LoggingConfig         `json:"logging_config,omitempty"`
	ImageConfig                *ImageConfig           `json:"image_config,omitempty"`
	FileSystemConfigs          []FileSystemConfig     `json:"file_system_configs,omitempty"`
	TenancyConfig              *TenancyConfig         `json:"tenancy_config,omitempty"`
	CapacityProviderConfig     map[string]interface{} `json:"capacity_provider_config,omitempty"`
	RuntimeVersionConfig       *RuntimeVersionConfig  `json:"runtime_version_config,omitempty"`
	DurableConfig              map[string]interface{} `json:"durable_config,omitempty"`

	Versions []Version        `json:"versions,omitempty"`
	Aliases  []Alias          `json:"aliases,omitempty"`
	Policies []FunctionPolicy `json:"policies,omitempty"`
	// PolicyRevisionId versions the resource-based policy independently of
	// the function revision: it is regenerated whenever the policy changes
	// (PutResourcePolicy, DeleteResourcePolicy, AddPermission,
	// RemovePermission) and compared against the RevisionId member of the
	// resource-policy operations.
	PolicyRevisionId string `json:"policy_revision_id,omitempty"`

	ReservedConcurrency    *int64                         `json:"reserved_concurrency,omitempty"`
	ProvisionedConcurrency []ProvisionedConcurrencyConfig `json:"provisioned_concurrency,omitempty"`
	EventInvokeConfigs     []EventInvokeConfig            `json:"event_invoke_configs,omitempty"`

	ContainerID      string `json:"container_id,omitempty"`
	ContainerImageID string `json:"container_image_id,omitempty"`

	latestVersionNum int                 `json:"-"`
	versionsByNum    map[string]*Version `json:"-"`
	aliasesByName    map[string]*Alias   `json:"-"`
}

// Version represents a published version of a Lambda function.
type Version struct {
	Version                  string             `json:"version"`
	FunctionArn              string             `json:"function_arn"`
	Runtime                  Runtime            `json:"runtime"`
	Role                     string             `json:"role"`
	Handler                  string             `json:"handler"`
	CodeSize                 int64              `json:"code_size"`
	CodeSha256               string             `json:"code_sha256"`
	CodeLocation             string             `json:"code_location"`
	ImageUri                 string             `json:"image_uri,omitempty"`
	Description              string             `json:"description,omitempty"`
	Timeout                  int32              `json:"timeout"`
	MemorySize               int32              `json:"memory_size"`
	EphemeralStorage         *EphemeralStorage  `json:"ephemeral_storage,omitempty"`
	Architectures            []string           `json:"architectures,omitempty"`
	KMSKeyArn                string             `json:"kms_key_arn,omitempty"`
	RevisionId               string             `json:"revision_id"`
	State                    State              `json:"state"`
	StateReason              string             `json:"state_reason,omitempty"`
	StateReasonCode          string             `json:"state_reason_code,omitempty"`
	LastUpdateStatus         LastUpdateStatus   `json:"last_update_status"`
	LastModified             time.Time          `json:"last_modified"`
	VpcConfig                *VpcConfig         `json:"vpc_config,omitempty"`
	Environment              *Environment       `json:"environment,omitempty"`
	DeadLetterConfig         *DeadLetterConfig  `json:"dead_letter_config,omitempty"`
	TracingConfig            *TracingConfig     `json:"tracing_config,omitempty"`
	Layers                   []LayerReference   `json:"layers,omitempty"`
	SnapStart                *SnapStart         `json:"snap_start,omitempty"`
	PackageType              string             `json:"package_type"`
	SigningProfileVersionArn string             `json:"signing_profile_version_arn,omitempty"`
	SigningJobArn            string             `json:"signing_job_arn,omitempty"`
	LoggingConfig            *LoggingConfig     `json:"logging_config,omitempty"`
	ImageConfig              *ImageConfig       `json:"image_config,omitempty"`
	FileSystemConfigs        []FileSystemConfig `json:"file_system_configs,omitempty"`

	ContainerID      string `json:"container_id,omitempty"`
	ContainerImageID string `json:"container_image_id,omitempty"`
}

// Alias represents an alias for a Lambda function.
type Alias struct {
	Name            string         `json:"name"`
	AliasArn        string         `json:"alias_arn"`
	FunctionVersion string         `json:"function_version"`
	Description     string         `json:"description,omitempty"`
	FunctionName    string         `json:"function_name"`
	RevisionId      string         `json:"revision_id"`
	RoutingConfig   *RoutingConfig `json:"routing_config,omitempty"`
}

// RoutingConfig defines the routing configuration for a Lambda function alias.
type RoutingConfig struct {
	AdditionalVersionWeights map[string]float64 `json:"additional_version_weights,omitempty"`
}

// LayerReference represents a reference to a Lambda layer.
type LayerReference struct {
	Arn      string `json:"arn"`
	CodeSize int64  `json:"code_size"`
}

// EphemeralStorage represents the ephemeral storage configuration for a Lambda function.
type EphemeralStorage struct {
	Size int32 `json:"size"`
}

// VpcConfig represents the VPC configuration for a Lambda function.
type VpcConfig struct {
	SubnetIds               []string `json:"subnet_ids,omitempty"`
	SecurityGroupIds        []string `json:"security_group_ids,omitempty"`
	VpcId                   string   `json:"vpc_id,omitempty"`
	Ipv6AllowedForDualStack bool     `json:"ipv6_allowed_for_dual_stack,omitempty"`
}

// Environment represents the environment variables for a Lambda function.
type Environment struct {
	Variables map[string]string `json:"variables,omitempty"`
	Error     *EnvironmentError `json:"error,omitempty"`
}

// EnvironmentError represents an error in the Lambda function environment.
type EnvironmentError struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// DeadLetterConfig represents the dead letter queue configuration for a Lambda function.
type DeadLetterConfig struct {
	TargetArn string `json:"target_arn,omitempty"`
}

// TracingConfig represents the tracing configuration for a Lambda function.
type TracingConfig struct {
	Mode string `json:"mode,omitempty"`
}

// SnapStart represents the SnapStart configuration for a Lambda function.
type SnapStart struct {
	ApplyOn            string `json:"apply_on,omitempty"`
	OptimizationStatus string `json:"optimization_status,omitempty"`
}

// LoggingConfig configures advanced Lambda logging settings.
type LoggingConfig struct {
	LogFormat           string `json:"log_format,omitempty"`
	ApplicationLogLevel string `json:"application_log_level,omitempty"`
	SystemLogLevel      string `json:"system_log_level,omitempty"`
	LogGroup            string `json:"log_group,omitempty"`
}

// ImageConfig configures container image entry point and command
// overrides for Lambda functions with PackageType=Image.
type ImageConfig struct {
	EntryPoint       []string `json:"entry_point,omitempty"`
	Command          []string `json:"command,omitempty"`
	WorkingDirectory string   `json:"working_directory,omitempty"`
}

// FileSystemConfig describes an EFS or S3 Files mount configuration.
type FileSystemConfig struct {
	Arn            string         `json:"arn,omitempty"`
	LocalMountPath string         `json:"local_mount_path,omitempty"`
	S3FilesConfig  *S3FilesConfig `json:"s3_files_config,omitempty"`
}

// S3FilesConfig controls how a function accesses data on an Amazon S3
// file system: direct bucket reads or reads through the file system.
type S3FilesConfig struct {
	DirectS3Read string `json:"direct_s3_read,omitempty"`
}

// DirectS3Read values (Smithy DirectS3Read enum): AUTO is the service
// default, ENABLED enforces direct bucket reads, DISABLED routes reads
// through the file system.
const (
	DirectS3ReadAuto     = "AUTO"
	DirectS3ReadEnabled  = "ENABLED"
	DirectS3ReadDisabled = "DISABLED"
)

// TenancyConfig configures the tenancy isolation mode for a function.
type TenancyConfig struct {
	TenantIsolationMode string `json:"tenant_isolation_mode,omitempty"`
}

// RuntimeVersionConfig records the runtime version assigned to a function.
type RuntimeVersionConfig struct {
	RuntimeVersionArn string                 `json:"runtime_version_arn,omitempty"`
	Error             map[string]interface{} `json:"error,omitempty"`
}

// FunctionUrlConfig represents the configuration for a Lambda function URL.
type FunctionUrlConfig struct {
	FunctionUrl      string      `json:"function_url"`
	FunctionArn      string      `json:"function_arn"`
	AuthType         string      `json:"auth_type"`
	Cors             *CorsConfig `json:"cors,omitempty"`
	InvokeMode       string      `json:"invoke_mode,omitempty"`
	CreationTime     time.Time   `json:"creation_time"`
	LastModifiedTime time.Time   `json:"last_modified_time"`
	Qualifier        string      `json:"qualifier,omitempty"`
}

// CorsConfig represents the CORS configuration for a Lambda function URL.
type CorsConfig struct {
	AllowCredentials bool     `json:"allow_credentials,omitempty"`
	AllowHeaders     []string `json:"allow_headers,omitempty"`
	AllowMethods     []string `json:"allow_methods,omitempty"`
	AllowOrigins     []string `json:"allow_origins,omitempty"`
	ExposeHeaders    []string `json:"expose_headers,omitempty"`
	MaxAge           int32    `json:"max_age,omitempty"`
}

// FunctionPolicy represents a resource-based policy for a Lambda function.
// A statement that arrived inside a PutResourcePolicy document also keeps
// its verbatim JSON in Raw so policy renderings reproduce the statement as
// submitted; statements added via AddPermission leave Raw empty and are
// rendered from the structured fields. Effect records the statement's
// effect ("Allow" or "Deny") so the eventbus dispatch authorisation can
// tell them apart. The Principal, Action and Resource fields carry the
// single-string form of those members for Allow-effect statements only:
// the matcher never evaluates a Deny statement as a grant, and an Allow
// statement whose members are not string-shaped cannot be evaluated by the
// matcher at all and grants nothing (its verbatim JSON still renders).
type FunctionPolicy struct {
	Id        string                 `json:"id"`
	Effect    string                 `json:"effect,omitempty"`
	Principal string                 `json:"principal"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Condition map[string]interface{} `json:"condition,omitempty"`
	Raw       string                 `json:"raw,omitempty"`
}

// Layer represents a Lambda layer configuration.
type Layer struct {
	LayerName             string          `json:"layer_name"`
	LayerArn              string          `json:"layer_arn"`
	LatestMatchingVersion *LayerVersion   `json:"latest_matching_version,omitempty"`
	Versions              []*LayerVersion `json:"versions,omitempty"`
	CreatedDate           time.Time       `json:"created_date"`
}

// LayerVersion represents a version of a Lambda layer.
type LayerVersion struct {
	Version                 int64         `json:"version"`
	LayerVersionArn         string        `json:"layer_version_arn"`
	RevisionId              string        `json:"revision_id,omitempty"`
	Description             string        `json:"description,omitempty"`
	CompatibleRuntimes      []Runtime     `json:"compatible_runtimes,omitempty"`
	LicenseInfo             string        `json:"license_info,omitempty"`
	CompatibleArchitectures []string      `json:"compatible_architectures,omitempty"`
	CreatedDate             time.Time     `json:"created_date"`
	CodeSize                int64         `json:"code_size"`
	CodeSha256              string        `json:"code_sha256"`
	CodeLocation            string        `json:"code_location"`
	SourceCodeHash          string        `json:"source_code_hash,omitempty"`
	Policies                []LayerPolicy `json:"policies,omitempty"`
}

// LayerPolicy represents a policy attached to a Lambda layer version.
type LayerPolicy struct {
	Id        string `json:"id"`
	Principal string `json:"principal"`
	Action    string `json:"action"`
}

// EventSourceMapping represents a mapping between an event source and a Lambda function.
type EventSourceMapping struct {
	UUID                           string                      `json:"uuid"`
	EventSourceMappingArn          string                      `json:"event_source_mapping_arn,omitempty"`
	BatchSize                      int32                       `json:"batch_size,omitempty"`
	MaximumBatchingWindowInSeconds int32                       `json:"maximum_batching_window_in_seconds,omitempty"`
	ParallelizationFactor          int32                       `json:"parallelization_factor,omitempty"`
	EventSourceArn                 string                      `json:"event_source_arn,omitempty"`
	FunctionArn                    string                      `json:"function_arn"`
	KMSKeyArn                      string                      `json:"kms_key_arn,omitempty"`
	LastModified                   time.Time                   `json:"last_modified"`
	LastProcessingResult           string                      `json:"last_processing_result,omitempty"`
	State                          string                      `json:"state"`
	StateTransitionReason          string                      `json:"state_transition_reason,omitempty"`
	DestinationConfig              *DestinationConfig          `json:"destination_config,omitempty"`
	Topics                         []string                    `json:"topics,omitempty"`
	Queues                         []string                    `json:"queues,omitempty"`
	SourceAccessConfigurations     []SourceAccessConfiguration `json:"source_access_configurations,omitempty"`
	SelfManagedEventSource         *SelfManagedEventSource     `json:"self_managed_event_source,omitempty"`
	MaximumRecordAgeInSeconds      int32                       `json:"maximum_record_age_in_seconds,omitempty"`
	BisectBatchOnFunctionError     bool                        `json:"bisect_batch_on_function_error,omitempty"`
	MaximumRetryAttempts           int32                       `json:"maximum_retry_attempts,omitempty"`
	TumblingWindowInSeconds        int32                       `json:"tumbling_window_in_seconds,omitempty"`
	FunctionResponseTypes          []string                    `json:"function_response_types,omitempty"`
	StartingPosition               string                      `json:"starting_position,omitempty"`
	StartingPositionTimestamp      time.Time                   `json:"starting_position_timestamp,omitempty"`
	FilterCriteria                 *FilterCriteria             `json:"filter_criteria,omitempty"`
	FunctionName                   string                      `json:"function_name"`
}

// DestinationConfig represents the destination configuration for asynchronous invocations.
type DestinationConfig struct {
	OnSuccess *OnSuccess `json:"on_success,omitempty"`
	OnFailure *OnFailure `json:"on_failure,omitempty"`
}

// OnSuccess represents the destination configuration for successful invocations.
type OnSuccess struct {
	Destination string `json:"destination,omitempty"`
}

// OnFailure represents the destination configuration for failed invocations.
type OnFailure struct {
	Destination string `json:"destination,omitempty"`
}

// SourceAccessConfiguration represents the configuration for accessing a self-managed event source.
type SourceAccessConfiguration struct {
	Type string `json:"type"`
	URI  string `json:"uri,omitempty"`
}

// SelfManagedEventSource represents a self-managed event source configuration.
type SelfManagedEventSource struct {
	Endpoints map[string][]string `json:"endpoints,omitempty"`
}

// FilterCriteria represents the filter criteria for Lambda event source mapping.
type FilterCriteria struct {
	Filters []Filter `json:"filters,omitempty"`
}

// Filter represents a single filter pattern for Lambda event source mapping.
type Filter struct {
	Pattern string `json:"pattern,omitempty"`
}

// ProvisionedConcurrencyConfig represents the provisioned concurrency configuration for a Lambda function.
type ProvisionedConcurrencyConfig struct {
	FunctionName                             string    `json:"function_name"`
	FunctionArn                              string    `json:"function_arn,omitempty"`
	Qualifier                                string    `json:"qualifier"`
	AllocatedProvisionedConcurrentExecutions int32     `json:"allocated_provisioned_concurrent_executions,omitempty"`
	AvailableProvisionedConcurrentExecutions int32     `json:"available_provisioned_concurrent_executions,omitempty"`
	RequestedProvisionedConcurrentExecutions int32     `json:"requested_provisioned_concurrent_executions,omitempty"`
	Status                                   string    `json:"status,omitempty"`
	StatusReason                             string    `json:"status_reason,omitempty"`
	LastModified                             time.Time `json:"last_modified"`
}

// EventInvokeConfig represents the event invoke configuration for a Lambda function.
type EventInvokeConfig struct {
	FunctionName             string             `json:"function_name"`
	Qualifier                string             `json:"qualifier"`
	LastModified             time.Time          `json:"last_modified"`
	DestinationConfig        *DestinationConfig `json:"destination_config,omitempty"`
	MaximumEventAgeInSeconds int32              `json:"maximum_event_age_in_seconds,omitempty"`
	MaximumRetryAttempts     int32              `json:"maximum_retry_attempts,omitempty"`
}

// InvocationResult represents the result of a Lambda function invocation.
type InvocationResult struct {
	StatusCode      int64  `json:"status_code"`
	ExecutedVersion string `json:"executed_version,omitempty"`
	Payload         []byte `json:"payload,omitempty"`
	FunctionError   string `json:"function_error,omitempty"`
	LogResult       string `json:"log_result,omitempty"`
}

// RuntimeImageMapping maps Lambda runtimes to their default container images.
var RuntimeImageMapping = map[Runtime]string{
	RuntimeNodejs24X:      "public.ecr.aws/lambda/nodejs:24",
	RuntimeNodejs22X:      "public.ecr.aws/lambda/nodejs:22",
	RuntimeNodejs20X:      "public.ecr.aws/lambda/nodejs:20",
	RuntimeNodejs18X:      "public.ecr.aws/lambda/nodejs:18",
	RuntimeNodejs16X:      "public.ecr.aws/lambda/nodejs:16",
	RuntimePython314:      "public.ecr.aws/lambda/python:3.14",
	RuntimePython313:      "public.ecr.aws/lambda/python:3.13",
	RuntimePython312:      "public.ecr.aws/lambda/python:3.12",
	RuntimePython311:      "public.ecr.aws/lambda/python:3.11",
	RuntimePython310:      "public.ecr.aws/lambda/python:3.10",
	RuntimePython39:       "public.ecr.aws/lambda/python:3.9",
	RuntimePython38:       "public.ecr.aws/lambda/python:3.8",
	RuntimeJava25:         "public.ecr.aws/lambda/java:25",
	RuntimeJava21:         "public.ecr.aws/lambda/java:21",
	RuntimeJava17:         "public.ecr.aws/lambda/java:17",
	RuntimeJava11:         "public.ecr.aws/lambda/java:11",
	RuntimeJava17Al2023:   "public.ecr.aws/lambda/java:17.al2023",
	RuntimeJava11Al2023:   "public.ecr.aws/lambda/java:11.al2023",
	RuntimeJava8Al2023:    "public.ecr.aws/lambda/java:8.al2023",
	RuntimeJava8Al2:       "public.ecr.aws/lambda/java:8.al2",
	RuntimeDotnet10:       "public.ecr.aws/lambda/dotnet:10",
	RuntimeDotnet8:        "public.ecr.aws/lambda/dotnet:8",
	RuntimeDotnet6:        "public.ecr.aws/lambda/dotnet:6",
	RuntimeRuby40:         "public.ecr.aws/lambda/ruby:4.0",
	RuntimeRuby34:         "public.ecr.aws/lambda/ruby:3.4",
	RuntimeRuby33:         "public.ecr.aws/lambda/ruby:3.3",
	RuntimeRuby32:         "public.ecr.aws/lambda/ruby:3.2",
	RuntimeProvidedAl2023: "public.ecr.aws/lambda/provided:al2023",
	RuntimeProvidedAl2:    "public.ecr.aws/lambda/provided:al2",
	RuntimeGo1X:           "public.ecr.aws/lambda/provided:al2",
}

// CanonicalRuntime normalises a runtime identifier to its canonical
// lowercase form and reports whether it identifies a currently supported
// runtime. The create, update and compatible-runtime paths persist the
// returned form so a case-variant request can never store a string the
// image lookup and the runtime-wrapper prefix checks would miss.
func CanonicalRuntime(runtime string) (Runtime, bool) {
	canonical := strings.ToLower(runtime)
	for _, r := range CurrentRuntimes {
		if string(r) == canonical {
			return r, true
		}
	}
	return "", false
}

// GetImageForRuntime returns the default container image for a runtime.
// The boolean reports whether the runtime has a mapped image; callers
// must treat a miss as an error rather than substituting a default, so
// an unmapped runtime can never silently execute as provided:al2.
func GetImageForRuntime(runtime Runtime) (string, bool) {
	image, ok := RuntimeImageMapping[runtime]
	return image, ok
}

// GenerateCodeHash generates a SHA-256 hash of the given data for code verification.
func GenerateCodeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}
