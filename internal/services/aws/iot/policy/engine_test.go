package policy

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEvaluateConditionStringEquals(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Connect"},
			Condition: map[string]map[string]string{
				"iot:ClientId": {"StringEquals": "trusted-device"},
			},
		}},
	}

	// Matching client — allowed.
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Connect",
		Resource: "anything",
		ClientID: "trusted-device",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected StringEquals match to allow")
	}

	// Non-matching client — denied by condition.
	req.ClientID = "untrusted"
	allowed, err = Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected StringEquals mismatch to deny")
	}
}

func TestEvaluateConditionStringNotEquals(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Publish"},
			Condition: map[string]map[string]string{
				"iot:ClientId": {"StringNotEquals": "banned-device"},
			},
		}},
	}

	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Publish",
		Resource: "topic/a",
		ClientID: "normal-device",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected StringNotEquals to allow non-banned client")
	}

	req.ClientID = "banned-device"
	allowed, err = Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected StringNotEquals to deny banned client")
	}
}

func TestEvaluateConditionStringLike(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Connect"},
			Condition: map[string]map[string]string{
				"iot:ClientId": {"StringLike": "sensor-*"},
			},
		}},
	}

	tests := []struct {
		clientID string
		want     bool
	}{
		{"sensor-1", true},
		{"sensor-abc", true},
		{"actuator-1", false},
		{"sensor", false}, // wildcard requires at least "sensor-" prefix plus char
	}

	for _, tt := range tests {
		req := &EvaluateRequest{
			Policies: []*PolicyVersion{pol},
			Action:   "iot:Connect",
			Resource: "x",
			ClientID: tt.clientID,
		}
		allowed, err := Evaluate(req)
		if err != nil {
			t.Fatalf("unexpected error for clientID=%s: %v", tt.clientID, err)
		}
		if allowed != tt.want {
			t.Errorf("StringLike clientID=%s: got %v, want %v", tt.clientID, allowed, tt.want)
		}
	}
}

func TestEvaluateConditionBool(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Connect"},
			Condition: map[string]map[string]string{
				"iot:SourceIp": {"Bool": "true"},
			},
		}},
	}

	// "true" == "true" → match.
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Connect",
		Resource: "x",
		SourceIP: "true",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected Bool match")
	}

	// "false" != "true" → no match.
	req.SourceIP = "false"
	allowed, err = Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected Bool mismatch to deny")
	}
}

func TestEvaluateConditionUnsupportedOperator(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Connect"},
			Condition: map[string]map[string]string{
				"iot:ClientId": {"TotallyBogusOperator": "42"},
			},
		}},
	}
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Connect",
		Resource: "x",
		ClientID: "42",
	}
	_, err := Evaluate(req)
	if err == nil {
		t.Fatal("expected error for unsupported condition operator")
	}
	if !strings.Contains(err.Error(), "unsupported condition operator") {
		t.Fatalf("expected 'unsupported condition operator' in error, got: %v", err)
	}
}

func TestEvaluateMultipleStatementsAllowAny(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{
			{
				Effect:   EffectAllow,
				Action:   StringOrSlice{"iot:Connect"},
				Resource: StringOrSlice{"client-a"},
			},
			{
				Effect:   EffectAllow,
				Action:   StringOrSlice{"iot:Connect"},
				Resource: StringOrSlice{"client-b"},
			},
		},
	}

	for _, res := range []string{"client-a", "client-b"} {
		req := &EvaluateRequest{
			Policies: []*PolicyVersion{pol},
			Action:   "iot:Connect",
			Resource: res,
		}
		allowed, err := Evaluate(req)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", res, err)
		}
		if !allowed {
			t.Errorf("expected allow for %s", res)
		}
	}

	// Neither statement matches.
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Connect",
		Resource: "client-c",
	}
	allowed, _ := Evaluate(req)
	if allowed {
		t.Fatal("expected deny for unmatched resource")
	}
}

