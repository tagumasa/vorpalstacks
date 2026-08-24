package sfn

import (
	"context"
	"strings"
	"testing"
)

// TestValidateArnRequiredUnicodeLengths pins that the Step Functions Arn
// shape @length(1,256) is counted in Unicode characters; the shape carries
// no pattern, so rune-legal multibyte ARNs must not be rejected on byte
// length.
func TestValidateArnRequiredUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateArnRequired(strings.Repeat(cjk, 256), "stateMachineArn"); err != nil {
		t.Errorf("256-character CJK ARN rejected: %v", err)
	}
	if err := validateArnRequired(strings.Repeat(cjk, 257), "stateMachineArn"); err == nil {
		t.Error("257-character CJK ARN accepted")
	}
	if err := validateArnRequired("", "stateMachineArn"); err == nil {
		t.Error("empty ARN accepted")
	}

	if err := validateRoleArnRequired(strings.Repeat(cjk, 256)); err != nil {
		t.Errorf("256-character CJK role ARN rejected: %v", err)
	}
	if err := validateRoleArnRequired(strings.Repeat(cjk, 257)); err == nil {
		t.Error("257-character CJK role ARN accepted")
	}
	if err := validateRoleArnOptional(strings.Repeat(cjk, 257)); err == nil {
		t.Error("257-character optional role ARN accepted")
	}
	if err := validateRoleArnOptional(""); err != nil {
		t.Errorf("empty optional role ARN rejected: %v", err)
	}
}

func TestValidateResourceNameRejectsDocumentedInvalidRunes(t *testing.T) {
	// The name contract lists the invalid character U+10FFFF separately
	// from the noncharacters U+FFFE-FFFF.
	if err := validateResourceName("bad\U0010FFFFname"); err == nil {
		t.Error("U+10FFFF accepted in a resource name")
	}
	if err := validateResourceName("\U0010FFFF"); err == nil {
		t.Error("U+10FFFF-only resource name accepted")
	}
	if err := validateResourceName("bad\uffffname"); err == nil {
		t.Error("U+FFFF accepted in a resource name")
	}
	if err := validateResourceName("\U0001D11E-legal"); err != nil {
		t.Errorf("musical symbol (a supplementary-plane letter) rejected: %v", err)
	}
}

func TestListExecutionsRejectsRedriveFilterWithStateMachineArn(t *testing.T) {
	svc := &StepFunctionService{}
	store, smArn, _, _ := newAliasTestStore(t)

	_, err := svc.listExecutionsCore(context.Background(), store, ListExecutionsInput{
		StateMachineArn: smArn,
		RedriveFilter:   "REDRIVEN",
	})
	requireAWSCode(t, err, "ValidationException")

	_, err = svc.listExecutionsCore(context.Background(), store, ListExecutionsInput{
		StateMachineArn: smArn,
		RedriveFilter:   "NOT_REDRIVEN",
	})
	requireAWSCode(t, err, "ValidationException")
}

// TestValidateWaitStates pins the Wait-state field contract at definition
// validation: exactly one of the wait fields, strict RFC3339-profile
// timestamps, and integer Seconds within the documented range.
func TestValidateWaitStates(t *testing.T) {
	waitDefinition := func(state string) string {
		// The structural validator requires a terminal state, so every
		// wrapped Wait state carries End unless the case supplies its own
		// transition shape.
		return `{"StartAt":"W","States":{"W":` + strings.TrimSuffix(state, "}") + `,"End":true}}}`
	}

	tests := []struct {
		name       string
		definition string
		wantErr    bool
	}{
		{"seconds zero accepted", waitDefinition(`{"Type":"Wait","Seconds":0}`), false},
		{"seconds ceiling accepted", waitDefinition(`{"Type":"Wait","Seconds":99999999}`), false},
		{"seconds above ceiling rejected", waitDefinition(`{"Type":"Wait","Seconds":100000000}`), true},
		{"negative seconds rejected", waitDefinition(`{"Type":"Wait","Seconds":-1}`), true},
		{"fractional seconds rejected", waitDefinition(`{"Type":"Wait","Seconds":1.5}`), true},
		{"string seconds rejected", waitDefinition(`{"Type":"Wait","Seconds":"10"}`), true},
		{"timestamp accepted", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00Z"}`), false},
		{"timestamp with offset accepted", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00+09:00"}`), false},
		{"timestamp with milliseconds accepted", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00.123Z"}`), false},
		{"offset-less timestamp rejected", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00"}`), true},
		{"lowercase t rejected", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14t01:59:00Z"}`), true},
		{"lowercase z rejected", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00z"}`), true},
		{"two-digit fraction rejected", waitDefinition(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00.12Z"}`), true},
		{"timestamp path accepted", waitDefinition(`{"Type":"Wait","TimestampPath":"$.expiry"}`), false},
		{"seconds path accepted", waitDefinition(`{"Type":"Wait","SecondsPath":"$.delay"}`), false},
		{"empty timestamp path rejected", waitDefinition(`{"Type":"Wait","TimestampPath":""}`), true},
		{"numeric timestamp path rejected", waitDefinition(`{"Type":"Wait","TimestampPath":5}`), true},
		{"no wait field rejected", waitDefinition(`{"Type":"Wait","Next":"W"}`), true},
		{"two wait fields rejected", waitDefinition(`{"Type":"Wait","Seconds":1,"Timestamp":"2024-03-14T01:59:00Z"}`), true},
		{"non-wait state untouched", `{"StartAt":"P","States":{"P":{"Type":"Pass","End":true}}}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDefinitionStructure(tt.definition, "STANDARD")
			if tt.wantErr && err == nil {
				t.Errorf("definition accepted: %s", tt.definition)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("definition rejected: %v", err)
			}
		})
	}
}

// TestValidateWaitStatesJSONata pins the JSONata Wait contract: exactly
// one of Seconds or Timestamp, path fields unsupported, and expressions
// accepted by presence (their values are only checkable at run time).
func TestValidateWaitStatesJSONata(t *testing.T) {
	jsonataWait := func(state string) string {
		return `{"QueryLanguage":"JSONata","StartAt":"W","States":{"W":` + strings.TrimSuffix(state, "}") + `,"End":true}}}`
	}

	tests := []struct {
		name       string
		definition string
		wantErr    bool
	}{
		{"literal seconds accepted", jsonataWait(`{"Type":"Wait","Seconds":10}`), false},
		{"seconds expression accepted", jsonataWait(`{"Type":"Wait","Seconds":"{% $wait %}"}`), false},
		{"literal timestamp accepted", jsonataWait(`{"Type":"Wait","Timestamp":"2024-03-14T01:59:00Z"}`), false},
		{"timestamp expression accepted", jsonataWait(`{"Type":"Wait","Timestamp":"{% $expiry %}"}`), false},
		{"invalid literal timestamp rejected", jsonataWait(`{"Type":"Wait","Timestamp":"2024-03-14 01:59:00"}`), true},
		{"both fields rejected", jsonataWait(`{"Type":"Wait","Seconds":1,"Timestamp":"2024-03-14T01:59:00Z"}`), true},
		{"no field rejected", jsonataWait(`{"Type":"Wait"}`), true},
		{"seconds path rejected", jsonataWait(`{"Type":"Wait","SecondsPath":"$.delay"}`), true},
		{"timestamp path rejected", jsonataWait(`{"Type":"Wait","TimestampPath":"$.expiry"}`), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDefinitionStructure(tt.definition, "STANDARD")
			if tt.wantErr && err == nil {
				t.Errorf("definition accepted: %s", tt.definition)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("definition rejected: %v", err)
			}
		})
	}
}
