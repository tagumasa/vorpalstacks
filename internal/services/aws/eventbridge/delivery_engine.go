package eventbridge

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// errPermanentDelivery marks delivery failures that no retry can clear:
// the target type has no delivery path in this deployment. Retrying such a
// failure would burn the full retry budget on every event to that target,
// so the delivery fails straight through to its terminal handling.
var errPermanentDelivery = errors.New("permanent delivery failure")

// permanentDeliveryError builds a delivery error classified as permanent.
func permanentDeliveryError(format string, args ...interface{}) error {
	return fmt.Errorf("%w: "+format, append([]interface{}{errPermanentDelivery}, args...)...)
}

// deliveryJob is one target delivery together with its retry state. The
// first attempt runs on the caller's thread (the bus handler or the
// synchronous dispatch loop); the engine takes over only when that attempt
// fails and the retry budget is not yet spent.
type deliveryJob struct {
	region  string
	event   *eventsstore.Event
	target  eventsstore.Target
	payload []byte

	// ruleARN feeds the dead-letter attribute envelope (RULE_ARN); the
	// bus path carries the delivering rule's ARN, the synchronous arm the
	// rule in scope at fan-out.
	ruleARN string
	// traceHeader rides to the dead-letter write as the AWSTraceHeader
	// message attribute; it never appears in the delivered payload.
	traceHeader string
	// partitionKey is the Kinesis partition key resolved from the original
	// event by whoever holds the full event (the publisher on the bus
	// path, this constructor on the synchronous path); empty means the
	// event-ID default applies at the put.
	partitionKey string

	// attempts counts delivery attempts made so far, including the first.
	// The retry budget allows maxRetries retries after the initial
	// attempt, so the job is exhausted once attempts > maxRetries.
	attempts   int32
	maxRetries int32
	deadline   time.Time
	backoff    time.Duration
	maxBackoff time.Duration

	nextDue time.Time
}

// exhausted reports whether the retry budget is spent after a failed
// attempt.
func (j *deliveryJob) exhausted() bool {
	return j.attempts > j.maxRetries || time.Now().After(j.deadline)
}

// backoffSleep computes the delay before the next retry: exponential
// backoff capped at maxBackoff, plus jitter of up to half the backoff.
// EventBridge documents its retry schedule as "an exponential back off and
// jitter, or randomized delay" (user guide, "How EventBridge retries
// delivering events"). backoff must be positive — a zero value would make
// rand.Int63n panic on a zero operand.
func backoffSleep(backoff, maxBackoff time.Duration) time.Duration {
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return backoff + time.Duration(rand.Int63n(int64(backoff)/2))
}

// deliveryQueue is a min-heap of jobs ordered by nextDue.
type deliveryQueue []*deliveryJob

func (q deliveryQueue) Len() int            { return len(q) }
func (q deliveryQueue) Less(i, j int) bool  { return q[i].nextDue.Before(q[j].nextDue) }
func (q deliveryQueue) Swap(i, j int)       { q[i], q[j] = q[j], q[i] }
func (q *deliveryQueue) Push(x interface{}) { *q = append(*q, x.(*deliveryJob)) }
func (q *deliveryQueue) Pop() interface{} {
	old := *q
	n := len(old)
	job := old[n-1]
	old[n-1] = nil
	*q = old[:n-1]
	return job
}

// deliveryEngine owns the retry lifecycle for target deliveries whose
// first attempt failed. Workers never sleep between retries: each retry is
// scheduled onto the due-time heap and dispatched when its time arrives,
// so a cohort of unreachable targets cannot occupy the worker set the way
// the previous in-handler sleep loop occupied the bus's async workers.
//
// The engine is scheduling only; delivery semantics live in the injected
// attempt and terminal callbacks (the EventsService methods), which keeps
// the retry mechanics testable against counted stubs.
type deliveryEngine struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	mu       sync.Mutex
	queue    deliveryQueue
	wake     chan struct{}
	ready    chan *deliveryJob
	attempt  func(context.Context, *deliveryJob) error
	terminal func(context.Context, *deliveryJob, error) error
}

// eitherContext returns a context that ends when either input context
// ends. The watcher goroutine exits once the merged context is done, so
// each call costs one short-lived goroutine for the lifetime of the
// delivery.
func eitherContext(a, b context.Context) (context.Context, context.CancelFunc) {
	merged, cancel := context.WithCancel(a)
	go func() {
		select {
		case <-a.Done():
		case <-b.Done():
		case <-merged.Done():
		}
		cancel()
	}()
	return merged, cancel
}

// deliveryEngineWorkers is the engine's attempt concurrency. Attempts are
// bounded invoker calls, not sleeps, so the workers cannot be starved by
// the retry schedule itself.
const deliveryEngineWorkers = 8

