package scheduler

import (
	"context"
	crand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/core/logs"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Retry and dead-letter lifecycle: retry policy defaults, backoff
// computation, persisted background retries, and DLQ routing.

// retryDefaults returns the effective MaximumRetryAttempts and
// MaximumEventAgeInSeconds for a target, applying the model's maximum
// values as defaults when the RetryPolicy is nil or individual fields are
// unset.
func retryDefaults(target *schedulerstore.Target) (maxRetries int, maxAgeSeconds int) {
	maxRetries = MaxRetryPolicyAttempts
	maxAgeSeconds = MaxRetryPolicyEventAgeSeconds
	if target.RetryPolicy != nil {
		if target.RetryPolicy.MaximumRetryAttempts != nil && *target.RetryPolicy.MaximumRetryAttempts >= 0 {
			maxRetries = *target.RetryPolicy.MaximumRetryAttempts
		}
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil && *target.RetryPolicy.MaximumEventAgeInSeconds >= MinRetryPolicyEventAgeSeconds {
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
// is maintained.
//
// Post-execution actions are handled at lifecycle completion via
// maybeActionAfterCompletion:
//   - success on attempt 1 or 2 → applied here
//   - maxRetries=0 DLQ route    → applied here
//   - retry record cannot be persisted → DLQ route applied here: without a
//     record no background retry will ever run, so the delivery lifecycle
//     ends at that point
//   - retry exhausted / event age exceeded → applied in processRetryRecord
//
// For ActionAfterCompletion=DELETE see the AWS blog: "the schedule is
// deleted shortly after its last target invocation" — NOT at fire time;
// for a one-time schedule with no action the completion marker ends its
// firing lifecycle.
func (e *Engine) deliverWithRetry(ctx context.Context, schedule *schedulerstore.Schedule, target *schedulerstore.Target) {
	maxRetries, _ := retryDefaults(target)

	// First attempt (immediate).
	err := e.deliverToTarget(ctx, schedule, target)
	if err == nil {
		e.maybeActionAfterCompletion(ctx, schedule)
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
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	logs.Warn("Scheduler delivery failed, retrying",
		logs.String("schedule", schedule.Name),
		logs.String("target", target.Arn),
		logs.Err(err))

	// Immediate retry (attempt 2).
	err = e.deliverToTarget(ctx, schedule, target)
	if err == nil {
		e.maybeActionAfterCompletion(ctx, schedule)
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
		e.routeToDLQ(ctx, schedule, target, "retry record serialisation failed: "+mErr.Error())
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}

	record := &schedulerstore.RetryRecord{
		ID:                    fmt.Sprintf("retry-%s-%d", uuidString(), now.UnixNano()),
		ScheduleName:          schedule.Name,
		GroupName:             schedule.GroupName,
		Region:                region,
		Target:                string(targetJSON),
		AttemptCount:          2, // Two immediate attempts already made.
		CreatedAt:             now,
		NextAttemptAt:         now.Add(computeRetryBackoff(3)),
		ActionAfterCompletion: string(schedule.ActionAfterCompletion),
		ScheduleExpression:    schedule.ScheduleExpression,
	}

	rs, rsErr := e.getRetryStore(region)
	if rsErr != nil {
		logs.Error("No retry store available for region, routing to DLQ",
			logs.String("schedule", schedule.Name),
			logs.String("region", region),
			logs.Err(rsErr))
		e.routeToDLQ(ctx, schedule, target, "retry store unavailable: "+rsErr.Error())
		e.maybeActionAfterCompletion(ctx, schedule)
		return
	}
	if sErr := rs.SaveRetryRecord(record); sErr != nil {
		logs.Error("Failed to persist retry record",
			logs.String("schedule", schedule.Name),
			logs.Err(sErr))
		// No record means no background retry: the delivery lifecycle ends
		// here and must reach the same terminal handling as retry
		// exhaustion instead of being silently dropped.
		e.routeToDLQ(ctx, schedule, target, "retry record persistence failed: "+sErr.Error())
		e.maybeActionAfterCompletion(ctx, schedule)
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
		rs, rsErr := e.getRetryStore(region)
		if rsErr != nil {
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
			e.wg.Add(1)
			go func() {
				defer e.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						logs.Error("scheduler: panic processing retry record",
							logs.String("recordId", record.ID),
							logs.Any("panic", r))
					}
				}()
				// Propagate e.ctx so retry processing stops when the engine
				// shuts down. Previously e.ctx was used inside
				// processRetryRecord but the passed context was ignored,
				// causing retry goroutines to outlive the engine.
				e.processRetryRecord(e.ctx, rs, record, now)
			}()
		}
	}
}

// processRetryRecord attempts redelivery of a single RetryRecord and handles
// the outcome: success -> delete, failure -> update next attempt or route to DLQ.
// The ctx parameter is propagated from checkRetries so that engine shutdown
// cancels in-flight retry processing.
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
		ScheduleExpression:    record.ScheduleExpression,
		Target:                &target,
		ActionAfterCompletion: schedulerstore.ActionAfterCompletion(record.ActionAfterCompletion),
	}

	// Verify the schedule still exists before retrying. With delayed
	// auto-deletion (AWS-compliant), the schedule remains alive during the
	// entire retry lifecycle, so this check passes for ActionAfterCompletion
	// = DELETE schedules and the retry continues. If the user manually
	// deletes the schedule mid-retry, this check fails and the retry record
	// is discarded — matching the user's intent to cancel. Only a confirmed
	// not-found discards: a storage failure retains the record for the next
	// sweep, because deleting it here would drop a pending delivery without
	// DLQ routing and break the at-least-once guarantee.
	schedStore := e.getStoreForSchedule(schedule)
	if schedStore != nil {
		if _, err := schedStore.GetSchedule(ctx, record.GroupName, record.ScheduleName); err != nil {
			if !errors.Is(err, schedulerstore.ErrScheduleNotFound) {
				logs.Error("Schedule existence probe failed, retaining retry record",
					logs.String("schedule", record.ScheduleName),
					logs.String("recordId", record.ID),
					logs.Err(err))
				return
			}
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
		e.maybeActionAfterCompletion(ctx, schedule)
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
		e.maybeActionAfterCompletion(ctx, schedule)
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
		e.maybeActionAfterCompletion(ctx, schedule)
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

	// The bus-less direct-delivery path cannot reach the SQS invoker
	// either; without this guard the invoker lookup below would
	// dereference a nil bus.
	if e.bus == nil {
		logs.Error("Scheduler engine has no event bus for DLQ delivery",
			logs.String("schedule", schedule.Name),
			logs.String("dlqArn", dlqArn))
		return
	}

	message := scheduleInput(target, schedule.Name)

	logs.Warn("Routing failed schedule delivery to DLQ",
		logs.String("schedule", schedule.Name),
		logs.String("dlqArn", dlqArn),
		logs.String("reason", reason))

	// DeadLetterConfig must be an SQS queue (AWS specification); one parse
	// supplies both the service check and the delivery region.
	_, dlqService, dlqRegion, _, _ := svcarn.SplitARN(dlqArn)
	if dlqService != "sqs" {
		logs.Error("DeadLetterConfig ARN must reference an SQS queue",
			logs.String("dlqArn", dlqArn),
			logs.String("service", dlqService))
		return
	}

	sqsInvoker := e.bus.SQSInvoker()
	if sqsInvoker == nil {
		logs.Error("SQS invoker not available for DLQ delivery",
			logs.String("dlqArn", dlqArn))
		return
	}
	queueName := svcarn.ExtractQueueNameFromARN(dlqArn)
	queueURL, qErr := sqsInvoker.GetQueueByName(ctx, dlqRegion, queueName)
	if qErr != nil {
		logs.Error("Failed to resolve DLQ queue URL",
			logs.String("dlqArn", dlqArn),
			logs.Err(qErr))
		return
	}
	sendOpts := sqsFifoSendOptions(target, queueName, schedule.Name)
	if _, _, err := sqsInvoker.SendMessage(ctx, dlqRegion, queueURL, message, sendOpts); err != nil {
		logs.Error("Failed to send to DLQ",
			logs.String("dlqArn", dlqArn),
			logs.Err(err))
	}
}

// uuidString returns a UUID-like unique string for retry record IDs.
// Uses crypto/rand for uniqueness without external UUID dependency.
func uuidString() string {
	b := make([]byte, 16)
	_, _ = crand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
