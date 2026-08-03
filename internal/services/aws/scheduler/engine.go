package scheduler

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// schedulerTickerInterval is the interval between schedule checks.
// Reduced to 1 second in test mode for faster test execution.
var schedulerTickerInterval = 1 * time.Minute

func init() {
	if os.Getenv("TEST_MODE") == "true" {
		schedulerTickerInterval = 1 * time.Second
	}
}

// Engine manages scheduled task execution for EventBridge Scheduler.
type Engine struct {
	storageManager *storage.RegionStorageManager
	accountID      string
	bus            eventbus.Bus
	stores         sync.Map // region → *schedulerstore.SchedulerStore

	// retryStores holds per-region RetryStores for persisted retry records.
	// Records survive server restarts so the at-least-once delivery
	// guarantee is maintained (S-B10).
	retryStores sync.Map // region → *schedulerstore.RetryStore

	// lastFired tracks the last execution time per schedule (key: groupName/name)
	// to prevent duplicate firing when the ticker polls multiple times within
	// the same execution window.
	lastFired sync.Map // string → time.Time

	running   bool
	runningMu sync.RWMutex
	stopChan  chan struct{}
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewEngine creates a new scheduler engine with the given store dependencies.
func NewEngine(
	storageManager *storage.RegionStorageManager,
	accountID string,
) *Engine {
	return &Engine{
		storageManager: storageManager,
		accountID:      accountID,
		stopChan:       make(chan struct{}),
	}
}

// SetEventBus injects the event bus for publishing scheduler lifecycle events.
func (e *Engine) SetEventBus(bus eventbus.Bus) {
	e.bus = bus
}

// Start starts the scheduler engine.
func (e *Engine) Start() error {
	e.runningMu.Lock()
	defer e.runningMu.Unlock()
	if e.running {
		return nil
	}
	e.running = true
	e.stopChan = make(chan struct{})
	e.ctx, e.cancel = context.WithCancel(context.Background())

	e.wg.Add(1)
	go e.run()

	logs.Debug("Scheduler engine started")
	return nil
}

// Stop stops the scheduler engine.
func (e *Engine) Stop() error {
	e.runningMu.Lock()
	defer e.runningMu.Unlock()
	if !e.running {
		return nil
	}
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	close(e.stopChan)
	e.wg.Wait()

	logs.Debug("Scheduler engine stopped")
	return nil
}

func (e *Engine) run() {
	defer e.wg.Done()
	defer func() { resilience.RecoverAndRestart("scheduler engine", &e.wg, e.run) }()

	ticker := time.NewTicker(schedulerTickerInterval)
	defer ticker.Stop()

	e.checkSchedules()

	for {
		select {
		case <-e.stopChan:
			return
		case <-ticker.C:
			e.checkSchedules()
			// Process pending retries from previous failed deliveries (S-B10).
			e.checkRetries()
		}
	}
}

// lastFiredKey builds the deduplication key used by the lastFired map.
// Region is part of the key so that schedules sharing the same
// group/name across regions do not interfere with each other's
// deduplication state (C-2).
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
	// share a legacy (region-less) key with the current region (C-2).
	allActiveKeys := make(map[string]bool)

	for _, region := range regions {
		storage, err := e.storageManager.GetStorage(region)
		if err != nil {
			logs.Debug("Failed to get storage for region", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}
		store := schedulerstore.NewSchedulerStore(storage, e.accountID, region)
		if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
			store = actual.(*schedulerstore.SchedulerStore)
		}
		schedules, err := store.GetAllEnabledSchedules(e.ctx)
		if err != nil {
			logs.Debug("Failed to get enabled schedules", logs.String("region", region), logs.String("error", err.Error()))
			continue
		}

		now := time.Now().UTC()

		for _, schedule := range schedules {
			schedule.Region = region
			allActiveKeys[lastFiredKey(region, schedule.GroupName, schedule.Name)] = true
			if e.shouldExecute(schedule, now) {
				// Provisionally reserve the dedup slot so a concurrent tick
				// does not double-fire while the goroutine is still in flight.
				// If executeSchedule fails or panics the goroutine releases
				// the reservation so the next tick can retry the schedule.
				dedupKey := lastFiredKey(region, schedule.GroupName, schedule.Name)
				e.lastFired.Store(dedupKey, now)
				e.wg.Add(1)
				go func(sch *schedulerstore.Schedule, key string) {
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
						}
					}
				}(schedule, dedupKey)
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

func (e *Engine) shouldExecute(schedule *schedulerstore.Schedule, now time.Time) bool {
	// Convert "now" to the schedule's evaluation timezone so that rate/cron/at
	// expressions are evaluated in the timezone the user configured (S-B3).
	loc := resolveScheduleLocation(schedule)
	nowLocal := now.In(loc)

	// AWS: "When you configure a one-time schedule, EventBridge Scheduler
	// ignores the StartDate and EndDate you specify for the schedule."
	// Only rate()/cron() honour StartDate/EndDate (S-B4).
	isAtExpression := strings.HasPrefix(schedule.ScheduleExpression, "at(")
	if !isAtExpression {
		if schedule.StartDate != nil && nowLocal.Before(schedule.StartDate.In(loc)) {
			return false
		}
		if schedule.EndDate != nil && nowLocal.After(schedule.EndDate.In(loc)) {
			return false
		}
	}

	nextTime, err := e.getNextExecutionTime(schedule, nowLocal)
	if err != nil {
		logs.Debug("Failed to calculate next execution time",
			logs.String("schedule", schedule.Name),
			logs.String("error", err.Error()))
		return false
	}

	// Prevent duplicate firing: skip if already executed for this interval.
	// Key includes region so multi-region schedules do not share state (C-2).
	// shouldExecute is a pure predicate; the caller reserves the dedup slot
	// after this returns true (H-2).
	key := lastFiredKey(schedule.Region, schedule.GroupName, schedule.Name)
	if last, ok := e.lastFired.Load(key); ok {
		if lastTime, ok := last.(time.Time); ok && !lastTime.Before(nextTime) {
			return false
		}
	}

	if schedule.FlexibleTimeWindow != nil && schedule.FlexibleTimeWindow.Mode == schedulerstore.FlexibleTimeWindowModeFlexible {
		maxWindow := 1
		if schedule.FlexibleTimeWindow.MaximumWindowInMinutes != nil {
			maxWindow = *schedule.FlexibleTimeWindow.MaximumWindowInMinutes
		}
		// AWS flexible time window starts at the scheduled time and extends
		// forward for MaximumWindowInMinutes. Execution must NOT occur before
		// the scheduled time.
		windowEnd := nextTime.Add(time.Duration(maxWindow) * time.Minute)
		if !nowLocal.Before(nextTime) && nowLocal.Before(windowEnd) {
			return true
		}
		return false
	}

	diff := nowLocal.Sub(nextTime)
	if diff >= 0 && diff < time.Minute {
		return true
	}
	return false
}

func (e *Engine) getNextExecutionTime(schedule *schedulerstore.Schedule, now time.Time) (time.Time, error) {
	expr := schedule.ScheduleExpression

	if strings.HasPrefix(expr, "at(") {
		matches := regexp.MustCompile(`^at\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\)$`).FindStringSubmatch(expr)
		if len(matches) == 2 {
			t, err := time.Parse(timeutils.ISO8601NoZFormat, matches[1])
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
	}

	if strings.HasPrefix(expr, "rate(") {
		matches := regexp.MustCompile(`^rate\((\d+)\s+(minute|minutes|hour|hours|day|days)\)$`).FindStringSubmatch(expr)
		if len(matches) == 3 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]

			var duration time.Duration
			switch strings.TrimSuffix(unit, "s") {
			case "minute":
				duration = time.Duration(value) * time.Minute
			case "hour":
				duration = time.Duration(value) * time.Hour
			case "day":
				duration = time.Duration(value) * 24 * time.Hour
			}

			creationTime := schedule.CreationDate
			if schedule.StartDate != nil {
				creationTime = *schedule.StartDate
			}
			elapsed := now.Sub(creationTime)
			periods := int(elapsed / duration)
			nextTime := creationTime.Add(time.Duration(periods) * duration)
			return nextTime, nil
		}
	}

	if strings.HasPrefix(expr, "cron(") {
		return e.parseCronNextTime(expr, now)
	}

	return time.Time{}, fmt.Errorf("unsupported schedule expression: %s", expr)
}

func (e *Engine) parseCronNextTime(expr string, now time.Time) (time.Time, error) {
	matches := regexp.MustCompile(`^cron\((.+)\)$`).FindStringSubmatch(expr)
	if len(matches) != 2 {
		return time.Time{}, fmt.Errorf("invalid cron expression: %s", expr)
	}

	fields := strings.Fields(matches[1])
	if len(fields) != 6 {
		return time.Time{}, fmt.Errorf("AWS cron expression must have 6 fields, got %d", len(fields))
	}

	// AWS cron supports L (last), W (nearest weekday), and # (nth occurrence)
	// wildcards that robfig/cron cannot parse. When any field contains these,
	// fall back to the manual evaluator (S-B9).
	needsManualEval := false
	for _, f := range fields {
		if strings.Contains(f, "L") || strings.Contains(f, "W") || strings.Contains(f, "#") {
			needsManualEval = true
			break
		}
	}

	if needsManualEval {
		return awsCronNextTime(fields, now)
	}

	standardCron := fmt.Sprintf("%s %s %s %s %s",
		convertAWSCronField(fields[0]),
		convertAWSCronField(fields[1]),
		convertAWSCronField(fields[2]),
		convertAWSCronField(fields[3]),
		convertAWSCronField(fields[4]),
	)

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(standardCron)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	nextTime := schedule.Next(now.Add(-time.Minute))

	// Honour the AWS cron year field (fields[5]). robfig/cron does not
	// support year, so we filter manually after calculating the next time.
	yearField := convertAWSCronField(fields[5])
	if yearField != "*" && yearField != "" {
		for y := nextTime.Year(); y <= now.Year()+100; y++ {
			if isYearAllowed(yearField, y) {
				if y == nextTime.Year() {
					return nextTime, nil
				}
				// Advance to January 1st of the next allowed year and
				// recalculate the cron next-time from there.
				baseline := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
				return schedule.Next(baseline.Add(-time.Minute)), nil
			}
		}
		return time.Time{}, fmt.Errorf("cron year field %q has no valid year in range", fields[5])
	}

	return nextTime, nil
}

// awsCronNextTime evaluates AWS cron expressions that contain L, W, or #
// wildcards by scanning forward minute-by-minute from the given time. This is
// necessary because robfig/cron does not support these AWS-specific wildcards.
// The scan is bounded to one year to prevent infinite loops on impossible
// expressions (e.g. February 31).
func awsCronNextTime(fields []string, now time.Time) (time.Time, error) {
	minuteMatcher, err := buildCronFieldMatcher(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron minutes field %q: %w", fields[0], err)
	}
	hourMatcher, err := buildCronFieldMatcher(fields[1], 0, 23)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron hours field %q: %w", fields[1], err)
	}
	domMatcher, err := buildDayOfMonthMatcher(fields[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron day-of-month field %q: %w", fields[2], err)
	}
	monthMatcher, err := buildCronFieldMatcher(convertAWSCronField(fields[3]), 1, 12)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron month field %q: %w", fields[3], err)
	}
	dowMatcher, err := buildDayOfWeekMatcher(fields[4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron day-of-week field %q: %w", fields[4], err)
	}
	yearMatcher, err := buildCronFieldMatcher(convertAWSCronField(fields[5]), 1970, 2199)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron year field %q: %w", fields[5], err)
	}

	domIsAny := fields[2] == "*" || fields[2] == "?"
	dowIsAny := fields[4] == "*" || fields[4] == "?"

	// Start scanning from the current minute boundary so that the
	// returned time can land within the current minute (shouldExecute
	// checks 0 <= diff < 1min). The legacy robfig path achieves this
	// via schedule.Next(now.Add(-time.Minute)).
	candidate := now.Truncate(time.Minute)
	limit := candidate.AddDate(1, 0, 0) // scan at most one year ahead

	for candidate.Before(limit) {
		if !minuteMatcher(candidate.Minute()) {
			candidate = candidate.Add(time.Minute)
			continue
		}
		if !hourMatcher(candidate.Hour()) {
			candidate = time.Date(candidate.Year(), candidate.Month(), candidate.Day(), candidate.Hour()+1, 0, 0, 0, candidate.Location())
			continue
		}
		if !monthMatcher(int(candidate.Month())) {
			candidate = time.Date(candidate.Year(), candidate.Month()+1, 1, 0, 0, 0, 0, candidate.Location())
			continue
		}
		if !yearMatcher(candidate.Year()) {
			// Year doesn't match — advance to Jan 1 of next year
			// instead of erroring. A cron like cron(0 9 * * ? 2026)
			// evaluated in 2025 should fire once 2026 arrives.
			candidate = time.Date(candidate.Year()+1, 1, 1, 0, 0, 0, 0, candidate.Location())
			continue
		}

		// Day matching: AWS uses OR semantics when both day-of-month and
		// day-of-week are specified (non-wildcard). When either is "?" or
		// "*", only the other is checked.
		dayMatches := false
		if domIsAny && dowIsAny {
			dayMatches = true
		} else if domIsAny {
			dayMatches = dowMatcher(candidate)
		} else if dowIsAny {
			dayMatches = domMatcher(candidate)
		} else {
			dayMatches = domMatcher(candidate) || dowMatcher(candidate)
		}

		if !dayMatches {
			candidate = time.Date(candidate.Year(), candidate.Month(), candidate.Day()+1, 0, 0, 0, 0, candidate.Location())
			continue
		}

		return candidate, nil
	}

	return time.Time{}, fmt.Errorf("no matching time found within one year for cron expression")
}

