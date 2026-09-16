package eventbridge

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// The accept/deliver matrix at PutTargets: a target ARN is accepted only
// when the platform has a delivery path for its service AND resource form;
// every rejected form names its reason.
func TestTargetAcceptanceMatrix(t *testing.T) {
	cases := []struct {
		arn      string
		accepted bool
	}{
		{"arn:aws:lambda:us-east-1:000000000000:function:fn", true},
		{"arn:aws:sqs:us-east-1:000000000000:queue", true},
		{"arn:aws:sns:us-east-1:000000000000:topic/t", true},
		{"arn:aws:states:us-east-1:000000000000:stateMachine:sm", true},
		{"arn:aws:logs:us-east-1:000000000000:log-group:g", true},
		{"arn:aws:kinesis:us-east-1:000000000000:stream/s", true},
		{"arn:aws:firehose:us-east-1:000000000000:deliverystream/ds", true},
		{"arn:aws:events:us-east-1:000000000000:event-bus/default", true},
		{"arn:aws:events:us-east-1:000000000000:event-bus/custom-bus", true},
		// API destinations are invoked over HTTPS through their connection.
		{"arn:aws:events:us-east-1:000000000000:api-destination/dest/id-1", true},
		// Every other events-service resource form is the silent-drop class:
		// no rule listing would ever match the phantom bus it names.
		{"arn:aws:events:us-east-1:000000000000:rule/x", false},
		{"arn:aws:events:us-east-1:000000000000:archive/a", false},
		{"arn:aws:events:us-east-1:000000000000:bare-name", false},
		// AppSync: the GraphQL endpoint forms name the API to invoke.
		{"arn:aws:appsync:us-east-1:000000000000:apis/abc123", true},
		{"arn:aws:appsync:us-east-1:000000000000:apis/abc123/endpoints/GRAPHQL", true},
		{"arn:aws:appsync:us-east-1:000000000000:functions/f", false},
		// SSM Run Command and ECS have no delivery path on this platform.
		{"arn:aws:ssm:us-east-1:000000000000:document/AWS-RunShellScript", false},
		{"arn:aws:ecs:us-east-1:000000000000:cluster/c", false},
		{"arn:aws:sagemaker:us-east-1:000000000000:pipeline/p", false},
		{"", false},
	}
	for _, tc := range cases {
		err := validateTargetARNAcceptance(tc.arn)
		if tc.accepted && err != nil {
			t.Errorf("target ARN %q must be accepted, got rejection: %v", tc.arn, err)
		}
		if !tc.accepted && err == nil {
			t.Errorf("target ARN %q must be rejected", tc.arn)
		}
	}
}

// Rejections carry a reason that names the unsupported surface, so the
// per-entry failure explains itself.
func TestTargetAcceptanceRejectionsNameReasons(t *testing.T) {
	for arn, want := range map[string]string{
		"arn:aws:events:us-east-1:000000000000:rule/x":                   "event-bus/ and api-destination/",
		"arn:aws:ssm:us-east-1:000000000000:document/AWS-RunShellScript": "Parameter Store only",
		"arn:aws:ecs:us-east-1:000000000000:cluster/c":                   "ECS",
	} {
		err := validateTargetARNAcceptance(arn)
		if err == nil {
			t.Fatalf("target ARN %q must be rejected", arn)
		}
		if !strings.Contains(err.Error(), want) {
			t.Errorf("rejection of %q must mention %q, got: %v", arn, want, err)
		}
	}
}

