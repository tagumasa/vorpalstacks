package scheduler

import (
	"context"
	"errors"
	"time"

	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// Due-schedule sweep: the region scan, the firing deduplication state, the
// boundary decision for rate()/cron()/at() expressions, and the DELETING
// group cascade.

// lastFiredEntry records the most recent fire of a schedule: when it
// fired and under which expression. The expression scopes the dedup —
// only boundaries of the same expression are suppressed.
type lastFiredEntry struct {
	firedAt time.Time
	expr    string
}

// lastFiredKey builds the deduplication key used by the lastFired map.
// Region is part of the key so that schedules sharing the same
// group/name across regions do not interfere with each other's
// deduplication state.
func lastFiredKey(region, groupName, name string) string {
	return region + "/" + groupName + "/" + name
}

func (e *Engine) checkSchedules() {
	if e.storageManager == nil {
		return
	}

	regions := e.storageManager.GetActiveRegions()

	// Collect active schedule keys across ALL regions so the lazy cleanup
	// pass at the end does not destroy entries belonging to a region that
	// has not yet been visited, or to a region whose schedules happen to
	// share a legacy (region-less) key with the current region.
	allActiveKeys := make(map[string]bool)

	for _, region := range regions {
		store := e.storeForRegion(region)
		if store == nil {
			continue
		}
		schedules, err := store.GetAllEnabledSchedules(e.ctx)
		if err != nil {
			logs.Debug("Failed to get enabled schedules", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}

		e.processDeletingGroups(e.ctx, store)

		now := time.Now().UTC()

		for _, schedule := range schedules {
			schedule.Region = region
			allActiveKeys[lastFiredKey(region, schedule.GroupName, schedule.Name)] = true
			boundary, due := e.dueBoundary(schedule, now)
			if due {
				// Provisionally reserve the dedup slot so a concurrent tick
				// does not double-fire while the goroutine is still in flight.
				// If executeSchedule fails or panics the goroutine releases
				// the reservation so the next tick can retry the schedule.
				dedupKey := lastFiredKey(region, schedule.GroupName, schedule.Name)
				e.lastFired.Store(dedupKey, lastFiredEntry{firedAt: now, expr: schedule.ScheduleExpression})
				e.wg.Add(1)
				go func(sch *schedulerstore.Schedule, key string, st *schedulerstore.SchedulerStore, fired time.Time) {
					defer e.wg.Done()
					defer func() {
						if r := recover(); r != nil {
							logs.Error("scheduler: panic executing schedule", logs.String("name", sch.Name), logs.Any("panic", r))
							// Release the dedup slot so the next tick can retry the schedule.
							e.lastFired.Delete(key)
						}
					}()
					select {
					case <-e.ctx.Done():
						return
					default:
						if err := e.executeSchedule(e.ctx, sch); err != nil {
							// Release the dedup slot so the next tick can retry the schedule.
							e.lastFired.Delete(key)
							return
						}
						// Persist the delivered boundary on the schedule record
						// so a restart cannot deliver it twice. The in-memory
						// slot then becomes redundant (the persisted marker
						// suppresses the boundary) and is released for the
						// next boundary.
						if err := st.TouchScheduleLastFired(e.ctx, sch.GroupName, sch.Name, fired); err != nil && !errors.Is(err, schedulerstore.ErrScheduleNotFound) {
							// A schedule deleted by its own completion has
							// nothing left to guard; any other failure keeps
							// the in-flight slot, which is then the only
							// duplicate guard until a restart.
							logs.Warn("Failed to persist the delivered boundary of schedule", logs.String("name", sch.Name), logs.Err(err))
							return
						}
						e.lastFired.Delete(key)
					}
				}(schedule, dedupKey, store, boundary)
			}
		}
	}

	// Remove lastFired entries for schedules that no longer exist so the
	// dedup map does not grow unbounded across schedule create/delete
	// cycles. Runs once after the full region sweep so that other regions'
	// entries are not destroyed while their sweep is still iterating.
	e.lastFired.Range(func(key, _ interface{}) bool {
		if k, ok := key.(string); ok {
			if !allActiveKeys[k] {
				e.lastFired.Delete(k)
			}
		}
		return true
	})
}

// resolveScheduleLocation returns the time.Location in which the schedule
// expression should be evaluated. When ScheduleExpressionTimezone is empty or
// invalid, UTC is used (matching the AWS default of UTC).
func resolveScheduleLocation(schedule *schedulerstore.Schedule) *time.Location {
	if schedule.ScheduleExpressionTimezone != "" {
		if loc, err := time.LoadLocation(schedule.ScheduleExpressionTimezone); err == nil {
			return loc
		}
		logs.Debug("Invalid schedule timezone, falling back to UTC",
			logs.String("schedule", schedule.Name),
			logs.String("timezone", schedule.ScheduleExpressionTimezone))
	}
	return time.UTC
}

// dueBoundary returns the schedule boundary this evaluation should deliver,
// if any. It is the decision core of the sweep: checkSchedules is its only
// production caller, and the firing tests wrap it in a local boolean helper.
func (e *Engine) dueBoundary(schedule *schedulerstore.Schedule, now time.Time) (time.Time, bool) {
	// Convert "now" to the schedule's evaluation timezone so that rate/cron/at
	// expressions are evaluated in the timezone the user configured.
	loc := resolveScheduleLocation(schedule)
	nowLocal := now.In(loc)

	// AWS: "When you configure a one-time schedule, EventBridge Scheduler
	// ignores the StartDate and EndDate you specify for the schedule."
	// Only rate()/cron() honour StartDate/EndDate.
	isAtExpression := scheduleexpr.IsAtExpression(schedule.ScheduleExpression)
	if !isAtExpression {
		if schedule.StartDate != nil && nowLocal.Before(schedule.StartDate.In(loc)) {
			return time.Time{}, false
		}
		if schedule.EndDate != nil && nowLocal.After(schedule.EndDate.In(loc)) {
			return time.Time{}, false
		}
	}

	boundary, elapsed := scheduleexpr.ElapsedExecutionTime(schedule.ScheduleExpression, nowLocal, schedule.CreationDate, schedule.StartDate, scheduleexpr.RateFiresAtAnchor)
	if !elapsed {
		return time.Time{}, false
	}

	// The persisted delivered-boundary marker is the source of truth across
	// restarts (the in-memory map dies with the process): an occurrence
	// already recorded as delivered must not fire again.
	if schedule.LastFiredAt != nil && !schedule.LastFiredAt.Before(boundary) {
		return time.Time{}, false
	}

	// Prevent duplicate firing: skip if already executed for this interval.
	// Key includes region so multi-region schedules do not share state.
	// dueBoundary is a pure predicate; the caller reserves the dedup slot
	// after this returns the boundary.
	key := lastFiredKey(schedule.Region, schedule.GroupName, schedule.Name)
	if last, ok := e.lastFired.Load(key); ok {
		// Only a fire under the same expression suppresses: an expression
		// change starts a new firing lifecycle (a completed one-time
		// schedule updated to a new past at() must fire again).
		if entry, ok := last.(lastFiredEntry); ok && entry.expr == schedule.ScheduleExpression && !entry.firedAt.Before(boundary) {
			return time.Time{}, false
		}
	}

	if schedule.FlexibleTimeWindow != nil && schedule.FlexibleTimeWindow.Mode == schedulerstore.FlexibleTimeWindowModeFlexible {
		maxWindow := 1
		if schedule.FlexibleTimeWindow.MaximumWindowInMinutes != nil {
			maxWindow = *schedule.FlexibleTimeWindow.MaximumWindowInMinutes
		}
		// AWS flexible time window starts at the scheduled time and extends
		// forward for MaximumWindowInMinutes. Execution never occurs before
		// the scheduled time (the resolved boundary already lies at or
		// before now); an occurrence whose window has closed is skipped,
		// not resurrected.
		windowEnd := boundary.Add(time.Duration(maxWindow) * time.Minute)
		if !nowLocal.Before(windowEnd) {
			return time.Time{}, false
		}
	}

	// The boundary lies at or before now: fire on this evaluation. A late
	// evaluation (ticker gap longer than the boundary interval) still fires
	// the pending boundary once instead of silently skipping it, and the
	// dedup checks above cap each boundary at one fire.
	return boundary, true
}

// processDeletingGroups completes the cascade deletion of schedule groups
// in the DELETING state: delete all member schedules, then purge the group
// record and its tags. A failure retries on the next sweep, which makes the
// cascade restart-safe.
func (e *Engine) processDeletingGroups(ctx context.Context, store *schedulerstore.SchedulerStore) {
	groups, err := store.ListDeletingScheduleGroups(ctx)
	if err != nil {
		logs.Debug("Failed to list deleting schedule groups", logs.String("error", err.Error()))
		return
	}
	for _, group := range groups {
		if err := store.DeleteSchedulesInGroup(ctx, group.Name); err != nil {
			logs.Warn("Failed to delete schedules in deleting group, will retry",
				logs.String("group", group.Name),
				logs.Err(err))
			continue
		}
		if err := store.PurgeDeletedScheduleGroup(ctx, group.Name); err != nil {
			logs.Warn("Failed to purge deleting schedule group, will retry",
				logs.String("group", group.Name),
				logs.Err(err))
		}
	}
}
