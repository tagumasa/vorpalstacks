package scheduler

import (
	"context"
	"sort"
	"testing"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/storage"
)

// The response-key vocabulary pin for every registered operation: each
// operation's emitted top-level member keys, with exact casing, against
// the output-shape member names of the vendored scheduler model
// (2021-06-30, transcribed into the table below — test files never read
// the vendored model). The rule: an emitted key the output shape does not
// model is a wire-contract violation however harmless it looks (SDKs
// silently drop it, and the drift becomes unpinnable), and a modelled
// member the scenario guarantees must actually be present. The Target
// tree is checked the same way through GetSchedule, one schedule per
// parameter family — the full-fidelity ECS schedule exercises every
// EcsParameters member including the model's deliberately lowercase
// member names (awsvpcConfiguration; capacityProvider/weight/base;
// placement expression/type/field), where a casing slip is the classic
// silent drift.
//
// The rows run as one scenario in table order against the real handler
// path (wire parse → Core → serialisation) over a real per-region store,
// so a key added to any serialiser without a model member — or dropped
// from one — fails here before it can drift further. Any fix that
// changes response construction extends this table in the same commit.

type vocabRow struct {
	name string
	call func(svc *SchedulerService, reqCtx *request.RequestContext) (interface{}, error)

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
	// after runs row-specific deep assertions (the Target tree, list-item
	// structures) once the top-level vocabulary has passed.
	after func(t *testing.T, resp map[string]interface{})
}

// The Target shape (com.amazonaws.scheduler#Target): exactly ten members.
var modelledTargetMembers = []string{
	"Arn", "DeadLetterConfig", "EcsParameters", "EventBridgeParameters",
	"Input", "KinesisParameters", "RetryPolicy", "RoleArn",
	"SageMakerPipelineParameters", "SqsParameters",
}

// The EcsParameters shape: exactly fourteen members.
var modelledEcsMembers = []string{
	"CapacityProviderStrategy", "EnableECSManagedTags", "EnableExecuteCommand",
	"Group", "LaunchType", "NetworkConfiguration", "PlacementConstraints",
	"PlacementStrategy", "PlatformVersion", "PropagateTags", "ReferenceId",
	"Tags", "TaskCount", "TaskDefinitionArn",
}