// Delivery-side permanent classification: target forms with no delivery
// path terminate on the first attempt instead of consuming the retry
// budget.
func TestDispatchPermanentClassification(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()

	permanent := []struct {
		name   string
		target eventsstore.Target
	}{
		{"firehose before the service exists", eventsstore.Target{ARN: "arn:aws:firehose:us-east-1:000000000000:deliverystream/ds"}},
		{"ecs not implemented", eventsstore.Target{ARN: "arn:aws:ecs:us-east-1:000000000000:cluster/c"}},
		{"ssm run command not implemented", eventsstore.Target{ARN: "arn:aws:ssm:us-east-1:000000000000:document/AWS-RunShellScript"}},
		{"events resource form without a delivery path", eventsstore.Target{ARN: "arn:aws:events:us-east-1:000000000000:rule/x"}},
		{"api destination that does not exist", eventsstore.Target{ARN: "arn:aws:events:us-east-1:000000000000:api-destination/ghost/id"}},
		{"appsync target without an operation", eventsstore.Target{
			ARN:               "arn:aws:appsync:us-east-1:000000000000:apis/abc123",
			AppSyncParameters: &eventsstore.AppSyncParameters{},
		}},
		{"appsync ARN not naming an API endpoint", eventsstore.Target{
			ARN:               "arn:aws:appsync:us-east-1:000000000000:functions/f",
			AppSyncParameters: &eventsstore.AppSyncParameters{GraphQLOperation: "mutation { m }"},
		}},
	}
	for _, tc := range permanent {
		job := svc.newDeliveryJob("us-east-1", &eventsstore.Event{ID: "evt-matrix"}, tc.target, []byte(`{}`))
		err := svc.attemptDelivery(context.Background(), job)
		if err == nil {
			t.Fatalf("%s: expected a failure", tc.name)
		}
		if !errors.Is(err, errPermanentDelivery) {
			t.Errorf("%s: failure must be classified permanent, got: %v", tc.name, err)
		}
	}

	// An AppSync target with an operation but no configured invoker is a
	// deployment gap, not a target defect: the failure stays transient so
	// the retry budget keeps the delivery alive across wiring recovery.
	job := svc.newDeliveryJob("us-east-1", &eventsstore.Event{ID: "evt-matrix"},
		eventsstore.Target{
			ARN:               "arn:aws:appsync:us-east-1:000000000000:apis/abc123",
			AppSyncParameters: &eventsstore.AppSyncParameters{GraphQLOperation: "mutation { m }"},
		}, []byte(`{}`))
	err = svc.attemptDelivery(context.Background(), job)
	if err == nil || errors.Is(err, errPermanentDelivery) {
		t.Fatalf("a missing invoker must be transient, got: %v", err)
	}
}

// recordingAppSyncInvoker captures the mutation document and variables of
// every invocation.
type recordingAppSyncInvoker struct {
	mu         sync.Mutex
	apiIDs     []string
	operations []string
	variables  [][]byte
}

func (r *recordingAppSyncInvoker) ExecuteGraphQLMutation(_ context.Context, region, apiID, operation string, variablesJSON []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.apiIDs = append(r.apiIDs, apiID)
	r.operations = append(r.operations, operation)
	r.variables = append(r.variables, append([]byte(nil), variablesJSON...))
	return nil
}

func (r *recordingAppSyncInvoker) calls() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.operations)
}

var _ invokers.AppSyncInvoker = (*recordingAppSyncInvoker)(nil)

// The mutation document crosses the bus on the delivery event: the handler
// reconstructs AppSyncParameters without re-reading the stored target and
// invokes the API named by the ARN with the transformed payload as the
// operation's variables.
func TestHandleBusDeliveryInvokesAppSyncTarget(t *testing.T) {
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	svc := NewEventsService(mgr, "000000000000")
	defer svc.Close()
	invoker := &recordingAppSyncInvoker{}
	bus := eventbus.NewEventBus()
	bus.SetAppSyncInvoker(invoker)
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}

	operation := "mutation Create($id: ID!) { pushEvent(id: $id) { id } }"
	evt := &eventbus.EventBridgeDeliveryEvent{
		TargetARN:               "arn:aws:appsync:us-east-1:000000000000:apis/abc123",
		Input:                   []byte(`{"id":"evt-appsync"}`),
		EventBridgeEventID:      "evt-appsync",
		AppSyncGraphQLOperation: operation,
	}
	evt.Region = "us-east-1"

	if res := svc.handleBusDelivery(context.Background(), evt); res.Error != nil {
		t.Fatalf("the AppSync delivery must succeed: %v", res.Error)
	}
	if invoker.calls() != 1 {
		t.Fatalf("exactly one mutation must run, got %d", invoker.calls())
	}
	if invoker.apiIDs[0] != "abc123" {
		t.Fatalf("the API ID must come from the ARN, got %q", invoker.apiIDs[0])
	}
	if invoker.operations[0] != operation {
		t.Fatalf("the mutation document must cross the bus unchanged, got %q", invoker.operations[0])
	}
	if string(invoker.variables[0]) != `{"id":"evt-appsync"}` {
		t.Fatalf("the payload must become the variables, got %q", invoker.variables[0])
	}
}
