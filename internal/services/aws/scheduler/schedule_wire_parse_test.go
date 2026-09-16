package scheduler

import (
	"testing"

	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// The wire-parse pins for the modelled member names and strict typing:
// every parser binds exactly the member names the model declares (Pascal
// at the Target level and inside EcsParameters/AwsVpcConfiguration/
// FlexibleTimeWindow; the model's own lowerCamel names for
// awsvpcConfiguration, the capacity-provider and placement items), a
// present member of the wrong JSON type is a ValidationException rather
// than a silently dropped value, and the FlexibleTimeWindow string branch
// binds the same members as the map branch.

// A non-model spelling is an unknown member: the parse reads only the
// modelled key, so the value stays unbound and the Core validation rejects
// the required member.
func TestParseTargetModelledKeysOnly(t *testing.T) {
	target, err := parseTargetFromMap(map[string]interface{}{
		"arn":     "arn:aws:lambda:us-east-1:123456789012:function:f",
		"roleArn": "arn:aws:iam::123456789012:role/r",
	})
	if err != nil {
		t.Fatalf("lowercase dialect parse: %v", err)
	}
	if target.Arn != "" || target.RoleArn != "" {
		t.Fatalf("non-model spellings bound values: %+v", target)
	}

	// The modelled Pascal names bind.
	target, err = parseTargetFromMap(map[string]interface{}{
		"Arn":     "arn:aws:lambda:us-east-1:123456789012:function:f",
		"RoleArn": "arn:aws:iam::123456789012:role/r",
		"Input":   `{"pin":true}`,
	})
	if err != nil {
		t.Fatalf("Pascal parse: %v", err)
	}
	if target.Arn == "" || target.RoleArn == "" || target.Input != `{"pin":true}` {
		t.Fatalf("modelled spellings not bound: %+v", target)
	}
}

// A present SageMakerPipelineParameters member is rejected like every other
// sub-parameter the accepted target's service cannot use: SageMaker targets
// are permanently out of scope, so the member is never a carryable value —
// reported here, not silently dropped where the Core validator cannot see it.
func TestParseTargetRejectsSageMakerPipelineParameters(t *testing.T) {
	_, err := parseTargetFromMap(map[string]interface{}{
		"Arn":                         "arn:aws:lambda:us-east-1:123456789012:function:f",
		"RoleArn":                     "arn:aws:iam::123456789012:role/r",
		"SageMakerPipelineParameters": map[string]interface{}{},
	})
	if err == nil {
		t.Fatal("a present SageMakerPipelineParameters member was silently dropped")
	}
	assertValidationMessage(t, err, "SageMaker targets are not supported")
}

// A present member whose JSON type contradicts the shape is reported with
// the member's path, in every member class.
func TestParseTargetStrictTypes(t *testing.T) {
	cases := []struct {
		name string
		m    map[string]interface{}
		frag string
	}{
		{"string member as number", map[string]interface{}{"Arn": 123}, "Target.Arn must be a string"},
		{"structure member as string", map[string]interface{}{"RetryPolicy": "none"}, "Target.RetryPolicy must be a structure"},
		{"list member as string", map[string]interface{}{
			"EcsParameters": map[string]interface{}{"Tags": "none"},
		}, "EcsParameters.Tags must be a list"},
		{"number member as string", map[string]interface{}{
			"RetryPolicy": map[string]interface{}{"MaximumEventAgeInSeconds": "soon"},
		}, "RetryPolicy.MaximumEventAgeInSeconds must be a number"},
		{"boolean member as string", map[string]interface{}{
			"EcsParameters": map[string]interface{}{"EnableExecuteCommand": "yes"},
		}, "EcsParameters.EnableExecuteCommand must be a boolean"},
		// A fractional JSON number is a wire violation for a modelled
		// integer shape — reported, not truncated into a reshaped value.
		{"integer member as fraction", map[string]interface{}{
			"EcsParameters": map[string]interface{}{"TaskCount": 5.7},
		}, "EcsParameters.TaskCount must be an integer"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseTargetFromMap(tc.m)
			if err == nil {
				t.Fatalf("wrong-typed member accepted")
			}
			assertValidationMessage(t, err, tc.frag)
		})
	}

	// A non-string entry inside a modelled string list is reported with
	// its index.
	_, err := parseEcsParameters(map[string]interface{}{
		"TaskDefinitionArn": "arn:aws:ecs:us-east-1:123456789012:task-definition/f:1",
		"NetworkConfiguration": map[string]interface{}{
			"awsvpcConfiguration": map[string]interface{}{
				"Subnets": []interface{}{"subnet-1", 42},
			},
		},
	})
	if err == nil {
		t.Fatal("non-string subnet accepted")
	}
	assertValidationMessage(t, err, "Subnets[1] must be a string")

	// A non-structure entry inside a modelled structure list is reported.
	_, err = parseEcsParameters(map[string]interface{}{
		"TaskDefinitionArn": "arn:aws:ecs:us-east-1:123456789012:task-definition/f:1",
		"PlacementStrategy": []interface{}{"spread"},
	})
	if err == nil {
		t.Fatal("non-map placement strategy item accepted")
	}
	assertValidationMessage(t, err, "PlacementStrategy[0] must be a structure")
}

