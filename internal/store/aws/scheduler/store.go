// Package scheduler provides EventBridge Scheduler storage functionality for vorpalstacks.
package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// SchedulerStore provides EventBridge Scheduler storage functionality.
type SchedulerStore struct {
	*common.BaseStore
	schedulesStore *common.BaseStore
	*common.TagStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	createMu   sync.Mutex
}

// NewSchedulerStore creates a new Scheduler store instance.
//
// Parameters:
//   - store: The storage instance
//   - accountID: The AWS account ID
//   - region: The AWS region
//
// Returns:
//   - *SchedulerStore: A new scheduler store instance
func NewSchedulerStore(store storage.BasicStorage, accountID, region string) *SchedulerStore {
	return &SchedulerStore{
		BaseStore:      common.NewBaseStore(store.Bucket("scheduler-groups-"+region), "scheduler-groups"),
		schedulesStore: common.NewBaseStore(store.Bucket("scheduler-schedules-"+region), "scheduler-schedules"),
		TagStore:       common.NewTagStoreWithRegion(store, "scheduler", region),
		arnBuilder:     svcarn.NewARNBuilder(accountID, region),
		accountID:      accountID,
		region:         region,
	}
}

// GetAccountID returns the AWS account ID associated with this store.
func (s *SchedulerStore) GetAccountID() string {
	return s.accountID
}

// GetRegion returns the AWS region associated with this store.
func (s *SchedulerStore) GetRegion() string {
	return s.region
}

func (s *SchedulerStore) buildScheduleGroupARN(name string) string {
	return s.arnBuilder.Scheduler().ScheduleGroup(name)
}

func (s *SchedulerStore) buildScheduleARN(groupName, scheduleName string) string {
	return s.arnBuilder.Scheduler().Schedule(groupName, scheduleName)
}

// BuildScheduleARNFromName builds an ARN for a schedule with the default group.
func (s *SchedulerStore) BuildScheduleARNFromName(name string) string {
	return s.buildScheduleARN("default", name)
}

// ScheduleGroup operations

// CreateScheduleGroup creates a new schedule group.
//
// Parameters:
//   - ctx: The context
//   - group: The schedule group to create
//
// Returns:
//   - error: An error if creation fails
func (s *SchedulerStore) CreateScheduleGroup(ctx context.Context, group *ScheduleGroup) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if group.Name == "" {
		return ErrInvalidARN
	}

	arn := s.buildScheduleGroupARN(group.Name)
	if s.Exists(arn) {
		return ErrScheduleGroupAlreadyExists
	}

	now := time.Now().UTC()
	group.ARN = arn
	group.State = ScheduleGroupStateActive
	group.CreationDate = now
	group.LastModificationDate = now

	return s.Put(arn, group)
}

// GetScheduleGroup retrieves a schedule group by name.
//
// Parameters:
//   - ctx: The context
//   - name: The schedule group name
//
// Returns:
//   - *ScheduleGroup: The schedule group if found
//   - error: An error if not found
func (s *SchedulerStore) GetScheduleGroup(ctx context.Context, name string) (*ScheduleGroup, error) {
	arn := s.buildScheduleGroupARN(name)
	var group ScheduleGroup
	if err := s.BaseStore.Get(arn, &group); err != nil {
		return nil, ErrScheduleGroupNotFound
	}
	return &group, nil
}

// DeleteScheduleGroup deletes a schedule group by name.
//
// Parameters:
//   - ctx: The context
//   - name: The schedule group name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *SchedulerStore) DeleteScheduleGroup(ctx context.Context, name string) error {
	arn := s.buildScheduleGroupARN(name)
	if !s.Exists(arn) {
		return ErrScheduleGroupNotFound
	}

	opts := common.ListOptions{MaxItems: 1000}
	result, err := common.List[Schedule](s.schedulesStore, opts, func(sch *Schedule) bool {
		return s.buildScheduleGroupARN(sch.GroupName) == arn
	})
	if err != nil {
		return err
	}
	if len(result.Items) > 0 {
		return ErrScheduleGroupNotEmpty
	}

	return s.BaseStore.Delete(arn)
}