// buildCronFieldMatcher builds a matcher function for standard cron fields
// (minutes, hours, month, year). Supports: *, comma-separated values, ranges,
// step values (*/n, a-b/n, a/n).
func buildCronFieldMatcher(field string, min, max int) (func(int) bool, error) {
	if field == "*" {
		return func(v int) bool { return v >= min && v <= max }, nil
	}

	allowed := make(map[int]bool)
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "*" {
			for v := min; v <= max; v++ {
				allowed[v] = true
			}
			continue
		}
		if strings.Contains(part, "/") {
			slashParts := strings.SplitN(part, "/", 2)
			base := slashParts[0]
			step, err := strconv.Atoi(slashParts[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step: %s", slashParts[1])
			}
			start := min
			end := max
			if base != "*" {
				if strings.Contains(base, "-") {
					rangeParts := strings.SplitN(base, "-", 2)
					start, _ = strconv.Atoi(rangeParts[0])
					end, _ = strconv.Atoi(rangeParts[1])
				} else {
					start, _ = strconv.Atoi(base)
				}
			}
			for v := start; v <= end; v += step {
				if v >= min && v <= max {
					allowed[v] = true
				}
			}
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, err
			}
			for v := start; v <= end; v++ {
				if v >= min && v <= max {
					allowed[v] = true
				}
			}
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		if val >= min && val <= max {
			allowed[val] = true
		}
	}

	return func(v int) bool { return allowed[v] }, nil
}

