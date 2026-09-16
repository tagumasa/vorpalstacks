package eventbridge

import (
	"context"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func newEventBusCoreTestStore(t *testing.T) *eventsstore.EventsStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return eventsstore.NewEventsStore(st, "000000000000", "us-east-1")
}

func newEventBusCoreTestService() *EventsService {
	return NewEventsService(nil, "000000000000")
}

// TestCreateEventBusLogConfigEnumRejected pins the enum contract of the
// LogConfig members on create: every provided value — including the
// explicitly provided empty string, which is not an enum member any more
// than an arbitrary word is — is rejected unless it is NONE/FULL for
// IncludeDetail and OFF/ERROR/INFO/TRACE for Level.
func TestCreateEventBusLogConfigEnumRejected(t *testing.T) {
	cases := []struct {
		name      string
		logConfig *BusLogConfigInput
		wantErr   string
	}{
		{
			name:      "empty IncludeDetail",
			logConfig: &BusLogConfigInput{IncludeDetailSet: true},
			wantErr:   "LogConfig.IncludeDetail must be one of: NONE, FULL",
		},
		{
			name:      "out-of-enum IncludeDetail",
			logConfig: &BusLogConfigInput{IncludeDetailSet: true, IncludeDetail: "BOGUS"},
			wantErr:   "LogConfig.IncludeDetail must be one of: NONE, FULL",
		},
		{
			name:      "empty Level",
			logConfig: &BusLogConfigInput{LevelSet: true},
			wantErr:   "LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE",
		},
		{
			name:      "out-of-enum Level",
			logConfig: &BusLogConfigInput{LevelSet: true, Level: "BOGUS"},
			wantErr:   "LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store := newEventBusCoreTestStore(t)
			svc := newEventBusCoreTestService()
			_, err := svc.createEventBusCore(context.Background(), store, CreateEventBusInput{
				Name:      "enum_bus",
				LogConfig: c.logConfig,
			})
			if err == nil {
				t.Fatalf("expected ValidationException, got nil")
			}
			if !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", c.wantErr, err)
			}
		})
	}

	t.Run("valid values stored", func(t *testing.T) {
		store := newEventBusCoreTestStore(t)
		svc := newEventBusCoreTestService()
		_, err := svc.createEventBusCore(context.Background(), store, CreateEventBusInput{
			Name: "enum_bus",
			LogConfig: &BusLogConfigInput{
				IncludeDetailSet: true,
				IncludeDetail:    "NONE",
				LevelSet:         true,
				Level:            "INFO",
			},
		})
		if err != nil {
			t.Fatalf("create with valid log config: %v", err)
		}
		bus, err := store.GetEventBus(context.Background(), "enum_bus")
		if err != nil {
			t.Fatalf("get event bus: %v", err)
		}
		if bus.LogConfig == nil || bus.LogConfig.IncludeDetail != "NONE" || bus.LogConfig.Level != "INFO" {
			t.Fatalf("stored log config mismatch: %+v", bus.LogConfig)
		}
	})
}

