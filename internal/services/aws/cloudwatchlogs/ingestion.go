package cloudwatchlogs

import (
	"fmt"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// Every log-event write path funnels through putLogEventsCore: the HTTP
// PutLogEvents handler, the bus handlers (Lambda execution logs, direct
// puts from EventBridge/Scheduler/SFN targets, API Gateway access logs)
// and the cross-service logs invoker. Validation (batch size,
// chronological order, age windows), the store write, and the metric and
// subscription filter fan-out happen in that one seam, so no ingestion
// surface can skip any of the three. The S3 import worker is the one
// deliberate exception: AWS's import is not a PutLogEvents write (the
// documented import-role permission list carries CreateLogGroup and
// CreateLogStream but not PutLogEvents, and the destination is a
// managed log group), so imported backfill writes through the store
// directly and fires no filters — the adjudication lives in the plan
// register.

const (
	// filterFanOutWorkers bounds the goroutines that evaluate metric and
	// subscription filters after a write. One unbounded goroutine per
	// PutLogEvents call meant a burst of large batches spawned as many
	// goroutines as batches; the pool keeps that count flat regardless of
	// write rate.
	filterFanOutWorkers = 4

	// filterFanOutQueueDepth is the number of written batches that may
	// await filter evaluation. A queue full past the submit window
	// below drops that batch's filter evaluation with a warning instead
	// of blocking the write path (whose callers include Lambda
	// invocation threads the pool itself may have provoked).
	filterFanOutQueueDepth = 256

	// filterFanOutSubmitTimeout bounds how long a submission waits for a
	// queue slot: a full queue is a burst exceeding the pool's drain
	// rate, and the batch's metric and subscription deliveries are
	// observable outputs — the wait drains ordinary bursts and only a
	// queue still saturated after the window drops.
	filterFanOutSubmitTimeout = 2 * time.Second
)

// filterFanOutJob carries one written batch to the bounded workers. The
// events slice is private to the job (dispatchFilterFanOut copies), so the
// submitting write path cannot mutate it under evaluation.
type filterFanOutJob struct {
	store       *logsstore.Store
	region      string
	logGroup    string
	logStream   string
	events      []logsstore.LogEntry
	transformed map[string]string
}

// startFilterFanOutWorkers launches the bounded filter-evaluation pool.
// The workers share the service lifecycle: they exit on context
// cancellation and are waited on by Stop. The job channel is a pool-start
// constant — the workers capture it by value, so the field may be swapped
// (tests isolate the queue this way) without a racing worker re-parking
// on the replacement.
func (s *LogsService) startFilterFanOutWorkers() {
	jobs := make(chan filterFanOutJob, filterFanOutQueueDepth)
	s.filterJobs = jobs
	for i := 0; i < filterFanOutWorkers; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			s.runFilterFanOutWorker(jobs)
		}()
	}
}

// runFilterFanOutWorker is one pool worker's serve loop. Shutdown does
// not abandon the accepted batches: each queued job carries metric
// emissions and subscription deliveries, and a delivery dropped with no
// retry trail is lost silently, so a worker that sees the context
// cancelled drains what the queue still holds — bounded by the queue
// depth — before returning (a sender racing this drain drains its own
// landing through the shared drain helper).
func (s *LogsService) runFilterFanOutWorker(jobs <-chan filterFanOutJob) {
	for {
		select {
		case <-s.ctx.Done():
			s.drainFilterFanOutQueue(jobs)
			return
		case job := <-jobs:
			s.runFilterFanOutJob(job)
		}
	}
}

// dispatchFilterFanOut submits a written batch for metric and subscription
// filter evaluation. Submission waits a bounded window for a queue slot
// (a full queue is a burst the pool has not drained yet, and dropping a
// batch silently loses its metric emissions and subscription
// deliveries); only a queue still full after the window drops the
// batch, with a warning, rather than stalling ingestion indefinitely.
// A dropped batch is not silent: subscription delivery is a durable
// output, so the drop path evaluates the subscription filters inline —
// the deliveries either dispatch or ride the retry window like any
// other failed delivery — while the metric evaluation is dropped with
// the warning (metric data is recomputed per matched event and a burst
// past the pool's drain rate loses those data points visibly).
func (s *LogsService) dispatchFilterFanOut(store *logsstore.Store, region, logGroup, logStream string, events []logsstore.LogEntry, transformed map[string]string) {
	eventsCopy := make([]logsstore.LogEntry, len(events))
	copy(eventsCopy, events)
	job := filterFanOutJob{
		store:       store,
		region:      region,
		logGroup:    logGroup,
		logStream:   logStream,
		events:      eventsCopy,
		transformed: transformed,
	}
	deliverDropped := func(why string) {
		logs.Warn("Metric/subscription filter fan-out dropped one batch's evaluation; delivering subscriptions inline",
			logs.String("logGroup", logGroup),
			logs.String("logStream", logStream),
			logs.String("reason", why),
			logs.Int("queueDepth", filterFanOutQueueDepth))
		s.deliverSubscriptionEvents(store, region, logGroup, logStream, eventsCopy, transformed)
	}
	timer := time.NewTimer(filterFanOutSubmitTimeout)
	defer timer.Stop()
	select {
	case s.filterJobs <- job:
		// A send that lands while the service is shutting down may have
		// no consumer: the workers drain-until-empty on cancellation and
		// a send racing that drain can win the queue slot after every
		// worker has already returned. The sender that raced therefore
		// drains the queue itself — each job is consumed exactly once
		// whichever side takes it, and the shutdown contract (accepted
		// batches are evaluated, never silently lost) holds.
		if s.ctx.Err() != nil {
			s.drainFilterFanOutQueue(s.filterJobs)
		}
	case <-timer.C:
		deliverDropped("queue full past the submit window")
	case <-s.ctx.Done():
		deliverDropped("service shutting down")
	}
}

