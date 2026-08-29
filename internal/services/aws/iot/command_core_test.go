package iot

import (
	"errors"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// TestParseTimeFilterLayouts pins the accepted TimeFilter date-time
// layouts: the documented yyyy-MM-dd'T'HH:mm format and the RFC 3339 form,
// with anything else rejected.
func TestParseTimeFilterLayouts(t *testing.T) {
	cases := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"documented minute format", map[string]interface{}{"after": "2026-01-15T10:30"}, false},
		{"rfc3339", map[string]interface{}{"after": "2026-01-15T10:30:00Z"}, false},
		{"garbage", map[string]interface{}{"after": "not-a-timestamp"}, true},
		{"non-string", map[string]interface{}{"after": 42}, true},
		{"absent", nil, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := parseTimeFilter(tc.value)
			if tc.wantErr && !errors.Is(err, iotstore.ErrValidation) {
				t.Fatalf("expected ErrValidation, got %v", err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected acceptance, got %v", err)
			}
		})
	}
}

// TestListCommandExecutionsFilterExclusivity pins the documented exclusivity
// rules: providing both time filters, or both the command and target ARN,
// is rejected before any store access.
func TestListCommandExecutionsFilterExclusivity(t *testing.T) {
	svc := &IoTService{}
	if _, err := svc.listCommandExecutionsCore(nil, ListCommandExecutionsInput{
		StartedTimeFilter:   map[string]interface{}{"after": "2026-01-01T00:00"},
		CompletedTimeFilter: map[string]interface{}{"before": "2030-01-01T00:00"},
	}); !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for both time filters, got %v", err)
	}
	if _, err := svc.listCommandExecutionsCore(nil, ListCommandExecutionsInput{
		TargetArn:  "arn:aws:iot:us-east-1:000000000000:thing/dev",
		CommandArn: "arn:aws:iot:us-east-1:000000000000:command/reboot",
	}); !errors.Is(err, iotstore.ErrValidation) {
		t.Fatalf("expected ErrValidation for both ARNs, got %v", err)
	}
}

func newCommandCoreTestStore(t *testing.T) iotstore.IotStoreInterface {
	t.Helper()
	rawStore, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { rawStore.Close() })
	return iotstore.NewIotStore(rawStore, "000000000000", "us-east-1", nil)
}

// TestListCommandsDefaultDescendingOrder pins the documented default: with
// sortOrder omitted, commands are listed in descending creation order.
func TestListCommandsDefaultDescendingOrder(t *testing.T) {
	store := newCommandCoreTestStore(t)
	svc := &IoTService{}
	if _, err := svc.createCommandCore(store, CreateCommandInput{CommandID: "cmd-early"}); err != nil {
		t.Fatalf("create early: %v", err)
	}
	// createdAt carries second granularity, so keep the two records in
	// distinct seconds before comparing their default order.
	time.Sleep(1100 * time.Millisecond)
	if _, err := svc.createCommandCore(store, CreateCommandInput{CommandID: "cmd-late"}); err != nil {
		t.Fatalf("create late: %v", err)
	}
	commands, err := svc.listCommandsCore(store, ListCommandsInput{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(commands) != 2 {
		t.Fatalf("expected 2 commands, got %d", len(commands))
	}
	if commands[0]["commandId"] != "cmd-late" || commands[1]["commandId"] != "cmd-early" {
		t.Fatalf("expected the later-created command first, got %v then %v",
			commands[0]["commandId"], commands[1]["commandId"])
	}
}

// TestListCommandExecutionsDefaultOrderAndSortKey pins the documented
// default ordering: executions are listed in descending order keyed on the
// start time, or on the completion time when the completed filter carries
// the window. Command executions cannot be created through any SDK
// operation on this platform, so the records are seeded directly.
func TestListCommandExecutionsDefaultOrderAndSortKey(t *testing.T) {
	store := newCommandCoreTestStore(t)
	svc := &IoTService{}
	seed := func(id string, startedAt, completedAt int64) {
		rec := map[string]interface{}{
			"executionId": id,
			"commandArn":  "arn:aws:iot:us-east-1:000000000000:command/reboot",
			"targetArn":   "arn:aws:iot:us-east-1:000000000000:thing/dev",
			"status":      "SUCCEEDED",
			"createdAt":   startedAt,
			"startedAt":   startedAt,
			"completedAt": completedAt,
		}
		if err := store.PutGeneric("iot-command-execution/"+id, rec); err != nil {
			t.Fatalf("seed %s: %v", id, err)
		}
	}
	seed("exec-a", 100, 500)
	seed("exec-b", 300, 200)
	seed("exec-c", 200, 400)

	// No filter: descending by the start time (b=300, c=200, a=100).
	got, err := svc.listCommandExecutionsCore(store, ListCommandExecutionsInput{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if order := executionIDs(got); order != "exec-b,exec-c,exec-a" {
		t.Fatalf("expected descending start-time order, got %s", order)
	}

	// Started filter with an explicit ascending order: a=100, c=200, b=300.
	got, err = svc.listCommandExecutionsCore(store, ListCommandExecutionsInput{
		SortOrder:         "ASCENDING",
		StartedTimeFilter: map[string]interface{}{"after": "1970-01-01T00:00"},
	})
	if err != nil {
		t.Fatalf("list ascending: %v", err)
	}
	if order := executionIDs(got); order != "exec-a,exec-c,exec-b" {
		t.Fatalf("expected ascending start-time order, got %s", order)
	}

	// Completed filter carries the window, so the completion time keys the
	// default descending order (a=500, c=400, b=200) — a different
	// permutation than the start-time order above.
	got, err = svc.listCommandExecutionsCore(store, ListCommandExecutionsInput{
		CompletedTimeFilter: map[string]interface{}{"after": "1970-01-01T00:00"},
	})
	if err != nil {
		t.Fatalf("list by completion: %v", err)
	}
	if order := executionIDs(got); order != "exec-a,exec-c,exec-b" {
		t.Fatalf("expected descending completion-time order, got %s", order)
	}
}

func executionIDs(records []map[string]interface{}) string {
	ids := ""
	for i, rec := range records {
		if i > 0 {
			ids += ","
		}
		id, _ := rec["executionId"].(string)
		ids += id
	}
	return ids
}
