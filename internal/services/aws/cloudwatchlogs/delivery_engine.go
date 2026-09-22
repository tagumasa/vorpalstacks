package cloudwatchlogs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"sync/atomic"
	"time"

	"vorpalstacks/internal/common/worker"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The vended-logs delivery engine: one background pass per delivery
// cadence reads every configured region's deliveries, fetches the source
// group's events past the delivery's persisted cursor, and writes them
// to the delivery's destination. The cadence is the documented one —
// "One or more log files are created every five minutes in the specified
// bucket" (Logs sent to Amazon S3, CloudWatch Logs User Guide) — and the
// S3 prefix is the documented granted prefix AWSLogs/<account-id>/.

// deliveryTickerInterval is the delivery cadence. TEST_MODE shortens it
// so the SDK-level delivery pins observe a full pass in seconds.
var deliveryTickerInterval = worker.Cadence(5*time.Minute, time.Second)

// startDeliveryWorker launches the delivery-cadence worker under the
// service lifecycle (Stop cancels the context and waits).
func (s *LogsService) startDeliveryWorker() {
	s.startRespawningWorker("delivery worker",
		worker.TickerLoop(s.ctx.Done(), deliveryTickerInterval, s.tickDeliveries))
}

// tickDeliveries runs one delivery pass over every configured region.
// The pass guard is a service-instance field: a package-global guard
// serialised the passes of every service instance (every account), so a
// losing CAS skipped that account's whole pass for the cadence.
func (s *LogsService) tickDeliveries() {
	if !atomic.CompareAndSwapUint32(&s.delivering, 0, 1) {
		return
	}
	defer atomic.StoreUint32(&s.delivering, 0)

	for _, region := range worker.ActiveRegions(s.storageManager) {
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			logs.Warn("Delivery pass skipped: region store unavailable",
				logs.String("region", region), logs.Err(err))
			continue
		}
		deliveries, err := store.ListDeliveries()
		if err != nil {
			logs.Warn("Delivery pass skipped: delivery listing failed",
				logs.String("region", region), logs.Err(err))
			continue
		}
		for _, delivery := range deliveries {
			s.runOneDelivery(s.ctx, store, region, delivery)
		}
	}
}