func TestResponseKeyVocabulary(t *testing.T) {
	ctx := context.Background()
	mgr, err := storage.NewRegionStorageManager(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	reqCtx := request.NewRequestContext(ctx, mgr, "000000000000", "us-east-1")

	svc := NewSchedulerService(mgr, "000000000000")
	t.Cleanup(func() { _ = svc.StopEngine() })

	const roleArn = "arn:aws:iam::000000000000:role/vocab-role"
	const kmsKeyArn = "arn:aws:kms:us-east-1:000000000000:key/1111aabb-2222-cc33-4d44-5555eee66666"

	// Captured wire values later rows depend on; the rows that use them
	// read the variable at call time, after the capturing row has run.
	var groupArn string

	params := func(kv ...interface{}) map[string]interface{} {
		m := make(map[string]interface{}, len(kv)/2)
		for i := 0; i+1 < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		return m
	}
	call := func(op string, kv ...interface{}) func(*SchedulerService, *request.RequestContext) (interface{}, error) {
		return func(svc *SchedulerService, reqCtx *request.RequestContext) (interface{}, error) {
			return dispatchVocab(context.Background(), svc, op, params(kv...), reqCtx)
		}
	}
	deferred := func(op string, kv func() []interface{}) func(*SchedulerService, *request.RequestContext) (interface{}, error) {
		return func(svc *SchedulerService, reqCtx *request.RequestContext) (interface{}, error) {
			return dispatchVocab(context.Background(), svc, op, params(kv()...), reqCtx)
		}
	}

	// ecsTarget builds the full-fidelity ECS target: every Target member
	// the platform carries and all fourteen EcsParameters members,
	// including the lowercase nested member names the model defines.
	ecsTarget := func() map[string]interface{} {
		return map[string]interface{}{
			"Arn":     "arn:aws:ecs:us-east-1:000000000000:task-definition/vocab-family",
			"RoleArn": roleArn,
			"Input":   `{"vocab":true}`,
			"DeadLetterConfig": map[string]interface{}{
				"Arn": "arn:aws:sqs:us-east-1:000000000000:vocab-dlq",
			},
			"RetryPolicy": map[string]interface{}{
				"MaximumEventAgeInSeconds": 3600,
				"MaximumRetryAttempts":     24,
			},
			"EcsParameters": map[string]interface{}{
				"TaskDefinitionArn":    "arn:aws:ecs:us-east-1:000000000000:task-definition/vocab-family:1",
				"LaunchType":           "FARGATE",
				"PlatformVersion":      "LATEST",
				"Group":                "vocab-task-group",
				"TaskCount":            2,
				"PropagateTags":        "TASK_DEFINITION",
				"ReferenceId":          "vocab-reference",
				"EnableECSManagedTags": true,
				"EnableExecuteCommand": true,
				"NetworkConfiguration": map[string]interface{}{
					"awsvpcConfiguration": map[string]interface{}{
						"Subnets":        []interface{}{"subnet-vocab0001"},
						"SecurityGroups": []interface{}{"sg-vocab0001"},
						"AssignPublicIp": "ENABLED",
					},
				},
				"CapacityProviderStrategy": []interface{}{map[string]interface{}{
					"capacityProvider": "FARGATE", "weight": 1, "base": 1,
				}},
				"PlacementConstraints": []interface{}{map[string]interface{}{
					"type": "memberOf", "expression": "attribute:ecs.ami-id",
				}},
				"PlacementStrategy": []interface{}{map[string]interface{}{
					"type": "spread", "field": "attribute:ecs.availability-zone",
				}},
				// EcsParameters.Tags is a list of TagMap shapes — each item
				// is a map of tag key to value, not a Key/Value structure.
				// This scenario emits single-pair entries; the bounds
				// validation admits multi-pair TagMaps.
				"Tags": []interface{}{map[string]interface{}{"env": "vocab"}},
			},
		}
	}
	scheduleBody := func(name string, target map[string]interface{}, ftw map[string]interface{}) []interface{} {
		return []interface{}{
			"Name", name,
			"GroupName", "vocab-group",
			"ScheduleExpression", "rate(1 hour)",
			"ScheduleExpressionTimezone", "UTC",
			"Description", "vocabulary pin schedule " + name,
			"State", "ENABLED",
			"KmsKeyArn", kmsKeyArn,
			"StartDate", "2030-01-01T00:00:00Z",
			"EndDate", "2031-01-01T00:00:00Z",
			"ActionAfterCompletion", "NONE",
			"FlexibleTimeWindow", ftw,
			"Target", target,
		}
	}
	flexible := map[string]interface{}{"Mode": "FLEXIBLE", "MaximumWindowInMinutes": 5}
	off := map[string]interface{}{"Mode": "OFF"}

	// assertTarget pins one emitted Target map against the Target shape's
	// member vocabulary, then recurses into every parameter family the
	// scenario emitted.
	assertTarget := func(context string, target map[string]interface{}, always []string) {
		assertKeySubsetAndAlways(t, context, target, modelledTargetMembers, always)
		if rp, ok := target["RetryPolicy"].(map[string]interface{}); ok {
			assertKeySubsetAndAlways(t, context+".RetryPolicy", rp,
				[]string{"MaximumEventAgeInSeconds", "MaximumRetryAttempts"},
				[]string{"MaximumEventAgeInSeconds", "MaximumRetryAttempts"})
		}
		if dlc, ok := target["DeadLetterConfig"].(map[string]interface{}); ok {
			assertKeySubsetAndAlways(t, context+".DeadLetterConfig", dlc, []string{"Arn"}, []string{"Arn"})
		}
		if sqs, ok := target["SqsParameters"].(map[string]interface{}); ok {
			assertKeySubsetAndAlways(t, context+".SqsParameters", sqs, []string{"MessageGroupId"}, []string{"MessageGroupId"})
		}
		if eb, ok := target["EventBridgeParameters"].(map[string]interface{}); ok {
			assertKeySubsetAndAlways(t, context+".EventBridgeParameters", eb,
				[]string{"DetailType", "Source"}, []string{"DetailType", "Source"})
		}
		if kin, ok := target["KinesisParameters"].(map[string]interface{}); ok {
			assertKeySubsetAndAlways(t, context+".KinesisParameters", kin,
				[]string{"PartitionKey"}, []string{"PartitionKey"})
		}
		if ecs, ok := target["EcsParameters"].(map[string]interface{}); ok {
			assertKeySubsetAndAlways(t, context+".EcsParameters", ecs, modelledEcsMembers, modelledEcsMembers)
			if nc, ok := ecs["NetworkConfiguration"].(map[string]interface{}); ok {
				assertKeySubsetAndAlways(t, context+".NetworkConfiguration", nc,
					[]string{"awsvpcConfiguration"}, []string{"awsvpcConfiguration"})
				if vpc, ok := nc["awsvpcConfiguration"].(map[string]interface{}); ok {
					assertKeySubsetAndAlways(t, context+".awsvpcConfiguration", vpc,
						[]string{"AssignPublicIp", "SecurityGroups", "Subnets"},
						[]string{"AssignPublicIp", "SecurityGroups", "Subnets"})
				}
			}
			for _, member := range []string{"CapacityProviderStrategy", "PlacementConstraints", "PlacementStrategy"} {
				list, ok := ecs[member].([]map[string]interface{})
				if !ok || len(list) == 0 {
					t.Errorf("%s.%s absent or empty", context, member)
					continue
				}
				var modelled, always []string
				switch member {
				case "CapacityProviderStrategy":
					modelled, always = []string{"base", "capacityProvider", "weight"}, []string{"base", "capacityProvider", "weight"}
				case "PlacementConstraints":
					modelled, always = []string{"expression", "type"}, []string{"expression", "type"}
				case "PlacementStrategy":
					modelled, always = []string{"field", "type"}, []string{"field", "type"}
				}
				for _, item := range list {
					assertKeySubsetAndAlways(t, context+"."+member+"[]", item, modelled, always)
				}
			}
			// EcsParameters.Tags items are TagMap shapes. This scenario
			// emits single-pair entries; multi-pair TagMaps are valid input
			// (adjudicated against the model and documented constraints).
			if tags, ok := ecs["Tags"].([]map[string]string); ok {
				for i, tag := range tags {
					if len(tag) != 1 {
						t.Errorf("%s.Tags[%d] has %d pairs, want the single pair this scenario emits", context, i, len(tag))
					}
				}
			} else {
				t.Errorf("%s.Tags absent or not a TagMap list (got %T)", context, ecs["Tags"])
			}
		}
	}

	rows := []vocabRow{
		{name: "CreateScheduleGroup", call: call("CreateScheduleGroup", "Name", "vocab-group",
			"Tags", []interface{}{map[string]interface{}{"Key": "env", "Value": "vocab"}}),
			modelled: []string{"ScheduleGroupArn"},
			always:   []string{"ScheduleGroupArn"},
			after: func(t *testing.T, resp map[string]interface{}) {
				groupArn, _ = resp["ScheduleGroupArn"].(string)
			}},
		{name: "CreateScheduleGroup (deletable)", call: call("CreateScheduleGroup", "Name", "vocab-todelete"),
			modelled: []string{"ScheduleGroupArn"},
			always:   []string{"ScheduleGroupArn"}},
		{name: "GetScheduleGroup", call: call("GetScheduleGroup", "Name", "vocab-group"),
			modelled: []string{"Arn", "CreationDate", "LastModificationDate", "Name", "State"},
			always:   []string{"Arn", "CreationDate", "LastModificationDate", "Name", "State"}},
		{name: "ListScheduleGroups", call: call("ListScheduleGroups"),
			modelled:    []string{"NextToken", "ScheduleGroups"},
			always:      []string{"ScheduleGroups"},
			items:       map[string][]string{"ScheduleGroups": {"Arn", "CreationDate", "LastModificationDate", "Name", "State"}},
			itemsAlways: map[string][]string{"ScheduleGroups": {"Arn", "CreationDate", "LastModificationDate", "Name", "State"}}},
		{name: "CreateSchedule (full ECS target)", call: call("CreateSchedule", scheduleBody("vocab-ecs", ecsTarget(), flexible)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		{name: "CreateSchedule (sqs target)", call: call("CreateSchedule", scheduleBody("vocab-sqs", map[string]interface{}{
			"Arn": "arn:aws:sqs:us-east-1:000000000000:vocab-target", "RoleArn": roleArn, "Input": `{"vocab":true}`,
			"SqsParameters": map[string]interface{}{"MessageGroupId": "vocab-group-id"},
		}, flexible)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		{name: "CreateSchedule (events target)", call: call("CreateSchedule", scheduleBody("vocab-events", map[string]interface{}{
			"Arn": "arn:aws:events:us-east-1:000000000000:event-bus/vocab-bus", "RoleArn": roleArn,
			"EventBridgeParameters": map[string]interface{}{"Source": "vocab.test", "DetailType": "VocabEvent"},
		}, flexible)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		{name: "CreateSchedule (kinesis target)", call: call("CreateSchedule", scheduleBody("vocab-kinesis", map[string]interface{}{
			"Arn": "arn:aws:kinesis:us-east-1:000000000000:stream/vocab-stream", "RoleArn": roleArn,
			"KinesisParameters": map[string]interface{}{"PartitionKey": "vocab-key"},
		}, flexible)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		{name: "CreateSchedule (minimal lambda target)", call: call("CreateSchedule", scheduleBody("vocab-min", map[string]interface{}{
			"Arn": "arn:aws:lambda:us-east-1:000000000000:function:vocab-function", "RoleArn": roleArn,
		}, off)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		// An EcsParameters whose NetworkConfiguration carries no
		// awsvpcConfiguration member: the absent structure must vanish on
		// the round trip — omitted from the response, never emitted as a
		// null member.
		{name: "CreateSchedule (ECS target, empty networkConfiguration)", call: call("CreateSchedule", scheduleBody("vocab-ecsnc", map[string]interface{}{
			"Arn":     "arn:aws:ecs:us-east-1:000000000000:task-definition/vocab-family",
			"RoleArn": roleArn,
			"EcsParameters": map[string]interface{}{
				"TaskDefinitionArn":    "arn:aws:ecs:us-east-1:000000000000:task-definition/vocab-family:1",
				"NetworkConfiguration": map[string]interface{}{},
			},
		}, off)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		// SqsParameters without MessageGroupId: the empty structure is
		// valid input (delivery falls back to the schedule name), and the
		// unset member is omitted from the read-back rather than emitted
		// as an empty string.
		{name: "CreateSchedule (sqs target, empty SqsParameters)", call: call("CreateSchedule", scheduleBody("vocab-sqsnc", map[string]interface{}{
			"Arn": "arn:aws:sqs:us-east-1:000000000000:vocab-target-nc", "RoleArn": roleArn,
			"SqsParameters": map[string]interface{}{},
		}, flexible)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		{name: "GetSchedule (full ECS target)", call: call("GetSchedule", "Name", "vocab-ecs", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertTarget("GetSchedule(full).Target", target,
					[]string{"Arn", "DeadLetterConfig", "EcsParameters", "Input", "RetryPolicy", "RoleArn"})
				ftw, _ := resp["FlexibleTimeWindow"].(map[string]interface{})
				assertKeySubsetAndAlways(t, "GetSchedule(full).FlexibleTimeWindow", ftw,
					[]string{"MaximumWindowInMinutes", "Mode"},
					[]string{"MaximumWindowInMinutes", "Mode"})
			}},
		{name: "GetSchedule (minimal lambda target)", call: call("GetSchedule", "Name", "vocab-min", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{"Arn", "CreationDate", "GroupName", "LastModificationDate", "Name",
				"ScheduleExpression", "State", "Target", "FlexibleTimeWindow"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertTarget("GetSchedule(min).Target", target, []string{"Arn", "RoleArn"})
				ftw, _ := resp["FlexibleTimeWindow"].(map[string]interface{})
				assertKeySubsetAndAlways(t, "GetSchedule(min).FlexibleTimeWindow", ftw,
					[]string{"MaximumWindowInMinutes", "Mode"}, []string{"Mode"})
			}},
		{name: "GetSchedule (sqs target)", call: call("GetSchedule", "Name", "vocab-sqs", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{"Target"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertTarget("GetSchedule(sqs).Target", target, []string{"Arn", "Input", "RoleArn", "SqsParameters"})
			}},
		{name: "GetSchedule (events target)", call: call("GetSchedule", "Name", "vocab-events", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{"Target"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertTarget("GetSchedule(events).Target", target, []string{"Arn", "EventBridgeParameters", "RoleArn"})
			}},
		{name: "GetSchedule (kinesis target)", call: call("GetSchedule", "Name", "vocab-kinesis", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{"Target"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertTarget("GetSchedule(kinesis).Target", target, []string{"Arn", "KinesisParameters", "RoleArn"})
			}},
		{name: "GetSchedule (ECS target, empty networkConfiguration)", call: call("GetSchedule", "Name", "vocab-ecsnc", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{"Target"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertKeySubsetAndAlways(t, "GetSchedule(ecsnc).Target", target,
					modelledTargetMembers, []string{"Arn", "EcsParameters", "RoleArn"})
				ecs, ok := target["EcsParameters"].(map[string]interface{})
				if !ok {
					t.Fatalf("GetSchedule(ecsnc).EcsParameters is %T, want a map", target["EcsParameters"])
				}
				assertKeySubsetAndAlways(t, "GetSchedule(ecsnc).EcsParameters", ecs,
					modelledEcsMembers, []string{"TaskDefinitionArn"})
				if _, ok := ecs["NetworkConfiguration"]; ok {
					t.Errorf("NetworkConfiguration must be omitted when awsvpcConfiguration is absent; got %v (%T)",
						ecs["NetworkConfiguration"], ecs["NetworkConfiguration"])
				}
			}},
		{name: "GetSchedule (sqs target, empty SqsParameters)", call: call("GetSchedule", "Name", "vocab-sqsnc", "groupName", "vocab-group"),
			modelled: []string{
				"ActionAfterCompletion", "Arn", "CreationDate", "Description", "EndDate",
				"FlexibleTimeWindow", "GroupName", "KmsKeyArn", "LastModificationDate", "Name",
				"ScheduleExpression", "ScheduleExpressionTimezone", "StartDate", "State", "Target"},
			always: []string{"Target"},
			after: func(t *testing.T, resp map[string]interface{}) {
				target, _ := resp["Target"].(map[string]interface{})
				assertKeySubsetAndAlways(t, "GetSchedule(sqsnc).Target", target,
					modelledTargetMembers, []string{"Arn", "RoleArn", "SqsParameters"})
				sqs, ok := target["SqsParameters"].(map[string]interface{})
				if !ok {
					t.Fatalf("GetSchedule(sqsnc).SqsParameters is %T, want a map", target["SqsParameters"])
				}
				assertKeySubsetAndAlways(t, "GetSchedule(sqsnc).SqsParameters", sqs,
					[]string{"MessageGroupId"}, nil)
				if _, ok := sqs["MessageGroupId"]; ok {
					t.Errorf("MessageGroupId must be omitted when unset; got %v (%T)",
						sqs["MessageGroupId"], sqs["MessageGroupId"])
				}
			}},
		{name: "UpdateSchedule", call: call("UpdateSchedule", scheduleBody("vocab-min", map[string]interface{}{
			"Arn": "arn:aws:lambda:us-east-1:000000000000:function:vocab-function", "RoleArn": roleArn,
		}, off)...),
			modelled: []string{"ScheduleArn"},
			always:   []string{"ScheduleArn"}},
		{name: "ListSchedules", call: call("ListSchedules", "ScheduleGroup", "vocab-group"),
			modelled:    []string{"NextToken", "Schedules"},
			always:      []string{"Schedules"},
			items:       map[string][]string{"Schedules": {"Arn", "CreationDate", "GroupName", "LastModificationDate", "Name", "State", "Target"}},
			itemsAlways: map[string][]string{"Schedules": {"Arn", "CreationDate", "GroupName", "LastModificationDate", "Name", "State", "Target"}},
			after: func(t *testing.T, resp map[string]interface{}) {
				// ScheduleSummary.Target is a TargetSummary: exactly one
				// member, Arn.
				list, _ := resp["Schedules"].([]map[string]interface{})
				for i, item := range list {
					summary, ok := item["Target"].(map[string]interface{})
					if !ok {
						t.Fatalf("Schedules[%d].Target is %T, want a map", i, item["Target"])
					}
					assertKeySubsetAndAlways(t, "Schedules[].Target", summary, []string{"Arn"}, []string{"Arn"})
				}
			}},
		{name: "TagResource", call: deferred("TagResource", func() []interface{} {
			return []interface{}{"ResourceArn", groupArn,
				"Tags", []interface{}{map[string]interface{}{"Key": "team", "Value": "vocab"}}}
		}),
			modelled: []string{},
			always:   []string{}},
		{name: "ListTagsForResource", call: deferred("ListTagsForResource", func() []interface{} {
			return []interface{}{"ResourceArn", groupArn}
		}),
			modelled:    []string{"Tags"},
			always:      []string{"Tags"},
			items:       map[string][]string{"Tags": {"Key", "Value"}},
			itemsAlways: map[string][]string{"Tags": {"Key", "Value"}}},
		{name: "UntagResource", call: deferred("UntagResource", func() []interface{} {
			return []interface{}{"ResourceArn", groupArn, "TagKeys", []interface{}{"team"}}
		}),
			modelled: []string{},
			always:   []string{}},
		{name: "DeleteSchedule", call: call("DeleteSchedule", "Name", "vocab-min", "groupName", "vocab-group"),
			modelled: []string{},
			always:   []string{}},
		{name: "DeleteScheduleGroup", call: call("DeleteScheduleGroup", "Name", "vocab-todelete", "clientToken", "vocab-delete-token"),
			modelled: []string{},
			always:   []string{}},
	}

	for _, row := range rows {
		resp, err := row.call(svc, reqCtx)
		if err != nil {
			t.Fatalf("%s: handler: %v", row.name, err)
		}
		assertVocabRow(t, row, resp)
		if row.after != nil {
			if m, ok := resp.(map[string]interface{}); ok {
				row.after(t, m)
			}
		}
	}
}

// dispatchVocab routes a vocab row to the registered handler method by
// operation name — the same dispatch the HTTP plane performs.
func dispatchVocab(ctx context.Context, svc *SchedulerService, op string, params map[string]interface{}, reqCtx *request.RequestContext) (interface{}, error) {
	req := &request.ParsedRequest{Parameters: params}
	switch op {
	case "CreateSchedule":
		return svc.CreateSchedule(ctx, reqCtx, req)
	case "DeleteSchedule":
		return svc.DeleteSchedule(ctx, reqCtx, req)
	case "GetSchedule":
		return svc.GetSchedule(ctx, reqCtx, req)
	case "UpdateSchedule":
		return svc.UpdateSchedule(ctx, reqCtx, req)
	case "ListSchedules":
		return svc.ListSchedules(ctx, reqCtx, req)
	case "CreateScheduleGroup":
		return svc.CreateScheduleGroup(ctx, reqCtx, req)
	case "DeleteScheduleGroup":
		return svc.DeleteScheduleGroup(ctx, reqCtx, req)
	case "GetScheduleGroup":
		return svc.GetScheduleGroup(ctx, reqCtx, req)
	case "ListScheduleGroups":
		return svc.ListScheduleGroups(ctx, reqCtx, req)
	case "TagResource":
		return svc.TagResource(ctx, reqCtx, req)
	case "UntagResource":
		return svc.UntagResource(ctx, reqCtx, req)
	case "ListTagsForResource":
		return svc.ListTagsForResource(ctx, reqCtx, req)
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
	assertKeySubsetAndAlways(t, row.name, m, row.modelled, row.always)
	for member, itemMembers := range row.items {
		list, ok := m[member].([]map[string]interface{})
		if !ok || len(list) == 0 {
			t.Errorf("%s: %s list absent or empty (got %T)", row.name, member, m[member])
			continue
		}
		for _, item := range list {
			assertKeySubsetAndAlways(t, row.name+"."+member+"[]", item, itemMembers, row.itemsAlways[member])
		}
	}
}

// assertKeySubsetAndAlways enforces the vocabulary rule on one map: every
// emitted key must be modelled, and every always member must be emitted.
func assertKeySubsetAndAlways(t *testing.T, context string, m map[string]interface{}, modelled, always []string) {
	t.Helper()
	if m == nil {
		t.Fatalf("%s: map is nil", context)
	}
	modelledSet := make(map[string]bool, len(modelled))
	for _, k := range modelled {
		modelledSet[k] = true
	}
	var emitted []string
	for k := range m {
		emitted = append(emitted, k)
		if !modelledSet[k] {
			t.Errorf("%s: emitted key %q is not modelled in the shape (modelled: %v)", context, k, modelled)
		}
	}
	sort.Strings(emitted)
	for _, k := range always {
		if _, ok := m[k]; !ok {
			t.Errorf("%s: modelled member %q missing from the response (emitted: %v)", context, k, emitted)
		}
	}
}
