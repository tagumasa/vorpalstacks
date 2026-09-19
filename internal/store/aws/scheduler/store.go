// Package scheduler provides EventBridge Scheduler storage functionality for vorpalstacks.
package scheduler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// SchedulerStore provides EventBridge Scheduler storage functionality.
type SchedulerStore struct {
	*common.BaseStore
	schedulesStore *common.BaseStore
	*common.TagStore
	clientTokens *ClientTokenStore
	arnBuilder   *svcarn.ARNBuilder
	region       string
}

// scheduleRecordWriteMu serialises read-modify-write cycles on schedule
// records across SchedulerStore instances. The engine and the service now
// share one store cache, so one instance per region exists in practice;
// the package scope keeps CompleteSchedule's read-modify-write from losing
// a concurrent UpdateSchedule write (or vice versa) independent of that
// invariant. Every record writer — group records included — takes this one
// lock: the group-state transitions (mark, purge) and the check-then-write
// creates (schedules, groups) must not interleave, or a schedule could be
// admitted into a group being deleted, or a group record recreated between
// a purge's emptiness check and its delete.
var scheduleRecordWriteMu sync.Mutex

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
		TagStore:       common.NewTagStoreWithRegion(store, "scheduler", region, common.StandardTagBudget("ValidationException")),
		clientTokens:   NewClientTokenStore(common.NewBaseStore(store.Bucket("scheduler-tokens-"+region), "scheduler-tokens")),
		arnBuilder:     svcarn.NewARNBuilder(accountID, region),
		region:         region,
	}
}

// ClientTokens returns the per-region ClientToken idempotency store.
func (s *SchedulerStore) ClientTokens() *ClientTokenStore {
	return s.clientTokens
}

