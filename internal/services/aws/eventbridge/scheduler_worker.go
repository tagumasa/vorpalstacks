package eventbridge

import (
	"context"
	"errors"
	"sync"
	"time"

	"vorpalstacks/internal/common/scheduleexpr"
	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// schedulerTickInterval is the granularity at which the scheduler evaluates
// rules. AWS EventBridge uses 1-minute minimum granularity for rate/cron.
const schedulerTickInterval = 1 * time.Minute

// scheduleFireDedup tracks the last fired boundary per scheduled rule so
// each boundary fires exactly once. The zero value is ready to use; the
// EventsService owns the single instance, keeping scheduler state off the
// package level. After a restart the map is re-seeded from each rule's
// persisted LastFiredAt marker (seedLastFire), and every successful fire
// is persisted back so the boundary survives the next restart.
type scheduleFireDedup struct {
	last sync.Map // rule ARN → time.Time
}

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
			s.tickScheduledRulesGuarded(ctx, now.UTC())
		}
	}
}

// tickScheduledRulesGuarded runs one scheduler tick behind a panic boundary:
// a panic while decoding a rule record or evaluating an expression must not
// kill the scheduler goroutine for the rest of the process lifetime — the
// tick is reported and the next minute's tick proceeds. The per-target
// delivery goroutines carry their own guards.
func (s *EventsService) tickScheduledRulesGuarded(ctx context.Context, now time.Time) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("eventbridge scheduler: panic during tick",
				logs.Any("panic", r),
				logs.String("tick", now.Format(time.RFC3339)))
		}
	}()
	s.tickScheduledRules(ctx, now)
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
			s.fireDedup.seedLastFire(rule)
			if !s.fireDedup.shouldFireSchedule(rule.ARN, rule.ScheduleExpression, now, rule.CreatedAt) {
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
		Resources:    []string{rule.ARN},
		Detail:       map[string]interface{}{},
		EventBusName: rule.EventBusName,
	}

	// A scheduled fire is a bus event: "EventBridge itself emits the
	// following events. These events are automatically sent to the default
	// event bus as with any other AWS service." and "EventBridge sends the
	// following schedule events to the default event bus." (the EventBridge
	// events detail reference) — the reference's sample envelope carries
	// the firing rule's ARN in resources and an empty detail object. The
	// fire therefore archives (the archives user guide's stated filter is
	// the event pattern alone, never the ingress path, so a bus archive
	// whose pattern matches aws.events captures scheduled fires too) and
	// fans out through bus-wide matching to every enabled rule whose
	// pattern matches. Sibling scheduled rules carry empty patterns and
	// never match bus events (deliverEventToBusRules skips them), so one
	// boundary fire still reaches exactly its own rule's targets plus the
	// pattern matches — the firing rule's own targets are dispatched
	// directly beside the fan-out.
	if err := s.deliverEvent(ctx, store, event, rule.EventBusName, region, 0); err != nil {
		logs.Warn("eventbridge scheduler: failed to deliver scheduled event",
			logs.String("rule", rule.Name),
			logs.String("region", region),
			logs.Err(err))
		return
	}
	if err := s.dispatchRuleTargets(ctx, region, event, rule, store, 0); err != nil {
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
	if boundary, ok := s.fireDedup.getLastFire(rule.ARN); ok {
		if err := store.TouchRuleLastFired(ctx, rule.EventBusName, rule.Name, boundary); err != nil {
			logs.Debug("eventbridge scheduler: failed to persist the fired boundary",
				logs.String("rule", rule.Name),
				logs.String("region", region),
				logs.Err(err))
		}
	}
}

// unevaluableScheduleLogged records the rule-and-expression pairs whose
// schedule expression failed evaluation, so the warning is logged once per
// pair instead of on every tick. Keyed with the expression so a corrected
// expression that later fails again logs anew.
var unevaluableScheduleLogged sync.Map