// TestUpdateEventBusLogConfigEnumRejected pins the same enum contract on
// update: the members are validated whenever the LogConfig member is set,
// before any value reaches the store.
func TestUpdateEventBusLogConfigEnumRejected(t *testing.T) {
	cases := []struct {
		name      string
		logConfig *BusLogConfigInput
		wantErr   string
	}{
		{
			name:      "empty IncludeDetail",
			logConfig: &BusLogConfigInput{IncludeDetailSet: true},
			wantErr:   "LogConfig.IncludeDetail must be one of: NONE, FULL",
		},
		{
			name:      "out-of-enum IncludeDetail",
			logConfig: &BusLogConfigInput{IncludeDetailSet: true, IncludeDetail: "BOGUS"},
			wantErr:   "LogConfig.IncludeDetail must be one of: NONE, FULL",
		},
		{
			name:      "empty Level",
			logConfig: &BusLogConfigInput{LevelSet: true},
			wantErr:   "LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE",
		},
		{
			name:      "out-of-enum Level",
			logConfig: &BusLogConfigInput{LevelSet: true, Level: "BOGUS"},
			wantErr:   "LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			store := newEventBusCoreTestStore(t)
			svc := newEventBusCoreTestService()
			ctx := context.Background()
			if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "enum_bus"}); err != nil {
				t.Fatalf("create: %v", err)
			}
			_, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
				Name:         "enum_bus",
				LogConfigSet: true,
				LogConfig:    c.logConfig,
			})
			if err == nil {
				t.Fatalf("expected ValidationException, got nil")
			}
			if !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", c.wantErr, err)
			}
			bus, err := store.GetEventBus(ctx, "enum_bus")
			if err != nil {
				t.Fatalf("get event bus: %v", err)
			}
			if bus.LogConfig != nil {
				t.Fatalf("rejected update must not persist a log config, got %+v", bus.LogConfig)
			}
		})
	}

	t.Run("valid values round-trip", func(t *testing.T) {
		store := newEventBusCoreTestStore(t)
		svc := newEventBusCoreTestService()
		ctx := context.Background()
		if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "enum_bus"}); err != nil {
			t.Fatalf("create: %v", err)
		}
		if _, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
			Name:         "enum_bus",
			LogConfigSet: true,
			LogConfig: &BusLogConfigInput{
				IncludeDetailSet: true,
				IncludeDetail:    "FULL",
				LevelSet:         true,
				Level:            "TRACE",
			},
		}); err != nil {
			t.Fatalf("update with valid log config: %v", err)
		}
		bus, err := store.GetEventBus(ctx, "enum_bus")
		if err != nil {
			t.Fatalf("get event bus: %v", err)
		}
		if bus.LogConfig == nil || bus.LogConfig.IncludeDetail != "FULL" || bus.LogConfig.Level != "TRACE" {
			t.Fatalf("stored log config mismatch: %+v", bus.LogConfig)
		}
	})
}

// TestEventBusCrudValidationSymmetry pins that the same members are
// validated on create and update at the Core layer (no check lives in a
// handler), and that describeEventBusCore owns the Name presence/ARN
// resolution: an absent Name addresses the default bus, an explicitly
// empty one is rejected, and the ARN form resolves to the canonical name.
func TestEventBusCrudValidationSymmetry(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	// The Description bound (Smithy @length(0,512)) is a Core-level check
	// on create, not a handler check.
	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{
		Name:        "sym-bus",
		Description: strings.Repeat("d", 513),
	}); err == nil || !strings.Contains(err.Error(), "Description must be at most") {
		t.Fatalf("over-long description on create: got %v, want rejection", err)
	}

	if _, err := svc.createEventBusCore(ctx, store, CreateEventBusInput{Name: "sym-bus"}); err != nil {
		t.Fatalf("create bus: %v", err)
	}

	// KmsKeyIdentifier is validated on update exactly like create.
	invalidKMS := "arn:aws:events:us-east-1:000000000000:not-a-key"
	if _, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
		Name: "sym-bus", KmsKeyIdentifierSet: true, KmsKeyIdentifier: invalidKMS,
	}); err == nil || !strings.Contains(err.Error(), "KmsKeyIdentifier") {
		t.Fatalf("invalid KmsKeyIdentifier on update: got %v, want rejection", err)
	}
	validKMS := "arn:aws:kms:us-east-1:000000000000:key/12345678-1234-1234-1234-123456789012"
	if _, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
		Name: "sym-bus", KmsKeyIdentifierSet: true, KmsKeyIdentifier: validKMS,
	}); err != nil {
		t.Fatalf("valid KmsKeyIdentifier on update: %v", err)
	}

	// describeEventBusCore: absent Name → default bus (the platform
	// auto-creates it on access paths that write; here the read of a
	// missing default bus is the plain not-found), provided empty Name →
	// rejection, ARN form → canonical name.
	if _, err := svc.describeEventBusCore(ctx, store, "", true); err == nil ||
		!strings.Contains(err.Error(), "EventBusName must not be empty") {
		t.Fatalf("provided empty Name: got %v, want rejection", err)
	}
	res, err := svc.describeEventBusCore(ctx, store, "arn:aws:events:us-east-1:000000000000:event-bus/sym-bus", true)
	if err != nil {
		t.Fatalf("describe via ARN form: %v", err)
	}
	if res.EventBus.Name != "sym-bus" {
		t.Fatalf("ARN form resolved to %q, want sym-bus", res.EventBus.Name)
	}
}