func newDeliveryEngine(attempt func(context.Context, *deliveryJob) error, terminal func(context.Context, *deliveryJob, error) error) *deliveryEngine {
	ctx, cancel := context.WithCancel(context.Background())
	e := &deliveryEngine{
		ctx:      ctx,
		cancel:   cancel,
		wake:     make(chan struct{}, 1),
		ready:    make(chan *deliveryJob, deliveryEngineWorkers*2),
		attempt:  attempt,
		terminal: terminal,
	}
	e.wg.Add(1 + deliveryEngineWorkers)
	go e.schedulerLoop()
	for i := 0; i < deliveryEngineWorkers; i++ {
		go e.workerLoop()
	}
	return e
}

// Close cancels the engine context — aborting in-flight attempts and the
// scheduler wait — and waits for its goroutines. Pending retries die with
// the process; their events have already lived their at-least-once
// transport guarantee (the outbox entry completed at handoff). The
// scheduler's exit drain and schedule's post-push check log the stranded
// retries so the loss stays observable.
func (e *deliveryEngine) Close() {
	e.cancel()
	e.wg.Wait()
}

// reportStranded logs the retries still queued at scheduler exit. The jobs
// die with the process by design (Close's contract); the report makes the
// loss observable.
func (e *deliveryEngine) reportStranded() {
	e.mu.Lock()
	n := e.queue.Len()
	e.mu.Unlock()
	if n > 0 {
		logs.Warn("delivery engine shut down with pending retries; they die with the process",
			logs.Int("pendingRetries", n))
	}
}

// schedule places a job for its next retry and wakes the scheduler. It
// fails only when the engine is shut down; the caller then reports the
// failure so the transport-level retry can re-drive the delivery.
func (e *deliveryEngine) schedule(job *deliveryJob) error {
	select {
	case <-e.ctx.Done():
		return e.ctx.Err()
	default:
	}
	e.mu.Lock()
	heap.Push(&e.queue, job)
	e.mu.Unlock()
	select {
	case e.wake <- struct{}{}:
	default:
	}
	// A cancellation racing the check above strands the job on a queue no
	// scheduler will drain; the error keeps that loss on the reporting
	// path the contract promises (the scheduler's exit drain covers the
	// jobs queued before it left).
	select {
	case <-e.ctx.Done():
		return e.ctx.Err()
	default:
	}
	return nil
}

// pending returns the number of jobs waiting for a retry. Test inspection.
func (e *deliveryEngine) pending() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.queue.Len()
}

func (e *deliveryEngine) schedulerLoop() {
	defer e.wg.Done()
	defer e.reportStranded()
	for {
		var due []*deliveryJob
		var next time.Time
		e.mu.Lock()
		now := time.Now()
		for e.queue.Len() > 0 && !e.queue[0].nextDue.After(now) {
			due = append(due, e.queue[0])
			heap.Pop(&e.queue)
		}
		if e.queue.Len() > 0 {
			next = e.queue[0].nextDue
		}
		e.mu.Unlock()

		// Hand off outside the lock: the send blocks while every worker is
		// busy, and workers re-take the lock when they schedule the next
		// retry of a job they just ran.
		for _, job := range due {
			select {
			case e.ready <- job:
			case <-e.ctx.Done():
				return
			}
		}

		if next.IsZero() {
			select {
			case <-e.wake:
			case <-e.ctx.Done():
				return
			}
			continue
		}
		timer := time.NewTimer(time.Until(next))
		select {
		case <-timer.C:
		case <-e.wake:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
		case <-e.ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (e *deliveryEngine) workerLoop() {
	defer e.wg.Done()
	for {
		select {
		case <-e.ctx.Done():
			return
		case job := <-e.ready:
			e.runJob(job)
		}
	}
}

// runJob attempts one scheduled retry and either finishes the delivery,
// terminates it, or re-schedules it with the next backoff. An error
// carrying a Retry-After hint (API destination endpoints) lengthens the
// delay to whichever of the policy backoff and the requested delay is more
// conservative.
func (e *deliveryEngine) runJob(job *deliveryJob) {
	err := e.attempt(e.ctx, job)
	if err == nil {
		return
	}
	job.attempts++
	if errors.Is(err, errPermanentDelivery) || job.exhausted() {
		// The engine cannot hand the dead-letter failure back to a caller —
		// terminalDelivery already logs the loss at full strength.
		_ = e.terminal(e.ctx, job, err)
		return
	}
	delay := backoffSleep(job.backoff, job.maxBackoff)
	var hint errRetryAfterHint
	if errors.As(err, &hint) {
		if requested := hint.RetryAfterDelay(); requested > delay {
			delay = requested
		}
	}
	job.nextDue = time.Now().Add(delay)
	job.backoff *= 2
	if err := e.schedule(job); err != nil {
		// The schedule contract makes the caller report the failure; the
		// engine's internal caller is the last reporter on this path.
		logs.Warn("delivery retry lost at engine shutdown",
			logs.String("targetArn", job.target.ARN),
			logs.String("eventId", job.event.ID),
			logs.Err(err))
	}
}
