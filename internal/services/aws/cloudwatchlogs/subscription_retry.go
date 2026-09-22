package cloudwatchlogs

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"vorpalstacks/internal/common/worker"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The subscription-delivery retry window and its cadence. The window is
// documented: "Throttled deliverables are retried for up to 24 hours.
// After 24 hours, the failed deliverables are dropped" (CloudWatch Logs
// User Guide, subscription filters). The tick and the backoff cap are
// internal scheduling beneath it: the tick bounds how long a recovered
// destination waits for its backlog, and the backoff cap bounds the
// record rewrite rate against a persistently failing destination.
const (
	subscriptionDeliveryRetryWindow = 24 * time.Hour
	subscriptionDeliveryRetryTick   = 30 * time.Second
	subscriptionDeliveryMaxBackoff  = 5 * time.Minute
)

// subscriptionRetryBackoff spaces re-drives of a failing delivery: one
// second after the first failure, doubling per attempt up to the cap, so
// a flapping destination is not hammered while a briefly throttled one
// recovers within a few ticks.
func subscriptionRetryBackoff(attempts int) time.Duration {
	d := time.Second << attempts
	if d > subscriptionDeliveryMaxBackoff || d <= 0 {
		return subscriptionDeliveryMaxBackoff
	}
	return d
}

// pendingDeliveryID derives a batch's stable identity from its
// addressing and payload: a re-drive rewrites the same record instead of
// piling up a duplicate of the same batch.
func pendingDeliveryID(region, destArn, logGroup, logStream string, payload []byte) string {
	h := sha256.Sum256(append([]byte(region+"\x00"+destArn+"\x00"+logGroup+"\x00"+logStream+"\x00"), payload...))
	return hex.EncodeToString(h[:])
}

// retryFailedDelivery carries a failed subscription delivery into the
// retry window, preserving the first attempt and attempt count of a
// batch that already failed before.
func (s *LogsService) retryFailedDelivery(store *logsstore.Store, region, destArn, logGroup, logStream, distribution string, compressed []byte, cause error) {
	now := time.Now()
	id := pendingDeliveryID(region, destArn, logGroup, logStream, compressed)
	pd := &logsstore.PendingDelivery{
		ID:           id,
		Region:       region,
		DestArn:      destArn,
		LogGroup:     logGroup,
		LogStream:    logStream,
		Distribution: distribution,
		Payload:      compressed,
		FirstAttempt: now.UnixMilli(),
		Deadline:     now.Add(subscriptionDeliveryRetryWindow).UnixMilli(),
		NextAttempt:  now.Add(subscriptionRetryBackoff(0)).UnixMilli(),
	}
	if err := store.UpsertPendingDelivery(id, func(existing *logsstore.PendingDelivery) *logsstore.PendingDelivery {
		if existing != nil {
			pd.FirstAttempt = existing.FirstAttempt
			pd.Deadline = existing.Deadline
			pd.Attempts = existing.Attempts + 1
			pd.NextAttempt = now.Add(subscriptionRetryBackoff(pd.Attempts)).UnixMilli()
		}
		return pd
	}); err != nil {
		logs.Error("Failed to persist pending subscription delivery",
			logs.String("destinationArn", destArn),
			logs.Err(err))
		return
	}
	logs.Warn("Subscription delivery failed; batch rides the retry window",
		logs.String("destinationArn", destArn),
		logs.Int64("deadline", pd.Deadline),
		logs.Err(cause))
}

// drainDueSubscriptionDeliveries re-drives one region's pending
// deliveries: due batches redispatch, successes clear their records,
// deadline-expired batches drop with a visible warning ("After 24 hours,
// the failed deliverables are dropped"), and failures reschedule under
// the backoff.
func (s *LogsService) drainDueSubscriptionDeliveries(store *logsstore.Store) {
	pending, err := store.ListPendingDeliveries()
	if err != nil {
		logs.Error("Failed to list pending subscription deliveries", logs.Err(err))
		return
	}
	now := time.Now()
	for _, pd := range pending {
		if now.UnixMilli() > pd.Deadline {
			if err := store.DeletePendingDelivery(pd.ID); err != nil {
				logs.Error("Failed to drop expired pending subscription delivery",
					logs.String("destinationArn", pd.DestArn), logs.Err(err))
				continue
			}
			logs.Warn("Subscription delivery dropped after the retry window",
				logs.String("destinationArn", pd.DestArn),
				logs.String("logGroup", pd.LogGroup),
				logs.Int("attempts", pd.Attempts))
			continue
		}
		if now.UnixMilli() < pd.NextAttempt {
			continue
		}
		if err := s.dispatchSubscriptionDelivery(pd.Region, pd.DestArn, pd.LogGroup, pd.LogStream, pd.Distribution, pd.Payload); err != nil {
			// The reschedule merges into the persisted record: a record
			// that left the window while this re-drive ran (delivered or
			// dropped by another pass) stays gone — a failed re-drive
			// must not resurrect it.
			if err := store.UpsertPendingDelivery(pd.ID, func(existing *logsstore.PendingDelivery) *logsstore.PendingDelivery {
				if existing == nil {
					return nil
				}
				existing.Attempts++
				existing.NextAttempt = now.Add(subscriptionRetryBackoff(existing.Attempts)).UnixMilli()
				return existing
			}); err != nil {
				logs.Error("Failed to reschedule pending subscription delivery",
					logs.String("destinationArn", pd.DestArn), logs.Err(err))
			}
			continue
		}
		if err := store.DeletePendingDelivery(pd.ID); err != nil {
			logs.Error("Failed to clear delivered pending subscription delivery",
				logs.String("destinationArn", pd.DestArn), logs.Err(err))
		}
	}
}

// startSubscriptionDeliveryRetryLoop runs the retry drain across the
// configured regions under the service's panic-respawn worker policy.
func (s *LogsService) startSubscriptionDeliveryRetryLoop() {
	s.startRespawningWorker("subscription delivery retry",
		worker.TickerLoop(s.ctx.Done(), subscriptionDeliveryRetryTick, func() {
			for _, region := range worker.ActiveRegions(s.storageManager) {
				store, err := s.getLogsStoreByRegion(region)
				if err != nil {
					continue
				}
				s.drainDueSubscriptionDeliveries(store)
			}
		}))
}

// emitControlMessage delivers the documented reachability record to a
// subscription destination: "Sometimes CloudWatch Logs may emit Amazon
// Kinesis Data Streams records with a 'CONTROL_MESSAGE' type, mainly for
// checking if the destination is reachable" (the Lambda section carries
// the same sentence for Lambda invocations). The probe rides the same
// dispatch and payload shape as a data batch, with the CONTROL_MESSAGE
// type and no log events; its own delivery failure is the answer the
// probe sought, so it warns and does not enter the retry window.
func (s *LogsService) emitControlMessage(region, destArn, logGroup, filterName, distribution string) {
	payload := map[string]interface{}{
		"owner":               s.accountID,
		"logGroup":            logGroup,
		"logStream":           "",
		"subscriptionFilters": []string{filterName},
		"messageType":         "CONTROL_MESSAGE",
		"logEvents":           []map[string]interface{}{},
	}
	compressed, err := compressJSON(payload)
	if err != nil {
		return
	}
	if err := s.dispatchSubscriptionDelivery(region, destArn, logGroup, "", distribution, compressed); err != nil {
		logs.Warn("Subscription destination reachability probe could not be delivered",
			logs.String("destinationArn", destArn), logs.Err(err))
	}
}