// TestUpdateEventBusNameResolution pins the optional-Name contract on
// UpdateEventBus: the model member carries no required trait and the API
// reference marks it Required: No, so an omitted name updates the default
// bus (the family-wide semantics of optional bus-name members) while an
// explicitly provided empty name is rejected per the EventBusName length
// floor.
func TestUpdateEventBusNameResolution(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()
	ensureDefaultEventBus(store)

	updated, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
		DescriptionSet: true,
		Description:    "updated without naming a bus",
	})
	if err != nil {
		t.Fatalf("update without Name: %v", err)
	}
	if updated.Name != "default" {
		t.Fatalf("absent Name updated %q, want the default bus", updated.Name)
	}
	bus, err := store.GetEventBus(ctx, "default")
	if err != nil {
		t.Fatalf("get default bus: %v", err)
	}
	if bus.Description != "updated without naming a bus" {
		t.Fatalf("default bus description = %q, want the update applied", bus.Description)
	}

	if _, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
		Name:         "",
		NameProvided: true,
	}); err == nil || !strings.Contains(err.Error(), "Event bus name must not be empty") {
		t.Fatalf("provided empty Name: got %v, want ValidationException", err)
	}
}

// TestResolveEventBusNameCore pins the EventBusName presence semantics: an
// absent member addresses the default bus while an explicitly provided
// empty value is rejected.
func TestResolveEventBusNameCore(t *testing.T) {
	name, err := resolveEventBusNameCore("", false)
	if err != nil || name != "default" {
		t.Fatalf("absent member: got (%q, %v), want (\"default\", nil)", name, err)
	}
	if _, err := resolveEventBusNameCore("", true); err == nil ||
		!strings.Contains(err.Error(), "EventBusName must not be empty") {
		t.Fatalf("provided empty member: got %v, want ValidationException", err)
	}
	name, err = resolveEventBusNameCore("bus", true)
	if err != nil || name != "bus" {
		t.Fatalf("provided name: got (%q, %v), want (\"bus\", nil)", name, err)
	}
}

// TestBusNameFormContracts pins the shape contract of every bus-reference
// member that does NOT target EventBusNameOrArn: the plain EventBusName
// members (DeleteEventBus/UpdateEventBus Name) reject the ARN form as a
// pattern violation, and the permission family's NonPartnerEventBusName
// member rejects both the ARN form and the slash-bearing partner names
// before any policy mutation can happen.
func TestBusNameFormContracts(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()
	busARN := "arn:aws:events:us-east-1:000000000000:event-bus/form-bus"

	if err := svc.deleteEventBusCore(ctx, store, DeleteEventBusInput{Name: busARN}); err == nil ||
		!strings.Contains(err.Error(), "Event bus name must match the pattern") {
		t.Fatalf("DeleteEventBus with ARN form: got %v, want the pattern ValidationException", err)
	}
	if _, err := svc.updateEventBusCore(ctx, store, UpdateEventBusInput{
		Name: busARN, NameProvided: true, DescriptionSet: true, Description: "d",
	}); err == nil || !strings.Contains(err.Error(), "Event bus name must match the pattern") {
		t.Fatalf("UpdateEventBus with ARN form: got %v, want the pattern ValidationException", err)
	}

	if _, err := resolveNonPartnerEventBusNameCore(busARN, true); err == nil ||
		!strings.Contains(err.Error(), "Event bus name must match the pattern") {
		t.Fatalf("permission resolver with ARN form: got %v, want the pattern ValidationException", err)
	}
	if _, err := resolveNonPartnerEventBusNameCore("aws.partner/example", true); err == nil {
		t.Fatalf("permission resolver with a partner-style name must be rejected")
	}
	name, err := resolveNonPartnerEventBusNameCore("form-bus", true)
	if err != nil || name != "form-bus" {
		t.Fatalf("permission resolver with a bare name: got (%q, %v), want (\"form-bus\", nil)", name, err)
	}
	if name, err = resolveNonPartnerEventBusNameCore("", false); err != nil || name != "default" {
		t.Fatalf("permission resolver with an absent member: got (%q, %v), want (\"default\", nil)", name, err)
	}

	// Through the real core paths: neither permission operation may touch
	// the store when the bus reference is an ARN form.
	if err := svc.putPermissionCore(ctx, store, PutPermissionInput{
		BusName: busARN, BusNameProvided: true, Principal: "*", StatementId: "s1", Action: "events:PutEvents",
	}); err == nil || !strings.Contains(err.Error(), "Event bus name must match the pattern") {
		t.Fatalf("PutPermission with ARN form: got %v, want the pattern ValidationException", err)
	}
	if err := svc.removePermissionCore(ctx, store, RemovePermissionInput{
		BusName: busARN, BusNameProvided: true, StatementId: "s1",
	}); err == nil || !strings.Contains(err.Error(), "Event bus name must match the pattern") {
		t.Fatalf("RemovePermission with ARN form: got %v, want the pattern ValidationException", err)
	}
}

