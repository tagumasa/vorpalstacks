package scheduler

import (
	"context"

	"vorpalstacks/internal/utils/aws/types"
)

// SchedulerStoreInterface defines operations for managing EventBridge Scheduler schedules.
type SchedulerStoreInterface interface {
	GetAccountID() string
	GetRegion() string
	BuildScheduleARNFromName(name string) string
	CreateScheduleGroup(ctx context.Context, group *ScheduleGroup) error
	GetScheduleGroup(ctx context.Context, name string) (*ScheduleGroup, error)
	DeleteScheduleGroup(ctx context.Context, name string) error
	ListScheduleGroups(ctx context.Context, namePrefix string, limit int32, nextToken string) (*ScheduleGroupListResult, error)
	UpdateScheduleGroup(ctx context.Context, group *ScheduleGroup) error
	CreateSchedule(ctx context.Context, schedule *Schedule) error
	GetSchedule(ctx context.Context, groupName, name string) (*Schedule, error)
	UpdateSchedule(ctx context.Context, schedule *Schedule) error
	DeleteSchedule(ctx context.Context, groupName, name string) error
	ListSchedules(ctx context.Context, groupName, namePrefix string, state ScheduleState, limit int32, nextToken string) (*ScheduleListResult, error)
	GetAllEnabledSchedules(ctx context.Context) ([]*Schedule, error)
	EnsureDefaultGroup(ctx context.Context) error
	TagFromSlice(arn string, tags []types.Tag) error
	Untag(arn string, tagKeys []string) error
	ListAsSlice(arn string) ([]types.Tag, error)
	Raw() *SchedulerStore
}

var _ SchedulerStoreInterface = (*SchedulerStore)(nil)

// Raw returns the underlying scheduler store.
func (s *SchedulerStore) Raw() *SchedulerStore {
	return s
}