// runOneDelivery executes one delivery: fetch the source group's events
// past the cursor, shape them per the delivery's configuration, and
// write them to the destination. The cursor advances only after a
// successful write, so a failed pass re-delivers at least once; a
// delivery whose source or destination disappeared is skipped with a
// warning (the API plane reports those states, the engine keeps them
// visible).
func (s *LogsService) runOneDelivery(ctx context.Context, store *logsstore.Store, region string, delivery *logsstore.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in delivery pass",
				logs.String("deliveryId", delivery.Id),
				logs.Any("panic", r))
		}
	}()

	source, err := store.GetDeliverySource(delivery.DeliverySourceName)
	if err != nil {
		logs.Warn("Delivery source missing for delivery pass",
			logs.String("deliveryId", delivery.Id),
			logs.String("deliverySourceName", delivery.DeliverySourceName))
		return
	}
	sourceGroup := resolveLogGroupIdentifier(source.ResourceArn)
	if _, err := store.GetLogGroup(sourceGroup); err != nil {
		// The source's underlying resource is gone: the source reports
		// INACTIVE/RESOURCE_DELETED and there is nothing to deliver. The
		// skip is the source's documented steady state, not a fault, so it
		// logs at debug level — one line per pass for the operator tracing
		// why no files appear, no warning spam.
		logs.Debug("Delivery pass skipped: source group no longer exists",
			logs.String("deliveryId", delivery.Id),
			logs.String("sourceGroup", sourceGroup))
		return
	}

	events, late, scanMark, err := fetchDeliveryEvents(store, sourceGroup, delivery)
	if err != nil {
		// No cursor advance on a read failure: the failed pass must
		// re-deliver at least once, never strand the unread window.
		logs.Error("Delivery pass aborted: source stream read failed",
			logs.String("deliveryId", delivery.Id), logs.Err(err))
		return
	}
	if len(events) == 0 && len(late) == 0 {
		advanceDeliveryIngestionMark(store, delivery.Id, scanMark)
		return
	}

	destName := svcarn.ExtractDeliveryDestinationNameFromARN(delivery.DeliveryDestinationArn)
	dest, err := store.GetDeliveryDestination(destName)
	if err != nil {
		logs.Warn("Delivery destination missing for delivery pass",
			logs.String("deliveryId", delivery.Id),
			logs.String("deliveryDestinationArn", delivery.DeliveryDestinationArn))
		return
	}

	// The late window leads the batch: its events sort strictly before
	// the normal window's first event, so the batch stays chronological
	// for the destination's per-stream batch cutting.
	batch := late
	batch = append(batch, events...)
	switch dest.DeliveryDestinationType {
	case "CWL":
		if !s.deliverEventsToGroup(region, dest, source, delivery, batch) {
			return
		}
	case "S3":
		if !s.deliverEventsToS3(ctx, region, dest, source, delivery, batch) {
			return
		}
	default:
		logs.Warn("Delivery destination type is not deliverable on this platform",
			logs.String("deliveryId", delivery.Id),
			logs.String("deliveryDestinationType", dest.DeliveryDestinationType))
		return
	}

	if ctx.Err() != nil {
		return
	}
	if len(events) == 0 {
		// A late-only pass delivered no cursor-window events: the cursor
		// triple stands and the ingestion mark advances alone.
		advanceDeliveryIngestionMark(store, delivery.Id, scanMark)
		return
	}
	last := events[len(events)-1]
	// The cursor merge runs inside the delivery mutate seam: a concurrent
	// UpdateDeliveryConfiguration's shaping fields survive this write,
	// and the cursor cannot revert a shaping update that committed while
	// the pass was delivering. The count records how many events share
	// the cursor triple — byte-identical events at the same millisecond
	// are distinct store records and each must deliver exactly once per
	// pass boundary.
	trailing := 0
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Timestamp != last.Timestamp || events[i].LogStreamName != last.LogStreamName ||
			events[i].Digest != last.Digest {
			break
		}
		trailing++
	}
	if err := store.MutateDelivery(delivery.Id, func(d *logsstore.Delivery) error {
		extendsCursor := last.Timestamp == d.CursorTime && last.LogStreamName == d.CursorStream &&
			last.Digest == d.CursorDigest
		d.CursorTime = last.Timestamp
		d.CursorStream = last.LogStreamName
		d.CursorDigest = last.Digest
		if extendsCursor {
			// The pass's last event sits at the previous cursor triple:
			// the fetch already skipped that triple's first CursorCount
			// events, so the counts add. Replacing the count would skip
			// too few on every later pass and re-deliver the triple's
			// earlier events forever.
			d.CursorCount = d.CursorCount + trailing
		} else {
			d.CursorCount = trailing
		}
		d.CursorIngestionMark = scanMark
		return nil
	}); err != nil {
		logs.Error("Failed to persist delivery cursor",
			logs.String("deliveryId", delivery.Id), logs.Err(err))
	}
}

// advanceDeliveryIngestionMark moves the delivery's late-window watermark
// on a pass that had nothing to deliver: the late scan already observed
// that no qualifying chunk exists, and leaving the mark behind would make
// every later pass re-scan the same empty window. The merge rides the
// delivery mutate seam so a concurrent shaping update survives it.
func advanceDeliveryIngestionMark(store *logsstore.Store, deliveryId string, scanMark int64) {
	if err := store.MutateDelivery(deliveryId, func(d *logsstore.Delivery) error {
		if scanMark > d.CursorIngestionMark {
			d.CursorIngestionMark = scanMark
		}
		return nil
	}); err != nil {
		logs.Error("Failed to persist delivery ingestion mark",
			logs.String("deliveryId", deliveryId), logs.Err(err))
	}
}

// deliveryEvent is one source event positioned in the delivery's sort
// order: (event timestamp, log stream, message digest).
type deliveryEvent struct {
	Timestamp     int64
	LogStreamName string
	Digest        string
	Message       string
}