// The nested ECS shapes bind the names the model declares: Pascal inside
// EcsParameters and AwsVpcConfiguration, the model's own lowerCamel names
// for the awsvpcConfiguration member and the capacity-provider and
// placement items.
func TestParseEcsNestedModelledNames(t *testing.T) {
	ecs, err := parseEcsParameters(map[string]interface{}{
		"TaskDefinitionArn": "arn:aws:ecs:us-east-1:123456789012:task-definition/f:1",
		"LaunchType":        "FARGATE",
		"NetworkConfiguration": map[string]interface{}{
			"awsvpcConfiguration": map[string]interface{}{
				"AssignPublicIp": "ENABLED",
				"Subnets":        []interface{}{"subnet-0abc"},
			},
		},
		"CapacityProviderStrategy": []interface{}{map[string]interface{}{
			"capacityProvider": "FARGATE", "weight": 1.0, "base": 1.0,
		}},
		"PlacementConstraints": []interface{}{map[string]interface{}{
			"type": "memberOf", "expression": "attribute:ecs.ami-id",
		}},
	})
	if err != nil {
		t.Fatalf("modelled ECS parse: %v", err)
	}
	if ecs.LaunchType != "FARGATE" {
		t.Fatalf("LaunchType not bound: %+v", ecs)
	}
	if ecs.NetworkConfiguration == nil || ecs.NetworkConfiguration.AwsVpcConfiguration == nil {
		t.Fatalf("awsvpcConfiguration not bound: %+v", ecs.NetworkConfiguration)
	}
	if ecs.NetworkConfiguration.AwsVpcConfiguration.AssignPublicIp != "ENABLED" {
		t.Fatalf("AssignPublicIp not bound")
	}
	if len(ecs.CapacityProviderStrategy) != 1 || ecs.CapacityProviderStrategy[0].CapacityProvider != "FARGATE" {
		t.Fatalf("capacityProviderStrategy item not bound: %+v", ecs.CapacityProviderStrategy)
	}
	if w := ecs.CapacityProviderStrategy[0].Weight; w == nil || *w != 1 {
		t.Fatalf("weight not bound: %+v", ecs.CapacityProviderStrategy[0])
	}
	if len(ecs.PlacementConstraints) != 1 || ecs.PlacementConstraints[0].Expression != "attribute:ecs.ami-id" {
		t.Fatalf("placement expression not bound: %+v", ecs.PlacementConstraints)
	}

	// The non-model Pascal spellings of the model's lowerCamel members are
	// unknown members and bind nothing.
	ecs, err = parseEcsParameters(map[string]interface{}{
		"TaskDefinitionArn": "arn:aws:ecs:us-east-1:123456789012:task-definition/f:1",
		"NetworkConfiguration": map[string]interface{}{
			"AwsvpcConfiguration": map[string]interface{}{
				"Subnets": []interface{}{"subnet-0abc"},
			},
		},
		"PlacementStrategy": []interface{}{map[string]interface{}{"Field": "host"}},
	})
	if err != nil {
		t.Fatalf("non-model ECS spellings parse: %v", err)
	}
	if ecs.NetworkConfiguration != nil {
		t.Fatalf("AwsvpcConfiguration (non-model) bound: %+v", ecs.NetworkConfiguration)
	}
	if len(ecs.PlacementStrategy) != 1 {
		t.Fatalf("PlacementStrategy item not parsed: %+v", ecs.PlacementStrategy)
	}
	if ecs.PlacementStrategy[0].Field != "" {
		t.Fatalf("Field (non-model) bound: %+v", ecs.PlacementStrategy[0])
	}
}

