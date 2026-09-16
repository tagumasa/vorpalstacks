package eventbridge

import (
	"testing"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cloudwatchevents"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// newAdminTestHandler returns an admin handler over a fresh per-test store
// cached for the header-less default region (us-east-1).
func newAdminTestHandler(t *testing.T) (*AdminHandler, *eventsstore.EventsStore) {
	t.Helper()
	svc := NewEventsService(nil, "000000000000")
	store := newEventBusCoreTestStore(t)
	svc.SetEventsStore("us-east-1", store)
	return NewAdminHandler(svc), store
}

// TestAdminCreateEventBusAppliesTags pins the admin-plane CreateEventBus
// Tags pass-through: the proto tags reach the same tag store the HTTP
// plane's CreateEventBus writes through, so a bus created from the console
// carries its tags on every later ListTags — on both planes.
func TestAdminCreateEventBusAppliesTags(t *testing.T) {
	ctx := t.Context()
	h, store := newAdminTestHandler(t)

	resp, err := h.CreateEventBus(ctx, connect.NewRequest(&pb.CreateEventBusRequest{
		Name: "admin-tags-bus",
		Tags: []*pb.Tag{{Key: "env", Value: "admin"}},
	}))
	if err != nil {
		t.Fatalf("admin CreateEventBus: %v", err)
	}
	busARN := resp.Msg.GetEventbusarn()
	if busARN == "" {
		t.Fatal("admin CreateEventBus returned no ARN")
	}

	tags, err := store.TagStore.ListAsSlice(busARN)
	if err != nil {
		t.Fatalf("list tags: %v", err)
	}
	if len(tags) != 1 || tags[0].Key != "env" || tags[0].Value != "admin" {
		t.Fatalf("bus tags = %+v, want the single env=admin tag applied at creation", tags)
	}
}

// TestAdminListOmitsExhaustedToken pins the NextToken emission contract
// shared with the HTTP plane: an exhausted page leaves the member unset
// instead of carrying an empty string.
func TestAdminListOmitsExhaustedToken(t *testing.T) {
	ctx := t.Context()
	h, _ := newAdminTestHandler(t)

	buses, err := h.ListEventBuses(ctx, connect.NewRequest(&pb.ListEventBusesRequest{}))
	if err != nil {
		t.Fatalf("admin ListEventBuses: %v", err)
	}
	if got := buses.Msg.GetNexttoken(); got != "" {
		t.Fatalf("exhausted ListEventBuses page carries NextToken %q, want the member unset", got)
	}

	rules, err := h.ListRules(ctx, connect.NewRequest(&pb.ListRulesRequest{}))
	if err != nil {
		t.Fatalf("admin ListRules: %v", err)
	}
	if got := rules.Msg.GetNexttoken(); got != "" {
		t.Fatalf("exhausted ListRules page carries NextToken %q, want the member unset", got)
	}
}