// shouldFireSchedule determines whether a schedule expression should fire at
// the given time. It uses the dedup cache to ensure each rule fires at most
// once per evaluation. creationTime anchors rate() period boundaries.
func (d *scheduleFireDedup) shouldFireSchedule(ruleARN, expr string, now, creationTime time.Time) bool {
	boundary, ok := scheduleexpr.ElapsedExecutionTime(expr, now, creationTime, nil, scheduleexpr.RateFiresAfterFirstInterval)
	if !ok {
		// ok=false covers both an expression that cannot be evaluated and
		// a valid expression whose next boundary has not elapsed yet; only
		// the former is reported — and once per rule-and-expression pair,
		// not per tick. Only reachable for records that predate
		// schedule-expression validation.
		if !scheduleexpr.ValidateExpression(expr) {
			if _, seen := unevaluableScheduleLogged.LoadOrStore(ruleARN+"\x00"+expr, true); !seen {
				logs.Warn("eventbridge scheduler: schedule expression failed evaluation; the rule will never fire",
					logs.String("ruleArn", ruleARN),
					logs.String("schedule", expr))
			}
		}
		return false
	}
	// Fire the latest elapsed boundary exactly once: an evaluation that
	// arrives late (a ticker gap longer than the boundary interval)
	// still fires the pending boundary instead of skipping it silently.
	// EventBridge rate rules do not fire immediately on creation — the
	// first fire happens one full interval after the rule was created —
	// and the boundaries stay pinned to the creation time instead of
	// drifting forward with each fire.
	if last, ok := d.getLastFire(ruleARN); ok && !boundary.After(last) {
		return false
	}
	d.setLastFire(ruleARN, boundary)
	return true
}

func (d *scheduleFireDedup) getLastFire(ruleARN string) (time.Time, bool) {
	v, ok := d.last.Load(ruleARN)
	if !ok {
		return time.Time{}, false
	}
	t, ok := v.(time.Time)
	return t, ok
}

func (d *scheduleFireDedup) setLastFire(ruleARN string, t time.Time) {
	d.last.Store(ruleARN, t)
}

// deleteLastFire drops a rule's dedup entry when the rule (or its bus) is
// deleted, so a re-created rule with the same ARN starts unreserved.
func (d *scheduleFireDedup) deleteLastFire(ruleARN string) {
	d.last.Delete(ruleARN)
}

// seedLastFire re-seeds the in-memory dedup cache from a rule's
// persisted fire marker. The persisted marker only ever advances the
// cached value, never regresses it.
func (d *scheduleFireDedup) seedLastFire(rule *eventsstore.Rule) {
	if rule.LastFiredAt.IsZero() {
		return
	}
	if last, ok := d.getLastFire(rule.ARN); !ok || rule.LastFiredAt.After(last) {
		d.setLastFire(rule.ARN, rule.LastFiredAt)
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
			s.purgeExpiredReplays(ctx, now.UTC())
		}
	}
}

// purgeExpiredReplays deletes replay records older than the documented
// retention: "EventBridge deletes replays after 90 days." (the archives user
// guide), counted from the record's creation. A record without a creation
// stamp predates the field and is left in place — the sweep never fabricates
// a timestamp to delete by. A record that vanishes concurrently is a
// legitimate outcome of a user delete racing the sweep.
func (s *EventsService) purgeExpiredReplays(ctx context.Context, now time.Time) {
	s.eventsStores.Range(func(key, value any) bool {
		store, ok := value.(*eventsstore.EventsStore)
		if !ok {
			return true
		}
		cutoff := now.AddDate(0, 0, -eventsstore.ReplayRetentionDays)
		nextToken := ""
		for {
			result, err := store.ListReplays(ctx, "", "", "", eventsstore.ListLimitMaximum, nextToken)
			if err != nil {
				logs.Warn("eventbridge replay retention: failed to list replays", logs.Err(err))
				return true
			}
			for _, replay := range result.Replays {
				if replay.CreatedAt.IsZero() || replay.CreatedAt.After(cutoff) {
					continue
				}
				if err := store.DeleteReplay(ctx, replay.Name); err != nil && !errors.Is(err, eventsstore.ErrReplayNotFound) {
					logs.Warn("eventbridge replay retention: failed to delete replay",
						logs.String("replay", replay.Name),
						logs.Err(err))
				}
			}
			if result.NextToken == "" {
				break
			}
			nextToken = result.NextToken
		}
		return true
	})
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