// ListScheduleGroups lists schedule groups with optional filtering.
//
// Parameters:
//   - ctx: The context
//   - namePrefix: Optional name prefix filter
//   - limit: Maximum number of results
//   - nextToken: Pagination token
//
// Returns:
//   - *ScheduleGroupListResult: The list result with schedule groups and next token
//   - error: An error if listing fails
func (s *SchedulerStore) ListScheduleGroups(ctx context.Context, namePrefix string, limit int32, nextToken string) (*ScheduleGroupListResult, error) {
	opts := common.ListOptions{
		Prefix:   "",
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[ScheduleGroup](s.BaseStore, opts, func(g *ScheduleGroup) bool {
		if namePrefix == "" {
			return true
		}
		return len(g.Name) >= len(namePrefix) && g.Name[:len(namePrefix)] == namePrefix
	})
	if err != nil {
		return nil, err
	}

	summaries := make([]ScheduleGroupSummary, len(result.Items))
	for i, g := range result.Items {
		summaries[i] = ScheduleGroupSummary{
			Arn:                  g.ARN,
			Name:                 g.Name,
			State:                g.State,
			CreationDate:         &g.CreationDate,
			LastModificationDate: &g.LastModificationDate,
		}
	}

	return &ScheduleGroupListResult{
		ScheduleGroups: summaries,
		NextToken:      result.NextMarker,
	}, nil
}

// UpdateScheduleGroup updates an existing schedule group.
//
// Parameters:
//   - ctx: The context
//   - group: The schedule group to update
//
// Returns:
//   - error: An error if update fails
func (s *SchedulerStore) UpdateScheduleGroup(ctx context.Context, group *ScheduleGroup) error {
	if !s.Exists(group.ARN) {
		return ErrScheduleGroupNotFound
	}
	group.LastModificationDate = time.Now().UTC()
	return s.Put(group.ARN, group)
}

// Schedule operations

func (s *SchedulerStore) buildScheduleKey(groupName, scheduleName string) string {
	return fmt.Sprintf("%s:%s", groupName, scheduleName)
}

// CreateSchedule creates a new schedule.
//
// Parameters:
//   - ctx: The context
//   - schedule: The schedule to create
//
// Returns:
//   - error: An error if creation fails
func (s *SchedulerStore) CreateSchedule(ctx context.Context, schedule *Schedule) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	if schedule.Name == "" {
		return ErrInvalidARN
	}

	groupName := schedule.GroupName
	if groupName == "" {
		groupName = "default"
	}
	schedule.GroupName = groupName

	key := s.buildScheduleKey(groupName, schedule.Name)
	if s.schedulesStore.Exists(key) {
		return ErrScheduleAlreadyExists
	}

	now := time.Now().UTC()
	schedule.ARN = s.buildScheduleARN(groupName, schedule.Name)
	schedule.CreationDate = now
	schedule.LastModificationDate = now
	if schedule.State == "" {
		schedule.State = ScheduleStateEnabled
	}

	return s.schedulesStore.Put(key, schedule)
}

// GetSchedule retrieves a schedule by group name and schedule name.
//
// Parameters:
//   - ctx: The context
//   - groupName: The schedule group name
//   - name: The schedule name
//
// Returns:
//   - *Schedule: The schedule if found
//   - error: An error if not found
func (s *SchedulerStore) GetSchedule(ctx context.Context, groupName, name string) (*Schedule, error) {
	if groupName == "" {
		groupName = "default"
	}
	key := s.buildScheduleKey(groupName, name)
	var schedule Schedule
	if err := s.schedulesStore.Get(key, &schedule); err != nil {
		return nil, ErrScheduleNotFound
	}
	return &schedule, nil
}