func TestEvaluateMultiplePolicies(t *testing.T) {
	pol1 := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   EffectAllow,
			Action:   StringOrSlice{"iot:Publish"},
			Resource: StringOrSlice{"topic/sensors/#"},
		}},
	}
	pol2 := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   EffectAllow,
			Action:   StringOrSlice{"iot:Subscribe"},
			Resource: StringOrSlice{"topic/commands/#"},
		}},
	}

	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol1, pol2},
		Action:   "iot:Publish",
		Resource: "topic/sensors/temp",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected pol1 to allow publish")
	}

	req.Action = "iot:Subscribe"
	req.Resource = "topic/commands/reboot"
	allowed, err = Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected pol2 to allow subscribe")
	}

	req.Action = "iot:Publish"
	req.Resource = "topic/commands/reboot"
	allowed, _ = Evaluate(req)
	if allowed {
		t.Fatal("expected deny for publish to commands (no matching policy)")
	}
}

func TestEvaluateEmptyPolicies(t *testing.T) {
	req := &EvaluateRequest{
		Policies: nil,
		Action:   "iot:Connect",
		Resource: "*",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected deny with no policies")
	}
}

func TestEvaluateDenyInSecondPolicy(t *testing.T) {
	pol1 := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   EffectAllow,
			Action:   StringOrSlice{"iot:*"},
			Resource: StringOrSlice{"*"},
		}},
	}
	pol2 := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   EffectDeny,
			Action:   StringOrSlice{"iot:Publish"},
			Resource: StringOrSlice{"topic/secret/*"},
		}},
	}
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol1, pol2},
		Action:   "iot:Publish",
		Resource: "topic/secret/data",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected cross-policy Deny to override Allow")
	}
}

func TestParsePolicyVersionValid(t *testing.T) {
	doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:*","Resource":"*"}]}`
	pv, err := ParsePolicyVersion([]byte(doc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pv.Version != "2012-10-17" {
		t.Errorf("version = %s, want 2012-10-17", pv.Version)
	}
	if len(pv.Statement) != 1 {
		t.Fatalf("statements = %d, want 1", len(pv.Statement))
	}
	if pv.Statement[0].Effect != EffectAllow {
		t.Errorf("effect = %s, want Allow", pv.Statement[0].Effect)
	}
}

func TestParsePolicyVersionMissingVersion(t *testing.T) {
	doc := `{"Statement":[{"Effect":"Allow","Action":"iot:*","Resource":"*"}]}`
	_, err := ParsePolicyVersion([]byte(doc))
	if err == nil {
		t.Fatal("expected error for missing Version")
	}
}

func TestParsePolicyVersionMissingEffect(t *testing.T) {
	doc := `{"Version":"2012-10-17","Statement":[{"Action":"iot:*","Resource":"*"}]}`
	_, err := ParsePolicyVersion([]byte(doc))
	if err == nil {
		t.Fatal("expected error for missing Effect")
	}
}

func TestParsePolicyVersionInvalidJSON(t *testing.T) {
	_, err := ParsePolicyVersion([]byte(`{invalid json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParsePolicyVersionArrayAction(t *testing.T) {
	doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["iot:Connect","iot:Publish"],"Resource":"*"}]}`
	pv, err := ParsePolicyVersion([]byte(doc))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pv.Statement[0].Action) != 2 {
		t.Errorf("action count = %d, want 2", len(pv.Statement[0].Action))
	}
}

func TestStringOrSliceUnmarshalSingle(t *testing.T) {
	var s StringOrSlice
	if err := json.Unmarshal([]byte(`"iot:Connect"`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s) != 1 || s[0] != "iot:Connect" {
		t.Errorf("got %v, want [iot:Connect]", []string(s))
	}
}

func TestStringOrSliceUnmarshalArray(t *testing.T) {
	var s StringOrSlice
	if err := json.Unmarshal([]byte(`["iot:Connect","iot:Publish"]`), &s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s) != 2 {
		t.Errorf("len = %d, want 2", len(s))
	}
}

func TestStringOrSliceUnmarshalInvalid(t *testing.T) {
	var s StringOrSlice
	err := json.Unmarshal([]byte(`123`), &s)
	if err == nil {
		t.Fatal("expected error for numeric value")
	}
}

func TestEvaluateResourceARNMatching(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   EffectAllow,
			Action:   StringOrSlice{"iot:Publish"},
			Resource: StringOrSlice{"arn:aws:iot:us-east-1:123456789012:topic/sensors/#"},
		}},
	}

	tests := []struct {
		resource string
		want     bool
	}{
		{"arn:aws:iot:us-east-1:123456789012:topic/sensors/temp", true},
		{"arn:aws:iot:us-east-1:123456789012:topic/sensors/", true},
		{"arn:aws:iot:us-east-1:123456789012:topic/commands/reboot", false},
		{"sensors/temp", false},
	}

	for _, tt := range tests {
		req := &EvaluateRequest{
			Policies: []*PolicyVersion{pol},
			Action:   "iot:Publish",
			Resource: tt.resource,
		}
		allowed, err := Evaluate(req)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", tt.resource, err)
		}
		if allowed != tt.want {
			t.Errorf("resource=%s: got %v, want %v", tt.resource, allowed, tt.want)
		}
	}
}