// TestRuleTargetValidationPrecedence pins the validation order of the
// operations that carry both a primary required member and the optional
// EventBusName: the Smithy model declares the primary member first
// (PutTargetsRequest lists Rule ahead of EventBusName, PutRuleRequest
// lists Name first, ListRuleNamesByTargetRequest lists TargetArn first),
// so a doubly-invalid request must fail on the primary member. The
// operations without a primary required member resolve the bus name first.
func TestRuleTargetValidationPrecedence(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := context.Background()

	cases := []struct {
		name    string
		call    func() error
		wantMsg string
	}{
		{"PutRule", func() error {
			_, err := svc.putRuleCore(ctx, store, PutRuleInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "Rule name is required"},
		{"DeleteRule", func() error {
			return svc.deleteRuleCore(ctx, store, DeleteRuleInput{EventBusName: "", EventBusNameProvided: true})
		}, "Rule name is required"},
		{"DescribeRule", func() error {
			_, err := svc.describeRuleCore(ctx, store, DescribeRuleInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "Rule name is required"},
		{"EnableRule", func() error {
			return svc.setRuleStateCore(ctx, store, SetRuleStateInput{EventBusName: "", EventBusNameProvided: true, State: eventsstore.RuleStateEnabled})
		}, "Rule name is required"},
		{"PutTargets", func() error {
			_, err := svc.putTargetsCore(ctx, store, PutTargetsInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "Rule name is required"},
		{"RemoveTargets", func() error {
			_, err := svc.removeTargetsCore(ctx, store, RemoveTargetsInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "Rule name is required"},
		{"ListTargetsByRule", func() error {
			_, err := svc.listTargetsByRuleCore(ctx, store, ListTargetsByRuleInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "Rule name is required"},
		{"ListRuleNamesByTarget", func() error {
			_, err := svc.listRuleNamesByTargetCore(ctx, store, ListRuleNamesByTargetInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "TargetArn is required"},
		{"ListRules", func() error {
			_, err := svc.listRulesCore(ctx, store, ListRulesInput{EventBusName: "", EventBusNameProvided: true})
			return err
		}, "EventBusName must not be empty"},
		{"PutPermission", func() error {
			return svc.putPermissionCore(ctx, store, PutPermissionInput{BusName: "", BusNameProvided: true})
		}, "EventBusName must not be empty"},
		{"RemovePermission", func() error {
			return svc.removePermissionCore(ctx, store, RemovePermissionInput{BusName: "", BusNameProvided: true})
		}, "EventBusName must not be empty"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.call()
			if err == nil {
				t.Fatalf("expected ValidationException, got nil")
			}
			if !strings.Contains(err.Error(), c.wantMsg) {
				t.Fatalf("expected error containing %q, got: %v", c.wantMsg, err)
			}
		})
	}
}