// buildDayOfMonthMatcher handles the day-of-month field including the AWS W
// wildcard (nearest weekday) and L (last day of month).
func buildDayOfMonthMatcher(field string) (func(time.Time) bool, error) {
	converted := convertAWSCronField(field)
	if converted == "*" || converted == "?" {
		return func(t time.Time) bool { return true }, nil
	}

	hasL := strings.Contains(field, "L")
	hasW := strings.Contains(field, "W")

	if !hasL && !hasW {
		matcher, err := buildCronFieldMatcher(converted, 1, 31)
		if err != nil {
			return nil, err
		}
		return func(t time.Time) bool { return matcher(t.Day()) }, nil
	}

	// Field contains L or W — needs full time context for evaluation.
	return func(t time.Time) bool {
		for _, part := range strings.Split(field, ",") {
			part = strings.TrimSpace(part)
			if part == "*" || part == "?" {
				return true
			}
			if part == "L" {
				lastDay := lastDayOfMonth(t.Year(), t.Month())
				if t.Day() == lastDay {
					return true
				}
				continue
			}
			if strings.HasSuffix(part, "L") {
				// e.g. "L-3" means 3 days before the last day of month.
				offsetStr := strings.TrimSuffix(part, "L")
				offset := 0
				if offsetStr != "" {
					offset, _ = strconv.Atoi(offsetStr)
				}
				lastDay := lastDayOfMonth(t.Year(), t.Month())
				if t.Day() == lastDay+offset {
					return true
				}
				continue
			}
			if strings.HasSuffix(part, "W") {
				dayStr := strings.TrimSuffix(part, "W")
				targetDay, err := strconv.Atoi(dayStr)
				if err != nil {
					continue
				}
				if isNearestWeekday(t, targetDay) {
					return true
				}
				continue
			}
			val, err := strconv.Atoi(convertAWSCronField(part))
			if err == nil && val == t.Day() {
				return true
			}
		}
		return false
	}, nil
}

// buildDayOfWeekMatcher handles the day-of-week field including the AWS L
// (last occurrence of a weekday in the month) and # (nth occurrence) wildcards.
// AWS uses 1-7 (SUN=1) for both names and numerics; Go uses 0-6 (SUN=0).
func buildDayOfWeekMatcher(field string) (func(time.Time) bool, error) {
	if field == "*" || field == "?" {
		return func(t time.Time) bool { return true }, nil
	}

	hasL := strings.Contains(field, "L")
	hasHash := strings.Contains(field, "#")

	if !hasL && !hasHash {
		// Standard field — parse each comma-separated value and convert
		// from AWS convention (1=SUN..7=SAT) to Go convention (0=SUN..6=SAT).
		allowed := make(map[int]bool)
		for _, part := range strings.Split(field, ",") {
			part = strings.TrimSpace(part)
			if part == "*" || part == "?" {
				for i := 0; i < 7; i++ {
					allowed[i] = true
				}
				continue
			}
			// Handle step: base/step
			step := 1
			base := part
			if idx := strings.Index(part, "/"); idx != -1 {
				base = part[:idx]
				step, _ = strconv.Atoi(part[idx+1:])
				if step <= 0 {
					step = 1
				}
			}
			// Handle range: lo-hi
			if idx := strings.Index(base, "-"); idx != -1 {
				lo := parseAWSDoWValue(base[:idx])
				hi := parseAWSDoWValue(base[idx+1:])
				for d := lo; d <= hi; d += step {
					allowed[awsDowToGo(d)] = true
				}
			} else if base == "*" {
				for d := 1; d <= 7; d += step {
					allowed[awsDowToGo(d)] = true
				}
			} else {
				d := parseAWSDoWValue(base)
				allowed[awsDowToGo(d)] = true
			}
		}
		return func(t time.Time) bool { return allowed[int(t.Weekday())] }, nil
	}

	return func(t time.Time) bool {
		for _, part := range strings.Split(field, ",") {
			part = strings.TrimSpace(part)
			if part == "*" || part == "?" {
				return true
			}
			if strings.HasSuffix(part, "L") {
				// e.g. "6L" = last Friday of the month (AWS: 6=FRI).
				dayStr := strings.TrimSuffix(part, "L")
				dayNum := parseAWSDoWValue(dayStr)
				goDay := awsDowToGo(dayNum)
				if isLastWeekdayOfMonth(t, goDay) {
					return true
				}
				continue
			}
			if strings.Contains(part, "#") {
				// e.g. "3#2" = second Tuesday of the month (AWS: 3=TUE).
				hashParts := strings.SplitN(part, "#", 2)
				dayNum := parseAWSDoWValue(hashParts[0])
				weekNum, err := strconv.Atoi(hashParts[1])
				if err != nil || weekNum < 1 || weekNum > 5 {
					continue
				}
				goDay := awsDowToGo(dayNum)
				if isNthWeekdayOfMonth(t, goDay, weekNum) {
					return true
				}
				continue
			}
			val := parseAWSDoWValue(part)
			goDay := awsDowToGo(val)
			if int(t.Weekday()) == goDay {
				return true
			}
		}
		return false
	}, nil
}

