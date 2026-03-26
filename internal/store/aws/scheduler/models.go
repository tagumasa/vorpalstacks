package scheduler

import "time"

// ScheduleGroupState represents the state of a schedule group.
type ScheduleGroupState string

// Schedule group state constants.
const (
	// ScheduleGroupStateActive indicates the schedule group is active.
	ScheduleGroupStateActive ScheduleGroupState = "ACTIVE"
	// ScheduleGroupStateDeleting indicates the schedule group is being deleted.
	ScheduleGroupStateDeleting ScheduleGroupState = "DELETING"
)

// ScheduleState represents the state of a schedule.
type ScheduleState string

// Schedule state constants.
const (
	// ScheduleStateEnabled indicates the schedule is enabled.
	ScheduleStateEnabled ScheduleState = "ENABLED"
	// ScheduleStateDisabled indicates the schedule is disabled.
	ScheduleStateDisabled ScheduleState = "DISABLED"
)

// ActionAfterCompletion specifies the action to take after a schedule completes.
type ActionAfterCompletion string

// Action after completion constants.
const (
	// ActionAfterCompletionNone indicates no action after completion.
	ActionAfterCompletionNone ActionAfterCompletion = "NONE"
	// ActionAfterCompletionDelete indicates the schedule should be deleted after completion.
	ActionAfterCompletionDelete ActionAfterCompletion = "DELETE"
)

// FlexibleTimeWindowMode specifies the flexible time window mode.
type FlexibleTimeWindowMode string

// Flexible time window mode constants.
const (
	// FlexibleTimeWindowModeOff indicates no flexible time window.
	FlexibleTimeWindowModeOff FlexibleTimeWindowMode = "OFF"
	// FlexibleTimeWindowModeFlexible indicates a flexible time window is enabled.
	FlexibleTimeWindowModeFlexible FlexibleTimeWindowMode = "FLEXIBLE"
)

// ScheduleGroup represents a schedule group in Amazon EventBridge Scheduler.
type ScheduleGroup struct {
	Name                 string             `json:"name"`
	ARN                  string             `json:"arn"`
	State                ScheduleGroupState `json:"state"`
	CreationDate         time.Time          `json:"creationDate"`
	LastModificationDate time.Time          `json:"lastModificationDate"`
}

// FlexibleTimeWindow defines the flexible time window configuration for a schedule.
type FlexibleTimeWindow struct {
	Mode                   FlexibleTimeWindowMode `json:"mode"`
	MaximumWindowInMinutes *int                   `json:"maximumWindowInMinutes,omitempty"`
}

// DeadLetterConfig defines the dead-letter queue configuration for failed schedule invocations.
type DeadLetterConfig struct {
	Arn string `json:"arn,omitempty"`
}

// RetryPolicy defines the retry policy for schedule invocations.
type RetryPolicy struct {
	MaximumEventAgeInSeconds *int `json:"maximumEventAgeInSeconds,omitempty"`
	MaximumRetryAttempts     *int `json:"maximumRetryAttempts,omitempty"`
}

// AwsVpcConfiguration defines the VPC configuration for ECS tasks.
type AwsVpcConfiguration struct {
	Subnets        []string `json:"subnets,omitempty"`
	SecurityGroups []string `json:"securityGroups,omitempty"`
	AssignPublicIp string   `json:"assignPublicIp,omitempty"`
}

// NetworkConfiguration defines the network configuration for ECS tasks.
type NetworkConfiguration struct {
	AwsVpcConfiguration *AwsVpcConfiguration `json:"awsvpcConfiguration,omitempty"`
}

// PlacementConstraint defines a placement constraint for ECS tasks.
type PlacementConstraint struct {
	Type       string `json:"type,omitempty"`
	Expression string `json:"expression,omitempty"`
}

// PlacementStrategy defines a placement strategy for ECS tasks.
type PlacementStrategy struct {
	Type  string `json:"type,omitempty"`
	Field string `json:"field,omitempty"`
}

// CapacityProviderStrategyItem defines a capacity provider strategy item.
type CapacityProviderStrategyItem struct {
	CapacityProvider string `json:"capacityProvider"`
	Weight           *int   `json:"weight,omitempty"`
	Base             *int   `json:"base,omitempty"`
}

// EcsParameters defines the parameters for ECS target invocations.
type EcsParameters struct {
	TaskDefinitionArn        string                         `json:"taskDefinitionArn"`
	TaskCount                *int                           `json:"taskCount,omitempty"`
	LaunchType               string                         `json:"launchType,omitempty"`
	NetworkConfiguration     *NetworkConfiguration          `json:"networkConfiguration,omitempty"`
	PlatformVersion          string                         `json:"platformVersion,omitempty"`
	Group                    string                         `json:"group,omitempty"`
	CapacityProviderStrategy []CapacityProviderStrategyItem `json:"capacityProviderStrategy,omitempty"`
	EnableECSManagedTags     *bool                          `json:"enableECSManagedTags,omitempty"`
	EnableExecuteCommand     *bool                          `json:"enableExecuteCommand,omitempty"`
	PlacementConstraints     []PlacementConstraint          `json:"placementConstraints,omitempty"`
	PlacementStrategy        []PlacementStrategy            `json:"placementStrategy,omitempty"`
	PropagateTags            string                         `json:"propagateTags,omitempty"`
	ReferenceId              string                         `json:"referenceId,omitempty"`
	Tags                     []map[string]string            `json:"tags,omitempty"`
}

