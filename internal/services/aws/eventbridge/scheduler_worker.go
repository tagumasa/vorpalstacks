package eventbridge

import (
	"context"
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
// The map is the primary dedup state; after a restart it is re-seeded
// from each rule's persisted LastFiredAt marker (see seedLastFire), and
// every successful fire is persisted back so the boundary survives the
// next restart.
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
			// A restart empties the in-memory dedup cache; re-seed it from
			// the persisted marker so the boundary fired just before the
			// restart is not fired again.
			seedLastFire(rule)
			if !shouldFireSchedule(rule.ARN, rule.ScheduleExpression, now, rule.CreatedAt) {
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
		return
	}

	// Persist the fired boundary after the successful delivery so a
	// restart does not fire it again. A failed delivery keeps the
	// in-memory reservation for this process only: the boundary is
	// retried after a restart rather than lost.
	if boundary, ok := getLastFire(rule.ARN); ok {
		if err := store.TouchRuleLastFired(ctx, rule.EventBusName, rule.Name, boundary); err != nil {
			logs.Debug("eventbridge scheduler: failed to persist the fired boundary",
				logs.String("rule", rule.Name),
				logs.String("region", region),
				logs.Err(err))
		}
	}
}

// shouldFireSchedule determines whether a schedule expression should fire at
// the given time. It uses lastFireTimes to ensure each rule fires at most once
// per evaluation. creationTime anchors rate() period boundaries.
func shouldFireSchedule(ruleARN, expr string, now, creationTime time.Time) bool {
	boundary, ok := scheduleexpr.ElapsedExecutionTime(expr, now, creationTime, nil)
	if !ok {
		return false
	}
	// Fire the latest elapsed boundary exactly once: an evaluation that
	// arrives late (a ticker gap longer than the boundary interval)
	// still fires the pending boundary instead of skipping it silently.
	// EventBridge rate rules do not fire immediately on creation — the
	// first fire happens one full interval after the rule was created —
	// and the boundaries stay pinned to the creation time instead of
	// drifting forward with each fire.
	if last, ok := getLastFire(ruleARN); ok && !boundary.After(last) {
		return false
	}
	setLastFire(ruleARN, boundary)
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

// seedLastFire re-seeds the in-memory dedup cache from a rule's
// persisted fire marker. The persisted marker only ever advances the
// cached value, never regresses it.
func seedLastFire(rule *eventsstore.Rule) {
	if rule.LastFiredAt.IsZero() {
		return
	}
	if last, ok := getLastFire(rule.ARN); !ok || rule.LastFiredAt.After(last) {
		setLastFire(rule.ARN, rule.LastFiredAt)
	}
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
