package scheduler

import (
	"context"
	"testing"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/eventbus"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// recordingKinesisInvoker captures the PutRecord contract: which region and
// stream each delivery addressed.
type recordingKinesisInvoker struct {
	putRegions map[string]string // streamName → region of the last PutRecord
	calls      int
}

func (r *recordingKinesisInvoker) ListShards(context.Context, string, string) ([]invokers.ShardInfo, error) {
	return nil, nil
}

func (r *recordingKinesisInvoker) PutRecord(_ context.Context, region, streamName, _ string, _ []byte) (string, error) {
	r.calls++
	if r.putRegions == nil {
		r.putRegions = make(map[string]string)
	}
	r.putRegions[streamName] = region
	return "seq-1", nil
}

func (r *recordingKinesisInvoker) CreateShardIterator(context.Context, string, string, string, string, string, *time.Time) (string, error) {
	return "", nil
}

func (r *recordingKinesisInvoker) GetRecords(context.Context, string, string, string, string, int32, bool) ([]invokers.KinesisRecord, string, error) {
	return nil, "", nil
}

func (r *recordingKinesisInvoker) StreamExists(context.Context, string, string) (bool, error) {
	return false, nil
}

// TestKinesisDeliveryUsesTargetARNRegion pins the region contract of the
// Kinesis deliverer: the target ARN's region reaches the PutRecord call —
// the platform is multi-region, and a cross-region target must not be
// captured by the server's default-region store — and the delivery leaves
// the schedule record unmutated (the region rode a schedule-record
// side-effect before; it is a delivery-routing concern, not schedule state).
func TestKinesisDeliveryUsesTargetARNRegion(t *testing.T) {
	svc := newLifecycleService(t)

	bus := eventbus.NewEventBus()
	if err := bus.Start(context.Background()); err != nil {
		t.Fatalf("start bus: %v", err)
	}
	t.Cleanup(func() { _ = bus.Shutdown(context.Background()) })
	invoker := &recordingKinesisInvoker{}
	bus.SetKinesisInvoker(invoker)
	svc.engine.SetEventBus(bus)

	schedule := &schedulerstore.Schedule{
		Name:      "kinesis-region-pin",
		GroupName: "default",
		Target: &schedulerstore.Target{
			// A region other than the engine's default: the delivery must
			// address the ARN's region, not the server default.
			Arn: "arn:aws:kinesis:eu-west-1:000000000000:stream/region-pin-stream",
		},
	}
	if err := svc.engine.sendToKinesis(t.Context(), schedule, schedule.Target); err != nil {
		t.Fatalf("sendToKinesis: %v", err)
	}

	if invoker.calls != 1 {
		t.Fatalf("PutRecord calls = %d, want 1", invoker.calls)
	}
	if got := invoker.putRegions["region-pin-stream"]; got != "eu-west-1" {
		t.Errorf("PutRecord region = %q, want the target ARN's region eu-west-1", got)
	}
	if schedule.Region != "" {
		t.Errorf("sendToKinesis mutated schedule.Region to %q; the region rides the PutRecord call and the record must stay unmutated", schedule.Region)
	}
}