// The FlexibleTimeWindow string branch binds the same modelled members as
// the map branch — both dialects inside one function was the drift.
func TestParseFlexibleTimeWindowOneDialect(t *testing.T) {
	ftw, err := parseFlexibleTimeWindow(map[string]interface{}{
		"FlexibleTimeWindow": map[string]interface{}{"Mode": "FLEXIBLE", "MaximumWindowInMinutes": 5.0},
	})
	if err != nil || ftw == nil {
		t.Fatalf("map branch parse: %v", err)
	}
	if ftw.Mode != schedulerstore.FlexibleTimeWindowMode("FLEXIBLE") || ftw.MaximumWindowInMinutes == nil || *ftw.MaximumWindowInMinutes != 5 {
		t.Fatalf("map branch not bound: %+v", ftw)
	}

	ftw, err = parseFlexibleTimeWindow(map[string]interface{}{
		"FlexibleTimeWindow": `{"Mode":"OFF"}`,
	})
	if err != nil || ftw == nil {
		t.Fatalf("string branch parse: %v", err)
	}
	if ftw.Mode != schedulerstore.FlexibleTimeWindowMode("OFF") {
		t.Fatalf("string branch not bound: %+v", ftw)
	}

	// The store type's lowerCamel json tags are an internal serialisation
	// detail, not a wire dialect: the string branch no longer accepts them.
	ftw, err = parseFlexibleTimeWindow(map[string]interface{}{
		"FlexibleTimeWindow": `{"mode":"OFF"}`,
	})
	if err != nil {
		t.Fatalf("lowerCamel string parse: %v", err)
	}
	if ftw.Mode != "" {
		t.Fatalf("lowerCamel string dialect bound a Mode: %+v", ftw)
	}

	// The non-model lowercase body key is an unknown member.
	ftw, err = parseFlexibleTimeWindow(map[string]interface{}{
		"flexibleTimeWindow": map[string]interface{}{"Mode": "OFF"},
	})
	if err != nil || ftw != nil {
		t.Fatalf("lowercase body key bound: %+v, %v", ftw, err)
	}
}

// MessageGroupId carries @length(1, 128): an explicitly present empty
// string is out of bounds where presence is visible, while an absent
// member leaves the DTO empty (delivery falls back to the schedule name).
func TestParseSqsParametersPresentEmptyRejected(t *testing.T) {
	sqs, err := parseSqsParameters(map[string]interface{}{})
	if err != nil || sqs.MessageGroupId != "" {
		t.Fatalf("absent member: %+v, %v", sqs, err)
	}
	_, err = parseSqsParameters(map[string]interface{}{"MessageGroupId": ""})
	if err == nil {
		t.Fatal("present-empty MessageGroupId accepted")
	}
	assertValidationMessage(t, err, "SqsParameters.MessageGroupId must be 1-")
	_, err = parseSqsParameters(map[string]interface{}{"MessageGroupId": 42})
	if err == nil {
		t.Fatal("non-string MessageGroupId accepted")
	}
	assertValidationMessage(t, err, "SqsParameters.MessageGroupId must be a string")
}

// A NetworkConfiguration without awsvpcConfiguration is an absent
// structure, not an empty wrapper to persist: restJson1 omits it rather
// than round-tripping it as null.
func TestParseNetworkConfigurationAbsentVpcIsNil(t *testing.T) {
	nc, err := parseNetworkConfiguration(map[string]interface{}{})
	if err != nil || nc != nil {
		t.Fatalf("empty NetworkConfiguration: %+v, %v", nc, err)
	}
	ecs, err := parseEcsParameters(map[string]interface{}{
		"TaskDefinitionArn":    "arn:aws:ecs:us-east-1:123456789012:task-definition/f:1",
		"NetworkConfiguration": map[string]interface{}{},
	})
	if err != nil {
		t.Fatalf("ECS with empty NetworkConfiguration: %v", err)
	}
	if ecs.NetworkConfiguration != nil {
		t.Fatalf("empty NetworkConfiguration wrapper stored: %+v", ecs.NetworkConfiguration)
	}

	// A wrong-typed awsvpcConfiguration is a wire-format violation naming
	// the member, not an absent structure.
	_, err = parseEcsParameters(map[string]interface{}{
		"TaskDefinitionArn":    "arn:aws:ecs:us-east-1:123456789012:task-definition/f:1",
		"NetworkConfiguration": map[string]interface{}{"awsvpcConfiguration": "none"},
	})
	if err == nil {
		t.Fatal("wrong-typed awsvpcConfiguration accepted")
	}
	assertValidationMessage(t, err, "NetworkConfiguration.awsvpcConfiguration must be a structure")
}
