package eventbridge

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// schedulerTickInterval is the granularity at which the scheduler evaluates
// rules. AWS EventBridge uses 1-minute minimum granularity for rate/cron.
const schedulerTickInterval = 1 * time.Minute

// lastFireTimes tracks the last time each scheduled rule was fired.
// Keyed by rule ARN. Accessed only by the single scheduler goroutine.
var lastFireTimes sync.Map

// startScheduler launches a background goroutine that ticks every minute and
// fires ENABLED rules whose ScheduleExpression matches the current time.
// It also starts a retention worker that purges expired archive events hourly.
func (s *EventsService) startScheduler() {
	ctx, cancel := context.WithCancel(context.Background())
	s.schedCancel = cancel
	s.schedWg.Add(2)
	go func() {
		defer s.schedWg.Done()
		s.runScheduler(ctx)
	}()
	go func() {
		defer s.schedWg.Done()
		s.runRetentionWorker(ctx)
	}()
}

func (s *EventsService) runScheduler(ctx context.Context) {
	// Align the first tick to the next UTC minute boundary so that
	// rate/cron expressions are evaluated on the same boundary as AWS
	// EventBridge, regardless of when the server started.
	now := time.Now().UTC()
	elapsed := time.Duration(now.Second())*time.Second + time.Duration(now.Nanosecond())
	wait := schedulerTickInterval - elapsed
	if wait > 0 {
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return
		}
	}

	ticker := time.NewTicker(schedulerTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.tickScheduledRules(ctx, now.UTC())
		}
	}
}

func (s *EventsService) tickScheduledRules(ctx context.Context, now time.Time) {
	s.eventsStores.Range(func(key, value any) bool {
		region, _ := key.(string)
		store, ok := value.(*eventsstore.EventsStore)
		if !ok {
			return true
		}
		s.fireScheduledRulesForRegion(ctx, region, store, now)
		return true
	})
}

func (s *EventsService) fireScheduledRulesForRegion(ctx context.Context, region string, store *eventsstore.EventsStore, now time.Time) {
	nextToken := ""
	for {
		result, err := store.ListRules(ctx, "", "", 1000, nextToken)
		if err != nil {
			logs.Warn("eventbridge scheduler: failed to list rules",
				logs.String("region", region),
				logs.Err(err))
			return
		}
		for _, rule := range result.Rules {
			if rule.State != eventsstore.RuleStateEnabled && rule.State != eventsstore.RuleStateEnabledWithAllCloudtrailManagementEvents {
				continue
			}
			if rule.ScheduleExpression == "" {
				continue
			}
			if !shouldFireSchedule(rule.ARN, rule.ScheduleExpression, now) {
				continue
			}
			s.fireScheduledRule(ctx, region, store, rule, now)
		}
		if result.NextToken == "" {
			break
		}
		nextToken = result.NextToken
	}
}

func (s *EventsService) fireScheduledRule(ctx context.Context, region string, store *eventsstore.EventsStore, rule *eventsstore.Rule, now time.Time) {
	logs.Debug("eventbridge scheduler: firing scheduled rule",
		logs.String("rule", rule.Name),
		logs.String("region", region),
		logs.String("schedule", rule.ScheduleExpression))

	event := &eventsstore.Event{
		ID:           generateEventID(),
		Version:      "0",
		DetailType:   "Scheduled Event",
		Source:       "aws.events",
		Account:      s.accountID,
		Time:         now,
		Region:       region,
		Detail:       map[string]interface{}{},
		EventBusName: rule.EventBusName,
	}

	if err := s.deliverEventWithStore(ctx, region, event, rule.EventBusName, store); err != nil {
		logs.Warn("eventbridge scheduler: failed to deliver scheduled event",
			logs.String("rule", rule.Name),
			logs.String("region", region),
			logs.Err(err))
	}
}

