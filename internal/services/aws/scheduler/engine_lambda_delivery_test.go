package scheduler

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// failingLambdaInvoker stands in for the Lambda service on the far side of
// the bus: every trigger invocation reports a failed function execution
// over a succeeded transport — exactly the outcome the engine must fold
// into its retry lifecycle instead of counting as a delivered schedule.
type failingLambdaInvoker struct{ calls int }

func (f *failingLambdaInvoker) InvokeForGateway(ctx context.Context, functionName string, payload []byte) (int64, []byte, error) {
	return 200, nil, nil
}

func (f *failingLambdaInvoker) InvokeForTrigger(ctx context.Context, functionName string, payload []byte) (invokers.LambdaInvocation, error) {
	f.calls++
	return invokers.LambdaInvocation{StatusCode: 200, FunctionError: "Handled"}, nil
}

func (f *failingLambdaInvoker) GetFunctionARN(ctx context.Context, functionName string) (string, error) {
	return "arn:aws:lambda:us-east-1:000000000000:function:" + functionName, nil
}

// TestLambdaFunctionErrorPersistsRetryRecord pins the FunctionError
// contract end to end: a schedule whose Lambda function fails on every
// attempt must NOT count as delivered — the immediate retry runs, and
// once it also fails a RetryRecord is persisted in the region's retry
// store so the background retry/DLQ lifecycle owns the delivery.
func TestLambdaFunctionErrorPersistsRetryRecord(t *testing.T) {
	svc := newLifecycleService(t)

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	invoker := &failingLambdaInvoker{}
	bus.SetLambdaInvoker(invoker)
	svc.engine.SetEventBus(bus)

	store, err := svc.GetStoreForRegion(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get store: %v", err)
	}
	if err := store.CreateSchedule(t.Context(), &schedulerstore.Schedule{
		Name:                  "lambda-fn-error",
		GroupName:             "default",
		State:                 schedulerstore.ScheduleStateEnabled,
		ScheduleExpression:    "at(2030-01-01T00:00:00)",
		ActionAfterCompletion: "NONE",
		Target: &schedulerstore.Target{
			Arn: "arn:aws:lambda:us-east-1:000000000000:function:failing",
		},
	}); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	schedule, err := store.GetSchedule(t.Context(), "default", "lambda-fn-error")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}

	svc.engine.deliverWithRetry(t.Context(), schedule, schedule.Target)

	// Both immediate attempts must have run and both must have failed —
	// the transport succeeded, the function did not.
	if invoker.calls != 2 {
		t.Errorf("trigger invocations = %d, want 2 (initial + immediate retry)", invoker.calls)
	}

	rs, err := svc.engine.getRetryStore(defaults.DefaultRegion)
	if err != nil {
		t.Fatalf("get retry store: %v", err)
	}
	// The record's NextAttemptAt lies in the future (backoff); a cutoff
	// an hour out observes it without waiting.
	due, err := rs.GetDueRetryRecords(time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("list retry records: %v", err)
	}
	found := false
	for _, rec := range due {
		if rec.ScheduleName == "lambda-fn-error" {
			found = true
		}
	}
	if !found {
		t.Errorf("no RetryRecord persisted for a FunctionError delivery (records: %d)", len(due))
	}
}

// TestLambdaDeliveryNilInvokerGuard pins the nil-invoker guard: when the
// bus is wired but the Lambda invoker is not (a deployment without the
// Lambda service), the delivery returns an error that feeds the
// retry/DLQ lifecycle — it must not panic the firing goroutine.
func TestLambdaDeliveryNilInvokerGuard(t *testing.T) {
	svc := newLifecycleService(t)

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	// No LambdaInvoker registered.
	svc.engine.SetEventBus(bus)

	err := svc.engine.deliverToTarget(t.Context(),
		&schedulerstore.Schedule{
			Name:      "nil-invoker",
			GroupName: "default",
			Region:    defaults.DefaultRegion,
		},
		&schedulerstore.Target{
			Arn: "arn:aws:lambda:us-east-1:000000000000:function:unwired",
		})
	if err == nil {
		t.Fatal("delivery through an unwired Lambda invoker returned nil, want an error")
	}
}