// parseAWSDoWValue converts a single AWS day-of-week token (name or numeric)
// to the AWS numeric convention (1=SUN..7=SAT).
func parseAWSDoWValue(s string) int {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "SUN":
		return 1
	case "MON":
		return 2
	case "TUE":
		return 3
	case "WED":
		return 4
	case "THU":
		return 5
	case "FRI":
		return 6
	case "SAT":
		return 7
	}
	n, _ := strconv.Atoi(s)
	if n < 1 {
		return 1
	}
	if n > 7 {
		return 7
	}
	return n
}

// awsDowToGo converts an AWS day-of-week number (1=SUN..7=SAT) to Go's
// time.Weekday convention (0=SUN..6=SAT).
func awsDowToGo(awsDay int) int {
	if awsDay >= 1 && awsDay <= 7 {
		return awsDay - 1
	}
	return awsDay
}

// lastDayOfMonth returns the last calendar day of the given year/month.
func lastDayOfMonth(year int, month time.Month) int {
	firstOfNext := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return firstOfNext.AddDate(0, 0, -1).Day()
}

// isNearestWeekday checks whether t falls on the nearest weekday (Monday to
// Friday) to targetDay in t's month. AWS "W" rule: if targetDay is a Saturday,
// the nearest weekday is Friday (targetDay-1); if a Sunday, Monday
// (targetDay+1). Edge cases at month boundaries are handled.
func isNearestWeekday(t time.Time, targetDay int) bool {
	lastDay := lastDayOfMonth(t.Year(), t.Month())
	if targetDay < 1 || targetDay > lastDay {
		targetDay = lastDay
	}

	target := time.Date(t.Year(), t.Month(), targetDay, 0, 0, 0, 0, t.Location())
	wd := target.Weekday()

	var nearestDay int
	switch wd {
	case time.Saturday:
		if targetDay-1 >= 1 {
			nearestDay = targetDay - 1
		} else {
			nearestDay = targetDay + 2
		}
	case time.Sunday:
		if targetDay+1 <= lastDay {
			nearestDay = targetDay + 1
		} else {
			nearestDay = targetDay - 2
		}
	default:
		nearestDay = targetDay
	}

	return t.Day() == nearestDay
}

// isLastWeekdayOfMonth checks whether t is the last occurrence of the given
// weekday in its month.
func isLastWeekdayOfMonth(t time.Time, goDay int) bool {
	if int(t.Weekday()) != goDay {
		return false
	}
	nextWeek := t.AddDate(0, 0, 7)
	return nextWeek.Month() != t.Month()
}

// isNthWeekdayOfMonth checks whether t is the nth occurrence of the given
// weekday in its month.
func isNthWeekdayOfMonth(t time.Time, goDay, weekNum int) bool {
	if int(t.Weekday()) != goDay {
		return false
	}
	weekOfMonth := (t.Day()-1)/7 + 1
	return weekOfMonth == weekNum
}

var awsCronReplacements = map[string]string{
	"?":   "*",
	"SUN": "0", "MON": "1", "TUE": "2", "WED": "3", "THU": "4", "FRI": "5", "SAT": "6",
	"JAN": "1", "FEB": "2", "MAR": "3", "APR": "4", "MAY": "5", "JUN": "6",
	"JUL": "7", "AUG": "8", "SEP": "9", "OCT": "10", "NOV": "11", "DEC": "12",
}

func convertAWSCronField(field string) string {
	for old, new := range awsCronReplacements {
		field = strings.ReplaceAll(field, old, new)
	}
	return field
}

// isYearAllowed checks whether the given year satisfies the AWS cron year
// field expression. Supports: "*", single years ("2027"), comma-separated
// lists ("2026,2028"), ranges ("2025-2030"), and step values ("*/2").
func isYearAllowed(field string, year int) bool {
	if field == "*" || field == "" {
		return true
	}
	now := time.Now().Year()
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "*" {
			return true
		}
		if strings.Contains(part, "/") {
			slashParts := strings.SplitN(part, "/", 2)
			base := slashParts[0]
			step, _ := strconv.Atoi(slashParts[1])
			if step <= 0 {
				step = 1
			}
			startYear := now
			if base != "*" {
				startYear, _ = strconv.Atoi(base)
			}
			if year >= startYear && (year-startYear)%step == 0 {
				return true
			}
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			startYear, _ := strconv.Atoi(rangeParts[0])
			endYear, _ := strconv.Atoi(rangeParts[1])
			if year >= startYear && year <= endYear {
				return true
			}
			continue
		}
		y, err := strconv.Atoi(part)
		if err == nil && y == year {
			return true
		}
	}
	return false
}