// shouldFireSchedule determines whether a schedule expression should fire at
// the given time. It uses lastFireTimes to ensure each rule fires at most once
// per evaluation.
func shouldFireSchedule(ruleARN, expr string, now time.Time) bool {
	if strings.HasPrefix(expr, "rate(") {
		return shouldFireRate(ruleARN, expr, now)
	}
	if strings.HasPrefix(expr, "cron(") {
		return shouldFireCron(ruleARN, expr, now)
	}
	return false
}

// shouldFireRate parses rate(value unit) expressions and fires when the
// elapsed time since the last fire is >= the rate duration.
func shouldFireRate(ruleARN, expr string, now time.Time) bool {
	inner := strings.TrimPrefix(expr, "rate(")
	inner = strings.TrimSuffix(inner, ")")
	inner = strings.TrimSpace(inner)

	parts := strings.Fields(inner)
	if len(parts) != 2 {
		return false
	}
	value, err := strconv.Atoi(parts[0])
	if err != nil || value <= 0 {
		return false
	}

	var duration time.Duration
	switch unit := strings.ToLower(parts[1]); unit {
	case "minute", "minutes":
		duration = time.Duration(value) * time.Minute
	case "hour", "hours":
		duration = time.Duration(value) * time.Hour
	case "day", "days":
		duration = time.Duration(value) * 24 * time.Hour
	case "week", "weeks":
		duration = time.Duration(value) * 7 * 24 * time.Hour
	default:
		return false
	}

	last, ok := getLastFire(ruleARN)
	if !ok {
		setLastFire(ruleARN, now)
		return true
	}
	if now.Sub(last) >= duration {
		setLastFire(ruleARN, now)
		return true
	}
	return false
}

// shouldFireCron evaluates a cron(...) expression through the shared
// schedule expression engine instead of a service-local matcher. The
// engine returns the first scheduled minute strictly after the previous
// minute; the rule fires when that minute is the current one, capped at
// once per minute. AWS cron format: cron(minutes hours day-of-month
// month day-of-week year) — including the L, W and # day wildcards.
func shouldFireCron(ruleARN, expr string, now time.Time) bool {
	next, err := scheduleexpr.NextExecutionTime(expr, now, time.Time{}, nil)
	if err != nil {
		return false
	}
	if !next.Equal(now.Truncate(time.Minute)) {
		return false
	}

	// Ensure we only fire once per minute.
	last, ok := getLastFire(ruleARN)
	if ok && last.Truncate(time.Minute).Equal(now.Truncate(time.Minute)) {
		return false
	}
	setLastFire(ruleARN, now)
	return true
}

func getLastFire(ruleARN string) (time.Time, bool) {
	v, ok := lastFireTimes.Load(ruleARN)
	if !ok {
		return time.Time{}, false
	}
	t, ok := v.(time.Time)
	return t, ok
}

func setLastFire(ruleARN string, t time.Time) {
	lastFireTimes.Store(ruleARN, t)
}

// retentionTickInterval controls how often the retention worker runs.
const retentionTickInterval = 1 * time.Hour

func (s *EventsService) runRetentionWorker(ctx context.Context) {
	ticker := time.NewTicker(retentionTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.purgeExpiredArchiveEvents(ctx, now.UTC())
		}
	}
}

func (s *EventsService) purgeExpiredArchiveEvents(ctx context.Context, now time.Time) {
	s.eventsStores.Range(func(key, value any) bool {
		store, ok := value.(*eventsstore.EventsStore)
		if !ok {
			return true
		}
		token := ""
		for {
			result, err := store.ListArchives(ctx, "", "", "", 1000, token)
			if err != nil {
				return true
			}
			for _, archive := range result.Archives {
				if archive.RetentionDays <= 0 {
					continue
				}
				cutoff := now.AddDate(0, 0, -int(archive.RetentionDays))
				if err := store.DeleteExpiredArchiveEvents(ctx, archive.Name, cutoff); err != nil {
					logs.Warn("eventbridge retention: failed to purge expired events",
						logs.String("archive", archive.Name),
						logs.Err(err))
				}
			}
			if result.NextToken == "" {
				break
			}
			token = result.NextToken
		}
		return true
	})
}
