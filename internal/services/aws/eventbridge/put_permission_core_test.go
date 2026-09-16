package eventbridge

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// putPermissionConditionBlock extracts the nested IAM condition element
// {"Type": {"Key": Value}} from the first statement with the given Sid.
func putPermissionConditionBlock(t *testing.T, policy, sid string) map[string]map[string]string {
	t.Helper()
	var doc struct {
		Statement []struct {
			Sid       string                            `json:"Sid"`
			Condition map[string]map[string]interface{} `json:"Condition"`
		} `json:"Statement"`
	}
	if err := json.Unmarshal([]byte(policy), &doc); err != nil {
		t.Fatalf("unmarshal policy: %v", err)
	}
	for _, s := range doc.Statement {
		if s.Sid != sid {
			continue
		}
		out := make(map[string]map[string]string, len(s.Condition))
		for op, kv := range s.Condition {
			row := make(map[string]string, len(kv))
			for k, v := range kv {
				row[k], _ = v.(string)
			}
			out[op] = row
		}
		return out
	}
	return nil
}

// TestPutPermissionConditionStructure pins the typed Condition contract:
// the modelled {Key, Type, Value} structure renders into the resource
// policy as the nested IAM condition element, partial structures are
// rejected, and the documented single supported pair
// (Type=StringEquals, Key=aws:PrincipalOrgID) is enforced.
func TestPutPermissionConditionStructure(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()
	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "pp-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	if err := svc.putPermissionCore(ctx, store, PutPermissionInput{
		BusName: "pp-bus", BusNameProvided: true,
		Principal:   "*",
		StatementId: "OrgStatement",
		Condition:   &PutPermissionCondition{Type: "StringEquals", Key: "aws:PrincipalOrgID", Value: "o-1234567890"},
	}); err != nil {
		t.Fatalf("put permission with condition: %v", err)
	}
	bus, err := store.GetEventBus(ctx, "pp-bus")
	if err != nil {
		t.Fatalf("get bus: %v", err)
	}
	block := putPermissionConditionBlock(t, bus.Policy, "OrgStatement")
	if block == nil {
		t.Fatalf("statement OrgStatement missing from policy: %s", bus.Policy)
	}
	orgID, ok := block["StringEquals"]["aws:PrincipalOrgID"]
	if !ok || orgID != "o-1234567890" {
		t.Fatalf("condition block wrong: %+v", block)
	}

	cases := []struct {
		name      string
		condition *PutPermissionCondition
		wantErr   string
	}{
		{"partial structure", &PutPermissionCondition{Type: "StringEquals", Key: "aws:PrincipalOrgID"}, "Condition must contain Type, Key, and Value"},
		{"unsupported type", &PutPermissionCondition{Type: "ArnLike", Key: "aws:PrincipalOrgID", Value: "o-1"}, "Condition.Type supports only StringEquals"},
		{"unsupported key", &PutPermissionCondition{Type: "StringEquals", Key: "aws:SourceAccount", Value: "123"}, "Condition.Key supports only aws:PrincipalOrgID"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := svc.putPermissionCore(ctx, store, PutPermissionInput{
				BusName: "pp-bus", BusNameProvided: true,
				Principal:   "111122223333",
				StatementId: "Stmt" + strings.ReplaceAll(strings.Title(c.name), " ", ""),
				Condition:   c.condition,
			})
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, c.wantErr)
			}
		})
	}
}

