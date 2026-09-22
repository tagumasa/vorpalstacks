package cloudwatchlogs

import (
	"time"

	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The CWL-typed destination family of the vended-logs delivery engine:
// the pass writes the delivered events into the destination log group
// through the shared ingestion seam, cutting each stream's window into
// batches the seam accepts and dropping events the destination would
// permanently reject as aged out.

// deliverEventsToGroup delivers a CWL-typed delivery through the shared
// ingestion seam: the destination is a log group, and the delivered
// events ingest into it with their originating stream names. A delivery
// whose destination is the source group itself is skipped: it would
// re-ingest every event into the group it read from, an endless loop
// with no reader. Each stream's window is cut into batches the
// ingestion seam accepts — a window past the batch contract rejects
// outright, so an uncut delivery could never succeed and its cursor
// could never advance — and events the destination would permanently
// reject as aged out are dropped with a warning: their age only grows,
// so no later pass could deliver them either. The boolean reports full
// success; the cursor advances only then.
func (s *LogsService) deliverEventsToGroup(region string, dest *logsstore.DeliveryDestination, source *logsstore.DeliverySource, delivery *logsstore.Delivery, events []deliveryEvent) bool {
	destRegion, destGroup := deliveryGroupDestination(region, dest)
	sourceGroup := resolveLogGroupIdentifier(source.ResourceArn)
	if destRegion == region && destGroup == sourceGroup {
		logs.Warn("Delivery destination is the delivery source group; skipping self-delivery",
			logs.String("deliveryId", delivery.Id),
			logs.String("logGroup", sourceGroup))
		return false
	}

	dropCutoff := s.deliveryDropCutoff(destRegion, destGroup)

	// The delivered form is shaped like any other destination's: the
	// destination's outputFormat and the delivery's recordFields render
	// each event ("Example access log sent to CloudWatch Logs" is a JSON
	// record of the selected fields) — the raw message re-ingests only
	// under the raw format.
	fields := deliveryRecordFields(delivery)

	byStream := make(map[string][]logsstore.LogEntry)
	var order []string
	for _, e := range events {
		line, err := deliveryEventLine(dest.OutputFormat, delivery.FieldDelimiter, fields, e)
		if err != nil {
			logs.Error("Failed to render delivery event",
				logs.String("deliveryId", delivery.Id), logs.Err(err))
			return false
		}
		if _, ok := byStream[e.LogStreamName]; !ok {
			order = append(order, e.LogStreamName)
		}
		byStream[e.LogStreamName] = append(byStream[e.LogStreamName], logsstore.LogEntry{
			Timestamp: e.Timestamp,
			Message:   line,
		})
	}
	for _, stream := range order {
		var kept []logsstore.LogEntry
		dropped := 0
		for _, entry := range byStream[stream] {
			if entry.Timestamp < dropCutoff {
				dropped++
				continue
			}
			kept = append(kept, entry)
		}
		if dropped > 0 {
			logs.Warn("Delivery drops events the destination horizon permanently rejects",
				logs.String("deliveryId", delivery.Id),
				logs.String("destinationGroup", destGroup),
				logs.String("logStream", stream),
				logs.Int("droppedEvents", dropped))
		}
		if len(kept) == 0 {
			continue
		}
		// The destination group is customer-managed — "Create a delivery
		// destination for an existing log group" (vended-logs setup), with
		// the delivery flow's granted actions on it exactly the stream
		// creation and event writes below. A deleted destination group
		// fails the pass instead of being silently recreated with default
		// settings: the customer's deletion stands until the customer
		// recreates the group or removes the delivery.
		if !s.prepareDeliveryDestinationStream(destRegion, destGroup, stream) {
			return false
		}
		for _, batch := range splitDeliveryBatches(kept) {
			if err := s.IngestLogEvents(destRegion, destGroup, stream, batch); err != nil {
				logs.Warn("Delivery to destination group failed",
					logs.String("deliveryId", delivery.Id),
					logs.String("destinationGroup", destGroup), logs.Err(err))
				return false
			}
		}
	}
	return true
}

// prepareDeliveryDestinationStream verifies the destination group still
// exists and creates the delivery's log stream in it — the two actions the
// delivery flow holds on the customer-managed destination group
// ("logs:CreateLogStream" and "logs:PutLogEvents", the granted delivery
// principal permissions, "Logs sent to CloudWatch Logs"). A deleted
// destination group fails the pass: the engine never recreates it with
// default settings the customer never chose.
func (s *LogsService) prepareDeliveryDestinationStream(region, logGroup, logStream string) bool {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		logs.Error("Failed to resolve logs store",
			logs.String("logGroup", logGroup),
			logs.String("region", region),
			logs.Err(err))
		return false
	}
	if _, err := store.GetLogGroup(logGroup); err != nil {
		logs.Warn("Delivery destination log group no longer exists; the pass fails until the group is recreated",
			logs.String("logGroup", logGroup),
			logs.String("region", region))
		return false
	}
	return createLogStreamIfAbsent(store, logGroup, logStream)
}

// deliveryDropCutoff is the timestamp horizon past which the destination
// ingestion permanently rejects events: the 14-day rule, tightened to
// the destination group's retention horizon when that is the stricter
// bound ("Events older than 14 days or preceding the log group's
// retention period are rejected while processing remaining valid
// events", PutLogEvents operation documentation).
func (s *LogsService) deliveryDropCutoff(region, group string) int64 {
	now := time.Now().UnixMilli()
	cutoff := now - tooOldThreshold
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return cutoff
	}
	lg, err := store.GetLogGroup(group)
	if err != nil {
		return cutoff
	}
	if retention := retentionCutoffMillis(now, lg.RetentionInDays); retention > cutoff {
		cutoff = retention
	}
	return cutoff
}

// splitDeliveryBatches cuts one stream's pending window into batches the
// shared ingestion seam accepts: at most MaxBatchLogEvents events, at most
// MaxPutLogEventsBatchBytes of batch bytes ("The maximum batch size is
// 1,048,576 bytes", PutLogEvents API reference), and a timestamp span
// within the documented 24 hours. The caller's entries are in
// chronological order, so every prefix cut keeps each batch ordered.
func splitDeliveryBatches(entries []logsstore.LogEntry) [][]logsstore.LogEntry {
	var batches [][]logsstore.LogEntry
	start := 0
	batchBytes := 0
	for i := range entries {
		eventBytes := len(entries[i].Message) + logsstore.PutLogEventOverheadBytes
		if i > start &&
			(entries[i].Timestamp-entries[start].Timestamp > maxEventsTimeSpan ||
				batchBytes+eventBytes > logsstore.MaxPutLogEventsBatchBytes ||
				i-start+1 > logsstore.MaxBatchLogEvents) {
			batches = append(batches, entries[start:i])
			start = i
			batchBytes = eventBytes
			continue
		}
		batchBytes += eventBytes
	}
	if start < len(entries) {
		batches = append(batches, entries[start:])
	}
	return batches
}

// deliveryGroupDestination resolves the CWL destination's region and
// group name from its destination resource ARN.
func deliveryGroupDestination(fallbackRegion string, dest *logsstore.DeliveryDestination) (string, string) {
	return deliveryResourceRegion(dest.DestinationResourceArn, fallbackRegion), resolveLogGroupIdentifier(dest.DestinationResourceArn)
}
