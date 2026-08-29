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
