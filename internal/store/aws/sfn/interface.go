package sfn

import (
	"context"
	"time"
)

// StepFunctionStoreInterface defines operations for managing Step Functions state machines and executions.
type StepFunctionStoreInterface interface {
	GetAccountID() string
	GetRegion() string
	CreateStateMachine(ctx context.Context, sm *StateMachine) error
	GetStateMachine(ctx context.Context, arn string) (*StateMachine, error)
	GetStateMachineByName(ctx context.Context, name string) (*StateMachine, error)
	UpdateStateMachine(ctx context.Context, sm *StateMachine) error
	DeleteStateMachine(ctx context.Context, arn string) error
	ListStateMachines(ctx context.Context, limit int32, nextToken string) (*StateMachineListResult, error)
	CreateExecution(ctx context.Context, exec *Execution) error
	GetExecution(ctx context.Context, arn string) (*Execution, error)
	UpdateExecution(ctx context.Context, exec *Execution) error
	ListExecutions(ctx context.Context, stateMachineArn string, statusFilter string, limit int32, nextToken string) (*ExecutionListResult, error)
	AddExecutionHistoryEvent(ctx context.Context, event *ExecutionHistoryEvent) error
	GetExecutionHistory(ctx context.Context, executionArn string, limit int32, nextToken string) ([]*ExecutionHistoryEvent, string, error)
	CreateActivity(ctx context.Context, activity *Activity) error
	GetActivity(ctx context.Context, arn string) (*Activity, error)
	DeleteActivity(ctx context.Context, arn string) error
	ListActivities(ctx context.Context, limit int32, nextToken string) (*ActivityListResult, error)
	CreateActivityTask(task *ActivityTask) error
	GetActivityTask(ctx context.Context, activityArn string, workerName string) (*ActivityTask, error)
	GetActivityTaskByToken(taskToken string) (*ActivityTask, error)
	CompleteActivityTask(taskToken string, output string) error
	FailActivityTask(taskToken string, errorMsg string, cause string) error
	WaitForTaskResult(ctx context.Context, taskToken string, timeout time.Duration) (*ActivityTaskResult, error)
	RegisterExecution(executionArn string, cancel context.CancelFunc)
	CancelExecution(executionArn string) bool
	UnregisterExecution(executionArn string)
	CancelAllExecutions()
	PublishStateMachineVersion(ctx context.Context, smArn string, description string) (*StateMachineVersion, error)
	GetStateMachineVersion(ctx context.Context, arn string) (*StateMachineVersion, error)
	DeleteStateMachineVersion(ctx context.Context, arn string) error
	ListStateMachineVersions(ctx context.Context, smArn string, limit int32, nextToken string) (*StateMachineVersionListResult, error)
	CreateStateMachineAlias(ctx context.Context, alias *StateMachineAlias) error
	GetStateMachineAlias(ctx context.Context, arn string) (*StateMachineAlias, error)
	UpdateStateMachineAlias(ctx context.Context, alias *StateMachineAlias) error
	DeleteStateMachineAlias(ctx context.Context, arn string) error
	ListStateMachineAliases(ctx context.Context, smArn string, limit int32, nextToken string) (*StateMachineAliasListResult, error)
}

var _ StepFunctionStoreInterface = (*StepFunctionStore)(nil)