func scheduleInput(target *schedulerstore.Target, scheduleName string) string {
	if target.Input != "" {
		return target.Input
	}
	msgPayload := map[string]interface{}{
		"schedule":  scheduleName,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	if msgBytes, err := json.Marshal(msgPayload); err == nil {
		return string(msgBytes)
	}
	return "{}"
}

func (e *Engine) executeSchedule(ctx context.Context, schedule *schedulerstore.Schedule) error {
	logs.Debug("Executing schedule",
		logs.String("name", schedule.Name),
		logs.String("group", schedule.GroupName))

	if schedule.Target == nil {
		logs.Debug("Schedule has no target", logs.String("schedule", schedule.Name))
		return nil
	}

	target := schedule.Target
	targetArn := target.Arn
	region := schedule.Region
	if region == "" {
		region = defaults.DefaultRegion
	}

	if e.bus != nil {
		input := scheduleInput(target, schedule.Name)
		schedEvt := &eventbus.ScheduleFiredEvent{
			ScheduleName:          schedule.Name,
			ScheduleArn:           schedule.ARN,
			GroupName:             schedule.GroupName,
			TargetArn:             targetArn,
			Input:                 input,
			ActionAfterCompletion: string(schedule.ActionAfterCompletion),
		}
		// Serialise the full target so the bus handler has access to all
		// sub-parameters (SqsParameters, KinesisParameters, etc.) without
		// needing to re-fetch the schedule from the store.
		if payloadBytes, err := json.Marshal(target); err == nil {
			schedEvt.TargetPayload = string(payloadBytes)
		}
		schedEvt.Region = region
		if err := e.bus.Publish(context.Background(), schedEvt); err != nil {
			logs.Warn("Failed to publish schedule fired event to event bus",
				logs.String("schedule", schedule.Name),
				logs.String("target", targetArn),
				logs.Err(err))
			return err
		}
	} else {
		// Direct delivery path: attempt delivery with retry (S-B10/S-B11).
		e.deliverWithRetry(ctx, schedule, target)
	}

	// NOTE: ActionAfterCompletion=DELETE is handled at delivery lifecycle
	// completion (success or retry exhaustion) by maybeAutoDelete, not here.
	// AWS deletes the schedule "shortly after its last target invocation",
	// i.e. after the retry policy terminates. Deleting here would orphan
	// retry records and break at-least-once delivery semantics.

	return nil
}

func (e *Engine) getStoreForSchedule(schedule *schedulerstore.Schedule) *schedulerstore.SchedulerStore {
	region := schedule.Region
	if region == "" {
		region = defaults.DefaultRegion
	}
	if cached, ok := e.stores.Load(region); ok {
		return cached.(*schedulerstore.SchedulerStore)
	}
	storage, err := e.storageManager.GetStorage(region)
	if err != nil {
		return nil
	}
	store := schedulerstore.NewSchedulerStore(storage, e.accountID, region)
	if actual, loaded := e.stores.LoadOrStore(region, store); loaded {
		return actual.(*schedulerstore.SchedulerStore)
	}
	return store
}

// maybeAutoDelete deletes the schedule when ActionAfterCompletion is DELETE.
// It is invoked at delivery lifecycle completion points (success or retry
// exhaustion) so that the schedule remains alive during the entire retry
// lifecycle, matching AWS EventBridge Scheduler semantics.
//
// AWS documentation: "the schedule is deleted shortly after its last target
// invocation" — meaning the final successful attempt or the last unsuccessful
// retry once the retry policy is exhausted.
//
// Idempotent: calling twice (e.g. success path racing with retry exhaustion)
// is safe — the second DeleteSchedule call returns NotFound and the error is
// logged at Debug level.
func (e *Engine) maybeAutoDelete(ctx context.Context, schedule *schedulerstore.Schedule) {
	if schedule == nil {
		return
	}
	if schedule.ActionAfterCompletion != schedulerstore.ActionAfterCompletionDelete {
		return
	}
	store := e.getStoreForSchedule(schedule)
	if store == nil {
		logs.Warn("No store available for auto-delete after completion",
			logs.String("schedule", schedule.Name),
			logs.String("group", schedule.GroupName),
			logs.String("region", schedule.Region))
		return
	}
	if err := store.DeleteSchedule(ctx, schedule.GroupName, schedule.Name); err != nil {
		logs.Debug("Failed to auto-delete schedule after completion",
			logs.String("schedule", schedule.Name),
			logs.String("group", schedule.GroupName),
			logs.String("error", err.Error()))
	}
}

// getRetryStore returns the RetryStore for the given region, creating it
// lazily on first access.
func (e *Engine) getRetryStore(region string) *schedulerstore.RetryStore {
	if region == "" {
		region = defaults.DefaultRegion
	}
	if cached, ok := e.retryStores.Load(region); ok {
		return cached.(*schedulerstore.RetryStore)
	}
	storage, err := e.storageManager.GetStorage(region)
	if err != nil {
		return nil
	}
	rs := schedulerstore.NewRetryStore(storage, region)
	actual, _ := e.retryStores.LoadOrStore(region, rs)
	return actual.(*schedulerstore.RetryStore)
}

// deliverToTarget dispatches a schedule delivery to the appropriate target
// type and returns an error on failure. This is the single entry point for
// all target deliveries, used by both the direct path and the bus path.
func (e *Engine) deliverToTarget(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	targetArn := target.Arn
	switch {
	case strings.Contains(targetArn, ":lambda:"):
		return e.invokeLambda(ctx, schedule, target)
	case strings.Contains(targetArn, ":sqs:"):
		return e.sendToSQS(ctx, schedule, target)
	case strings.Contains(targetArn, ":sns:"):
		return e.publishToSNS(ctx, schedule, target)
	case strings.Contains(targetArn, ":kinesis:"):
		return e.sendToKinesis(ctx, schedule, target)
	case strings.Contains(targetArn, ":states:"):
		return e.startStepFunctionExecution(ctx, schedule, target)
	case strings.Contains(targetArn, ":events:"):
		return e.sendToEventBridge(ctx, schedule, target)
	case strings.Contains(targetArn, ":logs:"):
		return e.sendToCloudWatchLogs(ctx, schedule, target)
	default:
		return fmt.Errorf("unsupported target type: %s", targetArn)
	}
}

// retryDefaults returns the effective MaximumRetryAttempts and
// MaximumEventAgeInSeconds for a target, applying AWS defaults when the
// RetryPolicy is nil or individual fields are unset.
// Defaults: 185 retries, 86400 seconds (24 hours).
func retryDefaults(target *schedulerstore.Target) (maxRetries int, maxAgeSeconds int) {
	maxRetries = 185
	maxAgeSeconds = 86400
	if target.RetryPolicy != nil {
		if target.RetryPolicy.MaximumRetryAttempts != nil && *target.RetryPolicy.MaximumRetryAttempts >= 0 {
			maxRetries = *target.RetryPolicy.MaximumRetryAttempts
		}
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil && *target.RetryPolicy.MaximumEventAgeInSeconds >= 60 {
			maxAgeSeconds = *target.RetryPolicy.MaximumEventAgeInSeconds
		}
	}
	return
}

// computeRetryBackoff calculates the delay before the next retry attempt
// using exponential backoff with jitter. The base interval doubles with each
// attempt, capped at 1 hour. Jitter is up to 50% of the base interval.
func computeRetryBackoff(attemptCount int) time.Duration {
	if attemptCount < 1 {
		attemptCount = 1
	}
	// Cap the shift to prevent int64 overflow. 1<<63 overflows to a
	// negative value, causing crand.Int to panic. Since the cap below
	// already limits base to 1 hour, attempts beyond ~10 are equivalent.
	shift := attemptCount - 1
	if shift > 12 {
		shift = 12 // 1<<12 seconds = 4096s ≈ 68min, capped to 1h below
	}
	base := time.Duration(1<<uint(shift)) * time.Second
	if base > time.Hour {
		base = time.Hour
	}
	// Add up to 50% jitter using crypto/rand for unpredictability.
	maxJitter := int64(base / 2)
	if maxJitter > 0 {
		n, err := crand.Int(crand.Reader, big.NewInt(maxJitter))
		if err == nil {
			return base + time.Duration(n.Int64())
		}
	}
	return base
}

// deliverWithRetry attempts delivery with an immediate retry, then persists a
// RetryRecord for background retries if the second attempt also fails.
// The RetryRecord survives server restarts so that at-least-once delivery
// is maintained (S-B10).
//
// ActionAfterCompletion=DELETE is handled at lifecycle completion:
//   - success on attempt 1 or 2 → maybeAutoDelete here
//   - maxRetries=0 DLQ route    → maybeAutoDelete here
//   - retry exhausted / event age exceeded → maybeAutoDelete in processRetryRecord
//
// See the AWS blog: "the schedule is deleted shortly after its last target
// invocation" — NOT at fire time.
func (e *Engine) deliverWithRetry(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	maxRetries, _ := retryDefaults(target)

	// First attempt (immediate).
	err := e.deliverToTarget(ctx, schedule, target)
	if err == nil {
		e.maybeAutoDelete(ctx, schedule)
		return
	}

	// AWS: MaximumRetryAttempts=0 means no retries — only the initial
	// attempt. Route to DLQ immediately on failure.
	if maxRetries == 0 {
		logs.Warn("Scheduler delivery failed, no retries configured",
			logs.String("schedule", schedule.Name),
			logs.String("target", target.Arn),
			logs.Err(err))
		e.routeToDLQ(ctx, schedule, target, "delivery failed (MaximumRetryAttempts=0)")
		// Retry lifecycle is now complete (no retries configured).
		e.maybeAutoDelete(ctx, schedule)
		return
	}

	logs.Warn("Scheduler delivery failed, retrying",
		logs.String("schedule", schedule.Name),
		logs.String("target", target.Arn),
		logs.Err(err))

	// Immediate retry (attempt 2).
	err = e.deliverToTarget(ctx, schedule, target)
	if err == nil {
		e.maybeAutoDelete(ctx, schedule)
		return
	}

	// Both immediate attempts failed — persist for background retry.
	region := schedule.Region
	if region == "" {
		region = defaults.DefaultRegion
	}
	now := time.Now()

	targetJSON, mErr := json.Marshal(target)
	if mErr != nil {
		logs.Error("Failed to serialise target for retry record",
			logs.String("schedule", schedule.Name),
			logs.Err(mErr))
		return
	}

	record := &schedulerstore.RetryRecord{
		ID:                    fmt.Sprintf("retry-%s-%d", uuidString(), now.UnixNano()),
		ScheduleName:          schedule.Name,
		GroupName:             schedule.GroupName,
		Region:                region,
		Target:                string(targetJSON),
		Input:                 scheduleInput(target, schedule.Name),
		AttemptCount:          2, // Two immediate attempts already made.
		CreatedAt:             now,
		NextAttemptAt:         now.Add(computeRetryBackoff(3)),
		ActionAfterCompletion: string(schedule.ActionAfterCompletion),
	}

	rs := e.getRetryStore(region)
	if rs == nil {
		logs.Error("No retry store available for region",
			logs.String("schedule", schedule.Name),
			logs.String("region", region))
		return
	}
	if sErr := rs.SaveRetryRecord(record); sErr != nil {
		logs.Error("Failed to persist retry record",
			logs.String("schedule", schedule.Name),
			logs.Err(sErr))
	}
}

// checkRetries scans all regions for due RetryRecords and attempts redelivery.
// Called on each ticker cycle by the background worker.
func (e *Engine) checkRetries() {
	if e.storageManager == nil {
		return
	}
	regions := e.storageManager.GetActiveRegions()
	now := time.Now()

	for _, region := range regions {
		rs := e.getRetryStore(region)
		if rs == nil {
			continue
		}
		due, err := rs.GetDueRetryRecords(now)
		if err != nil {
			logs.Debug("Failed to get due retry records",
				logs.String("region", region),
				logs.Err(err))
			continue
		}
		for _, record := range due {
			rec := record // capture for goroutine
			e.wg.Add(1)
			go func() {
				defer e.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						logs.Error("scheduler: panic processing retry record",
							logs.String("recordId", rec.ID),
							logs.Any("panic", r))
					}
				}()
				// Propagate e.ctx so retry processing stops when the engine
				// shuts down (H8). Previously e.ctx was used inside
				// processRetryRecord but the passed context was ignored,
				// causing retry goroutines to outlive the engine.
				e.processRetryRecord(e.ctx, rs, rec, now)
			}()
		}
	}
}