// TestPutPermissionStatementMemberValidation pins the documented statement
// mode member patterns: Principal (\d{12}|\*), StatementId
// [a-zA-Z0-9-_]+ (1-64) and Action events:[a-zA-Z]+ (1-64).
func TestPutPermissionStatementMemberValidation(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()
	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "pp-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	cases := []struct {
		name    string
		input   PutPermissionInput
		wantErr string
	}{
		{"malformed principal", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "not-an-account", StatementId: "S1",
		}, "Principal must be a 12-digit account ID"},
		{"principal too short", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "123", StatementId: "S2",
		}, "Principal must be a 12-digit account ID"},
		{"statement id bad chars", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "111122223333", StatementId: "bad sid!",
		}, "StatementId must match"},
		{"statement id too long", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "111122223333", StatementId: strings.Repeat("s", 65),
		}, "StatementId must match"},
		{"non-events action", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "111122223333", StatementId: "S5", Action: "sqs:SendMessage",
		}, "Action must match"},
		{"action bad chars", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "111122223333", StatementId: "S6", Action: "events:Put Events",
		}, "Action must match"},
		{"action too long", PutPermissionInput{
			BusName: "pp-bus", BusNameProvided: true,
			Principal: "111122223333", StatementId: "S7", Action: "events:" + strings.Repeat("a", 65),
		}, "Action must match"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := svc.putPermissionCore(ctx, store, c.input)
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("got %v, want error containing %q", err, c.wantErr)
			}
		})
	}

	// The documented forms are accepted: the wildcard principal, the
	// account principal and a non-PutEvents events action.
	for i, in := range []PutPermissionInput{
		{BusName: "pp-bus", BusNameProvided: true, Principal: "*", StatementId: "OkWildcard"},
		{BusName: "pp-bus", BusNameProvided: true, Principal: "111122223333", StatementId: "OkAccount"},
		{BusName: "pp-bus", BusNameProvided: true, Principal: "111122223333", StatementId: "OkAction", Action: "events:CustomAction"},
	} {
		if err := svc.putPermissionCore(ctx, store, in); err != nil {
			t.Fatalf("valid case %d rejected: %v", i, err)
		}
	}
}

// TestPutPermissionPolicySizeCeilingBothModes pins that the 10 KB resource
// policy ceiling ("The permission policy on the event bus cannot exceed
// 10 KB in size") is enforced in both modes: on the supplied document in
// Policy mode, and on the merged document in statement mode — where the
// rejection aborts the merge atomically, leaving the stored policy intact.
func TestPutPermissionPolicySizeCeilingBothModes(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()
	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "pp-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	// Policy mode rejects an oversized document up front.
	oversized := `{"Version":"2012-10-17","Statement":[{"Sid":"x","Note":"` + strings.Repeat("a", 11000) + `"}]}`
	err := svc.putPermissionCore(ctx, store, PutPermissionInput{
		BusName: "pp-bus", BusNameProvided: true,
		PolicySet: true, Policy: oversized,
	})
	if err == nil || !strings.Contains(err.Error(), "10240") {
		t.Fatalf("oversized policy document: got %v, want PolicyLengthExceededException", err)
	}

	// An explicitly empty Policy is an invalid policy document, not a
	// silent fall-through into statement mode.
	err = svc.putPermissionCore(ctx, store, PutPermissionInput{
		BusName: "pp-bus", BusNameProvided: true,
		PolicySet: true, Policy: "",
	})
	if err == nil || !strings.Contains(err.Error(), "Invalid policy document") {
		t.Fatalf("empty policy document: got %v, want invalid-document rejection", err)
	}

	// Seed a near-ceiling policy through Policy mode, then let a
	// statement-mode merge push it past the ceiling: the merge aborts and
	// the stored policy is unchanged.
	filler := `{"Version":"2012-10-17","Statement":[{"Sid":"filler","Note":"` + strings.Repeat("a", 10150) + `"}]}`
	if err := svc.putPermissionCore(ctx, store, PutPermissionInput{
		BusName: "pp-bus", BusNameProvided: true,
		PolicySet: true, Policy: filler,
	}); err != nil {
		t.Fatalf("seed policy: %v", err)
	}
	before, err := store.GetEventBus(ctx, "pp-bus")
	if err != nil {
		t.Fatalf("get bus: %v", err)
	}

	err = svc.putPermissionCore(ctx, store, PutPermissionInput{
		BusName: "pp-bus", BusNameProvided: true,
		Principal: "111122223333", StatementId: "OverTheTop",
	})
	if err == nil || !strings.Contains(err.Error(), "10240") {
		t.Fatalf("statement-mode overflow: got %v, want PolicyLengthExceededException", err)
	}
	after, err := store.GetEventBus(ctx, "pp-bus")
	if err != nil {
		t.Fatalf("get bus after rejection: %v", err)
	}
	if after.Policy != before.Policy {
		t.Fatalf("rejected statement merge must not alter the stored policy")
	}
}
