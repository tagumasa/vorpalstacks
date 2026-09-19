// Package stepfunction provides Step Functions storage functionality for vorpalstacks.
package sfn

import (
	"fmt"
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// StepFunctionStore provides Step Functions state machine and execution storage.
type StepFunctionStore struct {
	*common.BaseStore
	executionsStore       *common.BaseStore
	executionHistoryStore *common.BaseStore
	activitiesStore       *common.BaseStore
	tasksStore            *common.BaseStore
	versionsStore         *common.BaseStore
	aliasesStore          *common.BaseStore
	mapRunsStore          *common.BaseStore
	*common.TagStore
	arnBuilder         *svcarn.ARNBuilder
	accountID          string
	region             string
	pendingTasks       map[string]chan *ActivityTaskResult
	pendingTasksMu     sync.RWMutex
	activityQueues     map[string]chan *ActivityTask
	activityQueuesMu   sync.RWMutex
	tasksMu            sync.Mutex
	versionCounters    map[string]int64
	versionCountersMu  sync.Mutex
	activeExecutions   map[string]*executionRegistration
	activeExecutionsMu sync.RWMutex
	createMu           sync.Mutex
	mapRunSeq          int64
}

// NewStepFunctionStore creates a new StepFunctionStore instance.
func NewStepFunctionStore(store storage.BasicStorage, accountID, region string) *StepFunctionStore {
	return &StepFunctionStore{
		BaseStore:             common.NewBaseStore(store.Bucket("stepfunction-statemachines-"+region), "stepfunction-statemachines"),
		executionsStore:       common.NewBaseStore(store.Bucket("stepfunction-executions-"+region), "stepfunction-executions"),
		executionHistoryStore: common.NewBaseStore(store.Bucket("stepfunction-history-"+region), "stepfunction-history"),
		activitiesStore:       common.NewBaseStore(store.Bucket("stepfunction-activities-"+region), "stepfunction-activities"),
		tasksStore:            common.NewBaseStore(store.Bucket("stepfunction-tasks-"+region), "stepfunction-tasks"),
		versionsStore:         common.NewBaseStore(store.Bucket("stepfunction-versions-"+region), "stepfunction-versions"),
		aliasesStore:          common.NewBaseStore(store.Bucket("stepfunction-aliases-"+region), "stepfunction-aliases"),
		mapRunsStore:          common.NewBaseStore(store.Bucket("stepfunction-mapruns-"+region), "stepfunction-mapruns"),
		TagStore:              common.NewTagStoreWithRegion(store, "stepfunction", region, common.StandardTagBudget("TooManyTags")),
		arnBuilder:            svcarn.NewARNBuilder(accountID, region),
		accountID:             accountID,
		region:                region,
		pendingTasks:          make(map[string]chan *ActivityTaskResult),
		activityQueues:        make(map[string]chan *ActivityTask),
		activeExecutions:      make(map[string]*executionRegistration),
		versionCounters:       make(map[string]int64),
	}
}

// GetAccountID returns the AWS account ID.
func (s *StepFunctionStore) GetAccountID() string { return s.accountID }

// GetRegion returns the AWS region.
func (s *StepFunctionStore) GetRegion() string { return s.region }

func (s *StepFunctionStore) buildStateMachineARN(name string) string {
	return s.arnBuilder.StepFunctions().StateMachine(name)
}

func (s *StepFunctionStore) buildActivityARN(name string) string {
	return s.arnBuilder.StepFunctions().Activity(name)
}

// buildExecutionHistoryKey renders the storage key for one history event.
// The event ID is zero-padded to the width of the largest int64 so the
// lexicographic key order equals the numeric event-ID order; an unpadded
// decimal would sort event 10 between 1 and 2.
func (s *StepFunctionStore) buildExecutionHistoryKey(executionArn string, eventId int64) string {
	return fmt.Sprintf("%s:%019d", executionArn, eventId)
}