// drainFilterFanOutQueue consumes every job currently buffered on the
// fan-out queue, mirroring the workers' own shutdown drain. Both a
// worker and a racing sender may drain concurrently; the channel hands
// each job to exactly one consumer.
func (s *LogsService) drainFilterFanOutQueue(jobs <-chan filterFanOutJob) {
	for {
		select {
		case job := <-jobs:
			s.runFilterFanOutJob(job)
		default:
			return
		}
	}
}

// runFilterFanOutJob evaluates one batch against the group's metric and
// subscription filters. A panicking job is recovered and logged: it must
// not take the worker (and with it a quarter of the pool) down.
func (s *LogsService) runFilterFanOutJob(job filterFanOutJob) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error("Panic in metric/subscription filter fan-out",
				logs.String("logGroup", job.logGroup),
				logs.Any("panic", r))
		}
	}()
	s.evaluateMetricFilters(job.store, job.region, job.logGroup, job.logStream, job.events, job.transformed)
	s.deliverSubscriptionEvents(job.store, job.region, job.logGroup, job.logStream, job.events, job.transformed)
}

// IngestLogEvents is the cross-service entry to the shared ingestion seam:
// the logs invoker adapter calls it so platform writers (Lambda execution
// logs and any other invoker consumer) carry the same validation and
// filter fan-out as PutLogEvents API writes.
func (s *LogsService) IngestLogEvents(region, logGroupName, logStreamName string, entries []logsstore.LogEntry) error {
	_, err := s.putLogEventsCore(PutLogEventsInput{
		LogGroupName:  logGroupName,
		LogStreamName: logStreamName,
		Events:        wrapIngestedEvents(entries),
		Region:        region,
	})
	return err
}

// wrapIngestedEvents adapts already-typed internal entries to the parsed
// input form: the internal wire forms always carry the timestamp field, so
// every entry is timestamp-set.
func wrapIngestedEvents(entries []logsstore.LogEntry) []PutLogEvent {
	events := make([]PutLogEvent, len(entries))
	for i, e := range entries {
		events[i] = PutLogEvent{LogEntry: e, TimestampSet: true}
	}
	return events
}

// ingestBusEvents is the bus-plane leg of the ingestion seam. Platform
// internal writers (Lambda execution logs, EventBridge/Scheduler/SFN
// direct puts, API Gateway access logs) auto-create their log group and
// stream and log-and-swallow ingestion errors — the bus has no caller to
// surface an error to — but validation, the write and the filter fan-out
// are the shared seam's, identical to the API plane.
func (s *LogsService) ingestBusEvents(source, region, logGroup, logStream, accountID string, events []eventbus.LogEntry) {
	if s.ensureLogGroupAndStream(region, logGroup, logStream, accountID) == nil {
		return
	}
	if _, err := s.putLogEventsCore(PutLogEventsInput{
		LogGroupName:  logGroup,
		LogStreamName: logStream,
		Events:        wrapIngestedEvents(convertBusLogEntries(events)),
		Region:        region,
	}); err != nil {
		logs.Error(fmt.Sprintf("Failed to ingest %s", source),
			logs.String("logGroup", logGroup),
			logs.String("logStream", logStream),
			logs.Err(err))
	}
}

// busAccessLogEntry wraps one formatted access-log line as a bus log entry
// carrying the write time, matching the timestamp the platform's other
// single-message writers use.
func busAccessLogEntry(message string) eventbus.LogEntry {
	return eventbus.LogEntry{
		Timestamp: time.Now().UnixMilli(),
		Message:   message,
	}
}