// EventBridgeParameters defines the parameters for EventBridge target invocations.
type EventBridgeParameters struct {
	DetailType string `json:"detailType"`
	Source     string `json:"source"`
}

// KinesisParameters defines the parameters for Kinesis target invocations.
type KinesisParameters struct {
	PartitionKey string `json:"partitionKey"`
}

// SageMakerPipelineParameter defines a parameter for SageMaker pipeline invocations.
type SageMakerPipelineParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SageMakerPipelineParameters defines the parameters for SageMaker pipeline invocations.
type SageMakerPipelineParameters struct {
	PipelineParameterList []SageMakerPipelineParameter `json:"pipelineParameterList,omitempty"`
}

// SqsParameters defines the parameters for SQS target invocations.
type SqsParameters struct {
	MessageGroupId string `json:"messageGroupId,omitempty"`
}

// Target defines the target for a schedule invocation.
type Target struct {
	Arn                         string                       `json:"arn"`
	RoleArn                     string                       `json:"roleArn"`
	Input                       string                       `json:"input,omitempty"`
	DeadLetterConfig            *DeadLetterConfig            `json:"deadLetterConfig,omitempty"`
	RetryPolicy                 *RetryPolicy                 `json:"retryPolicy,omitempty"`
	EcsParameters               *EcsParameters               `json:"ecsParameters,omitempty"`
	EventBridgeParameters       *EventBridgeParameters       `json:"eventBridgeParameters,omitempty"`
	KinesisParameters           *KinesisParameters           `json:"kinesisParameters,omitempty"`
	SageMakerPipelineParameters *SageMakerPipelineParameters `json:"sageMakerPipelineParameters,omitempty"`
	SqsParameters               *SqsParameters               `json:"sqsParameters,omitempty"`
}

// Schedule represents a schedule in Amazon EventBridge Scheduler.
type Schedule struct {
	Name                       string                `json:"name"`
	GroupName                  string                `json:"groupName"`
	ARN                        string                `json:"arn"`
	Region                     string                `json:"region"`
	ScheduleExpression         string                `json:"scheduleExpression"`
	ScheduleExpressionTimezone string                `json:"scheduleExpressionTimezone,omitempty"`
	State                      ScheduleState         `json:"state"`
	Description                string                `json:"description,omitempty"`
	StartDate                  *time.Time            `json:"startDate,omitempty"`
	EndDate                    *time.Time            `json:"endDate,omitempty"`
	KmsKeyArn                  string                `json:"kmsKeyArn,omitempty"`
	Target                     *Target               `json:"target"`
	FlexibleTimeWindow         *FlexibleTimeWindow   `json:"flexibleTimeWindow"`
	ActionAfterCompletion      ActionAfterCompletion `json:"actionAfterCompletion,omitempty"`
	CreationDate               time.Time             `json:"creationDate"`
	LastModificationDate       time.Time             `json:"lastModificationDate"`
}

// ScheduleGroupSummary represents a summary of a schedule group.
type ScheduleGroupSummary struct {
	Arn                  string             `json:"arn,omitempty"`
	Name                 string             `json:"name,omitempty"`
	State                ScheduleGroupState `json:"state,omitempty"`
	CreationDate         *time.Time         `json:"creationDate,omitempty"`
	LastModificationDate *time.Time         `json:"lastModificationDate,omitempty"`
}

// ScheduleSummary represents a summary of a schedule.
type ScheduleSummary struct {
	Arn                  string         `json:"arn,omitempty"`
	Name                 string         `json:"name,omitempty"`
	GroupName            string         `json:"groupName,omitempty"`
	State                ScheduleState  `json:"state,omitempty"`
	CreationDate         *time.Time     `json:"creationDate,omitempty"`
	LastModificationDate *time.Time     `json:"lastModificationDate,omitempty"`
	Target               *TargetSummary `json:"target,omitempty"`
}

// TargetSummary represents a summary of a schedule target.
type TargetSummary struct {
	Arn string `json:"arn"`
}

// ScheduleGroupListResult represents the result of listing schedule groups.
type ScheduleGroupListResult struct {
	ScheduleGroups []ScheduleGroupSummary `json:"scheduleGroups"`
	NextToken      string                 `json:"nextToken,omitempty"`
}

// ScheduleListResult represents the result of listing schedules.
type ScheduleListResult struct {
	Schedules []ScheduleSummary `json:"schedules"`
	NextToken string            `json:"nextToken,omitempty"`
}
