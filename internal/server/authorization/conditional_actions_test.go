package authorization

import (
	"testing"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/request"
)

// TestConditionalExtraActionsTrigger pins the parameter-conditional
// escalation trigger: only a states TestState request with revealSecrets
// set requires the additional states:RevealSecrets action.
func TestConditionalExtraActionsTrigger(t *testing.T) {
	tests := []struct {
		name      string
		service   string
		operation string
		params    map[string]interface{}
		want      int
	}{
		{"reveal secrets on test state", "states", "TestState", map[string]interface{}{"revealSecrets": true}, 1},
		{"reveal secrets false", "states", "TestState", map[string]interface{}{"revealSecrets": false}, 0},
		{"reveal secrets absent", "states", "TestState", nil, 0},
		{"reveal secrets on another operation", "states", "StartExecution", map[string]interface{}{"revealSecrets": true}, 0},
		{"reveal secrets on another service", "lambda", "Invoke", map[string]interface{}{"revealSecrets": true}, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsed := &request.ParsedRequest{Operation: tc.operation, Parameters: tc.params}
			got := conditionalExtraActions(parsed, tc.service)
			if len(got) != tc.want {
				t.Fatalf("got %v, want %d actions", got, tc.want)
			}
			if tc.want > 0 && got[0] != "states:RevealSecrets" {
				t.Fatalf("action = %q, want states:RevealSecrets", got[0])
			}
		})
	}
}

// TestRevealSecretsPolicyEvaluation pins that the policy evaluator denies
// the states:RevealSecrets action unless a policy statement allows it:
// the access-denied contract TestState documents for revealSecrets.
func TestRevealSecretsPolicyEvaluation(t *testing.T) {
	withPermission := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["states:TestState","states:RevealSecrets"],"Resource":"*"}]}`
	withoutPermission := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["states:TestState"],"Resource":"*"}]}`

	evaluator := policy.NewPolicyEvaluator()

	for _, tc := range []struct {
		name     string
		document string
		allowed  bool
	}{
		{"policy carrying the action allows", withPermission, true},
		{"policy without the action denies", withoutPermission, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := policy.ParseDocument(tc.document)
			if err != nil {
				t.Fatalf("policy parse failed: %v", err)
			}
			evalCtx := &policy.EvaluationContext{
				Principal: "arn:aws:iam::000000000000:user/tester",
				Action:    "states:RevealSecrets",
				Resource:  "*",
			}
			decision := evaluator.Evaluate(evalCtx, []*policy.Document{doc})
			if tc.allowed && decision.Effect != policy.DecisionEffectAllow {
				t.Fatalf("decision = %s, want Allow", decision.Effect)
			}
			if !tc.allowed && decision.Effect == policy.DecisionEffectAllow {
				t.Fatalf("decision = %s, want a denial", decision.Effect)
			}
		})
	}
}