// processRetryRecord attempts redelivery of a single RetryRecord and handles
// the outcome: success -> delete, failure -> update next attempt or route to DLQ.
// The ctx parameter is propagated from checkRetries so that engine shutdown
// cancels in-flight retry processing (H8).
func (e *Engine) processRetryRecord(ctx context.Context, rs *schedulerstore.RetryStore, record *schedulerstore.RetryRecord, now time.Time) {
	var target schedulerstore.Target
	if err := json.Unmarshal([]byte(record.Target), &target); err != nil {
		logs.Error("Failed to deserialise target from retry record",
			logs.String("recordId", record.ID),
			logs.Err(err))
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		return
	}

	schedule := &schedulerstore.Schedule{
		Name:                  record.ScheduleName,
		GroupName:             record.GroupName,
		Region:                record.Region,
		Target:                &target,
		ActionAfterCompletion: schedulerstore.ActionAfterCompletion(record.ActionAfterCompletion),
	}

	// Verify the schedule still exists before retrying. With delayed
	// auto-deletion (AWS-compliant), the schedule remains alive during the
	// entire retry lifecycle, so this check passes for ActionAfterCompletion
	// = DELETE schedules and the retry continues. If the user manually
	// deletes the schedule mid-retry, this check fails and the retry record
	// is discarded — matching the user's intent to cancel.
	schedStore := e.getStoreForSchedule(schedule)
	if schedStore != nil {
		if _, err := schedStore.GetSchedule(ctx, record.GroupName, record.ScheduleName); err != nil {
			logs.Debug("Schedule no longer exists, discarding retry record",
				logs.String("schedule", record.ScheduleName),
				logs.String("recordId", record.ID))
			_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
			return
		}
	}

	maxRetries, maxAgeSeconds := retryDefaults(&target)

	// Check if MaximumEventAgeInSeconds has been exceeded.
	if now.Sub(record.CreatedAt) > time.Duration(maxAgeSeconds)*time.Second {
		logs.Warn("Retry record expired (MaximumEventAgeInSeconds exceeded)",
			logs.String("schedule", record.ScheduleName),
			logs.String("recordId", record.ID),
			logs.Int("attempts", record.AttemptCount))
		e.routeToDLQ(ctx, schedule, &target, "maximum event age exceeded")
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		// Retry lifecycle complete — auto-delete if configured.
		e.maybeAutoDelete(ctx, schedule)
		return
	}

	// Check if MaximumRetryAttempts has been exceeded.
	// AWS: MaximumRetryAttempts=N → N+1 total delivery attempts
	// (1 initial + N retries). Use > not >= so we don't lose one retry.
	if record.AttemptCount > maxRetries {
		logs.Warn("Retry policy exhausted",
			logs.String("schedule", record.ScheduleName),
			logs.String("recordId", record.ID),
			logs.Int("attempts", record.AttemptCount),
			logs.Int("maxRetries", maxRetries))
		e.routeToDLQ(ctx, schedule, &target, "retry policy exhausted")
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		// Retry lifecycle complete — auto-delete if configured.
		e.maybeAutoDelete(ctx, schedule)
		return
	}

	// Attempt delivery.
	err := e.deliverToTarget(ctx, schedule, &target)
	if err == nil {
		logs.Debug("Retry delivery succeeded",
			logs.String("schedule", record.ScheduleName),
			logs.String("recordId", record.ID),
			logs.Int("attempts", record.AttemptCount+1))
		_ = rs.DeleteRetryRecord(record.ID, record.NextAttemptAt)
		// Retry lifecycle complete — auto-delete if configured.
		e.maybeAutoDelete(ctx, schedule)
		return
	}

	// Delivery failed — persist the updated record with the new NextAttemptAt
	// BEFORE deleting the old key (keyed by old NextAttemptAt). This order
	// prevents data loss if Save fails after Delete succeeds (e.g. disk full).
	// If Save fails, the old record remains and will be retried at its
	// original NextAttemptAt, which is safe (at-least-once).
	oldNextAttempt := record.NextAttemptAt
	record.AttemptCount++
	record.NextAttemptAt = now.Add(computeRetryBackoff(record.AttemptCount))
	if sErr := rs.SaveRetryRecord(record); sErr != nil {
		logs.Error("Failed to update retry record (old record retained)",
			logs.String("recordId", record.ID),
			logs.Err(sErr))
		return
	}
	_ = rs.DeleteRetryRecord(record.ID, oldNextAttempt)
}

