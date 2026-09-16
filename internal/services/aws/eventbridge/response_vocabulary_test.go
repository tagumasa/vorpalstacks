package eventbridge

import (
	"context"
	"sort"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// The response-key vocabulary pin for every registered operation: each
// operation's emitted top-level member keys, with exact casing, against the
// output shape member names of the vendored eventbridge model
// (2015-10-07, transcribed into the table below — test files never read
// the vendored model). The rule: an emitted key that the output shape does
// not model is a wire-contract violation no matter how harmless it looks
// (SDKs silently drop it, and the drift is unpinnable), and a modelled
// member the scenario guarantees must actually be present. List members
// carry their item shapes the same way.
//
// The rows run as one scenario in table order against the real handler
// path (wire parse → Core → serialisation) over a real per-region store,
// so a key added to any serialiser without a model member — or dropped
// from one — fails here before it can drift further.

type vocabRow struct {
	name string
	call func(svc *EventsService, reqCtx *request.RequestContext) (interface{}, error)

	// modelled holds every member of the operation's output shape; every
	// emitted key must be one of these.
	modelled []string
	// always holds the modelled members this scenario guarantees; each
	// must be present in the emitted map.
	always []string
	// items maps a list member to its item shape's modelled member names;
	// every emitted item key must be one of these, and every item must
	// carry the itemsAlways members for that list.
	items       map[string][]string
	itemsAlways map[string][]string
}

func TestResponseKeyVocabulary(t *testing.T) {
	ctx := context.Background()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(ctx, mgr, "000000000000", "us-east-1")

	svc := NewEventsService(nil, "000000000000")
	bus := eventbus.NewEventBus()
	bus.SetSecretsManagerInvoker(newFakeSecretsInvoker())
	if err := svc.SetEventBus(bus); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(svc.Close)

	const busARN = "arn:aws:events:us-east-1:000000000000:event-bus/vocab-bus"
	const sqsTargetARN = "arn:aws:sqs:us-east-1:000000000000:vocab-target"

	// Captured wire values later rows depend on; the rows that use them
	// read the variables at call time, after the capturing row has run.
	var connectionArn, archiveArn string

	params := func(kv ...interface{}) map[string]interface{} {
		m := make(map[string]interface{}, len(kv)/2)
		for i := 0; i+1 < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		return m
	}
	call := func(op string, kv ...interface{}) func(*EventsService, *request.RequestContext) (interface{}, error) {
		return func(svc *EventsService, reqCtx *request.RequestContext) (interface{}, error) {
			return dispatchVocab(context.Background(), svc, op, params(kv...), reqCtx)
		}
	}
	deferred := func(op string, kv func() []interface{}) func(*EventsService, *request.RequestContext) (interface{}, error) {
		return func(svc *EventsService, reqCtx *request.RequestContext) (interface{}, error) {
			return dispatchVocab(context.Background(), svc, op, params(kv()...), reqCtx)
		}
	}

	rows := []vocabRow{
		{name: "CreateEventBus", call: call("CreateEventBus", "Name", "vocab-bus", "Description", "vocab bus"),
			modelled: []string{"EventBusArn", "Description", "KmsKeyIdentifier", "DeadLetterConfig", "LogConfig"},
			always:   []string{"EventBusArn", "Description"}},
		{name: "DescribeEventBus", call: call("DescribeEventBus", "Name", "vocab-bus"),
			modelled: []string{"Name", "Arn", "Description", "KmsKeyIdentifier", "DeadLetterConfig", "Policy", "LogConfig", "CreationTime", "LastModifiedTime"},
			always:   []string{"Name", "Arn", "CreationTime", "LastModifiedTime"}},
		{name: "UpdateEventBus", call: call("UpdateEventBus", "Name", "vocab-bus", "Description", "vocab bus v2"),
			modelled: []string{"Arn", "Name", "KmsKeyIdentifier", "Description", "DeadLetterConfig", "LogConfig"},
			always:   []string{"Arn", "Name"}},
		{name: "ListEventBuses", call: call("ListEventBuses"),
			modelled:    []string{"EventBuses", "NextToken"},
			always:      []string{"EventBuses"},
			items:       map[string][]string{"EventBuses": {"Name", "Arn", "Description", "Policy", "CreationTime", "LastModifiedTime"}},
			itemsAlways: map[string][]string{"EventBuses": {"Name", "Arn", "CreationTime", "LastModifiedTime"}}},
		{name: "PutRule", call: call("PutRule", "EventBusName", "vocab-bus", "Name", "vocab-rule", "EventPattern", `{"source":["vocab.test"]}`),
			modelled: []string{"RuleArn"},
			always:   []string{"RuleArn"}},
		{name: "DescribeRule", call: call("DescribeRule", "EventBusName", "vocab-bus", "Name", "vocab-rule"),
			modelled: []string{"Name", "Arn", "EventPattern", "ScheduleExpression", "State", "Description", "RoleArn", "ManagedBy", "EventBusName", "CreatedBy"},
			always:   []string{"Name", "Arn", "EventBusName", "State", "EventPattern"}},
		{name: "ListRules", call: call("ListRules", "EventBusName", "vocab-bus"),
			modelled:    []string{"Rules", "NextToken"},
			always:      []string{"Rules"},
			items:       map[string][]string{"Rules": {"Name", "Arn", "EventPattern", "State", "Description", "ScheduleExpression", "RoleArn", "ManagedBy", "EventBusName"}},
			itemsAlways: map[string][]string{"Rules": {"Name", "Arn", "EventBusName", "State"}}},
		{name: "PutTargets", call: call("PutTargets", "EventBusName", "vocab-bus", "Rule", "vocab-rule",
			"Targets", []interface{}{map[string]interface{}{"Id": "t1", "Arn": sqsTargetARN}}),
			modelled: []string{"FailedEntryCount", "FailedEntries"},
			always:   []string{"FailedEntryCount", "FailedEntries"}},
		{name: "CreateArchive", call: call("CreateArchive", "ArchiveName", "vocab-archive", "EventSourceArn", busARN,
			"EventPattern", `{"source":["vocab.test"]}`, "Description", "vocab archive"),
			modelled: []string{"ArchiveArn", "State", "StateReason", "CreationTime"},
			always:   []string{"ArchiveArn", "State", "CreationTime"}},
		{name: "PutEvents", call: call("PutEvents", "Entries", []interface{}{map[string]interface{}{
			"EventBusName": "vocab-bus", "Source": "vocab.test", "DetailType": "VocabEvent", "Detail": `{"k":"v"}`}}),
			modelled:    []string{"FailedEntryCount", "Entries"},
			always:      []string{"FailedEntryCount", "Entries"},
			items:       map[string][]string{"Entries": {"EventId", "ErrorCode", "ErrorMessage"}},
			itemsAlways: map[string][]string{"Entries": {"EventId"}}},
		{name: "TestEventPattern", call: call("TestEventPattern", "EventPattern", `{"source":["vocab.test"]}`,
			"Event", `{"id":"1","account":"000000000000","source":"vocab.test","time":"2026-01-01T00:00:00Z","region":"us-east-1","resources":[],"detail-type":"VocabEvent","detail":{"k":"v"}}`),
			modelled: []string{"Result"},
			always:   []string{"Result"}},
		{name: "TagResource", call: call("TagResource", "ResourceARN", busARN,
			"Tags", []interface{}{map[string]interface{}{"Key": "env", "Value": "vocab"}}),
			modelled: []string{},
			always:   []string{}},
		{name: "ListTagsForResource", call: call("ListTagsForResource", "ResourceARN", busARN),
			modelled:    []string{"Tags"},
			always:      []string{"Tags"},
			items:       map[string][]string{"Tags": {"Key", "Value"}},
			itemsAlways: map[string][]string{"Tags": {"Key", "Value"}}},
		{name: "UntagResource", call: call("UntagResource", "ResourceARN", busARN, "TagKeys", []interface{}{"env"}),
			modelled: []string{},
			always:   []string{}},
		{name: "PutPermission", call: call("PutPermission", "EventBusName", "vocab-bus",
			"StatementId", "vocab-stmt", "Principal", "*", "Action", "events:PutEvents"),
			modelled: []string{},
			always:   []string{}},
		{name: "RemovePermission", call: call("RemovePermission", "EventBusName", "vocab-bus", "StatementId", "vocab-stmt"),
			modelled: []string{},
			always:   []string{}},
		{name: "DescribeArchive", call: call("DescribeArchive", "ArchiveName", "vocab-archive"),
			modelled: []string{"ArchiveArn", "ArchiveName", "EventSourceArn", "Description", "EventPattern", "State", "StateReason", "KmsKeyIdentifier", "RetentionDays", "SizeBytes", "EventCount", "CreationTime"},
			always:   []string{"ArchiveArn", "ArchiveName", "EventSourceArn", "Description", "EventPattern", "State", "CreationTime", "EventCount", "SizeBytes"}},
		{name: "UpdateArchive", call: call("UpdateArchive", "ArchiveName", "vocab-archive", "Description", "vocab archive v2"),
			modelled: []string{"ArchiveArn", "State", "StateReason", "CreationTime"},
			always:   []string{"ArchiveArn", "State", "CreationTime"}},
		{name: "ListArchives", call: call("ListArchives"),
			modelled:    []string{"Archives", "NextToken"},
			always:      []string{"Archives"},
			items:       map[string][]string{"Archives": {"ArchiveName", "EventSourceArn", "State", "StateReason", "RetentionDays", "SizeBytes", "EventCount", "CreationTime"}},
			itemsAlways: map[string][]string{"Archives": {"ArchiveName", "EventSourceArn", "State", "CreationTime", "EventCount", "SizeBytes"}}},
		{name: "StartReplay", call: deferred("StartReplay", func() []interface{} {
			return []interface{}{"ReplayName", "vocab-replay", "EventSourceArn", archiveArn,
				"EventStartTime", "2020-01-01T00:00:00Z", "EventEndTime", "2030-01-01T00:00:00Z",
				"Destination", map[string]interface{}{"Arn": busARN}}
		}),
			modelled: []string{"ReplayArn", "State", "StateReason", "ReplayStartTime"},
			always:   []string{"ReplayArn", "State"}},
		{name: "DescribeReplay", call: call("DescribeReplay", "ReplayName", "vocab-replay"),
			modelled: []string{"ReplayName", "ReplayArn", "Description", "State", "StateReason", "EventSourceArn", "Destination", "EventStartTime", "EventEndTime", "EventLastReplayedTime", "ReplayStartTime", "ReplayEndTime"},
			always:   []string{"ReplayName", "ReplayArn", "State", "EventSourceArn", "Destination", "EventStartTime", "EventEndTime"}},
		{name: "ListReplays", call: call("ListReplays"),
			modelled:    []string{"Replays", "NextToken"},
			always:      []string{"Replays"},
			items:       map[string][]string{"Replays": {"ReplayName", "EventSourceArn", "State", "StateReason", "EventStartTime", "EventEndTime", "EventLastReplayedTime", "ReplayStartTime", "ReplayEndTime"}},
			itemsAlways: map[string][]string{"Replays": {"ReplayName", "State", "EventSourceArn", "EventStartTime", "EventEndTime"}}},
		{name: "CreateConnection", call: call("CreateConnection", "Name", "vocab-conn", "AuthorizationType", "BASIC",
			"AuthParameters", map[string]interface{}{"BasicAuthParameters": map[string]interface{}{"Username": "u", "Password": "p"}}),
			modelled: []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime"},
			always:   []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime"}},
		{name: "DescribeConnection", call: call("DescribeConnection", "Name", "vocab-conn"),
			modelled: []string{"ConnectionArn", "Name", "Description", "InvocationConnectivityParameters", "ConnectionState", "StateReason", "AuthorizationType", "SecretArn", "KmsKeyIdentifier", "AuthParameters", "CreationTime", "LastModifiedTime", "LastAuthorizedTime"},
			always:   []string{"ConnectionArn", "Name", "ConnectionState", "AuthorizationType", "CreationTime"}},
		{name: "UpdateConnection", call: call("UpdateConnection", "Name", "vocab-conn", "Description", "vocab conn v2"),
			modelled: []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime", "LastAuthorizedTime"},
			always:   []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime"}},
		{name: "DeauthorizeConnection", call: call("DeauthorizeConnection", "Name", "vocab-conn"),
			modelled: []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime", "LastAuthorizedTime"},
			always:   []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime"}},
		{name: "ListConnections", call: call("ListConnections"),
			modelled:    []string{"Connections", "NextToken"},
			always:      []string{"Connections"},
			items:       map[string][]string{"Connections": {"ConnectionArn", "Name", "ConnectionState", "StateReason", "AuthorizationType", "CreationTime", "LastModifiedTime", "LastAuthorizedTime"}},
			itemsAlways: map[string][]string{"Connections": {"ConnectionArn", "Name", "ConnectionState", "AuthorizationType", "CreationTime"}}},
		{name: "CreateApiDestination", call: deferred("CreateApiDestination", func() []interface{} {
			return []interface{}{"Name", "vocab-dest", "ConnectionArn", connectionArn,
				"InvocationEndpoint", "https://example.com/", "HttpMethod", "POST", "Description", "vocab destination"}
		}),
			modelled: []string{"ApiDestinationArn", "ApiDestinationState", "CreationTime", "LastModifiedTime"},
			always:   []string{"ApiDestinationArn", "ApiDestinationState", "CreationTime", "LastModifiedTime"}},
		{name: "DescribeApiDestination", call: call("DescribeApiDestination", "Name", "vocab-dest"),
			modelled: []string{"ApiDestinationArn", "Name", "Description", "ApiDestinationState", "ConnectionArn", "InvocationEndpoint", "HttpMethod", "InvocationRateLimitPerSecond", "CreationTime", "LastModifiedTime"},
			always:   []string{"ApiDestinationArn", "Name", "Description", "ApiDestinationState", "ConnectionArn", "InvocationEndpoint", "HttpMethod", "CreationTime"}},
		{name: "UpdateApiDestination", call: call("UpdateApiDestination", "Name", "vocab-dest", "InvocationEndpoint", "https://example.com/v2"),
			modelled: []string{"ApiDestinationArn", "ApiDestinationState", "CreationTime", "LastModifiedTime"},
			always:   []string{"ApiDestinationArn", "ApiDestinationState", "CreationTime", "LastModifiedTime"}},
		{name: "ListApiDestinations", call: call("ListApiDestinations"),
			modelled:    []string{"ApiDestinations", "NextToken"},
			always:      []string{"ApiDestinations"},
			items:       map[string][]string{"ApiDestinations": {"ApiDestinationArn", "Name", "ApiDestinationState", "ConnectionArn", "InvocationEndpoint", "HttpMethod", "InvocationRateLimitPerSecond", "CreationTime", "LastModifiedTime"}},
			itemsAlways: map[string][]string{"ApiDestinations": {"ApiDestinationArn", "Name", "ApiDestinationState", "ConnectionArn", "InvocationEndpoint", "HttpMethod", "CreationTime"}}},
		{name: "ListRuleNamesByTarget", call: call("ListRuleNamesByTarget", "EventBusName", "vocab-bus", "TargetArn", sqsTargetARN),
			modelled: []string{"RuleNames", "NextToken"},
			always:   []string{"RuleNames"}},
		{name: "ListTargetsByRule", call: call("ListTargetsByRule", "EventBusName", "vocab-bus", "Rule", "vocab-rule"),
			modelled:    []string{"Targets", "NextToken"},
			always:      []string{"Targets"},
			items:       map[string][]string{"Targets": {"Id", "Arn", "RoleArn", "Input", "InputPath", "InputTransformer", "KinesisParameters", "RunCommandParameters", "EcsParameters", "BatchParameters", "SqsParameters", "HttpParameters", "RedshiftDataParameters", "SageMakerPipelineParameters", "DeadLetterConfig", "RetryPolicy", "AppSyncParameters"}},
			itemsAlways: map[string][]string{"Targets": {"Id", "Arn"}}},
		{name: "RemoveTargets", call: call("RemoveTargets", "EventBusName", "vocab-bus", "Rule", "vocab-rule", "Ids", []interface{}{"t1"}),
			modelled: []string{"FailedEntryCount", "FailedEntries"},
			always:   []string{"FailedEntryCount", "FailedEntries"}},
		{name: "DisableRule", call: call("DisableRule", "EventBusName", "vocab-bus", "Name", "vocab-rule"),
			modelled: []string{},
			always:   []string{}},
		{name: "EnableRule", call: call("EnableRule", "EventBusName", "vocab-bus", "Name", "vocab-rule"),
			modelled: []string{},
			always:   []string{}},
		{name: "DeleteArchive", call: call("DeleteArchive", "ArchiveName", "vocab-archive"),
			modelled: []string{},
			always:   []string{}},
		{name: "DeleteApiDestination", call: call("DeleteApiDestination", "Name", "vocab-dest"),
			modelled: []string{},
			always:   []string{}},
		{name: "DeleteConnection", call: call("DeleteConnection", "Name", "vocab-conn"),
			modelled: []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime", "LastAuthorizedTime"},
			always:   []string{"ConnectionArn", "ConnectionState", "CreationTime", "LastModifiedTime"}},
		{name: "DeleteRule", call: call("DeleteRule", "EventBusName", "vocab-bus", "Name", "vocab-rule"),
			modelled: []string{},
			always:   []string{}},
		{name: "DeleteEventBus", call: call("DeleteEventBus", "Name", "vocab-bus"),
			modelled: []string{},
			always:   []string{}},
	}

	// The capture loop reads the ARNs the CreateArchive/CreateConnection
	// responses emit; the deferred rows above consume them at call time.
	for _, row := range rows {
		resp, err := row.call(svc, reqCtx)
		if err != nil {
			t.Fatalf("%s: handler: %v", row.name, err)
		}
		assertVocabRow(t, row, resp)

		// Capture the wire values later rows reference.
		if m, ok := resp.(map[string]interface{}); ok {
			if v, ok := m["ArchiveArn"].(string); ok && archiveArn == "" {
				archiveArn = v
			}
			if v, ok := m["ConnectionArn"].(string); ok && connectionArn == "" {
				connectionArn = v
			}
		}
	}
}

// dispatchVocab routes a vocab row to the registered handler method by
// operation name — the same dispatch the HTTP plane performs.
func dispatchVocab(ctx context.Context, svc *EventsService, op string, params map[string]interface{}, reqCtx *request.RequestContext) (interface{}, error) {
	req := &request.ParsedRequest{Parameters: params}
	switch op {
	case "CreateEventBus":
		return svc.CreateEventBus(ctx, reqCtx, req)
	case "DeleteEventBus":
		return svc.DeleteEventBus(ctx, reqCtx, req)
	case "DescribeEventBus":
		return svc.DescribeEventBus(ctx, reqCtx, req)
	case "ListEventBuses":
		return svc.ListEventBuses(ctx, reqCtx, req)
	case "UpdateEventBus":
		return svc.UpdateEventBus(ctx, reqCtx, req)
	case "PutRule":
		return svc.PutRule(ctx, reqCtx, req)
	case "DeleteRule":
		return svc.DeleteRule(ctx, reqCtx, req)
	case "DescribeRule":
		return svc.DescribeRule(ctx, reqCtx, req)
	case "ListRules":
		return svc.ListRules(ctx, reqCtx, req)
	case "EnableRule":
		return svc.EnableRule(ctx, reqCtx, req)
	case "DisableRule":
		return svc.DisableRule(ctx, reqCtx, req)
	case "PutTargets":
		return svc.PutTargets(ctx, reqCtx, req)
	case "RemoveTargets":
		return svc.RemoveTargets(ctx, reqCtx, req)
	case "ListTargetsByRule":
		return svc.ListTargetsByRule(ctx, reqCtx, req)
	case "ListRuleNamesByTarget":
		return svc.ListRuleNamesByTarget(ctx, reqCtx, req)
	case "PutEvents":
		return svc.PutEvents(ctx, reqCtx, req)
	case "TagResource":
		return svc.TagResource(ctx, reqCtx, req)
	case "UntagResource":
		return svc.UntagResource(ctx, reqCtx, req)
	case "ListTagsForResource":
		return svc.ListTagsForResource(ctx, reqCtx, req)
	case "CreateArchive":
		return svc.CreateArchive(ctx, reqCtx, req)
	case "DeleteArchive":
		return svc.DeleteArchive(ctx, reqCtx, req)
	case "DescribeArchive":
		return svc.DescribeArchive(ctx, reqCtx, req)
	case "UpdateArchive":
		return svc.UpdateArchive(ctx, reqCtx, req)
	case "ListArchives":
		return svc.ListArchives(ctx, reqCtx, req)
	case "StartReplay":
		return svc.StartReplay(ctx, reqCtx, req)
	case "DescribeReplay":
		return svc.DescribeReplay(ctx, reqCtx, req)
	case "ListReplays":
		return svc.ListReplays(ctx, reqCtx, req)
	case "CancelReplay":
		return svc.CancelReplay(ctx, reqCtx, req)
	case "CreateConnection":
		return svc.CreateConnection(ctx, reqCtx, req)
	case "DeleteConnection":
		return svc.DeleteConnection(ctx, reqCtx, req)
	case "DescribeConnection":
		return svc.DescribeConnection(ctx, reqCtx, req)
	case "UpdateConnection":
		return svc.UpdateConnection(ctx, reqCtx, req)
	case "DeauthorizeConnection":
		return svc.DeauthorizeConnection(ctx, reqCtx, req)
	case "ListConnections":
		return svc.ListConnections(ctx, reqCtx, req)
	case "CreateApiDestination":
		return svc.CreateApiDestination(ctx, reqCtx, req)
	case "DeleteApiDestination":
		return svc.DeleteApiDestination(ctx, reqCtx, req)
	case "DescribeApiDestination":
		return svc.DescribeApiDestination(ctx, reqCtx, req)
	case "UpdateApiDestination":
		return svc.UpdateApiDestination(ctx, reqCtx, req)
	case "ListApiDestinations":
		return svc.ListApiDestinations(ctx, reqCtx, req)
	case "TestEventPattern":
		return svc.TestEventPattern(ctx, reqCtx, req)
	case "PutPermission":
		return svc.PutPermission(ctx, reqCtx, req)
	case "RemovePermission":
		return svc.RemovePermission(ctx, reqCtx, req)
	}
	return nil, nil
}

// assertVocabRow checks one response against its model transcriptions.
func assertVocabRow(t *testing.T, row vocabRow, resp interface{}) {
	t.Helper()
	m, ok := resp.(map[string]interface{})
	if !ok {
		t.Fatalf("%s: response is %T, want a map", row.name, resp)
	}
	modelled := make(map[string]bool, len(row.modelled))
	for _, k := range row.modelled {
		modelled[k] = true
	}
	var emitted []string
	for k := range m {
		emitted = append(emitted, k)
		if !modelled[k] {
			t.Errorf("%s: emitted key %q is not modelled in the output shape (modelled: %v)", row.name, k, row.modelled)
		}
	}
	sort.Strings(emitted)
	for _, k := range row.always {
		if _, ok := m[k]; !ok {
			t.Errorf("%s: modelled member %q missing from the response (emitted: %v)", row.name, k, emitted)
		}
	}
	for member, itemMembers := range row.items {
		// List members serialise as []map[string]interface{}; the tag
		// family's shared builder emits []map[string]string.
		type stringItem struct {
			keys   []string
			hasKey func(string) bool
		}
		var items []stringItem
		switch list := m[member].(type) {
		case []map[string]interface{}:
			for _, item := range list {
				item := item
				var keys []string
				for k := range item {
					keys = append(keys, k)
				}
				items = append(items, stringItem{keys, func(k string) bool { _, ok := item[k]; return ok }})
			}
		case []map[string]string:
			for _, item := range list {
				item := item
				var keys []string
				for k := range item {
					keys = append(keys, k)
				}
				items = append(items, stringItem{keys, func(k string) bool { _, ok := item[k]; return ok }})
			}
		default:
			t.Errorf("%s: %s list absent or empty (got %T)", row.name, member, m[member])
			continue
		}
		if len(items) == 0 {
			t.Errorf("%s: %s list is empty", row.name, member)
			continue
		}
		allowed := make(map[string]bool, len(itemMembers))
		for _, k := range itemMembers {
			allowed[k] = true
		}
		for i, item := range items {
			for _, k := range item.keys {
				if !allowed[k] {
					t.Errorf("%s: %s item %d emits unmodelled key %q (item shape: %v)", row.name, member, i, k, itemMembers)
				}
			}
			for _, k := range row.itemsAlways[member] {
				if !item.hasKey(k) {
					t.Errorf("%s: %s item %d missing modelled member %q", row.name, member, i, k)
				}
			}
		}
	}
}

// TestCancelReplayResponseVocabulary pins the CancelReplay response keys on
// a replay no worker owns (seeded directly in the store, still STARTING),
// keeping the row deterministic: cancel flips it to CANCELLING and the
// CANCELLED terminal write belongs to the worker that does not exist here.
func TestCancelReplayResponseVocabulary(t *testing.T) {
	ctx := context.Background()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(ctx, mgr, "000000000000", "us-east-1")
	svc := NewEventsService(nil, "000000000000")
	t.Cleanup(svc.Close)

	store, err := svc.store(reqCtx)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.CreateReplayCapped(ctx, &eventsstore.Replay{
		Name: "vocab-cancel", State: eventsstore.ReplayStateStarting,
	}, eventsstore.MaxConcurrentReplays); err != nil {
		t.Fatal(err)
	}

	resp, err := svc.CancelReplay(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"ReplayName": "vocab-cancel",
	}})
	if err != nil {
		t.Fatal(err)
	}
	assertVocabRow(t, vocabRow{
		name:     "CancelReplay",
		modelled: []string{"ReplayArn", "State", "StateReason"},
		always:   []string{"ReplayArn", "State", "StateReason"},
	}, resp)
}
