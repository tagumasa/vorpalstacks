package eventbridge

import (
	"strings"
	"testing"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestArchiveRetentionDaysRange pins the archive RetentionDays member to
// the model's @range(min:0) on both the create and the update path: zero
// (the documented default, indefinite retention) is accepted, negatives
// are rejected before any record is written.
func TestArchiveRetentionDaysRange(t *testing.T) {
	ctx := t.Context()
	svc := newEventBusCoreTestService()
	store := newEventBusCoreTestStore(t)
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	// The archive's event source is the event bus ARN.
	sourceBusARN := "arn:aws:events:us-east-1:000000000000:event-bus/default"

	_, err := svc.createArchiveCore(ctx, store, CreateArchiveInput{
		ArchiveName:    "retention-probe",
		EventSourceArn: sourceBusARN,
		ArchiveMergeMembers: ArchiveMergeMembers{
			RetentionDaysSet: true,
			RetentionDays:    -1,
		},
	})
	if err == nil || !strings.Contains(err.Error(), "RetentionDays") {
		t.Fatalf("create with RetentionDays -1: got %v, want the RetentionDays range rejection", err)
	}

	created, err := svc.createArchiveCore(ctx, store, CreateArchiveInput{
		ArchiveName:    "retention-probe",
		EventSourceArn: sourceBusARN,
		ArchiveMergeMembers: ArchiveMergeMembers{
			RetentionDaysSet: true,
			RetentionDays:    0,
		},
	})
	if err != nil {
		t.Fatalf("RetentionDays 0 (indefinite retention) must be accepted at create: %v", err)
	}
	if created.RetentionDays != 0 {
		t.Fatalf("created archive RetentionDays = %d, want 0", created.RetentionDays)
	}

	_, err = svc.updateArchiveCore(ctx, store, UpdateArchiveInput{
		ArchiveName: "retention-probe",
		ArchiveMergeMembers: ArchiveMergeMembers{
			RetentionDaysSet: true,
			RetentionDays:    -7,
		},
	})
	if err == nil || !strings.Contains(err.Error(), "RetentionDays") {
		t.Fatalf("update with RetentionDays -7: got %v, want the RetentionDays range rejection", err)
	}

	updated, err := svc.updateArchiveCore(ctx, store, UpdateArchiveInput{
		ArchiveName: "retention-probe",
		ArchiveMergeMembers: ArchiveMergeMembers{
			RetentionDaysSet: true,
			RetentionDays:    90,
		},
	})
	if err != nil {
		t.Fatalf("a positive RetentionDays must be accepted at update: %v", err)
	}
	if updated.RetentionDays != 90 {
		t.Fatalf("updated archive RetentionDays = %d, want 90", updated.RetentionDays)
	}
}

// TestArchiveEventSourceArnForm pins the EventBusArn shape contract of the
// archive family's EventSourceArn member: a value that is not an event-bus
// ARN is a pattern violation on both CreateArchive (required member) and
// ListArchives (optional filter), never a missing-bus lookup or a silent
// empty listing.
func TestArchiveEventSourceArnForm(t *testing.T) {
	ctx := t.Context()
	svc := newEventBusCoreTestService()
	store := newEventBusCoreTestStore(t)
	if err := store.CreateEventBus(ctx, &eventsstore.EventBus{Name: "default"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	_, err := svc.createArchiveCore(ctx, store, CreateArchiveInput{
		ArchiveName:    "form-probe",
		EventSourceArn: "default",
	})
	if err == nil || !strings.Contains(err.Error(), "EventSourceArn must match the event bus ARN pattern") {
		t.Fatalf("create with a bare name: got %v, want the EventSourceArn pattern rejection", err)
	}
	_, err = svc.createArchiveCore(ctx, store, CreateArchiveInput{
		ArchiveName:    "form-probe",
		EventSourceArn: "arn:aws:sqs:us-east-1:000000000000:queue/q",
	})
	if err == nil || !strings.Contains(err.Error(), "EventSourceArn must match the event bus ARN pattern") {
		t.Fatalf("create with a foreign-service ARN: got %v, want the EventSourceArn pattern rejection", err)
	}

	_, err = svc.listArchivesCore(ctx, store, ListArchivesInput{EventSourceArn: "default"})
	if err == nil || !strings.Contains(err.Error(), "EventSourceArn must match the event bus ARN pattern") {
		t.Fatalf("list with a bare name: got %v, want the EventSourceArn pattern rejection", err)
	}
	if _, err := svc.listArchivesCore(ctx, store, ListArchivesInput{
		EventSourceArn: "arn:aws:events:us-east-1:000000000000:event-bus/default",
	}); err != nil {
		t.Fatalf("list with a well-formed bus ARN must pass the filter: %v", err)
	}
}