// routeToDLQ sends a failed delivery to the DeadLetterConfig ARN. If no DLQ
// is configured, the message is discarded with an error log (AWS-compliant).
func (e *Engine) routeToDLQ(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target, reason string) {
	if target.DeadLetterConfig == nil || target.DeadLetterConfig.Arn == "" {
		logs.Error("Scheduler delivery permanently failed, no DLQ configured — discarding",
			logs.String("schedule", schedule.Name),
			logs.String("target", target.Arn),
			logs.String("reason", reason))
		return
	}

	dlqArn := target.DeadLetterConfig.Arn
	message := scheduleInput(target, schedule.Name)

	logs.Warn("Routing failed schedule delivery to DLQ",
		logs.String("schedule", schedule.Name),
		logs.String("dlqArn", dlqArn),
		logs.String("reason", reason))

	if strings.Contains(dlqArn, ":sqs:") {
		sqsInvoker := e.bus.SQSInvoker()
		if sqsInvoker == nil {
			logs.Error("SQS invoker not available for DLQ delivery",
				logs.String("dlqArn", dlqArn))
			return
		}
		queueName := svcarn.ExtractQueueNameFromARN(dlqArn)
		queueURL, qErr := sqsInvoker.GetQueueByName(ctx, queueName)
		if qErr != nil {
			logs.Error("Failed to resolve DLQ queue URL",
				logs.String("dlqArn", dlqArn),
				logs.Err(qErr))
			return
		}
		// FIFO queues require MessageGroupId. Propagate from the
		// target's SqsParameters if available; otherwise, when the
		// destination queue ends with the FIFO suffix, fall back to
		// the schedule name so the SendMessage call satisfies the
		// SQS FIFO requirement.
		sendOpts := eventbus.SQSSendOptions{}
		if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
			sendOpts.MessageGroupID = target.SqsParameters.MessageGroupId
		} else if strings.HasSuffix(queueName, ".fifo") {
			sendOpts.MessageGroupID = schedule.Name
		}
		if _, _, err := sqsInvoker.SendMessage(ctx, queueURL, message, sendOpts); err != nil {
			logs.Error("Failed to send to DLQ",
				logs.String("dlqArn", dlqArn),
				logs.Err(err))
		}
	} else if strings.Contains(dlqArn, ":sns:") {
		dlqTarget := &schedulerstore.Target{Arn: dlqArn, Input: message}
		if err := e.publishToSNS(ctx, schedule, dlqTarget); err != nil {
			logs.Error("Failed to publish to SNS DLQ",
				logs.String("dlqArn", dlqArn),
				logs.Err(err))
		}
	} else {
		logs.Error("Unsupported DLQ type",
			logs.String("dlqArn", dlqArn))
	}
}