func TestEvaluateActionWildcardMatching(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:*"},
		}},
	}

	actions := []string{"iot:Connect", "iot:Publish", "iot:Subscribe", "iot:Receive"}
	for _, action := range actions {
		req := &EvaluateRequest{
			Policies: []*PolicyVersion{pol},
			Action:   action,
			Resource: "anything",
		}
		allowed, err := Evaluate(req)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", action, err)
		}
		if !allowed {
			t.Errorf("expected iot:* to allow %s", action)
		}
	}

	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "s3:GetObject",
		Resource: "anything",
	}
	allowed, _ := Evaluate(req)
	if allowed {
		t.Fatal("expected iot:* to deny s3:GetObject")
	}
}

func TestEvaluateMultipleActionPatterns(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Connect", "iot:Publish"},
		}},
	}

	for _, action := range []string{"iot:Connect", "iot:Publish"} {
		req := &EvaluateRequest{
			Policies: []*PolicyVersion{pol},
			Action:   action,
			Resource: "x",
		}
		allowed, _ := Evaluate(req)
		if !allowed {
			t.Errorf("expected allow for %s", action)
		}
	}

	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Subscribe",
		Resource: "x",
	}
	allowed, _ := Evaluate(req)
	if allowed {
		t.Error("expected deny for iot:Subscribe (not in action list)")
	}
}

func TestSubstituteResourcesAllVariables(t *testing.T) {
	resources := []string{
		"client/${iot:ClientId}",
		"topic/${iot:topic}",
		"from/${iot:SourceIp}",
	}
	req := &EvaluateRequest{
		ClientID: "dev-001",
		Topic:    "sensors/temp",
		SourceIP: "10.0.0.1",
	}
	result := substituteResources(resources, req)

	expected := []string{"client/dev-001", "topic/sensors/temp", "from/10.0.0.1"}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("substituteResources[%d] = %s, want %s", i, result[i], want)
		}
	}
}

func TestSubstituteResourcesEmpty(t *testing.T) {
	result := substituteResources(nil, &EvaluateRequest{})
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestEvaluateNoActionSpecified(t *testing.T) {
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect:   EffectAllow,
			Resource: StringOrSlice{"#"},
		}},
	}
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "anything",
		Resource: "anything",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected allow when Action is omitted (matches all)")
	}
}

func TestEvaluateNoResourceSpecified(t *testing.T) {
	// AWS: If Resource is omitted, the statement matches all resources.
	pol := &PolicyVersion{
		Version: "2012-10-17",
		Statement: []PolicyStatement{{
			Effect: EffectAllow,
			Action: StringOrSlice{"iot:Connect"},
		}},
	}
	req := &EvaluateRequest{
		Policies: []*PolicyVersion{pol},
		Action:   "iot:Connect",
		Resource: "anything",
	}
	allowed, err := Evaluate(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected allow when Resource is omitted (matches all)")
	}
}