// UpdateSchedule updates an existing schedule.
//
// Parameters:
//   - ctx: The context
//   - schedule: The schedule to update
//
// Returns:
//   - error: An error if update fails
func (s *SchedulerStore) UpdateSchedule(ctx context.Context, schedule *Schedule) error {
	groupName := schedule.GroupName
	if groupName == "" {
		groupName = "default"
	}
	key := s.buildScheduleKey(groupName, schedule.Name)
	if !s.schedulesStore.Exists(key) {
		return ErrScheduleNotFound
	}
	schedule.LastModificationDate = time.Now().UTC()
	return s.schedulesStore.Put(key, schedule)
}

// DeleteSchedule deletes a schedule by group name and schedule name.
//
// Parameters:
//   - ctx: The context
//   - groupName: The schedule group name
//   - name: The schedule name to delete
//
// Returns:
//   - error: An error if deletion fails
func (s *SchedulerStore) DeleteSchedule(ctx context.Context, groupName, name string) error {
	if groupName == "" {
		groupName = "default"
	}
	key := s.buildScheduleKey(groupName, name)
	if !s.schedulesStore.Exists(key) {
		return ErrScheduleNotFound
	}
	return s.schedulesStore.Delete(key)
}

// ListSchedules lists schedules with optional filtering.
//
// Parameters:
//   - ctx: The context
//   - groupName: The schedule group name filter
//   - namePrefix: Optional name prefix filter
//   - state: Optional state filter
//   - limit: Maximum number of results
//   - nextToken: Pagination token
//
// Returns:
//   - *ScheduleListResult: The list result with schedules and next token
//   - error: An error if listing fails
func (s *SchedulerStore) ListSchedules(ctx context.Context, groupName, namePrefix string, state ScheduleState, limit int32, nextToken string) (*ScheduleListResult, error) {
	prefix := ""
	if groupName != "" {
		prefix = groupName + ":"
	}

	opts := common.ListOptions{
		Prefix:   prefix,
		Marker:   nextToken,
		MaxItems: int(limit),
	}

	result, err := common.List[Schedule](s.schedulesStore, opts, func(s *Schedule) bool {
		if namePrefix != "" {
			if len(s.Name) < len(namePrefix) || s.Name[:len(namePrefix)] != namePrefix {
				return false
			}
		}
		if state != "" && s.State != state {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	summaries := make([]ScheduleSummary, len(result.Items))
	for i, s := range result.Items {
		summary := ScheduleSummary{
			Arn:                  s.ARN,
			Name:                 s.Name,
			GroupName:            s.GroupName,
			State:                s.State,
			CreationDate:         &s.CreationDate,
			LastModificationDate: &s.LastModificationDate,
		}
		if s.Target != nil {
			summary.Target = &TargetSummary{Arn: s.Target.Arn}
		}
		summaries[i] = summary
	}

	return &ScheduleListResult{
		Schedules: summaries,
		NextToken: result.NextMarker,
	}, nil
}

// GetAllEnabledSchedules retrieves all enabled schedules.
//
// Parameters:
//   - ctx: The context
//
// Returns:
//   - []*Schedule: The list of enabled schedules
//   - error: An error if retrieval fails
func (s *SchedulerStore) GetAllEnabledSchedules(ctx context.Context) ([]*Schedule, error) {
	opts := common.ListOptions{
		Prefix:   "",
		Marker:   "",
		MaxItems: 0,
	}

	result, err := common.List[Schedule](s.schedulesStore, opts, func(sch *Schedule) bool {
		return sch.State == ScheduleStateEnabled
	})
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// EnsureDefaultGroup creates the default schedule group if it doesn't exist.
//
// Parameters:
//   - ctx: The context
//
// Returns:
//   - error: An error if creation fails
func (s *SchedulerStore) EnsureDefaultGroup(ctx context.Context) error {
	s.createMu.Lock()
	defer s.createMu.Unlock()
	arn := s.buildScheduleGroupARN("default")
	if s.Exists(arn) {
		return nil
	}

	now := time.Now().UTC()
	group := &ScheduleGroup{
		Name:                 "default",
		ARN:                  arn,
		State:                ScheduleGroupStateActive,
		CreationDate:         now,
		LastModificationDate: now,
	}
	return s.Put(arn, group)
}