// uuidString returns a UUID-like unique string for retry record IDs.
// Uses crypto/rand for uniqueness without external UUID dependency.
func uuidString() string {
	b := make([]byte, 16)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (e *Engine) invokeLambda(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for Lambda invocation", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	input := scheduleInput(target, schedule.Name)

	functionName := svcarn.ExtractFunctionNameFromARN(target.Arn)
	if functionName == "" {
		logs.Debug("Failed to extract function name from ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid Lambda ARN: %s", target.Arn)
	}

	logs.Debug("Invoking Lambda for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName))

	statusCode, _, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, target.Arn, []byte(input))
	if err != nil {
		logs.Debug("Failed to invoke Lambda",
			logs.String("schedule", schedule.Name),
			logs.String("function", functionName),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Lambda invocation completed",
		logs.String("schedule", schedule.Name),
		logs.String("function", functionName),
		logs.Int("statusCode", int(statusCode)))
	return nil
}

func (e *Engine) sendToSQS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for SQS delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	sqsInvoker := e.bus.SQSInvoker()
	if sqsInvoker == nil {
		return fmt.Errorf("SQS invoker not available")
	}

	queueName := svcarn.ExtractQueueNameFromARN(target.Arn)
	if queueName == "" {
		logs.Debug("Invalid SQS ARN", logs.String("arn", target.Arn))
		return fmt.Errorf("invalid SQS ARN: %s", target.Arn)
	}

	queueURL, qErr := sqsInvoker.GetQueueByName(ctx, queueName)
	if qErr != nil {
		logs.Debug("SQS queue not found", logs.String("queue", queueName), logs.Err(qErr))
		return qErr
	}

	messageBody := scheduleInput(target, schedule.Name)

	// Honour SqsParameters.MessageGroupId for FIFO queues. AWS requires
	// MessageGroupId when the target is a FIFO queue. Fall back to the
	// schedule name if not explicitly set, matching routeToDLQ behaviour (M10).
	sendOpts := eventbus.SQSSendOptions{}
	if target.SqsParameters != nil && target.SqsParameters.MessageGroupId != "" {
		sendOpts.MessageGroupID = target.SqsParameters.MessageGroupId
	} else if strings.HasSuffix(queueName, ".fifo") {
		sendOpts.MessageGroupID = schedule.Name
	}

	logs.Debug("Sending to SQS for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("queue", queueName))

	if _, _, err := sqsInvoker.SendMessage(ctx, queueURL, messageBody, sendOpts); err != nil {
		logs.Debug("Failed to send to SQS",
			logs.String("schedule", schedule.Name),
			logs.String("queue", queueName),
			logs.Err(err))
		return err
	}
	return nil
}

func (e *Engine) publishToSNS(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for SNS delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	snsInvoker := e.bus.SNSInvoker()
	if snsInvoker == nil {
		return fmt.Errorf("SNS invoker not available")
	}

	message := scheduleInput(target, schedule.Name)

	// Delegate to the SNS service's Publish API. This ensures the SNS
	// service handles topic policy evaluation, subscription filtering,
	// message persistence, fan-out to all subscription endpoints (SQS,
	// Lambda, HTTP, etc.), and EventBridge bus event publication (H7).
	messageID, err := snsInvoker.PublishToTopic(ctx, target.Arn, message, "", nil)
	if err != nil {
		logs.Debug("Failed to publish to SNS topic",
			logs.String("schedule", schedule.Name),
			logs.String("topicArn", target.Arn),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("SNS delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("topic", target.Arn),
		logs.String("messageId", messageID))
	return nil
}

func (e *Engine) sendToCloudWatchLogs(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping CloudWatch Logs delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	logGroup := svcarn.ExtractLogGroupNameFromARN(target.Arn)
	if logGroup == "" {
		logs.Debug("Failed to extract log group from CloudWatch Logs ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid CloudWatch Logs ARN: %s", target.Arn)
	}

	_, _, region, _, resource := svcarn.SplitARN(target.Arn)
	var logStream string
	if idx := strings.LastIndex(resource, ":log-stream:"); idx != -1 {
		logStream = resource[idx+12:]
	} else {
		logStream = fmt.Sprintf("scheduler-%s", schedule.Name)
	}

	message := scheduleInput(target, schedule.Name)

	evt := &eventbus.CloudWatchLogsPutEvent{
		LogGroup:  logGroup,
		LogStream: logStream,
		LogEvents: []eventbus.LogEntry{
			{Timestamp: time.Now().UnixMilli(), Message: message},
		},
	}
	evt.Region = region
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to deliver schedule to CloudWatch Logs",
			logs.String("schedule", schedule.Name),
			logs.String("logGroup", logGroup),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to CloudWatch Logs",
		logs.String("schedule", schedule.Name),
		logs.String("logGroup", logGroup))
	return nil
}

func (e *Engine) sendToKinesis(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured for Kinesis delivery", logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	kinesisInvoker := e.bus.KinesisInvoker()
	if kinesisInvoker == nil {
		logs.Debug("Kinesis invoker not available", logs.String("schedule", schedule.Name))
		return fmt.Errorf("Kinesis invoker not available")
	}

	// Extract stream name from ARN: arn:aws:kinesis:<region>:<account>:stream/<name>
	_, _, kRegion, _, resource := svcarn.SplitARN(target.Arn)
	streamName := ""
	if idx := strings.LastIndex(resource, "/"); idx != -1 {
		streamName = resource[idx+1:]
	} else {
		streamName = resource
	}
	if streamName == "" {
		logs.Debug("Failed to extract stream name from Kinesis ARN",
			logs.String("schedule", schedule.Name),
			logs.String("arn", target.Arn))
		return fmt.Errorf("invalid Kinesis ARN: %s", target.Arn)
	}

	// Use PartitionKey from KinesisParameters if provided, otherwise fall
	// back to the schedule name (AWS uses a similar default behaviour).
	partitionKey := schedule.Name
	if target.KinesisParameters != nil && target.KinesisParameters.PartitionKey != "" {
		partitionKey = target.KinesisParameters.PartitionKey
	}

	data := []byte(scheduleInput(target, schedule.Name))

	logs.Debug("Sending to Kinesis for schedule",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))

	if _, err := kinesisInvoker.PutRecord(ctx, streamName, partitionKey, data); err != nil {
		logs.Debug("Failed to send to Kinesis",
			logs.String("schedule", schedule.Name),
			logs.String("stream", streamName),
			logs.Err(err))
		return err
	}

	// Propagate region from the target ARN if the schedule does not carry one.
	if schedule.Region == "" && kRegion != "" {
		schedule.Region = kRegion
	}

	logs.Debug("Kinesis delivery completed",
		logs.String("schedule", schedule.Name),
		logs.String("stream", streamName))
	return nil
}

func (e *Engine) startStepFunctionExecution(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping Step Functions delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	_, _, smRegion, _, _ := svcarn.SplitARN(target.Arn)
	if smRegion == "" {
		smRegion = schedule.Region
	}

	input := scheduleInput(target, schedule.Name)

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: target.Arn,
		Input:           input,
	}
	evt.Region = smRegion
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to start Step Functions execution from schedule",
			logs.String("schedule", schedule.Name),
			logs.String("stateMachineArn", target.Arn),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to Step Functions",
		logs.String("schedule", schedule.Name),
		logs.String("stateMachineArn", target.Arn))
	return nil
}

func (e *Engine) sendToEventBridge(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) error {
	if e.bus == nil {
		logs.Debug("event bus not configured, skipping EventBridge delivery",
			logs.String("schedule", schedule.Name))
		return fmt.Errorf("event bus not configured")
	}

	_, _, ebRegion, _, resource := svcarn.SplitARN(target.Arn)
	if ebRegion == "" {
		ebRegion = schedule.Region
	}

	eventBusName := "default"
	if idx := strings.LastIndex(resource, ":event-bus/"); idx != -1 {
		eventBusName = resource[idx+len(":event-bus/"):]
	}

	input := scheduleInput(target, schedule.Name)

	evt := &eventbus.EventBridgePutEventsEvent{
		EventBusName: eventBusName,
		Input:        input,
	}
	// Populate DetailType and Source from EventBridgeParameters so that
	// EventBridge rules can match on them (S-B5).
	if target.EventBridgeParameters != nil {
		evt.DetailType = target.EventBridgeParameters.DetailType
		evt.Source = target.EventBridgeParameters.Source
	}
	evt.Region = ebRegion
	evt.AccountID = e.accountID

	if err := e.bus.Publish(ctx, evt); err != nil {
		logs.Debug("Failed to deliver schedule to EventBridge",
			logs.String("schedule", schedule.Name),
			logs.String("eventBus", eventBusName),
			logs.String("error", err.Error()))
		return err
	}

	logs.Debug("Schedule delivered to EventBridge",
		logs.String("schedule", schedule.Name),
		logs.String("eventBus", eventBusName))
	return nil
}
