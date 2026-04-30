// Package lambda provides AWS Lambda store functionality for vorpalstacks.
package lambda

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"vorpalstacks/internal/utils/aws/types"
)

// Runtime represents the Lambda function runtime environment.
type Runtime string

// Runtime constants define supported Lambda runtime environments.
const (
	RuntimeNodejs22X      Runtime = "nodejs22.x"
	RuntimeNodejs20X      Runtime = "nodejs20.x"
	RuntimeNodejs18X      Runtime = "nodejs18.x"
	RuntimeNodejs16X      Runtime = "nodejs16.x"
	RuntimePython313      Runtime = "python3.13"
	RuntimePython312      Runtime = "python3.12"
	RuntimePython311      Runtime = "python3.11"
	RuntimePython310      Runtime = "python3.10"
	RuntimePython39       Runtime = "python3.9"
	RuntimePython38       Runtime = "python3.8"
	RuntimeJava21         Runtime = "java21"
	RuntimeJava17         Runtime = "java17"
	RuntimeJava11         Runtime = "java11"
	RuntimeJava8Al2       Runtime = "java8.al2"
	RuntimeDotnet8        Runtime = "dotnet8"
	RuntimeDotnet6        Runtime = "dotnet6"
	RuntimeRuby34         Runtime = "ruby3.4"
	RuntimeRuby33         Runtime = "ruby3.3"
	RuntimeRuby32         Runtime = "ruby3.2"
	RuntimeProvidedAl2023 Runtime = "provided.al2023"
	RuntimeProvidedAl2    Runtime = "provided.al2"
	RuntimeGo1X           Runtime = "go1.x"
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
	FunctionName             string             `json:"function_name"`
	FunctionArn              string             `json:"function_arn"`
	Runtime                  Runtime            `json:"runtime"`
	Role                     string             `json:"role"`
	Handler                  string             `json:"handler"`
	CodeSize                 int64              `json:"code_size"`
	CodeSha256               string             `json:"code_sha256"`
	CodeLocation             string             `json:"code_location"`
	ImageUri                 string             `json:"image_uri,omitempty"`
	SourceCodeHash           string             `json:"source_code_hash,omitempty"`
	Description              string             `json:"description,omitempty"`
	Timeout                  int32              `json:"timeout"`
	MemorySize               int32              `json:"memory_size"`
	EphemeralStorage         *EphemeralStorage  `json:"ephemeral_storage,omitempty"`
	Publish                  bool               `json:"publish"`
	Architectures            []string           `json:"architectures,omitempty"`
	KMSKeyArn                string             `json:"kms_key_arn,omitempty"`
	RevisionId               string             `json:"revision_id"`
	State                    State              `json:"state"`
	StateReason              string             `json:"state_reason,omitempty"`
	StateReasonCode          string             `json:"state_reason_code,omitempty"`
	LastUpdateStatus         LastUpdateStatus   `json:"last_update_status"`
	LastUpdateReason         string             `json:"last_update_reason,omitempty"`
	LastUpdateStatusReason   string             `json:"last_update_status_reason,omitempty"`
	LastModified             time.Time          `json:"last_modified"`
	LastModifiedUser         string             `json:"last_modified_user,omitempty"`
	VpcConfig                *VpcConfig         `json:"vpc_config,omitempty"`
	Environment              *Environment       `json:"environment,omitempty"`
	DeadLetterConfig         *DeadLetterConfig  `json:"dead_letter_config,omitempty"`
	TracingConfig            *TracingConfig     `json:"tracing_config,omitempty"`
	Layers                   []LayerReference   `json:"layers,omitempty"`
	Tags                     []types.Tag        `json:"tags,omitempty"`
	SnapStart                *SnapStart         `json:"snap_start,omitempty"`
	PackageType              string             `json:"package_type"`
	SigningProfileVersionArn string             `json:"signing_profile_version_arn,omitempty"`
	SigningJobArn            string             `json:"signing_job_arn,omitempty"`
	UrlConfig                *FunctionUrlConfig `json:"url_config,omitempty"`
	CodeSigningConfigArn     string             `json:"code_signing_config_arn,omitempty"`

	Versions       []Version        `json:"versions,omitempty"`
	Aliases        []Alias          `json:"aliases,omitempty"`
	Policies       []FunctionPolicy `json:"policies,omitempty"`
	CurrentVersion string           `json:"current_version"`

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
	Version                  string            `json:"version"`
	FunctionArn              string            `json:"function_arn"`
	Runtime                  Runtime           `json:"runtime"`
	Role                     string            `json:"role"`
	Handler                  string            `json:"handler"`
	CodeSize                 int64             `json:"code_size"`
	CodeSha256               string            `json:"code_sha256"`
	CodeLocation             string            `json:"code_location"`
	ImageUri                 string            `json:"image_uri,omitempty"`
	Description              string            `json:"description,omitempty"`
	Timeout                  int32             `json:"timeout"`
	MemorySize               int32             `json:"memory_size"`
	EphemeralStorage         *EphemeralStorage `json:"ephemeral_storage,omitempty"`
	Architectures            []string          `json:"architectures,omitempty"`
	KMSKeyArn                string            `json:"kms_key_arn,omitempty"`
	RevisionId               string            `json:"revision_id"`
	State                    State             `json:"state"`
	StateReason              string            `json:"state_reason,omitempty"`
	StateReasonCode          string            `json:"state_reason_code,omitempty"`
	LastUpdateStatus         LastUpdateStatus  `json:"last_update_status"`
	LastModified             time.Time         `json:"last_modified"`
	VpcConfig                *VpcConfig        `json:"vpc_config,omitempty"`
	Environment              *Environment      `json:"environment,omitempty"`
	DeadLetterConfig         *DeadLetterConfig `json:"dead_letter_config,omitempty"`
	TracingConfig            *TracingConfig    `json:"tracing_config,omitempty"`
	Layers                   []LayerReference  `json:"layers,omitempty"`
	SnapStart                *SnapStart        `json:"snap_start,omitempty"`
	PackageType              string            `json:"package_type"`
	SigningProfileVersionArn string            `json:"signing_profile_version_arn,omitempty"`
	SigningJobArn            string            `json:"signing_job_arn,omitempty"`

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
type FunctionPolicy struct {
	Id        string                 `json:"id"`
	Statement string                 `json:"statement"`
	Principal string                 `json:"principal"`
	Action    string                 `json:"action"`
	Resource  string                 `json:"resource"`
	Condition map[string]interface{} `json:"condition,omitempty"`
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
	Statement string `json:"statement"`
	Principal string `json:"principal"`
	Action    string `json:"action"`
}

// EventSourceMapping represents a mapping between an event source and a Lambda function.
type EventSourceMapping struct {
	UUID                           string                      `json:"uuid"`
	BatchSize                      int32                       `json:"batch_size,omitempty"`
	MaximumBatchingWindowInSeconds int32                       `json:"maximum_batching_window_in_seconds,omitempty"`
	ParallelizationFactor          int32                       `json:"parallelization_factor,omitempty"`
	EventSourceArn                 string                      `json:"event_source_arn,omitempty"`
	FunctionArn                    string                      `json:"function_arn"`
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

// CodeSigningConfig represents the code signing configuration for a Lambda function.
type CodeSigningConfig struct {
	CodeSigningConfigId  string               `json:"code_signing_config_id"`
	CodeSigningConfigArn string               `json:"code_signing_config_arn"`
	Description          string               `json:"description,omitempty"`
	AllowedPublishers    *AllowedPublishers   `json:"allowed_publishers,omitempty"`
	CodeSigningPolicies  *CodeSigningPolicies `json:"code_signing_policies,omitempty"`
	LastModified         time.Time            `json:"last_modified"`
}

// AllowedPublishers represents the allowed publishers for code signing.
type AllowedPublishers struct {
	SigningProfileVersionArns []string `json:"signing_profile_version_arns,omitempty"`
}

// CodeSigningPolicies represents the code signing policies for a Lambda function.
type CodeSigningPolicies struct {
	UntrustedArtifactOnDeployment string `json:"untrusted_artifact_on_deployment,omitempty"`
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

// FunctionListResult represents the result of listing Lambda functions.
type FunctionListResult struct {
	Functions   []*Function
	NextMarker  string
	IsTruncated bool
}

// RuntimeImageMapping maps Lambda runtimes to their default container images.
var RuntimeImageMapping = map[Runtime]string{
	RuntimeNodejs22X:      "public.ecr.aws/lambda/nodejs:22",
	RuntimeNodejs20X:      "public.ecr.aws/lambda/nodejs:20",
	RuntimeNodejs18X:      "public.ecr.aws/lambda/nodejs:18",
	RuntimeNodejs16X:      "public.ecr.aws/lambda/nodejs:16",
	RuntimePython313:      "public.ecr.aws/lambda/python:3.13",
	RuntimePython312:      "public.ecr.aws/lambda/python:3.12",
	RuntimePython311:      "public.ecr.aws/lambda/python:3.11",
	RuntimePython310:      "public.ecr.aws/lambda/python:3.10",
	RuntimePython39:       "public.ecr.aws/lambda/python:3.9",
	RuntimePython38:       "public.ecr.aws/lambda/python:3.8",
	RuntimeJava21:         "public.ecr.aws/lambda/java:21",
	RuntimeJava17:         "public.ecr.aws/lambda/java:17",
	RuntimeJava11:         "public.ecr.aws/lambda/java:11",
	RuntimeJava8Al2:       "public.ecr.aws/lambda/java:8.al2",
	RuntimeDotnet8:        "public.ecr.aws/lambda/dotnet:8",
	RuntimeDotnet6:        "public.ecr.aws/lambda/dotnet:6",
	RuntimeRuby34:         "public.ecr.aws/lambda/ruby:3.4",
	RuntimeRuby33:         "public.ecr.aws/lambda/ruby:3.3",
	RuntimeRuby32:         "public.ecr.aws/lambda/ruby:3.2",
	RuntimeProvidedAl2023: "public.ecr.aws/lambda/provided:al2023",
	RuntimeProvidedAl2:    "public.ecr.aws/lambda/provided:al2",
	RuntimeGo1X:           "public.ecr.aws/lambda/provided:al2",
}

// GetImageForRuntime returns the default container image for a given Lambda runtime.
func GetImageForRuntime(runtime Runtime) string {
	if image, ok := RuntimeImageMapping[runtime]; ok {
		return image
	}
	return "public.ecr.aws/lambda/provided:al2"
}

// GenerateCodeHash generates a SHA-256 hash of the given data for code verification.
func GenerateCodeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}