// fetchDeliveryEvents reads the source group's events past the delivery
// cursor in the delivery sort order. The fetch re-reads from the
// cursor's timestamp inclusive and drops everything before the cursor
// triple plus the first CursorCount events at it, so byte-identical
// events sharing the cursor timestamp deliver once each and events
// arriving later still deliver. A stream that cannot be read fails the
// whole pass: the cursor advances from the delivered events and a later
// fetch never looks below it again, so skipping the stream here would
// make its window permanently undeliverable.
//
// The late window is the second half of that promise: events whose
// timestamps fall below the cursor but whose ingestion time exceeds the
// delivery's CursorIngestionMark were put after a pass had already
// passed their position, and no timestamp-window fetch can ever see
// them again. They return in the separate late slice, sorted by the
// same key, and every element sorts strictly before the normal window's
// first element. scanMark is the pass's scan-start clock — the caller
// persists it as the next CursorIngestionMark once the pass's deliveries
// have written, which the group lock's stamp-before-commit ordering
// keeps sound (an event invisible to this scan was ingested after
// scanMark and stays above the mark for the next one).
func fetchDeliveryEvents(store *logsstore.Store, sourceGroup string, delivery *logsstore.Delivery) (events, late []deliveryEvent, scanMark int64, err error) {
	scanMark = time.Now().UnixMilli()
	streams, err := fetchAllLogStreams(store, sourceGroup, "")
	if err != nil {
		return nil, nil, 0, err
	}
	for _, stream := range streams {
		raw, err := fetchAllLogEvents(store, sourceGroup, stream.Name, delivery.CursorTime, 0)
		if err != nil {
			return nil, nil, 0, err
		}
		for _, evt := range raw {
			events = append(events, deliveryEvent{
				Timestamp:     evt.Timestamp,
				LogStreamName: stream.Name,
				Digest:        eventDigest(evt.Message),
				Message:       evt.Message,
			})
		}
		lateEntries, err := store.LateIngestionEvents(sourceGroup, stream.Name, delivery.CursorIngestionMark, delivery.CursorTime)
		if err != nil {
			return nil, nil, 0, err
		}
		for _, evt := range lateEntries {
			late = append(late, deliveryEvent{
				Timestamp:     evt.Timestamp,
				LogStreamName: stream.Name,
				Digest:        eventDigest(evt.Message),
				Message:       evt.Message,
			})
		}
	}
	sortDeliveryEvents(events)
	sortDeliveryEvents(late)
	past := 0
	atCursor := 0
	for _, e := range events {
		if afterDeliveryCursor(e, delivery) {
			break
		}
		if deliveryCursorEqual(e, delivery) {
			atCursor++
			if atCursor <= delivery.CursorCount {
				past++
				continue
			}
			// An event at the cursor triple beyond the delivered count
			// is a later-arriving byte-identical event: it delivers.
			break
		}
		past++
	}
	// The slice is sorted; everything at or before the cursor is the
	// leading run, so one cut removes it.
	return events[past:], late, scanMark, nil
}

// sortDeliveryEvents orders a delivery window by the engine's global
// key: event timestamp, then log stream, then message digest.
func sortDeliveryEvents(events []deliveryEvent) {
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Timestamp != events[j].Timestamp {
			return events[i].Timestamp < events[j].Timestamp
		}
		if events[i].LogStreamName != events[j].LogStreamName {
			return events[i].LogStreamName < events[j].LogStreamName
		}
		return events[i].Digest < events[j].Digest
	})
}

// afterDeliveryCursor reports whether the event sorts strictly after the
// delivery's cursor triple.
func afterDeliveryCursor(e deliveryEvent, delivery *logsstore.Delivery) bool {
	if delivery.CursorTime == 0 {
		return true
	}
	if e.Timestamp != delivery.CursorTime {
		return e.Timestamp > delivery.CursorTime
	}
	if e.LogStreamName != delivery.CursorStream {
		return e.LogStreamName > delivery.CursorStream
	}
	return e.Digest > delivery.CursorDigest
}

// deliveryCursorEqual reports whether the event sits exactly at the
// cursor triple.
func deliveryCursorEqual(e deliveryEvent, delivery *logsstore.Delivery) bool {
	return delivery.CursorTime != 0 &&
		e.Timestamp == delivery.CursorTime &&
		e.LogStreamName == delivery.CursorStream &&
		e.Digest == delivery.CursorDigest
}

// eventDigest is the message tiebreak component of the delivery order.
func eventDigest(message string) string {
	sum := sha256.Sum256([]byte(message))
	return hex.EncodeToString(sum[:8])
}

// deliveryResourceRegion resolves the region a delivery resource ARN
// addresses, deferring to the fallback when the ARN carries none. The
// admission checks and the engine's destination writes share this one
// rule, so they always address the same region's store.
func deliveryResourceRegion(arn, fallbackRegion string) string {
	_, _, region, _, _ := svcarn.SplitARN(arn)
	if region == "" {
		return fallbackRegion
	}
	return region
}