// Close stops the background cleanup goroutine in the ClientTokenStore.
func (s *SchedulerStore) Close() {
	if s.clientTokens != nil {
		s.clientTokens.Stop()
	}
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

// BuildScheduleARN builds an ARN for a schedule with an explicit group.
func (s *SchedulerStore) BuildScheduleARN(groupName, name string) string {
	return s.buildScheduleARN(groupName, name)
}

// BuildScheduleGroupARN builds an ARN for a schedule group.
func (s *SchedulerStore) BuildScheduleGroupARN(name string) string {
	return s.buildScheduleGroupARN(name)
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
	// The record lock keeps the name-existence check and the write atomic
	// against the other group writers (mark, purge): recreating a record a
	// purge is mid-way through deleting must not interleave with it.
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()
	if group.Name == "" {
		return ErrInvalidName
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
//   - error: ErrScheduleGroupNotFound when the group is absent; any other
//     error is a storage failure returned as-is so callers can tell a
//     client-facing not-found from a server fault
func (s *SchedulerStore) GetScheduleGroup(ctx context.Context, name string) (*ScheduleGroup, error) {
	arn := s.buildScheduleGroupARN(name)
	var group ScheduleGroup
	if err := s.BaseStore.Get(arn, &group); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrScheduleGroupNotFound
		}
		return nil, err
	}
	return &group, nil
}

// MarkScheduleGroupDeleting transitions a schedule group to the DELETING
// state. Deleting a group cascades: per the DeleteScheduleGroup model
// documentation, the group remains in DELETING until all of its schedules
// are deleted; the engine's sweep completes the removal.
func (s *SchedulerStore) MarkScheduleGroupDeleting(ctx context.Context, name string) error {
	arn := s.buildScheduleGroupARN(name)
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()
	if !s.Exists(arn) {
		return ErrScheduleGroupNotFound
	}
	var group ScheduleGroup
	if err := s.BaseStore.Get(arn, &group); err != nil {
		return err
	}
	group.State = ScheduleGroupStateDeleting
	group.LastModificationDate = time.Now().UTC()
	return s.Put(arn, &group)
}

// ListDeletingScheduleGroups returns every schedule group currently in the
// DELETING state (the engine's cascade sweep input). The walk is
// unbounded: the paged List form clamps a zero MaxItems to the default
// page size and would silently starve the groups beyond the first page.
func (s *SchedulerStore) ListDeletingScheduleGroups(ctx context.Context) ([]*ScheduleGroup, error) {
	return common.ListMatching[ScheduleGroup](s.BaseStore, "", func(g *ScheduleGroup) bool {
		return g.State == ScheduleGroupStateDeleting
	})
}

// DeleteSchedulesInGroup deletes every schedule that belongs to the group.
func (s *SchedulerStore) DeleteSchedulesInGroup(ctx context.Context, groupName string) error {
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()
	schedules, err := common.ListMatching[Schedule](s.schedulesStore, groupName+":", nil)
	if err != nil {
		return err
	}
	for _, sch := range schedules {
		key := s.buildScheduleKey(sch.GroupName, sch.Name)
		if err := s.schedulesStore.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// PurgeDeletedScheduleGroup removes the group record and its tags once the
// group has no member schedules left.
func (s *SchedulerStore) PurgeDeletedScheduleGroup(ctx context.Context, name string) error {
	arn := s.buildScheduleGroupARN(name)
	// The record lock keeps the emptiness check and the delete in one
	// critical section against every record writer: a CreateSchedule cannot
	// add a member and a CreateScheduleGroup cannot recreate the record
	// between the two.
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()
	if !s.Exists(arn) {
		return nil
	}
	schedules, err := common.ListMatching[Schedule](s.schedulesStore, name+":", nil)
	if err != nil {
		return err
	}
	if len(schedules) > 0 {
		return ErrScheduleGroupNotEmpty
	}
	// Delete the primary resource first so tag metadata I/O errors never
	// block resource lifecycle. Tag cleanup is best-effort.
	if err := s.BaseStore.Delete(arn); err != nil {
		return err
	}
	if err := s.TagStore.Delete(arn); err != nil {
		logs.Warn("Failed to clean up tags for deleted schedule group (orphaned tags may remain)",
			logs.String("arn", arn),
			logs.Err(err))
	}
	return nil
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
		return strings.HasPrefix(g.Name, namePrefix)
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
	// The group name arrives resolved: the service layer owns the
	// default-group semantics, so an empty name is a malformed identity (it
	// would persist the record under a phantom key, not the implicit
	// default group).
	if schedule.Name == "" || schedule.GroupName == "" {
		return ErrInvalidName
	}

	// The record lock holds the duplicate-name check, the group-state
	// check, and the write in one critical section against the group
	// writers: a schedule can neither be admitted into a group already
	// DELETING nor slip into the window between a purge's emptiness check
	// and its delete.
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()

	key := s.buildScheduleKey(schedule.GroupName, schedule.Name)
	if s.schedulesStore.Exists(key) {
		return ErrScheduleAlreadyExists
	}
	if err := s.requireGroupNotDeleting(schedule.GroupName); err != nil {
		return err
	}

	now := time.Now().UTC()
	schedule.ARN = s.buildScheduleARN(schedule.GroupName, schedule.Name)
	schedule.CreationDate = now
	schedule.LastModificationDate = now

	return s.schedulesStore.Put(key, schedule)
}

// requireGroupNotDeleting reports whether groupName may receive schedule
// records: the group record must exist and must not be in the DELETING
// state, whose engine cascade would silently destroy any schedule admitted
// now. The implicit default group is the exception — it can never be
// deleted, so an unseeded record is not a refusal for it.
func (s *SchedulerStore) requireGroupNotDeleting(groupName string) error {
	groupArn := s.buildScheduleGroupARN(groupName)
	var group ScheduleGroup
	if err := s.BaseStore.Get(groupArn, &group); err != nil {
		if !common.IsNotFound(err) {
			return err
		}
		if groupName == DefaultGroupName {
			return nil
		}
		return ErrScheduleGroupNotFound
	}
	if group.State == ScheduleGroupStateDeleting {
		return ErrScheduleGroupDeleting
	}
	return nil
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
//   - error: ErrScheduleNotFound when the schedule is absent; any other
//     error is a storage failure returned as-is so callers can tell a
//     client-facing not-found from a server fault
func (s *SchedulerStore) GetSchedule(ctx context.Context, groupName, name string) (*Schedule, error) {
	key := s.buildScheduleKey(groupName, name)
	var schedule Schedule
	if err := s.schedulesStore.Get(key, &schedule); err != nil {
		if common.IsNotFound(err) {
			return nil, ErrScheduleNotFound
		}
		return nil, err
	}
	return &schedule, nil
}

// MutateSchedule applies fn to the schedule record inside the write mutex
// so the whole read-modify-write cycle is atomic with respect to every
// other record writer — the engine and the service share one store cache
// (one instance per region in practice), and the package-scope lock keeps
// the cycle atomic independent of that invariant. fn mutates the record in
// place; callers decide whether to stamp LastModificationDate —
// user-initiated updates do, internal markers must not.
//
// Parameters:
//   - ctx: The context
//   - groupName: The schedule group (resolved by the caller)
//   - name: The schedule name
//   - fn: The mutation applied to the record read under the lock
//
// Returns:
//   - error: ErrScheduleNotFound if the schedule does not exist
func (s *SchedulerStore) MutateSchedule(ctx context.Context, groupName, name string, fn func(*Schedule) error) error {
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()

	return s.mutateScheduleLocked(ctx, groupName, name, fn)
}

// MutateScheduleInLiveGroup applies fn to the schedule record with the
// group-state check held inside the same critical section as the write:
// an update may not land in a group whose deletion cascade started after
// the caller's fast-path probe, because the cascade would destroy the
// acknowledged write. The engine's own markers (completion, delivered
// boundaries) keep the plain MutateSchedule — those must keep working on
// the schedules of a group already being deleted, exactly until the
// cascade removes them.
//
// Returns:
//   - error: ErrScheduleGroupDeleting or ErrScheduleGroupNotFound from the
//     in-lock check; ErrScheduleNotFound if the schedule does not exist
func (s *SchedulerStore) MutateScheduleInLiveGroup(ctx context.Context, groupName, name string, fn func(*Schedule) error) error {
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()

	if err := s.requireGroupNotDeleting(groupName); err != nil {
		return err
	}
	return s.mutateScheduleLocked(ctx, groupName, name, fn)
}

// mutateScheduleLocked performs the locked read-modify-write; the caller
// holds scheduleRecordWriteMu.
func (s *SchedulerStore) mutateScheduleLocked(ctx context.Context, groupName, name string, fn func(*Schedule) error) error {
	schedule, err := s.GetSchedule(ctx, groupName, name)
	if err != nil {
		return err
	}
	if err := fn(schedule); err != nil {
		return err
	}
	key := s.buildScheduleKey(groupName, name)
	return s.schedulesStore.Put(key, schedule)
}

// CompleteSchedule marks a one-time schedule's execution lifecycle as
// ended. The schedule keeps its wire state — the AWS ScheduleState enum
// has no COMPLETED value — but the engine stops firing it, including
// after restarts. The operation is idempotent.
func (s *SchedulerStore) CompleteSchedule(ctx context.Context, groupName, name string) error {
	return s.MutateSchedule(ctx, groupName, name, func(schedule *Schedule) error {
		if schedule.CompletionDate != nil {
			return nil
		}
		now := time.Now().UTC()
		schedule.CompletionDate = &now
		// The completion marker is an internal field, so this write must
		// not stamp LastModificationDate the way a user-initiated update
		// does.
		return nil
	})
}

// TouchScheduleLastFired records the boundary of the most recently
// delivered occurrence so a restart cannot deliver it again. The marker
// only advances — an older boundary never overwrites a newer one — and,
// like the completion marker, the write does not stamp
// LastModificationDate.
func (s *SchedulerStore) TouchScheduleLastFired(ctx context.Context, groupName, name string, boundary time.Time) error {
	return s.MutateSchedule(ctx, groupName, name, func(schedule *Schedule) error {
		if schedule.LastFiredAt != nil && !schedule.LastFiredAt.Before(boundary) {
			return nil
		}
		fired := boundary
		schedule.LastFiredAt = &fired
		return nil
	})
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
	key := s.buildScheduleKey(groupName, name)
	// The record lock closes the resurrection window: without it a
	// MutateSchedule cycle that read before the delete could write the
	// record back after it.
	scheduleRecordWriteMu.Lock()
	if !s.schedulesStore.Exists(key) {
		scheduleRecordWriteMu.Unlock()
		return ErrScheduleNotFound
	}
	// Schedules carry no tag metadata — only schedule groups are taggable —
	// so deleting the record is the whole removal.
	err := s.schedulesStore.Delete(key)
	scheduleRecordWriteMu.Unlock()
	if err != nil {
		return err
	}
	return nil
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
		if namePrefix != "" && !strings.HasPrefix(s.Name, namePrefix) {
			return false
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
//
// The walk is unbounded — the same form the deletion cascade uses: the
// paged List clamps a zero MaxItems to the default page size and discards
// the marker, which would silently starve the enabled schedules beyond
// the first page of the firing sweep.
func (s *SchedulerStore) GetAllEnabledSchedules(ctx context.Context) ([]*Schedule, error) {
	return common.ListMatching[Schedule](s.schedulesStore, "", func(sch *Schedule) bool {
		if sch.State != ScheduleStateEnabled {
			return false
		}
		// A one-time schedule whose execution lifecycle has ended never
		// fires again — the persisted completion marker survives
		// restarts, unlike the engine's in-memory dedup map.
		if sch.CompletionDate != nil {
			return false
		}
		// at() expressions ignore StartDate/EndDate (AWS spec).
		// For rate()/cron() expressions, filter out permanently expired
		// schedules as defence-in-depth alongside the engine's dueBoundary
		// decision. This avoids fetching and evaluating schedules whose
		// EndDate has already passed.
		if !scheduleexpr.IsAtExpression(sch.ScheduleExpression) {
			if sch.EndDate != nil && time.Now().UTC().After(*sch.EndDate) {
				return false
			}
		}
		return true
	})
}

// EnsureDefaultGroup creates the default schedule group if it doesn't exist.
//
// Parameters:
//   - ctx: The context
//
// Returns:
//   - error: An error if creation fails
func (s *SchedulerStore) EnsureDefaultGroup(ctx context.Context) error {
	scheduleRecordWriteMu.Lock()
	defer scheduleRecordWriteMu.Unlock()
	arn := s.buildScheduleGroupARN(DefaultGroupName)
	if s.Exists(arn) {
		return nil
	}

	now := time.Now().UTC()
	group := &ScheduleGroup{
		Name:                 DefaultGroupName,
		ARN:                  arn,
		State:                ScheduleGroupStateActive,
		CreationDate:         now,
		LastModificationDate: now,
	}
	return s.Put(arn, group)
}
